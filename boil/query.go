package boil

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"

	"github.com/nullbio/sqlboiler/strmangle"
)

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

// Query holds the state for the built up query
type Query struct {
	executor        Executor
	plainSQL        plainSQL
	delete          bool
	update          map[string]interface{}
	selectCols      []string
	count           bool
	from            []string
	innerJoins      []join
	outerJoins      []join
	leftOuterJoins  []join
	rightOuterJoins []join
	where           []where
	groupBy         []string
	orderBy         []string
	having          []string
	limit           int
	offset          int
}

func buildQuery(q *Query) (string, []interface{}) {
	var buf *bytes.Buffer
	var args []interface{}

	switch {
	case len(q.plainSQL.sql) != 0:
		return q.plainSQL.sql, q.plainSQL.args
	case q.delete:
		buf, args = buildDeleteQuery(q)
	case len(q.update) > 0:
		buf, args = buildUpdateQuery(q)
	default:
		buf, args = buildSelectQuery(q)
	}

	return buf.String(), args
}

func buildSelectQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	buf.WriteString("SELECT ")

	if q.count {
		buf.WriteString("COUNT(")
	}
	if len(q.selectCols) > 0 {
		buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.selectCols), `, `))
	} else {
		buf.WriteByte('*')
	}
	// close sql COUNT function
	if q.count {
		buf.WriteString(")")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.from), ","))

	where, args := whereClause(q)
	buf.WriteString(where)

	if len(q.orderBy) != 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(q.orderBy, ","))
	}

	if q.limit != 0 {
		fmt.Fprintf(buf, " LIMIT %d", q.limit)
	}
	if q.offset != 0 {
		fmt.Fprintf(buf, " OFFSET %d", q.offset)
	}

	buf.WriteByte(';')
	return buf, args
}

func buildDeleteQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	buf.WriteString("DELETE FROM ")
	fmt.Fprintf(buf, `"%s"`, q.from)

	where, args := whereClause(q)
	buf.WriteString(where)

	buf.WriteByte(';')

	return buf, args
}

func buildUpdateQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	buf.WriteByte(';')
	return buf, nil
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

// SetSelect on the query.
func SetSelect(q *Query, columns ...string) {
	cols := make([]string, len(columns))
	copy(cols, columns)

	q.selectCols = append(q.selectCols, cols...)
}

// Select returns the select columns in the query.
func Select(q *Query) []string {
	cols := make([]string, len(q.selectCols))
	copy(cols, q.selectCols)
	return cols
}

// SetFrom on the query.
func SetFrom(q *Query, from ...string) {
	q.from = append(q.from, from...)
}

// SetInnerJoin on the query.
func SetInnerJoin(q *Query, on string, args ...interface{}) {
	q.innerJoins = append(q.innerJoins, join{on: on, args: args})
}

// SetOuterJoin on the query.
func SetOuterJoin(q *Query, on string, args ...interface{}) {
	q.outerJoins = append(q.outerJoins, join{on: on, args: args})
}

// SetLeftOuterJoin on the query.
func SetLeftOuterJoin(q *Query, on string, args ...interface{}) {
	q.leftOuterJoins = append(q.leftOuterJoins, join{on: on, args: args})
}

// SetRightOuterJoin on the query.
func SetRightOuterJoin(q *Query, on string, args ...interface{}) {
	q.rightOuterJoins = append(q.rightOuterJoins, join{on: on, args: args})
}

// SetWhere on the query.
func SetWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

// SetLastWhereAsOr sets the or seperator for the last element in the where slice
func SetLastWhereAsOr(q *Query) {
	q.where[len(q.where)-1].orSeperator = true
}

// SetGroupBy on the query.
func SetGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

// SetOrderBy on the query.
func SetOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}

// SetHaving on the query.
func SetHaving(q *Query, clause string) {
	q.having = append(q.having, clause)
}

// SetLimit on the query.
func SetLimit(q *Query, limit int) {
	q.limit = limit
}

// SetOffset on the query.
func SetOffset(q *Query, offset int) {
	q.offset = offset
}

func whereClause(q *Query) (string, []interface{}) {
	if len(q.where) == 0 {
		return "", nil
	}

	buf := &bytes.Buffer{}
	var args []interface{}

	buf.WriteString(" WHERE ")
	for i := 0; i < len(q.where); i++ {
		buf.WriteString(fmt.Sprintf("%s", q.where[i].clause))
		args = append(args, q.where[i].args...)
		if i >= len(q.where)-1 {
			continue
		}
		if q.where[i].orSeperator {
			buf.WriteString(" OR ")
		} else {
			buf.WriteString(" AND ")
		}
	}

	return buf.String(), args
}
