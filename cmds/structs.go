package cmds

import (
	"bytes"
	"go/format"
	"os"
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

type tplData struct {
	TableName string
	TableData interface{}
}

func structRun(cmd *cobra.Command, args []string) {
	out := generateStructs()

	for _, v := range out {
		os.Stdout.Write(v)
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

	var structOutputs [][]byte

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

		structOutputs = append(structOutputs, out)
	}

	return structOutputs
}
