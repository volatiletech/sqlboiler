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
