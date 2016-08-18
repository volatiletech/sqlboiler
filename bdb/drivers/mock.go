package drivers

import (
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/strmangle"
)

// MockDriver is a mock implementation of the bdb driver Interface
type MockDriver int

// TableNames returns a list of mock table names
func (MockDriver) TableNames(exclude []string) ([]string, error) {
	tables := []string{"pilots", "jets", "airports", "licenses", "pilots_jets_tags"}
	return strmangle.SetComplement(tables, exclude), nil
}

// Columns returns a list of mock columns
func (MockDriver) Columns(tableName string) ([]bdb.Column, error) {
	return map[string][]bdb.Column{
		"pilots":   {{Name: "id", Type: "int32"}},
		"airports": {{Name: "id", Type: "int32", Nullable: true}},
		"jets": {
			{Name: "id", Type: "int32"},
			{Name: "pilot_id", Type: "int32", Nullable: true, Unique: true},
			{Name: "airport_id", Type: "int32"},
		},
		"licenses": {
			{Name: "pilot_id", Type: "int32"},
			{Name: "source_id", Type: "int32", Nullable: true},
		},
		"pilots_jets_tags": {
			{Name: "pilot_id", Type: "int32"},
			{Name: "jet_id", Type: "int32"},
		},
	}[tableName], nil
}

// ForeignKeyInfo returns a list of mock foreignkeys
func (MockDriver) ForeignKeyInfo(tableName string) ([]bdb.ForeignKey, error) {
	return map[string][]bdb.ForeignKey{
		"jets": {
			{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id"},
		},
		"licenses": {
			{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			{Name: "licenses_source_id_fk", Column: "source_id", ForeignTable: "pilots", ForeignColumn: "id"},
		},
		"pilots_jets_tags": {
			{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			{Name: "jet_id_fk", Column: "jet_id", ForeignTable: "jets", ForeignColumn: "id"},
		},
	}[tableName], nil
}

// TranslateColumnType converts a column to its "null." form if it is nullable
func (MockDriver) TranslateColumnType(c bdb.Column) bdb.Column {
	if c.Nullable {
		c.Type = "null." + strmangle.TitleCase(c.Type)
	}
	return c
}

// PrimaryKeyInfo returns mock primary key info for the passed in table name
func (MockDriver) PrimaryKeyInfo(tableName string) (*bdb.PrimaryKey, error) {
	return map[string]*bdb.PrimaryKey{
		"pilots_jets_tags": {
			Name:    "pilot_jet_id_pkey",
			Columns: []string{"pilot_id", "jet_id"},
		},
	}[tableName], nil
}

// UseLastInsertID returns a database mock LastInsertID compatability flag
func (MockDriver) UseLastInsertID() bool { return false }

// Open mimics a database open call and returns nil for no error
func (MockDriver) Open() error { return nil }

// Close mimics a database close call
func (MockDriver) Close() {}
