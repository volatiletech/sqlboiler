package queries

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
	case len(q.rawSQL.sql) != 0:
		return q.rawSQL.sql, q.rawSQL.args
	case q.delete:
		buf, args = buildDeleteQuery(q)
	case len(q.update) > 0:
		buf, args = buildUpdateQuery(q)
	default:
		buf, args = buildSelectQuery(q)
	}

	defer strmangle.PutBuffer(buf)

	// Cache the generated query for query object re-use
	bufStr := buf.String()
	q.rawSQL.sql = bufStr
	q.rawSQL.args = args

	return bufStr, args
}

func buildSelectQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := strmangle.GetBuffer()
	var args []interface{}

	buf.WriteString("SELECT ")

	if q.count {
		buf.WriteString("COUNT(")
	}

	hasSelectCols := len(q.selectCols) != 0
	hasJoins := len(q.joins) != 0
	if hasJoins && hasSelectCols && !q.count {
		selectColsWithAs := writeAsStatements(q)
		// Don't identQuoteSlice - writeAsStatements does this
		buf.WriteString(strings.Join(selectColsWithAs, ", "))
	} else if hasSelectCols {
		buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.dialect.LQ, q.dialect.RQ, q.selectCols), ", "))
	} else if hasJoins && !q.count {
		selectColsWithStars := writeStars(q)
		buf.WriteString(strings.Join(selectColsWithStars, ", "))
	} else {
		buf.WriteByte('*')
	}

	// close SQL COUNT function
	if q.count {
		buf.WriteByte(')')
	}

	fmt.Fprintf(buf, " FROM %s", strings.Join(strmangle.IdentQuoteSlice(q.dialect.LQ, q.dialect.RQ, q.from), ", "))

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
		var resp string
		if q.dialect.IndexPlaceholders {
			resp, _ = convertQuestionMarks(joinBuf.String(), argsLen+1)
		} else {
			resp = joinBuf.String()
		}
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
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.dialect.LQ, q.dialect.RQ, q.from), ", "))

	where, whereArgs := whereClause(q, 1)
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

func buildUpdateQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := strmangle.GetBuffer()

	buf.WriteString("UPDATE ")
	buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.dialect.LQ, q.dialect.RQ, q.from), ", "))

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
		cols[i] = strmangle.IdentQuote(q.dialect.LQ, q.dialect.RQ, cols[i])
	}

	buf.WriteString(fmt.Sprintf(
		" SET (%s) = (%s)",
		strings.Join(cols, ", "),
		strmangle.Placeholders(q.dialect.IndexPlaceholders, len(cols), 1, 1)),
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

// BuildUpsertQueryMySQL builds a SQL statement string using the upsertData provided.
func BuildUpsertQueryMySQL(dia Dialect, tableName string, update, whitelist []string) string {
	whitelist = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, whitelist)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	if len(update) == 0 {
		fmt.Fprintf(
			buf,
			"INSERT IGNORE INTO %s (%s) VALUES (%s)",
			tableName,
			strings.Join(whitelist, ", "),
			strmangle.Placeholders(dia.IndexPlaceholders, len(whitelist), 1, 1),
		)
		return buf.String()
	}

	fmt.Fprintf(
		buf,
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE ",
		tableName,
		strings.Join(whitelist, ", "),
		strmangle.Placeholders(dia.IndexPlaceholders, len(whitelist), 1, 1),
	)

	for i, v := range update {
		if i != 0 {
			buf.WriteByte(',')
		}
		quoted := strmangle.IdentQuote(dia.LQ, dia.RQ, v)
		buf.WriteString(quoted)
		buf.WriteString(" = VALUES(")
		buf.WriteString(quoted)
		buf.WriteByte(')')
	}

	return buf.String()
}

// BuildUpsertQueryPostgres builds a SQL statement string using the upsertData provided.
func BuildUpsertQueryPostgres(dia Dialect, tableName string, updateOnConflict bool, ret, update, conflict, whitelist []string) string {
	conflict = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, conflict)
	whitelist = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, whitelist)
	ret = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, ret)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	fmt.Fprintf(
		buf,
		"INSERT INTO %s (%s) VALUES (%s) ON CONFLICT ",
		tableName,
		strings.Join(whitelist, ", "),
		strmangle.Placeholders(dia.IndexPlaceholders, len(whitelist), 1, 1),
	)

	if !updateOnConflict || len(update) == 0 {
		buf.WriteString("DO NOTHING")
	} else {
		buf.WriteByte('(')
		buf.WriteString(strings.Join(conflict, ", "))
		buf.WriteString(") DO UPDATE SET ")

		for i, v := range update {
			if i != 0 {
				buf.WriteByte(',')
			}
			quoted := strmangle.IdentQuote(dia.LQ, dia.RQ, v)
			buf.WriteString(quoted)
			buf.WriteString(" = EXCLUDED.")
			buf.WriteString(quoted)
		}
	}

	if len(ret) != 0 {
		buf.WriteString(" RETURNING ")
		buf.WriteString(strings.Join(ret, ", "))
	}

	return buf.String()
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
		var resp string
		if q.dialect.IndexPlaceholders {
			resp, _ = convertQuestionMarks(havingBuf.String(), argsLen+1)
		} else {
			resp = havingBuf.String()
		}
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

	if len(q.forlock) != 0 {
		fmt.Fprintf(buf, " FOR %s", q.forlock)
	}
}

func writeStars(q *Query) []string {
	cols := make([]string, len(q.from))
	for i, f := range q.from {
		toks := strings.Split(f, " ")
		if len(toks) == 1 {
			cols[i] = fmt.Sprintf(`%s.*`, strmangle.IdentQuote(q.dialect.LQ, q.dialect.RQ, toks[0]))
			continue
		}

		alias, name, ok := parseFromClause(toks)
		if !ok {
			return nil
		}

		if len(alias) != 0 {
			name = alias
		}
		cols[i] = fmt.Sprintf(`%s.*`, strmangle.IdentQuote(q.dialect.LQ, q.dialect.RQ, name))
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
			cols[i] = strmangle.IdentQuote(q.dialect.LQ, q.dialect.RQ, col)
			continue
		}

		asParts := make([]string, len(toks))
		for j, tok := range toks {
			asParts[j] = strings.Trim(tok, `"`)
		}

		cols[i] = fmt.Sprintf(`%s as "%s"`, strmangle.IdentQuote(q.dialect.LQ, q.dialect.RQ, col), strings.Join(asParts, "."))
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

	var resp string
	if q.dialect.IndexPlaceholders {
		resp, _ = convertQuestionMarks(buf.String(), startAt)
	} else {
		resp = buf.String()
	}

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
			clause, count := convertInQuestionMarks(q.dialect.IndexPlaceholders, in.clause, startAt, 1, ln)
			buf.WriteString(clause)
			startAt = startAt + count
		} else {
			leftSide := strings.TrimSpace(matches[1])
			rightSide := strings.TrimSpace(matches[2])
			// If matches are found, we have to parse the left side (column side)
			// of the clause to determine how many columns they are using.
			// This number determines the groupAt for the convert function.
			cols := strings.Split(leftSide, ",")
			cols = strmangle.IdentQuoteSlice(q.dialect.LQ, q.dialect.RQ, cols)
			groupAt := len(cols)

			var leftClause string
			var leftCount int
			if q.dialect.IndexPlaceholders {
				leftClause, leftCount = convertQuestionMarks(strings.Join(cols, ","), startAt)
			} else {
				// Count the number of cols that are question marks, so we know
				// how much to offset convertInQuestionMarks by
				for _, v := range cols {
					if v == "?" {
						leftCount++
					}
				}
				leftClause = strings.Join(cols, ",")
			}
			rightClause, rightCount := convertInQuestionMarks(q.dialect.IndexPlaceholders, rightSide, startAt+leftCount, groupAt, ln-leftCount)
			buf.WriteString(leftClause)
			buf.WriteString(" IN ")
			buf.WriteString(rightClause)
			startAt = startAt + leftCount + rightCount
		}

		args = append(args, in.args...)
	}

	return buf.String(), args
}

// convertInQuestionMarks finds the first unescaped occurrence of ? and swaps it
// with a list of numbered placeholders, starting at startAt.
// It uses groupAt to determine how many placeholders should be in each group,
// for example, groupAt 2 would result in: (($1,$2),($3,$4))
// and groupAt 1 would result in ($1,$2,$3,$4)
func convertInQuestionMarks(indexPlaceholders bool, clause string, startAt, groupAt, total int) (string, int) {
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
	paramBuf.WriteString(strmangle.Placeholders(indexPlaceholders, total, startAt, groupAt))
	paramBuf.WriteByte(')')
	paramBuf.WriteString(clause[foundAt+1:])

	// Remove all backslashes from escaped question-marks
	ret := strings.Replace(paramBuf.String(), `\?`, "?", -1)
	return ret, total
}

// convertQuestionMarks converts each occurrence of ? with $<number>
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

// parseFromClause will parse something that looks like
// a
// a b
// a as b
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
