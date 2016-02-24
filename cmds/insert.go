package cmds

import (
	"fmt"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
	"github.com/spf13/cobra"
)

// init the "insert" command
func init() {
	SQLBoiler.AddCommand(insertCmd)
	insertCmd.Run = insertRun
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

// insertRun executes the insert command, and generates the insert statement
// boilerplate from the template file.
func insertRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateInserts())
	if err != nil {
		errorQuit(err)
	}
}

// generateInserts returns a slice of each template execution result.
// Each of these results holds an insert statement generated from the insert template.
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

// makeGoInsertParamNames takes a []DBTable and returns a comma seperated
// list of parameter names for the insert statement template.
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

// makeGoInsertParamFlags takes a []DBTable and returns a comma seperated
// list of parameter flags for the insert statement template.
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
