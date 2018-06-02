package strmangle

import (
	"reflect"
	"testing"
)

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
