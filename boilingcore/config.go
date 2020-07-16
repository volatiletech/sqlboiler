package boilingcore

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/volatiletech/sqlboiler/v4/drivers"
	"github.com/volatiletech/sqlboiler/v4/importers"
)

// Config for the running of the commands
type Config struct {
	DriverName   string         `toml:"driver_name,omitempty" json:"driver_name,omitempty"`
	DriverConfig drivers.Config `toml:"driver_config,omitempty" json:"driver_config,omitempty"`

	PkgName           string   `toml:"pkg_name,omitempty" json:"pkg_name,omitempty"`
	OutFolder         string   `toml:"out_folder,omitempty" json:"out_folder,omitempty"`
	TemplateDirs      []string `toml:"template_dirs,omitempty" json:"template_dirs,omitempty"`
	Tags              []string `toml:"tags,omitempty" json:"tags,omitempty"`
	Replacements      []string `toml:"replacements,omitempty" json:"replacements,omitempty"`
	Debug             bool     `toml:"debug,omitempty" json:"debug,omitempty"`
	AddGlobal         bool     `toml:"add_global,omitempty" json:"add_global,omitempty"`
	AddPanic          bool     `toml:"add_panic,omitempty" json:"add_panic,omitempty"`
	AddSoftDeletes    bool     `toml:"add_soft_deletes,omitempty" json:"add_soft_deletes,omitempty"`
	NoContext         bool     `toml:"no_context,omitempty" json:"no_context,omitempty"`
	NoTests           bool     `toml:"no_tests,omitempty" json:"no_tests,omitempty"`
	NoHooks           bool     `toml:"no_hooks,omitempty" json:"no_hooks,omitempty"`
	NoAutoTimestamps  bool     `toml:"no_auto_timestamps,omitempty" json:"no_auto_timestamps,omitempty"`
	NoRowsAffected    bool     `toml:"no_rows_affected,omitempty" json:"no_rows_affected,omitempty"`
	NoDriverTemplates bool     `toml:"no_driver_templates,omitempty" json:"no_driver_templates,omitempty"`
	NoBackReferencing bool     `toml:"no_back_reference,omitempty" json:"no_back_reference,omitempty"`
	Wipe              bool     `toml:"wipe,omitempty" json:"wipe,omitempty"`
	StructTagCasing   string   `toml:"struct_tag_casing,omitempty" json:"struct_tag_casing,omitempty"`
	RelationTag       string   `toml:"relation_tag,omitempty" json:"relation_tag,omitempty"`
	TagIgnore         []string `toml:"tag_ignore,omitempty" json:"tag_ignore,omitempty"`

	Imports importers.Collection `toml:"imports,omitempty" json:"imports,omitempty"`

	Aliases      Aliases       `toml:"aliases,omitempty" json:"aliases,omitempty"`
	TypeReplaces []TypeReplace `toml:"type_replaces,omitempty" json:"type_replaces,omitempty"`

	Version string `toml:"version" json:"version"`
}

// TypeReplace replaces a column type with something else
type TypeReplace struct {
	Tables  []string       `toml:"tables,omitempty" json:"tables,omitempty"`
	Match   drivers.Column `toml:"match,omitempty" json:"match,omitempty"`
	Replace drivers.Column `toml:"replace,omitempty" json:"replace,omitempty"`
	Imports importers.Set  `toml:"imports,omitempty" json:"imports,omitempty"`
}

// OutputDirDepth returns depth of output directory
func (c *Config) OutputDirDepth() int {
	d := filepath.ToSlash(filepath.Clean(c.OutFolder))
	if d == "." {
		return 0
	}

	return strings.Count(d, "/") + 1
}

// ConvertAliases is necessary because viper
//
// It also supports two different syntaxes, because of viper:
//
//   [aliases.tables.table_name]
//   fields... = "values"
//     [aliases.tables.columns]
//     colname = "alias"
//     [aliases.tables.relationships.fkey_name]
//     local   = "x"
//     foreign = "y"
//
// Or alternatively (when toml key names or viper's
// lowercasing of key names gets in the way):
//
//   [[aliases.tables]]
//   name = "table_name"
//   fields... = "values"
//     [[aliases.tables.columns]]
//     name  = "colname"
//     alias = "alias"
//     [[aliases.tables.relationships]]
//     name    = "fkey_name"
//     local   = "x"
//     foreign = "y"
func ConvertAliases(i interface{}) (a Aliases) {
	if i == nil {
		return a
	}

	topLevel := cast.ToStringMap(i)

	tablesIntf := topLevel["tables"]

	iterateMapOrSlice(tablesIntf, func(name string, tIntf interface{}) {
		if a.Tables == nil {
			a.Tables = make(map[string]TableAlias)
		}

		t := cast.ToStringMap(tIntf)

		var ta TableAlias

		if s := t["name_singular"]; s != nil {
			ta.NameSingular = s.(string)
		}
		if s := t["up_plural"]; s != nil {
			ta.UpPlural = s.(string)
		}
		if s := t["up_singular"]; s != nil {
			ta.UpSingular = s.(string)
		}
		if s := t["down_plural"]; s != nil {
			ta.DownPlural = s.(string)
		}
		if s := t["down_singular"]; s != nil {
			ta.DownSingular = s.(string)
		}

		if colsIntf, ok := t["columns"]; ok {
			ta.Columns = make(map[string]string)

			iterateMapOrSlice(colsIntf, func(name string, colIntf interface{}) {
				var alias string
				switch col := colIntf.(type) {
				case map[string]interface{}, map[interface{}]interface{}:
					cmap := cast.ToStringMap(colIntf)
					alias = cmap["alias"].(string)
				case string:
					alias = col
				}
				ta.Columns[name] = alias
			})
		}

		relationshipsIntf, ok := t["relationships"]
		if ok {
			iterateMapOrSlice(relationshipsIntf, func(name string, rIntf interface{}) {
				if ta.Relationships == nil {
					ta.Relationships = make(map[string]RelationshipAlias)
				}

				var ra RelationshipAlias
				rel := cast.ToStringMap(rIntf)

				if s := rel["local"]; s != nil {
					ra.Local = s.(string)
				}
				if s := rel["foreign"]; s != nil {
					ra.Foreign = s.(string)
				}

				ta.Relationships[name] = ra
			})
		}

		a.Tables[name] = ta
	})

	return a
}

func iterateMapOrSlice(mapOrSlice interface{}, fn func(name string, obj interface{})) {
	switch t := mapOrSlice.(type) {
	case map[string]interface{}, map[interface{}]interface{}:
		tmap := cast.ToStringMap(mapOrSlice)
		for name, table := range tmap {
			fn(name, table)
		}
	case []interface{}:
		for _, intf := range t {
			obj := cast.ToStringMap(intf)
			name := obj["name"].(string)
			fn(name, intf)
		}
	}
}

// ConvertTypeReplace is necessary because viper
func ConvertTypeReplace(i interface{}) []TypeReplace {
	if i == nil {
		return nil
	}

	intfArray := i.([]interface{})
	var replaces []TypeReplace
	for _, r := range intfArray {
		replaceIntf := cast.ToStringMap(r)
		replace := TypeReplace{}

		if replaceIntf["match"] == nil || replaceIntf["replace"] == nil {
			panic("replace types must specify both match and replace")
		}

		replace.Match = columnFromInterface(replaceIntf["match"])
		replace.Replace = columnFromInterface(replaceIntf["replace"])

		replace.Tables = tablesOfTypeReplace(replaceIntf["match"])

		if imps := replaceIntf["imports"]; imps != nil {
			imps = cast.ToStringMap(imps)
			var err error
			replace.Imports, err = importers.SetFromInterface(imps)
			if err != nil {
				panic(err)
			}
		}

		replaces = append(replaces, replace)
	}

	return replaces
}

func tablesOfTypeReplace(i interface{}) []string {
	tables := []string{}

	m := cast.ToStringMap(i)
	if s := m["tables"]; s != nil {
		tables = cast.ToStringSlice(s)
	}

	return tables
}

func columnFromInterface(i interface{}) (col drivers.Column) {
	m := cast.ToStringMap(i)
	if s := m["name"]; s != nil {
		col.Name = s.(string)
	}
	if s := m["type"]; s != nil {
		col.Type = s.(string)
	}
	if s := m["db_type"]; s != nil {
		col.DBType = s.(string)
	}
	if s := m["udt_name"]; s != nil {
		col.UDTName = s.(string)
	}
	if s := m["full_db_type"]; s != nil {
		col.FullDBType = s.(string)
	}
	if s := m["arr_type"]; s != nil {
		col.ArrType = new(string)
		*col.ArrType = s.(string)
	}
	if s := m["domain_name"]; s != nil {
		col.DomainName = new(string)
		*col.DomainName = s.(string)
	}
	if b := m["auto_generated"]; b != nil {
		col.AutoGenerated = b.(bool)
	}
	if b := m["nullable"]; b != nil {
		col.Nullable = b.(bool)
	}

	return col
}
