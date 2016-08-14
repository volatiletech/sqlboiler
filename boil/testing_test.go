package boil

import (
	"testing"
	"time"

	"gopkg.in/nullbio/null.v4"
)

func TestRandomizeStruct(t *testing.T) {
	t.Parallel()

	var testStruct = struct {
		Int       int
		Int64     int64
		Float64   float64
		Bool      bool
		Time      time.Time
		String    string
		ByteSlice []byte
		Interval  string

		Ignore int

		NullInt      null.Int
		NullFloat64  null.Float64
		NullBool     null.Bool
		NullString   null.String
		NullTime     null.Time
		NullInterval null.String
	}{}

	fieldTypes := map[string]string{
		"Int":          "integer",
		"Int64":        "bigint",
		"Float64":      "decimal",
		"Bool":         "boolean",
		"Time":         "date",
		"String":       "character varying",
		"ByteSlice":    "bytea",
		"Interval":     "interval",
		"Ignore":       "integer",
		"NullInt":      "integer",
		"NullFloat64":  "numeric",
		"NullBool":     "boolean",
		"NullString":   "character",
		"NullTime":     "time",
		"NullInterval": "interval",
	}

	err := RandomizeStruct(&testStruct, fieldTypes, true, "Ignore")
	if err != nil {
		t.Fatal(err)
	}

	if testStruct.Ignore != 0 {
		t.Error("blacklisted value was filled in:", testStruct.Ignore)
	}

	if testStruct.Int == 0 &&
		testStruct.Int64 == 0 &&
		testStruct.Float64 == 0 &&
		testStruct.Bool == false &&
		testStruct.Time.IsZero() &&
		testStruct.String == "" &&
		testStruct.Interval == "" &&
		testStruct.ByteSlice == nil {
		t.Errorf("the regular values are not being randomized: %#v", testStruct)
	}

	if testStruct.NullInt.Valid == false &&
		testStruct.NullFloat64.Valid == false &&
		testStruct.NullBool.Valid == false &&
		testStruct.NullString.Valid == false &&
		testStruct.NullInterval.Valid == false &&
		testStruct.NullTime.Valid == false {
		t.Errorf("the null values are not being randomized: %#v", testStruct)
	}
}
