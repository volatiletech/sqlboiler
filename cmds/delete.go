package cmds

import (
	"text/template"

	"github.com/spf13/cobra"
)

// init the "delete" command
func init() {
	SQLBoiler.AddCommand(deleteCmd)
	deleteCmd.Run = deleteRun
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

// deleteRun executes the delete command, and generates the delete statement
// boilerplate from the template file.
func deleteRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateDeletes())
	if err != nil {
		errorQuit(err)
	}
}

// generateDeletes returns a slice of each template execution result.
// Each of these results holds a delete statement generated from the delete template.
func generateDeletes() [][]byte {
	t, err := template.New("delete.tpl").Funcs(template.FuncMap{
		"makeGoColName": makeGoColName,
	}).ParseFiles("templates/delete.tpl")

	if err != nil {
		errorQuit(err)
	}

	outputs, err := processTemplate(t)
	if err != nil {
		errorQuit(err)
	}

	return outputs
}
