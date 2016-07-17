package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/nullbio/sqlboiler/bdb"
	"github.com/nullbio/sqlboiler/strmangle"
)

// templateData for sqlboiler templates
type templateData struct {
	Tables     []bdb.Table
	Table      bdb.Table
	DriverName string
	PkgName    string

	StringFuncs map[string]func(string) string
}

type templateList struct {
	*template.Template
}

type templateNameList []string

func (t templateNameList) Len() int {
	return len(t)
}

func (t templateNameList) Swap(k, j int) {
	t[k], t[j] = t[j], t[k]
}

func (t templateNameList) Less(k, j int) bool {
	// Make sure "struct" goes to the front
	if t[k] == "struct.tpl" {
		return true
	}

	res := strings.Compare(t[k], t[j])
	if res <= 0 {
		return true
	}

	return false
}

// Templates returns the name of all the templates defined in the template list
func (t templateList) Templates() []string {
	tplList := t.Template.Templates()

	if len(tplList) == 0 {
		return nil
	}

	ret := make([]string, 0, len(tplList))
	for _, tpl := range tplList {
		if name := tpl.Name(); strings.HasSuffix(name, ".tpl") {
			ret = append(ret, name)
		}
	}

	sort.Sort(templateNameList(ret))

	return ret
}

// loadTemplates loads all of the template files in the specified directory.
func loadTemplates(dir string) (*templateList, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(wd, dir, "*.tpl")
	tpl, err := template.New("").Funcs(templateFunctions).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	return &templateList{Template: tpl}, err
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
	"trimPrefix": func(pre, str string) string { return strings.TrimPrefix(str, pre) },
	"remove":     func(rem, str string) string { return strings.Replace(str, rem, "", -1) },
	"replace":    func(rep, with, str string) string { return strings.Replace(str, rep, with, -1) },
	"prefix":     func(add, str string) string { return fmt.Sprintf("%s%s", add, str) },
	"quoteWrap":  func(a string) string { return fmt.Sprintf(`"%s"`, a) },
	"substring":  strmangle.Substring,
	"id":         strmangle.Identifier,

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
	"joinSlices":        strmangle.JoinSlices,
	"stringMap":         strmangle.StringMap,
	"hasElement":        strmangle.HasElement,
	"prefixStringSlice": strmangle.PrefixStringSlice,

	// String Map ops
	"makeStringMap": strmangle.MakeStringMap,

	// Database related mangling
	"whereClause": strmangle.WhereClause,

	// Text helpers
	"textsFromForeignKey":   textsFromForeignKey,
	"textsFromRelationship": textsFromRelationship,

	// dbdrivers ops
	"driverUsesLastInsertID":       strmangle.DriverUsesLastInsertID,
	"makeDBName":                   strmangle.MakeDBName,
	"filterColumnsByDefault":       bdb.FilterColumnsByDefault,
	"filterColumnsBySimpleDefault": bdb.FilterColumnsBySimpleDefault,
	"filterColumnsByAutoIncrement": bdb.FilterColumnsByAutoIncrement,
	"autoIncPrimaryKey":            bdb.AutoIncPrimaryKey,
	"sqlColDefinitions":            bdb.SQLColDefinitions,
	"sqlColDefStrings":             bdb.SQLColDefStrings,
	"columnNames":                  bdb.ColumnNames,
	"columnTypes":                  bdb.ColumnTypes,
	"columnDBTypes":                bdb.ColumnDBTypes,
	"defaultValues":                bdb.DefaultValues,
}
