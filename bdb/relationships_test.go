package bdb

import (
	"reflect"
	"testing"
)

func TestToOneRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name:    "pilots",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "name", Unique: true}}},
		{
			Name:    "airports",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "size", Unique: true}},
		},
		{
			Name:    "jets",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "pilot_id", Unique: true}, {Name: "airport_id", Unique: true}},
			FKeys: []ForeignKey{
				{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Unique: true},
				{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id", Unique: true},
			},
		},
		{
			Name:    "licenses",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "pilot_id", Unique: true}},
			FKeys: []ForeignKey{
				{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Unique: true},
			},
		},
		{
			Name:    "hangars",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "name", Unique: true}},
		},
		{
			Name:    "languages",
			Columns: []Column{{Name: "id", Unique: true}, {Name: "language", Unique: true}},
		},
		{
			Name:        "pilot_languages",
			IsJoinTable: true,
			Columns:     []Column{{Name: "pilot_id", Unique: true}, {Name: "language_id", Unique: true}},
			FKeys: []ForeignKey{
				{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Unique: true},
				{Name: "language_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id", Unique: true},
			},
		},
	}

	relationships := ToOneRelationships("pilots", tables)

	expected := []ToOneRelationship{
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "jets",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   true,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "licenses",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   true,
		},
	}

	if len(relationships) != 2 {
		t.Error("wrong # of relationships", len(relationships))
	}

	for i, v := range relationships {
		if !reflect.DeepEqual(v, expected[i]) {
			t.Errorf("[%d] Mismatch between relationships:\n\nwant:%#v\n\ngot:%#v\n\n", i, expected[i], v)
		}
	}
}

func TestToManyRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name:    "pilots",
			Columns: []Column{{Name: "id"}, {Name: "name"}},
		},
		{
			Name:    "airports",
			Columns: []Column{{Name: "id"}, {Name: "size"}},
		},
		{
			Name:    "jets",
			Columns: []Column{{Name: "id"}, {Name: "pilot_id"}, {Name: "airport_id"}},
			FKeys: []ForeignKey{
				{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
				{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id"},
			},
		},
		{
			Name:    "licenses",
			Columns: []Column{{Name: "id"}, {Name: "pilot_id"}},
			FKeys: []ForeignKey{
				{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			},
		},
		{
			Name:    "hangars",
			Columns: []Column{{Name: "id"}, {Name: "name"}},
		},
		{
			Name:    "languages",
			Columns: []Column{{Name: "id"}, {Name: "language"}},
		},
		{
			Name:        "pilot_languages",
			IsJoinTable: true,
			Columns:     []Column{{Name: "pilot_id"}, {Name: "language_id"}},
			FKeys: []ForeignKey{
				{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
				{Name: "language_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id"},
			},
		},
	}

	relationships := ToManyRelationships("pilots", tables)

	expected := []ToManyRelationship{
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "jets",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "licenses",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "languages",
			ForeignColumn:         "id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: true,
			JoinTable:   "pilot_languages",

			JoinLocalColumn:         "pilot_id",
			JoinLocalColumnNullable: false,
			JoinLocalColumnUnique:   false,

			JoinForeignColumn:         "language_id",
			JoinForeignColumnNullable: false,
			JoinForeignColumnUnique:   false,
		},
	}

	if len(relationships) != 3 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	for i, v := range relationships {
		if !reflect.DeepEqual(v, expected[i]) {
			t.Errorf("[%d] Mismatch between relationships:\n\nwant:%#v\n\ngot:%#v\n\n", i, expected[i], v)
		}
	}
}

func TestToManyRelationshipsNull(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name:    "pilots",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "name", Nullable: true}}},
		{
			Name:    "airports",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "size", Nullable: true}},
		},
		{
			Name:    "jets",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "pilot_id", Nullable: true}, {Name: "airport_id", Nullable: true}},
			FKeys: []ForeignKey{
				{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, ForeignColumnNullable: true},
				{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id", Nullable: true, ForeignColumnNullable: true},
			},
		},
		{
			Name:    "licenses",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "pilot_id", Nullable: true}},
			FKeys: []ForeignKey{
				{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, ForeignColumnNullable: true},
			},
		},
		{
			Name:    "hangars",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "name", Nullable: true}},
		},
		{
			Name:    "languages",
			Columns: []Column{{Name: "id", Nullable: true}, {Name: "language", Nullable: true}},
		},
		{
			Name:        "pilot_languages",
			IsJoinTable: true,
			Columns:     []Column{{Name: "pilot_id", Nullable: true}, {Name: "language_id", Nullable: true}},
			FKeys: []ForeignKey{
				{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, ForeignColumnNullable: true},
				{Name: "language_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id", Nullable: true, ForeignColumnNullable: true},
			},
		},
	}

	relationships := ToManyRelationships("pilots", tables)
	if len(relationships) != 3 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	expected := []ToManyRelationship{
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: true,
			Unique:   false,

			ForeignTable:          "jets",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: true,
			Unique:   false,

			ForeignTable:          "licenses",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: true,
			Unique:   false,

			ForeignTable:          "languages",
			ForeignColumn:         "id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   false,

			ToJoinTable: true,
			JoinTable:   "pilot_languages",

			JoinLocalColumn:         "pilot_id",
			JoinLocalColumnNullable: true,
			JoinLocalColumnUnique:   false,

			JoinForeignColumn:         "language_id",
			JoinForeignColumnNullable: true,
			JoinForeignColumnUnique:   false,
		},
	}

	for i, v := range relationships {
		if !reflect.DeepEqual(v, expected[i]) {
			t.Errorf("[%d] Mismatch between relationships:\n\nwant:%#v\n\ngot:%#v\n\n", i, expected[i], v)
		}
	}
}
