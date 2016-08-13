// Package strmangle is used exclusively by the templates in sqlboiler.
// There are many helper functions to deal with bdb.* values as well
// as string manipulation. Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

import (
	"fmt"
	"math"
	"regexp"
	"strings"

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

	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(s, ".")
	for i, split := range splits {
		if i != 0 {
			buf.WriteByte('.')
		}

		if strings.HasPrefix(split, `"`) || strings.HasSuffix(split, `"`) || split == "*" {
			buf.WriteString(split)
			continue
		}

		buf.WriteByte('"')
		buf.WriteString(split)
		buf.WriteByte('"')
	}

	return buf.String()
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
			buf.WriteString(inflection.Plural(splits[len(splits)-1]))
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
			buf.WriteString(inflection.Singular(splits[len(splits)-1]))
			break
		}

		buf.WriteString(splits[i])
	}

	return buf.String()
}

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
func TitleCase(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(name, "_")

	for _, split := range splits {
		if uppercaseWords.MatchString(split) {
			buf.WriteString(strings.ToUpper(split))
			continue
		}

		buf.WriteString(strings.Title(split))
	}

	return buf.String()
}

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func CamelCase(name string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	splits := strings.Split(name, "_")

	for i, split := range splits {
		if i == 0 {
			buf.WriteString(split)
			continue
		}

		if uppercaseWords.MatchString(split) {
			buf.WriteString(strings.ToUpper(split))
			continue
		}

		buf.WriteString(strings.Title(split))
	}

	return buf.String()
}

// MakeStringMap converts a map[string]string into the format:
// "key": "value", "key": "value"
func MakeStringMap(types map[string]string) string {
	buf := GetBuffer()
	defer PutBuffer(buf)

	c := 0
	for k, v := range types {
		buf.WriteString(fmt.Sprintf(`"%s": "%s"`, k, v))
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

// Placeholders generates the SQL statement placeholders for in queries.
// For example, ($1, $2, $3), ($4, $5, $6) etc.
// It will start counting placeholders at "start".
func Placeholders(count int, start int, group int) string {
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
		buf.WriteString(fmt.Sprintf("$%d", start+i))
	}
	if group > 1 {
		buf.WriteByte(')')
	}

	return buf.String()
}

// SetParamNames takes a slice of columns and returns a comma separated
// list of parameter names for a template statement SET clause.
// eg: "col1"=$1, "col2"=$2, "col3"=$3
func SetParamNames(columns []string) string {
	names := make([]string, 0, len(columns))
	counter := 0
	for _, c := range columns {
		counter++
		names = append(names, fmt.Sprintf(`"%s"=$%d`, c, counter))
	}
	return strings.Join(names, ", ")
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(start int, cols []string) string {
	if start == 0 {
		panic("0 is not a valid start number for whereClause")
	}

	ret := make([]string, len(cols))
	for i, c := range cols {
		ret[i] = fmt.Sprintf(`"%s"=$%d`, c, start+i)
	}

	return strings.Join(ret, " AND ")
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
