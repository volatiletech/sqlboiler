package boilingcore

import (
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/drivers"
)

func TestAliasesTables(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
			Columns: []drivers.Column{
				{Name: "id"},
				{Name: "name"},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect := TableAlias{
			UpPlural:     "Videos",
			UpSingular:   "Video",
			DownPlural:   "videos",
			DownSingular: "video",
			Columns: map[string]string{
				"id":   "ID",
				"name": "Name",
			},
		}

		a := Aliases{}
		FillAliases(&a, tables)

		if got := a.Tables["videos"]; !reflect.DeepEqual(expect, got) {
			t.Errorf("it should fill in the blanks: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect := TableAlias{
			UpPlural:     "NotVideos",
			UpSingular:   "NotVideo",
			DownPlural:   "notVideos",
			DownSingular: "notVideo",
			Columns: map[string]string{
				"id":   "NotID",
				"name": "NotName",
			},
		}

		a := Aliases{}
		a.Tables = map[string]TableAlias{"videos": expect}
		FillAliases(&a, tables)

		if !reflect.DeepEqual(expect, a.Tables["videos"]) {
			t.Error("it should not alter things that were specified by user")
		}
	})
}

func TestAliasesRelationships(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
			Columns: []drivers.Column{
				{Name: "id"},
				{Name: "name"},
			},
			FKeys: []drivers.ForeignKey{
				{
					Name:          "fkey_1",
					Table:         "videos",
					Column:        "user_id",
					ForeignTable:  "users",
					ForeignColumn: "id",
				},
				{
					Name:          "fkey_2",
					Table:         "videos",
					Column:        "publisher_id",
					ForeignTable:  "users",
					ForeignColumn: "id",
				},
				{
					Name:          "fkey_3",
					Table:         "videos",
					Column:        "one_id",
					Unique:        true,
					ForeignTable:  "ones",
					ForeignColumn: "one",
				},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "User",
			Foreign: "Videos",
		}
		expect2 := RelationshipAlias{
			Local:   "Publisher",
			Foreign: "PublisherVideos",
		}
		expect3 := RelationshipAlias{
			Local:   "One",
			Foreign: "Video",
		}

		a := Aliases{}
		FillAliases(&a, tables)

		if got := a.Relationships["fkey_1"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fkey_2"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fkey_3"]; !reflect.DeepEqual(expect3, got) {
			t.Errorf("bad values: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "TheUser",
			Foreign: "Videos",
		}
		expect2 := RelationshipAlias{
			Local:   "Publisher",
			Foreign: "PublishedVideos",
		}
		expect3 := RelationshipAlias{
			Local:   "TheOne",
			Foreign: "AwesomeOneVideo",
		}

		a := Aliases{
			Relationships: map[string]RelationshipAlias{
				"fkey_1": {Local: "TheUser"},
				"fkey_2": {Foreign: "PublishedVideos"},
				"fkey_3": {Local: "TheOne", Foreign: "AwesomeOneVideo"},
			},
		}
		FillAliases(&a, tables)

		if got := a.Relationships["fkey_1"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fkey_2"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fkey_3"]; !reflect.DeepEqual(expect3, got) {
			t.Errorf("bad values: %#v", got)
		}
	})
}

func TestAliasesRelationshipsJoinTable(t *testing.T) {
	t.Parallel()

	tables := []drivers.Table{
		{
			Name: "videos",
			ToManyRelationships: []drivers.ToManyRelationship{
				{
					Table:         "videos",
					ForeignTable:  "tags",
					Column:        "id",
					ForeignColumn: "id",

					ToJoinTable: true,
					JoinTable:   "video_tags",

					JoinLocalFKeyName:   "fk_video_id",
					JoinLocalColumn:     "video_id",
					JoinForeignFKeyName: "fk_tag_id",
					JoinForeignColumn:   "tag_id",
				},
			},
		},
	}

	t.Run("Automatic", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "Tags",
			Foreign: "Videos",
		}
		expect2 := RelationshipAlias{
			Local:   "Videos",
			Foreign: "Tags",
		}

		a := Aliases{}
		FillAliases(&a, tables)

		if got := a.Relationships["fk_video_id"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fk_tag_id"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
	})

	t.Run("UserOverride", func(t *testing.T) {
		expect1 := RelationshipAlias{
			Local:   "NotTags",
			Foreign: "NotVideos",
		}
		expect2 := RelationshipAlias{
			Local:   "NotVideos",
			Foreign: "NotTags",
		}

		a := Aliases{
			Relationships: map[string]RelationshipAlias{
				"fk_video_id": {Local: "NotTags", Foreign: "NotVideos"},
			},
		}
		FillAliases(&a, tables)

		if got := a.Relationships["fk_video_id"]; !reflect.DeepEqual(expect1, got) {
			t.Errorf("bad values: %#v", got)
		}
		if got := a.Relationships["fk_tag_id"]; !reflect.DeepEqual(expect2, got) {
			t.Errorf("bad values: %#v", got)
		}
	})
}
