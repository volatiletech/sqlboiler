package bdb

// Table metadata from the database schema.
type Table struct {
	Name    string
	Columns []Column

	PKey  *PrimaryKey
	FKeys []ForeignKey

	IsJoinTable bool
}
