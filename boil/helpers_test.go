package boil

import (
	"testing"
	"time"
)

type testObj struct {
	ID       int
	Name     string `db:"TestHello"`
	HeadSize int
}

func TestWherePrimaryKeyIn(t *testing.T) {
	t.Parallel()

	x := WherePrimaryKeyIn(1, "aa")
	expect := `("aa") IN ($1)`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(2, "aa")
	expect = `("aa") IN ($1,$2)`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(3, "aa")
	expect = `("aa") IN ($1,$2,$3)`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(1, "aa", "bb")
	expect = `("aa","bb") IN (($1,$2))`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(2, "aa", "bb")
	expect = `("aa","bb") IN (($1,$2),($3,$4))`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(3, "aa", "bb")
	expect = `("aa","bb") IN (($1,$2),($3,$4),($5,$6))`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(4, "aa", "bb")
	expect = `("aa","bb") IN (($1,$2),($3,$4),($5,$6),($7,$8))`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}

	x = WherePrimaryKeyIn(4, "aa", "bb", "cc")
	expect = `("aa","bb","cc") IN (($1,$2,$3),($4,$5,$6),($7,$8,$9),($10,$11,$12))`
	if x != expect {
		t.Errorf("Expected %s, got %s\n", expect, x)
	}
}

func TestGoVarToSQLName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In, Out string
	}{
		{"IDStruct", "id_struct"},
		{"WigglyBits", "wiggly_bits"},
		{"HoboIDFriend3333", "hobo_id_friend3333"},
		{"3333friend", "3333friend"},
		{"ID3ID", "id3_id"},
		{"Wei3rd", "wei3rd"},
		{"He3I3Test", "he3_i3_test"},
		{"He3ID3Test", "he3_id3_test"},
		{"HelloFriendID", "hello_friend_id"},
	}

	for i, test := range tests {
		if out := goVarToSQLName(test.In); out != test.Out {
			t.Errorf("%d) from: %q, want: %q, got: %q", i, test.In, test.Out, out)
		}
	}
}

func TestSelectNames(t *testing.T) {
	t.Parallel()

	o := testObj{
		Name:     "bob",
		ID:       5,
		HeadSize: 23,
	}

	result := SelectNames(o)
	if result != `id, TestHello, head_size` {
		t.Error("Result was wrong, got:", result)
	}
}

func TestWhereClause(t *testing.T) {
	t.Parallel()

	columns := map[string]interface{}{
		"name": "bob",
		"id":   5,
		"date": time.Now(),
	}

	result := WhereClause(columns)

	if result != `date=$1 AND id=$2 AND name=$3` {
		t.Error("Result was wrong, got:", result)
	}
}

func TestWhereParams(t *testing.T) {
	t.Parallel()

	columns := map[string]interface{}{
		"name": "bob",
		"id":   5,
	}

	result := WhereParams(columns)

	if result[0].(int) != 5 {
		t.Error("Result[0] was wrong, got:", result[0])
	}

	if result[1].(string) != "bob" {
		t.Error("Result[1] was wrong, got:", result[1])
	}
}
