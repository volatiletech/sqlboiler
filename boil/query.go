package boil

import (
	"database/sql"
	"fmt"
)

// joinKind is the type of join
type joinKind int

// Join type constants
const (
	JoinInner joinKind = iota
	JoinOuterLeft
	JoinOuterRight
	JoinNatural
)

// Query holds the state for the built up query
type Query struct {
	executor    Executor
	plainSQL    plainSQL
	delete      bool
	update      map[string]interface{}
	selectCols  []string
	modFunction string
	from        []string
	joins       []join
	where       []where
	groupBy     []string
	orderBy     []string
	having      []string
	limit       int
	offset      int
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
	kind   joinKind
	clause string
	args   []interface{}
}

// SQL makes a plainSQL query, usually for use with bind
func SQL(query string, args ...interface{}) *Query {
	return &Query{
		plainSQL: plainSQL{
			sql:  query,
			args: args,
		},
	}
}

// ExecQuery executes a query that does not need a row returned
func ExecQuery(q *Query) (sql.Result, error) {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
	}
	return q.executor.Exec(qs, args...)
}

// ExecQueryOne executes the query for the One finisher and returns a row
func ExecQueryOne(q *Query) *sql.Row {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
	}
	return q.executor.QueryRow(qs, args...)
}

// ExecQueryAll executes the query for the All finisher and returns multiple rows
func ExecQueryAll(q *Query) (*sql.Rows, error) {
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
	q.modFunction = "COUNT"
}

// SetAvg on the query.
func SetAvg(q *Query) {
	q.modFunction = "AVG"
}

// SetMax on the query.
func SetMax(q *Query) {
	q.modFunction = "MAX"
}

// SetMin on the query.
func SetMin(q *Query) {
	q.modFunction = "MIN"
}

// SetSum on the query.
func SetSum(q *Query) {
	q.modFunction = "SUM"
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
func AppendInnerJoin(q *Query, clause string, args ...interface{}) {
	q.joins = append(q.joins, join{clause: clause, kind: JoinInner, args: args})
}

// SetInnerJoin on the query.
func SetInnerJoin(q *Query, clause string, args ...interface{}) {
	q.joins = append([]join(nil), join{clause: clause, kind: JoinInner, args: args})
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
