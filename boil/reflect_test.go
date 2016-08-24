package boil

import (
	"database/sql/driver"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gopkg.in/nullbio/null.v4"
)

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

type testLoadedStruct struct{}

func (r *testLoadedStruct) LoadTestOne(exec Executor, singular bool, obj interface{}) error {
	loadFunctionCalled = true
	return nil
}

func TestLoadRelationshipsSlice(t *testing.T) {
	// t.Parallel() Function uses globals
	loadFunctionCalled = false

	testSlice := []*struct {
		ID     int
		Loaded *testLoadedStruct
	}{}

	q := Query{load: []string{"TestOne"}, executor: nil}
	if err := q.loadRelationships(&testSlice, false); err != nil {
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
		ID     int
		Loaded *testLoadedStruct
	}{}

	q := Query{load: []string{"TestOne"}, executor: nil}
	if err := q.loadRelationships(&testSingular, true); err != nil {
		t.Error(err)
	}

	if !loadFunctionCalled {
		t.Errorf("Load function was not called for singular")
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

func TestBind_InnerJoinSelect(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		Happy struct {
			ID int
		} `boil:"h,bind"`
		Fun struct {
			ID int
		} `boil:",bind"`
	}{}

	query := &Query{
		selectCols: []string{"fun.id", "h.id"},
		from:       []string{"fun"},
		joins:      []join{{kind: JoinInner, clause: "happy as h on fun.happy_id = h.id"}},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"fun.id", "h.id"})
	ret.AddRow(driver.Value(int64(10)), driver.Value(int64(11)))
	ret.AddRow(driver.Value(int64(12)), driver.Value(int64(13)))
	mock.ExpectQuery(`SELECT "fun"."id" as "fun.id", "h"."id" as "h.id" FROM "fun" INNER JOIN happy as h on fun.happy_id = h.id;`).WillReturnRows(ret)

	SetExecutor(query, db)
	err = query.Bind(&testResults)
	if err != nil {
		t.Error(err)
	}

	if len(testResults) != 2 {
		t.Fatal("wrong number of results:", len(testResults))
	}
	if id := testResults[0].Happy.ID; id != 11 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[0].Fun.ID; id != 10 {
		t.Error("wrong ID:", id)
	}

	if id := testResults[1].Happy.ID; id != 13 {
		t.Error("wrong ID:", id)
	}
	if id := testResults[1].Fun.ID; id != 12 {
		t.Error("wrong ID:", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
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
	t.Parallel()

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

	vals := GetSliceValues(in, "id", "name")
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
