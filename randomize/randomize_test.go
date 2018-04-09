package randomize

import (
	"reflect"
	"testing"
	"time"

	null "gopkg.in/volatiletech/null.v6"
)

func TestRandomizeStruct(t *testing.T) {
	t.Parallel()

	s := NewSeed()

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
		"NullBool":     "boolean",
		"NullString":   "character",
		"NullTime":     "time",
		"NullInterval": "interval",
	}

	err := Struct(s, &testStruct, fieldTypes, true, "Ignore")
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

func TestRandomizeField(t *testing.T) {
	t.Parallel()

	type RandomizeTest struct {
		In   interface{}
		Typs []string
		Out  interface{}
	}

	s := NewSeed()
	inputs := []RandomizeTest{
		{In: &null.Bool{}, Out: null.Bool{}, Typs: []string{"boolean"}},
		{In: &null.String{}, Out: null.String{}, Typs: []string{"character", "uuid", "interval", "numeric"}},
		{In: &null.Time{}, Out: null.Time{}, Typs: []string{"time"}},
		{In: &null.Float32{}, Out: null.Float32{}, Typs: []string{"real"}},
		{In: &null.Float64{}, Out: null.Float64{}, Typs: []string{"decimal"}},
		{In: &null.Int{}, Out: null.Int{}, Typs: []string{"integer"}},
		{In: &null.Int8{}, Out: null.Int8{}, Typs: []string{"integer"}},
		{In: &null.Int16{}, Out: null.Int16{}, Typs: []string{"smallint"}},
		{In: &null.Int32{}, Out: null.Int32{}, Typs: []string{"integer"}},
		{In: &null.Int64{}, Out: null.Int64{}, Typs: []string{"bigint"}},
		{In: &null.Uint{}, Out: null.Uint{}, Typs: []string{"integer"}},
		{In: &null.Uint8{}, Out: null.Uint8{}, Typs: []string{"integer"}},
		{In: &null.Uint16{}, Out: null.Uint16{}, Typs: []string{"integer"}},
		{In: &null.Uint32{}, Out: null.Uint32{}, Typs: []string{"integer"}},
		{In: &null.Uint64{}, Out: null.Uint64{}, Typs: []string{"integer"}},

		{In: new(float32), Out: float32(0), Typs: []string{"real"}},
		{In: new(int), Out: int(0), Typs: []string{"integer"}},
		{In: new(int8), Out: int8(0), Typs: []string{"integer"}},
		{In: new(int16), Out: int16(0), Typs: []string{"smallserial"}},
		{In: new(int32), Out: int32(0), Typs: []string{"integer"}},
		{In: new(int64), Out: int64(0), Typs: []string{"bigserial"}},
		{In: new(uint), Out: uint(0), Typs: []string{"integer"}},
		{In: new(uint8), Out: uint8(0), Typs: []string{"integer"}},
		{In: new(uint16), Out: uint16(0), Typs: []string{"integer"}},
		{In: new(uint32), Out: uint32(0), Typs: []string{"integer"}},
		{In: new(uint64), Out: uint64(0), Typs: []string{"integer"}},

		{In: new(bool), Out: false},
		{In: new(string), Out: ""},
		{In: new([]byte), Out: new([]byte)},
		{In: &time.Time{}, Out: &time.Time{}},
	}

	for i := 0; i < len(inputs); i++ {
		for _, typ := range inputs[i].Typs {
			val := reflect.Indirect(reflect.ValueOf(&inputs[i]))
			field := val.FieldByName("In").Elem().Elem()

			// Make sure we never get back values that would be considered null
			// by the boil whitelist generator, or by the database driver
			if err := randomizeField(s, field, typ, false); err != nil {
				t.Errorf("%d) %s", i, err)
			}

			if inputs[i].In == inputs[i].Out {
				t.Errorf("%d) Field should not be null, got: %v -- type: %s\n", i, inputs[i].In, typ)
			}
		}
	}
}

func TestRandEnumValue(t *testing.T) {
	t.Parallel()

	s := NewSeed()

	enum1 := "enum.workday('monday','tuesday')"
	enum2 := "enum('monday','tuesday')"
	enum3 := "enum('monday')"

	r1, err := randEnumValue(s, enum1)
	if err != nil {
		t.Error(err)
	}

	if r1 != "monday" && r1 != "tuesday" {
		t.Errorf("Expected monday or tuesday, got: %q", r1)
	}

	r2, err := randEnumValue(s, enum2)
	if err != nil {
		t.Error(err)
	}

	if r2 != "monday" && r2 != "tuesday" {
		t.Errorf("Expected monday or tuesday, got: %q", r2)
	}

	r3, err := randEnumValue(s, enum3)
	if err != nil {
		t.Error(err)
	}

	if r3 != "monday" {
		t.Errorf("Expected monday got: %q", r3)
	}
}
