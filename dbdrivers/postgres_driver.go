package dbdrivers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	connStr string
	dbConn  *sql.DB
}

func NewPostgresDriver(user, pass, dbname, host string, port int) *PostgresDriver {
	driver := PostgresDriver{
		connStr: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d",
			user, pass, dbname, host, port),
	}

	return &driver
}

func (d *PostgresDriver) Open() error {
	var err error
	d.dbConn, err = sql.Open("postgres", d.connStr)
	if err != nil {
		return err
	}

	return nil
}

func (d *PostgresDriver) Close() {
	d.dbConn.Close()
}

func (d *PostgresDriver) GetAllTableNames() ([]string, error) {
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

func (d *PostgresDriver) GetTableInfo(tableName string) ([]DBTable, error) {
	var tableInfo []DBTable

	rows, err := d.dbConn.Query(`select column_name, data_type from
    information_schema.columns where table_name=$1`, tableName)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var colName, colType string
		if err := rows.Scan(&colName, &colType); err != nil {
			return nil, err
		}
		tableInfo = append(tableInfo, d.ParseTableInfo(colName, colType))
	}

	return tableInfo, nil
}

func (d *PostgresDriver) ParseTableInfo(colName, colType string) DBTable {
	t := DBTable{}

	t.ColName = colName
	switch colType {
	case "bigint", "bigserial", "integer", "smallint", "smallserial", "serial":
		t.ColType = "int64"
	case "bit", "bit varying", "character", "character varying", "cidr", "inet", "json", "macaddr", "text", "uuid", "xml":
		t.ColType = "string"
	case "bytea":
		t.ColType = "[]byte"
	case "boolean":
		t.ColType = "bool"
	case "date", "interval", "time", "timestamp without time zone", "timestamp with time zone":
		t.ColType = "time.Time"
	case "double precision", "money", "numeric", "real":
		t.ColType = "float64"
	default:
		t.ColType = "string"
	}

	return t
}
