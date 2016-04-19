package boil

import (
	"bytes"
	"database/sql"
	"fmt"
	"strings"
)

type where struct {
	clause string
	args   []interface{}
}

type join struct {
	on   string
	args []interface{}
}

type Query struct {
	executor        Executor
	delete          bool
	update          map[string]interface{}
	selectCols      []string
	from            string
	innerJoins      []join
	outerJoins      []join
	leftOuterJoins  []join
	rightOuterJoins []join
	where           []where
	groupBy         []string
	orderBy         []string
	having          []string
	limit           int
}

func buildQuery(q *Query) (string, []interface{}) {
	var buf *bytes.Buffer
	var args []interface{}

	switch {
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

	if len(q.selectCols) > 0 {
		buf.WriteString(strings.Join(q.selectCols, ","))
	} else {
		buf.WriteByte('*')
	}

	buf.WriteString(" FROM ")
	fmt.Fprintf(buf, `"%s"`, q.from)

	return buf, []interface{}{}
}

func buildDeleteQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	return buf, nil
}

func buildUpdateQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	return buf, nil
}

func ExecQuery(q *Query) error {
	return nil
}

func ExecQueryOne(q *Query) (*sql.Row, error) {
	return nil, nil
}

func ExecQueryAll(q *Query) (*sql.Rows, error) {
	return nil, nil
}

func Apply(q *Query, mods ...func(q *Query)) {
	for _, mod := range mods {
		mod(q)
	}
}

func SetDelete(q *Query, flag bool) {
	q.delete = flag
}

func SetUpdate(q *Query, cols map[string]interface{}) {
	q.update = cols
}

func SetExecutor(q *Query, exec Executor) {
	q.executor = exec
}

func SetSelect(q *Query, columns ...string) {
	q.selectCols = append(q.selectCols, columns...)
}

func SetFrom(q *Query, table string) {
	q.from = table
}

func SetInnerJoin(q *Query, on string, args ...interface{}) {
	q.innerJoins = append(q.innerJoins, join{on: on, args: args})
}

func SetOuterJoin(q *Query, on string, args ...interface{}) {
	q.outerJoins = append(q.outerJoins, join{on: on, args: args})
}

func SetLeftOuterJoin(q *Query, on string, args ...interface{}) {
	q.leftOuterJoins = append(q.leftOuterJoins, join{on: on, args: args})
}

func SetRightOuterJoin(q *Query, on string, args ...interface{}) {
	q.rightOuterJoins = append(q.rightOuterJoins, join{on: on, args: args})
}

func SetWhere(q *Query, clause string, args ...interface{}) {
	q.where = append(q.where, where{clause: clause, args: args})
}

func SetGroupBy(q *Query, clause string) {
	q.groupBy = append(q.groupBy, clause)
}

func SetOrderBy(q *Query, clause string) {
	q.orderBy = append(q.orderBy, clause)
}

func SetHaving(q *Query, clause string) {
	q.having = append(q.having, clause)
}

func SetLimit(q *Query, limit int) {
	q.limit = limit
}
