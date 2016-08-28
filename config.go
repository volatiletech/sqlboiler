package main

// Config for the running of the commands
type Config struct {
	DriverName       string   `toml:"driver_name"`
	PkgName          string   `toml:"pkg_name"`
	OutFolder        string   `toml:"out_folder"`
	BaseDir          string   `toml:"base_dir"`
	ExcludeTables    []string `toml:"exclude"`
	NoHooks          bool     `toml:"no_hooks"`
	NoAutoTimestamps bool     `toml:"no_auto_timestamps"`

	Postgres PostgresConfig `toml:"postgres"`
}

// PostgresConfig configures a postgres database
type PostgresConfig struct {
	User    string `toml:"user"`
	Pass    string `toml:"pass"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	DBName  string `toml:"dbname"`
	SSLMode string `toml:"sslmode"`
}
