package cmds

import (
	"text/template"

	"github.com/spf13/cobra"
)

type importList []string

// imports defines the optional standard imports and
// thirdparty imports (from github for example)
type imports struct {
	standard   importList
	thirdparty importList
}

// sqlBoilerDefaultImports defines the list of default template imports.
var sqlBoilerDefaultImports = imports{
	standard: importList{
		`"errors"`,
		`"fmt"`,
	},
	thirdparty: importList{
		`"github.com/pobri19/sqlboiler/boil"`,
	},
}

// sqlBoilerDefaultTestImports defines the list of default test template imports.
var sqlBoilerDefaultTestImports = imports{
	standard: importList{
		`"testing"`,
	},
}

// sqlBoilerConditionalTypeImports imports are only included in the template output
// if the database requires one of the following special types. Check ParseTableInfo
// to see the type assignments.
var sqlBoilerConditionalTypeImports = map[string]imports{
	"null.Int": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.String": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Bool": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Float": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Time": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"time.Time": imports{
		standard: importList{`"time"`},
	},
}

var sqlBoilerCustomImports map[string]imports
var sqlBoilerCustomTestImports map[string]imports

// sqlBoilerCommands points each command to its cobra.Command variable.
//
// If you would like to add your own custom command, add it to this
// map, and point it to your <commandName>Cmd variable.
//
// Command names should match the template file name (without the file extension).
var sqlBoilerCommands = map[string]*cobra.Command{
	// Command to generate all commands
	"boil": boilCmd,
	// Struct commands
	"struct": structCmd,
	// Insert commands
	"insert": insertCmd,
	// Select commands
	"all":         allCmd,
	"where":       whereCmd,
	"select":      selectCmd,
	"selectwhere": selectWhereCmd,
	"find":        findCmd,
	"findselect":  findSelectCmd,
	// Delete commands
	"delete": deleteCmd,
	// Update commands
	"update": updateCmd,
}

// sqlBoilerCommandRuns points each command to its custom run function.
// If a run function is not defined here, it will use the defaultRun() default run function.
var sqlBoilerCommandRuns = map[string]CobraRunFunc{
	"boil": boilRun,
}

// sqlBoilerTemplateFuncs is a map of all the functions that get passed into the templates.
// If you wish to pass a new function into your own template, add a pointer to it here.
var sqlBoilerTemplateFuncs = template.FuncMap{
	"singular":             singular,
	"plural":               plural,
	"titleCase":            titleCase,
	"titleCaseSingular":    titleCaseSingular,
	"titleCasePlural":      titleCasePlural,
	"camelCase":            camelCase,
	"camelCaseSingular":    camelCaseSingular,
	"camelCasePlural":      camelCasePlural,
	"makeDBName":           makeDBName,
	"selectParamNames":     selectParamNames,
	"insertParamNames":     insertParamNames,
	"insertParamFlags":     insertParamFlags,
	"insertParamVariables": insertParamVariables,
	"scanParamNames":       scanParamNames,
	"hasPrimaryKey":        hasPrimaryKey,
	"getPrimaryKey":        getPrimaryKey,
	"updateParamNames":     updateParamNames,
	"updateParamVariables": updateParamVariables,
}

/* Struct commands */

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Generate structs from table definitions",
}

/* Insert commands */

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

/* Select commands */

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate a helper to select all records",
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Generate a helper to select all records with specific column values",
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Generate a helper to select specific fields of all records",
}

var selectWhereCmd = &cobra.Command{
	Use:   "selectwhere",
	Short: "Generate a helper to select specific fields of records with specific column values",
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Generate a helper to select a single record by ID",
}

var findSelectCmd = &cobra.Command{
	Use:   "findselect",
	Short: "Generate a helper to select specific fields of a record by ID",
}

/* Delete commands */

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

/* Update commands */

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Generate update statement helpers from table definitions",
}
