package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	// Side-effect import sql driver

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/bdb"
	"github.com/volatiletech/sqlboiler/strmangle"
	"regexp"
)

var re = regexp.MustCompile(`\(([^\)]+)\)`)

// CockroachDriver holds the database connection string and a handle
// to the database connection.
type CockroachDriver struct {
	connStr string
	dbConn  *sql.DB
}

// NewCockroachDriver takes the database connection details as parameters and
// returns a pointer to a CockroachDriver object. Note that it is required to
// call CockroachDriver.Open() and CockroachDriver.Close() to open and close
// the database connection once an object has been obtained.
func NewCockroachDriver(user, pass, dbname, host string, port int, sslmode string) *CockroachDriver {
	driver := CockroachDriver{
		connStr: CockroachBuildQueryString(user, pass, dbname, host, port, sslmode),
	}

	return &driver
}

// CockroachBuildQueryString builds a query string.
func CockroachBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", user, pass, host, port, dbname, sslmode)
}

// Open opens the database connection using the connection string
func (p *CockroachDriver) Open() error {
	var err error
	p.dbConn, err = sql.Open("postgres", p.connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (p *CockroachDriver) Close() {
	p.dbConn.Close()
}

// UseLastInsertID returns false for Cockroach
func (p *CockroachDriver) UseLastInsertID() bool {
	return false
}

// UseTopClause returns false to indicate PSQL doesnt support SQL TOP clause
func (m *CockroachDriver) UseTopClause() bool {
	return false
}

// TableNames connects to the Cockroach database and
// retrieves all table names from the information_schema where the
// table schema is schema. It uses a whitelist and blacklist.
func (p *CockroachDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
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
func (p *CockroachDriver) Columns(schema, tableName string) ([]bdb.Column, error) {
	var columns []bdb.Column

	rows, err := p.dbConn.Query(`
		select
		distinct c.column_name,
		c.data_type,
		c.column_default,
		(case when c.is_nullable = 'NO' then FALSE
			else TRUE end) as is_nullable,
		(case when kcu.constraint_name is not null then TRUE
			else False END) as is_unique
		from information_schema.columns as c
			left join information_schema.key_column_usage kcu on c.table_name = kcu.table_name
				and c.table_schema = kcu.table_schema and c.column_name = kcu.column_name
		where c.table_schema = $1 and c.table_name = $2;
	`, schema, tableName)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, udtName string
		var defaultValue, arrayType *string
		var nullable, unique bool
		if err := rows.Scan(&colName, &colType, &defaultValue, &nullable, &unique); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		dbType := strings.ToLower(re.ReplaceAllString(colType, ""))
		// todo find a better way to infer this
		tmp := strings.Replace(dbType, "[]", "", 1)
		if dbType != tmp{
			arrayType = &tmp
			dbType = "array"
		}

		column := bdb.Column{
			Name:     colName,
			DBType:   dbType,
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
func (p *CockroachDriver) PrimaryKeyInfo(schema, tableName string) (*bdb.PrimaryKey, error) {
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
	where  constraint_name = $1 and table_schema = $2 and table_name = $3;`

	var rows *sql.Rows
	if rows, err = p.dbConn.Query(queryColumns, pkey.Name, schema, tableName); err != nil {
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
func (p *CockroachDriver) ForeignKeyInfo(schema, tableName string) ([]bdb.ForeignKey, error) {
	var fkeys []bdb.ForeignKey

	query := `
	select
	  distinct pgcon.conname,
	  pgc.relname as source_table,
	  kcu.column_name as source_column,
	  dstlookupname.relname as dest_table,
	  pgadst.attname as dest_column
	from pg_namespace pgn
	  inner join pg_class pgc on pgn.oid = pgc.relnamespace and pgc.relkind = 'r'
	  inner join pg_constraint pgcon on pgn.oid = pgcon.connamespace and pgc.oid = pgcon.conrelid
	  inner join pg_class dstlookupname on pgcon.confrelid = dstlookupname.oid
	  left join information_schema.key_column_usage kcu on pgcon.conname = kcu.constraint_name and pgc.relname = kcu.table_name
	  left join information_schema.key_column_usage kcudst on pgcon.conname = kcu.constraint_name and dstlookupname.relname = kcu.table_name
	  inner join pg_attribute pgadst on pgcon.confrelid = pgadst.attrelid and pgadst.attnum = ANY(pgcon.confkey)
	where pgn.nspname = $2 and pgc.relname = $1 and pgcon.contype = 'f';
	`

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

// TranslateColumnType converts Cockroach database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a Column object.
func (p *CockroachDriver) TranslateColumnType(c bdb.Column) bdb.Column {
	// parse DB type
	if c.Nullable {
		switch c.DBType {
		case "int", "serial":
			c.Type = "null.Int64"
		case "int4":
			c.Type = "null.Int32"
		case "int2":
			c.Type = "null.Int16"
		case "decimal", "double precision":
			c.Type = "null.Float64"
		case "real":
			c.Type = "null.Float32"
		case "string", "uuid", "collate":
			c.Type = "null.String"
		case "bytes":
			c.Type = "null.Bytes"
		case "bool":
			c.Type = "null.Bool"
		case "date", "timestamp with time zone", "timestamp without time zone":
			c.Type = "null.Time"
		case "array":
			if c.ArrType == nil {
				panic("unable to get Cockroach ARRAY underlying type")
			}
			c.Type = getCockroachArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType = c.DBType + *c.ArrType
		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "int", "serial":
			c.Type = "int64"
		case "int4":
			c.Type = "int32"
		case "int2":
			c.Type = "int16"
		case "decimal", "double precision":
			c.Type = "float64"
		case "real":
			c.Type = "float32"
		case "string", "uuid", "collate":
			c.Type = "string"
		case "bytes":
			c.Type = "[]byte"
		case "bool":
			c.Type = "bool"
		case "date", "timestamp with time zone", "timestamp without time zone":
			c.Type = "time.Time"
		case "array":
			if c.ArrType == nil {
				panic("unable to get Cockroach ARRAY underlying type")
			}
			c.Type = getCockroachArrayType(c)
			// Make DBType something like ARRAYinteger for parsing with randomize.Struct
			c.DBType = c.DBType + *c.ArrType
		default:
			c.Type = "string"
		}
	}
	return c
}

// getCockroachArrayType returns the correct boil.Array type for each database type
func getCockroachArrayType(c bdb.Column) string {
	switch *c.ArrType {
	case "int", "serial":
		return "types.Int64Array"
	case "bytes":
		return "types.BytesArray"
	case "string", "uuid", "collate":
		return "types.StringArray"
	case "bool":
		return "types.BoolArray"
	case "decimal", "numeric", "double precision", "real":
		return "types.Float64Array"
	default:
		return "types.StringArray"
	}
}

// RightQuote is the quoting character for the right side of the identifier
func (p *CockroachDriver) RightQuote() byte {
	return '"'
}

// LeftQuote is the quoting character for the left side of the identifier
func (p *CockroachDriver) LeftQuote() byte {
	return '"'
}

// IndexPlaceholders returns true to indicate PSQL supports indexed placeholders
func (p *CockroachDriver) IndexPlaceholders() bool {
	return true
}
