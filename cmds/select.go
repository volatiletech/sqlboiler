package cmds

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
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
	out := generateSelects()

	for _, v := range out {
		os.Stdout.Write(v)
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
