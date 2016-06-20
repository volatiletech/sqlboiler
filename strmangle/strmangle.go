// Package strmangle is used exclusively by the templates in sqlboiler.
// There are many helper functions to deal with dbdrivers.* values as well
// as string manipulation. Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/nullbio/sqlboiler/dbdrivers"
)

var rgxAutoIncColumn = regexp.MustCompile(`^nextval\(.*\)`)

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
		if split == "id" {
			splits[i] = "ID"
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
		if split == "id" && i > 0 {
			splits[i] = "ID"
			continue
		}

		if i == 0 {
			continue
		}

		splits[i] = strings.Title(split)
	}

	return strings.Join(splits, "")
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

// PrimaryKeyFuncSig generates the function signature parameters.
// example: id int64, thingName string
func PrimaryKeyFuncSig(cols []dbdrivers.Column, pkeyCols []string) string {
	ret := make([]string, len(pkeyCols))

	for i, pk := range pkeyCols {
		for _, c := range cols {
			if pk != c.Name {
				continue
			}

			ret[i] = fmt.Sprintf("%s %s", CamelCase(pk), c.Type)
		}
	}

	return strings.Join(ret, ", ")
}

// GenerateParamFlags generates the SQL statement parameter flags
// For example, $1,$2,$3 etc. It will start counting at startAt.
func GenerateParamFlags(colCount int, startAt int) string {
	cols := make([]string, 0, colCount)

	for i := startAt; i < colCount+startAt; i++ {
		cols = append(cols, fmt.Sprintf("$%d", i))
	}

	return strings.Join(cols, ",")
}

// WherePrimaryKey returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WherePrimaryKey(pkeyCols []string, start int) string {
	if start == 0 {
		panic("0 is not a valid start number for wherePrimaryKey")
	}

	cols := make([]string, len(pkeyCols))
	for i, c := range pkeyCols {
		cols[i] = fmt.Sprintf("%s=$%d", c, start+i)
	}

	return strings.Join(cols, " AND ")
}

// AutoIncPrimaryKey returns the auto-increment primary key column name or an
// empty string.
func AutoIncPrimaryKey(cols []dbdrivers.Column, pkey *dbdrivers.PrimaryKey) string {
	if pkey == nil {
		return ""
	}

	for _, pkeyColumn := range pkey.Columns {
		for _, c := range cols {
			if c.Name != pkeyColumn {
				continue
			}

			if !rgxAutoIncColumn.MatchString(c.Default) || c.IsNullable ||
				!(strings.HasPrefix(c.Type, "int") || strings.HasPrefix(c.Type, "uint")) {
				continue
			}

			return pkeyColumn
		}
	}

	return ""
}

// ColumnNames of the columns.
func ColumnNames(cols []dbdrivers.Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
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

// FilterColumnsByDefault generates the list of columns that have default values
func FilterColumnsByDefault(columns []dbdrivers.Column, defaults bool) string {
	var cols []string

	for _, c := range columns {
		if (defaults && len(c.Default) != 0) || (!defaults && len(c.Default) == 0) {
			cols = append(cols, fmt.Sprintf(`"%s"`, c.Name))
		}
	}

	return strings.Join(cols, `,`)
}

// FilterColumnsByAutoIncrement generates the list of auto increment columns
func FilterColumnsByAutoIncrement(columns []dbdrivers.Column) string {
	var cols []string

	for _, c := range columns {
		if rgxAutoIncColumn.MatchString(c.Default) {
			cols = append(cols, fmt.Sprintf(`"%s"`, c.Name))
		}
	}

	return strings.Join(cols, `,`)
}

// Substring returns a substring of str starting at index start and going
// to end-1.
func Substring(start, end int, str string) string {
	return str[start:end]
}
