package drivers

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/bdb"
)

// MSSQLDriver holds the database connection string and a handle
// to the database connection.
type MSSQLDriver struct {
	connStr string
	dbConn  *sql.DB
}

// NewMSSQLDriver takes the database connection details as parameters and
// returns a pointer to a MSSQLDriver object. Note that it is required to
// call MSSQLDriver.Open() and MSSQLDriver.Close() to open and close
// the database connection once an object has been obtained.
func NewMSSQLDriver(user, pass, dbname, host string, port int, sslmode string) *MSSQLDriver {
	driver := MSSQLDriver{
		connStr: MSSQLBuildQueryString(user, pass, dbname, host, port, sslmode),
	}

	return &driver
}

// MSSQLBuildQueryString builds a query string for MSSQL.
func MSSQLBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {

	query := url.Values{}
	query.Add("database", dbname)
	query.Add("encrypt", sslmode)

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%d", host, port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}

	return u.String()
}

// Open opens the database connection using the connection string
func (m *MSSQLDriver) Open() error {
	var err error
	m.dbConn, err = sql.Open("mssql", m.connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (m *MSSQLDriver) Close() {
	m.dbConn.Close()
}

// UseLastInsertID returns false for mssql
func (m *MSSQLDriver) UseLastInsertID() bool {
	return false
}

// UseTopClause returns true to indicate MS SQL supports SQL TOP clause
func (m *MSSQLDriver) UseTopClause() bool {
	return true
}

// TableNames connects to the postgres database and
// retrieves all table names from the information_schema where the
// table schema is public.
func (m *MSSQLDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	var names []string

	query := fmt.Sprintf(`select table_name from information_schema.tables where table_schema = ? and table_type = 'BASE TABLE'`)
	args := []interface{}{schema}
	if len(whitelist) > 0 {
		query += fmt.Sprintf(" and table_name in (%s);", strings.Repeat(",?", len(whitelist))[1:])
		for _, w := range whitelist {
			args = append(args, w)
		}
	} else if len(blacklist) > 0 {
		query += fmt.Sprintf(" and table_name not in (%s);", strings.Repeat(",?", len(blacklist))[1:])
		for _, b := range blacklist {
			args = append(args, b)
		}
	}

	rows, err := m.dbConn.Query(query, args...)

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
func (m *MSSQLDriver) Columns(schema, tableName string) ([]bdb.Column, error) {
	var columns []bdb.Column

	rows, err := m.dbConn.Query(`
	SELECT column_name,
       CASE
         WHEN CHARACTER_MAXIMUM_LENGTH IS NULL THEN data_type
         ELSE data_type + '(' + CAST(CHARACTER_MAXIMUM_LENGTH AS VARCHAR) + ')'
       END AS full_type,
       data_type,
	   column_default,
       CASE
         WHEN IS_NULLABLE = 'YES' THEN 1
         ELSE 0
       END AS is_nullable,
       CASE
         WHEN EXISTS (SELECT c.column_name
                      FROM information_schema.table_constraints tc
                        INNER JOIN information_schema.key_column_usage kcu
                                ON tc.constraint_name = kcu.constraint_name
                               AND tc.table_name = kcu.table_name
                               AND tc.table_schema = kcu.table_schema
                      WHERE c.column_name = kcu.column_name
                      AND   tc.table_name = c.table_name
                      AND   (tc.constraint_type = 'PRIMARY KEY' OR tc.constraint_type = 'UNIQUE')
                      AND   (SELECT COUNT(*)
                             FROM information_schema.key_column_usage
                             WHERE table_schema = kcu.table_schema
                             AND   table_name = tc.table_name
                             AND   constraint_name = tc.constraint_name) = 1) THEN 1
         ELSE 0
       END AS is_unique
	FROM INFORMATION_SCHEMA.columns c
	where table_name = ? and table_schema = ?;
	`, tableName, schema)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var colName, colType, colFullType string
		var nullable, unique bool
		var defaultValue *string
		if err := rows.Scan(&colName, &colFullType, &colType, &defaultValue, &nullable, &unique); err != nil {
			return nil, errors.Wrapf(err, "unable to scan for table %s", tableName)
		}

		column := bdb.Column{
			Name:       colName,
			FullDBType: colFullType, // example: tinyint(1) instead of tinyint
			DBType:     colType,
			Nullable:   nullable,
			Unique:     unique,
		}

		if defaultValue != nil && *defaultValue != "NULL" {
			column.Default = *defaultValue
		}

		columns = append(columns, column)
	}

	return columns, nil
}

// PrimaryKeyInfo looks up the primary key for a table.
func (m *MSSQLDriver) PrimaryKeyInfo(schema, tableName string) (*bdb.PrimaryKey, error) {
	pkey := &bdb.PrimaryKey{}
	var err error

	query := `
	select tc.constraint_name
	from information_schema.table_constraints as tc
	where tc.table_name = ? and tc.constraint_type = 'PRIMARY KEY' and tc.table_schema = ?;`

	row := m.dbConn.QueryRow(query, tableName, schema)
	if err = row.Scan(&pkey.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	queryColumns := `
	select kcu.column_name
	from   information_schema.key_column_usage as kcu
	where  table_name = ? and constraint_name = ? and table_schema = ?;`

	var rows *sql.Rows
	if rows, err = m.dbConn.Query(queryColumns, tableName, pkey.Name, schema); err != nil {
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
func (m *MSSQLDriver) ForeignKeyInfo(schema, tableName string) ([]bdb.ForeignKey, error) {
	var fkeys []bdb.ForeignKey

	query := `
	SELECT 
		ccu.constraint_name AS SourceConstraint
		,ccu.table_name AS SourceTable
		,ccu.column_name AS SourceColumn
		,kcu.table_name AS TargetTable
		,kcu.column_name AS TargetColumn
	FROM INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE ccu
    INNER JOIN INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS rc
        ON ccu.CONSTRAINT_NAME = rc.CONSTRAINT_NAME 
    INNER JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
        ON kcu.CONSTRAINT_NAME = rc.UNIQUE_CONSTRAINT_NAME  
	where ccu.table_schema = ? and ccu.constraint_schema = ? and ccu.table_name = ?
	`

	var rows *sql.Rows
	var err error
	if rows, err = m.dbConn.Query(query, schema, schema, tableName); err != nil {
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
func (m *MSSQLDriver) TranslateColumnType(c bdb.Column) bdb.Column {
	if c.Nullable {
		switch c.DBType {
		case "bit":
			c.Type = "bool"
		case "tinyint":
			c.Type = "int8"
		case "smallint":
			c.Type = "null.Int16"
		case "mediumint":
			c.Type = "null.Int32"
		case "int", "integer":
			c.Type = "null.Int"
		case "bigint":
			c.Type = "null.Int64"
		case "float":
			c.Type = "null.Float32"
		case "double", "double precision", "real":
			c.Type = "null.Float64"
		case "boolean", "bool":
			c.Type = "null.Bool"
		case "date", "datetime", "timestamp", "time":
			c.Type = "null.Time"
		case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
			c.Type = "null.Bytes"
		case "json":
			c.Type = "types.JSON"
		default:
			c.Type = "null.String"
		}
	} else {
		switch c.DBType {
		case "bit":
			c.Type = "bool"
		case "tinyint":
			c.Type = "int8"
		case "smallint":
			c.Type = "int16"
		case "mediumint":
			c.Type = "int32"
		case "int", "integer":
			c.Type = "int"
		case "bigint":
			c.Type = "int64"
		case "float":
			c.Type = "float32"
		case "double", "double precision", "real":
			c.Type = "float64"
		case "boolean", "bool":
			c.Type = "bool"
		case "date", "datetime", "timestamp", "time":
			c.Type = "time.Time"
		case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
			c.Type = "[]byte"
		case "json":
			c.Type = "types.JSON"
		default:
			c.Type = "string"
		}
	}

	return c
}

// RightQuote is the quoting character for the right side of the identifier
func (m *MSSQLDriver) RightQuote() byte {
	return '"'
}

// LeftQuote is the quoting character for the left side of the identifier
func (m *MSSQLDriver) LeftQuote() byte {
	return '"'
}

// IndexPlaceholders returns true to indicate MS SQL supports indexed placeholders
func (m *MSSQLDriver) IndexPlaceholders() bool {
	return true
}
