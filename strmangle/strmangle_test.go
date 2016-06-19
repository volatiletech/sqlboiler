package strmangle

import (
	"testing"

	"github.com/nullbio/sqlboiler/dbdrivers"
)

var testColumns = []dbdrivers.Column{
	{Name: "friend_column", Type: "int", IsNullable: false},
	{Name: "enemy_column_thing", Type: "string", IsNullable: true},
}

func TestCommaList(t *testing.T) {
	t.Parallel()

	cols := []string{
		"test1",
	}

	x := CommaList(cols)

	if x != `"test1"` {
		t.Errorf(`Expected "test1" - got %s`, x)
	}

	cols = append(cols, "test2")

	x = CommaList(cols)

	if x != `"test1", "test2"` {
		t.Errorf(`Expected "test1", "test2" - got %s`, x)
	}

	cols = append(cols, "test3")

	x = CommaList(cols)

	if x != `"test1", "test2", "test3"` {
		t.Errorf(`Expected "test1", "test2", "test3" - got %s`, x)
	}
}

func TestTitleCaseCommaList(t *testing.T) {
	t.Parallel()

	cols := []string{
		"test_id",
		"test_thing",
		"test_stuff_thing",
		"test",
	}

	x := TitleCaseCommaList("", cols)
	expected := `TestID, TestThing, TestStuffThing, Test`
	if x != expected {
		t.Errorf("Expected %s, got %s", expected, x)
	}

	x = TitleCaseCommaList("o.", cols)
	expected = `o.TestID, o.TestThing, o.TestStuffThing, o.Test`
	if x != expected {
		t.Errorf("Expected %s, got %s", expected, x)
	}
}

func TestCamelCaseCommaList(t *testing.T) {
	t.Parallel()

	cols := []string{
		"test_id",
		"test_thing",
		"test_stuff_thing",
		"test",
	}

	x := CamelCaseCommaList("", cols)
	expected := `testID, testThing, testStuffThing, test`
	if x != expected {
		t.Errorf("Expected %s, got %s", expected, x)
	}

	x = CamelCaseCommaList("o.", cols)
	expected = `o.testID, o.testThing, o.testStuffThing, o.test`
	if x != expected {
		t.Errorf("Expected %s, got %s", expected, x)
	}
}

func TestAutoIncPrimaryKey(t *testing.T) {
	t.Parallel()

	var pkey *dbdrivers.PrimaryKey
	var cols []dbdrivers.Column

	r := AutoIncPrimaryKey(cols, pkey)
	if r != "" {
		t.Errorf("Expected empty string, got %s", r)
	}

	pkey = &dbdrivers.PrimaryKey{
		Columns: []string{
			"col1", "auto",
		},
		Name: "",
	}

	cols = []dbdrivers.Column{
		{
			Name:       "thing",
			IsNullable: true,
			Type:       "int64",
			Default:    "nextval('abc'::regclass)",
		},
		{
			Name:       "stuff",
			IsNullable: false,
			Type:       "string",
			Default:    "nextval('abc'::regclass)",
		},
		{
			Name:       "other",
			IsNullable: false,
			Type:       "int64",
			Default:    "nextval",
		},
	}

	r = AutoIncPrimaryKey(cols, pkey)
	if r != "" {
		t.Errorf("Expected empty string, got %s", r)
	}

	cols = append(cols, dbdrivers.Column{
		Name:       "auto",
		IsNullable: false,
		Type:       "int64",
		Default:    "nextval('abc'::regclass)",
	})

	r = AutoIncPrimaryKey(cols, pkey)
	if r != "auto" {
		t.Errorf("Expected empty string, got %s", r)
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

func TestMakeDBName(t *testing.T) {
	t.Parallel()

	if out := MakeDBName("a", "b"); out != "a_b" {
		t.Error("Out was wrong:", out)
	}
}

func TestUpdateParamNames(t *testing.T) {
	t.Parallel()

	var testCols = []dbdrivers.Column{
		{Name: "id", Type: "int", IsNullable: false},
		{Name: "friend_column", Type: "int", IsNullable: false},
		{Name: "enemy_column_thing", Type: "string", IsNullable: true},
	}

	out := UpdateParamNames(testCols, []string{"id"})
	if out != "friend_column=$1,enemy_column_thing=$2" {
		t.Error("Wrong output:", out)
	}
}

func TestUpdateParamVariables(t *testing.T) {
	t.Parallel()

	var testCols = []dbdrivers.Column{
		{Name: "id", Type: "int", IsNullable: false},
		{Name: "friend_column", Type: "int", IsNullable: false},
		{Name: "enemy_column_thing", Type: "string", IsNullable: true},
	}

	out := UpdateParamVariables("o.", testCols, []string{"id"})
	if out != "o.FriendColumn, o.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamNames(t *testing.T) {
	t.Parallel()

	out := InsertParamNames(testColumns)
	if out != "friend_column, enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamFlags(t *testing.T) {
	t.Parallel()

	out := InsertParamFlags(testColumns)
	if out != "$1, $2" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamVariables(t *testing.T) {
	t.Parallel()

	out := InsertParamVariables("o.", testColumns)
	if out != "o.FriendColumn, o.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestSelectParamFlags(t *testing.T) {
	t.Parallel()

	out := SelectParamNames("table", testColumns)
	if out != "friend_column AS table_friend_column, enemy_column_thing AS table_enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}

func TestScanParams(t *testing.T) {
	t.Parallel()

	out := ScanParamNames("object", testColumns)
	if out != "&object.FriendColumn, &object.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestHasPrimaryKey(t *testing.T) {
	t.Parallel()

	var pkey *dbdrivers.PrimaryKey
	if HasPrimaryKey(pkey) {
		t.Errorf("1) Expected false, got true")
	}

	pkey = &dbdrivers.PrimaryKey{}
	if HasPrimaryKey(pkey) {
		t.Errorf("2) Expected false, got true")
	}

	pkey.Columns = append(pkey.Columns, "test")
	if !HasPrimaryKey(pkey) {
		t.Errorf("3) Expected true, got false")
	}
}

func TestParamsPrimaryKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Pkey   dbdrivers.PrimaryKey
		Prefix string
		Should string
	}{
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one"}},
			Prefix: "o.", Should: "o.ColOne",
		},
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one", "col_two"}},
			Prefix: "o.", Should: "o.ColOne, o.ColTwo",
		},
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one", "col_two", "col_three"}},
			Prefix: "o.", Should: "o.ColOne, o.ColTwo, o.ColThree",
		},
	}

	for i, test := range tests {
		r := ParamsPrimaryKey(test.Prefix, test.Pkey.Columns, true)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}

	tests2 := []struct {
		Pkey   dbdrivers.PrimaryKey
		Prefix string
		Should string
	}{
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one"}},
			Prefix: "o.", Should: "o.col_one",
		},
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one", "col_two"}},
			Prefix: "o.", Should: "o.col_one, o.col_two",
		},
		{
			Pkey:   dbdrivers.PrimaryKey{Columns: []string{"col_one", "col_two", "col_three"}},
			Prefix: "o.", Should: "o.col_one, o.col_two, o.col_three",
		},
	}

	for i, test := range tests2 {
		r := ParamsPrimaryKey(test.Prefix, test.Pkey.Columns, false)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
}

func TestWherePrimaryKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Pkey   dbdrivers.PrimaryKey
		Start  int
		Should string
	}{
		{Pkey: dbdrivers.PrimaryKey{Columns: []string{"col1"}}, Start: 2, Should: "col1=$2"},
		{Pkey: dbdrivers.PrimaryKey{Columns: []string{"col1", "col2"}}, Start: 4, Should: "col1=$4 AND col2=$5"},
		{Pkey: dbdrivers.PrimaryKey{Columns: []string{"col1", "col2", "col3"}}, Start: 4, Should: "col1=$4 AND col2=$5 AND col3=$6"},
	}

	for i, test := range tests {
		r := WherePrimaryKey(test.Pkey.Columns, test.Start)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
}

func TestFilterColumnsByDefault(t *testing.T) {
	t.Parallel()

	cols := []dbdrivers.Column{
		{
			Name:    "col1",
			Default: "",
		},
		{
			Name:    "col2",
			Default: "things",
		},
		{
			Name:    "col3",
			Default: "",
		},
		{
			Name:    "col4",
			Default: "things2",
		},
	}

	res := FilterColumnsByDefault(cols, false)
	if res != `"col1","col3"` {
		t.Errorf("Invalid result: %s", res)
	}

	res = FilterColumnsByDefault(cols, true)
	if res != `"col2","col4"` {
		t.Errorf("Invalid result: %s", res)
	}

	res = FilterColumnsByDefault([]dbdrivers.Column{}, false)
	if res != `` {
		t.Errorf("Invalid result: %s", res)
	}
}

func TestFilterColumnsByAutoIncrement(t *testing.T) {
	t.Parallel()

	cols := []dbdrivers.Column{
		{
			Name:    "col1",
			Default: `nextval("thing"::thing)`,
		},
		{
			Name:    "col2",
			Default: "things",
		},
		{
			Name:    "col3",
			Default: "",
		},
		{
			Name:    "col4",
			Default: `nextval("thing"::thing)`,
		},
	}

	res := FilterColumnsByAutoIncrement(cols)
	if res != `"col1","col4"` {
		t.Errorf("Invalid result: %s", res)
	}

	res = FilterColumnsByAutoIncrement([]dbdrivers.Column{})
	if res != `` {
		t.Errorf("Invalid result: %s", res)
	}
}
