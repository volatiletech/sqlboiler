package cmds

import (
	"reflect"
	"sort"
	"testing"
	"text/template"
)

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

func TestSortTemplates(t *testing.T) {
	templs := templater{
		template.New("bob.tpl"),
		template.New("all.tpl"),
		template.New("struct.tpl"),
		template.New("ttt.tpl"),
	}

	expected := []string{"bob.tpl", "all.tpl", "struct.tpl", "ttt.tpl"}

	for i, v := range templs {
		if v.Name() != expected[i] {
			t.Errorf("Order mismatch, expected: %s, got: %s", expected[i], v.Name())
		}
	}

	expected = []string{"struct.tpl", "all.tpl", "bob.tpl", "ttt.tpl"}

	sort.Sort(templs)

	for i, v := range templs {
		if v.Name() != expected[i] {
			t.Errorf("Order mismatch, expected: %s, got: %s", expected[i], v.Name())
		}
	}
}
