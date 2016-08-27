package bdb

import (
	"reflect"
	"testing"
)

func TestToManyRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name: "pilots",
			Columns: []Column{
				{Name: "id"},
				{Name: "name"},
			},
		},
		{
			Name: "airports",
			Columns: []Column{
				{Name: "id"},
				{Name: "size"},
			},
		},
		{
			Name: "jets",
			Columns: []Column{
				{Name: "id"},
				{Name: "pilot_id"},
				{Name: "airport_id"},
			},
			FKeys: []ForeignKey{
				{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
				{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id"},
			},
		},
		{
			Name: "licenses",
			Columns: []Column{
				{Name: "id"},
				{Name: "pilot_id"},
			},
			FKeys: []ForeignKey{
				{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
			},
		},
		{
			Name: "hangars",
			Columns: []Column{
				{Name: "id"},
				{Name: "name"},
			},
		},
		{
			Name: "languages",
			Columns: []Column{
				{Name: "id"},
				{Name: "language"},
			},
		},
		{
			Name:        "pilot_languages",
			IsJoinTable: true,
			Columns: []Column{
				{Name: "pilot_id"},
				{Name: "language_id"},
			},
			FKeys: []ForeignKey{
				{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id"},
				{Name: "language_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id"},
			},
		},
	}

	relationships := ToManyRelationships("pilots", tables)

	expected := []ToManyRelationship{
		{
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
			t.Errorf("[%d] Mismatch between relationships:\n\n%#v\n\n%#v\n\n", i, v, expected[i])
		}
	}
}

func TestToManyRelationshipsNull(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{
			Name: "pilots",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "name", Nullable: true, Unique: true},
			},
		},
		{
			Name: "airports",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "size", Nullable: true, Unique: true},
			},
		},
		{
			Name: "jets",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "pilot_id", Nullable: true, Unique: true},
				{Name: "airport_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "jets_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, Unique: true},
				{Name: "jets_airport_id_fk", Column: "airport_id", ForeignTable: "airports", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
		{
			Name: "licenses",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "pilot_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "licenses_pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
		{
			Name: "hangars",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "name", Nullable: true, Unique: true},
			},
		},
		{
			Name: "languages",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "language", Nullable: true, Unique: true},
			},
		},
		{
			Name:        "pilot_languages",
			IsJoinTable: true,
			Columns: []Column{
				{Name: "pilot_id", Nullable: true, Unique: true},
				{Name: "language_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "pilot_id_fk", Column: "pilot_id", ForeignTable: "pilots", ForeignColumn: "id", Nullable: true, Unique: true},
				{Name: "language_id_fk", Column: "language_id", ForeignTable: "languages", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
	}

	relationships := ToManyRelationships("pilots", tables)
	if len(relationships) != 3 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	expected := []ToManyRelationship{
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "jets",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "licenses",
			ForeignColumn:         "pilot_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: false,
		},
		{
			Table:    "pilots",
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "languages",
			ForeignColumn:         "id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: true,
			JoinTable:   "pilot_languages",

			JoinLocalColumn:         "pilot_id",
			JoinLocalColumnNullable: true,
			JoinLocalColumnUnique:   true,

			JoinForeignColumn:         "language_id",
			JoinForeignColumnNullable: true,
			JoinForeignColumnUnique:   true,
		},
	}

	for i, v := range relationships {
		if !reflect.DeepEqual(v, expected[i]) {
			t.Errorf("[%d] Mismatch between relationships null:\n\n%#v\n\n%#v\n\n", i, v, expected[i])
		}
	}
}
