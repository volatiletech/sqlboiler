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

// tplData is used to pass data to the template
type tplData struct {
	TableName string
	TableData []dbdrivers.DBTable
}

// errorQuit displays an error message and then exits the application.
func errorQuit(err error) {
	fmt.Println(fmt.Sprintf("Error: %s\n---\n\nRun 'sqlboiler --help' for usage.", err))
	os.Exit(-1)
}

// processTemplate takes a template and returns a slice of byte slices.
// Each byte slice in the slice of bytes is the output of the template execution.
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

// outHandler loops over the slice of byte slices, outputting them to either
// the OutFile if it is specified with a flag, or to Stdout if no flag is specified.
func outHandler(data [][]byte) error {
	nl := []byte{'\n'}

	// Use stdout if no outfile is specified
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

// makeGoColName takes a column name in the format of "column_name" and converts
// it into a go styled object variable name of "ColumnName".
// makeGoColName also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
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

// makeGoVarName takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// makeGoVarName also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func makeGoVarName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {
		// Only uppercase if not a single word variable
		if s[i] == "id" && i > 0 {
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

// makeDBColName takes a table name in the format of "table_name" and a
// column name in the format of "column_name" and returns a name used in the
// `db:""` component of an object in the format of "table_name_column_name"
func makeDBColName(tableName, colName string) string {
	return tableName + "_" + colName
}
