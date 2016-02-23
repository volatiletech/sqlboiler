package cmds

import (
	"bytes"
	"go/format"
	"os"
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
	out := generateDeletes()

	for _, v := range out {
		os.Stdout.Write(v)
	}
}

func generateDeletes() [][]byte {
	t, err := template.New("delete.tpl").Funcs(template.FuncMap{
		"makeGoColName": makeGoColName,
	}).ParseFiles("templates/delete.tpl")

	if err != nil {
		errorQuit(err)
	}

	var outputs [][]byte

	for i := 0; i < len(cmdData.TablesInfo); i++ {
		data := tplData{
			TableName: cmdData.TableNames[i],
			TableData: cmdData.TablesInfo[i],
		}

		var buf bytes.Buffer
		if err = t.Execute(&buf, data); err != nil {
			errorQuit(err)
		}

		out, err := format.Source(buf.Bytes())
		if err != nil {
			errorQuit(err)
		}

		outputs = append(outputs, out)
	}

	return outputs
}
