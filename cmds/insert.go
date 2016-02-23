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
	SQLBoiler.AddCommand(insertCmd)
	insertCmd.Run = insertRun
}

var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "Generate insert statement helpers from table definitions",
}

func insertRun(cmd *cobra.Command, args []string) {
	out := generateInserts()

	for _, v := range out {
		os.Stdout.Write(v)
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
