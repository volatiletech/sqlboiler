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
		Aliases             *Aliases

		LocalFn   string
		ForeignFn string
	}{
		{"jets", "airport_id", false, "airports", "id", true, &Aliases{}, "Jets", "Airport"},
		{"jets", "airport_id", true, "airports", "id", true, &Aliases{}, "Jet", "Airport"},

		{"jets", "holiday_id", false, "airports", "id", true, &Aliases{}, "HolidayJets", "Holiday"},
		{"jets", "holiday_id", true, "airports", "id", true, &Aliases{}, "HolidayJet", "Holiday"},

		{"jets", "holiday_airport_id", false, "airports", "id", true, &Aliases{}, "HolidayAirportJets", "HolidayAirport"},
		{"jets", "holiday_airport_id", true, "airports", "id", true, &Aliases{}, "HolidayAirportJet", "HolidayAirport"},

		{"jets", "jet_id", false, "jets", "id", true, &Aliases{}, "Jets", "Jet"},
		{"jets", "jet_id", true, "jets", "id", true, &Aliases{}, "Jet", "Jet"},
		{"jets", "plane_id", false, "jets", "id", true, &Aliases{}, "PlaneJets", "Plane"},
		{"jets", "plane_id", true, "jets", "id", true, &Aliases{}, "PlaneJet", "Plane"},

		{"videos", "user_id", false, "users", "id", true, &Aliases{}, "Videos", "User"},
		{"videos", "producer_id", false, "users", "id", true, &Aliases{}, "ProducerVideos", "Producer"},
		{"videos", "user_id", true, "users", "id", true, &Aliases{}, "Video", "User"},
		{"videos", "producer_id", true, "users", "id", true, &Aliases{}, "ProducerVideo", "Producer"},

		{"videos", "user", false, "users", "id", true, &Aliases{}, "Videos", "VideoUser"},
		{"videos", "created_by", false, "users", "id", true, &Aliases{}, "CreatedByVideos", "CreatedByUser"},
		{"videos", "director", false, "users", "id", true, &Aliases{}, "DirectorVideos", "DirectorUser"},
		{"videos", "user", true, "users", "id", true, &Aliases{}, "Video", "VideoUser"},
		{"videos", "created_by", true, "users", "id", &Aliases{}, true, "CreatedByVideo", "CreatedByUser"},
		{"videos", "director", true, "users", "id", true, &Aliases{}, "DirectorVideo", "DirectorUser"},

		{"industries", "industry_id", false, "industries", "id", true, &Aliases{}, "Industries", "Industry"},
		{"industries", "parent_id", false, "industries", "id", true, &Aliases{}, "ParentIndustries", "Parent"},
		{"industries", "industry_id", true, "industries", "id", true, &Aliases{}, "Industry", "Industry"},
		{"industries", "parent_id", true, "industries", "id", true, &Aliases{}, "ParentIndustry", "Parent"},

		{"race_result_scratchings", "results_id", false, "race_results", "id", true, &Aliases{}, "ResultRaceResultScratchings", "Result"},
		{"race_result_scratchings", "results_id", false, "race_results", "id", true,
			&Aliases{Tables: map[string]TableAlias{"race_results": {DownSingular: "result"}}}, "RaceResultScratchings", "Result"},
	}

	for i, test := range tests {
		fk := drivers.ForeignKey{
			Table: test.Table, Column: test.Column, Unique: test.Unique,
			ForeignTable: test.ForeignTable, ForeignColumn: test.ForeignColumn, ForeignColumnUnique: test.ForeignColumnUnique,
		}
		a := test.Aliases

		local, foreign := txtNameToOne(fk, a)
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
