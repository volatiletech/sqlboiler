package boil

import (
	"testing"
	"time"

	"gopkg.in/nullbio/null.v4"
)

func TestIsZeroValue(t *testing.T) {
	t.Parallel()

	o := struct {
		A []byte
		B time.Time
		C null.Time
		D null.Int64
		E int64
	}{}

	if errs := IsZeroValue(o, true, "A", "B", "C", "D", "E"); errs != nil {
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}

	colNames := []string{"A", "B", "C", "D", "E"}
	for _, c := range colNames {
		if err := IsZeroValue(o, true, c); err != nil {
			t.Errorf("Expected %s to be zero value: %s", c, err[0])
		}
	}

	o.A = []byte("asdf")
	o.B = time.Now()
	o.C = null.NewTime(time.Now(), false)
	o.D = null.NewInt64(2, false)
	o.E = 5

	if errs := IsZeroValue(o, false, "A", "B", "C", "D", "E"); errs != nil {
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}

	for _, c := range colNames {
		if err := IsZeroValue(o, false, c); err != nil {
			t.Errorf("Expected %s to be non-zero value: %s", c, err[0])
		}
	}
}

func TestIsValueMatch(t *testing.T) {
	t.Parallel()

	var errs []error
	var values []interface{}

	o := struct {
		A []byte
		B time.Time
		C null.Time
		D null.Int64
		E int64
	}{}

	values = []interface{}{
		[]byte(nil),
		time.Time{},
		null.Time{},
		null.Int64{},
		int64(0),
	}

	cols := []string{"A", "B", "C", "D", "E"}
	errs = IsValueMatch(o, cols, values)
	if errs != nil {
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}

	values = []interface{}{
		[]byte("hi"),
		time.Date(2007, 11, 2, 1, 1, 1, 1, time.UTC),
		null.NewTime(time.Date(2007, 11, 2, 1, 1, 1, 1, time.UTC), true),
		null.NewInt64(5, false),
		int64(6),
	}

	errs = IsValueMatch(o, cols, values)
	// Expect 6 errors
	// 5 for each column and an additional 1 for the invalid Valid field match
	if len(errs) != 6 {
		t.Errorf("Expected 6 errors, got: %d", len(errs))
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}

	o.A = []byte("hi")
	o.B = time.Date(2007, 11, 2, 1, 1, 1, 1, time.UTC)
	o.C = null.NewTime(time.Date(2007, 11, 2, 1, 1, 1, 1, time.UTC), true)
	o.D = null.NewInt64(5, false)
	o.E = 6

	errs = IsValueMatch(o, cols, values)
	if errs != nil {
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}

	o.B = time.Date(2007, 11, 2, 2, 2, 2, 2, time.UTC)
	errs = IsValueMatch(o, cols, values)
	if errs != nil {
		for _, e := range errs {
			t.Errorf("%s", e)
		}
	}
}

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

func TestRandomizeValidatedStruct(t *testing.T) {
	t.Parallel()

	var testStruct = struct {
		Int1     int
		NullInt1 null.Int
		UUID1    string
		UUID2    string
	}{}

	validatedCols := []string{
		"uuid1",
		"uuid2",
	}
	fieldTypes := map[string]string{
		"Int":     "integer",
		"NullInt": "integer",
		"UUID1":   "uuid",
		"UUID2":   "uuid",
	}

	err := RandomizeValidatedStruct(&testStruct, validatedCols, fieldTypes)
	if err != nil {
		t.Fatal(err)
	}
	if testStruct.Int1 != 0 || testStruct.NullInt1.Int != 0 ||
		testStruct.NullInt1.Valid != false {
		t.Errorf("the regular values are being randomized when they should be zero vals: %#v", testStruct)
	}

	if testStruct.UUID1 == "" || testStruct.UUID2 == "" {
		t.Errorf("the validated values should be set: %#v", testStruct)
	}
}
