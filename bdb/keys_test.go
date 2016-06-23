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

func TestAutoIncPrimaryKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		Ok      bool
		Expect  Column
		Pkey    *PrimaryKey
		Columns []Column
	}{
		"nillcase": {
			Ok:      false,
			Pkey:    nil,
			Columns: nil,
		},
		"easycase": {
			Ok:      true,
			Expect:  Column{Name: "one", Type: "int32", IsNullable: false, Default: `nextval('abc'::regclass)`},
			Pkey:    &PrimaryKey{Name: "pkey", Columns: []string{"one"}},
			Columns: []Column{Column{Name: "one", Type: "int32", IsNullable: false, Default: `nextval('abc'::regclass)`}},
		},
		"missingcase": {
			Ok:      false,
			Pkey:    &PrimaryKey{Name: "pkey", Columns: []string{"two"}},
			Columns: []Column{Column{Name: "one", Type: "int32", IsNullable: false, Default: `nextval('abc'::regclass)`}},
		},
		"wrongtype": {
			Ok:      false,
			Pkey:    &PrimaryKey{Name: "pkey", Columns: []string{"one"}},
			Columns: []Column{Column{Name: "one", Type: "string", IsNullable: false, Default: `nextval('abc'::regclass)`}},
		},
		"nodefault": {
			Ok:      false,
			Pkey:    &PrimaryKey{Name: "pkey", Columns: []string{"one"}},
			Columns: []Column{Column{Name: "one", Type: "string", IsNullable: false, Default: ``}},
		},
		"nullable": {
			Ok:      false,
			Pkey:    &PrimaryKey{Name: "pkey", Columns: []string{"one"}},
			Columns: []Column{Column{Name: "one", Type: "string", IsNullable: true, Default: `nextval('abc'::regclass)`}},
		},
	}

	for testName, test := range tests {
		pkey := AutoIncPrimaryKey(test.Columns, test.Pkey)
		if ok := (pkey != nil); ok != test.Ok {
			t.Errorf("%s) found state was wrong, want: %t, got: %t", testName, test.Ok, ok)
		} else if test.Ok && *pkey != test.Expect {
			t.Errorf("%s) wrong primary key, want: %#v, got %#v", testName, test.Expect, pkey)
		}
	}
}
