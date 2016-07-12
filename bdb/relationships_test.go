package bdb

import "testing"

func TestToManyRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		Table{Name: "users", Columns: []Column{{Name: "id"}}},
		Table{Name: "contests", Columns: []Column{{Name: "id"}}},
		Table{
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
		Table{
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
		Table{
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
	if len(relationships) != 4 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	r := relationships[0]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.Nullable {
		t.Error("should not be nullable")
	}
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ForeignColumnNullable {
		t.Error("should not be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[1]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.Nullable {
		t.Error("should not be nullable")
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ForeignColumnNullable {
		t.Error("should not be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[2]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.Nullable {
		t.Error("should not be nullable")
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "source_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ForeignColumnNullable {
		t.Error("should not be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[3]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.Nullable {
		t.Error("should not be nullable")
	}
	if r.ForeignColumn != "id" {
		t.Error("wrong foreign column:", r.Column)
	}
	if r.ForeignColumnNullable {
		t.Error("should not be nullable")
	}
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.JoinTable != "users_video_tags" {
		t.Error("wrong join table:", r.ForeignTable)
	}
	if r.JoinLocalColumn != "user_id" {
		t.Error("wrong local join column:", r.JoinLocalColumn)
	}
	if r.JoinLocalColumnNullable {
		t.Error("should not be nullable")
	}
	if r.JoinForeignColumn != "video_id" {
		t.Error("wrong foreign join column:", r.JoinForeignColumn)
	}
	if r.JoinForeignColumnNullable {
		t.Error("should not be nullable")
	}
	if !r.ToJoinTable {
		t.Error("expected a join table")
	}
}

func TestToManyRelationshipsNull(t *testing.T) {
	t.Parallel()

	tables := []Table{
		Table{Name: "users", Columns: []Column{{Name: "id", Nullable: true}}},
		Table{Name: "contests", Columns: []Column{{Name: "id", Nullable: true}}},
		Table{
			Name: "videos",
			Columns: []Column{
				{Name: "id", Nullable: true},
				{Name: "user_id", Nullable: true},
				{Name: "contest_id", Nullable: true},
			},
			FKeys: []ForeignKey{
				{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true},
				{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id", Nullable: true},
			},
		},
		Table{
			Name: "notifications",
			Columns: []Column{
				{Name: "user_id", Nullable: true},
				{Name: "source_id", Nullable: true},
			},
			FKeys: []ForeignKey{
				{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true},
				{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true},
			},
		},
		Table{
			Name:        "users_video_tags",
			IsJoinTable: true,
			Columns: []Column{
				{Name: "user_id", Nullable: true},
				{Name: "video_id", Nullable: true},
			},
			FKeys: []ForeignKey{
				{Name: "user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id", Nullable: true},
				{Name: "video_id_fk", Column: "video_id", ForeignTable: "videos", ForeignColumn: "id", Nullable: true},
			},
		},
	}

	relationships := ToManyRelationships("users", tables)
	if len(relationships) != 4 {
		t.Error("wrong # of relationships:", len(relationships))
	}

	r := relationships[0]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if !r.Nullable {
		t.Error("should be nullable")
	}
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if !r.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[1]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if !r.Nullable {
		t.Error("should be nullable")
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if !r.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[2]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if !r.Nullable {
		t.Error("should be nullable")
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "source_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if !r.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[3]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if !r.Nullable {
		t.Error("should be nullable")
	}
	if r.ForeignColumn != "id" {
		t.Error("wrong foreign column:", r.Column)
	}
	if !r.ForeignColumnNullable {
		t.Error("should be nullable")
	}
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.JoinTable != "users_video_tags" {
		t.Error("wrong join table:", r.ForeignTable)
	}
	if r.JoinLocalColumn != "user_id" {
		t.Error("wrong local join column:", r.JoinLocalColumn)
	}
	if !r.JoinLocalColumnNullable {
		t.Error("should be nullable")
	}
	if r.JoinForeignColumn != "video_id" {
		t.Error("wrong foreign join column:", r.JoinForeignColumn)
	}
	if !r.JoinForeignColumnNullable {
		t.Error("should be nullable")
	}
	if !r.ToJoinTable {
		t.Error("expected a join table")
	}
}
