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

func TestSQLDefStrings(t *testing.T) {
	t.Parallel()

	cols := []Column{
		{Name: "one", Type: "int64"},
		{Name: "two", Type: "string"},
		{Name: "three", Type: "string"},
	}

	defs := SQLColDefinitions(cols, []string{"one", "three"})
	strs := SQLColDefStrings(defs)

	if got := strs[0]; got != "one int64" {
		t.Error("wrong str:", got)
	}
	if got := strs[1]; got != "three string" {
		t.Error("wrong str:", got)
	}
}
