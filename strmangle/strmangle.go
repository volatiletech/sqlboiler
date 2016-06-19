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

// TitleCaseSingular changes a snake-case variable name
// to a go styled object variable name of "ColumnName".
// titleCaseSingular also converts the last word in the
// variable name to a singularized version of itself.
func TitleCaseSingular(name string) string {
	return TitleCase(Singular(name))
}

// TitleCasePlural changes a snake-case variable name
// to a go styled object variable name of "ColumnName".
// titleCasePlural also converts the last word in the
// variable name to a pluralized version of itself.
func TitleCasePlural(name string) string {
	return TitleCase(Plural(name))
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

// CamelCaseSingular takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// CamelCaseSingular also converts the last word in the
// variable name to a singularized version of itself.
func CamelCaseSingular(name string) string {
	return CamelCase(Singular(name))
}

// CamelCasePlural takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// CamelCasePlural also converts the last word in the
// variable name to a pluralized version of itself.
func CamelCasePlural(name string) string {
	return CamelCase(Plural(name))
}

// CamelCaseCommaList generates a list of comma seperated camel cased column names
// example: thingName, o.stuffName, etc
func CamelCaseCommaList(prefix string, cols []string) string {
	var output []string

	for _, c := range cols {
		output = append(output, prefix+CamelCase(c))
	}

	return strings.Join(output, ", ")
}

// TitleCaseCommaList generates a list of comma seperated title cased column names
// example: o.ThingName, o.Stuff, ThingStuff, etc
func TitleCaseCommaList(prefix string, cols []string) string {
	var output []string

	for _, c := range cols {
		output = append(output, prefix+TitleCase(c))
	}

	return strings.Join(output, ", ")
}

// MakeDBName takes a table name in the format of "table_name" and a
// column name in the format of "column_name" and returns a name used in the
// `db:""` component of an object in the format of "table_name_column_name"
func MakeDBName(tableName, colName string) string {
	return tableName + "_" + colName
}

// UpdateParamNames takes a []Column and returns a comma seperated
// list of parameter names for the update statement template SET clause.
// eg: col1=$1,col2=$2,col3=$3
// Note: updateParamNames will exclude the PRIMARY KEY column.
func UpdateParamNames(columns []dbdrivers.Column, pkeyColumns []string) string {
	names := make([]string, 0, len(columns))
	counter := 0
	for _, c := range columns {
		if IsPrimaryKey(c.Name, pkeyColumns) {
			continue
		}
		counter++
		names = append(names, fmt.Sprintf("%s=$%d", c.Name, counter))
	}
	return strings.Join(names, ",")
}

// UpdateParamVariables takes a prefix and a []Columns and returns a
// comma seperated list of parameter variable names for the update statement.
// eg: prefix("o."), column("name_id") -> "o.NameID, ..."
// Note: updateParamVariables will exclude the PRIMARY KEY column.
func UpdateParamVariables(prefix string, columns []dbdrivers.Column, pkeyColumns []string) string {
	names := make([]string, 0, len(columns))

	for _, c := range columns {
		if IsPrimaryKey(c.Name, pkeyColumns) {
			continue
		}
		names = append(names, fmt.Sprintf("%s%s", prefix, TitleCase(c.Name)))
	}

	return strings.Join(names, ", ")
}

// IsPrimaryKey checks if the column is found in the primary key columns
func IsPrimaryKey(col string, pkeyCols []string) bool {
	for _, pkey := range pkeyCols {
		if pkey == col {
			return true
		}
	}

	return false
}

// InsertParamNames takes a []Column and returns a comma seperated
// list of parameter names for the insert statement template.
func InsertParamNames(columns []dbdrivers.Column) string {
	names := make([]string, len(columns))
	for i, c := range columns {
		names[i] = c.Name
	}
	return strings.Join(names, ", ")
}

// InsertParamFlags takes a []Column and returns a comma seperated
// list of parameter flags for the insert statement template.
func InsertParamFlags(columns []dbdrivers.Column) string {
	params := make([]string, len(columns))
	for i := range columns {
		params[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(params, ", ")
}

// InsertParamVariables takes a prefix and a []Columns and returns a
// comma seperated list of parameter variable names for the insert statement.
// For example: prefix("o."), column("name_id") -> "o.NameID, ..."
func InsertParamVariables(prefix string, columns []dbdrivers.Column) string {
	names := make([]string, len(columns))

	for i, c := range columns {
		names[i] = prefix + TitleCase(c.Name)
	}

	return strings.Join(names, ", ")
}

// SelectParamNames takes a []Column and returns a comma seperated
// list of parameter names with for the select statement template.
// It also uses the table name to generate the "AS" part of the statement, for
// example: var_name AS table_name_var_name, ...
func SelectParamNames(tableName string, columns []dbdrivers.Column) string {
	selects := make([]string, len(columns))
	for i, c := range columns {
		selects[i] = fmt.Sprintf("%s AS %s", c.Name, MakeDBName(tableName, c.Name))
	}

	return strings.Join(selects, ", ")
}

// ScanParamNames takes a []Column and returns a comma seperated
// list of parameter names for use in a db.Scan() call.
func ScanParamNames(object string, columns []dbdrivers.Column) string {
	scans := make([]string, len(columns))
	for i, c := range columns {
		scans[i] = fmt.Sprintf("&%s.%s", object, TitleCase(c.Name))
	}

	return strings.Join(scans, ", ")
}

// HasPrimaryKey returns true if one of the columns passed in is a primary key
func HasPrimaryKey(pKey *dbdrivers.PrimaryKey) bool {
	if pKey == nil || len(pKey.Columns) == 0 {
		return false
	}

	return true
}

// PrimaryKeyFuncSig generates the function signature parameters.
// example: id int64, thingName string
func PrimaryKeyFuncSig(cols []dbdrivers.Column, pkeyCols []string) string {
	var output []string

	for _, pk := range pkeyCols {
		for _, c := range cols {
			if pk == c.Name {
				output = append(output, fmt.Sprintf("%s %s", CamelCase(pk), c.Type))
				break
			}
		}
	}

	return strings.Join(output, ", ")
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
	var output string

	// 0 is not a valid start number
	if start == 0 {
		start = 1
	}

	cols := make([]string, len(pkeyCols))
	copy(cols, pkeyCols)

	for i, c := range cols {
		output = fmt.Sprintf("%s%s=$%d", output, c, start)
		start++

		if i < len(cols)-1 {
			output = fmt.Sprintf("%s AND ", output)
		}
	}

	return output
}

// AutoIncPrimaryKey returns the auto-increment primary key column name or an
// empty string.
func AutoIncPrimaryKey(cols []dbdrivers.Column, pkey *dbdrivers.PrimaryKey) string {
	if pkey == nil {
		return ""
	}

	for _, c := range cols {
		if rgxAutoIncColumn.MatchString(c.Default) &&
			c.IsNullable == false &&
			(strings.HasPrefix(c.Type, "int") || strings.HasPrefix(c.Type, "uint")) {
			for _, p := range pkey.Columns {
				if c.Name == p {
					return p
				}
			}
		}
	}

	return ""
}

// ColumnsToStrings changes the columns into a list of column names
func ColumnsToStrings(cols []dbdrivers.Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
}

// CommaList returns a comma seperated list: "col1", "col2", "col3"
func CommaList(cols []string) string {
	return fmt.Sprintf(`"%s"`, strings.Join(cols, `", "`))
}

// ParamsPrimaryKey returns the parameters for the sql statement $ flags
// For example, if prefix was "o.", and titleCase was true: "o.ColumnName1, o.ColumnName2"
func ParamsPrimaryKey(prefix string, columns []string, shouldTitleCase bool) string {
	names := make([]string, 0, len(columns))

	for _, c := range columns {
		var n string
		if shouldTitleCase {
			n = prefix + TitleCase(c)
		} else {
			n = prefix + c
		}
		names = append(names, n)
	}

	return strings.Join(names, ", ")
}

// PrimaryKeyFlagIndex generates the primary key column flag number for the query params
func PrimaryKeyFlagIndex(regularCols []dbdrivers.Column, pkeyCols []string) int {
	return len(regularCols) - len(pkeyCols) + 1
}

// SupportsResultObject returns whether the database driver supports the
// sql.Results interface, i.e. LastReturnId and RowsAffected
func SupportsResultObject(driverName string) bool {
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

// AddID to the end of the string
func AddID(str string) string {
	return str + "_id"
}

// RemoveID from the end of the string
func RemoveID(str string) string {
	return strings.TrimSuffix(str, "_id")
}

// Substring returns a substring of str starting at index start and going
// to end-1.
func Substring(start, end int, str string) string {
	return str[start:end]
}
