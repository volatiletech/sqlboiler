package boilingcore

import (
	"testing"

	"github.com/volatiletech/sqlboiler/v4/drivers"
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
		NameSingular        string

		LocalFn   string
		ForeignFn string
	}{
		{"jets", "airport_id", false, "airports", "id", true, "", "Jets", "Airport"},
		{"jets", "airport_id", true, "airports", "id", true, "", "Jet", "Airport"},

		{"jets", "holiday_id", false, "airports", "id", true, "", "HolidayJets", "Holiday"},
		{"jets", "holiday_id", true, "airports", "id", true, "", "HolidayJet", "Holiday"},

		{"jets", "holiday_airport_id", false, "airports", "id", true, "", "HolidayAirportJets", "HolidayAirport"},
		{"jets", "holiday_airport_id", true, "airports", "id", true, "", "HolidayAirportJet", "HolidayAirport"},

		{"jets", "jet_id", false, "jets", "id", true, "", "Jets", "Jet"},
		{"jets", "jet_id", true, "jets", "id", true, "", "Jet", "Jet"},
		{"jets", "plane_id", false, "jets", "id", true, "", "PlaneJets", "Plane"},
		{"jets", "plane_id", true, "jets", "id", true, "", "PlaneJet", "Plane"},

		{"videos", "user_id", false, "users", "id", true, "", "Videos", "User"},
		{"videos", "producer_id", false, "users", "id", true, "", "ProducerVideos", "Producer"},
		{"videos", "user_id", true, "users", "id", true, "", "Video", "User"},
		{"videos", "producer_id", true, "users", "id", true, "", "ProducerVideo", "Producer"},

		{"industries", "industry_id", false, "industries", "id", true, "", "Industries", "Industry"},
		{"industries", "parent_id", false, "industries", "id", true, "", "ParentIndustries", "Parent"},
		{"industries", "industry_id", true, "industries", "id", true, "", "Industry", "Industry"},
		{"industries", "parent_id", true, "industries", "id", true, "", "ParentIndustry", "Parent"},

		{"race_result_scratchings", "results_id", false, "race_results", "id", true, "", "ResultRaceResultScratchings", "Result"},
		{"race_result_scratchings", "results_id", false, "race_results", "id", true, "result", "RaceResultScratchings", "Result"},
	}

	for i, test := range tests {
		fk := drivers.ForeignKey{
			Table: test.Table, Column: test.Column, Unique: test.Unique,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn, ForeignColumnUnique: test.ForeignColumnUnique,
		}

		local, foreign := txtNameToOne(fk, test.NameSingular)
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
		LHSTable  string
		LHSColumn string

		RHSTable  string
		RHSColumn string

		LHSFn string
		RHSFn string
	}{
		{"pilots", "pilot_id", "languages", "language_id", "Pilots", "Languages"},
		{"pilots", "captain_id", "languages", "lingo_id", "CaptainPilots", "LingoLanguages"},

		{"pilots", "pilot_id", "pilots", "mentor_id", "Pilots", "MentorPilots"},
		{"pilots", "mentor_id", "pilots", "pilot_id", "MentorPilots", "Pilots"},
		{"pilots", "captain_id", "pilots", "mentor_id", "CaptainPilots", "MentorPilots"},

		{"videos", "video_id", "tags", "tag_id", "Videos", "Tags"},
		{"tags", "tag_id", "videos", "video_id", "Tags", "Videos"},
	}

	for i, test := range tests {
		lhsFk := drivers.ForeignKey{
			ForeignTable: test.LHSTable,
			Column:       test.LHSColumn,
		}
		rhsFk := drivers.ForeignKey{
			ForeignTable: test.RHSTable,
			Column:       test.RHSColumn,
		}

		lhs, rhs := txtNameToMany(lhsFk, rhsFk)
		if lhs != test.LHSFn {
			t.Error(i, "local wrong:", lhs, "want:", test.LHSFn)
		}
		if rhs != test.RHSFn {
			t.Error(i, "foreign wrong:", rhs, "want:", test.RHSFn)
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
