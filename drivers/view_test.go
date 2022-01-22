package drivers

import "testing"

func TestGetView(t *testing.T) {
	t.Parallel()

	views := []View{
		{Name: "one"},
	}

	tbl := GetView(views, "one")

	if tbl.Name != "one" {
		t.Error("didn't get column")
	}
}

func TestGetViewMissing(t *testing.T) {
	t.Parallel()

	views := []View{
		{Name: "one"},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic failure")
		}
	}()

	GetView(views, "missing")
}

func TestGetViewColumn(t *testing.T) {
	t.Parallel()

	view := View{
		Columns: []Column{
			{Name: "one"},
		},
	}

	c := view.GetColumn("one")

	if c.Name != "one" {
		t.Error("didn't get column")
	}
}

func TestGetViewColumnMissing(t *testing.T) {
	t.Parallel()

	view := View{
		Columns: []Column{
			{Name: "one"},
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic failure")
		}
	}()

	view.GetColumn("missing")
}
