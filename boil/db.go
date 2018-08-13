package boil

import (
	"context"
	"database/sql"
)

// Executor can perform SQL queries.
type Executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	Executor

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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

// Begin a transaction with the current global database handle.
func Begin() (Transactor, error) {
	creator, ok := currentDB.(Beginner)
	if !ok {
		panic("database does not support transactions")
	}

	return creator.Begin()
}

// ContextTransactor can commit and rollback, on top of being able to execute
// context-aware queries.
type ContextTransactor interface {
	Commit() error
	Rollback() error

	ContextExecutor
}

// ContextBeginner allows creation of context aware transactions with options.
type ContextBeginner interface {
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

// BeginTx begins a transaction with the current global database handle.
func BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	creator, ok := currentDB.(ContextBeginner)
	if !ok {
		panic("database does not support context-aware transactions")
	}

	return creator.BeginTx(ctx, opts)
}
