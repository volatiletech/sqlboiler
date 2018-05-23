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
		if got := IdentQuote('"', '"', test.In); got != test.Out {
			t.Errorf("want: %s, got: %s", test.Out, got)
		}
	}
}

func TestIdentQuoteSlice(t *testing.T) {
	t.Parallel()

	ret := IdentQuoteSlice('"', '"', []string{`thing`, `null`})
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

func TestQuoteCharacter(t *testing.T) {
	t.Parallel()

	if QuoteCharacter('[') != "[" {
		t.Error("want just the normal quote character")
	}
	if QuoteCharacter('`') != "`" {
		t.Error("want just the normal quote character")
	}
	if QuoteCharacter('"') != `\"` {
		t.Error("want an escaped character")
	}
}

func TestPlaceholders(t *testing.T) {
	t.Parallel()

	x := Placeholders(true, 1, 2, 1)
	want := "$2"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(true, 5, 1, 1)
	want = "$1,$2,$3,$4,$5"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(false, 5, 1, 1)
	want = "?,?,?,?,?"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(true, 6, 1, 2)
	want = "($1,$2),($3,$4),($5,$6)"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(true, 6, 1, 2)
	want = "($1,$2),($3,$4),($5,$6)"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(false, 9, 1, 3)
	want = "(?,?,?),(?,?,?),(?,?,?)"
	if want != x {
		t.Errorf("want %s, got %s", want, x)
	}

	x = Placeholders(true, 7, 1, 3)
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
		{"areas", "area"},
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
		{"area", "areas"},
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
		{"____a____a___", "AA"},
		{"_a_a_", "AA"},
		{"fun_id", "FunID"},
		{"_fun_id", "FunID"},
		{"__fun____id_", "FunID"},
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
		{"gvzxc", "GVZXC"},
		{"id_trgb_id", "IDTRGBID"},
		{"vzxx_vxccb_nmx", "VZXXVXCCBNMX"},
		{"thing_zxc_stuff_vxz", "ThingZXCStuffVXZ"},
		{"zxc_thing_vxz_stuff", "ZXCThingVXZStuff"},
		{"zxc_vdf9c9_hello9", "ZXCVDF9C9Hello9"},
		{"id9_uid911_guid9e9", "ID9UID911GUID9E9"},
		{"zxc_vdf0c0_hello0", "ZXCVDF0C0Hello0"},
		{"id0_uid000_guid0e0", "ID0UID000GUID0E0"},
		{"ab_5zxc5d5", "Ab5ZXC5D5"},
		{"Identifier", "Identifier"},
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
		{"a_", "a"},
		{"aaa_", "aaa"},
		{"a___", "a"},
		{"aaa___", "aaa"},
		{"_a_", "a"},
		{"_aaa_", "aaa"},
		{"____a___", "a"},
		{"___aaa____", "aaa"},
		{"___a", "a"},
		{"___aa", "aa"},
		{"____a____a___", "aA"},
		{"_a_a_", "aA"},
		{"_fun_id", "funID"},
		{"__fun____id_", "funID"},
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

func TestTitleCaseIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"hello", "Hello"},
		{"hello.world", "Hello.World"},
		{"hey.id.world", "Hey.ID.World"},
	}

	for i, test := range tests {
		if out := TitleCaseIdentifier(test.In); out != test.Out {
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

	e1 := "`TestOne`: `interval`, `TestTwo`: `integer`"
	e2 := "`TestTwo`: `integer`, `TestOne`: `interval`"

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

func TestPrefixStringSlice(t *testing.T) {
	t.Parallel()

	slice := PrefixStringSlice("o.", []string{"one", "two"})
	if got := strings.Join(slice, " "); got != "o.one o.two" {
		t.Error("wrong output:", got)
	}
}

func TestSetParamNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Cols   []string
		Start  int
		Should string
	}{
		{Cols: []string{"col1", "col2"}, Start: 0, Should: `"col1"=?,"col2"=?`},
		{Cols: []string{"col1"}, Start: 2, Should: `"col1"=$2`},
		{Cols: []string{"col1", "col2"}, Start: 4, Should: `"col1"=$4,"col2"=$5`},
		{Cols: []string{"col1", "col2", "col3"}, Start: 4, Should: `"col1"=$4,"col2"=$5,"col3"=$6`},
	}

	for i, test := range tests {
		r := SetParamNames(`"`, `"`, test.Start, test.Cols)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
}

func TestWhereClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Cols   []string
		Start  int
		Should string
	}{
		{Cols: []string{"col1", "col2"}, Start: 0, Should: `"col1"=? AND "col2"=?`},
		{Cols: []string{"col1"}, Start: 2, Should: `"col1"=$2`},
		{Cols: []string{"col1", "col2"}, Start: 4, Should: `"col1"=$4 AND "col2"=$5`},
		{Cols: []string{"col1", "col2", "col3"}, Start: 4, Should: `"col1"=$4 AND "col2"=$5 AND "col3"=$6`},
	}

	for i, test := range tests {
		r := WhereClause(`"`, `"`, test.Start, test.Cols)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
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

func TestStringSliceMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		a      []string
		b      []string
		expect bool
	}{
		{
			a:      []string{},
			b:      []string{},
			expect: true,
		},
		{
			a:      []string{"a"},
			b:      []string{},
			expect: false,
		},
		{
			a:      []string{"a"},
			b:      []string{"a"},
			expect: true,
		},
		{
			a:      []string{},
			b:      []string{"b"},
			expect: false,
		},
		{
			a:      []string{"c", "d"},
			b:      []string{"b", "d"},
			expect: false,
		},
		{
			a:      []string{"b", "d"},
			b:      []string{"c", "d"},
			expect: false,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"c", "b", "a"},
			expect: true,
		},
		{
			a:      []string{"a", "b", "c"},
			b:      []string{"a", "b", "c"},
			expect: true,
		},
	}

	for i, test := range tests {
		if StringSliceMatch(test.a, test.b) != test.expect {
			t.Errorf("%d) Expected match to return %v, but got %v", i, test.expect, !test.expect)
		}
	}
}

func TestContainsAny(t *testing.T) {
	t.Parallel()

	a := []string{"hello", "friend"}
	if ContainsAny([]string{}, "x") {
		t.Errorf("Should not contain x")
	}

	if ContainsAny(a, "x") {
		t.Errorf("Should not contain x")
	}

	if !ContainsAny(a, "hello") {
		t.Errorf("Should contain hello")
	}

	if !ContainsAny(a, "friend") {
		t.Errorf("Should contain friend")
	}

	if !ContainsAny(a, "hello", "friend") {
		t.Errorf("Should contain hello and friend")
	}

	if ContainsAny(a) {
		t.Errorf("Should not return true")
	}
}

func TestGenerateTags(t *testing.T) {
	tags := GenerateTags([]string{}, "col_name")
	if tags != "" {
		t.Errorf("Expected empty string, got %s", tags)
	}

	tags = GenerateTags([]string{"xml"}, "col_name")
	exp := `xml:"col_name" `
	if tags != exp {
		t.Errorf("expected %s, got %s", exp, tags)
	}

	tags = GenerateTags([]string{"xml", "db"}, "col_name")
	exp = `xml:"col_name" db:"col_name" `
	if tags != exp {
		t.Errorf("expected %s, got %s", exp, tags)
	}
}

func TestGenerateIgnoreTags(t *testing.T) {
	tags := GenerateIgnoreTags([]string{})
	if tags != "" {
		t.Errorf("Expected empty string, got %s", tags)
	}

	tags = GenerateIgnoreTags([]string{"xml"})
	exp := `xml:"-" `
	if tags != exp {
		t.Errorf("expected %s, got %s", exp, tags)
	}

	tags = GenerateIgnoreTags([]string{"xml", "db"})
	exp = `xml:"-" db:"-" `
	if tags != exp {
		t.Errorf("expected %s, got %s", exp, tags)
	}
}

func TestParseEnum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Enum string
		Name string
		Vals []string
	}{
		{"enum('one')", "", []string{"one"}},
		{"enum('one','two')", "", []string{"one", "two"}},
		{"enum.working('one')", "working", []string{"one"}},
		{"enum.wor_king('one','two')", "wor_king", []string{"one", "two"}},
	}

	for i, test := range tests {
		name := ParseEnumName(test.Enum)
		vals := ParseEnumVals(test.Enum)
		if name != test.Name {
			t.Errorf("%d) name was wrong, want: %s got: %s (%s)", i, test.Name, name, test.Enum)
		}
		for j, v := range test.Vals {
			if v != vals[j] {
				t.Errorf("%d.%d) value was wrong, want: %s got: %s (%s)", i, j, v, vals[j], test.Enum)
			}
		}
	}
}

func TestIsEnumNormal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Vals []string
		Ok   bool
	}{
		{[]string{"o1ne", "two2"}, true},
		{[]string{"one", "t#wo2"}, false},
		{[]string{"1one", "two2"}, false},
	}

	for i, test := range tests {
		if got := IsEnumNormal(test.Vals); got != test.Ok {
			t.Errorf("%d) want: %t got: %t, %#v", i, test.Ok, got, test.Vals)
		}
	}
}

func TestShouldTitleCaseEnum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Val string
		Ok  bool
	}{
		{"hello_there0", true},
		{"hEllo", false},
		{"_hello", false},
		{"0hello", false},
	}

	for i, test := range tests {
		if got := ShouldTitleCaseEnum(test.Val); got != test.Ok {
			t.Errorf("%d) want: %t got: %t, %v", i, test.Ok, got, test.Val)
		}
	}
}

func TestReplaceReservedWords(t *testing.T) {
	tests := []struct {
		Word    string
		Replace bool
	}{
		{"break", true},
		{"id", false},
		{"type", true},
	}

	for i, test := range tests {
		got := ReplaceReservedWords(test.Word)
		if test.Replace && !strings.HasSuffix(got, "_") {
			t.Errorf("%d) want suffixed (%s), got: %s", i, test.Word, got)
		} else if !test.Replace && strings.HasSuffix(got, "_") {
			t.Errorf("%d) want normal (%s), got: %s", i, test.Word, got)
		}
	}
}
