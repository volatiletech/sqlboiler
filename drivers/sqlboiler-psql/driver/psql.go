// Package driver implements an sqlboiler driver.
// It can be used by either building the main.go in the same project
// and using as a binary or using the side effect import.
package driver

import (
	"database/sql"
	"embed"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/importers"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/strmangle"

	// Side-effect import sql driver
	_ "github.com/lib/pq"

	_ "embed"
)

//go:embed override
var templates embed.FS

func init() {
	drivers.RegisterFromInit("psql", &PostgresDriver{})
}

// Assemble is more useful for calling into the library so you don't
// have to instantiate an empty type.
func Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	driver := PostgresDriver{}
	return driver.Assemble(config)
}

// PostgresDriver holds the database connection string and a handle
// to the database connection.
type PostgresDriver struct {
	connStr string
	conn    *sql.DB
	version int
}

// Templates that should be added/overridden
func (p PostgresDriver) Templates() (map[string]string, error) {
	tpls := make(map[string]string)
	fs.WalkDir(templates, "override", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(templates, path)
		if err != nil {
			return err
		}
		tpls[strings.Replace(path, "override/", "", 1)] = base64.StdEncoding.EncodeToString(b)

		return nil
	})

	return tpls, nil
}

// Assemble all the information we need to provide back to the driver
func (p *PostgresDriver) Assemble(config drivers.Config) (dbinfo *drivers.DBInfo, err error) {
	defer func() {
		if r := recover(); r != nil && err == nil {
			dbinfo = nil
			err = r.(error)
		}
	}()

	user := config.MustString(drivers.ConfigUser)
	pass, _ := config.String(drivers.ConfigPass)
	dbname := config.MustString(drivers.ConfigDBName)
	host := config.MustString(drivers.ConfigHost)
	port := config.DefaultInt(drivers.ConfigPort, 5432)
	sslmode := config.DefaultString(drivers.ConfigSSLMode, "require")
	schema := config.DefaultString(drivers.ConfigSchema, "public")
	whitelist, _ := config.StringSlice(drivers.ConfigWhitelist)
	blacklist, _ := config.StringSlice(drivers.ConfigBlacklist)

	useSchema := schema != "public"

	p.connStr = PSQLBuildQueryString(user, pass, dbname, host, port, sslmode)
	p.conn, err = sql.Open("postgres", p.connStr)
	if err != nil {
		return nil, errors.Wrap(err, "sqlboiler-psql failed to connect to database")
	}

	defer func() {
		if e := p.conn.Close(); e != nil {
			dbinfo = nil
			err = e
		}
	}()

	p.version, err = p.getVersion()
	if err != nil {
		return nil, errors.Wrap(err, "sqlboiler-psql failed to get database version")
	}

	dbinfo = &drivers.DBInfo{
		Schema: schema,
		Dialect: drivers.Dialect{
			LQ: '"',
			RQ: '"',

			UseIndexPlaceholders: true,
			UseSchema:            useSchema,
			UseDefaultKeyword:    true,
		},
	}
	dbinfo.Tables, err = drivers.Tables(p, schema, whitelist, blacklist)
	if err != nil {
		return nil, err
	}

	return dbinfo, err
}

// PSQLBuildQueryString builds a query string.
func PSQLBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	parts := []string{}
	if len(user) != 0 {
		parts = append(parts, fmt.Sprintf("user=%s", user))
	}
	if len(pass) != 0 {
		parts = append(parts, fmt.Sprintf("password=%s", pass))
	}
	if len(dbname) != 0 {
		parts = append(parts, fmt.Sprintf("dbname=%s", dbname))
	}
	if len(host) != 0 {
		parts = append(parts, fmt.Sprintf("host=%s", host))
	}
	if port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", port))
	}
	if len(sslmode) != 0 {
		parts = append(parts, fmt.Sprintf("sslmode=%s", sslmode))
	}

	return strings.Join(parts, " ")
}

// TableNames connects to the postgres database and
// retrieves all table names from the information_schema where the
// table schema is schema. It uses a whitelist and blacklist.
func (p *PostgresDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var names []string

	query := fmt.Sprintf(`select table_name from information_schema.tables where table_schema = $1 and table_type = 'BASE TABLE'`)
	args := []interface{}{schema}
	if len(whitelist) > 0 {
		tables := drivers.TablesFromList(whitelist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" and table_name in (%s)", strmangle.Placeholders(true, len(tables), 2, 1))
			for _, w := range tables {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		tables := drivers.TablesFromList(blacklist)
		if len(tables) > 0 {
			query += fmt.Sprintf(" and table_name not in (%s)", strmangle.Placeholders(true, len(tables), 2, 1))
			for _, b := range tables {
				args = append(args, b)
			}
		}
	}

	query += ` order by table_name;`

	rows, err := p.conn.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil
}

// Columns takes a table name and attempts to retrieve the table information
// from the database information_schema.columns. It retrieves the column names
// and column types and returns those as a []Column after TranslateColumnType()
// converts the SQL types to Go types, for example: "varchar" to "string"
func (p *PostgresDriver) Columns(schema, tableName string, whitelist, blacklist []string) ([]drivers.Column, error) {
	var columns []drivers.Column
	args := []interface{}{schema, tableName}

	query := `
	select
		c.column_name,
		ct.column_type,
		(
			case when c.character_maximum_length != 0
			then
			(
				ct.column_type || '(' || c.character_maximum_length || ')'
			)
			else c.udt_name
			end
		) as column_full_type,

		c.udt_name,
		e.data_type as array_type,
		c.domain_name,
		c.column_default,

		COALESCE(col_description(('"'||c.table_schema||'"."'||c.table_name||'"')::regclass::oid, ordinal_position), '') as column_comment,

		c.is_nullable = 'YES' as is_nullable,
		(case
			when (select
		    case
			    when column_name = 'is_identity' then (select c.is_identity = 'YES' as is_identity)
		    else
			    false
		    end as is_identity from information_schema.columns
		    WHERE table_schema='information_schema' and table_name='columns' and column_name='is_identity') IS NULL then 'NO' else is_identity end) = 'YES' as is_identity,
		(select exists(
			select 1
			from information_schema.table_constraints tc
			inner join information_schema.constraint_column_usage as ccu on tc.constraint_name = ccu.constraint_name
			where tc.table_schema = $1 and tc.constraint_type = 'UNIQUE' and ccu.constraint_schema = $1 and ccu.table_name = c.table_name and ccu.column_name = c.column_name and
				(select count(*) from information_schema.constraint_column_usage where constraint_schema = $1 and constraint_name = tc.constraint_name) = 1
		)) OR
		(select exists(
			select 1
			from pg_indexes pgix
			inner join pg_class pgc on pgix.indexname = pgc.relname and pgc.relkind = 'i' and pgc.relnatts = 1
			inner join pg_index pgi on pgi.indexrelid = pgc.oid
			inner join pg_attribute pga on pga.attrelid = pgi.indrelid and pga.attnum = ANY(pgi.indkey)
			where
				pgix.schemaname = $1 and pgix.tablename = c.table_name and pga.attname = c.column_name and pgi.indisunique = true
		)) as is_unique

		from information_schema.columns as c
		inner join pg_namespace as pgn on pgn.nspname = c.udt_schema
		left join pg_type pgt on c.data_type = 'USER-DEFINED' and pgn.oid = pgt.typnamespace and c.udt_name = pgt.typname
		left join information_schema.element_types e
			on ((c.table_catalog, c.table_schema, c.table_name, 'TABLE', c.dtd_identifier)
			= (e.object_catalog, e.object_schema, e.object_name, e.object_type, e.collection_type_identifier)),
		lateral (select
			(
				case when pgt.typtype = 'e'
				then
				(
					select 'enum.' || c.udt_name || '(''' || string_agg(labels.label, ''',''') || ''')'
					from (
						select pg_enum.enumlabel as label
						from pg_enum
						where pg_enum.enumtypid =
						(
							select typelem
							from pg_type
							inner join pg_namespace ON pg_type.typnamespace = pg_namespace.oid
							where pg_type.typtype = 'b' and pg_type.typname = ('_' || c.udt_name) and pg_namespace.nspname=$1
							limit 1
						)
						order by pg_enum.enumsortorder
					) as labels
				)
				else c.data_type
				end
			) as column_type
		) ct
		where c.table_name = $2 and c.table_schema = $1 and c.is_generated = 'NEVER'`

	if len(whitelist) > 0 {
		cols := drivers.ColumnsFromList(whitelist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name in (%s)", strmangle.Placeholders(true, len(cols), 3, 1))
			for _, w := range cols {
				args = append(args, w)
			}
		}
	} else if len(blacklist) > 0 {
		cols := drivers.ColumnsFromList(blacklist, tableName)
		if len(cols) > 0 {
			query += fmt.Sprintf(" and c.column_name not in (%s)", strmangle.Placeholders(true, len(cols), 3, 1))
			for _, w := range cols {
				args = append(args, w)
			}
		}
	}

	query += ` order by c.ordinal_position;`

	rows, err := p.conn.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, colFullType, udtName, comment string
		var defaultValue, arrayType, domainName *string
		var nullable, identity, unique bool
		if err := rows.Scan(&colName, &colType, &colFullType, &udtName, &arrayType, &domainName, &defaultValue, &comment, &nullable, &identity, &unique); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		column := drivers.Column{
			Name:       colName,
			DBType:     colType,
			FullDBType: colFullType,
			ArrType:    arrayType,
			DomainName: domainName,
			UDTName:    udtName,
			Comment:    comment,
			Nullable:   nullable,
			Unique:     unique,
		}
		if defaultValue != nil {
			column.Default = *defaultValue
		}

		if identity != false {
			column.Default = "IDENTITY"
		}

		columns = append(columns, column)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (p *PostgresDriver) PrimaryKeyInfo(schema, tableName string) (*drivers.PrimaryKey, error) {
	pkey := &drivers.PrimaryKey{}
	var err error

	query := `
	select tc.constraint_name
	from information_schema.table_constraints as tc
	where tc.table_name = $1 and tc.constraint_type = 'PRIMARY KEY' and tc.table_schema = $2;`

	row := p.conn.QueryRow(query, tableName, schema)
	if err = row.Scan(&pkey.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryColumns := `
	select kcu.column_name
	from   information_schema.key_column_usage as kcu
	where  constraint_name = $1 and table_name = $2 and table_schema = $3
	order by kcu.ordinal_position;`

	var rows *sql.Rows
	if rows, err = p.conn.Query(queryColumns, pkey.Name, tableName, schema); err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var column string

		err = rows.Scan(&column)
		if err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	pkey.Columns = columns

	return pkey, nil
}

// ForeignKeyInfo retrieves the foreign keys for a given table name.
func (p *PostgresDriver) ForeignKeyInfo(schema, tableName string) ([]drivers.ForeignKey, error) {
	var fkeys []drivers.ForeignKey

	whereConditions := []string{"pgn.nspname = $2", "pgc.relname = $1", "pgcon.contype = 'f'"}
	if p.version >= 120000 {
		whereConditions = append(whereConditions, "pgasrc.attgenerated = ''", "pgadst.attgenerated = ''")
	}

	query := fmt.Sprintf(`
	select
		pgcon.conname,
		pgc.relname as source_table,
		pgasrc.attname as source_column,
		dstlookupname.relname as dest_table,
		pgadst.attname as dest_column
	from pg_namespace pgn
		inner join pg_class pgc on pgn.oid = pgc.relnamespace and pgc.relkind = 'r'
		inner join pg_constraint pgcon on pgn.oid = pgcon.connamespace and pgc.oid = pgcon.conrelid
		inner join pg_class dstlookupname on pgcon.confrelid = dstlookupname.oid
		inner join pg_attribute pgasrc on pgc.oid = pgasrc.attrelid and pgasrc.attnum = ANY(pgcon.conkey)
		inner join pg_attribute pgadst on pgcon.confrelid = pgadst.attrelid and pgadst.attnum = ANY(pgcon.confkey)
	where %s
	order by pgcon.conname, source_table, source_column, dest_table, dest_column`,
		strings.Join(whereConditions, " and "),
	)

	var rows *sql.Rows
	var err error
	if rows, err = p.conn.Query(query, tableName, schema); err != nil {
		return nil, err
	}

	for rows.Next() {
		var fkey drivers.ForeignKey
		var sourceTable string

		fkey.Table = tableName
		err = rows.Scan(&fkey.Name, &sourceTable, &fkey.Column, &fkey.ForeignTable, &fkey.ForeignColumn)
		if err != nil {
			return nil, err
		}

		fkeys = append(fkeys, fkey)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return fkeys, nil
}

// TranslateColumnType converts postgres database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a Column object.
func (p *PostgresDriver) TranslateColumnType(c drivers.Column) drivers.Column {
	if c.Nullable {
		switch c.DBType {
		case "bigint", "bigserial":
			c.Type = "null.Int64"
		case "integer", "serial":
			c.Type = "null.Int"
		case "oid":
			c.Type = "null.Uint32"
		case "smallint", "smallserial":
			c.Type = "null.Int16"
		case "decimal", "numeric":
			c.Type = "types.NullDecimal"
		case "double precision":
			c.Type = "null.Float64"
		case "real":
			c.Type = "null.Float32"
		case "bit", "interval", "bit varying", "character", "money", "character varying", "cidr", "inet", "macaddr", "text", "uuid", "xml":
			c.Type = "null.String"
		case `"char"`:
			c.Type = "null.Byte"
		case "bytea":
			c.Type = "null.Bytes"
		case "json", "jsonb":
			c.Type = "null.JSON"
		case "boolean":
			c.Type = "null.Bool"
		case "date", "time", "timestamp without time zone", "timestamp with time zone", "time without time zone", "time with time zone":
			c.Type = "null.Time"
		case "point":
			c.Type = "pgeo.NullPoint"
		case "line":
			c.Type = "pgeo.NullLine"
		case "lseg":
			c.Type = "pgeo.NullLseg"
		case "box":
			c.Type = "pgeo.NullBox"
		case "path":
			c.Type = "pgeo.NullPath"
		case "polygon":
			c.Type = "pgeo.NullPolygon"
		case "circle":
			c.Type = "pgeo.NullCircle"
		case "ARRAY":
			var dbType string
			c.Type, dbType = getArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType += dbType
		case "USER-DEFINED":
			switch c.UDTName {
			case "hstore":
				c.Type = "types.HStore"
				c.DBType = "hstore"
			case "citext":
				c.Type = "null.String"
			default:
				c.Type = "string"
				fmt.Fprintf(os.Stderr, "warning: incompatible data type detected: %s\n", c.UDTName)
			}
		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "bigint", "bigserial":
			c.Type = "int64"
		case "integer", "serial":
			c.Type = "int"
		case "oid":
			c.Type = "uint32"
		case "smallint", "smallserial":
			c.Type = "int16"
		case "decimal", "numeric":
			c.Type = "types.Decimal"
		case "double precision":
			c.Type = "float64"
		case "real":
			c.Type = "float32"
		case "bit", "interval", "uuint", "bit varying", "character", "money", "character varying", "cidr", "inet", "macaddr", "text", "uuid", "xml":
			c.Type = "string"
		case `"char"`:
			c.Type = "types.Byte"
		case "json", "jsonb":
			c.Type = "types.JSON"
		case "bytea":
			c.Type = "[]byte"
		case "boolean":
			c.Type = "bool"
		case "date", "time", "timestamp without time zone", "timestamp with time zone", "time without time zone", "time with time zone":
			c.Type = "time.Time"
		case "point":
			c.Type = "pgeo.Point"
		case "line":
			c.Type = "pgeo.Line"
		case "lseg":
			c.Type = "pgeo.Lseg"
		case "box":
			c.Type = "pgeo.Box"
		case "path":
			c.Type = "pgeo.Path"
		case "polygon":
			c.Type = "pgeo.Polygon"
		case "circle":
			c.Type = "pgeo.Circle"
		case "ARRAY":
			var dbType string
			c.Type, dbType = getArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType += dbType
		case "USER-DEFINED":
			switch c.UDTName {
			case "hstore":
				c.Type = "types.HStore"
				c.DBType = "hstore"
			case "citext":
				c.Type = "string"
			default:
				c.Type = "string"
				fmt.Fprintf(os.Stderr, "warning: incompatible data type detected: %s\n", c.UDTName)
			}
		default:
			c.Type = "string"
		}
	}

	return c
}

// getArrayType returns the correct boil.Array type for each database type
func getArrayType(c drivers.Column) (string, string) {
	// If a domain is created with a statement like this: "CREATE DOMAIN
	// text_array AS TEXT[] CHECK ( ... )" then the array type will be null,
	// but the udt name will be whatever the underlying type is with a leading
	// underscore. Note that this code handles some types, but not nearly all
	// the possibities. Notably, an array of a user-defined type ("CREATE
	// DOMAIN my_array AS my_type[]") will be treated as an array of strings,
	// which is not guaranteed to be correct.
	if c.ArrType != nil {
		switch *c.ArrType {
		case "bigint", "bigserial", "integer", "serial", "smallint", "smallserial", "oid":
			return "types.Int64Array", *c.ArrType
		case "bytea":
			return "types.BytesArray", *c.ArrType
		case "bit", "interval", "uuint", "bit varying", "character", "money", "character varying", "cidr", "inet", "macaddr", "text", "uuid", "xml":
			return "types.StringArray", *c.ArrType
		case "boolean":
			return "types.BoolArray", *c.ArrType
		case "decimal", "numeric":
			return "types.DecimalArray", *c.ArrType
		case "double precision", "real":
			return "types.Float64Array", *c.ArrType
		default:
			return "types.StringArray", *c.ArrType
		}
	} else {
		switch c.UDTName {
		case "_int4", "_int8":
			return "types.Int64Array", c.UDTName
		case "_bytea":
			return "types.BytesArray", c.UDTName
		case "_bit", "_interval", "_varbit", "_char", "_money", "_varchar", "_cidr", "_inet", "_macaddr", "_citext", "_text", "_uuid", "_xml":
			return "types.StringArray", c.UDTName
		case "_bool":
			return "types.BoolArray", c.UDTName
		case "_numeric":
			return "types.DecimalArray", c.UDTName
		case "_float4", "_float8":
			return "types.Float64Array", c.UDTName
		default:
			return "types.StringArray", c.UDTName
		}
	}
}

// Imports for the postgres driver
func (p PostgresDriver) Imports() (importers.Collection, error) {
	var col importers.Collection

	col.All = importers.Set{
		Standard: importers.List{
			`"strconv"`,
		},
	}
	col.Singleton = importers.Map{
		"psql_upsert": {
			Standard: importers.List{
				`"fmt"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/volatiletech/strmangle"`,
				`"github.com/volatiletech/sqlboiler/v4/drivers"`,
			},
		},
	}
	col.TestSingleton = importers.Map{
		"psql_suites_test": {
			Standard: importers.List{
				`"testing"`,
			},
		},
		"psql_main_test": {
			Standard: importers.List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"regexp"`,
				`"strings"`,
			},
			ThirdParty: importers.List{
				`"github.com/kat-co/vala"`,
				`"github.com/friendsofgo/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql/driver"`,
				`"github.com/volatiletech/randomize"`,
				`_ "github.com/lib/pq"`,
			},
		},
	}
	col.BasedOnType = importers.Map{
		"null.Float32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Float64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Int": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Int8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Int16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Int32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Int64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Uint": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Uint8": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Uint16": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Uint32": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Uint64": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.String": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Bool": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Time": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.JSON": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"null.Bytes": {
			ThirdParty: importers.List{`"github.com/volatiletech/null/v8"`},
		},
		"time.Time": {
			Standard: importers.List{`"time"`},
		},
		"types.JSON": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.Decimal": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.BytesArray": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.Int64Array": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.Float64Array": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.BoolArray": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.StringArray": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.DecimalArray": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"types.HStore": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"pgeo.Point": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.Line": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.Lseg": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.Box": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.Path": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.Polygon": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"types.NullDecimal": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types"`},
		},
		"pgeo.Circle": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullPoint": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullLine": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullLseg": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullBox": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullPath": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullPolygon": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
		"pgeo.NullCircle": {
			ThirdParty: importers.List{`"github.com/volatiletech/sqlboiler/v4/types/pgeo"`},
		},
	}

	return col, nil
}

// getVersion gets the version of underlying database
func (p *PostgresDriver) getVersion() (int, error) {
	type versionInfoType struct {
		ServerVersionNum int `json:"server_version_num"`
	}
	versionInfo := &versionInfoType{}

	row := p.conn.QueryRow("SHOW server_version_num")
	if err := row.Scan(&versionInfo.ServerVersionNum); err != nil {
		return 0, err
	}

	return versionInfo.ServerVersionNum, nil
}
