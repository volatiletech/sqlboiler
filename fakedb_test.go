package main

import (
	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/strmangle"
)

type fakeDB int

func (fakeDB) TableNames() ([]string, error) {
	return []string{"users", "videos", "contests", "notifications", "users_videos_tags"}, nil
}
func (fakeDB) Columns(tableName string) ([]bdb.Column, error) {
	return map[string][]bdb.Column{
		"users":    {{Name: "id", Type: "int32"}},
		"contests": {{Name: "id", Type: "int32", Nullable: true}},
		"videos": {
			{Name: "id", Type: "int32"},
			{Name: "user_id", Type: "int32", Nullable: true, Unique: true},
			{Name: "contest_id", Type: "int32"},
		},
		"notifications": {
			{Name: "user_id", Type: "int32"},
			{Name: "source_id", Type: "int32", Nullable: true},
		},
		"users_videos_tags": {
			{Name: "user_id", Type: "int32"},
			{Name: "video_id", Type: "int32"},
		},
	}[tableName], nil
}
func (fakeDB) ForeignKeyInfo(tableName string) ([]bdb.ForeignKey, error) {
	return map[string][]bdb.ForeignKey{
		"videos": {
			{Name: "videos_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
			{Name: "videos_contest_id_fk", Column: "contest_id", ForeignTable: "contests", ForeignColumn: "id"},
		},
		"notifications": {
			{Name: "notifications_user_id_fk", Column: "user_id", ForeignTable: "users", ForeignColumn: "id"},
			{Name: "notifications_source_id_fk", Column: "source_id", ForeignTable: "users", ForeignColumn: "id"},
		},
		"users_videos_tags": {
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
		"users_videos_tags": {
			Name:    "user_video_id_pkey",
			Columns: []string{"user_id", "video_id"},
		},
	}[tableName], nil
}
func (fakeDB) UseLastInsertID() bool { return false }
func (fakeDB) Open() error           { return nil }
func (fakeDB) Close()                {}
