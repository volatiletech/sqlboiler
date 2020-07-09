package boilingcore

import (
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/v4/drivers"
)

func TestOrders(t *testing.T) {
	t.Parallel()

	columns := []drivers.Column{
		{Name: "deleted_at"},
		{Name: "name"},
		{Name: "created_at"},
		{Name: "updated_at"},
		{Name: "code"},
		{Name: "id"},
	}

	t.Run("Sort", func(t *testing.T) {
		expect := []drivers.Column{
			{Name: "id"},
			{Name: "name"},
			{Name: "code"},
			{Name: "created_at"},
			{Name: "updated_at"},
			{Name: "deleted_at"},
		}

		o := Orders{
			Columns: map[string]int{
				"id":         -1,
				"created_at": 3,
				"updated_at": 4,
				"deleted_at": 5,
			},
		}
		if got := o.sortColumns(columns); !reflect.DeepEqual(expect, got) {
			t.Errorf("it should sorted: %#v", got)
		}
	})

	t.Run("NoSort", func(t *testing.T) {
		expect := []drivers.Column{
			{Name: "deleted_at"},
			{Name: "name"},
			{Name: "created_at"},
			{Name: "updated_at"},
			{Name: "code"},
			{Name: "id"},
		}

		o := Orders{}
		if got := o.sortColumns(columns); !reflect.DeepEqual(expect, got) {
			t.Errorf("it should sorted: %#v", got)
		}
	})
}
