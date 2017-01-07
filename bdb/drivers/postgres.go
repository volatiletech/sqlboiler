package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	// Side-effect import sql driver

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/strmangle"
)

// PostgresDriver holds the database connection string and a handle
// to the database connection.
type PostgresDriver struct {
	connStr string
	dbConn  *sql.DB
}

// NewPostgresDriver takes the database connection details as parameters and
// returns a pointer to a PostgresDriver object. Note that it is required to
// call PostgresDriver.Open() and PostgresDriver.Close() to open and close
// the database connection once an object has been obtained.
func NewPostgresDriver(user, pass, dbname, host string, port int, sslmode string) *PostgresDriver {
	driver := PostgresDriver{
		connStr: PostgresBuildQueryString(user, pass, dbname, host, port, sslmode),
	}

	return &driver
}

// PostgresBuildQueryString builds a query string.
func PostgresBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
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

// Open opens the database connection using the connection string
func (p *PostgresDriver) Open() error {
	var err error
	p.dbConn, err = sql.Open("postgres", p.connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (p *PostgresDriver) Close() {
	p.dbConn.Close()
}

// UseLastInsertID returns false for postgres
func (p *PostgresDriver) UseLastInsertID() bool {
	return false
}

// TableNames connects to the postgres database and
// retrieves all table names from the information_schema where the
// table schema is schema. It uses a whitelist and blacklist.
func (p *PostgresDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var names []string

	query := fmt.Sprintf(`select table_name from information_schema.tables where table_schema = $1`)
	args := []interface{}{schema}
	if len(whitelist) > 0 {
		query += fmt.Sprintf(" and table_name in (%s);", strmangle.Placeholders(true, len(whitelist), 2, 1))
		for _, w := range whitelist {
			args = append(args, w)
		}
	} else if len(blacklist) > 0 {
		query += fmt.Sprintf(" and table_name not in (%s);", strmangle.Placeholders(true, len(blacklist), 2, 1))
		for _, b := range blacklist {
			args = append(args, b)
		}
	}

	rows, err := p.dbConn.Query(query, args...)

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
func (p *PostgresDriver) Columns(schema, tableName string) ([]bdb.Column, error) {
	var columns []bdb.Column

	rows, err := p.dbConn.Query(`
		select
		c.column_name,
		(
			case when c.data_type = 'USER-DEFINED' and c.udt_name <> 'hstore'
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
						where pg_type.typtype = 'b' and pg_type.typname = ('_' || c.udt_name)
						limit 1
					)
					order by pg_enum.enumsortorder
				) as labels
			)
			else c.data_type
			end
		) as column_type,

		c.udt_name,
		e.data_type as array_type,
		c.column_default,

		c.is_nullable = 'YES' as is_nullable,
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
		left join information_schema.element_types e
			on ((c.table_catalog, c.table_schema, c.table_name, 'TABLE', c.dtd_identifier)
			= (e.object_catalog, e.object_schema, e.object_name, e.object_type, e.collection_type_identifier))
		where c.table_name = $2 and c.table_schema = $1;
	`, schema, tableName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, udtName string
		var defaultValue, arrayType *string
		var nullable, unique bool
		if err := rows.Scan(&colName, &colType, &udtName, &arrayType, &defaultValue, &nullable, &unique); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		column := bdb.Column{
			Name:     colName,
			DBType:   colType,
			ArrType:  arrayType,
			UDTName:  udtName,
			Nullable: nullable,
			Unique:   unique,
		}
		if defaultValue != nil {
			column.Default = *defaultValue
		}

		columns = append(columns, column)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (p *PostgresDriver) PrimaryKeyInfo(schema, tableName string) (*bdb.PrimaryKey, error) {
	pkey := &bdb.PrimaryKey{}
	var err error

	query := `
	select tc.constraint_name
	from information_schema.table_constraints as tc
	where tc.table_name = $1 and tc.constraint_type = 'PRIMARY KEY' and tc.table_schema = $2;`

	row := p.dbConn.QueryRow(query, tableName, schema)
	if err = row.Scan(&pkey.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryColumns := `
	select kcu.column_name
	from   information_schema.key_column_usage as kcu
	where  constraint_name = $1 and table_schema = $2;`

	var rows *sql.Rows
	if rows, err = p.dbConn.Query(queryColumns, pkey.Name, schema); err != nil {
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
func (p *PostgresDriver) ForeignKeyInfo(schema, tableName string) ([]bdb.ForeignKey, error) {
	var fkeys []bdb.ForeignKey

	query := `
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
	where pgn.nspname = $2 and pgc.relname = $1 and pgcon.contype = 'f'`

	var rows *sql.Rows
	var err error
	if rows, err = p.dbConn.Query(query, tableName, schema); err != nil {
		return nil, err
	}

	for rows.Next() {
		var fkey bdb.ForeignKey
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
func (p *PostgresDriver) TranslateColumnType(c bdb.Column) bdb.Column {
	if c.Nullable {
		switch c.DBType {
		case "bigint", "bigserial":
			c.Type = "null.Int64"
		case "integer", "serial":
			c.Type = "null.Int"
		case "smallint", "smallserial":
			c.Type = "null.Int16"
		case "decimal", "numeric", "double precision":
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
		case "date", "time", "timestamp without time zone", "timestamp with time zone":
			c.Type = "null.Time"
		case "ARRAY":
			if c.ArrType == nil {
				panic("unable to get postgres ARRAY underlying type")
			}
			c.Type = getArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType = c.DBType + *c.ArrType
		case "USER-DEFINED":
			if c.UDTName == "hstore" {
				c.Type = "types.HStore"
				c.DBType = "hstore"
			} else {
				c.Type = "string"
				fmt.Printf("Warning: Incompatible data type detected: %s\n", c.UDTName)
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
		case "smallint", "smallserial":
			c.Type = "int16"
		case "decimal", "numeric", "double precision":
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
		case "date", "time", "timestamp without time zone", "timestamp with time zone":
			c.Type = "time.Time"
		case "ARRAY":
			c.Type = getArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType = c.DBType + *c.ArrType
		case "USER-DEFINED":
			if c.UDTName == "hstore" {
				c.Type = "types.HStore"
				c.DBType = "hstore"
			} else {
				c.Type = "string"
				fmt.Printf("Warning: Incompatible data type detected: %s\n", c.UDTName)
			}
		default:
			c.Type = "string"
		}
	}

	return c
}

// getArrayType returns the correct boil.Array type for each database type
func getArrayType(c bdb.Column) string {
	switch *c.ArrType {
	case "bigint", "bigserial", "integer", "serial", "smallint", "smallserial":
		return "types.Int64Array"
	case "bytea":
		return "types.BytesArray"
	case "bit", "interval", "uuint", "bit varying", "character", "money", "character varying", "cidr", "inet", "macaddr", "text", "uuid", "xml":
		return "types.StringArray"
	case "boolean":
		return "types.BoolArray"
	case "decimal", "numeric", "double precision", "real":
		return "types.Float64Array"
	default:
		return "types.StringArray"
	}
}

// RightQuote is the quoting character for the right side of the identifier
func (p *PostgresDriver) RightQuote() byte {
	return '"'
}

// LeftQuote is the quoting character for the left side of the identifier
func (p *PostgresDriver) LeftQuote() byte {
	return '"'
}

// IndexPlaceholders returns true to indicate PSQL supports indexed placeholders
func (p *PostgresDriver) IndexPlaceholders() bool {
	return true
}
