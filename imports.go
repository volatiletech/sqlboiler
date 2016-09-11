package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/vattle/sqlboiler/bdb"
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
		`"fmt"`,
		`"strings"`,
		`"database/sql"`,
		`"reflect"`,
		`"sync"`,
		`"time"`,
	},
	thirdParty: importList{
		`"github.com/pkg/errors"`,
		`"github.com/vattle/sqlboiler/boil"`,
		`"github.com/vattle/sqlboiler/boil/qm"`,
		`"github.com/vattle/sqlboiler/strmangle"`,
	},
}

var defaultSingletonTemplateImports = map[string]imports{
	"boil_queries": {
		thirdParty: importList{
			`"github.com/vattle/sqlboiler/boil"`,
			`"github.com/vattle/sqlboiler/boil/qm"`,
		},
	},
	"boil_types": {
		thirdParty: importList{
			`"github.com/pkg/errors"`,
			`"github.com/vattle/sqlboiler/strmangle"`,
		},
	},
}

var defaultTestTemplateImports = imports{
	standard: importList{
		`"testing"`,
		`"reflect"`,
	},
	thirdParty: importList{
		`"github.com/vattle/sqlboiler/boil"`,
		`"github.com/vattle/sqlboiler/boil/randomize"`,
		`"github.com/vattle/sqlboiler/strmangle"`,
	},
}

var defaultSingletonTestTemplateImports = map[string]imports{
	"boil_viper_test": {
		standard: importList{
			`"database/sql"`,
			`"os"`,
			`"path/filepath"`,
		},
		thirdParty: importList{
			`"github.com/spf13/viper"`,
		},
	},
	"boil_queries_test": {
		standard: importList{
			`"crypto/md5"`,
			`"fmt"`,
			`"os"`,
			`"strconv"`,
			`"math/rand"`,
		},
		thirdParty: importList{
			`"github.com/vattle/sqlboiler/boil"`,
		},
	},
	"boil_suites_test": {
		standard: importList{
			`"testing"`,
		},
	},
}

var defaultTestMainImports = map[string]imports{
	"postgres": {
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
			`"github.com/kat-co/vala"`,
			`"github.com/pkg/errors"`,
			`"github.com/spf13/viper"`,
			`"github.com/vattle/sqlboiler/boil"`,
			`"github.com/vattle/sqlboiler/bdb/drivers"`,
			`_ "github.com/lib/pq"`,
		},
	},
	"mysql": {
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
			`"github.com/kat-co/vala"`,
			`"github.com/pkg/errors"`,
			`"github.com/spf13/viper"`,
			`"github.com/vattle/sqlboiler/boil"`,
			`"github.com/vattle/sqlboiler/bdb/drivers"`,
			`_ "github.com/go-mysql-driver/mysql"`,
		},
	},
}

// importsBasedOnType imports are only included in the template output if the
// database requires one of the following special types. Check
// TranslateColumnType to see the type assignments.
var importsBasedOnType = map[string]imports{
	"null.Float32": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Float64": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Int": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Int8": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Int16": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Int32": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Int64": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Uint": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Uint8": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Uint16": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Uint32": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Uint64": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.String": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Bool": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Time": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.JSON": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"null.Bytes": {
		thirdParty: importList{`"gopkg.in/nullbio/null.v5"`},
	},
	"time.Time": {
		standard: importList{`"time"`},
	},
	"types.JSON": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
	"types.BytesArray": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
	"types.Int64Array": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
	"types.Float64Array": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
	"types.BoolArray": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
	"types.Hstore": {
		thirdParty: importList{`"github.com/vattle/sqlboiler/boil/types"`},
	},
}
