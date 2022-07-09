package drivers

import (
	"testing"

	"github.com/volatiletech/strmangle"
)

type testMockDriver struct{}

func (m testMockDriver) TranslateColumnType(c Column) Column { return c }
func (m testMockDriver) UseLastInsertID() bool               { return false }
func (m testMockDriver) UseTopClause() bool                  { return false }
func (m testMockDriver) Open() error                         { return nil }
func (m testMockDriver) Close()                              {}

func (m testMockDriver) TableNames(schema string, whitelist, blacklist []string) ([]string, error) {
	if len(whitelist) > 0 {
		return whitelist, nil
	}
	tables := []string{"pilots", "jets", "airports", "licenses", "hangars", "languages", "pilot_languages"}
	return strmangle.SetComplement(tables, blacklist), nil
}

// Columns returns a list of mock columns
func (m testMockDriver) Columns(schema, tableName string, whitelist, blacklist []string) ([]Column, error) {
	return map[string][]Column{
		"pilots": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "name", Type: "string", DBType: "character"},
		},
		"airports": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "size", Type: "null.Int", DBType: "integer", Nullable: true},
		},
		"jets": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "pilot_id", Type: "int", DBType: "integer", Nullable: true, Unique: true},
			{Name: "airport_id", Type: "int", DBType: "integer"},
			{Name: "name", Type: "string", DBType: "character", Nullable: false},
			{Name: "color", Type: "null.String", DBType: "character", Nullable: true},
			{Name: "uuid", Type: "string", DBType: "uuid", Nullable: true},
			{Name: "identifier", Type: "string", DBType: "uuid", Nullable: false},
			{Name: "cargo", Type: "[]byte", DBType: "bytea", Nullable: false},
			{Name: "manifest", Type: "[]byte", DBType: "bytea", Nullable: true, Unique: true},
		},
		"licenses": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "pilot_id", Type: "int", DBType: "integer"},
		},
		"hangars": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "name", Type: "string", DBType: "character", Nullable: true, Unique: true},
			{Name: "hangar_id", Type: "int", DBType: "integer", Nullable: true},
		},
		"languages": {
			{Name: "id", Type: "int", DBType: "integer"},
			{Name: "language", Type: "string", DBType: "character", Nullable: false, Unique: true},
		},
		"pilot_languages": {
			{Name: "pilot_id", Type: "int", DBType: "integer"},
			{Name: "language_id", Type: "int", DBType: "integer"},
		},
	}[tableName], nil
}

// ForeignKeyInfo returns a list of mock foreignkeys
func (m testMockDriver) ForeignKeyInfo(schema, tableName string) ([]ForeignKey, error) {
	return map[string][]ForeignKey{
		"jets": {
			{Table: "jets", Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", ForeignColumnUnique: true},
			{Table: "jets", Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id"},
		},
		"licenses": {
			{Table: "licenses", Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
		},
		"pilot_languages": {
			{Table: "pilot_languages", Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			{Table: "pilot_languages", Name: "jet_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id"},
		},
		"hangars": {
			{Table: "hangars", Name: "hangar_fk_id", Column: "hangar_id", ForeignTable: "hangars", ForeignColumn: "id"},
		},
	}[tableName], nil
}

// PrimaryKeyInfo returns mock primary key info for the passed in table name
func (m testMockDriver) PrimaryKeyInfo(schema, tableName string) (*PrimaryKey, error) {
	return map[string]*PrimaryKey{
		"pilots":          {Name: "pilot_id_pkey", Columns: []string{"id"}},
		"airports":        {Name: "airport_id_pkey", Columns: []string{"id"}},
		"jets":            {Name: "jet_id_pkey", Columns: []string{"id"}},
		"licenses":        {Name: "license_id_pkey", Columns: []string{"id"}},
		"hangars":         {Name: "hangar_id_pkey", Columns: []string{"id"}},
		"languages":       {Name: "language_id_pkey", Columns: []string{"id"}},
		"pilot_languages": {Name: "pilot_languages_pkey", Columns: []string{"pilot_id", "language_id"}},
	}[tableName], nil
}

// RightQuote is the quoting character for the right side of the identifier
func (m testMockDriver) RightQuote() byte {
	return '"'
}

// LeftQuote is the quoting character for the left side of the identifier
func (m testMockDriver) LeftQuote() byte {
	return '"'
}

// UseIndexPlaceholders returns true to indicate fake support of indexed placeholders
func (m testMockDriver) UseIndexPlaceholders() bool {
	return false
}

func TestTables(t *testing.T) {
	t.Parallel()

	tables, err := Tables(testMockDriver{}, "public", nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(tables) != 7 {
		t.Errorf("Expected len 7, got: %d\n", len(tables))
	}

	prev := ""
	for i := range tables {
		if prev >= tables[i].Name {
			t.Error("tables are not sorted")
		}
		prev = tables[i].Name
	}

	pilots := GetTable(tables, "pilots")
	if len(pilots.Columns) != 2 {
		t.Error()
	}
	if pilots.ToOneRelationships[0].ForeignTable != "jets" {
		t.Error("want a to many to jets")
	}
	if pilots.ToManyRelationships[0].ForeignTable != "licenses" {
		t.Error("want a to many to languages")
	}
	if pilots.ToManyRelationships[1].ForeignTable != "languages" {
		t.Error("want a to many to languages")
	}

	jets := GetTable(tables, "jets")
	if len(jets.ToManyRelationships) != 0 {
		t.Error("want no to many relationships")
	}

	languages := GetTable(tables, "pilot_languages")
	if !languages.IsJoinTable {
		t.Error("languages is a join table")
	}

	hangars := GetTable(tables, "hangars")
	if len(hangars.ToManyRelationships) != 1 || hangars.ToManyRelationships[0].ForeignTable != "hangars" {
		t.Error("want 1 to many relationships")
	}
	if len(hangars.FKeys) != 1 || hangars.FKeys[0].ForeignTable != "hangars" {
		t.Error("want one hangar foreign key to itself")
	}
}

func TestFilterForeignKeys(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name: "one",
			Columns: []Column{
				{Name: "id"},
				{Name: "two_id"},
				{Name: "three_id"},
				{Name: "four_id"},
			},
			FKeys: []ForeignKey{
				{Table: "one", Column: "two_id", ForeignTable: "two", ForeignColumn: "id"},
				{Table: "one", Column: "three_id", ForeignTable: "three", ForeignColumn: "id"},
				{Table: "one", Column: "four_id", ForeignTable: "four", ForeignColumn: "id"},
			},
		},
		{
			Name: "two",
			Columns: []Column{
				{Name: "id"},
			},
		},
		{
			Name: "three",
			Columns: []Column{
				{Name: "id"},
			},
		},
		{
			Name: "four",
			Columns: []Column{
				{Name: "id"},
			},
		},
	}

	tests := []struct {
		Whitelist   []string
		Blacklist   []string
		ExpectFkNum int
	}{
		{[]string{}, []string{}, 3},
		{[]string{"one", "two", "three"}, []string{}, 2},
		{[]string{"one.two_id", "two"}, []string{}, 1},
		{[]string{"*.two_id", "two"}, []string{}, 1},
		{[]string{}, []string{"three", "four"}, 1},
		{[]string{}, []string{"three.id"}, 2},
		{[]string{}, []string{"one.two_id"}, 2},
		{[]string{}, []string{"*.two_id"}, 2},
		{[]string{"one", "two"}, []string{"two"}, 0},
	}

	for i, test := range tests {
		table := tables[0]
		filterForeignKeys(&table, test.Whitelist, test.Blacklist)
		if fkNum := len(table.FKeys); fkNum != test.ExpectFkNum {
			t.Errorf("%d) want: %d, got: %d\nTest: %#v", i, test.ExpectFkNum, fkNum, test)
		}
	}
}

func TestKnownColumn(t *testing.T) {
	tests := []struct {
		table     string
		column    string
		whitelist []string
		blacklist []string
		expected  bool
	}{
		{"one", "id", []string{"one"}, []string{}, true},
		{"one", "id", []string{}, []string{"one"}, false},
		{"one", "id", []string{"one.id"}, []string{}, true},
		{"one", "id", []string{"one.id"}, []string{"one"}, false},
		{"one", "id", []string{"two"}, []string{}, false},
		{"one", "id", []string{"two"}, []string{"one"}, false},
		{"one", "id", []string{"two.id"}, []string{}, false},
		{"one", "id", []string{"*.id"}, []string{}, true},
		{"one", "id", []string{"*.id"}, []string{"*.id"}, false},
	}

	for i, test := range tests {
		known := knownColumn(test.table, test.column, test.whitelist, test.blacklist)
		if known != test.expected {
			t.Errorf("%d) want: %t, got: %t\nTest: %#v", i, test.expected, known, test)
		}
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

func TestSetForeignKeyConstraints(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name: "one",
			Columns: []Column{
				{Name: "id1", Type: "string", Nullable: false, Unique: false},
				{Name: "id2", Type: "string", Nullable: true, Unique: true},
			},
		},
		{
			Name: "other",
			Columns: []Column{
				{Name: "one_id_1", Type: "string", Nullable: false, Unique: false},
				{Name: "one_id_2", Type: "string", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Column: "one_id_1", ForeignTable: "one", ForeignColumn: "id1"},
				{Column: "one_id_2", ForeignTable: "one", ForeignColumn: "id2"},
			},
		},
	}

	setForeignKeyConstraints(&tables[0], tables)
	setForeignKeyConstraints(&tables[1], tables)

	first := tables[1].FKeys[0]
	second := tables[1].FKeys[1]
	if first.Nullable {
		t.Error("should not be nullable")
	}
	if first.Unique {
		t.Error("should not be unique")
	}
	if first.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if first.ForeignColumnUnique {
		t.Error("should be unique")
	}
	if !second.Nullable {
		t.Error("should be nullable")
	}
	if !second.Unique {
		t.Error("should be unique")
	}
	if !second.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if !second.ForeignColumnUnique {
		t.Error("should be unique")
	}
}

func TestSetRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name: "one",
			Columns: []Column{
				{Name: "id", Type: "string"},
			},
		},
		{
			Name: "other",
			Columns: []Column{
				{Name: "other_id", Type: "string"},
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
