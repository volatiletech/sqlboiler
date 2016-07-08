package bdb

import (
	"regexp"
	"strings"
)

// Column holds information about a database column.
// Types are Go types, converted by TranslateColumnType.
type Column struct {
	Name       string
	Type       string
	Default    string
	IsNullable bool
}

// ColumnNames of the columns.
func ColumnNames(cols []Column) []string {
	names := make([]string, len(cols))
	for i, c := range cols {
		names[i] = c.Name
	}

	return names
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
// A simple default value is one without a function call
func FilterColumnsBySimpleDefault(columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if len(c.Default) != 0 && !strings.ContainsAny(c.Default, "()") {
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

var (
	rgxRawDefaultValue   = regexp.MustCompile(`'(.*)'::`)
	rgxBoolDefaultValue  = regexp.MustCompile(`(?i)true|false`)
	rgxByteaDefaultValue = regexp.MustCompile(`(?i)\\x([0-9A-F]*)`)
)

// DefaultValues returns the Go converted values of the default value columns
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
			"uint", "uint8", "uint16", "uint32", "uint64",
			"int", "int8", "int16", "int32", "int64",
			"null.Float32", "null.Float64", "float32", "float64":
			dVals = append(dVals, dVal)
		case "null.Bool", "bool":
			m = rgxBoolDefaultValue.FindStringSubmatch(dVal)
			if len(m) == 0 {
				dVals = append(dVals, "false")
			}
			dVals = append(dVals, strings.ToLower(m[0]))
		case "null.Time", "time.Time", "null.String", "string":
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

// ZeroValue returns the zero value string of the column type
func ZeroValue(column Column) string {
	switch column.Type {
	case "null.Uint", "null.Uint8", "null.Uint16", "null.Uint32", "null.Uint64",
		"null.Int", "null.Int8", "null.Int16", "null.Int32", "null.Int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"int", "int8", "int16", "int32", "int64":
		return `0`
	case "null.Float32", "null.Float64", "float32", "float64":
		return `0.0`
	case "null.String", "string":
		return `""`
	case "null.Bool", "bool":
		return `false`
	case "null.Time", "time.Time":
		return `time.Time{}`
	case "[]byte":
		return `[]byte(nil)`
	default:
		return ""
	}
}
