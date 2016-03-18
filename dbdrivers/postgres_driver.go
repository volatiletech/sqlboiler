package dbdrivers

import (
	"database/sql"
	"fmt"

	// Import the postgres driver
	_ "github.com/lib/pq"
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
func NewPostgresDriver(user, pass, dbname, host string, port int) *PostgresDriver {
	driver := PostgresDriver{
		connStr: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
			user, pass, dbname, host, port),
	}

	return &driver
}

// Open opens the database connection using the connection string
func (d *PostgresDriver) Open() error {
	var err error
	d.dbConn, err = sql.Open("postgres", d.connStr)
	if err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func (d *PostgresDriver) Close() {
	d.dbConn.Close()
}

// GetAllTables connects to the postgres database and
// retrieves all table names from the information_schema where the
// table schema is public. It excludes common migration tool tables
// such as gorp_migrations
func (d *PostgresDriver) GetAllTables() ([]string, error) {
	var tableNames []string

	rows, err := d.dbConn.Query(`select table_name from
    information_schema.tables where table_schema='public'
    and table_name <> 'gorp_migrations'`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}

// GetTableInfo takes a table name and attempts to retrieve the table information
// from the database information_schema.columns. It retrieves the column names
// and column types and returns those as a []DBColumn after ParseTableInfo()
// converts the SQL types to Go types, for example: "varchar" to "string"
func (d *PostgresDriver) GetTableInfo(tableName string) ([]DBColumn, error) {
	var table []DBColumn

	rows, err := d.dbConn.Query(`
		SELECT c.column_name, c.data_type, c.is_nullable,
		CASE WHEN pk.column_name IS NOT NULL THEN 'PRIMARY KEY' ELSE '' END AS KeyType
		FROM information_schema.columns c
		LEFT JOIN (
		  SELECT ku.table_name, ku.column_name
		  FROM information_schema.table_constraints AS tc
		  INNER JOIN information_schema.key_column_usage AS ku
		    ON tc.constraint_type = 'PRIMARY KEY'
		    AND tc.constraint_name = ku.constraint_name
		) pk
		ON c.table_name = pk.table_name
		AND c.column_name = pk.column_name
		WHERE c.table_name=$1
	`, tableName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var colName, colType, isNullable, isPrimary string
		if err := rows.Scan(&colName, &colType, &isNullable, &isPrimary); err != nil {
			return nil, err
		}
		t := d.ParseTableInfo(colName, colType, isNullable == "YES", isPrimary == "PRIMARY KEY")
		table = append(table, t)
	}

	return table, nil
}

// ParseTableInfo converts postgres database types to Go types, for example
// "varchar" to "string" and "bigint" to "int64". It returns this parsed data
// as a DBColumn object.
func (d *PostgresDriver) ParseTableInfo(colName, colType string, isNullable bool, isPrimary bool) DBColumn {
	var t DBColumn

	t.Name = colName
	t.IsPrimaryKey = isPrimary
	t.IsNullable = isNullable

	if isNullable {
		switch colType {
		case "bigint", "bigserial", "integer", "smallint", "smallserial", "serial":
			t.Type = "null.Int"
		case "bit", "bit varying", "character", "character varying", "cidr", "inet", "json", "macaddr", "text", "uuid", "xml":
			t.Type = "null.String"
		case "boolean":
			t.Type = "null.Bool"
		case "date", "interval", "time", "timestamp without time zone", "timestamp with time zone":
			t.Type = "null.Time"
		case "double precision", "money", "numeric", "real":
			t.Type = "null.Float"
		default:
			t.Type = "null.String"
		}
	} else {
		switch colType {
		case "bigint", "bigserial", "integer", "smallint", "smallserial", "serial":
			t.Type = "int64"
		case "bit", "bit varying", "character", "character varying", "cidr", "inet", "json", "macaddr", "text", "uuid", "xml":
			t.Type = "string"
		case "bytea":
			t.Type = "[]byte"
		case "boolean":
			t.Type = "bool"
		case "date", "interval", "time", "timestamp without time zone", "timestamp with time zone":
			t.Type = "time.Time"
		case "double precision", "money", "numeric", "real":
			t.Type = "float64"
		default:
			t.Type = "string"
		}
	}

	return t
}
