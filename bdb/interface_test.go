package bdb

import (
	"reflect"
	"testing"
)

type testInterface struct{}

func (t testInterface) TableNames() ([]string, error) {
	return []string{"table1", "table2"}, nil
}

var testCols = []Column{
	Column{Name: "col1", Type: "character varying"},
	Column{Name: "col2", Type: "character varying", Nullable: true},
}

func (t testInterface) Columns(tableName string) ([]Column, error) {
	return testCols, nil
}

var testPkey = &PrimaryKey{Name: "pkey1", Columns: []string{"col1", "col2"}}

func (t testInterface) PrimaryKeyInfo(tableName string) (*PrimaryKey, error) {
	return testPkey, nil
}

var testFkeys = []ForeignKey{
	{
		Name:          "fkey1",
		Column:        "col1",
		ForeignTable:  "table2",
		ForeignColumn: "col2",
	},
	{
		Name:          "fkey2",
		Column:        "col2",
		ForeignTable:  "table1",
		ForeignColumn: "col1",
	},
}

func (t testInterface) ForeignKeyInfo(tableName string) ([]ForeignKey, error) {
	return testFkeys, nil
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

	if !reflect.DeepEqual(tables[0].Columns, testCols) {
		t.Errorf("Did not get expected columns, got:\n%#v\n%#v", tables[0].Columns, testCols)
	}

	if !tables[0].IsJoinTable || !tables[1].IsJoinTable {
		t.Errorf("Expected IsJoinTable to be true")
	}

	if !reflect.DeepEqual(tables[0].PKey, testPkey) {
		t.Errorf("Did not get expected PKey, got:\n#%v\n%#v", tables[0].PKey, testPkey)
	}

	if !reflect.DeepEqual(tables[0].FKeys, testFkeys) {
		t.Errorf("Did not get expected Fkey, got:\n%#v\n%#v", tables[0].FKeys, testFkeys)
	}

	if len(tables[0].ToManyRelationships) != 1 {
		t.Error("wanted a to many relationship")
	}
	if len(tables[1].ToManyRelationships) != 1 {
		t.Error("wanted a to many relationship")
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

func TestSetForeignKeyNullability(t *testing.T) {
	t.Parallel()

	table := &Table{
		Columns: []Column{
			Column{Name: "col1", Type: "string"},
			Column{Name: "col2", Type: "string", Nullable: true},
		},
		FKeys: []ForeignKey{
			{
				Column: "col1",
			},
			{
				Column: "col2",
			},
		},
	}

	setForeignKeyNullability(table)

	if table.FKeys[0].Nullable {
		t.Error("should not be nullable")
	}
	if !table.FKeys[1].Nullable {
		t.Error("should be nullable")
	}
}

func TestSetRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		Table{
			Name: "one",
			Columns: []Column{
				Column{Name: "id", Type: "string"},
			},
		},
		Table{
			Name: "other",
			Columns: []Column{
				Column{Name: "other_id", Type: "string"},
			},
			FKeys: []ForeignKey{{Column: "other_id", ForeignTable: "one", ForeignColumn: "id", Nullable: true}},
		},
	}

	setRelationships(&tables[0], tables)
	setRelationships(&tables[1], tables)

	if got := len(tables[0].ToManyRelationships); got != 1 {
		t.Error("should have a relationship:", got)
	}
	if got := len(tables[1].ToManyRelationships); got != 0 {
		t.Error("should have no to many relationships:", got)
	}

	rel := tables[0].ToManyRelationships[0]
	if rel.Column != "id" {
		t.Error("wrong column:", rel.Column)
	}
	if rel.ForeignTable != "other" {
		t.Error("wrong table:", rel.ForeignTable)
	}
	if rel.ForeignColumn != "other_id" {
		t.Error("wrong column:", rel.ForeignColumn)
	}
	if rel.ToJoinTable {
		t.Error("should not be a join table")
	}
}
