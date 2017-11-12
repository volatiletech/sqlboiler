package boilingcore

import "github.com/volatiletech/sqlboiler/importers"

// Config for the running of the commands
type Config struct {
	DriverName   string
	DriverConfig map[string]interface{}

	Schema           string
	PkgName          string
	OutFolder        string
	BaseDir          string
	Tags             []string
	Replacements     []string
	Debug            bool
	NoTests          bool
	NoHooks          bool
	NoAutoTimestamps bool
	Wipe             bool
	StructTagCasing  string

	Imports importers.Collection
}
