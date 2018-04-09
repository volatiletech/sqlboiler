package boilingcore

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ann-kilzer/sqlboiler/bdb"
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

type mapImports map[string]imports

type importer struct {
	Standard     imports
	TestStandard imports

	Singleton     mapImports
	TestSingleton mapImports

	TestMain mapImports

	BasedOnType mapImports
}

// newImporter returns an importer struct with default import values
func newImporter() importer {
	var imp importer

	imp.Standard = imports{
		standard: importList{
			`"bytes"`,
			`"database/sql"`,
			`"fmt"`,
			`"reflect"`,
			`"strings"`,
			`"sync"`,
			`"time"`,
		},
		thirdParty: importList{
			`"github.com/pkg/errors"`,
			`"github.com/ann-kilzer/sqlboiler/boil"`,
			`"github.com/ann-kilzer/sqlboiler/queries"`,
			`"github.com/ann-kilzer/sqlboiler/queries/qm"`,
			`"github.com/ann-kilzer/sqlboiler/strmangle"`,
		},
	}

	imp.Singleton = mapImports{
		"boil_queries": {
			thirdParty: importList{
				`"github.com/ann-kilzer/sqlboiler/boil"`,
				`"github.com/ann-kilzer/sqlboiler/queries"`,
				`"github.com/ann-kilzer/sqlboiler/queries/qm"`,
			},
		},
		"boil_types": {
			thirdParty: importList{
				`"github.com/pkg/errors"`,
				`"github.com/ann-kilzer/sqlboiler/strmangle"`,
			},
		},
	}

	imp.TestStandard = imports{
		standard: importList{
			`"bytes"`,
			`"reflect"`,
			`"testing"`,
		},
		thirdParty: importList{
			`"github.com/ann-kilzer/sqlboiler/boil"`,
			`"github.com/ann-kilzer/sqlboiler/randomize"`,
			`"github.com/ann-kilzer/sqlboiler/strmangle"`,
		},
	}

	imp.TestSingleton = mapImports{
		"boil_main_test": {
			standard: importList{
				`"database/sql"`,
				`"flag"`,
				`"fmt"`,
				`"math/rand"`,
				`"os"`,
				`"path/filepath"`,
				`"testing"`,
				`"time"`,
			},
			thirdParty: importList{
				`"github.com/kat-co/vala"`,
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/ann-kilzer/sqlboiler/boil"`,
			},
		},
		"boil_queries_test": {
			standard: importList{
				`"bytes"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"math/rand"`,
				`"regexp"`,
			},
			thirdParty: importList{
				`"github.com/ann-kilzer/sqlboiler/boil"`,
			},
		},
		"boil_suites_test": {
			standard: importList{
				`"testing"`,
			},
		},
	}

	imp.TestMain = mapImports{
		"postgres": {
			standard: importList{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			thirdParty: importList{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/ann-kilzer/sqlboiler/bdb/drivers"`,
				`"github.com/ann-kilzer/sqlboiler/randomize"`,
				`_ "github.com/lib/pq"`,
			},
		},
		"mysql": {
			standard: importList{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			thirdParty: importList{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/ann-kilzer/sqlboiler/bdb/drivers"`,
				`"github.com/ann-kilzer/sqlboiler/randomize"`,
				`_ "github.com/go-sql-driver/mysql"`,
			},
		},
		"mssql": {
			standard: importList{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			thirdParty: importList{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/ann-kilzer/sqlboiler/bdb/drivers"`,
				`"github.com/ann-kilzer/sqlboiler/randomize"`,
				`_ "github.com/denisenkom/go-mssqldb"`,
			},
		},
	}

	// basedOnType imports are only included in the template output if the
	// database requires one of the following special types. Check
	// TranslateColumnType to see the type assignments.
	imp.BasedOnType = mapImports{
		"null.Float32": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Float64": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int8": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int16": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int32": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int64": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint8": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint16": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint32": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint64": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.String": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Bool": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Time": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.JSON": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Bytes": {
			thirdParty: importList{`"gopkg.in/volatiletech/null.v6"`},
		},
		"time.Time": {
			standard: importList{`"time"`},
		},
		"types.JSON": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.BytesArray": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.Int64Array": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.Float64Array": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.BoolArray": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.StringArray": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
		"types.Hstore": {
			thirdParty: importList{`"github.com/ann-kilzer/sqlboiler/types"`},
		},
	}

	return imp
}

// Remove an import matching the match string under the specified key.
// Remove will search both standard and thirdParty import lists for a match.
func (m mapImports) Remove(key string, match string) {
	mp := m[key]
	for idx := 0; idx < len(mp.standard); idx++ {
		if mp.standard[idx] == match {
			mp.standard[idx] = mp.standard[len(mp.standard)-1]
			mp.standard = mp.standard[:len(mp.standard)-1]
			break
		}
	}
	for idx := 0; idx < len(mp.thirdParty); idx++ {
		if mp.thirdParty[idx] == match {
			mp.thirdParty[idx] = mp.thirdParty[len(mp.thirdParty)-1]
			mp.thirdParty = mp.thirdParty[:len(mp.thirdParty)-1]
			break
		}
	}

	// delete the key and return if both import lists are empty
	if len(mp.thirdParty) == 0 && len(mp.standard) == 0 {
		delete(m, key)
		return
	}

	m[key] = mp
}

// Add an import under the specified key. If the key does not exist, it
// will be created.
func (m mapImports) Add(key string, value string, thirdParty bool) {
	mp := m[key]
	if thirdParty {
		mp.thirdParty = append(mp.thirdParty, value)
	} else {
		mp.standard = append(mp.standard, value)
	}

	m[key] = mp
}

// Remove an import matching the match string under the specified key.
// Remove will search both standard and thirdParty import lists for a match.
func (i *imports) Remove(match string) {
	for idx := 0; idx < len(i.standard); idx++ {
		if i.standard[idx] == match {
			i.standard[idx] = i.standard[len(i.standard)-1]
			i.standard = i.standard[:len(i.standard)-1]
			break
		}
	}
	for idx := 0; idx < len(i.thirdParty); idx++ {
		if i.thirdParty[idx] == match {
			i.thirdParty[idx] = i.thirdParty[len(i.thirdParty)-1]
			i.thirdParty = i.thirdParty[:len(i.thirdParty)-1]
			break
		}
	}
}

// Add an import under the specified key. If the key does not exist, it
// will be created.
func (i *imports) Add(value string, thirdParty bool) {
	if thirdParty {
		i.thirdParty = append(i.thirdParty, value)
	} else {
		i.standard = append(i.standard, value)
	}
}
