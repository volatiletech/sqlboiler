package cmds

import (
	"text/template"

	"github.com/spf13/cobra"
)

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
	"all":          allCmd,
	"allby":        allByCmd,
	"fieldsall":    fieldsAllCmd,
	"fieldsallby":  fieldsAllByCmd,
	"find":         findCmd,
	"findby":       findByCmd,
	"fieldsfind":   fieldsFindCmd,
	"fieldsfindby": fieldsFindByCmd,
	// Delete commands
	"delete": deleteCmd,
}

// sqlBoilerCommandRuns points each command to its custom run function.
// If a run function is not defined here, it will use the defaultRun() default run function.
var sqlBoilerCommandRuns = map[string]CobraRunFunc{
	"boil": boilRun,
}

// sqlBoilerTemplateFuncs is a map of all the functions that get passed into the templates.
// If you wish to pass a new function into your own template, add a pointer to it here.
var sqlBoilerTemplateFuncs = template.FuncMap{
	"makeGoColName":          makeGoColName,
	"makeGoVarName":          makeGoVarName,
	"makeDBColName":          makeDBColName,
	"makeSelectParamNames":   makeSelectParamNames,
	"makeGoInsertParamNames": makeGoInsertParamNames,
	"makeGoInsertParamFlags": makeGoInsertParamFlags,
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Generate a helper to select all records",
}

var allByCmd = &cobra.Command{
	Use:   "allby",
	Short: "Generate a helper to select all records with specific column values",
}

var fieldsAllCmd = &cobra.Command{
	Use:   "fieldsall",
	Short: "Generate a helper to select specific fields of all records",
}

var fieldsAllByCmd = &cobra.Command{
	Use:   "fieldsallby",
	Short: "Generate a helper to select specific fields of records with specific column values",
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Generate a helper to select a single record by ID",
}

var findByCmd = &cobra.Command{
	Use:   "findby",
	Short: "Generate a helper to select a single record that has specific column values",
}

var fieldsFindCmd = &cobra.Command{
	Use:   "fieldsfind",
	Short: "Generate a helper to select specific fields of records by ID",
}

var fieldsFindByCmd = &cobra.Command{
	Use:   "fieldsfindby",
	Short: "Generate a helper to select specific fields of a single record that has specific column values",
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Generate structs from table definitions",
}
