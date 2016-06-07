package boil

import (
	"testing"
	"time"

	"gopkg.in/BlackBaronsTux/null-extended.v1"
)

func TestBind(t *testing.T) {

}

func TestBindOne(t *testing.T) {

}

func TestBindAll(t *testing.T) {

}

func TestGetStructValues(t *testing.T) {
	t.Parallel()
	timeThing := time.Now()
	o := struct {
		TitleThing string
		Name       string
		ID         int
		Stuff      int
		Things     int
		Time       time.Time
		NullBool   null.Bool
	}{
		TitleThing: "patrick",
		Stuff:      10,
		Things:     0,
		Time:       timeThing,
		NullBool:   null.NewBool(true, false),
	}

	vals := GetStructValues(&o, "title_thing", "name", "id", "stuff", "things", "time", "null_bool")
	if vals[0].(string) != "patrick" {
		t.Errorf("Want test, got %s", vals[0])
	}
	if vals[1].(string) != "" {
		t.Errorf("Want empty string, got %s", vals[1])
	}
	if vals[2].(int) != 0 {
		t.Errorf("Want 0, got %d", vals[2])
	}
	if vals[3].(int) != 10 {
		t.Errorf("Want 10, got %d", vals[3])
	}
	if vals[4].(int) != 0 {
		t.Errorf("Want 0, got %d", vals[4])
	}
	if !vals[5].(time.Time).Equal(timeThing) {
		t.Errorf("Want %s, got %s", o.Time, vals[5])
	}
	if !vals[6].(null.Bool).IsZero() {
		t.Errorf("Want %v, got %v", o.NullBool, vals[6])
	}
}

func TestGetStructPointers(t *testing.T) {
	t.Parallel()

	o := struct {
		Title string
		ID    *int
	}{
		Title: "patrick",
	}

	ptrs := GetStructPointers(&o, "title", "id")
	*ptrs[0].(*string) = "test"
	if o.Title != "test" {
		t.Errorf("Expected test, got %s", o.Title)
	}
	x := 5
	*ptrs[1].(**int) = &x
	if *o.ID != 5 {
		t.Errorf("Expected 5, got %d", *o.ID)
	}
}

func TestCheckType(t *testing.T) {
	t.Parallel()

	type Thing struct {
	}

	validTest := []struct {
		Input    interface{}
		IsSlice  bool
		TypeName string
	}{
		{&[]*Thing{}, true, "boil.Thing"},
		{[]Thing{}, false, ""},
		{&[]Thing{}, false, ""},
		{Thing{}, false, ""},
		{new(int), false, ""},
		{5, false, ""},
		{&Thing{}, false, "boil.Thing"},
	}

	for i, test := range validTest {
		typ, isSlice, err := checkType(test.Input)
		if err != nil {
			if len(test.TypeName) > 0 {
				t.Errorf("%d) Type: %T %#v - should have succeded but got err: %v", i, test.Input, test.Input, err)
			}
			continue
		}

		if isSlice != test.IsSlice {
			t.Errorf("%d) Type: %T %#v - succeded but wrong isSlice value: %t, want %t", i, test.Input, test.Input, isSlice, test.IsSlice)
		}

		if got := typ.String(); got != test.TypeName {
			t.Errorf("%d) Type: %T %#v - succeded but wrong type name: %s, want: %s", i, test.Input, test.Input, got, test.TypeName)
		}
	}
}

func TestRandomizeStruct(t *testing.T) {
	var testStruct = struct {
		Int       int
		Int64     int64
		Float64   float64
		Bool      bool
		Time      time.Time
		String    string
		ByteSlice []byte

		Ignore int

		NullInt     null.Int
		NullFloat64 null.Float64
		NullBool    null.Bool
		NullString  null.String
		NullTime    null.Time
	}{}

	err := RandomizeStruct(&testStruct, "Ignore")
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
		testStruct.ByteSlice == nil {
		t.Errorf("the regular values are not being randomized: %#v", testStruct)
	}

	if testStruct.NullInt.Valid == false &&
		testStruct.NullFloat64.Valid == false &&
		testStruct.NullBool.Valid == false &&
		testStruct.NullString.Valid == false &&
		testStruct.NullTime.Valid == false {
		t.Errorf("the null values are not being randomized: %#v", testStruct)
	}
}
