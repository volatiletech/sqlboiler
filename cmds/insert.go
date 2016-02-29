package cmds

import (
	"text/template"

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
