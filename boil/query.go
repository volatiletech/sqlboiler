package boil

import (
	"database/sql"
	"fmt"
)

// Query holds the state for the built up query
type Query struct {
	executor   Executor
	plainSQL   plainSQL
	delete     bool
	update     map[string]interface{}
	selectCols []string
	count      bool
	from       []string
	innerJoins []join
	where      []where
	groupBy    []string
	orderBy    []string
	having     []string
	limit      int
	offset     int
}

type where struct {
	clause      string
	orSeperator bool
	args        []interface{}
}

type plainSQL struct {
	sql  string
	args []interface{}
}

type join struct {
	on   string
	args []interface{}
}

// ExecStatement executes a query that does not need a row returned
func ExecStatement(q *Query) (sql.Result, error) {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
	}
	return q.executor.Exec(qs, args...)
}

// ExecQuery executes the query for the All finisher and returns multiple rows
func ExecQuery(q *Query) (*sql.Rows, error) {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
	}
	return q.executor.Query(qs, args...)
}

// SetSQL on the query.
func SetSQL(q *Query, sql string, args ...interface{}) {
	q.plainSQL = plainSQL{sql: sql, args: args}
}

// SetCount on the query.
func SetCount(q *Query) {
	q.count = true
}

// SetDelete on the query.
func SetDelete(q *Query) {
	q.delete = true
}

// SetUpdate on the query.
func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

// SetExecutor on the query.
func SetExecutor(q *Query, exec Executor) {
	q.executor = exec
}

// AppendSelect on the query.
func AppendSelect(q *Query, columns ...string) {
	q.selectCols = append(q.selectCols, columns...)
}

// SetSelect replaces the current select clause.
func SetSelect(q *Query, columns ...string) {
	q.selectCols = append([]string(nil), columns...)
}

// Select returns the select columns in the query.
func Select(q *Query) []string {
	cols := make([]string, len(q.selectCols))
	copy(cols, q.selectCols)
	return cols
}

// AppendFrom on the query.
func AppendFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

// SetFrom replaces the current from statements.
func SetFrom(q *Query, from ...string) {
	q.from = append([]string(nil), from...)
}

// AppendInnerJoin on the query.
func AppendInnerJoin(q *Query, on string, args ...interface{}) {
	q.innerJoins = append(q.innerJoins, join{on: on, args: args})
}

// SetInnerJoin on the query.
func SetInnerJoin(q *Query, on string, args ...interface{}) {
	q.innerJoins = append([]join(nil), join{on: on, args: args})
}

// AppendWhere on the query.
func AppendWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// SetWhere on the query.
func SetWhere(q *Query, clause string, args ...interface{}) {
	q.where = append([]where(nil), where{clause: clause, args: args})
}

// SetLastWhereAsOr sets the or seperator for the last element in the where slice
func SetLastWhereAsOr(q *Query) {
	q.where[len(q.where)-1].orSeperator = true
}

// ApplyGroupBy on the query.
func ApplyGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// SetGroupBy on the query.
func SetGroupBy(q *Query, clause string) {
	q.groupBy = append([]string(nil), clause)
}

// ApplyOrderBy on the query.
func ApplyOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}

// SetOrderBy on the query.
func SetOrderBy(q *Query, clause string) {
	q.orderBy = append([]string(nil), clause)
}

// ApplyHaving on the query.
func ApplyHaving(q *Query, clause string) {
	q.having = append(q.having, clause)
}

// SetHaving on the query.
func SetHaving(q *Query, clause string) {
	q.having = append([]string(nil), clause)
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}
