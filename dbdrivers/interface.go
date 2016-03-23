package dbdrivers

// Interface for a database driver. Functionality required to support a specific
// database type (eg, MySQL, Postgres etc.)
type Interface interface {
	// Tables connects to the database and retrieves the table metadata for
	// the given tables, or all tables if none are provided.
	Tables(names ...string) ([]Table, error)

	// TranslateColumnType takes a Database column type and returns a go column
	// type.
	TranslateColumnType(Column) Column

	// Open the database connection
	Open() error

	// Close the database connection
	Close()
}

// Table metadata from the database schema.
type Table struct {
	Name    string
	Columns []Column

	PKey  *PrimaryKey
	FKeys []ForeignKey

	IsJoinTable bool
}

// Column holds information about a database column.
// Types are Go types, converted by TranslateColumnType.
type Column struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNullable   bool
}

// PrimaryKey represents a primary key constraint in a database
type PrimaryKey struct {
	Name    string
	Columns []string
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	Name   string
	Column string

	ForeignTable  string
	ForeignColumn string
}
