package cmds

import (
	"fmt"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/pobri19/sqlboiler/strmangle"
)

// sqlBoilerTypeImports imports are only included in the template output if the database
// requires one of the following special types. Check TranslateColumnType to see the type assignments.
var sqlBoilerTypeImports = map[string]imports{
	"null.Int": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.String": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Bool": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Float": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"null.Time": imports{
		thirdparty: importList{`"gopkg.in/guregu/null.v3"`},
	},
	"time.Time": imports{
		standard: importList{`"time"`},
	},
}

// sqlBoilerImports defines the list of default template imports.
var sqlBoilerImports = imports{
	standard: importList{
		`"errors"`,
		`"fmt"`,
		`"strings"`,
	},
	thirdparty: importList{
		`"github.com/pobri19/sqlboiler/boil"`,
		`"github.com/pobri19/sqlboiler/boil/qs"`,
	},
}

var sqlBoilerSinglesImports = map[string]imports{
	"helpers": imports{
		standard: importList{},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
			`"github.com/pobri19/sqlboiler/boil/qs"`,
		},
	},
}

// sqlBoilerTestImports defines the list of default test template imports.
var sqlBoilerTestImports = imports{
	standard: importList{
		`"testing"`,
	},
}

var sqlBoilerTestMainImports = map[string]imports{
	"postgres": imports{
		standard: importList{
			`"testing"`,
			`"os"`,
			`"os/exec"`,
			`"fmt"`,
			`"io/ioutil"`,
			`"bytes"`,
			`"database/sql"`,
			`"time"`,
			`"math/rand"`,
		},
		thirdparty: importList{
			`"github.com/BurntSushi/toml"`,
			`_ "github.com/lib/pq"`,
		},
	},
}

// sqlBoilerTemplateFuncs is a map of all the functions that get passed into the templates.
// If you wish to pass a new function into your own template, add a pointer to it here.
var sqlBoilerTemplateFuncs = template.FuncMap{
	"singular":               strmangle.Singular,
	"plural":                 strmangle.Plural,
	"titleCase":              strmangle.TitleCase,
	"titleCaseSingular":      strmangle.TitleCaseSingular,
	"titleCasePlural":        strmangle.TitleCasePlural,
	"camelCase":              strmangle.CamelCase,
	"camelCaseSingular":      strmangle.CamelCaseSingular,
	"camelCasePlural":        strmangle.CamelCasePlural,
	"camelCaseCommaList":     strmangle.CamelCaseCommaList,
	"commaList":              strmangle.CommaList,
	"makeDBName":             strmangle.MakeDBName,
	"selectParamNames":       strmangle.SelectParamNames,
	"insertParamNames":       strmangle.InsertParamNames,
	"insertParamFlags":       strmangle.InsertParamFlags,
	"insertParamVariables":   strmangle.InsertParamVariables,
	"scanParamNames":         strmangle.ScanParamNames,
	"hasPrimaryKey":          strmangle.HasPrimaryKey,
	"primaryKeyFuncSig":      strmangle.PrimaryKeyFuncSig,
	"wherePrimaryKey":        strmangle.WherePrimaryKey,
	"paramsPrimaryKey":       strmangle.ParamsPrimaryKey,
	"primaryKeyFlagIndex":    strmangle.PrimaryKeyFlagIndex,
	"updateParamNames":       strmangle.UpdateParamNames,
	"updateParamVariables":   strmangle.UpdateParamVariables,
	"primaryKeyStrList":      strmangle.PrimaryKeyStrList,
	"supportsResultObject":   strmangle.SupportsResultObject,
	"filterColumnsByDefault": strmangle.FilterColumnsByDefault,
	"autoIncPrimaryKey":      strmangle.AutoIncPrimaryKey,
}

// LoadConfigFile loads the toml config file into the cfg object
func (c *CmdData) LoadConfigFile(filename string) error {
	cfg := &Config{}

	_, err := toml.DecodeFile(filename, &cfg)

	if os.IsNotExist(err) {
		return fmt.Errorf("Failed to find the toml configuration file %s: %s", filename, err)
	}

	if err != nil {
		return fmt.Errorf("Failed to decode toml configuration file: %s", err)
	}

	c.Config = cfg
	return nil
}
