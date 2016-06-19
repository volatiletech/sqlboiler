package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/nullbio/sqlboiler/dbdrivers"
	"github.com/nullbio/sqlboiler/strmangle"
)

// templateData for sqlboiler templates
type templateData struct {
	Tables     []dbdrivers.Table
	Table      dbdrivers.Table
	DriverName string
	PkgName    string
}

type templateList []*template.Template

func (t templateList) Len() int {
	return len(t)
}

func (t templateList) Swap(k, j int) {
	t[k], t[j] = t[j], t[k]
}

func (t templateList) Less(k, j int) bool {
	// Make sure "struct" goes to the front
	if t[k].Name() == "struct.tpl" {
		return true
	}

	res := strings.Compare(t[k].Name(), t[j].Name())
	if res <= 0 {
		return true
	}

	return false
}

// loadTemplates loads all of the template files in the specified directory.
func loadTemplates(dir string) (templateList, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, dir, "*.tpl")
	tpl, err := template.New("").Funcs(templateFunctions).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	templates := templateList(tpl.Templates())
	sort.Sort(templates)

	return templates, err
}

// loadTemplate loads a single template file.
func loadTemplate(dir string, filename string) (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, dir, filename)
	tpl, err := template.New("").Funcs(templateFunctions).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Lookup(filename), err
}

// templateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var templateFunctions = template.FuncMap{
	"tolower":                      strings.ToLower,
	"toupper":                      strings.ToUpper,
	"substring":                    strmangle.Substring,
	"singular":                     strmangle.Singular,
	"plural":                       strmangle.Plural,
	"titleCase":                    strmangle.TitleCase,
	"titleCaseSingular":            strmangle.TitleCaseSingular,
	"titleCasePlural":              strmangle.TitleCasePlural,
	"titleCaseCommaList":           strmangle.TitleCaseCommaList,
	"camelCase":                    strmangle.CamelCase,
	"camelCaseSingular":            strmangle.CamelCaseSingular,
	"camelCasePlural":              strmangle.CamelCasePlural,
	"camelCaseCommaList":           strmangle.CamelCaseCommaList,
	"columnsToStrings":             strmangle.ColumnsToStrings,
	"commaList":                    strmangle.CommaList,
	"makeDBName":                   strmangle.MakeDBName,
	"selectParamNames":             strmangle.SelectParamNames,
	"insertParamNames":             strmangle.InsertParamNames,
	"insertParamFlags":             strmangle.InsertParamFlags,
	"insertParamVariables":         strmangle.InsertParamVariables,
	"scanParamNames":               strmangle.ScanParamNames,
	"hasPrimaryKey":                strmangle.HasPrimaryKey,
	"primaryKeyFuncSig":            strmangle.PrimaryKeyFuncSig,
	"wherePrimaryKey":              strmangle.WherePrimaryKey,
	"paramsPrimaryKey":             strmangle.ParamsPrimaryKey,
	"primaryKeyFlagIndex":          strmangle.PrimaryKeyFlagIndex,
	"updateParamNames":             strmangle.UpdateParamNames,
	"updateParamVariables":         strmangle.UpdateParamVariables,
	"supportsResultObject":         strmangle.SupportsResultObject,
	"filterColumnsByDefault":       strmangle.FilterColumnsByDefault,
	"filterColumnsByAutoIncrement": strmangle.FilterColumnsByAutoIncrement,
	"autoIncPrimaryKey":            strmangle.AutoIncPrimaryKey,
	"addID":                        strmangle.AddID,
	"removeID":                     strmangle.RemoveID,

	"randDBStruct":      strmangle.RandDBStruct,
	"randDBStructSlice": strmangle.RandDBStructSlice,
}
