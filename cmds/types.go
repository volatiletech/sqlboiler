package cmds

import (
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

// CobraRunFunc declares the cobra.Command.Run function definition
type CobraRunFunc func(cmd *cobra.Command, args []string)

type templater []*template.Template

// CmdData holds the table schema a slice of (column name, column type) slices.
// It also holds a slice of all of the table names sqlboiler is generating against,
// the database driver chosen by the driver flag at runtime, and a pointer to the
// output file, if one is specified with a flag.
type CmdData struct {
	Tables     []dbdrivers.Table
	PkgName    string
	OutFolder  string
	Interface  dbdrivers.Interface
	DriverName string
	Config     *Config

	Templates templater
	// SingleTemplates are only created once, not per table
	SingleTemplates templater

	TestTemplates templater
	// SingleTestTemplates are only created once, not per table
	SingleTestTemplates templater
	//TestMainTemplate is only created once, not per table
	TestMainTemplate *template.Template
}

// tplData is used to pass data to the template
type tplData struct {
	Table      dbdrivers.Table
	DriverName string
	PkgName    string
	Tables     []string
}

type importList []string

// imports defines the optional standard imports and
// thirdparty imports (from github for example)
type imports struct {
	standard   importList
	thirdparty importList
}

type PostgresCfg struct {
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	DBName string `toml:"dbname"`
}

type Config struct {
	Postgres PostgresCfg `toml:"postgres"`
}
