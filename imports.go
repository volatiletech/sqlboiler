package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/nullbio/sqlboiler/bdb"
)

// imports defines the optional standard imports and
// thirdParty imports (from github for example)
type imports struct {
	standard   importList
	thirdParty importList
}

// importList is a list of import names
type importList []string

func (i importList) Len() int {
	return len(i)
}

func (i importList) Swap(k, j int) {
	i[k], i[j] = i[j], i[k]
}

func (i importList) Less(k, j int) bool {
	res := strings.Compare(strings.TrimLeft(i[k], "_ "), strings.TrimLeft(i[j], "_ "))
	if res <= 0 {
		return true
	}

	return false
}

func combineImports(a, b imports) imports {
	var c imports

	c.standard = removeDuplicates(combineStringSlices(a.standard, b.standard))
	c.thirdParty = removeDuplicates(combineStringSlices(a.thirdParty, b.thirdParty))

	sort.Sort(c.standard)
	sort.Sort(c.thirdParty)

	return c
}

func combineTypeImports(a imports, b map[string]imports, columns []bdb.Column) imports {
	tmpImp := imports{
		standard:   make(importList, len(a.standard)),
		thirdParty: make(importList, len(a.thirdParty)),
	}

	copy(tmpImp.standard, a.standard)
	copy(tmpImp.thirdParty, a.thirdParty)

	for _, col := range columns {
		for key, imp := range b {
			if col.Type == key {
				tmpImp.standard = append(tmpImp.standard, imp.standard...)
				tmpImp.thirdParty = append(tmpImp.thirdParty, imp.thirdParty...)
			}
		}
	}

	tmpImp.standard = removeDuplicates(tmpImp.standard)
	tmpImp.thirdParty = removeDuplicates(tmpImp.thirdParty)

	sort.Sort(tmpImp.standard)
	sort.Sort(tmpImp.thirdParty)

	return tmpImp
}

func buildImportString(imps imports) []byte {
	stdlen, thirdlen := len(imps.standard), len(imps.thirdParty)
	if stdlen+thirdlen < 1 {
		return []byte{}
	}

	if stdlen+thirdlen == 1 {
		var imp string
		if stdlen == 1 {
			imp = imps.standard[0]
		} else {
			imp = imps.thirdParty[0]
		}
		return []byte(fmt.Sprintf("import %s", imp))
	}

	buf := &bytes.Buffer{}
	buf.WriteString("import (")
	for _, std := range imps.standard {
		fmt.Fprintf(buf, "\n\t%s", std)
	}
	if stdlen != 0 && thirdlen != 0 {
		buf.WriteString("\n")
	}
	for _, third := range imps.thirdParty {
		fmt.Fprintf(buf, "\n\t%s", third)
	}
	buf.WriteString("\n)\n")

	return buf.Bytes()
}

func combineStringSlices(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	if len(a) > 0 {
		copy(c, a)
	}
	if len(b) > 0 {
		copy(c[len(a):], b)
	}

	return c
}

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}

var defaultTemplateImports = imports{
	standard: importList{
		`"errors"`,
		`"fmt"`,
		`"strings"`,
	},
	thirdParty: importList{
		`"github.com/nullbio/sqlboiler/boil"`,
		`"github.com/nullbio/sqlboiler/boil/qm"`,
	},
}

var defaultSingletonTemplateImports = map[string]imports{
	"helpers": imports{
		standard: importList{},
		thirdParty: importList{
			`"github.com/nullbio/sqlboiler/boil"`,
			`"github.com/nullbio/sqlboiler/boil/qm"`,
		},
	},
}

var defaultTestTemplateImports = imports{
	standard: importList{
		`"testing"`,
		`"reflect"`,
		`"time"`,
	},
	thirdParty: importList{
		`"gopkg.in/nullbio/null.v4"`,
		`"github.com/nullbio/sqlboiler/boil"`,
		`"github.com/nullbio/sqlboiler/boil/qm"`,
		`"github.com/nullbio/sqlboiler/strmangle"`,
	},
}

var defaultSingletonTestTemplateImports = map[string]imports{
	"main_helper_funcs": imports{
		standard: importList{
			`"database/sql"`,
			`"os"`,
			`"path/filepath"`,
		},
		thirdParty: importList{
			`"github.com/spf13/viper"`,
		},
	},
	"helper_funcs": imports{
		standard: importList{
			`"crypto/md5"`,
			`"fmt"`,
			`"os"`,
			`"strconv"`,
			`"math/rand"`,
			`"bytes"`,
		},
		thirdParty: importList{
			`"github.com/nullbio/sqlboiler/boil"`,
		},
	},
}

var defaultTestMainImports = map[string]imports{
	"postgres": imports{
		standard: importList{
			`"testing"`,
			`"os"`,
			`"os/exec"`,
			`"flag"`,
			`"fmt"`,
			`"io/ioutil"`,
			`"bytes"`,
			`"database/sql"`,
			`"path/filepath"`,
			`"time"`,
			`"math/rand"`,
		},
		thirdParty: importList{
			`"github.com/nullbio/sqlboiler/boil"`,
			`"github.com/nullbio/sqlboiler/bdb/drivers"`,
			`_ "github.com/lib/pq"`,
			`"github.com/spf13/viper"`,
			`"github.com/kat-co/vala"`,
		},
	},
}

// importsBasedOnType imports are only included in the template output if the
// database requires one of the following special types. Check
// TranslateColumnType to see the type assignments.
var importsBasedOnType = map[string]imports{
	"null.Float32": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Float64": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Int": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Int8": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Int16": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Int32": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Int64": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Uint": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Uint8": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Uint16": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Uint32": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Uint64": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.String": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Bool": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"null.Time": imports{
		thirdParty: importList{`"gopkg.in/nullbio/null.v4"`},
	},
	"time.Time": imports{
		standard: importList{`"time"`},
	},
}
