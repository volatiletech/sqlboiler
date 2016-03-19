package boil

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// SelectNames returns the column names for a select statement
// Eg: col1, col2, col3
func SelectNames(results interface{}) string {
	var names []string

	structValue := reflect.ValueOf(results)
	if structValue.Kind() == reflect.Ptr {
		structValue = structValue.Elem()
	}

	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		var name string

		if db := field.Tag.Get("db"); len(db) != 0 {
			name = db
		} else {
			name = goVarToSQLName(field.Name)
		}

		names = append(names, name)
	}

	return strings.Join(names, ", ")
}

// Where returns the where clause for an sql statement
// eg: col1=$1 AND col2=$2 AND col3=$3
func Where(columns map[string]interface{}) string {
	names := make([]string, 0, len(columns))

	for c := range columns {
		names = append(names, c)
	}

	sort.Strings(names)

	for i, c := range names {
		names[i] = fmt.Sprintf("%s=$%d", c, i+1)
	}

	return strings.Join(names, " AND ")
}

// Update returns the column list for an update statement SET clause
// eg: col1=$1,col2=$2,col3=$3
func Update(columns map[string]interface{}) string {
	names := make([]string, 0, len(columns))

	for c := range columns {
		names = append(names, c)
	}

	sort.Strings(names)

	for i, c := range names {
		names[i] = fmt.Sprintf("%s=$%d", c, i+1)
	}

	return strings.Join(names, ",")
}

// WhereParams returns a list of sql parameter values for the query
func WhereParams(columns map[string]interface{}) []interface{} {
	names := make([]string, 0, len(columns))
	results := make([]interface{}, 0, len(columns))

	for c := range columns {
		names = append(names, c)
	}

	sort.Strings(names)

	for _, c := range names {
		results = append(results, columns[c])
	}

	return results
}

// goVarToSQLName converts a go variable name to a column name
// example: HelloFriendID to hello_friend_id
func goVarToSQLName(name string) string {
	str := &bytes.Buffer{}
	isUpper, upperStreak := false, false

	for i := 0; i < len(name); i++ {
		c := rune(name[i])
		if unicode.IsDigit(c) || unicode.IsLower(c) {
			isUpper = false
			upperStreak = false

			str.WriteRune(c)
			continue
		}

		if isUpper {
			upperStreak = true
		} else if i != 0 {
			str.WriteByte('_')
		}
		isUpper = true

		if j := i + 1; j < len(name) && upperStreak && unicode.IsLower(rune(name[j])) {
			str.WriteByte('_')
		}

		str.WriteRune(unicode.ToLower(c))
	}

	return str.String()
}
