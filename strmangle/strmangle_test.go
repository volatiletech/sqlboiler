package strmangle

import (
	"strings"
	"testing"
)

func TestDriverUsesLastInsertID(t *testing.T) {
	t.Parallel()

	if DriverUsesLastInsertID("postgres") {
		t.Error("postgres does not support LastInsertId")
	}
	if !DriverUsesLastInsertID("mysql") {
		t.Error("postgres does support LastInsertId")
	}
}

func TestGenerateParamFlags(t *testing.T) {
	t.Parallel()

	x := GenerateParamFlags(5, 1)
	want := "$1,$2,$3,$4,$5"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}
}

func TestSingular(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"hello_people", "hello_person"},
		{"hello_person", "hello_person"},
		{"friends", "friend"},
	}

	for i, test := range tests {
		if out := Singular(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestPlural(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"hello_person", "hello_people"},
		{"friend", "friends"},
		{"friends", "friends"},
	}

	for i, test := range tests {
		if out := Plural(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestTitleCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"hello_there", "HelloThere"},
		{"", ""},
		{"fun_id", "FunID"},
	}

	for i, test := range tests {
		if out := TitleCase(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestCamelCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"hello_there_sunny", "helloThereSunny"},
		{"", ""},
		{"fun_id_times", "funIDTimes"},
	}

	for i, test := range tests {
		if out := CamelCase(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestStringMap(t *testing.T) {
	t.Parallel()

	mapped := StringMap(strings.ToLower, []string{"HELLO", "WORLD"})
	if got := strings.Join(mapped, " "); got != "hello world" {
		t.Errorf("mapped was wrong: %q", got)
	}
}

func TestMakeDBName(t *testing.T) {
	t.Parallel()

	if out := MakeDBName("a", "b"); out != "a_b" {
		t.Error("Out was wrong:", out)
	}
}

func TestHasElement(t *testing.T) {
	t.Parallel()

	elements := []string{"one", "two"}
	if got := HasElement("one", elements); !got {
		t.Error("should have found element key")
	}
	if got := HasElement("three", elements); got {
		t.Error("should not have found element key")
	}
}

func TestPrefixStringSlice(t *testing.T) {
	t.Parallel()

	slice := PrefixStringSlice("o.", []string{"one", "two"})
	if got := strings.Join(slice, " "); got != "o.one o.two" {
		t.Error("wrong output:", got)
	}
}

func TestWhereClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Cols   []string
		Start  int
		Should string
	}{
		{Cols: []string{"col1"}, Start: 2, Should: `"col1"=$2`},
		{Cols: []string{"col1", "col2"}, Start: 4, Should: `"col1"=$4 AND "col2"=$5`},
		{Cols: []string{"col1", "col2", "col3"}, Start: 4, Should: `"col1"=$4 AND "col2"=$5 AND "col3"=$6`},
	}

	for i, test := range tests {
		r := WhereClause(test.Cols, test.Start)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
}

func TestWherePrimaryKeyPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Error("did not panic")
		}
	}()

	WhereClause(nil, 0)
}

func TestSubstring(t *testing.T) {
	t.Parallel()

	str := "hello"

	if got := Substring(0, 5, str); got != "hello" {
		t.Errorf("substring was wrong: %q", got)
	}
	if got := Substring(1, 4, str); got != "ell" {
		t.Errorf("substring was wrong: %q", got)
	}
	if got := Substring(2, 3, str); got != "l" {
		t.Errorf("substring was wrong: %q", got)
	}
	if got := Substring(5, 5, str); got != "" {
		t.Errorf("substring was wrong: %q", got)
	}
}
