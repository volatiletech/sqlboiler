package boil

import (
	"testing"
	"time"

	"gopkg.in/nullbio/null.v4"
)

func TestBind(t *testing.T) {
	t.Skip("Not implemented")
}

func TestBindPtrs_Easy(t *testing.T) {
	t.Parallel()

	testStruct := struct {
		ID   int `boil:"identifier"`
		Date time.Time
	}{}

	cols := []string{"identifier", "date"}
	ptrs, err := bindPtrs(&testStruct, cols...)
	if err != nil {
		t.Error(err)
	}

	if ptrs[0].(*int) != &testStruct.ID {
		t.Error("id is the wrong pointer")
	}
	if ptrs[1].(*time.Time) != &testStruct.Date {
		t.Error("id is the wrong pointer")
	}
}

func TestBindPtrs_Recursive(t *testing.T) {
	t.Parallel()

	testStruct := struct {
		Happy struct {
			ID int `boil:"identifier"`
		}
		Fun struct {
			ID int
		} `boil:",bind"`
	}{}

	cols := []string{"id", "fun.id"}
	ptrs, err := bindPtrs(&testStruct, cols...)
	if err != nil {
		t.Error(err)
	}

	if ptrs[0].(*int) != &testStruct.Fun.ID {
		t.Error("id is the wrong pointer")
	}
	if ptrs[1].(*int) != &testStruct.Fun.ID {
		t.Error("id is the wrong pointer")
	}
}

func TestBindPtrs_RecursiveTags(t *testing.T) {
	testStruct := struct {
		Happy struct {
			ID int `boil:"identifier"`
		} `boil:",bind"`
		Fun struct {
			ID int `boil:"identification"`
		} `boil:",bind"`
	}{}

	cols := []string{"happy.identifier", "fun.identification"}
	ptrs, err := bindPtrs(&testStruct, cols...)
	if err != nil {
		t.Error(err)
	}

	if ptrs[0].(*int) != &testStruct.Happy.ID {
		t.Error("id is the wrong pointer")
	}
	if ptrs[1].(*int) != &testStruct.Fun.ID {
		t.Error("id is the wrong pointer")
	}
}

func TestBindPtrs_Ignore(t *testing.T) {
	t.Parallel()

	testStruct := struct {
		ID    int `boil:"-"`
		Happy struct {
			ID int
		} `boil:",bind"`
	}{}

	cols := []string{"id"}
	ptrs, err := bindPtrs(&testStruct, cols...)
	if err != nil {
		t.Error(err)
	}

	if ptrs[0].(*int) != &testStruct.Happy.ID {
		t.Error("id is the wrong pointer")
	}
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
