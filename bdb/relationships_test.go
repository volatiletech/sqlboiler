package bdb

import "testing"

func TestToManyRelationships(t *testing.T) {
	t.Parallel()

	tables := []Table{
		Table{
			Name: "videos",
			FKeys: []ForeignKey{
				{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
				{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id"},
			},
		},
		Table{
			Name: "notifications",
			FKeys: []ForeignKey{
				{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
				{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id"},
			},
		},
		Table{
			Name:        "users_video_tags",
			IsJoinTable: true,
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
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[1]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[2]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.ForeignTable != "notifications" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "source_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
	if r.ToJoinTable {
		t.Error("not a join table")
	}

	r = relationships[3]
	if r.Column != "id" {
		t.Error("wrong local column:", r.Column)
	}
	if r.ForeignColumn != "id" {
		t.Error("wrong foreign column:", r.Column)
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
	if r.JoinForeignColumn != "video_id" {
		t.Error("wrong foreign join column:", r.JoinForeignColumn)
	}
	if !r.ToJoinTable {
		t.Error("expected a join table")
	}
}
