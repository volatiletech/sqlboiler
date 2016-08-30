package sqlmock

import (
	"database/sql/driver"
)

type statement struct {
	conn  *sqlmock
	query string
	err   error
}

func (stmt *statement) Close() error {
	return stmt.err
}

func (stmt *statement) NumInput() int {
	return -1
}

func (stmt *statement) Exec(args []driver.Value) (driver.Result, error) {
	return stmt.conn.Exec(stmt.query, args)
}

func (stmt *statement) Query(args []driver.Value) (driver.Rows, error) {
	return stmt.conn.Query(stmt.query, args)
}
