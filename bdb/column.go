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

// DefaultValue returns the Go converted value of the default value column
func DefaultValue(column Column) string {
	defaultVal := ""

	// Attempt to strip out the raw default value if its contained
	// within a Postgres type cast statement
	m := rgxRawDefaultValue.FindStringSubmatch(column.Default)
	if len(m) > 1 {
		defaultVal = m[len(m)-1]
	} else {
		defaultVal = column.Default
	}

	switch column.Type {
	case "null.Uint", "null.Uint8", "null.Uint16", "null.Uint32", "null.Uint64",
		"null.Int", "null.Int8", "null.Int16", "null.Int32", "null.Int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"int", "int8", "int16", "int32", "int64",
		"null.Float32", "null.Float64", "float32", "float64":
		return defaultVal
	case "null.Bool", "bool":
		m = rgxBoolDefaultValue.FindStringSubmatch(defaultVal)
		if len(m) == 0 {
			return "false"
		}
		return strings.ToLower(m[0])
	case "null.Time", "time.Time", "null.String", "string":
		return `"` + defaultVal + `"`
	case "[]byte":
		m := rgxByteaDefaultValue.FindStringSubmatch(defaultVal)
		if len(m) != 2 {
			return `[]byte{}`
		}
		hexstr := m[1]
		bs := make([]string, len(hexstr)/2)
		c := 0
		for i := 0; i < len(hexstr); i += 2 {
			bs[c] = "0x" + hexstr[i:i+2]
			c++
		}
		return `[]byte{` + strings.Join(bs, ", ") + `}`
	default:
		return ""
	}
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
		return `[]byte{}`
	default:
		return ""
	}
}
