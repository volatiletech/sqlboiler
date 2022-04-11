package boil

import (
	"context"
	"database/sql"
)

// ContextExecutor can perform SQL queries with context
type ContextExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
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
	creator, ok := currentContextDB.(ContextBeginner)
	if !ok {
		panic("database does not support context-aware transactions")
	}

	return creator.BeginTx(ctx, opts)
}
