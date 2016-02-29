package cmds

import (
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
)

type CobraRunFunc func(cmd *cobra.Command, args []string)

var sqlBoilerCommands = map[string]*cobra.Command{
	"all":    allCmd,
	"insert": insertCmd,
	"delete": deleteCmd,
	"select": selectCmd,
	"struct": structCmd,
}

// sqlBoilerCommandRuns points each command to its custom run function.
// If a run function is not defined here, it will use the
// defaultRun run function.
var sqlBoilerCommandRuns = map[string]CobraRunFunc{
	"all": allRun,
}

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
	Short: "Generate all templates from table definitions",
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

func defaultRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateTemplate(cmd.Name()))
	if err != nil {
		errorQuit(fmt.Errorf("Unable to generate the template for command %s: %s", cmd.Name(), err))
	}
}

func generateTemplate(name string) [][]byte {
	var template *template.Template
	for _, t := range templates {
		fmt.Printf("File name: %s", t.Name())
		if t.Name() == name+".tpl" {
			template = t
			break
		}
	}

	outputs, err := processTemplate(template)
	if err != nil {
		errorQuit(fmt.Errorf("Unable to process the template: %s", err))
	}

	return outputs
}
