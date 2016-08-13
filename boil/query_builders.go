package boil

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/vattle/sqlboiler/strmangle"
)

var (
	rgxIdentifier = regexp.MustCompile(`^(?i)"?[a-z_][_a-z0-9]*"?(?:\."?[_a-z][_a-z0-9]*"?)*$`)
)

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

	defer strmangle.PutBuffer(buf)

	return buf.String(), args
}

func buildSelectQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := strmangle.GetBuffer()
	var args []interface{}

	buf.WriteString("SELECT ")

	// Wrap the select in the modifier function
	hasModFunc := len(q.modFunction) != 0
	if hasModFunc {
		fmt.Fprintf(buf, "%s(", q.modFunction)
	}

	hasSelectCols := len(q.selectCols) != 0
	hasJoins := len(q.joins) != 0
	if hasSelectCols && hasJoins && !hasModFunc {
		selectColsWithAs := writeAsStatements(q)
		// Don't identQuoteSlice - writeAsStatements does this
		buf.WriteString(strings.Join(selectColsWithAs, ", "))
	} else if hasSelectCols {
		buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.selectCols), ", "))
	} else if hasJoins {
		selectColsWithStars := writeStars(q)
		buf.WriteString(strings.Join(selectColsWithStars, ", "))
	} else {
		buf.WriteByte('*')
	}

	if hasModFunc {
		buf.WriteByte(')')
	}

	fmt.Fprintf(buf, " FROM %s", strings.Join(strmangle.IdentQuoteSlice(q.from), ", "))

	if len(q.joins) > 0 {
		argsLen := len(args)
		joinBuf := strmangle.GetBuffer()
		for _, j := range q.joins {
			if j.kind != JoinInner {
				panic("only inner joins are supported")
			}
			fmt.Fprintf(joinBuf, " INNER JOIN %s", j.clause)
			args = append(args, j.args...)
		}
		fmt.Fprintf(buf, convertQuestionMarks(joinBuf.String(), argsLen+1))
		strmangle.PutBuffer(joinBuf)
	}

	where, whereArgs := whereClause(q, len(args)+1)
	buf.WriteString(where)

	args = append(args, whereArgs...)
	writeModifiers(q, buf, &args)

	buf.WriteByte(';')
	return buf, args
}

func buildDeleteQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := strmangle.GetBuffer()

	buf.WriteString("DELETE FROM ")
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.from), ", "))

	where, args := whereClause(q, 1)
	buf.WriteString(where)

	writeModifiers(q, buf, &args)

	buf.WriteByte(';')

	return buf, args
}

func buildUpdateQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := strmangle.GetBuffer()

	buf.WriteString("UPDATE ")
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.from), ", "))

	cols := make(sort.StringSlice, len(q.update))
	var args []interface{}

	count := 0
	for name := range q.update {
		cols[count] = name
		count++
	}

	cols.Sort()

	for i := 0; i < len(cols); i++ {
		args = append(args, q.update[cols[i]])
		cols[i] = strmangle.IdentQuote(cols[i])
	}

	buf.WriteString(fmt.Sprintf(
		" SET (%s) = (%s)",
		strings.Join(cols, ", "),
		strmangle.Placeholders(len(cols), 1, 1)),
	)

	where, whereArgs := whereClause(q, len(args)+1)
	buf.WriteString(where)

	args = append(args, whereArgs...)
	writeModifiers(q, buf, &args)

	buf.WriteByte(';')

	return buf, args
}

func writeModifiers(q *Query, buf *bytes.Buffer, args *[]interface{}) {
	if len(q.groupBy) != 0 {
		fmt.Fprintf(buf, " GROUP BY %s", strings.Join(q.groupBy, ", "))
	}

	if len(q.having) != 0 {
		argsLen := len(*args)
		havingBuf := strmangle.GetBuffer()
		fmt.Fprintf(havingBuf, " HAVING ")
		for i, j := range q.having {
			if i > 0 {
				fmt.Fprintf(havingBuf, ", ")
			}
			fmt.Fprintf(havingBuf, j.clause)
			*args = append(*args, j.args...)
		}
		fmt.Fprintf(buf, convertQuestionMarks(havingBuf.String(), argsLen+1))
		strmangle.PutBuffer(havingBuf)
	}

	if len(q.orderBy) != 0 {
		buf.WriteString(" ORDER BY ")
		buf.WriteString(strings.Join(q.orderBy, ", "))
	}

	if q.limit != 0 {
		fmt.Fprintf(buf, " LIMIT %d", q.limit)
	}
	if q.offset != 0 {
		fmt.Fprintf(buf, " OFFSET %d", q.offset)
	}
}

func writeStars(q *Query) []string {
	cols := make([]string, 0, len(q.from))
	for _, f := range q.from {
		toks := strings.Split(f, " ")
		if len(toks) == 1 {
			cols = append(cols, fmt.Sprintf(`%s.*`, strmangle.IdentQuote(toks[0])))
			continue
		}

		alias, name, ok := parseFromClause(toks)
		if !ok {
		}

		if len(alias) != 0 {
			name = alias
		}
		cols = append(cols, fmt.Sprintf(`%s.*`, strmangle.IdentQuote(name)))
	}

	return cols
}

func writeAsStatements(q *Query) []string {
	cols := make([]string, len(q.selectCols))
	for i, col := range q.selectCols {
		if !rgxIdentifier.MatchString(col) {
			cols[i] = col
			continue
		}

		toks := strings.Split(col, ".")
		if len(toks) == 1 {
			cols[i] = strmangle.IdentQuote(col)
			continue
		}

		asParts := make([]string, len(toks))
		for j, tok := range toks {
			asParts[j] = strings.Trim(tok, `"`)
		}

		cols[i] = fmt.Sprintf(`%s as "%s"`, strmangle.IdentQuote(col), strings.Join(asParts, "."))
	}

	return cols
}

// whereClause parses a where slice and converts it into a
// single WHERE clause like:
// WHERE (a=$1) AND (b=$2)
//
// startAt specifies what number placeholders start at
func whereClause(q *Query, startAt int) (string, []interface{}) {
	if len(q.where) == 0 {
		return "", nil
	}

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)
	var args []interface{}

	buf.WriteString(" WHERE ")
	for i := 0; i < len(q.where); i++ {
		buf.WriteString(fmt.Sprintf("(%s)", q.where[i].clause))
		args = append(args, q.where[i].args...)

		// break on the last loop
		if i == len(q.where)-1 {
			break
		}

		if q.where[i].orSeparator {
			buf.WriteString(" OR ")
		} else {
			buf.WriteString(" AND ")
		}
	}

	return convertQuestionMarks(buf.String(), startAt), args
}

func convertQuestionMarks(clause string, startAt int) string {
	if startAt == 0 {
		panic("Not a valid start number.")
	}

	paramBuf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(paramBuf)
	paramIndex := 0

	for ; ; startAt++ {
		if paramIndex >= len(clause) {
			break
		}

		clause = clause[paramIndex:]
		paramIndex = strings.IndexByte(clause, '?')

		if paramIndex == -1 {
			paramBuf.WriteString(clause)
			break
		}

		paramBuf.WriteString(clause[:paramIndex] + fmt.Sprintf("$%d", startAt))
		paramIndex++
	}

	return paramBuf.String()
}

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

	for _, join := range q.joins {
		tokens := strings.Split(join.clause, " ")
		parseIdentifierClause(tokens, setID)
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
	alias, name, ok := parseFromClause(tokens)
	if !ok {
		panic("could not parse from statement")
	}

	if len(alias) > 0 {
		setID(alias, name)
	} else {
		setID(name, name)
	}
}

func parseFromClause(toks []string) (alias, name string, ok bool) {
	if len(toks) > 3 {
		toks = toks[:3]
	}

	sawIdent, sawAs := false, false
	for _, tok := range toks {
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
		ok = true
	}

	return alias, name, ok
}
