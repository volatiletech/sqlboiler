package cmds

import (
	"fmt"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(insertCmd)
	insertCmd.Run = insertRun
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

func insertRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateInserts())
	if err != nil {
		errorQuit(err)
	}
}

func generateInserts() [][]byte {
	t, err := template.New("insert.tpl").Funcs(template.FuncMap{
		"makeGoColName":          makeGoColName,
		"makeDBColName":          makeDBColName,
		"makeGoInsertParamNames": makeGoInsertParamNames,
		"makeGoInsertParamFlags": makeGoInsertParamFlags,
	}).ParseFiles("templates/insert.tpl")

	if err != nil {
		errorQuit(err)
	}

	outputs, err := processTemplate(t)
	if err != nil {
		errorQuit(err)
	}

	return outputs
}

// makeGoInsertParamNames takes a [][]DBData and returns a comma seperated
// list of parameter names for the insert statement
func makeGoInsertParamNames(data []dbdrivers.DBTable) string {
	var paramNames string
	for i := 0; i < len(data); i++ {
		paramNames = paramNames + data[i].ColName
		if len(data) != i+1 {
			paramNames = paramNames + ", "
		}
	}
	return paramNames
}

// makeGoInsertParamFlags takes a [][]DBData and returns a comma seperated
// list of parameter flags for the insert statement
func makeGoInsertParamFlags(data []dbdrivers.DBTable) string {
	var paramFlags string
	for i := 0; i < len(data); i++ {
		paramFlags = fmt.Sprintf("%s$%d", paramFlags, i+1)
		if len(data) != i+1 {
			paramFlags = paramFlags + ", "
		}
	}
	return paramFlags
}
