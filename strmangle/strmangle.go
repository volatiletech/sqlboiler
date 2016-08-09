// Package strmangle is used exclusively by the templates in sqlboiler.
// There are many helper functions to deal with bdb.* values as well
// as string manipulation. Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

var (
	idAlphabet     = []byte("abcdefghijklmnopqrstuvwxyz")
	uppercaseWords = regexp.MustCompile(`^(?i)(id|uid|uuid|guid|ssn|tz)[0-9]*$`)
	smartQuoteRgx  = regexp.MustCompile(`^(?i)"?[a-z_][_a-z0-9]*"?(\."?[_a-z][_a-z0-9]*"?)*(\.\*)?$`)
)

// IdentQuote attempts to quote simple identifiers in SQL tatements
func IdentQuote(s string) string {
	if strings.ToLower(s) == "null" {
		return s
	}

	if m := smartQuoteRgx.MatchString(s); m != true {
		return s
	}

	splits := strings.Split(s, ".")
	for i, split := range splits {
		if strings.HasPrefix(split, `"`) || strings.HasSuffix(split, `"`) || split == "*" {
			continue
		}

		splits[i] = fmt.Sprintf(`"%s"`, split)
	}

	return strings.Join(splits, ".")
}

// IdentQuoteSlice applies IdentQuote to a slice.
func IdentQuoteSlice(s []string) []string {
	if len(s) == 0 {
		return s
	}

	strs := make([]string, len(s))
	for i, str := range s {
		strs[i] = IdentQuote(str)
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

	cols := make([]byte, n)

	for i := 0; i < n; i++ {
		divisor := int(math.Pow(float64(ln), float64(n-i-1)))
		rem := in / divisor
		cols[i] = idAlphabet[rem]

		in -= rem * divisor
	}

	return string(cols)
}

// Plural converts singular words to plural words (eg: person to people)
func Plural(name string) string {
	splits := strings.Split(name, "_")
	splits[len(splits)-1] = inflection.Plural(splits[len(splits)-1])
	return strings.Join(splits, "_")
}

// Singular converts plural words to singular words (eg: people to person)
func Singular(name string) string {
	splits := strings.Split(name, "_")
	splits[len(splits)-1] = inflection.Singular(splits[len(splits)-1])
	return strings.Join(splits, "_")
}

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
func TitleCase(name string) string {
	splits := strings.Split(name, "_")

	for i, split := range splits {
		if uppercaseWords.MatchString(split) {
			splits[i] = strings.ToUpper(split)
			continue
		}

		splits[i] = strings.Title(split)
	}

	return strings.Join(splits, "")
}

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func CamelCase(name string) string {
	splits := strings.Split(name, "_")

	for i, split := range splits {
		if i == 0 {
			continue
		}

		if i > 0 {
			if uppercaseWords.MatchString(split) {
				splits[i] = strings.ToUpper(split)
				continue
			}
		}

		splits[i] = strings.Title(split)
	}

	return strings.Join(splits, "")
}

// SnakeCase converts TitleCase variable names to snake_case format.
func SnakeCase(name string) string {
	runes := []rune(name)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// MakeStringMap converts a map[string]string into the format:
// "key": "value", "key": "value"
func MakeStringMap(types map[string]string) string {
	var typArr []string
	for k, v := range types {
		typArr = append(typArr, fmt.Sprintf(`"%s": "%s"`, k, v))
	}

	return strings.Join(typArr, ", ")
}

// StringMap maps a function over a slice of strings.
func StringMap(modifier func(string) string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, str := range strs {
		ret[i] = modifier(str)
	}

	return ret
}

// MakeDBName takes a table name in the format of "table_name" and a
// column name in the format of "column_name" and returns a name used in the
// `db:""` component of an object in the format of "table_name_column_name"
func MakeDBName(tableName, colName string) string {
	return fmt.Sprintf("%s_%s", tableName, colName)
}

// HasElement checks to see if the string is found in the string slice
func HasElement(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}

// PrefixStringSlice with the given str.
func PrefixStringSlice(str string, strs []string) []string {
	ret := make([]string, len(strs))

	for i, s := range strs {
		ret[i] = fmt.Sprintf("%s%s", str, s)
	}

	return ret
}

// GenerateParamFlags generates the SQL statement parameter flags
// For example, $1,$2,$3 etc. It will start counting at startAt.
//
// If GroupAt is greater than 1, instead of returning $1,$2,$3
// it will return wrapped groups of param flags, for example:
//
//	GroupAt(1): $1,$2,$3,$4,$5,$6
//	GroupAt(2): ($1,$2),($3,$4),($5,$6)
//	GroupAt(3): ($1,$2,$3),($4,$5,$6),($7,$8,$9)
func GenerateParamFlags(colCount int, startAt int, groupAt int) string {
	var buf bytes.Buffer

	if groupAt > 1 {
		buf.WriteByte('(')
	}

	groupCounter := 0
	for i := startAt; i < colCount+startAt; i++ {
		groupCounter++
		buf.WriteString(fmt.Sprintf("$%d", i))
		if i+1 != colCount+startAt {
			if groupAt > 1 && groupCounter == groupAt {
				buf.WriteString("),(")
				groupCounter = 0
			} else {
				buf.WriteByte(',')
			}
		}
	}
	if groupAt > 1 {
		buf.WriteByte(')')
	}

	return buf.String()
}

// WhereClause is a version of Where that binds multiple checks together
// with an or statement.
// WhereMultiple(1, 2, "a", "b") = "(a=$1 and b=$2) or (a=$3 and b=$4)"
func WhereClause(start, count int, cols []string) string {
	if start == 0 {
		panic("0 is not a valid start number for whereMultiple")
	}

	buf := &bytes.Buffer{}
	for i := 0; i < count; i++ {
		if i != 0 {
			buf.WriteString(" OR ")
		}
		buf.WriteByte('(')
		for j, key := range cols {
			if j != 0 {
				buf.WriteString(" AND ")
			}
			fmt.Fprintf(buf, `"%s"=$%d`, key, start+i*len(cols)+j)
		}
		buf.WriteByte(')')
	}

	return buf.String()
}

// InClause generates SQL that could go inside an "IN ()" statement
// $1, $2, $3
func InClause(start, count int) string {
	if start == 0 {
		panic("0 is not a valid start number for inClause")
	}

	buf := &bytes.Buffer{}
	for i := 0; i < count; i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(buf, "$%d", i+start)
	}

	return buf.String()
}

// DriverUsesLastInsertID returns whether the database driver supports the
// sql.Result interface.
func DriverUsesLastInsertID(driverName string) bool {
	switch driverName {
	case "postgres":
		return false
	default:
		return true
	}
}

// Substring returns a substring of str starting at index start and going
// to end-1.
func Substring(start, end int, str string) string {
	return str[start:end]
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
