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
