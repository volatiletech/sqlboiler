package boil

import (
	"database/sql"
	"os"
)

var (
	// currentDB is a global database handle for the package
	currentDB Executor
)

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

// DebugMode is a flag controlling whether generated sql statements and
// debug information is outputted to the DebugWriter handle
//
// NOTE: This should be disabled in production to avoid leaking sensitive data
var DebugMode = false

// DebugWriter is where the debug output will be sent if DebugMode is true
var DebugWriter = os.Stdout

// Begin a transaction
func Begin() (Transactor, error) {
	creator, ok := currentDB.(Beginner)
	if !ok {
		panic("database does not support transactions")
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
