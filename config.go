package main

// Config for the running of the commands
type Config struct {
	DriverName       string
	PkgName          string
	OutFolder        string
	BaseDir          string
	ExcludeTables    []string
	Debug            bool
	NoTests          bool
	NoHooks          bool
	NoAutoTimestamps bool

	Postgres PostgresConfig
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
