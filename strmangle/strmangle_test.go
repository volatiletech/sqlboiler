package strmangle

import (
	"strings"
	"testing"
)

func TestIdentQuote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{In: `thing`, Out: `"thing"`},
		{In: `null`, Out: `null`},
		{In: `"thing"`, Out: `"thing"`},
		{In: `*`, Out: `*`},
		{In: ``, Out: ``},
		{In: `thing.thing`, Out: `"thing"."thing"`},
		{In: `"thing"."thing"`, Out: `"thing"."thing"`},
		{In: `thing.thing.thing.thing`, Out: `"thing"."thing"."thing"."thing"`},
		{In: `thing."thing".thing."thing"`, Out: `"thing"."thing"."thing"."thing"`},
		{In: `count(*) as ab, thing as bd`, Out: `count(*) as ab, thing as bd`},
		{In: `hello.*`, Out: `"hello".*`},
		{In: `hello.there.*`, Out: `"hello"."there".*`},
		{In: `"hello".there.*`, Out: `"hello"."there".*`},
		{In: `hello."there".*`, Out: `"hello"."there".*`},
	}

	for _, test := range tests {
		if got := IdentQuote(test.In); got != test.Out {
			t.Errorf("want: %s, got: %s", test.Out, got)
		}
	}
}

func TestIdentQuoteSlice(t *testing.T) {
	t.Parallel()

	ret := IdentQuoteSlice([]string{`thing`, `null`})
	if ret[0] != `"thing"` {
		t.Error(ret[0])
	}
	if ret[1] != `null` {
		t.Error(ret[1])
	}
}

func TestIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  int
		Out string
	}{
		{In: 0, Out: "a"},
		{In: 25, Out: "z"},
		{In: 26, Out: "ba"},
		{In: 52, Out: "ca"},
		{In: 675, Out: "zz"},
		{In: 676, Out: "baa"},
	}

	for _, test := range tests {
		if got := Identifier(test.In); got != test.Out {
			t.Errorf("[%d] want: %q, got: %q", test.In, test.Out, got)
		}
	}
}

func TestDriverUsesLastInsertID(t *testing.T) {
	t.Parallel()

	if DriverUsesLastInsertID("postgres") {
		t.Error("postgres does not support LastInsertId")
	}
	if !DriverUsesLastInsertID("mysql") {
		t.Error("postgres does support LastInsertId")
	}
}

func TestPlaceholders(t *testing.T) {
	t.Parallel()

	x := Placeholders(1, 2, 1)
	want := "$2"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(5, 1, 1)
	want = "$1,$2,$3,$4,$5"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(6, 1, 2)
	want = "($1,$2),($3,$4),($5,$6)"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(9, 1, 3)
	want = "($1,$2,$3),($4,$5,$6),($7,$8,$9)"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(7, 1, 3)
	want = "($1,$2,$3),($4,$5,$6),($7)"
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
		{"hello_there_people", "hello_there_person"},
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
		{"hello_there_person", "hello_there_people"},
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
		{"uid", "UID"},
		{"guid", "GUID"},
		{"uid", "UID"},
		{"uuid", "UUID"},
		{"ssn", "SSN"},
		{"tz", "TZ"},
		{"thing_guid", "ThingGUID"},
		{"guid_thing", "GUIDThing"},
		{"thing_guid_thing", "ThingGUIDThing"},
		{"id", "ID"},
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
		{"uid", "uid"},
		{"guid", "guid"},
		{"uid", "uid"},
		{"uuid", "uuid"},
		{"ssn", "ssn"},
		{"tz", "tz"},
		{"thing_guid", "thingGUID"},
		{"guid_thing", "guidThing"},
		{"thing_guid_thing", "thingGUIDThing"},
	}

	for i, test := range tests {
		if out := CamelCase(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestMakeStringMap(t *testing.T) {
	t.Parallel()

	var m map[string]string
	r := MakeStringMap(m)

	if r != "" {
		t.Errorf("Expected empty result, got: %s", r)
	}

	m = map[string]string{
		"TestOne": "interval",
		"TestTwo": "integer",
	}

	r = MakeStringMap(m)

	e1 := `"TestOne": "interval", "TestTwo": "integer"`
	e2 := `"TestTwo": "integer", "TestOne": "interval"`

	if r != e1 && r != e2 {
		t.Errorf("Got %s", r)
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
		r := WhereClause(test.Start, test.Cols)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
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

func TestJoinSlices(t *testing.T) {
	t.Parallel()

	ret := JoinSlices("", nil, nil)
	if ret != nil {
		t.Error("want nil, got:", ret)
	}

	ret = JoinSlices(" ", []string{"one", "two"}, []string{"three", "four"})
	if got := ret[0]; got != "one three" {
		t.Error("ret element was wrong:", got)
	}
	if got := ret[1]; got != "two four" {
		t.Error("ret element was wrong:", got)
	}
}

func TestJoinSlicesFail(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Error("did not panic")
		}
	}()

	JoinSlices("", nil, []string{"hello"})
}
