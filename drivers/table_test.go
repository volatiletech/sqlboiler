package drivers

import "testing"

func TestGetTable(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{Name: "one"},
	}

	tbl := GetTable(tables, "one")

	if tbl.Name != "one" {
		t.Error("didn't get column")
	}
}

func TestGetTableMissing(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{Name: "one"},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic failure")
		}
	}()

	GetTable(tables, "missing")
}

func TestGetColumn(t *testing.T) {
	t.Parallel()

	table := Table{
		Columns: []Column{
			{Name: "one"},
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
			{Name: "one"},
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic failure")
		}
	}()

	table.GetColumn("missing")
}

func TestCanLastInsertID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Can   bool
		PKeys []Column
	}{
		{true, []Column{
			{Name: "id", Type: "int64", Default: "a"},
		}},
		{true, []Column{
			{Name: "id", Type: "uint64", Default: "a"},
		}},
		{true, []Column{
			{Name: "id", Type: "int", Default: "a"},
		}},
		{true, []Column{
			{Name: "id", Type: "uint", Default: "a"},
		}},
		{true, []Column{
			{Name: "id", Type: "uint", Default: "a"},
		}},
		{false, []Column{
			{Name: "id", Type: "uint", Default: "a"},
			{Name: "id2", Type: "uint", Default: "a"},
		}},
		{false, []Column{
			{Name: "id", Type: "string", Default: "a"},
		}},
		{false, []Column{
			{Name: "id", Type: "int", Default: ""},
		}},
		{false, nil},
	}

	for i, test := range tests {
		table := Table{
			Columns: test.PKeys,
			PKey:    &PrimaryKey{},
		}

		var pkeyNames []string
		for _, pk := range test.PKeys {
			pkeyNames = append(pkeyNames, pk.Name)
		}
		table.PKey.Columns = pkeyNames

		if got := table.CanLastInsertID(); got != test.Can {
			t.Errorf("%d) wrong: %t", i, got)
		}
	}
}

func TestCanSoftDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Can     bool
		Columns []Column
	}{
		{true, []Column{
			{Name: "deleted_at", Type: "null.Time"},
		}},
		{false, []Column{
			{Name: "deleted_at", Type: "time.Time"},
		}},
		{false, []Column{
			{Name: "deleted_at", Type: "int"},
		}},
		{false, nil},
	}

	for i, test := range tests {
		table := Table{
			Columns: test.Columns,
		}

		if got := table.CanSoftDelete("deleted_at"); got != test.Can {
			t.Errorf("%d) wrong: %t", i, got)
		}
	}
}
