package boil

import "database/sql"

type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Transactor interface {
	Commit() error
	Rollback() error

	Executor
}

type Creator interface {
	Begin() (*sql.Tx, error)
}

var currentDB Executor

func Begin() (Transactor, error) {
	creator, ok := currentDB.(Creator)
	if !ok {
		panic("Your database does not support transactions.")
	}

	return creator.Begin()
}

// SetDB initializes the database handle for all template db interactions
func SetDB(db Executor) {
	currentDB = db
}

// GetDB retrieves the global state database handle
func GetDB() Executor {
	return currentDB
}
