// Package importers helps with dynamic imports for templating
package importers

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Collection of imports for various templating purposes
type Collection struct {
	All  Set `toml:"all"`
	Test Set `toml:"test"`

	Singleton     Map `toml:"singleton"`
	TestSingleton Map `toml:"test_singleton"`

	TestMain Map `toml:"test_main"`

	BasedOnType Map `toml:"based_on_type"`
}

// Set defines the optional standard imports and
// thirdParty imports (from github for example)
type Set struct {
	Standard   List `toml:"standard"`
	ThirdParty List `toml:"third_party"`
}

// Format the set into Go syntax (compatible with go imports)
func (s Set) Format() []byte {
	stdlen, thirdlen := len(s.Standard), len(s.ThirdParty)
	if stdlen+thirdlen < 1 {
		return []byte{}
	}

	if stdlen+thirdlen == 1 {
		var imp string
		if stdlen == 1 {
			imp = s.Standard[0]
		} else {
			imp = s.ThirdParty[0]
		}
		return []byte(fmt.Sprintf("import %s", imp))
	}

	buf := &bytes.Buffer{}
	buf.WriteString("import (")
	for _, std := range s.Standard {
		fmt.Fprintf(buf, "\n\t%s", std)
	}
	if stdlen != 0 && thirdlen != 0 {
		buf.WriteString("\n")
	}
	for _, third := range s.ThirdParty {
		fmt.Fprintf(buf, "\n\t%s", third)
	}
	buf.WriteString("\n)\n")

	return buf.Bytes()
}

// SetFromInterface creates a set from a theoretical map[string]interface{}.
// This is to load from a loosely defined configuration file.
func SetFromInterface(intf interface{}) (Set, error) {
	s := Set{}

	setIntf, ok := intf.(map[string]interface{})
	if !ok {
		return s, errors.New("import set should be map[string]interface{}")
	}

	standardIntf, ok := setIntf["standard"]
	if ok {
		standardsIntf, ok := standardIntf.([]interface{})
		if !ok {
			return s, errors.New("import set standards must be an slice")
		}

		s.Standard = List{}
		for i, intf := range standardsIntf {
			str, ok := intf.(string)
			if !ok {
				return s, errors.Errorf("import set standard slice element %d (%+v) must be string", i, s)
			}
			s.Standard = append(s.Standard, str)
		}
	}

	thirdPartyIntf, ok := setIntf["third_party"]
	if ok {
		thirdPartysIntf, ok := thirdPartyIntf.([]interface{})
		if !ok {
			return s, errors.New("import set third_party must be an slice")
		}

		s.ThirdParty = List{}
		for i, intf := range thirdPartysIntf {
			str, ok := intf.(string)
			if !ok {
				return s, errors.Errorf("import set third party slice element %d (%+v) must be string", i, intf)
			}
			s.ThirdParty = append(s.ThirdParty, str)
		}
	}

	return s, nil
}

// Map of type -> import
type Map map[string]Set

// MapFromInterface creates a Map from a theoretical map[string]interface{}.
// This is to load from a loosely defined configuration file.
func MapFromInterface(intf interface{}) (Map, error) {
	m := Map{}

	mapIntf, ok := intf.(map[string]interface{})
	if !ok {
		return m, errors.New("import map should be map[string]interface{}")
	}

	for k, v := range mapIntf {
		s, err := SetFromInterface(v)
		if err != nil {
			return nil, err
		}

		m[k] = s
	}

	return m, nil
}

// List of imports
type List []string

// Len implements sort.Interface.Len
func (l List) Len() int {
	return len(l)
}

// Swap implements sort.Interface.Swap
func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less implements sort.Interface.Less
func (l List) Less(i, j int) bool {
	res := strings.Compare(strings.TrimLeft(l[i], "_ "), strings.TrimLeft(l[j], "_ "))
	if res <= 0 {
		return true
	}

	return false
}

// NewDefaultImports returns a default Imports struct.
func NewDefaultImports() Collection {
	var col Collection

	col.All = Set{
		Standard: List{
			`"bytes"`,
			`"database/sql"`,
			`"fmt"`,
			`"reflect"`,
			`"strings"`,
			`"sync"`,
			`"time"`,
		},
		ThirdParty: List{
			`"github.com/pkg/errors"`,
			`"github.com/volatiletech/sqlboiler/boil"`,
			`"github.com/volatiletech/sqlboiler/queries"`,
			`"github.com/volatiletech/sqlboiler/queries/qm"`,
			`"github.com/volatiletech/sqlboiler/strmangle"`,
		},
	}

	col.Singleton = Map{
		"boil_queries": {
			ThirdParty: List{
				`"github.com/volatiletech/sqlboiler/boil"`,
				`"github.com/volatiletech/sqlboiler/drivers"`,
				`"github.com/volatiletech/sqlboiler/queries"`,
				`"github.com/volatiletech/sqlboiler/queries/qm"`,
			},
		},
		"boil_types": {
			ThirdParty: List{
				`"github.com/pkg/errors"`,
				`"github.com/volatiletech/sqlboiler/strmangle"`,
			},
		},
	}

	col.Test = Set{
		Standard: List{
			`"bytes"`,
			`"reflect"`,
			`"testing"`,
		},
		ThirdParty: List{
			`"github.com/volatiletech/sqlboiler/boil"`,
			`"github.com/volatiletech/sqlboiler/randomize"`,
			`"github.com/volatiletech/sqlboiler/strmangle"`,
		},
	}

	col.TestSingleton = Map{
		"boil_main_test": {
			Standard: List{
				`"database/sql"`,
				`"flag"`,
				`"fmt"`,
				`"math/rand"`,
				`"os"`,
				`"path/filepath"`,
				`"strings"`,
				`"testing"`,
				`"time"`,
			},
			ThirdParty: List{
				`"github.com/kat-co/vala"`,
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/boil"`,
			},
		},
		"boil_queries_test": {
			Standard: List{
				`"bytes"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"math/rand"`,
				`"regexp"`,
			},
			ThirdParty: List{
				`"github.com/volatiletech/sqlboiler/boil"`,
			},
		},
		"boil_suites_test": {
			Standard: List{
				`"testing"`,
			},
		},
	}

	col.TestMain = Map{
		"psql": {
			Standard: List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			ThirdParty: List{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/lib/pq"`,
			},
		},
		"mysql": {
			Standard: List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"io/ioutil"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			ThirdParty: List{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/go-sql-driver/mysql"`,
			},
		},
		"mssql": {
			Standard: List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			ThirdParty: List{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/denisenkom/go-mssqldb"`,
			},
		},
		"crdb": {
			Standard: List{
				`"bytes"`,
				`"database/sql"`,
				`"fmt"`,
				`"io"`,
				`"os"`,
				`"os/exec"`,
				`"strings"`,
			},
			ThirdParty: List{
				`"github.com/pkg/errors"`,
				`"github.com/spf13/viper"`,
				`"github.com/volatiletech/sqlboiler/drivers/sqlboiler-crdb/driver"`,
				`"github.com/volatiletech/sqlboiler/randomize"`,
				`_ "github.com/lib/pq"`,
			},
		},
	}

	// basedOnType imports are only included in the template output if the
	// database requires one of the following special types. Check
	// TranslateColumnType to see the type assignments.
	col.BasedOnType = Map{
		"null.Float32": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Float64": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int8": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int16": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int32": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Int64": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint8": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint16": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint32": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Uint64": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.String": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Bool": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Time": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.JSON": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"null.Bytes": {
			ThirdParty: List{`"gopkg.in/volatiletech/null.v6"`},
		},
		"time.Time": {
			Standard: List{`"time"`},
		},
		"types.JSON": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.BytesArray": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.Int64Array": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.Float64Array": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.BoolArray": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.StringArray": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
		"types.Hstore": {
			ThirdParty: List{`"github.com/volatiletech/sqlboiler/types"`},
		},
	}

	return col
}

func combineImports(a, b Set) Set {
	var c Set

	c.Standard = strmangle.RemoveDuplicates(combineStringSlices(a.Standard, b.Standard))
	c.ThirdParty = strmangle.RemoveDuplicates(combineStringSlices(a.ThirdParty, b.ThirdParty))

	sort.Sort(c.Standard)
	sort.Sort(c.ThirdParty)

	return c
}

func CombineTypeImports(a Set, b map[string]Set, columnTypes []string) Set {
	tmpImp := Set{
		Standard:   make(List, len(a.Standard)),
		ThirdParty: make(List, len(a.ThirdParty)),
	}

	copy(tmpImp.Standard, a.Standard)
	copy(tmpImp.ThirdParty, a.ThirdParty)

	for _, typ := range columnTypes {
		for key, imp := range b {
			if typ == key {
				tmpImp.Standard = append(tmpImp.Standard, imp.Standard...)
				tmpImp.ThirdParty = append(tmpImp.ThirdParty, imp.ThirdParty...)
			}
		}
	}

	tmpImp.Standard = strmangle.RemoveDuplicates(tmpImp.Standard)
	tmpImp.ThirdParty = strmangle.RemoveDuplicates(tmpImp.ThirdParty)

	sort.Sort(tmpImp.Standard)
	sort.Sort(tmpImp.ThirdParty)

	return tmpImp
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
