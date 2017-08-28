package boilingcore

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/bdb"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// templateData for sqlboiler templates
type templateData struct {
	Tables []bdb.Table
	Table  bdb.Table

	// Controls what names are output
	PkgName string
	Schema  string

	// Controls which code is output (mysql vs postgres ...)
	DriverName      string
	UseLastInsertID bool

	// Turn off auto timestamps or hook generation
	NoHooks          bool
	NoAutoTimestamps bool

	// Tags control which
	Tags []string

	// Generate struct tags as camelCase or snake_case
	StructTagCasing string

	// StringFuncs are usable in templates with stringMap
	StringFuncs map[string]func(string) string

	// Dialect controls quoting
	Dialect queries.Dialect
	LQ      string
	RQ      string
}

func (t templateData) Quotes(s string) string {
	return fmt.Sprintf("%s%s%s", t.LQ, s, t.RQ)
}

func (t templateData) SchemaTable(table string) string {
	return strmangle.SchemaTable(t.LQ, t.RQ, t.DriverName, t.Schema, table)
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
	pattern := filepath.Join(dir, "*.tpl")
	tpl, err := template.New("").Funcs(templateFunctions).ParseGlob(pattern)

	if err != nil {
		return nil, err
	}

	return &templateList{Template: tpl}, err
}

// loadTemplate loads a single template file
func loadTemplate(dir string, filename string) (*template.Template, error) {
	pattern := filepath.Join(dir, filename)
	tpl, err := template.New("").Funcs(templateFunctions).ParseFiles(pattern)

	if err != nil {
		return nil, err
	}

	return tpl.Lookup(filename), err
}

// replaceTemplate finds the template matching with name and replaces its
// contents with the contents of the template located at filename
func replaceTemplate(tpl *template.Template, name, filename string) error {
	if tpl == nil {
		return fmt.Errorf("template for %s is nil", name)
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "failed reading template file: %s", filename)
	}

	if tpl, err = tpl.New(name).Funcs(templateFunctions).Parse(string(b)); err != nil {
		return errors.Wrapf(err, "failed to parse template file: %s", filename)
	}

	return nil
}

// set is to stop duplication from named enums, allowing a template loop
// to keep some state
type once map[string]struct{}

func newOnce() once {
	return make(once)
}

func (o once) Has(s string) bool {
	_, ok := o[s]
	return ok
}

func (o once) Put(s string) bool {
	if _, ok := o[s]; ok {
		return false
	}

	o[s] = struct{}{}
	return true
}

// templateStringMappers are placed into the data to make it easy to use the
// stringMap function.
var templateStringMappers = map[string]func(string) string{
	// String ops
	"quoteWrap":       func(a string) string { return fmt.Sprintf(`"%s"`, a) },
	"replaceReserved": strmangle.ReplaceReservedWords,

	// Casing
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

// templateFunctions is a map of all the functions that get passed into the
// templates. If you wish to pass a new function into your own template,
// add a function pointer here.
var templateFunctions = template.FuncMap{
	// String ops
	"quoteWrap": func(s string) string { return fmt.Sprintf(`"%s"`, s) },
	"id":        strmangle.Identifier,

	// Pluralization
	"singular": strmangle.Singular,
	"plural":   strmangle.Plural,

	// Casing
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,

	// String Slice ops
	"join":               func(sep string, slice []string) string { return strings.Join(slice, sep) },
	"joinSlices":         strmangle.JoinSlices,
	"stringMap":          strmangle.StringMap,
	"prefixStringSlice":  strmangle.PrefixStringSlice,
	"containsAny":        strmangle.ContainsAny,
	"generateTags":       strmangle.GenerateTags,
	"generateIgnoreTags": strmangle.GenerateIgnoreTags,

	// Enum ops
	"parseEnumName":       strmangle.ParseEnumName,
	"parseEnumVals":       strmangle.ParseEnumVals,
	"isEnumNormal":        strmangle.IsEnumNormal,
	"shouldTitleCaseEnum": strmangle.ShouldTitleCaseEnum,
	"onceNew":             newOnce,
	"oncePut":             once.Put,
	"onceHas":             once.Has,

	// String Map ops
	"makeStringMap": strmangle.MakeStringMap,

	// Set operations
	"setInclude": strmangle.SetInclude,

	// Database related mangling
	"whereClause": strmangle.WhereClause,

	// Relationship text helpers
	"txtsFromFKey":     txtsFromFKey,
	"txtsFromOneToOne": txtsFromOneToOne,
	"txtsFromToMany":   txtsFromToMany,

	// dbdrivers ops
	"filterColumnsByAuto":    bdb.FilterColumnsByAuto,
	"filterColumnsByDefault": bdb.FilterColumnsByDefault,
	"filterColumnsByEnum":    bdb.FilterColumnsByEnum,
	"sqlColDefinitions":      bdb.SQLColDefinitions,
	"columnNames":            bdb.ColumnNames,
	"columnDBTypes":          bdb.ColumnDBTypes,
	"getTable":               bdb.GetTable,
}
