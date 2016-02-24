package cmds

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

type tplData struct {
	TableName string
	TableData []dbdrivers.DBTable
}

func errorQuit(err error) {
	fmt.Println(fmt.Sprintf("Error: %s\n---\n", err))
	structCmd.Help()
	os.Exit(-1)
}

func processTemplate(t *template.Template) ([][]byte, error) {
	var outputs [][]byte
	for i := 0; i < len(cmdData.TablesInfo); i++ {
		data := tplData{
			TableName: cmdData.TableNames[i],
			TableData: cmdData.TablesInfo[i],
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return nil, err
		}

		out, err := format.Source(buf.Bytes())
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, out)
	}

	return outputs, nil

}

func outHandler(data [][]byte) error {
	nl := []byte{'\n'}

	var out *os.File
	if cmdData.OutFile == nil {
		out = os.Stdout
	} else {
		out = cmdData.OutFile
	}

	for _, v := range data {
		if _, err := out.Write(v); err != nil {
			return err
		}
		if _, err := out.Write(nl); err != nil {
			return err
		}
	}

	return nil
}

func makeGoColName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {
		if s[i] == "id" {
			s[i] = "ID"
			continue
		}
		s[i] = strings.Title(s[i])
	}

	return strings.Join(s, "")
}

func makeGoVarName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {
		if s[i] == "id" {
			s[i] = "ID"
			continue
		}

		// Skip first word Title for variable names
		if i == 0 {
			continue
		}

		s[i] = strings.Title(s[i])
	}

	return strings.Join(s, "")
}

func makeDBColName(tableName, colName string) string {
	return tableName + "_" + colName
}
