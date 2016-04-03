package dbdrivers

import (
	"reflect"
	"testing"
)

type testInterface struct{}

func (t testInterface) TableNames() ([]string, error) {
	return []string{"table1", "table2"}, nil
}

func (t testInterface) Columns(tableName string) ([]Column, error) {
	return []Column{
		Column{Name: "col1", Type: "character varying"},
		Column{Name: "col2", Type: "character varying"},
	}, nil
}

func (t testInterface) PrimaryKeyInfo(tableName string) (*PrimaryKey, error) {
	return &PrimaryKey{Name: "pkey1", Columns: []string{"col1", "col2"}}, nil
}

func (t testInterface) ForeignKeyInfo(tableName string) ([]ForeignKey, error) {
	return []ForeignKey{
		{
			Name:          "fkey1",
			Column:        "col1",
			ForeignTable:  "table3",
			ForeignColumn: "col3",
		},
		{
			Name:          "fkey2",
			Column:        "col2",
			ForeignTable:  "table3",
			ForeignColumn: "col3",
		},
	}, nil
}

func (t testInterface) TranslateColumnType(column Column) Column {
	column.Type = "string"
	return column
}

func (t testInterface) Open() error {
	return nil
}

func (t testInterface) Close() {}

func TestTables(t *testing.T) {
	t.Parallel()

	tables, err := Tables(testInterface{})
	if err != nil {
		t.Error(err)
	}

	if len(tables) != 2 {
		t.Errorf("Expected len 2, got: %d\n", len(tables))
	}

	expectCols := []Column{
		Column{Name: "col1", Type: "string"},
		Column{Name: "col2", Type: "string"},
	}

	if !reflect.DeepEqual(tables[0].Columns, expectCols) {
		t.Errorf("Did not get expected columns, got:\n%#v\n%#v", tables[0].Columns, expectCols)
	}

	if !tables[0].IsJoinTable || !tables[1].IsJoinTable {
		t.Errorf("Expected IsJoinTable to be true")
	}

	expectPkey := &PrimaryKey{Name: "pkey1", Columns: []string{"col1", "col2"}}
	expectFkey := []ForeignKey{
		{
			Name:          "fkey1",
			Column:        "col1",
			ForeignTable:  "table3",
			ForeignColumn: "col3",
		},
		{
			Name:          "fkey2",
			Column:        "col2",
			ForeignTable:  "table3",
			ForeignColumn: "col3",
		},
	}

	if !reflect.DeepEqual(tables[0].FKeys, expectFkey) {
		t.Errorf("Did not get expected Fkey, got:\n%#v\n%#v", tables[0].FKeys, expectFkey)
	}

	if !reflect.DeepEqual(tables[0].PKey, expectPkey) {
		t.Errorf("Did not get expected PKey, got:\n#%v\n%#v", tables[0].PKey, expectPkey)
	}
}

func TestSetIsJoinTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Pkey   []string
		Fkey   []string
		Should bool
	}{
		{Pkey: []string{"one", "two"}, Fkey: []string{"one", "two"}, Should: true},
		{Pkey: []string{"two", "one"}, Fkey: []string{"one", "two"}, Should: true},

		{Pkey: []string{"one"}, Fkey: []string{"one"}, Should: false},
		{Pkey: []string{"one", "two", "three"}, Fkey: []string{"one", "two"}, Should: false},
		{Pkey: []string{"one", "two", "three"}, Fkey: []string{"one", "two", "three"}, Should: false},
		{Pkey: []string{"one"}, Fkey: []string{"one", "two"}, Should: false},
		{Pkey: []string{"one", "two"}, Fkey: []string{"one"}, Should: false},
	}

	for i, test := range tests {
		var table Table

		table.PKey = &PrimaryKey{Columns: test.Pkey}
		for _, k := range test.Fkey {
			table.FKeys = append(table.FKeys, ForeignKey{Column: k})
		}

		setIsJoinTable(&table)
		if is := table.IsJoinTable; is != test.Should {
			t.Errorf("%d) want: %t, got: %t\nTest: %#v", i, test.Should, is, test)
		}
	}
}
