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
	}

	relationships := ToManyRelationships("users", tables)
	r := relationships[0]
	if r.Name != "Videos" {
		t.Error("wrong name:", r.Name)
	}
	if r.ForeignTable != "videos" {
		t.Error("wrong foreign table:", r.ForeignTable)
	}
	if r.ForeignColumn != "user_id" {
		t.Error("wrong foreign column:", r.ForeignColumn)
	}
}
