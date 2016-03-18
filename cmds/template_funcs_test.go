package cmds

import (
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

var testColumns = []dbdrivers.DBColumn{
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
