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
	executor   Executor
	dialect    *Dialect
	plainSQL   plainSQL
	load       []string
	delete     bool
	update     map[string]interface{}
	selectCols []string
	count      bool
	from       []string
	joins      []join
	where      []where
	in         []in
	groupBy    []string
	orderBy    []string
	having     []having
	limit      int
	offset     int
	forlock    string
}

// Dialect holds values that direct the query builder
// how to build compatible queries for each database.
// Each database driver needs to implement functions
// that provide these values.
type Dialect struct {
	// The left quote character for SQL identifiers
	LQ string
	// The right quote character for SQL identifiers
	RQ string
	// Bool flag indicating whether indexed
	// placeholders ($1) are used, or ? placeholders.
	IndexPlaceholders bool
}

type where struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type in struct {
	clause      string
	orSeparator bool
	args        []interface{}
}

type having struct {
	clause string
	args   []interface{}
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
func SQL(exec Executor, query string, args ...interface{}) *Query {
	return &Query{
		executor: exec,
		plainSQL: plainSQL{
			sql:  query,
			args: args,
		},
	}
}

// SQLG makes a plainSQL query using the global Executor, usually for use with bind
func SQLG(query string, args ...interface{}) *Query {
	return SQL(GetDB(), query, args...)
}

// ExecQuery executes a query that does not need a row returned
func ExecQuery(q *Query) (sql.Result, error) {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
		fmt.Fprintln(DebugWriter, args)
	}
	return q.executor.Exec(qs, args...)
}

// ExecQueryOne executes the query for the One finisher and returns a row
func ExecQueryOne(q *Query) *sql.Row {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
		fmt.Fprintln(DebugWriter, args)
	}
	return q.executor.QueryRow(qs, args...)
}

// ExecQueryAll executes the query for the All finisher and returns multiple rows
func ExecQueryAll(q *Query) (*sql.Rows, error) {
	qs, args := buildQuery(q)
	if DebugMode {
		fmt.Fprintln(DebugWriter, qs)
		fmt.Fprintln(DebugWriter, args)
	}
	return q.executor.Query(qs, args...)
}

// SetExecutor on the query.
func SetExecutor(q *Query, exec Executor) {
	q.executor = exec
}

// GetExecutor on the query.
func GetExecutor(q *Query) Executor {
	return q.executor
}

// SetDialect on the query.
func SetDialect(q *Query, dialect *Dialect) {
	q.dialect = dialect
}

// SetSQL on the query.
func SetSQL(q *Query, sql string, args ...interface{}) {
	q.plainSQL = plainSQL{sql: sql, args: args}
}

// SetLoad on the query.
func SetLoad(q *Query, relationships ...string) {
	q.load = append([]string(nil), relationships...)
}

// SetCount on the query.
func SetCount(q *Query) {
	q.count = true
}

// SetDelete on the query.
func SetDelete(q *Query) {
	q.delete = true
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}

// SetFor on the query.
func SetFor(q *Query, clause string) {
	q.forlock = clause
}

// SetUpdate on the query.
func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

// AppendSelect on the query.
func AppendSelect(q *Query, columns ...string) {
	q.selectCols = append(q.selectCols, columns...)
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

// AppendHaving on the query.
func AppendHaving(q *Query, clause string, args ...interface{}) {
	q.having = append(q.having, having{clause: clause, args: args})
}

// AppendWhere on the query.
func AppendWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// AppendIn on the query.
func AppendIn(q *Query, clause string, args ...interface{}) {
	q.in = append(q.in, in{clause: clause, args: args})
}

// SetLastWhereAsOr sets the or separator for the tail "WHERE" in the slice
func SetLastWhereAsOr(q *Query) {
	if len(q.where) == 0 {
		return
	}

	q.where[len(q.where)-1].orSeparator = true
}

// SetLastInAsOr sets the or separator for the tail "IN" in the slice
func SetLastInAsOr(q *Query) {
	if len(q.in) == 0 {
		return
	}

	q.in[len(q.in)-1].orSeparator = true
}

// AppendGroupBy on the query.
func AppendGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// AppendOrderBy on the query.
func AppendOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}
