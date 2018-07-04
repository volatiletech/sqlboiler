package randomize

import (
	"reflect"
	"testing"
	"time"
)

type MagicType struct {
	Value      int
	Randomized bool
}

func (m *MagicType) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	m.Value = int(nextInt())
	m.Randomized = true
}

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

		Magic MagicType

		Ignore int
	}{}

	fieldTypes := map[string]string{
		"Int":       "integer",
		"Int64":     "bigint",
		"Float64":   "decimal",
		"Bool":      "boolean",
		"Time":      "date",
		"String":    "character varying",
		"ByteSlice": "bytea",
		"Interval":  "interval",
		"Magic":     "magic_type",
		"Ignore":    "integer",
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
		testStruct.Magic.Value == 0 &&
		testStruct.ByteSlice == nil {
		t.Errorf("the regular values are not being randomized: %#v", testStruct)
	}

	if !testStruct.Magic.Randomized {
		t.Error("The randomize interface should have been used")
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
		{In: new(float32), Out: float32(0), Typs: []string{"real"}},
		{In: new(float64), Out: float64(0), Typs: []string{"numeric"}},
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

func TestEnumValue(t *testing.T) {
	t.Parallel()

	s := NewSeed()

	enum1 := "enum.workday('monday','tuesday')"
	enum2 := "enum('monday','tuesday')"
	enum3 := "enum('monday')"

	r1, err := EnumValue(s.NextInt, enum1)
	if err != nil {
		t.Error(err)
	}

	if r1 != "monday" && r1 != "tuesday" {
		t.Errorf("Expected monday or tuesday, got: %q", r1)
	}

	r2, err := EnumValue(s.NextInt, enum2)
	if err != nil {
		t.Error(err)
	}

	if r2 != "monday" && r2 != "tuesday" {
		t.Errorf("Expected monday or tuesday, got: %q", r2)
	}

	r3, err := EnumValue(s.NextInt, enum3)
	if err != nil {
		t.Error(err)
	}

	if r3 != "monday" {
		t.Errorf("Expected monday got: %q", r3)
	}
}
