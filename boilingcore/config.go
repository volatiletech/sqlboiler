package boilingcore

// Config for the running of the commands
type Config struct {
	DriverName       string
	Schema           string
	PkgName          string
	OutFolder        string
	BaseDir          string
	WhitelistTables  []string
	BlacklistTables  []string
	Tags             []string
	Replacements     []string
	Debug            bool
	NoTests          bool
	NoHooks          bool
	NoAutoTimestamps bool
	Wipe             bool

	Postgres PostgresConfig
	MySQL    MySQLConfig
	MSSQL    MSSQLConfig
	SQLite   SQLiteConfig
}

// PostgresConfig configures a postgres database
type PostgresConfig struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DBName  string
	SSLMode string
}

// MySQLConfig configures a mysql database
type MySQLConfig struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DBName  string
	SSLMode string
}

// MSSQLConfig configures a mysql database
type MSSQLConfig struct {
	User    string
	Pass    string
	Host    string
	Port    int
	DBName  string
	SSLMode string
}

// SQLiteConfig configures a sqlite db
type SQLiteConfig struct {
	File string
}
