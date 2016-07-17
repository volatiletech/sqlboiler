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
		"users_videos_tags": []bdb.ForeignKey{
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
func (fakeDB) PrimaryKeyInfo(tableName string) (*bdb.PrimaryKey, error) {
	return map[string]*bdb.PrimaryKey{
		"users_videos_tags": &bdb.PrimaryKey{
			Name:    "user_video_id_pkey",
			Columns: []string{"user_id", "video_id"},
		},
	}[tableName], nil
}
func (fakeDB) Open() error { return nil }
func (fakeDB) Close()      {}

func TestTextsFromForeignKey(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(fakeDB(0))
	if err != nil {
		t.Fatal(err)
	}

	videos := bdb.GetTable(tables, "videos")
	texts := textsFromForeignKey("models", tables, videos, videos.FKeys[0])
	expect := RelationshipToOneTexts{}

	expect.ForeignKey = videos.FKeys[0]

	expect.LocalTable.NameGo = "Video"
	expect.LocalTable.ColumnNameGo = "UserID"

	expect.ForeignTable.Name = "users"
	expect.ForeignTable.NameGo = "User"
	expect.ForeignTable.ColumnName = "id"
	expect.ForeignTable.ColumnNameGo = "ID"

	expect.Function.PackageName = "models"
	expect.Function.Name = "User"
	expect.Function.Varname = "user"
	expect.Function.Receiver = "v"

	expect.Function.LocalAssignment = "UserID.Int32"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTextsFromRelationship(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(fakeDB(0))
	if err != nil {
		t.Fatal(err)
	}

	users := bdb.GetTable(tables, "users")
	texts := textsFromRelationship(tables, users, users.ToManyRelationships[0])
	expect := RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "User"
	expect.LocalTable.NameSingular = "user"

	expect.ForeignTable.NameGo = "Video"
	expect.ForeignTable.NameSingular = "video"
	expect.ForeignTable.NamePluralGo = "Videos"
	expect.ForeignTable.NameHumanReadable = "videos"
	expect.ForeignTable.Slice = "VideoSlice"

	expect.Function.Name = "Videos"
	expect.Function.Receiver = "u"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "UserID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, users, users.ToManyRelationships[1])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "User"
	expect.LocalTable.NameSingular = "user"

	expect.ForeignTable.NameGo = "Notification"
	expect.ForeignTable.NameSingular = "notification"
	expect.ForeignTable.NamePluralGo = "Notifications"
	expect.ForeignTable.NameHumanReadable = "notifications"
	expect.ForeignTable.Slice = "NotificationSlice"

	expect.Function.Name = "Notifications"
	expect.Function.Receiver = "u"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "UserID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, users, users.ToManyRelationships[2])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "User"
	expect.LocalTable.NameSingular = "user"

	expect.ForeignTable.NameGo = "Notification"
	expect.ForeignTable.NameSingular = "notification"
	expect.ForeignTable.NamePluralGo = "Notifications"
	expect.ForeignTable.NameHumanReadable = "notifications"
	expect.ForeignTable.Slice = "NotificationSlice"

	expect.Function.Name = "SourceNotifications"
	expect.Function.Receiver = "u"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "SourceID.Int32"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}

	texts = textsFromRelationship(tables, users, users.ToManyRelationships[3])
	expect = RelationshipToManyTexts{}
	expect.LocalTable.NameGo = "User"
	expect.LocalTable.NameSingular = "user"

	expect.ForeignTable.NameGo = "Video"
	expect.ForeignTable.NameSingular = "video"
	expect.ForeignTable.NamePluralGo = "Videos"
	expect.ForeignTable.NameHumanReadable = "videos"
	expect.ForeignTable.Slice = "VideoSlice"

	expect.Function.Name = "Videos"
	expect.Function.Receiver = "u"
	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}
