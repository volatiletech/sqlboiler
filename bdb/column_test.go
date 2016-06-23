package bdb

import (
	"strings"
	"testing"
)

func TestColumnNames(t *testing.T) {
	t.Parallel()

	cols := []Column{
		Column{Name: "one"},
		Column{Name: "two"},
		Column{Name: "three"},
	}

	out := strings.Join(ColumnNames(cols), " ")
	if out != "one two three" {
		t.Error("output was wrong:", out)
	}
}

func TestFilterColumnsByDefault(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "col1", Default: ""},
		{Name: "col2", Default: "things"},
		{Name: "col3", Default: ""},
		{Name: "col4", Default: "things2"},
	}

	res := FilterColumnsByDefault(cols, false)
	if res[0].Name != `col1` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col3` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByDefault(cols, true)
	if res[0].Name != `col2` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col4` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByDefault([]Column{}, false)
	if res != nil {
		t.Errorf("Invalid result: %#v", res)
	}
}

func TestFilterColumnsByAutoIncrement(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "col1", Default: `nextval("thing"::thing)`},
		{Name: "col2", Default: "things"},
		{Name: "col3", Default: ""},
		{Name: "col4", Default: `nextval("thing"::thing)`},
	}

	res := FilterColumnsByAutoIncrement(cols)
	if res[0].Name != `col1` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col4` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByAutoIncrement([]Column{})
	if res != nil {
		t.Errorf("Invalid result: %#v", res)
	}
}
