package cmds

import (
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(deleteCmd)
	deleteCmd.Run = deleteRun
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Generate delete statement helpers from table definitions",
}

func deleteRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateDeletes())
	if err != nil {
		errorQuit(err)
	}
}

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
