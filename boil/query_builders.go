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
	rgxInClause   = regexp.MustCompile(`^(?i)(.*[\s|\)|\?])IN([\s|\(|\?].*)$`)
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
		resp, _ := convertQuestionMarks(joinBuf.String(), argsLen+1)
		fmt.Fprintf(buf, resp)
		strmangle.PutBuffer(joinBuf)
	}

	where, whereArgs := whereClause(q, len(args)+1)
	buf.WriteString(where)
	if len(whereArgs) != 0 {
		args = append(args, whereArgs...)
	}

	in, inArgs := inClause(q, len(args)+1)
	buf.WriteString(in)
	if len(inArgs) != 0 {
		args = append(args, inArgs...)
	}

	writeModifiers(q, buf, &args)

	buf.WriteByte(';')
	return buf, args
}

func buildDeleteQuery(q *Query) (*bytes.Buffer, []interface{}) {
	var args []interface{}
	buf := strmangle.GetBuffer()

	buf.WriteString("DELETE FROM ")
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.from), ", "))

	where, whereArgs := whereClause(q, 1)
	if len(whereArgs) != 0 {
		args = append(args, whereArgs)
	}
	buf.WriteString(where)

	in, inArgs := inClause(q, len(args)+1)
	if len(inArgs) != 0 {
		args = append(args, inArgs...)
	}
	buf.WriteString(in)

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
	if len(whereArgs) != 0 {
		args = append(args, whereArgs...)
	}
	buf.WriteString(where)

	in, inArgs := inClause(q, len(args)+1)
	if len(inArgs) != 0 {
		args = append(args, inArgs...)
	}
	buf.WriteString(in)

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
		resp, _ := convertQuestionMarks(havingBuf.String(), argsLen+1)
		fmt.Fprintf(buf, resp)
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
	for i, where := range q.where {
		if i != 0 {
			if where.orSeparator {
				buf.WriteString(" OR ")
			} else {
				buf.WriteString(" AND ")
			}
		}

		buf.WriteString(fmt.Sprintf("(%s)", where.clause))
		args = append(args, where.args...)
	}

	resp, _ := convertQuestionMarks(buf.String(), startAt)
	return resp, args
}

// inClause parses an in slice and converts it into a
// single IN clause, like:
// WHERE ("a", "b") IN (($1,$2),($3,$4)).
func inClause(q *Query, startAt int) (string, []interface{}) {
	if len(q.in) == 0 {
		return "", nil
	}

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)
	var args []interface{}

	if len(q.where) == 0 {
		buf.WriteString(" WHERE ")
	}

	for i, in := range q.in {
		ln := len(in.args)
		// We only prefix the OR and AND separators after the first
		// clause has been generated UNLESS there is already a where
		// clause that we have to add on to.
		if i != 0 || len(q.where) > 0 {
			if in.orSeparator {
				buf.WriteString(" OR ")
			} else {
				buf.WriteString(" AND ")
			}
		}

		matches := rgxInClause.FindStringSubmatch(in.clause)
		// If we can't find any matches attempt a simple replace with 1 group.
		// Clauses that fit this criteria will not be able to contain ? in their
		// column name side, however if this case is being hit then the regexp
		// probably needs adjustment, or the user is passing in invalid clauses.
		if matches == nil {
			clause, count := convertInQuestionMarks(in.clause, startAt, 1, ln)
			buf.WriteString(clause)
			startAt = startAt + count
		} else {
			leftSide := strings.TrimSpace(matches[1])
			rightSide := strings.TrimSpace(matches[2])
			// If matches are found, we have to parse the left side (column side)
			// of the clause to determine how many columns they are using.
			// This number determines the groupAt for the convert function.
			cols := strings.Split(leftSide, ",")
			cols = strmangle.IdentQuoteSlice(cols)
			groupAt := len(cols)

			leftClause, leftCount := convertQuestionMarks(strings.Join(cols, ","), startAt)
			rightClause, rightCount := convertInQuestionMarks(rightSide, startAt+leftCount, groupAt, ln-leftCount)
			buf.WriteString(leftClause)
			buf.WriteString(" IN ")
			buf.WriteString(rightClause)
			startAt = startAt + leftCount + rightCount
		}

		args = append(args, in.args...)
	}

	return buf.String(), args
}

// convertInQuestionMarks finds the first unescaped occurence of ? and swaps it
// with a list of numbered placeholders, starting at startAt.
// It uses groupAt to determine how many placeholders should be in each group,
// for example, groupAt 2 would result in: (($1,$2),($3,$4))
// and groupAt 1 would result in ($1,$2,$3,$4)
func convertInQuestionMarks(clause string, startAt, groupAt, total int) (string, int) {
	if startAt == 0 || len(clause) == 0 {
		panic("Not a valid start number.")
	}

	paramBuf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(paramBuf)

	foundAt := -1
	for i := 0; i < len(clause); i++ {
		if (clause[i] == '?' && i == 0) || (clause[i] == '?' && clause[i-1] != '\\') {
			foundAt = i
			break
		}
	}

	if foundAt == -1 {
		return strings.Replace(clause, `\?`, "?", -1), 0
	}

	paramBuf.WriteString(clause[:foundAt])
	paramBuf.WriteByte('(')
	paramBuf.WriteString(strmangle.Placeholders(total, startAt, groupAt))
	paramBuf.WriteByte(')')
	paramBuf.WriteString(clause[foundAt+1:])

	// Remove all backslashes from escaped question-marks
	ret := strings.Replace(paramBuf.String(), `\?`, "?", -1)
	return ret, total
}

// convertQuestionMarks converts each occurence of ? with $<number>
// where <number> is an incrementing digit starting at startAt.
// If question-mark (?) is escaped using back-slash (\), it will be ignored.
func convertQuestionMarks(clause string, startAt int) (string, int) {
	if startAt == 0 {
		panic("Not a valid start number.")
	}

	paramBuf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(paramBuf)
	paramIndex := 0
	total := 0

	for {
		if paramIndex >= len(clause) {
			break
		}

		clause = clause[paramIndex:]
		paramIndex = strings.IndexByte(clause, '?')

		if paramIndex == -1 {
			paramBuf.WriteString(clause)
			break
		}

		escapeIndex := strings.Index(clause, `\?`)
		if escapeIndex != -1 && paramIndex > escapeIndex {
			paramBuf.WriteString(clause[:escapeIndex] + "?")
			paramIndex++
			continue
		}

		paramBuf.WriteString(clause[:paramIndex] + fmt.Sprintf("$%d", startAt))
		total++
		startAt++
		paramIndex++
	}

	return paramBuf.String(), total
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

// parseIdentifierClause takes a set of tokens and looks for something of the form:
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
