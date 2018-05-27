package boilingcore

import (
	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/importers"
)

// Config for the running of the commands
type Config struct {
	DriverName   string
	DriverConfig drivers.Config

	PkgName          string
	OutFolder        string
	BaseDir          string
	Tags             []string
	Replacements     []string
	Debug            bool
	NoTests          bool
	NoHooks          bool
	NoAutoTimestamps bool
	NoRowsAffected   bool
	Wipe             bool
	StructTagCasing  string

	Imports importers.Collection
}
