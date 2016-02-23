package dbdrivers

type DBDriver interface {
	GetAllTableNames() ([]string, error)
	GetTableInfo(tableName string) ([]DBTable, error)
	ParseTableInfo(colName, colType string) DBTable
	// Open the database connection
	Open() error
	// Close the database connection
	Close()
}

type DBTable struct {
	ColName string
	ColType string
}
