package main

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/nullbio/sqlboiler/bdb"
	"github.com/nullbio/sqlboiler/strmangle"
)

type fakeDB int

func (fakeDB) TableNames() ([]string, error) {
	return []string{"users", "videos", "contests", "notifications", "users_videos_tags"}, nil
}
func (fakeDB) Columns(tableName string) ([]bdb.Column, error) {
	return map[string][]bdb.Column{
		"users":    []bdb.Column{{Name: "id", Type: "int32"}},
		"contests": []bdb.Column{{Name: "id", Type: "int32", Nullable: true}},
		"videos": []bdb.Column{
			{Name: "id", Type: "int32"},
			{Name: "user_id", Type: "int32", Nullable: true},
			{Name: "contest_id", Type: "int32"},
		},
		"notifications": []bdb.Column{
			{Name: "user_id", Type: "int32"},
			{Name: "source_id", Type: "int32", Nullable: true},
		},
		"users_videos_tags": []bdb.Column{
			{Name: "user_id", Type: "int32"},
			{Name: "video_id", Type: "int32"},
		},
	}[tableName], nil
}
func (fakeDB) ForeignKeyInfo(tableName string) ([]bdb.ForeignKey, error) {
	return map[string][]bdb.ForeignKey{
		"videos": []bdb.ForeignKey{
			{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
			{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id"},
		},
		"notifications": []bdb.ForeignKey{
			{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
			{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id"},
		},
		"users_video_tags": []bdb.ForeignKey{
			{Name: "user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
			{Name: "video_id_fk", Column: "video_id", ForeignTable: "videos", ForeignColumn: "id"},
		},
	}[tableName], nil
}
func (fakeDB) TranslateColumnType(c bdb.Column) bdb.Column {
	if c.Nullable {
		c.Type = "null." + strmangle.TitleCase(c.Type)
	}
	return c
}
func (fakeDB) PrimaryKeyInfo(tableName string) (*bdb.PrimaryKey, error) { return nil, nil }
func (fakeDB) Open() error                                              { return nil }
func (fakeDB) Close()                                                   {}

func TestCreateTextsFromRelationship(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(fakeDB(0))
	if err != nil {
		t.Fatal(err)
	}

	users := bdb.GetTable(tables, "users")
	texts := createTextsFromRelationship(tables, users, users.ToManyRelationships[0])

	expect := RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "User"
	expect.LocalTable.NameSingular = "user"

	expect.ForeignTable.NameGo = "Video"
	expect.ForeignTable.NameSingular = "video"
	expect.ForeignTable.NamePluralGo = "Videos"
	expect.ForeignTable.NameHumanReadable = "videos"
	expect.ForeignTable.Slice = "videoSlice"

	expect.Function.Name = "Videos"
	expect.Function.Receiver = "u"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "UserID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}
