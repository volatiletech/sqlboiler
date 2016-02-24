package cmds

import (
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	SQLBoiler.AddCommand(structCmd)
	structCmd.Run = structRun
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Generate structs from table definitions",
}

func structRun(cmd *cobra.Command, args []string) {
	err := outHandler(generateStructs())
	if err != nil {
		errorQuit(err)
	}
}

func generateStructs() [][]byte {
	t, err := template.New("struct.tpl").Funcs(template.FuncMap{
		"makeGoColName": makeGoColName,
		"makeDBColName": makeDBColName,
	}).ParseFiles("templates/struct.tpl")

	if err != nil {
		errorQuit(err)
	}

	outputs, err := processTemplate(t)
	if err != nil {
		errorQuit(err)
	}

	return outputs
}
