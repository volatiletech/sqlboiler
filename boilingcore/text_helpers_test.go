package boilingcore

import (
	"testing"

	"github.com/volatiletech/sqlboiler/drivers"
)

func TestTxtNameToOne(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Table               string
		Column              string
		Unique              bool
		ForeignTable        string
		ForeignColumn       string
		ForeignColumnUnique bool

		LocalFn   string
		ForeignFn string
	}{
		{"jets", "airport_id", false, "airports", "id", true, "Jets", "Airport"},
		{"jets", "airport_id", true, "airports", "id", true, "Jet", "Airport"},

		{"jets", "holiday_id", false, "airports", "id", true, "HolidayJets", "Holiday"},
		{"jets", "holiday_id", true, "airports", "id", true, "HolidayJet", "Holiday"},

		{"jets", "holiday_airport_id", false, "airports", "id", true, "HolidayAirportJets", "HolidayAirport"},
		{"jets", "holiday_airport_id", true, "airports", "id", true, "HolidayAirportJet", "HolidayAirport"},

		{"jets", "jet_id", false, "jets", "id", true, "Jets", "Jet"},
		{"jets", "jet_id", true, "jets", "id", true, "Jet", "Jet"},
		{"jets", "plane_id", false, "jets", "id", true, "PlaneJets", "Plane"},
		{"jets", "plane_id", true, "jets", "id", true, "PlaneJet", "Plane"},

		{"race_result_scratchings", "results_id", false, "race_results", "id", true, "ResultRaceResultScratchings", "Result"},
	}

	for i, test := range tests {
		fk := drivers.ForeignKey{
			Table: test.Table, Column: test.Column, Unique: test.Unique,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn, ForeignColumnUnique: test.ForeignColumnUnique,
		}

		local, foreign := txtNameToOne(fk)
		if local != test.LocalFn {
			t.Error(i, "local wrong:", local, "want:", test.LocalFn)
		}
		if foreign != test.ForeignFn {
			t.Error(i, "foreign wrong:", foreign, "want:", test.ForeignFn)
		}
	}
}

func TestTxtNameToMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Table  string
		Column string

		ForeignTable  string
		ForeignColumn string

		ToJoinTable       bool
		JoinLocalColumn   string
		JoinForeignColumn string

		LocalFn   string
		ForeignFn string
	}{
		{"pilots", "id", "languages", "id", true, "pilot_id", "language_id", "Pilots", "Languages"},
		{"pilots", "id", "languages", "id", true, "captain_id", "lingo_id", "CaptainPilots", "LingoLanguages"},

		{"pilots", "id", "pilots", "id", true, "pilot_id", "mentor_id", "Pilots", "MentorPilots"},
		{"pilots", "id", "pilots", "id", true, "mentor_id", "pilot_id", "MentorPilots", "Pilots"},
		{"pilots", "id", "pilots", "id", true, "captain_id", "mentor_id", "CaptainPilots", "MentorPilots"},

		{"videos", "id", "tags", "id", true, "video_id", "tag_id", "Videos", "Tags"},
		{"tags", "id", "videos", "id", true, "tag_id", "video_id", "Tags", "Videos"},
	}

	for i, test := range tests {
		fk := drivers.ToManyRelationship{
			Table: test.Table, Column: test.Column,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn,
			ToJoinTable:     test.ToJoinTable,
			JoinLocalColumn: test.JoinLocalColumn, JoinForeignColumn: test.JoinForeignColumn,
		}

		local, foreign := txtNameToMany(fk)
		if local != test.LocalFn {
			t.Error(i, "local wrong:", local, "want:", test.LocalFn)
		}
		if foreign != test.ForeignFn {
			t.Error(i, "foreign wrong:", foreign, "want:", test.ForeignFn)
		}
	}
}

func TestTrimSuffixes(t *testing.T) {
	t.Parallel()

	for _, s := range identifierSuffixes {
		a := "hello" + s

		if z := trimSuffixes(a); z != "hello" {
			t.Errorf("got %s", z)
		}
	}
}
