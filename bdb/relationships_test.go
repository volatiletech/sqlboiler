package bdb

import (
	"reflect"
	"testing"
)

func TestToManyRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		{Name: "users", Columns: []Column{{Name: "id"}}},
		{Name: "contests", Columns: []Column{{Name: "id"}}},
		{
			Name: "videos",
			Columns: []Column{
				{Name: "id"},
				{Name: "user_id"},
				{Name: "contest_id"},
			},
			FKeys: []ForeignKey{
				{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
				{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id"},
			},
		},
		{
			Name: "notifications",
			Columns: []Column{
				{Name: "user_id"},
				{Name: "source_id"},
			},
			FKeys: []ForeignKey{
				{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
				{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id"},
			},
		},
		{
			Name:        "users_video_tags",
			IsJoinTable: true,
			Columns: []Column{
				{Name: "user_id"},
				{Name: "video_id"},
			},
			FKeys: []ForeignKey{
				{Name: "user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
				{Name: "video_id_fk", Column: "video_id", ForeignTable: "videos", ForeignColumn: "id"},
			},
		},
	}

	relationships := ToManyRelationships("users", tables)

	expected := []ToManyRelationship{
		{
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "videos",
			ForeignColumn:         "user_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "notifications",
			ForeignColumn:         "user_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "notifications",
			ForeignColumn:         "source_id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: false,
			Unique:   false,

			ForeignTable:          "videos",
			ForeignColumn:         "id",
			ForeignColumnNullable: false,
			ForeignColumnUnique:   false,

			ToJoinTable: true,
			JoinTable:   "users_video_tags",

			JoinLocalColumn:         "user_id",
			JoinLocalColumnNullable: false,
			JoinLocalColumnUnique:   false,

			JoinForeignColumn:         "video_id",
			JoinForeignColumnNullable: false,
			JoinForeignColumnUnique:   false,
		},
	}

	if len(relationships) != 4 {
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
		{Name: "users", Columns: []Column{{Name: "id", Nullable: true, Unique: true}}},
		{Name: "contests", Columns: []Column{{Name: "id", Nullable: true, Unique: true}}},
		{
			Name: "videos",
			Columns: []Column{
				{Name: "id", Nullable: true, Unique: true},
				{Name: "user_id", Nullable: true, Unique: true},
				{Name: "contest_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true, Unique: true},
				{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
		{
			Name: "notifications",
			Columns: []Column{
				{Name: "user_id", Nullable: true, Unique: true},
				{Name: "source_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true, Unique: true},
				{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
		{
			Name:        "users_video_tags",
			IsJoinTable: true,
			Columns: []Column{
				{Name: "user_id", Nullable: true, Unique: true},
				{Name: "video_id", Nullable: true, Unique: true},
			},
			FKeys: []ForeignKey{
				{Name: "user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true, Unique: true},
				{Name: "video_id_fk", Column: "video_id", ForeignTable: "videos", ForeignColumn: "id", Nullable: true, Unique: true},
			},
		},
	}

	relationships := ToManyRelationships("users", tables)
	if len(relationships) != 4 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	expected := []ToManyRelationship{
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "videos",
			ForeignColumn:         "user_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "notifications",
			ForeignColumn:         "user_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "notifications",
			ForeignColumn:         "source_id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: false,
		},
		{
			Column:   "id",
			Nullable: true,
			Unique:   true,

			ForeignTable:          "videos",
			ForeignColumn:         "id",
			ForeignColumnNullable: true,
			ForeignColumnUnique:   true,

			ToJoinTable: true,
			JoinTable:   "users_video_tags",

			JoinLocalColumn:         "user_id",
			JoinLocalColumnNullable: true,
			JoinLocalColumnUnique:   true,

			JoinForeignColumn:         "video_id",
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
