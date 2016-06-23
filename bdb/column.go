package bdb

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
func FilterColumnsByDefault(columns []Column, defaults bool) []Column {
	var cols []Column

	for _, c := range columns {
		if (defaults && len(c.Default) != 0) || (!defaults && len(c.Default) == 0) {
			cols = append(cols, c)
		}
	}

	return cols
}

// FilterColumnsByAutoIncrement generates the list of auto increment columns
func FilterColumnsByAutoIncrement(columns []Column) []Column {
	var cols []Column

	for _, c := range columns {
		if rgxAutoIncColumn.MatchString(c.Default) {
			cols = append(cols, c)
		}
	}

	return cols
}
