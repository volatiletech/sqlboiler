package dbdrivers

// DBDriver is an interface that handles the operation uniqueness of each
// type of database connection. For example, queries to obtain schema data
// will vary depending on what type of database software is in use.
// The goal of the DBDriver is to retrieve all table names in a database
// using GetAllTableNames() if no table names are provided via flags,
// to handle the database connection using Open() and Close(), and to
// build the table information using GetTableInfo() and ParseTableInfo()
type DBDriver interface {
	// GetAllTableNames connects to the database and retrieves all "public" table names
	GetAllTableNames() ([]string, error)

	// GetTableInfo builds an object of []DBTable containing the table information
	GetTableInfo(tableName string) ([]DBTable, error)

	// ParseTableInfo builds a DBTable out of a column name and column type.
	// Its main responsibility is to convert database types to Go types, for example
	// "varchar" to "string".
	ParseTableInfo(colName, colType, isNullable string) DBTable

	// Open the database connection
	Open() error

	// Close the database connection
	Close()
}

// DBTable holds a column name, for example "column_name", and a column type,
// for example "int64". Column types are Go types, converted by ParseTableInfo.
type DBTable struct {
	ColName    string
	ColType    string
	IsNullable string
}
