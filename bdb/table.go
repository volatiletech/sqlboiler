package bdb

import "fmt"

// Table metadata from the database schema.
type Table struct {
	Name string
	// For dbs with real schemas, like Postgres.
	// Example value: "schema_name"."table_name"
	SchemaName string
	Columns    []Column

	PKey  *PrimaryKey
	FKeys []ForeignKey

	IsJoinTable bool

	ToOneRelationships  []ToOneRelationship
	ToManyRelationships []ToManyRelationship
}

// GetTable by name. Panics if not found (for use in templates mostly).
func GetTable(tables []Table, name string) (tbl Table) {
	for _, t := range tables {
		if t.Name == name {
			return t
		}
	}

	panic(fmt.Sprintf("could not find table name: %s", name))
}

// GetColumn by name. Panics if not found (for use in templates mostly).
func (t Table) GetColumn(name string) (col Column) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}

	panic(fmt.Sprintf("could not find column name: %s", name))
}
