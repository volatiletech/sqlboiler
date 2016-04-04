package cmds

import (
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

var testColumns = []dbdrivers.Column{
	{Name: "friend_column", Type: "int", IsNullable: false},
	{Name: "enemy_column_thing", Type: "string", IsNullable: true},
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
		if out := singular(test.In); out != test.Out {
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
		if out := plural(test.In); out != test.Out {
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
		if out := titleCase(test.In); out != test.Out {
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
		if out := camelCase(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}

func TestMakeDBName(t *testing.T) {
	t.Parallel()

	if out := makeDBName("a", "b"); out != "a_b" {
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

	out := updateParamNames(testCols, []string{"id"})
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

	out := updateParamVariables("o.", testCols, []string{"id"})
	if out != "o.FriendColumn, o.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamNames(t *testing.T) {
	t.Parallel()

	out := insertParamNames(testColumns)
	if out != "friend_column, enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamFlags(t *testing.T) {
	t.Parallel()

	out := insertParamFlags(testColumns)
	if out != "$1, $2" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamVariables(t *testing.T) {
	out := insertParamVariables("o.", testColumns)
	if out != "o.FriendColumn, o.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestSelectParamFlags(t *testing.T) {
	t.Parallel()

	out := selectParamNames("table", testColumns)
	if out != "friend_column AS table_friend_column, enemy_column_thing AS table_enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}

func TestScanParams(t *testing.T) {
	t.Parallel()

	out := scanParamNames("object", testColumns)
	if out != "&object.FriendColumn, &object.EnemyColumnThing" {
		t.Error("Wrong output:", out)
	}
}

func TestHasPrimaryKey(t *testing.T) {
	t.Parallel()

	var pkey *dbdrivers.PrimaryKey
	if hasPrimaryKey(pkey) {
		t.Errorf("1) Expected false, got true")
	}

	pkey = &dbdrivers.PrimaryKey{}
	if hasPrimaryKey(pkey) {
		t.Errorf("2) Expected false, got true")
	}

	pkey.Columns = append(pkey.Columns, "test")
	if !hasPrimaryKey(pkey) {
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
		r := paramsPrimaryKey(test.Prefix, test.Pkey.Columns, true)
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
		r := paramsPrimaryKey(test.Prefix, test.Pkey.Columns, false)
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
		r := wherePrimaryKey(&test.Pkey, test.Start)
		if r != test.Should {
			t.Errorf("(%d) want: %s, got: %s\nTest: %#v", i, test.Should, r, test)
		}
	}
}
