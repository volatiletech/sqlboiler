package cmds

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

func TestCombineTypeImports(t *testing.T) {
	imports1 := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
		},
	}

	importsExpected := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
			`"gopkg.in/guregu/null.v3"`,
		},
	}

	cols := []dbdrivers.Column{
		dbdrivers.Column{
			Type: "null.Time",
		},
		dbdrivers.Column{
			Type: "null.Time",
		},
		dbdrivers.Column{
			Type: "time.Time",
		},
		dbdrivers.Column{
			Type: "null.Float",
		},
	}

	res1 := combineTypeImports(imports1, sqlBoilerTypeImports, cols)

	if !reflect.DeepEqual(res1, importsExpected) {
		t.Errorf("Expected res1 to match importsExpected, got:\n\n%#v\n", res1)
	}

	imports2 := imports{
		standard: importList{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		thirdparty: importList{
			`"github.com/pobri19/sqlboiler/boil"`,
			`"gopkg.in/guregu/null.v3"`,
		},
	}

	res2 := combineTypeImports(imports2, sqlBoilerTypeImports, cols)

	if !reflect.DeepEqual(res2, importsExpected) {
		t.Errorf("Expected res2 to match importsExpected, got:\n\n%#v\n", res1)
	}
}

func TestSortImports(t *testing.T) {
	t.Parallel()

	a1 := importList{
		`"fmt"`,
		`"errors"`,
	}
	a2 := importList{
		`_ "github.com/lib/pq"`,
		`_ "github.com/gorilla/n"`,
		`"github.com/gorilla/mux"`,
		`"github.com/gorilla/websocket"`,
	}

	a1Expected := importList{`"errors"`, `"fmt"`}
	a2Expected := importList{
		`"github.com/gorilla/mux"`,
		`_ "github.com/gorilla/n"`,
		`"github.com/gorilla/websocket"`,
		`_ "github.com/lib/pq"`,
	}

	sort.Sort(a1)
	if !reflect.DeepEqual(a1, a1Expected) {
		t.Errorf("Expected a1 to match a1Expected, got: %v", a1)
	}

	for i, v := range a1 {
		if v != a1Expected[i] {
			t.Errorf("Expected a1[%d] to match a1Expected[%d]:\n%s\n%s\n", i, i, v, a1Expected[i])
		}
	}

	sort.Sort(a2)
	if !reflect.DeepEqual(a2, a2Expected) {
		t.Errorf("Expected a2 to match a2expected, got: %v", a2)
	}

	for i, v := range a2 {
		if v != a2Expected[i] {
			t.Errorf("Expected a2[%d] to match a2Expected[%d]:\n%s\n%s\n", i, i, v, a1Expected[i])
		}
	}
}

func TestCombineImports(t *testing.T) {
	t.Parallel()

	a := imports{
		standard:   importList{"fmt"},
		thirdparty: importList{"github.com/pobri19/sqlboiler", "gopkg.in/guregu/null.v3"},
	}
	b := imports{
		standard:   importList{"os"},
		thirdparty: importList{"github.com/pobri19/sqlboiler"},
	}

	c := combineImports(a, b)

	if c.standard[0] != "fmt" && c.standard[1] != "os" {
		t.Errorf("Wanted: fmt, os got: %#v", c.standard)
	}
	if c.thirdparty[0] != "github.com/pobri19/sqlboiler" && c.thirdparty[1] != "gopkg.in/guregu/null.v3" {
		t.Errorf("Wanted: github.com/pobri19/sqlboiler, gopkg.in/guregu/null.v3 got: %#v", c.thirdparty)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	t.Parallel()

	hasDups := func(possible []string) error {
		for i := 0; i < len(possible)-1; i++ {
			for j := i + 1; j < len(possible); j++ {
				if possible[i] == possible[j] {
					return fmt.Errorf("found duplicate: %s [%d] [%d]", possible[i], i, j)
				}
			}
		}

		return nil
	}

	if len(removeDuplicates([]string{})) != 0 {
		t.Error("It should have returned an empty slice")
	}

	oneItem := []string{"patrick"}
	slice := removeDuplicates(oneItem)
	if ln := len(slice); ln != 1 {
		t.Error("Length was wrong:", ln)
	} else if oneItem[0] != slice[0] {
		t.Errorf("Slices differ: %#v %#v", oneItem, slice)
	}

	slice = removeDuplicates([]string{"hello", "patrick", "hello"})
	if ln := len(slice); ln != 2 {
		t.Error("Length was wrong:", ln)
	}
	if err := hasDups(slice); err != nil {
		t.Error(err)
	}

	slice = removeDuplicates([]string{"five", "patrick", "hello", "hello", "patrick", "hello", "hello"})
	if ln := len(slice); ln != 3 {
		t.Error("Length was wrong:", ln)
	}
	if err := hasDups(slice); err != nil {
		t.Error(err)
	}
}

func TestCombineStringSlices(t *testing.T) {
	t.Parallel()

	var a, b []string
	slice := combineStringSlices(a, b)
	if ln := len(slice); ln != 0 {
		t.Error("Len was wrong:", ln)
	}

	a = []string{"1", "2"}
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 2 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != a[0] || slice[1] != a[1] {
		t.Errorf("Slice mismatch: %#v %#v", a, slice)
	}

	b = a
	a = nil
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 2 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != b[0] || slice[1] != b[1] {
		t.Errorf("Slice mismatch: %#v %#v", b, slice)
	}

	a = b
	b = []string{"3", "4"}
	slice = combineStringSlices(a, b)
	if ln := len(slice); ln != 4 {
		t.Error("Len was wrong:", ln)
	} else if slice[0] != a[0] || slice[1] != a[1] || slice[2] != b[0] || slice[3] != b[1] {
		t.Errorf("Slice mismatch: %#v + %#v != #%v", a, b, slice)
	}
}
