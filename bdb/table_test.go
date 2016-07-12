package bdb

import "testing"

func TestGetColumn(t *testing.T) {
	t.Parallel()

	table := Table{
		Columns: []Column{
			Column{Name: "one"},
		},
	}

	c := table.GetColumn("one")

	if c.Name != "one" {
		t.Error("didn't get column")
	}
}

func TestGetColumnMissing(t *testing.T) {
	t.Parallel()

	table := Table{
		Columns: []Column{
			Column{Name: "one"},
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic failure")
		}
	}()

	table.GetColumn("missing")
}
