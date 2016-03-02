package cmds

import (
	"testing"

	"github.com/pobri19/sqlboiler/dbdrivers"
)

var testColumns = []dbdrivers.DBColumn{
	{Name: "friend_column", Type: "int", IsNullable: false},
	{Name: "enemy_column_thing", Type: "string", IsNullable: true},
}

func TestTitleCase(t *testing.T) {
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
	if out := makeDBName("a", "b"); out != "a_b" {
		t.Error("Out was wrong:", out)
	}
}

func TestInsertParamNames(t *testing.T) {
	out := insertParamNames(testColumns)
	if out != "friend_column, enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}

func TestInsertParamFlags(t *testing.T) {
	out := insertParamFlags(testColumns)
	if out != "$1, $2" {
		t.Error("Wrong output:", out)
	}
}

func TestSelectParamFlags(t *testing.T) {
	out := selectParamNames("table", testColumns)
	if out != "friend_column AS table_friend_column, enemy_column_thing AS table_enemy_column_thing" {
		t.Error("Wrong output:", out)
	}
}
