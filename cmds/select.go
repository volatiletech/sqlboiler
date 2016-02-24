package cmds

import (
	"fmt"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

// init the "select" command
func init() {
	SQLBoiler.AddCommand(selectCmd)
	selectCmd.Run = selectRun
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Generate select statement helpers from table definitions",
}

// selectRun executes the select command, and generates the select statement
// boilerplate from the select file.
func selectRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateSelects())
	if err != nil {
		errorQuit(err)
	}
}

// generateSelects returns a slice of each template execution result.
// Each of these results holds a select statement generated from the select template.
func generateSelects() [][]byte {
	t, err := template.New("select.tpl").Funcs(template.FuncMap{
		"makeGoColName":        makeGoColName,
		"makeGoVarName":        makeGoVarName,
		"makeSelectParamNames": makeSelectParamNames,
	}).ParseFiles("templates/select.tpl")

	if err != nil {
		errorQuit(err)
	}

	outputs, err := processTemplate(t)
	if err != nil {
		errorQuit(err)
	}

	return outputs
}

// makeSelectParamNames takes a []DBTable and returns a comma seperated
// list of parameter names with for the select statement template.
// It also uses the table name to generate the "AS" part of the statement, for
// example: var_name AS table_name_var_name, ...
func makeSelectParamNames(tableName string, data []dbdrivers.DBTable) string {
	var paramNames string
	for i := 0; i < len(data); i++ {
		paramNames = fmt.Sprintf("%s%s AS %s", paramNames, data[i].ColName,
			makeDBColName(tableName, data[i].ColName),
		)
		if len(data) != i+1 {
			paramNames = paramNames + ", "
		}
	}
	return paramNames
}
