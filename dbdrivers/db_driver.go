package dbdrivers

// DBDriver is an interface that handles the operation uniqueness of each
// type of database connection. For example, queries to obtain schema data
// will vary depending on what type of database software is in use.
// The goal of the DBDriver is to retrieve all table names in a database
// using GetAllTables() if no table names are provided via flags,
// to handle the database connection using Open() and Close(), and to
// build the table information using GetTableInfo() and ParseTableInfo()
type DBDriver interface {
	// GetAllTables connects to the database and retrieves all "public" table names
	GetAllTables() ([]string, error)

	// GetTableInfo retrieves column information about the table.
	GetTableInfo(tableName string) ([]DBColumn, error)

	// ParseTableInfo builds a DBColumn out of a column name and column type.
	// Its main responsibility is to convert database types to Go types, for example
	// "varchar" to "string".
	ParseTableInfo(name, colType string, isNullable bool, isPrimary bool) DBColumn

	// Open the database connection
	Open() error

	// Close the database connection
	Close()
}

// DBColumn holds information about a database column name.
// Column types are Go types, converted by ParseTableInfo.
type DBColumn struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNullable   bool
}
