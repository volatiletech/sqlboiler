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
	"all":    allCmd,
	"insert": insertCmd,
	"delete": deleteCmd,
	"select": selectCmd,
	"struct": structCmd,
}

// sqlBoilerCommandRuns points each command to its custom run function.
// If a run function is not defined here, it will use the defaultRun() default run function.
var sqlBoilerCommandRuns = map[string]CobraRunFunc{
	"all": allRun,
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

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Generate select statement helpers from table definitions",
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Generate structs from table definitions",
}
