package queries

import (
	"reflect"
	"testing"
	"time"

	"github.com/volatiletech/null"
)

type testObj struct {
	ID       int
	Name     string `db:"TestHello"`
	HeadSize int
}

func TestNonZeroDefaultSet(t *testing.T) {
	t.Parallel()

	type Anything struct {
		ID        int        `boil:"id"`
		Name      string     `boil:"name"`
		CreatedAt *time.Time `boil:"created_at"`
		UpdatedAt null.Time  `boil:"updated_at"`
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
