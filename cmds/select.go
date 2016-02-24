package cmds

import (
	"fmt"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(selectCmd)
	selectCmd.Run = selectRun
}

var selectCmd = &cobra.Command{
	Use:   "select",
	Short: "Generate select statement helpers from table definitions",
}

func selectRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateSelects())
	if err != nil {
		errorQuit(err)
	}
}

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
