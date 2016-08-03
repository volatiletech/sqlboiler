package bdb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/nullbio/sqlboiler/strmangle"
)

// Column holds information about a database column.
// Types are Go types, converted by TranslateColumnType.
type Column struct {
	Name     string
	Type     string
	DBType   string
	Default  string
	Nullable bool
	Unique   bool
	Enforced bool
}

// ColumnNames of the columns.
func ColumnNames(cols []Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
}

// ColumnTypes of the columns.
func ColumnTypes(cols []Column) []string {
	types := make([]string, len(cols))
	for i, c := range cols {
		types[i] = c.Type
	}

	return types
}

// ColumnDBTypes of the columns.
func ColumnDBTypes(cols []Column) map[string]string {
	types := map[string]string{}

	for _, c := range cols {
		types[strmangle.TitleCase(c.Name)] = c.DBType
	}

	return types
}

// FilterColumnsByDefault generates the list of columns that have default values
func FilterColumnsByDefault(defaults bool, columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if (defaults && len(c.Default) != 0) || (!defaults && len(c.Default) == 0) {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsBySimpleDefault generates a list of columns that have simple default values
// A simple default value is one without a function call and a non-enforced type
func FilterColumnsBySimpleDefault(columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if len(c.Default) != 0 && !strings.ContainsAny(c.Default, "()") && !c.Enforced {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsByAutoIncrement generates the list of auto increment columns
func FilterColumnsByAutoIncrement(autos bool, columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if (autos && rgxAutoIncColumn.MatchString(c.Default)) ||
			(!autos && !rgxAutoIncColumn.MatchString(c.Default)) {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsByEnforced generates the list of enforced columns
func FilterColumnsByEnforced(columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if c.Enforced == true {
			cols = append(cols, c)
		}
	}

	return cols
}

var (
	rgxRawDefaultValue   = regexp.MustCompile(`'(.*)'::`)
	rgxBoolDefaultValue  = regexp.MustCompile(`(?i)true|false`)
	rgxByteaDefaultValue = regexp.MustCompile(`(?i)\\x([0-9A-F]*)`)
)

// DefaultValues returns the Go converted values of the default value columns.
// For the time columns it will return time.Now() since we cannot extract
// the true time from the default value string.
func DefaultValues(columns []Column) []string {
	var dVals []string

	for _, c := range columns {
		var dVal string
		// Attempt to strip out the raw default value if its contained
		// within a Postgres type cast statement
		m := rgxRawDefaultValue.FindStringSubmatch(c.Default)
		if len(m) > 1 {
			dVal = m[len(m)-1]
		} else {
			dVal = c.Default
		}

		switch c.Type {
		case "null.Uint", "null.Uint8", "null.Uint16", "null.Uint32", "null.Uint64",
			"null.Int", "null.Int8", "null.Int16", "null.Int32", "null.Int64",
			"null.Float32", "null.Float64":
			dVals = append(dVals,
				fmt.Sprintf(`null.New%s(%s, true)`,
					strings.TrimPrefix(c.Type, "null."),
					dVal),
			)
		case "uint", "uint8", "uint16", "uint32", "uint64",
			"int", "int8", "int16", "int32", "int64", "float32", "float64":
			dVals = append(dVals, fmt.Sprintf(`%s(%s)`, c.Type, dVal))
		case "null.Bool":
			m = rgxBoolDefaultValue.FindStringSubmatch(dVal)
			if len(m) == 0 {
				dVals = append(dVals, `null.NewBool(false, true)`)
			}
			dVals = append(dVals, fmt.Sprintf(`null.NewBool(%s, true)`, strings.ToLower(dVal)))
		case "bool":
			m = rgxBoolDefaultValue.FindStringSubmatch(dVal)
			if len(m) == 0 {
				dVals = append(dVals, "false")
			}
			dVals = append(dVals, strings.ToLower(m[0]))
		case "null.Time":
			dVals = append(dVals, fmt.Sprintf(`null.NewTime(time.Now(), true)`))
		case "time.Time":
			dVals = append(dVals, `time.Now()`)
		case "null.String":
			dVals = append(dVals, fmt.Sprintf(`null.NewString("%s", true)`, dVal))
		case "string":
			dVals = append(dVals, `"`+dVal+`"`)
		case "[]byte":
			m := rgxByteaDefaultValue.FindStringSubmatch(dVal)
			if len(m) != 2 {
				dVals = append(dVals, `[]byte{}`)
			}
			hexstr := m[1]
			bs := make([]string, len(hexstr)/2)
			count := 0
			for i := 0; i < len(hexstr); i += 2 {
				bs[count] = "0x" + hexstr[i:i+2]
				count++
			}
			dVals = append(dVals, `[]byte{`+strings.Join(bs, ", ")+`}`)
		default:
			dVals = append(dVals, "")
		}
	}

	return dVals
}
