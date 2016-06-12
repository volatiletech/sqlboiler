package sqlboiler

import (
	"text/template"

	"github.com/nullbio/sqlboiler/dbdrivers"
)

// State holds the global data needed by most pieces to run
type State struct {
	Config *Config

	Driver dbdrivers.Interface
	Tables []dbdrivers.Table

	Templates              templateList
	TestTemplates          templateList
	SingletonTemplates     templateList
	SingletonTestTemplates templateList

	TestMainTemplate *template.Template
}

// Config for the running of the commands
type Config struct {
	DriverName string `toml:"driver_name"`
	PkgName    string `toml:"pkg_name"`
	OutFolder  string `toml:"out_folder"`
	TableName  string `toml:"table_name"`

	Postgres PostgresConfig `toml:"postgres"`
}

// PostgresConfig configures a postgres database
type PostgresConfig struct {
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	DBName string `toml:"dbname"`
}
