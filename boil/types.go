package boil

import "database/sql"

// DB implements the functions necessary for the templates to function.
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// M type is for providing where filters to Where helpers.
type M map[string]interface{}
