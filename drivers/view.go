package drivers

import (
	"fmt"
	"strings"
)

// View metadata from the database schema.
type View struct {
	Name string `json:"name"`
	// For dbs with real schemas, like Postgres.
	// Example value: "schema_name"."view_name"
	SchemaName string   `json:"schema_name"`
	Columns    []Column `json:"columns"`
}

// GetView by name. Panics if not found (for use in templates mostly).
func GetView(views []View, name string) (tbl View) {
	for _, t := range views {
		if t.Name == name {
			return t
		}
	}

	panic(fmt.Sprintf("could not find view name: %s", name))
}

// GetColumn by name. Panics if not found (for use in templates mostly).
func (t View) GetColumn(name string) (col Column) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}

	panic(fmt.Sprintf("could not find column name: %s", name))
}

func (v View) CanSoftDelete(deleteColumn string) bool {
	if deleteColumn == "" {
		deleteColumn = "deleted_at"
	}

	for _, column := range v.Columns {
		if column.Name == deleteColumn && column.Type == "null.Time" {
			return true
		}
	}
	return false
}

func ViewsHaveNullableEnums(views []View) bool {
	for _, view := range views {
		for _, col := range view.Columns {
			if col.Nullable &&
				(strings.HasPrefix(col.DBType, "enum.") || // postgresql
					strings.HasPrefix(col.DBType, "enum(")) { // mysql
				return true
			}
		}
	}
	return false
}
