package main

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/bdb/drivers"
)

func TestTextsFromForeignKey(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(drivers.MockDriver(0))
	if err != nil {
		t.Fatal(err)
	}

	jets := bdb.GetTable(tables, "jets")
	texts := textsFromForeignKey("models", tables, jets, jets.FKeys[0])
	expect := RelationshipToOneTexts{}

	expect.ForeignKey = jets.FKeys[0]

	expect.LocalTable.NameGo = "Jet"
	expect.LocalTable.ColumnNameGo = "PilotID"

	expect.ForeignTable.Name = "pilots"
	expect.ForeignTable.NameGo = "Pilot"
	expect.ForeignTable.NamePluralGo = "Pilots"
	expect.ForeignTable.ColumnName = "id"
	expect.ForeignTable.ColumnNameGo = "ID"

	expect.Function.PackageName = "models"
	expect.Function.Name = "Pilot"
	expect.Function.Varname = "pilot"
	expect.Function.Receiver = "j"
	expect.Function.ReverseInserts = false

	expect.Function.LocalAssignment = "PilotID.Int32"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTextsFromOneToOneRelationship(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(drivers.MockDriver(0))
	if err != nil {
		t.Fatal(err)
	}

	pilots := bdb.GetTable(tables, "pilots")
	texts := textsFromOneToOneRelationship("models", tables, pilots, pilots.ToManyRelationships[0])
	expect := RelationshipToOneTexts{}

	expect.ForeignKey = bdb.ForeignKey{
		Name:     "none",
		Column:   "id",
		Nullable: false,
		Unique:   false,

		ForeignTable:          "jets",
		ForeignColumn:         "pilot_id",
		ForeignColumnNullable: true,
		ForeignColumnUnique:   true,
	}

	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.ColumnNameGo = "ID"

	expect.ForeignTable.Name = "jets"
	expect.ForeignTable.NameGo = "Jet"
	expect.ForeignTable.NamePluralGo = "Jets"
	expect.ForeignTable.ColumnName = "pilot_id"
	expect.ForeignTable.ColumnNameGo = "PilotID"

	expect.Function.PackageName = "models"
	expect.Function.Name = "Jet"
	expect.Function.Varname = "jet"
	expect.Function.Receiver = "p"
	expect.Function.ReverseInserts = true

	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTextsFromRelationship(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(drivers.MockDriver(0))
	if err != nil {
		t.Fatal(err)
	}

	pilots := bdb.GetTable(tables, "pilots")
	texts := textsFromRelationship(tables, pilots, pilots.ToManyRelationships[0])
	expect := RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.NameSingular = "pilot"

	expect.ForeignTable.NameGo = "Jet"
	expect.ForeignTable.NameSingular = "jet"
	expect.ForeignTable.NamePluralGo = "Jets"
	expect.ForeignTable.NameHumanReadable = "jets"
	expect.ForeignTable.Slice = "JetSlice"

	expect.Function.Name = "Jets"
	expect.Function.Receiver = "p"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, pilots, pilots.ToManyRelationships[1])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.NameSingular = "pilot"

	expect.ForeignTable.NameGo = "License"
	expect.ForeignTable.NameSingular = "license"
	expect.ForeignTable.NamePluralGo = "Licenses"
	expect.ForeignTable.NameHumanReadable = "licenses"
	expect.ForeignTable.Slice = "LicenseSlice"

	expect.Function.Name = "Licenses"
	expect.Function.Receiver = "p"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "PilotID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, pilots, pilots.ToManyRelationships[2])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.NameSingular = "pilot"

	expect.ForeignTable.NameGo = "License"
	expect.ForeignTable.NameSingular = "license"
	expect.ForeignTable.NamePluralGo = "Licenses"
	expect.ForeignTable.NameHumanReadable = "licenses"
	expect.ForeignTable.Slice = "LicenseSlice"

	expect.Function.Name = "SourceLicenses"
	expect.Function.Receiver = "p"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "SourceID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, pilots, pilots.ToManyRelationships[3])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "Pilot"
	expect.LocalTable.NameSingular = "pilot"

	expect.ForeignTable.NameGo = "Jet"
	expect.ForeignTable.NameSingular = "jet"
	expect.ForeignTable.NamePluralGo = "Jets"
	expect.ForeignTable.NameHumanReadable = "jets"
	expect.ForeignTable.Slice = "JetSlice"

	expect.Function.Name = "Jets"
	expect.Function.Receiver = "p"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}
