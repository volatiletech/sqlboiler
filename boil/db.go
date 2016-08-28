package boil

import "database/sql"

// Executor can perform SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// Transactor can commit and rollback, on top of being able to execute queries.
type Transactor interface {
	Commit() error
	Rollback() error

	Executor
}

// Beginner begins transactions.
type Beginner interface {
	Begin() (*sql.Tx, error)
}

// Begin a transaction
func Begin() (Transactor, error) {
	creator, ok := currentDB.(Beginner)
	if !ok {
		panic("database does not support transactions")
	}

	return creator.Begin()
}
