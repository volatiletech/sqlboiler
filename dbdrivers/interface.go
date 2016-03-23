package dbdrivers

// Interface for a database driver. Functionality required to support a specific
// database type (eg, MySQL, Postgres etc.)
type Interface interface {
	// AllTables connects to the database and retrieves all "public" table names
	AllTables() ([]string, error)

	// Columns retrieves column information about the table.
	Columns(tableName string) ([]Column, error)

	// TranslateColumn builds a Column out of a column metadata.
	// Its main responsibility is to convert database types to Go types, for
	// example "varchar" to "string".
	TranslateColumn(Column) Column

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
// Types are Go types, converted by TranslateColumn.
type Column struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNullable   bool
}
