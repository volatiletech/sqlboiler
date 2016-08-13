package main

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/vattle/sqlboiler/bdb"
)

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
	expect.ForeignTable.NamePluralGo = "Users"
	expect.ForeignTable.ColumnName = "id"
	expect.ForeignTable.ColumnNameGo = "ID"

	expect.Function.PackageName = "models"
	expect.Function.Name = "User"
	expect.Function.Varname = "user"
	expect.Function.Receiver = "v"
	expect.Function.ReverseInserts = false

	expect.Function.LocalAssignment = "UserID.Int32"
	expect.Function.ForeignAssignment = "ID"

	if !reflect.DeepEqual(expect, texts) {
		t.Errorf("Want:\n%s\nGot:\n%s\n", spew.Sdump(expect), spew.Sdump(texts))
	}
}

func TestTextsFromOneToOneRelationship(t *testing.T) {
	t.Parallel()

	tables, err := bdb.Tables(fakeDB(0))
	if err != nil {
		t.Fatal(err)
	}

	users := bdb.GetTable(tables, "users")
	texts := textsFromOneToOneRelationship("models", tables, users, users.ToManyRelationships[0])
	expect := RelationshipToOneTexts{}

	expect.ForeignKey = bdb.ForeignKey{
		Name:     "none",
		Column:   "id",
		Nullable: false,
		Unique:   false,

		ForeignTable:          "videos",
		ForeignColumn:         "user_id",
		ForeignColumnNullable: true,
		ForeignColumnUnique:   true,
	}

	expect.LocalTable.NameGo = "User"
	expect.LocalTable.ColumnNameGo = "ID"

	expect.ForeignTable.Name = "videos"
	expect.ForeignTable.NameGo = "Video"
	expect.ForeignTable.NamePluralGo = "Videos"
	expect.ForeignTable.ColumnName = "user_id"
	expect.ForeignTable.ColumnNameGo = "UserID"

	expect.Function.PackageName = "models"
	expect.Function.Name = "Video"
	expect.Function.Varname = "video"
	expect.Function.Receiver = "u"
	expect.Function.ReverseInserts = true

	expect.Function.LocalAssignment = "ID"
	expect.Function.ForeignAssignment = "UserID.Int32"

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
