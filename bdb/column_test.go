package bdb

import (
	"strings"
	"testing"
)

func TestColumnNames(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "one"},
		{Name: "two"},
		{Name: "three"},
	}

	out := strings.Join(ColumnNames(cols), " ")
	if out != "one two three" {
		t.Error("output was wrong:", out)
	}
}

func TestColumnDBTypes(t *testing.T) {
	cols := []Column{
		{Name: "test_one", DBType: "integer"},
		{Name: "test_two", DBType: "interval"},
	}

	res := ColumnDBTypes(cols)
	if res["TestOne"] != "integer" {
		t.Errorf(`Expected res["TestOne"]="integer", got: %s`, res["TestOne"])
	}
	if res["TestTwo"] != "interval" {
		t.Errorf(`Expected res["TestOne"]="interval", got: %s`, res["TestOne"])
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

func TestFilterColumnsByEnum(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "col1", DBType: "enum('hello')"},
		{Name: "col2", DBType: "enum('hello','there')"},
		{Name: "col3", DBType: "enum"},
		{Name: "col4", DBType: ""},
		{Name: "col5", DBType: "int"},
	}

	res := FilterColumnsByEnum(cols)
	if res[0].Name != `col1` {
		t.Errorf("Invalid result: %#v", res)
	}
	if res[1].Name != `col2` {
		t.Errorf("Invalid result: %#v", res)
	}
}
