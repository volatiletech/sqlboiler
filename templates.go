package main

import (
	"fmt"
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

	StringFuncs map[string]func(string) string
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

// templateStringMappers are placed into the data to make it easy to use the
// stringMap function.
var templateStringMappers = map[string]func(string) string{
	// String ops
	"quoteWrap": func(a string) string { return fmt.Sprintf(`"%s"`, a) },

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"toLower":   strings.ToLower,
	"toUpper":   strings.ToUpper,
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

// templateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var templateFunctions = template.FuncMap{
	// String ops
	"substring": strmangle.Substring,
	"remove":    func(rem string, str string) string { return strings.Replace(str, rem, "", -1) },
	"prefix":    func(add string, str string) string { return fmt.Sprintf("%s%s", add, str) },
	"quoteWrap": func(a string) string { return fmt.Sprintf(`"%s"`, a) },

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"toLower":   strings.ToLower,
	"toUpper":   strings.ToUpper,
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,

	// String Slice ops
	"join":              func(sep string, slice []string) string { return strings.Join(slice, sep) },
	"stringMap":         strmangle.StringMap,
	"hasElement":        strmangle.HasElement,
	"prefixStringSlice": strmangle.PrefixStringSlice,

	// Database related mangling
	"wherePrimaryKey": strmangle.WherePrimaryKey,

	// dbdrivers ops
	"driverUsesLastInsertID":       strmangle.DriverUsesLastInsertID,
	"filterColumnsByDefault":       strmangle.FilterColumnsByDefault,
	"filterColumnsByAutoIncrement": strmangle.FilterColumnsByAutoIncrement,
	"autoIncPrimaryKey":            strmangle.AutoIncPrimaryKey,
	"primaryKeyFuncSig":            strmangle.PrimaryKeyFuncSig,
	"columnNames":                  strmangle.ColumnNames,
	"makeDBName":                   strmangle.MakeDBName,
}
