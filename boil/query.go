package boil

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"
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
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.from), ","))

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

var (
	rgxIdentifier      = regexp.MustCompile(`^(?i)"?[a-z_][_a-z0-9]*"?(?:\."?[_a-z][_a-z0-9]*"?)*$`)
	rgxJoinIdentifiers = regexp.MustCompile(`^(?i)(?:join|inner|natural|outer|left|right)$`)
)

// identifierMapping creates a map of all identifiers to potential model names
func identifierMapping(q *Query) map[string]string {
	var ids map[string]string
	setID := func(alias, name string) {
		if ids == nil {
			ids = make(map[string]string)
		}
		ids[alias] = name
	}

	for _, from := range q.from {
		tokens := strings.Split(from, " ")
		parseIdentifierClause(tokens, setID)
	}

	for _, join := range q.innerJoins {
		tokens := strings.Split(join.on, " ")
		discard := 0
		for rgxJoinIdentifiers.MatchString(tokens[discard]) {
			discard++
		}
		parseIdentifierClause(tokens[discard:], setID)
	}

	return ids
}

// parseBits takes a set of tokens and looks for something of the form:
// a b
// a as b
// where 'a' and 'b' are valid SQL identifiers
// It only evaluates the first 3 tokens (anything past that is superfluous)
// It stops parsing when it finds "on" or an invalid identifier
func parseIdentifierClause(tokens []string, setID func(string, string)) {
	var name, alias string
	sawIdent, sawAs := false, false

	if len(tokens) > 3 {
		tokens = tokens[:3]
	}

	for _, tok := range tokens {
		if t := strings.ToLower(tok); sawIdent && t == "as" {
			sawAs = true
			continue
		} else if sawIdent && t == "on" {
			break
		}

		if !rgxIdentifier.MatchString(tok) {
			break
		}

		if sawIdent || sawAs {
			alias = strings.Trim(tok, `"`)
			break
		}

		name = strings.Trim(tok, `"`)
		sawIdent = true
	}

	if len(alias) > 0 {
		setID(alias, name)
	} else {
		setID(name, name)
	}
}
