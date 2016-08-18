package bdb

import "testing"

func TestSQLColDefinitions(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "one", Type: "int64"},
		{Name: "two", Type: "string"},
		{Name: "three", Type: "string"},
	}

	defs := SQLColDefinitions(cols, []string{"one"})
	if len(defs) != 1 {
		t.Error("wrong number of defs:", len(defs))
	}
	if got := defs[0].String(); got != "one int64" {
		t.Error("wrong def:", got)
	}

	defs = SQLColDefinitions(cols, []string{"one", "three"})
	if len(defs) != 2 {
		t.Error("wrong number of defs:", len(defs))
	}
	if got := defs[0].String(); got != "one int64" {
		t.Error("wrong def:", got)
	}
	if got := defs[1].String(); got != "three string" {
		t.Error("wrong def:", got)
	}
}

func TestTypes(t *testing.T) {
	t.Parallel()

	defs := SQLColumnDefs{
		{Type: "thing1"},
		{Type: "thing2"},
	}

	ret := defs.Types()
	if ret[0] != "thing1" {
		t.Error("wrong type:", ret[0])
	}
	if ret[1] != "thing2" {
		t.Error("wrong type:", ret[1])
	}
}

func TestNames(t *testing.T) {
	t.Parallel()

	defs := SQLColumnDefs{
		{Name: "thing1"},
		{Name: "thing2"},
	}

	ret := defs.Names()
	if ret[0] != "thing1" {
		t.Error("wrong type:", ret[0])
	}
	if ret[1] != "thing2" {
		t.Error("wrong type:", ret[1])
	}
}
