package strmangle

import (
	"reflect"
	"testing"
)

func TestUpdateColumnSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Cols      []string
		PKeys     []string
		Whitelist []string
		Out       []string
	}{
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Out: []string{"b"}},
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Whitelist: []string{"a"}, Out: []string{"a"}},
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Whitelist: []string{"a", "b"}, Out: []string{"a", "b"}},
	}

	for i, test := range tests {
		set := UpdateColumnSet(test.Cols, test.PKeys, test.Whitelist)

		if !reflect.DeepEqual(set, test.Out) {
			t.Errorf("%d) set was wrong\nwant: %v\ngot:  %v", i, test.Out, set)
		}
	}
}

func TestInsertColumnSet(t *testing.T) {
	t.Parallel()

	columns := []string{"a", "b", "c"}
	defaults := []string{"a", "c"}
	nodefaults := []string{"b"}

	tests := []struct {
		Cols            []string
		Defaults        []string
		NoDefaults      []string
		NonZeroDefaults []string
		Whitelist       []string
		Set             []string
		Ret             []string
	}{
		// No whitelist
		{Set: []string{"b"}, Ret: []string{"a", "c"}},
		{Defaults: []string{}, NoDefaults: []string{"a", "b", "c"}, Set: []string{"a", "b", "c"}, Ret: []string{}},

		// No whitelist + Nonzero defaults
		{NonZeroDefaults: []string{"a"}, Set: []string{"a", "b"}, Ret: []string{"c"}},
		{NonZeroDefaults: []string{"c"}, Set: []string{"b", "c"}, Ret: []string{"a"}},

		// Whitelist
		{Whitelist: []string{"a"}, Set: []string{"a"}, Ret: []string{"c"}},
		{Whitelist: []string{"c"}, Set: []string{"c"}, Ret: []string{"a"}},
		{Whitelist: []string{"a", "c"}, Set: []string{"a", "c"}, Ret: []string{}},
		{Whitelist: []string{"a", "b", "c"}, Set: []string{"a", "b", "c"}, Ret: []string{}},

		// Whitelist + Nonzero defaults (shouldn't care, same results as above)
		{Whitelist: []string{"a"}, NonZeroDefaults: []string{"c"}, Set: []string{"a"}, Ret: []string{"c"}},
		{Whitelist: []string{"c"}, NonZeroDefaults: []string{"b"}, Set: []string{"c"}, Ret: []string{"a"}},
	}

	for i, test := range tests {
		if test.Cols == nil {
			test.Cols = columns
		}
		if test.Defaults == nil {
			test.Defaults = defaults
		}
		if test.NoDefaults == nil {
			test.NoDefaults = nodefaults
		}

		set, ret := InsertColumnSet(test.Cols, test.Defaults, test.NoDefaults, test.NonZeroDefaults, test.Whitelist)

		if !reflect.DeepEqual(set, test.Set) {
			t.Errorf("%d) set was wrong\nwant: %v\ngot:  %v", i, test.Set, set)
		}
		if !reflect.DeepEqual(ret, test.Ret) {
			t.Errorf("%d) ret was wrong\nwant: %v\ngot:  %v", i, test.Ret, ret)
		}
	}
}

func TestSetInclude(t *testing.T) {
	t.Parallel()

	elements := []string{"one", "two"}
	if got := SetInclude("one", elements); !got {
		t.Error("should have found element key")
	}
	if got := SetInclude("three", elements); got {
		t.Error("should not have found element key")
	}
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
