package importers

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestSetFromInterface(t *testing.T) {
	t.Parallel()

	setIntf := map[string]interface{}{
		"standard": []interface{}{
			"hello",
			"there",
		},
		"third_party": []interface{}{
			"there",
			"hello",
		},
	}

	set, err := SetFromInterface(setIntf)
	if err != nil {
		t.Error(err)
	}

	if set.Standard[0] != "hello" {
		t.Error("set was wrong:", set.Standard[0])
	}
	if set.Standard[1] != "there" {
		t.Error("set was wrong:", set.Standard[1])
	}
	if set.ThirdParty[0] != "there" {
		t.Error("set was wrong:", set.ThirdParty[0])
	}
	if set.ThirdParty[1] != "hello" {
		t.Error("set was wrong:", set.ThirdParty[1])
	}
}

func TestMapFromInterface(t *testing.T) {
	t.Parallel()

	mapIntf := map[string]interface{}{
		"test_main": map[string]interface{}{
			"standard": []interface{}{
				"hello",
				"there",
			},
			"third_party": []interface{}{
				"there",
				"hello",
			},
		},
	}

	mp, err := MapFromInterface(mapIntf)
	if err != nil {
		t.Error(err)
	}

	set, ok := mp["test_main"]
	if !ok {
		t.Error("could not find set 'test_main'")
	}

	if set.Standard[0] != "hello" {
		t.Error("set was wrong:", set.Standard[0])
	}
	if set.Standard[1] != "there" {
		t.Error("set was wrong:", set.Standard[1])
	}
	if set.ThirdParty[0] != "there" {
		t.Error("set was wrong:", set.ThirdParty[0])
	}
	if set.ThirdParty[1] != "hello" {
		t.Error("set was wrong:", set.ThirdParty[1])
	}
}

func TestMapFromInterfaceAltSyntax(t *testing.T) {
	t.Parallel()

	mapIntf := []interface{}{
		map[string]interface{}{
			"name": "test_main",
			"standard": []interface{}{
				"hello",
				"there",
			},
			"third_party": []interface{}{
				"there",
				"hello",
			},
		},
	}

	mp, err := MapFromInterface(mapIntf)
	if err != nil {
		t.Error(err)
	}

	set, ok := mp["test_main"]
	if !ok {
		t.Error("could not find set 'test_main'")
	}

	if set.Standard[0] != "hello" {
		t.Error("set was wrong:", set.Standard[0])
	}
	if set.Standard[1] != "there" {
		t.Error("set was wrong:", set.Standard[1])
	}
	if set.ThirdParty[0] != "there" {
		t.Error("set was wrong:", set.ThirdParty[0])
	}
	if set.ThirdParty[1] != "hello" {
		t.Error("set was wrong:", set.ThirdParty[1])
	}
}

func TestImportsSort(t *testing.T) {
	t.Parallel()

	a1 := List{
		`"fmt"`,
		`"errors"`,
	}
	a2 := List{
		`_ "github.com/lib/pq"`,
		`_ "github.com/gorilla/n"`,
		`"github.com/gorilla/mux"`,
		`"github.com/gorilla/websocket"`,
	}

	a1Expected := List{`"errors"`, `"fmt"`}
	a2Expected := List{
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

func TestAddTypeImports(t *testing.T) {
	t.Parallel()

	imports1 := Set{
		Standard: List{
			`"errors"`,
			`"fmt"`,
		},
		ThirdParty: List{
			`"github.com/volatiletech/sqlboiler/v4/boil"`,
		},
	}

	importsExpected := Set{
		Standard: List{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		ThirdParty: List{
			`"github.com/volatiletech/null/v8"`,
			`"github.com/volatiletech/sqlboiler/v4/boil"`,
		},
	}

	types := []string{
		"null.Time",
		"null.Time",
		"time.Time",
	}

	imps := NewDefaultImports()

	imps.BasedOnType = Map{
		"null.Time": Set{ThirdParty: List{`"github.com/volatiletech/null/v8"`}},
		"time.Time": Set{Standard: List{`"time"`}},
	}

	res1 := AddTypeImports(imports1, imps.BasedOnType, types)

	if !reflect.DeepEqual(res1, importsExpected) {
		t.Errorf("Expected res1 to match importsExpected, got:\n\n%#v\n", res1)
	}

	imports2 := Set{
		Standard: List{
			`"errors"`,
			`"fmt"`,
			`"time"`,
		},
		ThirdParty: List{
			`"github.com/volatiletech/null/v8"`,
			`"github.com/volatiletech/sqlboiler/v4/boil"`,
		},
	}

	res2 := AddTypeImports(imports2, imps.BasedOnType, types)

	if !reflect.DeepEqual(res2, importsExpected) {
		t.Errorf("Expected res2 to match importsExpected, got:\n\n%#v\n", res1)
	}
}

func TestMergeSet(t *testing.T) {
	t.Parallel()

	a := Set{
		Standard:   List{"fmt"},
		ThirdParty: List{"github.com/volatiletech/sqlboiler/v4", "github.com/volatiletech/null/v8"},
	}
	b := Set{
		Standard:   List{"os"},
		ThirdParty: List{"github.com/volatiletech/sqlboiler/v4"},
	}

	c := mergeSet(a, b)

	if c.Standard[0] != "fmt" && c.Standard[1] != "os" {
		t.Errorf("Wanted: fmt, os got: %#v", c.Standard)
	}
	if c.ThirdParty[0] != "github.com/volatiletech/null/v8" && c.ThirdParty[1] != "github.com/volatiletech/sqlboiler/v4" {
		t.Errorf("Wanted: github.com/volatiletech/sqlboiler, github.com/volatiletech/null/v8 got: %#v", c.ThirdParty)
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

func TestMerge(t *testing.T) {
	var a, b Collection

	a.All = Set{Standard: List{"aa"}, ThirdParty: List{"aa"}}
	a.Test = Set{Standard: List{"at"}, ThirdParty: List{"at"}}
	a.Singleton = Map{
		"a": {Standard: List{"as"}, ThirdParty: List{"as"}},
		"c": {Standard: List{"as"}, ThirdParty: List{"as"}},
	}
	a.TestSingleton = Map{
		"a": {Standard: List{"at"}, ThirdParty: List{"at"}},
		"c": {Standard: List{"at"}, ThirdParty: List{"at"}},
	}
	a.BasedOnType = Map{
		"a": {Standard: List{"abot"}, ThirdParty: List{"abot"}},
		"c": {Standard: List{"abot"}, ThirdParty: List{"abot"}},
	}

	b.All = Set{Standard: List{"bb"}, ThirdParty: List{"bb"}}
	b.Test = Set{Standard: List{"bt"}, ThirdParty: List{"bt"}}
	b.Singleton = Map{
		"b": {Standard: List{"bs"}, ThirdParty: List{"bs"}},
		"c": {Standard: List{"bs"}, ThirdParty: List{"bs"}},
	}
	b.TestSingleton = Map{
		"b": {Standard: List{"bt"}, ThirdParty: List{"bt"}},
		"c": {Standard: List{"bt"}, ThirdParty: List{"bt"}},
	}
	b.BasedOnType = Map{
		"b": {Standard: List{"bbot"}, ThirdParty: List{"bbot"}},
		"c": {Standard: List{"bbot"}, ThirdParty: List{"bbot"}},
	}

	c := Merge(a, b)

	setHas := func(s Set, a, b string) {
		t.Helper()
		if s.Standard[0] != a {
			t.Error("standard index 0, want:", a, "got:", s.Standard[0])
		}
		if s.Standard[1] != b {
			t.Error("standard index 1, want:", a, "got:", s.Standard[1])
		}
		if s.ThirdParty[0] != a {
			t.Error("third party index 0, want:", a, "got:", s.ThirdParty[0])
		}
		if s.ThirdParty[1] != b {
			t.Error("third party index 1, want:", a, "got:", s.ThirdParty[1])
		}
	}
	mapHas := func(m Map, key, a, b string) {
		t.Helper()
		setHas(m[key], a, b)
	}

	setHas(c.All, "aa", "bb")
	setHas(c.Test, "at", "bt")
	mapHas(c.Singleton, "c", "as", "bs")
	mapHas(c.TestSingleton, "c", "at", "bt")
	mapHas(c.BasedOnType, "c", "abot", "bbot")

	if t.Failed() {
		t.Logf("%#v\n", c)
	}
}

var testImportStringExpect = `import (
	"fmt"

	"github.com/friendsofgo/errors"
)`

func TestSetFormat(t *testing.T) {
	t.Parallel()

	s := Set{
		Standard: List{
			`"fmt"`,
		},
		ThirdParty: List{
			`"github.com/friendsofgo/errors"`,
		},
	}

	got := strings.TrimSpace(string(s.Format()))
	if got != testImportStringExpect {
		t.Error("want:\n", testImportStringExpect, "\ngot:\n", got)
	}
}
