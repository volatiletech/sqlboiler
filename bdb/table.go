package bdb

// Table metadata from the database schema.
type Table struct {
	Name    string
	Columns []Column

	PKey  *PrimaryKey
	FKeys []ForeignKey

	IsJoinTable bool
}

// GetColumn by name. Panics if not found (for use in templates).
func (t Table) GetColumn(name string) (col Column) {
	for _, c := range t.Columns {
		if c.Name == name {
			return c
		}
	}

	panic("hello")
}
