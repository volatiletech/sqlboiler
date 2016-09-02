package boil

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/nullbio/null.v4"
)

func bin64(i uint64) string {
	str := strconv.FormatUint(i, 2)
	pad := 64 - len(str)
	if pad > 0 {
		str = strings.Repeat("0", pad) + str
	}

	var newStr string
	for i := 0; i < len(str); i += 8 {
		if i != 0 {
			newStr += " "
		}
		newStr += str[i : i+8]
	}

	return newStr
}

type mockRowMaker struct {
	int
	rows []driver.Value
}

func TestBind(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		ID   int
		Name string `boil:"test"`
	}{}

	query := &Query{
		from: []string{"fun"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(12)), driver.Value("cat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetExecutor(query, db)
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[0].Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if id := testResults[1].ID; id != 12 {
		t.Error("wrong ID:", id)
	}
	if name := testResults[1].Name; name != "cat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func testMakeMapping(byt ...byte) uint64 {
	var x uint64
	for i, b := range byt {
		x |= uint64(b) << (uint(i) * 8)
	}
	x |= uint64(255) << uint(len(byt)*8)
	return x
}

func TestMakeStructMapping(t *testing.T) {
	t.Parallel()

	var testStruct = struct {
		LastName    string `boil:"different"`
		AwesomeName string `boil:"awesome_name"`
		Face        string `boil:"-"`
		Nose        string

		Nested struct {
			LastName    string `boil:"different"`
			AwesomeName string `boil:"awesome_name"`
			Face        string `boil:"-"`
			Nose        string

			Nested2 struct {
				Nose string
			} `boil:",bind"`
		} `boil:",bind"`
	}{}

	got := makeStructMapping(reflect.TypeOf(testStruct), nil)

	expectMap := map[string]uint64{
		"Different":           testMakeMapping(0),
		"AwesomeName":         testMakeMapping(1),
		"Nose":                testMakeMapping(3),
		"Nested.Different":    testMakeMapping(4, 0),
		"Nested.AwesomeName":  testMakeMapping(4, 1),
		"Nested.Nose":         testMakeMapping(4, 3),
		"Nested.Nested2.Nose": testMakeMapping(4, 4, 0),
	}

	for expName, expVal := range expectMap {
		gotVal, ok := got[expName]
		if !ok {
			t.Errorf("%s) had no value", expName)
			continue
		}

		if gotVal != expVal {
			t.Errorf("%s) wrong value,\nwant: %x (%s)\ngot:  %x (%s)", expName, expVal, bin64(expVal), gotVal, bin64(gotVal))
		}
	}
}

func TestPtrFromMapping(t *testing.T) {
	t.Parallel()

	type NestedPtrs struct {
		Int         int
		IntP        *int
		NestedPtrsP *NestedPtrs
	}

	val := &NestedPtrs{
		Int:  5,
		IntP: new(int),
		NestedPtrsP: &NestedPtrs{
			Int:  6,
			IntP: new(int),
		},
	}

	v := ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(0))
	if got := *v.Interface().(*int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(1))
	if got := *v.Interface().(*int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 0))
	if got := *v.Interface().(*int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 1))
	if got := *v.Interface().(*int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
}

func TestGetBoilTag(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		FirstName   string `boil:"test_one,bind"`
		LastName    string `boil:"test_two"`
		MiddleName  string `boil:"middle_name,bind"`
		AwesomeName string `boil:"awesome_name"`
		Age         string `boil:",bind"`
		Face        string `boil:"-"`
		Nose        string
	}

	var testTitleCases = map[string]string{
		"test_one":     "TestOne",
		"test_two":     "TestTwo",
		"middle_name":  "MiddleName",
		"awesome_name": "AwesomeName",
		"age":          "Age",
		"face":         "Face",
		"nose":         "Nose",
	}

	var structFields []reflect.StructField
	typ := reflect.TypeOf(TestStruct{})
	removeOk := func(thing reflect.StructField, ok bool) reflect.StructField {
		if !ok {
			panic("Exploded")
		}
		return thing
	}
	structFields = append(structFields, removeOk(typ.FieldByName("FirstName")))
	structFields = append(structFields, removeOk(typ.FieldByName("LastName")))
	structFields = append(structFields, removeOk(typ.FieldByName("MiddleName")))
	structFields = append(structFields, removeOk(typ.FieldByName("AwesomeName")))
	structFields = append(structFields, removeOk(typ.FieldByName("Age")))
	structFields = append(structFields, removeOk(typ.FieldByName("Face")))
	structFields = append(structFields, removeOk(typ.FieldByName("Nose")))

	expect := []struct {
		Name    string
		Recurse bool
	}{
		{"TestOne", true},
		{"TestTwo", false},
		{"MiddleName", true},
		{"AwesomeName", false},
		{"Age", true},
		{"-", false},
		{"Nose", false},
	}
	for i, s := range structFields {
		name, recurse := getBoilTag(s, testTitleCases)
		if expect[i].Name != name {
			t.Errorf("Invalid name, expect %q, got %q", expect[i].Name, name)
		}
		if expect[i].Recurse != recurse {
			t.Errorf("Invalid recurse, expect %v, got %v", !recurse, recurse)
		}
	}
}

func TestBindSingular(t *testing.T) {
	t.Parallel()

	testResults := struct {
		ID   int
		Name string `boil:"test"`
	}{}

	query := &Query{
		from: []string{"fun"},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	SetExecutor(query, db)
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if id := testResults.ID; id != 35 {
		t.Error("wrong ID:", id)
	}
	if name := testResults.Name; name != "pat" {
		t.Error("wrong name:", name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

var loadFunctionCalled bool
var loadFunctionNestedCalled int

type testRStruct struct {
}
type testLStruct struct {
}

type testNestedStruct struct {
	ID int
	R  *testNestedRStruct
	L  testNestedLStruct
}
type testNestedRStruct struct {
	ToEagerLoad *testNestedStruct
}
type testNestedLStruct struct {
}

type testNestedSlice struct {
	ID int
	R  *testNestedRSlice
	L  testNestedLSlice
}
type testNestedRSlice struct {
	ToEagerLoad []*testNestedSlice
}
type testNestedLSlice struct {
}

func (testLStruct) LoadTestOne(exec Executor, singular bool, obj interface{}) error {
	loadFunctionCalled = true
	return nil
}

func (testNestedLStruct) LoadToEagerLoad(exec Executor, singular bool, obj interface{}) error {
	switch x := obj.(type) {
	case *testNestedStruct:
		x.R = &testNestedRStruct{
			&testNestedStruct{ID: 4},
		}
	case *[]*testNestedStruct:
		for _, r := range *x {
			r.R = &testNestedRStruct{
				&testNestedStruct{ID: 4},
			}
		}
	}
	loadFunctionNestedCalled++
	return nil
}

func (testNestedLSlice) LoadToEagerLoad(exec Executor, singular bool, obj interface{}) error {

	switch x := obj.(type) {
	case *testNestedSlice:
		x.R = &testNestedRSlice{
			[]*testNestedSlice{&testNestedSlice{ID: 5}},
		}
	case *[]*testNestedSlice:
		for _, r := range *x {
			r.R = &testNestedRSlice{
				[]*testNestedSlice{&testNestedSlice{ID: 5}},
			}
		}
	}
	loadFunctionNestedCalled++
	return nil
}

func TestLoadRelationshipsSlice(t *testing.T) {
	// t.Parallel() Function uses globals
	loadFunctionCalled = false

	testSlice := []*struct {
		ID int
		R  *testRStruct
		L  testLStruct
	}{}

	if err := loadRelationships(nil, []string{"TestOne"}, &testSlice, false); err != nil {
		t.Error(err)
	}

	if !loadFunctionCalled {
		t.Errorf("Load function was not called for testSlice")
	}
}

func TestLoadRelationshipsSingular(t *testing.T) {
	// t.Parallel() Function uses globals
	loadFunctionCalled = false

	testSingular := struct {
		ID int
		R  *testRStruct
		L  testLStruct
	}{}

	if err := loadRelationships(nil, []string{"TestOne"}, &testSingular, true); err != nil {
		t.Error(err)
	}

	if !loadFunctionCalled {
		t.Errorf("Load function was not called for singular")
	}
}

func TestLoadRelationshipsSliceNested(t *testing.T) {
	// t.Parallel() Function uses globals
	testSlice := []*testNestedStruct{
		{
			ID: 2,
		},
	}
	loadFunctionNestedCalled = 0
	if err := loadRelationships(nil, []string{"ToEagerLoad", "ToEagerLoad", "ToEagerLoad"}, &testSlice, false); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}

	testSliceSlice := []*testNestedSlice{
		{
			ID: 2,
		},
	}
	loadFunctionNestedCalled = 0
	if err := loadRelationships(nil, []string{"ToEagerLoad", "ToEagerLoad", "ToEagerLoad"}, &testSliceSlice, false); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}
}

func TestLoadRelationshipsSingularNested(t *testing.T) {
	// t.Parallel() Function uses globals
	testSingular := testNestedStruct{
		ID: 3,
	}
	loadFunctionNestedCalled = 0
	if err := loadRelationships(nil, []string{"ToEagerLoad", "ToEagerLoad", "ToEagerLoad"}, &testSingular, true); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}

	testSingularSlice := testNestedSlice{
		ID: 3,
	}
	loadFunctionNestedCalled = 0
	if err := loadRelationships(nil, []string{"ToEagerLoad", "ToEagerLoad", "ToEagerLoad"}, &testSingularSlice, true); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}
}

func TestBind_InnerJoin(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		Happy struct {
			ID int `boil:"identifier"`
		} `boil:",bind"`
		Fun struct {
			ID int `boil:"id"`
		} `boil:",bind"`
	}{}

	query := &Query{
		from:  []string{"fun"},
		joins: []join{{kind: JoinInner, clause: "happy as h on fun.id = h.fun_id"}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id"})
	ret.AddRow(driver.Value(int64(10)))
	ret.AddRow(driver.Value(int64(11)))
	mock.ExpectQuery(`SELECT "fun"\.\* FROM "fun" INNER JOIN happy as h on fun.id = h.fun_id;`).WillReturnRows(ret)

	SetExecutor(query, db)
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].Happy.ID; id != 0 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[0].Fun.ID; id != 10 {
		t.Error("wrong ID:", id)
	}

	if id := testResults[1].Happy.ID; id != 0 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[1].Fun.ID; id != 11 {
		t.Error("wrong ID:", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

// func TestBind_InnerJoinSelect(t *testing.T) {
// 	t.Parallel()
//
// 	testResults := []*struct {
// 		Happy struct {
// 			ID int
// 		} `boil:"h,bind"`
// 		Fun struct {
// 			ID int
// 		} `boil:",bind"`
// 	}{}
//
// 	query := &Query{
// 		selectCols: []string{"fun.id", "h.id"},
// 		from:       []string{"fun"},
// 		joins:      []join{{kind: JoinInner, clause: "happy as h on fun.happy_id = h.id"}},
// 	}
//
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	ret := sqlmock.NewRows([]string{"fun.id", "h.id"})
// 	ret.AddRow(driver.Value(int64(10)), driver.Value(int64(11)))
// 	ret.AddRow(driver.Value(int64(12)), driver.Value(int64(13)))
// 	mock.ExpectQuery(`SELECT "fun"."id" as "fun.id", "h"."id" as "h.id" FROM "fun" INNER JOIN happy as h on fun.happy_id = h.id;`).WillReturnRows(ret)
//
// 	SetExecutor(query, db)
// 	err = query.Bind(&testResults)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if len(testResults) != 2 {
// 		t.Fatal("wrong number of results:", len(testResults))
// 	}
// 	if id := testResults[0].Happy.ID; id != 11 {
// 		t.Error("wrong ID:", id)
// 	}
// 	if id := testResults[0].Fun.ID; id != 10 {
// 		t.Error("wrong ID:", id)
// 	}
//
// 	if id := testResults[1].Happy.ID; id != 13 {
// 		t.Error("wrong ID:", id)
// 	}
// 	if id := testResults[1].Fun.ID; id != 12 {
// 		t.Error("wrong ID:", id)
// 	}
//
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestBindPtrs_Easy(t *testing.T) {
// 	t.Parallel()
//
// 	testStruct := struct {
// 		ID   int `boil:"identifier"`
// 		Date time.Time
// 	}{}
//
// 	cols := []string{"identifier", "date"}
// 	ptrs, err := bindPtrs(&testStruct, nil, cols...)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if ptrs[0].(*int) != &testStruct.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// 	if ptrs[1].(*time.Time) != &testStruct.Date {
// 		t.Error("id is the wrong pointer")
// 	}
// }
//
// func TestBindPtrs_Recursive(t *testing.T) {
// 	t.Parallel()
//
// 	testStruct := struct {
// 		Happy struct {
// 			ID int `boil:"identifier"`
// 		}
// 		Fun struct {
// 			ID int
// 		} `boil:",bind"`
// 	}{}
//
// 	cols := []string{"id", "fun.id"}
// 	ptrs, err := bindPtrs(&testStruct, nil, cols...)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if ptrs[0].(*int) != &testStruct.Fun.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// 	if ptrs[1].(*int) != &testStruct.Fun.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// }
//
// func TestBindPtrs_RecursiveTags(t *testing.T) {
// 	t.Parallel()
//
// 	testStruct := struct {
// 		Happy struct {
// 			ID int `boil:"identifier"`
// 		} `boil:",bind"`
// 		Fun struct {
// 			ID int `boil:"identification"`
// 		} `boil:",bind"`
// 	}{}
//
// 	cols := []string{"happy.identifier", "fun.identification"}
// 	ptrs, err := bindPtrs(&testStruct, nil, cols...)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if ptrs[0].(*int) != &testStruct.Happy.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// 	if ptrs[1].(*int) != &testStruct.Fun.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// }
//
// func TestBindPtrs_Ignore(t *testing.T) {
// 	t.Parallel()
//
// 	testStruct := struct {
// 		ID    int `boil:"-"`
// 		Happy struct {
// 			ID int
// 		} `boil:",bind"`
// 	}{}
//
// 	cols := []string{"id"}
// 	ptrs, err := bindPtrs(&testStruct, nil, cols...)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if ptrs[0].(*int) != &testStruct.Happy.ID {
// 		t.Error("id is the wrong pointer")
// 	}
// }

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

	vals := GetStructValues(&o, nil, "title_thing", "name", "id", "stuff", "things", "time", "null_bool")
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

func TestGetSliceValues(t *testing.T) {
	t.Parallel()

	o := []struct {
		ID   int
		Name string
	}{
		{5, "a"},
		{6, "b"},
	}

	in := make([]interface{}, len(o))
	in[0] = o[0]
	in[1] = o[1]

	vals := GetSliceValues(in, nil, "id", "name")
	if got := vals[0].(int); got != 5 {
		t.Error(got)
	}
	if got := vals[1].(string); got != "a" {
		t.Error(got)
	}
	if got := vals[2].(int); got != 6 {
		t.Error(got)
	}
	if got := vals[3].(string); got != "b" {
		t.Error(got)
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

	ptrs := GetStructPointers(&o, nil, "title", "id")
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
