package cmds

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

// generateTemplate generates the template associated to the passed in command name.
func generateTemplate(commandName string) [][]byte {
	var template *template.Template

	// Find the template that matches the passed in command name
	for _, t := range templates {
		if t.Name() == commandName+".tpl" {
			template = t
			break
		}
	}

	if template == nil {
		errorQuit(fmt.Errorf("Unable to find the template: %s", commandName+".tpl"))
	}

	outputs, err := processTemplate(template)
	if err != nil {
		errorQuit(fmt.Errorf("Unable to process the template: %s", err))
	}

	return outputs
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

		if s[i] == "id" && i > 0 {
			s[i] = "ID"
			continue
		}

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

// makeGoInsertParamNames takes a []DBTable and returns a comma seperated
// list of parameter names for the insert statement template.
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

// makeGoInsertParamFlags takes a []DBTable and returns a comma seperated
// list of parameter flags for the insert statement template.
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

// makeSelectParamNames takes a []DBTable and returns a comma seperated
// list of parameter names with for the select statement template.
// It also uses the table name to generate the "AS" part of the statement, for
// example: var_name AS table_name_var_name, ...
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
