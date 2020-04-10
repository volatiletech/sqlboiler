// Package strmangle is a collection of string manipulation functions.
// Primarily used by boil and templates for code generation.
// Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	idAlphabet    = []byte("abcdefghijklmnopqrstuvwxyz")
	smartQuoteRgx = regexp.MustCompile(`^(?i)"?[a-z_][_a-z0-9\-]*"?(\."?[_a-z][_a-z0-9]*"?)*(\.\*)?$`)

	rgxEnum            = regexp.MustCompile(`^enum(\.[a-z0-9_]+)?\((,?'[^']+')+\)$`)
	rgxEnumIsOK        = regexp.MustCompile(`^(?i)[a-z][a-z0-9_\s]*$`)
	rgxEnumShouldTitle = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	rgxWhitespace      = regexp.MustCompile(`\s`)
)

var uppercaseWords = map[string]struct{}{
	"acl":   {},
	"api":   {},
	"ascii": {},
	"cpu":   {},
	"eof":   {},
	"guid":  {},
	"id":    {},
	"ip":    {},
	"json":  {},
	"ram":   {},
	"sla":   {},
	"udp":   {},
	"ui":    {},
	"uid":   {},
	"uuid":  {},
	"uri":   {},
	"url":   {},
	"utf8":  {},
}

var reservedWords = map[string]struct{}{
	"break":       {},
	"case":        {},
	"chan":        {},
	"const":       {},
	"continue":    {},
	"default":     {},
	"defer":       {},
	"else":        {},
	"fallthrough": {},
	"for":         {},
	"func":        {},
	"go":          {},
	"goto":        {},
	"if":          {},
	"import":      {},
	"interface":   {},
	"map":         {},
	"package":     {},
	"range":       {},
	"return":      {},
	"select":      {},
	"struct":      {},
	"switch":      {},
	"type":        {},
	"var":         {},
}

func init() {
	// Our Boil inflection Ruleset does not include uncountable inflections.
	// This way, people using words like Sheep will not have
	// collisions with their model name (Sheep) and their
	// function name (Sheep()). Instead, it will
	// use the regular inflection rules: Sheep, Sheeps().
	boilRuleset = newBoilRuleset()
}

// SchemaTable returns a table name with a schema prefixed if
// using a database that supports real schemas, for example,
// for Postgres: "schema_name"."table_name",
// for MS SQL: [schema_name].[table_name], versus
// simply "table_name" for MySQL (because it does not support real schemas)
func SchemaTable(lq, rq string, useSchema bool, schema string, table string) string {
	if useSchema {
		return fmt.Sprintf(`%s%s%s.%s%s%s`, lq, schema, rq, lq, table, rq)
	}

	return fmt.Sprintf(`%s%s%s`, lq, table, rq)
}

// IdentQuote attempts to quote simple identifiers in SQL statements
func IdentQuote(lq rune, rq rune, s string) string {
	if strings.EqualFold(s, "null") || s == "?" {
		return s
	}

	if m := smartQuoteRgx.MatchString(s); m != true {
		return s
	}

	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(s, ".")
	for i, split := range splits {
		if i != 0 {
			buf.WriteByte('.')
		}

		if rune(split[0]) == lq || rune(split[len(split)-1]) == rq || split == "*" {
			buf.WriteString(split)
			continue
		}

		buf.WriteRune(lq)
		buf.WriteString(split)
		buf.WriteRune(rq)
	}

	return buf.String()
}

// IdentQuoteSlice applies IdentQuote to a slice.
func IdentQuoteSlice(lq rune, rq rune, s []string) []string {
	if len(s) == 0 {
		return s
	}

	strs := make([]string, len(s))
	for i, str := range s {
		strs[i] = IdentQuote(lq, rq, str)
	}

	return strs
}

// Identifier is a base conversion from Base 10 integers to Base 26
// integers that are represented by an alphabet from a-z
// See tests for example outputs.
func Identifier(in int) string {
	ln := len(idAlphabet)
	var n int
	if in == 0 {
		n = 1
	} else {
		n = 1 + int(math.Log(float64(in))/math.Log(float64(ln)))
	}

	cols := GetBuffer()
	defer PutBuffer(cols)

	for i := 0; i < n; i++ {
		divisor := int(math.Pow(float64(ln), float64(n-i-1)))
		rem := in / divisor
		cols.WriteByte(idAlphabet[rem])

		in -= rem * divisor
	}

	return cols.String()
}

// QuoteCharacter returns a string that allows the quote character
// to be embedded into a Go string that uses double quotes:
func QuoteCharacter(q rune) string {
	if q == '"' {
		return `\"`
	}

	return string(q)
}

// Plural converts singular words to plural words (eg: person to people)
func Plural(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(name, "_")

	for i := 0; i < len(splits); i++ {
		if i != 0 {
			buf.WriteByte('_')
		}

		if i == len(splits)-1 {
			buf.WriteString(boilRuleset.Pluralize(splits[len(splits)-1]))
			break
		}

		buf.WriteString(splits[i])
	}

	return buf.String()
}

// Singular converts plural words to singular words (eg: people to person)
func Singular(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(name, "_")

	for i := 0; i < len(splits); i++ {
		if i != 0 {
			buf.WriteByte('_')
		}

		if i == len(splits)-1 {
			buf.WriteString(boilRuleset.Singularize(splits[len(splits)-1]))
			break
		}

		buf.WriteString(splits[i])
	}

	return buf.String()
}

// titleCaseCache holds the mapping of title cases.
// Example: map["MyWord"] == "my_word"
var (
	mut            sync.RWMutex
	titleCaseCache = map[string]string{}
)

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
//
// Note: This method is ugly because it has been highly optimized,
// we found that it was a fairly large bottleneck when we were using regexp.
func TitleCase(n string) string {
	// Attempt to fetch from cache
	mut.RLock()
	val, ok := titleCaseCache[n]
	mut.RUnlock()
	if ok {
		return val
	}

	ln := len(n)
	name := []byte(n)
	buf := GetBuffer()

	start := 0
	end := 0
	for start < ln {
		// Find the start and end of the underscores to account
		// for the possibility of being multiple underscores in a row.
		if end < ln {
			if name[start] == '_' {
				start++
				end++
				continue
				// Once we have found the end of the underscores, we can
				// find the end of the first full word.
			} else if name[end] != '_' {
				end++
				continue
			}
		}

		word := name[start:end]
		wordLen := len(word)
		var vowels bool

		numStart := wordLen
		for i, c := range word {
			vowels = vowels || (c == 97 || c == 101 || c == 105 || c == 111 || c == 117 || c == 121)

			if c > 47 && c < 58 && numStart == wordLen {
				numStart = i
			}
		}

		_, match := uppercaseWords[string(word[:numStart])]

		if match || !vowels {
			// Uppercase all a-z characters
			for _, c := range word {
				if c > 96 && c < 123 {
					buf.WriteByte(c - 32)
				} else {
					buf.WriteByte(c)
				}
			}
		} else {
			if c := word[0]; c > 96 && c < 123 {
				buf.WriteByte(word[0] - 32)
				buf.Write(word[1:])
			} else {
				buf.Write(word)
			}
		}

		start = end + 1
		end = start
	}

	ret := buf.String()
	PutBuffer(buf)

	// Cache the title case result
	mut.Lock()
	titleCaseCache[n] = ret
	mut.Unlock()

	return ret
}

// Ignore sets "-" for the tags values, so the fields will be ignored during parsing
func Ignore(table, column string, ignoreList map[string]struct{}) bool {
	_, ok := ignoreList[column]
	if ok {
		return true
	}
	_, ok = ignoreList[fmt.Sprintf("%s.%s", table, column)]
	if ok {
		return true
	}
	return false
}

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID". It will also lowercase the first letter
// of the name in the case where it's fed something that starts with uppercase.
func CamelCase(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	// Discard all leading '_'
	index := -1
	for i := 0; i < len(name); i++ {
		if name[i] != '_' {
			index = i
			break
		}
	}

	if index != -1 {
		name = name[index:]
	} else {
		return ""
	}

	index = -1
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			index = i
			break
		}
	}

	if index == -1 {
		buf.WriteString(strings.ToLower(string(name[0])))
		if len(name) > 1 {
			buf.WriteString(name[1:])
		}
	} else {
		buf.WriteString(strings.ToLower(string(name[0])))
		if len(name) > 1 {
			buf.WriteString(name[1:index])
			buf.WriteString(TitleCase(name[index+1:]))
		}
	}

	return buf.String()
}

// TitleCaseIdentifier splits on dots and then titlecases each fragment.
// map titleCase (split c ".")
func TitleCaseIdentifier(id string) string {
	nextDot := strings.IndexByte(id, '.')
	if nextDot < 0 {
		return TitleCase(id)
	}

	buf := GetBuffer()
	defer PutBuffer(buf)
	lastDot := 0
	ln := len(id)
	addDots := false

	for i := 0; nextDot >= 0; i++ {
		fragment := id[lastDot:nextDot]

		titled := TitleCase(fragment)

		if addDots {
			buf.WriteByte('.')
		}
		buf.WriteString(titled)
		addDots = true

		if nextDot == ln {
			break
		}

		lastDot = nextDot + 1
		if nextDot = strings.IndexByte(id[lastDot:], '.'); nextDot >= 0 {
			nextDot += lastDot
		} else {
			nextDot = ln
		}
	}

	return buf.String()
}

// MakeStringMap converts a map[string]string into the format:
// "key": "value", "key": "value"
func MakeStringMap(types map[string]string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	keys := make([]string, 0, len(types))
	for k := range types {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	c := 0
	for _, k := range keys {
		v := types[k]
		buf.WriteString(fmt.Sprintf("`%s`: `%s`", k, v))
		if c < len(types)-1 {
			buf.WriteString(", ")
		}

		c++
	}

	return buf.String()
}

// StringMap maps a function over a slice of strings.
func StringMap(modifier func(string) string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, str := range strs {
		ret[i] = modifier(str)
	}

	return ret
}

// PrefixStringSlice with the given str.
func PrefixStringSlice(str string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, s := range strs {
		ret[i] = fmt.Sprintf("%s%s", str, s)
	}

	return ret
}

// Placeholders generates the SQL statement placeholders for in queries.
// For example, ($1,$2,$3),($4,$5,$6) etc.
// It will start counting placeholders at "start".
// If useIndexPlaceholders is false, it will convert to ? instead of $1 etc.
func Placeholders(useIndexPlaceholders bool, count int, start int, group int) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	if start == 0 || group == 0 {
		panic("Invalid start or group numbers supplied.")
	}

	if group > 1 {
		buf.WriteByte('(')
	}
	for i := 0; i < count; i++ {
		if i != 0 {
			if group > 1 && i%group == 0 {
				buf.WriteString("),(")
			} else {
				buf.WriteByte(',')
			}
		}
		if useIndexPlaceholders {
			buf.WriteString(fmt.Sprintf("$%d", start+i))
		} else {
			buf.WriteByte('?')
		}
	}
	if group > 1 {
		buf.WriteByte(')')
	}

	return buf.String()
}

// SetParamNames takes a slice of columns and returns a comma separated
// list of parameter names for a template statement SET clause.
// eg: "col1"=$1, "col2"=$2, "col3"=$3
func SetParamNames(lq, rq string, start int, columns []string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	for i, c := range columns {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s%s%s=$%d`, lq, c, rq, i+start))
		} else {
			buf.WriteString(fmt.Sprintf(`%s%s%s=?`, lq, c, rq))
		}

		if i < len(columns)-1 {
			buf.WriteByte(',')
		}
	}

	return buf.String()
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(lq, rq string, start int, cols []string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	for i, c := range cols {
		if start != 0 {
			buf.WriteString(fmt.Sprintf(`%s%s%s=$%d`, lq, c, rq, start+i))
		} else {
			buf.WriteString(fmt.Sprintf(`%s%s%s=?`, lq, c, rq))
		}

		if i < len(cols)-1 {
			buf.WriteString(" AND ")
		}
	}

	return buf.String()
}

// WhereClauseRepeated returns the where clause repeated with OR clause using start as the $ flag index
// For example, if start was 2 output would be: "(colthing=$2 AND colstuff=$3) OR (colthing=$4 AND colstuff=$5)"
func WhereClauseRepeated(lq, rq string, start int, cols []string, count int) string {
	var startIndex int
	buf := GetBuffer()
	defer PutBuffer(buf)
	buf.WriteByte('(')
	for i := 0; i < count; i++ {
		if i != 0 {
			buf.WriteString(") OR (")
		}

		startIndex = 0
		if start > 0 {
			startIndex = start + i*len(cols)
		}

		buf.WriteString(WhereClause(lq, rq, startIndex, cols))
	}
	buf.WriteByte(')')

	return buf.String()
}

// JoinSlices merges two string slices of equal length
func JoinSlices(sep string, a, b []string) []string {
	lna, lnb := len(a), len(b)
	if lna != lnb {
		panic("joinSlices: can only merge slices of same length")
	} else if lna == 0 {
		return nil
	}

	ret := make([]string, len(a))
	for i, elem := range a {
		ret[i] = fmt.Sprintf("%s%s%s", elem, sep, b[i])
	}

	return ret
}

// StringSliceMatch returns true if the length of both
// slices is the same, and the elements of both slices are the same.
// The elements can be in any order.
func StringSliceMatch(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, aval := range a {
		found := false
		for _, bval := range b {
			if bval == aval {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// ContainsAny returns true if any of the passed in strings are
// found in the passed in string slice
func ContainsAny(a []string, finds ...string) bool {
	for _, s := range a {
		for _, find := range finds {
			if s == find {
				return true
			}
		}
	}

	return false
}

// GenerateTags converts a slice of tag strings into tags that
// can be passed onto the end of a struct, for example:
// tags: ["xml", "db"] convert to: xml:"column_name" db:"column_name"
func GenerateTags(tags []string, columnName string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	for _, tag := range tags {
		buf.WriteString(tag)
		buf.WriteString(`:"`)
		buf.WriteString(columnName)
		buf.WriteString(`" `)
	}

	return buf.String()
}

// GenerateIgnoreTags converts a slice of tag strings into
// ignore tags that can be passed onto the end of a struct, for example:
// tags: ["xml", "db"] convert to: xml:"-" db:"-"
func GenerateIgnoreTags(tags []string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	for _, tag := range tags {
		buf.WriteString(tag)
		buf.WriteString(`:"-" `)
	}

	return buf.String()
}

// ParseEnumVals returns the values from an enum string
//
// Postgres and MySQL drivers return different values
// psql:  enum.enum_name('values'...)
// mysql: enum('values'...)
func ParseEnumVals(s string) []string {
	if !rgxEnum.MatchString(s) {
		return nil
	}

	startIndex := strings.IndexByte(s, '(')
	s = s[startIndex+2 : len(s)-2]
	return strings.Split(s, "','")
}

// ParseEnumName returns the name portion of an enum if it exists
//
// Postgres and MySQL drivers return different values
// psql:  enum.enum_name('values'...)
// mysql: enum('values'...)
// In the case of mysql, the name will never return anything
func ParseEnumName(s string) string {
	if !rgxEnum.MatchString(s) {
		return ""
	}

	endIndex := strings.IndexByte(s, '(')
	s = s[:endIndex]
	startIndex := strings.IndexByte(s, '.')
	if startIndex < 0 {
		return ""
	}

	return s[startIndex+1:]
}

// IsEnumNormal checks a set of eval values to see if they're "normal"
func IsEnumNormal(values []string) bool {
	for _, v := range values {
		if !rgxEnumIsOK.MatchString(v) {
			return false
		}
	}

	return true
}

//StripWhitespace removes all whitespace from a string
func StripWhitespace(value string) string {
	return rgxWhitespace.ReplaceAllString(value, "")
}

// ShouldTitleCaseEnum checks a value to see if it's title-case-able
func ShouldTitleCaseEnum(value string) bool {
	return rgxEnumShouldTitle.MatchString(value)
}

// ReplaceReservedWords takes a word and replaces it with word_ if it's found
// in the list of reserved words.
func ReplaceReservedWords(word string) string {
	if _, ok := reservedWords[word]; ok {
		return word + "_"
	}
	return word
}

// RemoveDuplicates from a string slice
func RemoveDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}
