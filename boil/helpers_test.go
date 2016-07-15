package boil

import (
	"reflect"
	"testing"
	"time"

	"gopkg.in/nullbio/null.v4"
)

type testObj struct {
	ID       int
	Name     string `db:"TestHello"`
	HeadSize int
}

func TestSetComplement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		A []string
		B []string
		C []string
	}{
		{
			[]string{"thing1", "thing2", "thing3"},
			[]string{"thing2", "otherthing", "stuff"},
			[]string{"thing1", "thing3"},
		},
		{
			[]string{},
			[]string{"thing1", "thing2"},
			[]string{},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{},
			[]string{"thing1", "thing2"},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{"thing1", "thing2"},
			[]string{},
		},
	}

	for i, test := range tests {
		c := SetComplement(test.A, test.B)
		if !reflect.DeepEqual(test.C, c) {
			t.Errorf("[%d] mismatch:\nWant: %#v\nGot:  %#v", i, test.C, c)
		}
	}
}

func TestSetIntersect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		A []string
		B []string
		C []string
	}{
		{
			[]string{"thing1", "thing2", "thing3"},
			[]string{"thing2", "otherthing", "stuff"},
			[]string{"thing2"},
		},
		{
			[]string{},
			[]string{"thing1", "thing2"},
			[]string{},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{},
			[]string{},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{"thing1", "thing2"},
			[]string{"thing1", "thing2"},
		},
	}

	for i, test := range tests {
		c := SetIntersect(test.A, test.B)
		if !reflect.DeepEqual(test.C, c) {
			t.Errorf("[%d] mismatch:\nWant: %#v\nGot:  %#v", i, test.C, c)
		}
	}
}

func TestSetMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		A []string
		B []string
		C []string
	}{
		{
			[]string{"thing1", "thing2", "thing3"},
			[]string{"thing1", "thing3", "thing4"},
			[]string{"thing1", "thing2", "thing3", "thing4"},
		},
		{
			[]string{},
			[]string{"thing1", "thing2"},
			[]string{"thing1", "thing2"},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{},
			[]string{"thing1", "thing2"},
		},
		{
			[]string{"thing1", "thing2", "thing3"},
			[]string{"thing1", "thing2", "thing3"},
			[]string{"thing1", "thing2", "thing3"},
		},
		{
			[]string{"thing1", "thing2"},
			[]string{"thing3", "thing4"},
			[]string{"thing1", "thing2", "thing3", "thing4"},
		},
	}

	for i, test := range tests {
		m := SetMerge(test.A, test.B)
		if !reflect.DeepEqual(test.C, m) {
			t.Errorf("[%d] mismatch:\nWant: %#v\nGot: %#v", i, test.C, m)
		}
	}
}

func TestNonZeroDefaultSet(t *testing.T) {
	t.Parallel()

	type Anything struct {
		ID        int
		Name      string
		CreatedAt *time.Time
		UpdatedAt null.Time
	}

	now := time.Now()

	tests := []struct {
		Defaults []string
		Obj      interface{}
		Ret      []string
	}{
		{
			[]string{"id"},
			Anything{Name: "hi", CreatedAt: nil, UpdatedAt: null.Time{Valid: false}},
			[]string{},
		},
		{
			[]string{"id"},
			Anything{ID: 5, Name: "hi", CreatedAt: nil, UpdatedAt: null.Time{Valid: false}},
			[]string{"id"},
		},
		{
			[]string{},
			Anything{ID: 5, Name: "hi", CreatedAt: nil, UpdatedAt: null.Time{Valid: false}},
			[]string{},
		},
		{
			[]string{"id", "created_at", "updated_at"},
			Anything{ID: 5, Name: "hi", CreatedAt: nil, UpdatedAt: null.Time{Valid: false}},
			[]string{"id"},
		},
		{
			[]string{"id", "created_at", "updated_at"},
			Anything{ID: 5, Name: "hi", CreatedAt: &now, UpdatedAt: null.Time{Valid: true, Time: time.Now()}},
			[]string{"id", "created_at", "updated_at"},
		},
	}

	for i, test := range tests {
		z := NonZeroDefaultSet(test.Defaults, test.Obj)
		if !reflect.DeepEqual(test.Ret, z) {
			t.Errorf("[%d] mismatch:\nWant: %#v\nGot:  %#v", i, test.Ret, z)
		}
	}
}

func TestSortByKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Keys []string
		Strs []string
		Ret  []string
	}{
		{
			[]string{"id", "name", "thing", "stuff"},
			[]string{"thing", "stuff", "name", "id"},
			[]string{"id", "name", "thing", "stuff"},
		},
		{
			[]string{"id", "name", "thing", "stuff"},
			[]string{"id", "name", "thing", "stuff"},
			[]string{"id", "name", "thing", "stuff"},
		},
		{
			[]string{"id", "name", "thing", "stuff"},
			[]string{"stuff", "thing"},
			[]string{"thing", "stuff"},
		},
	}

	for i, test := range tests {
		z := SortByKeys(test.Keys, test.Strs)
		if !reflect.DeepEqual(test.Ret, z) {
			t.Errorf("[%d] mismatch:\nWant: %#v\nGot:  %#v", i, test.Ret, z)
		}
	}
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

	columns := []string{
		"id",
		"name",
		"date",
	}

	result := WhereClause(columns)

	if result != `id=$1 AND name=$2 AND date=$3` {
		t.Error("Result was wrong, got:", result)
	}
}
