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

	res := FilterColumnsByDefault(false, cols)
	if res[0].Name != `col1` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col3` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByDefault(true, cols)
	if res[0].Name != `col2` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col4` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByDefault(false, []Column{})
	if res != nil {
		t.Errorf("Invalid result: %#v", res)
	}
}

func TestDefaultValues(t *testing.T) {
	t.Parallel()

	c := Column{}

	c.Default = `\x12345678`
	c.Type = "[]byte"

	res := DefaultValues([]Column{c})
	if len(res) != 1 {
		t.Errorf("Expected res len 1, got %d", len(res))
	}
	if res[0] != `[]byte{0x12, 0x34, 0x56, 0x78}` {
		t.Errorf("Invalid result: %#v", res)
	}

	c.Default = `\x`

	res = DefaultValues([]Column{c})
	if res[0] != `[]byte{}` {
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

	res := FilterColumnsByAutoIncrement(true, cols)
	if res[0].Name != `col1` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col4` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByAutoIncrement(false, cols)
	if res[0].Name != `col2` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col3` {
		t.Errorf("Invalid result: %#v", res)
	}

	res = FilterColumnsByAutoIncrement(true, []Column{})
	if res != nil {
		t.Errorf("Invalid result: %#v", res)
	}
}
