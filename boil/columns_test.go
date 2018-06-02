package boil

import (
	"reflect"
	"testing"
)

func TestColumns(t *testing.T) {
	t.Parallel()

	list := Whitelist("a", "b")
	if list.Kind != ColumnsWhitelist {
		t.Error(list.Kind)
	}
	if list.Cols[0] != "a" || list.Cols[1] != "b" {
		t.Error("columns were wrong")
	}
	list = Blacklist("a", "b")
	if list.Kind != ColumnsBlacklist {
		t.Error(list.Kind)
	}
	if list.Cols[0] != "a" || list.Cols[1] != "b" {
		t.Error("columns were wrong")
	}
	list = Greylist("a", "b")
	if list.Kind != ColumnsGreylist {
		t.Error(list.Kind)
	}
	if list.Cols[0] != "a" || list.Cols[1] != "b" {
		t.Error("columns were wrong")
	}

	list = Infer()
	if list.Kind != ColumnsInfer {
		t.Error(list.Kind)
	}
	if len(list.Cols) != 0 {
		t.Error("non zero length columns")
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
		Columns         Columns
		Set             []string
		Ret             []string
	}{
		// Infer
		{Columns: Infer(), Set: []string{"b"}, Ret: []string{"a", "c"}},
		{Columns: Infer(), Defaults: []string{}, NoDefaults: []string{"a", "b", "c"}, Set: []string{"a", "b", "c"}, Ret: []string{}},

		// Infer with non-zero defaults
		{Columns: Infer(), NonZeroDefaults: []string{"a"}, Set: []string{"a", "b"}, Ret: []string{"c"}},
		{Columns: Infer(), NonZeroDefaults: []string{"c"}, Set: []string{"b", "c"}, Ret: []string{"a"}},

		// Whitelist
		{Columns: Whitelist("a"), Set: []string{"a"}, Ret: []string{"c"}},
		{Columns: Whitelist("c"), Set: []string{"c"}, Ret: []string{"a"}},
		{Columns: Whitelist("a", "c"), Set: []string{"a", "c"}, Ret: []string{}},
		{Columns: Whitelist("a", "b", "c"), Set: []string{"a", "b", "c"}, Ret: []string{}},

		// Whitelist + Nonzero defaults (shouldn't care, same results as above)
		{Columns: Whitelist("a"), NonZeroDefaults: []string{"c"}, Set: []string{"a"}, Ret: []string{"c"}},
		{Columns: Whitelist("c"), NonZeroDefaults: []string{"b"}, Set: []string{"c"}, Ret: []string{"a"}},

		// Blacklist
		{Columns: Blacklist("b"), NonZeroDefaults: []string{"c"}, Set: []string{"c"}, Ret: []string{"a"}},
		{Columns: Blacklist("c"), NonZeroDefaults: []string{"c"}, Set: []string{"b"}, Ret: []string{"a", "c"}},

		// Greylist
		{Columns: Greylist("c"), NonZeroDefaults: []string{}, Set: []string{"b", "c"}, Ret: []string{"a"}},
		{Columns: Greylist("a"), NonZeroDefaults: []string{}, Set: []string{"a", "b"}, Ret: []string{"c"}},
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

		set, ret := InsertColumnSet(test.Cols, test.Defaults, test.NoDefaults, test.NonZeroDefaults, test.Columns)

		if !reflect.DeepEqual(set, test.Set) {
			t.Errorf("%d) set was wrong\nwant: %v\ngot:  %v", i, test.Set, set)
		}
		if !reflect.DeepEqual(ret, test.Ret) {
			t.Errorf("%d) ret was wrong\nwant: %v\ngot:  %v", i, test.Ret, ret)
		}
	}
}

func TestUpdateColumnSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Cols    []string
		PKeys   []string
		Columns Columns
		Out     []string
	}{
		// Infer
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Columns: Infer(), Out: []string{"b"}},

		// Whitelist
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Columns: Whitelist("a"), Out: []string{"a"}},
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Columns: Whitelist("a", "b"), Out: []string{"a", "b"}},

		// Blacklist
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Columns: Blacklist("b"), Out: []string{}},

		// Greylist
		{Cols: []string{"a", "b"}, PKeys: []string{"a"}, Columns: Greylist("a"), Out: []string{"a", "b"}},
	}

	for i, test := range tests {
		set := UpdateColumnSet(test.Cols, test.PKeys, test.Columns)

		if !reflect.DeepEqual(set, test.Out) {
			t.Errorf("%d) set was wrong\nwant: %v\ngot:  %v", i, test.Out, set)
		}
	}
}
