package boil

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/nullbio/sqlboiler/strmangle"
)

var (
	rgxIdentifier      = regexp.MustCompile(`^(?i)"?[a-z_][_a-z0-9]*"?(?:\."?[_a-z][_a-z0-9]*"?)*$`)
	rgxJoinIdentifiers = regexp.MustCompile(`^(?i)(?:join|inner|natural|outer|left|right)$`)
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

	return buf.String(), args
}

func buildSelectQuery(q *Query) (*bytes.Buffer, []interface{}) {
	buf := &bytes.Buffer{}

	buf.WriteString("SELECT ")

	// Wrap the select in the modifier function
	hasModFunc := len(q.modFunction) != 0
	if hasModFunc {
		fmt.Fprintf(buf, "%s(", q.modFunction)
	}

	hasSelectCols := len(q.selectCols) != 0
	if len(q.innerJoins) != 0 && hasSelectCols && !hasModFunc {
		writeComplexSelect(q, buf)
	} else if hasSelectCols {
		buf.WriteString(strings.Join(strmangle.IdentQuoteSlice(q.selectCols), `, `))
	} else {
		buf.WriteByte('*')
	}

	if hasModFunc {
		buf.WriteString(")")
	}

	fmt.Fprintf(buf, " FROM %s", strings.Join(strmangle.IdentQuoteSlice(q.from), ","))

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

func writeComplexSelect(q *Query, buf *bytes.Buffer) {
	cols := make([]string, len(q.selectCols))
	for _, col := range q.selectCols {
		if !rgxIdentifier.Match {
			cols
		}
	}
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
