package queries

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/drivers"

	"github.com/DATA-DOG/go-sqlmock"
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

func TestBindStruct(t *testing.T) {
	t.Parallel()

	testResults := struct {
		ID   int
		Name string `boil:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(65)), driver.Value("hat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	err = query.Bind(nil, db, &testResults)
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

func TestBindSlice(t *testing.T) {
	t.Parallel()

	testResults := []struct {
		ID   int
		Name string `boil:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(12)), driver.Value("cat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	err = query.Bind(nil, db, &testResults)
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

func TestBindPtrSlice(t *testing.T) {
	t.Parallel()

	testResults := []*struct {
		ID   int
		Name string `boil:"test"`
	}{}

	query := &Query{
		from:    []string{"fun"},
		dialect: &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	ret.AddRow(driver.Value(int64(12)), driver.Value("cat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	err = query.Bind(context.Background(), db, &testResults)
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

	testStruct := struct {
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

	got := MakeStructMapping(reflect.TypeOf(testStruct))

	expectMap := map[string]uint64{
		"different":           testMakeMapping(0),
		"awesome_name":        testMakeMapping(1),
		"nose":                testMakeMapping(3),
		"nested.different":    testMakeMapping(4, 0),
		"nested.awesome_name": testMakeMapping(4, 1),
		"nested.nose":         testMakeMapping(4, 3),
		"nested.nested2.nose": testMakeMapping(4, 4, 0),
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
			Int: 6,
		},
	}

	v := ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(0), true)
	if got := *v.Interface().(*int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(1), true)
	if got := *v.Interface().(*int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 0), true)
	if got := *v.Interface().(*int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	v = ptrFromMapping(reflect.Indirect(reflect.ValueOf(val)), testMakeMapping(2, 1), true)
	if got := *v.Interface().(*int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
}

func TestValuesFromMapping(t *testing.T) {
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
	mapping := []uint64{testMakeMapping(0), testMakeMapping(1), testMakeMapping(2, 0), testMakeMapping(2, 1), 0}
	v := ValuesFromMapping(reflect.Indirect(reflect.ValueOf(val)), mapping)

	if got := v[0].(int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	if got := v[1].(int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	if got := v[2].(int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	if got := v[3].(int); got != 0 {
		t.Error("nested pointer was wrong:", got)
	}
	if got := *v[4].(*interface{}); got != nil {
		t.Error("nil pointer was not be ignored:", got)
	}
}

func TestPtrsFromMapping(t *testing.T) {
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

	mapping := []uint64{testMakeMapping(0), testMakeMapping(1), testMakeMapping(2, 0), testMakeMapping(2, 1)}
	v := PtrsFromMapping(reflect.Indirect(reflect.ValueOf(val)), mapping)

	if got := *v[0].(*int); got != 5 {
		t.Error("flat int was wrong:", got)
	}
	if got := *v[1].(*int); got != 0 {
		t.Error("flat pointer was wrong:", got)
	}
	if got := *v[2].(*int); got != 6 {
		t.Error("nested int was wrong:", got)
	}
	if got := *v[3].(*int); got != 0 {
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
		{"test_one", true},
		{"test_two", false},
		{"middle_name", true},
		{"awesome_name", false},
		{"", true},
		{"-", false},
		{"", false},
	}
	for i, s := range structFields {
		name, recurse := getBoilTag(s)
		if expect[i].Name != name {
			t.Errorf("Invalid name, expect %q, got %q", expect[i].Name, name)
		}
		if expect[i].Recurse != recurse {
			t.Errorf("Invalid recurse, expect %v, got %v", !recurse, recurse)
		}
	}
}

func TestBindChecks(t *testing.T) {
	t.Parallel()

	type useless struct{}

	tests := []struct {
		BKind bindKind
		Fail  bool
		Obj   interface{}
	}{
		{BKind: kindStruct, Fail: false, Obj: &useless{}},
		{BKind: kindSliceStruct, Fail: false, Obj: &[]useless{}},
		{BKind: kindPtrSliceStruct, Fail: false, Obj: &[]*useless{}},
		{Fail: true, Obj: 5},
		{Fail: true, Obj: useless{}},
		{Fail: true, Obj: []useless{}},
	}

	for i, test := range tests {
		str, sli, bk, err := bindChecks(test.Obj)

		if err != nil {
			if !test.Fail {
				t.Errorf("%d) should not fail, got: %v", i, err)
			}
			continue
		} else if test.Fail {
			t.Errorf("%d) should fail, got: %v", i, bk)
			continue
		}

		if s := str.Kind(); s != reflect.Struct {
			t.Error("struct kind was wrong:", s)
		}
		if test.BKind != kindStruct {
			if s := sli.Kind(); s != reflect.Slice {
				t.Error("slice kind was wrong:", s)
			}
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
		from:    []string{"fun"},
		dialect: &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id", "test"})
	ret.AddRow(driver.Value(int64(35)), driver.Value("pat"))
	mock.ExpectQuery(`SELECT \* FROM "fun";`).WillReturnRows(ret)

	err = query.Bind(nil, db, &testResults)
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
		from:    []string{"fun"},
		joins:   []join{{kind: JoinInner, clause: "happy as h on fun.id = h.fun_id"}},
		dialect: &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	ret := sqlmock.NewRows([]string{"id"})
	ret.AddRow(driver.Value(int64(10)))
	ret.AddRow(driver.Value(int64(11)))
	mock.ExpectQuery(`SELECT "fun"\.\* FROM "fun" INNER JOIN happy as h on fun.id = h.fun_id;`).WillReturnRows(ret)

	err = query.Bind(nil, db, &testResults)
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
		dialect:    &drivers.Dialect{LQ: '"', RQ: '"', UseIndexPlaceholders: true},
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

	err = query.Bind(nil, db, &testResults)
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

func TestEqual(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		A    interface{}
		B    interface{}
		Want bool
	}{
		{A: int(5), B: int(5), Want: true},
		{A: int(5), B: int(6), Want: false},
		{A: int(5), B: int32(5), Want: true},
		{A: []byte("hello"), B: []byte("hello"), Want: true},
		{A: []byte("hello"), B: []byte("world"), Want: false},
		{A: "hello", B: sql.NullString{String: "hello", Valid: true}, Want: true},
		{A: "hello", B: sql.NullString{Valid: false}, Want: false},
		{A: now, B: now, Want: true},
		{A: now, B: now.Add(time.Hour), Want: false},
		{A: null.Uint64From(uint64(9223372036854775808)), B: uint64(9223372036854775808), Want: true},
		{A: null.Uint64From(uint64(9223372036854775808)), B: uint64(9223372036854775809), Want: false},
	}

	for i, test := range tests {
		if got := Equal(test.A, test.B); got != test.Want {
			t.Errorf("%d) compare %#v and %#v resulted in wrong value, want: %t, got %t", i, test.A, test.B, test.Want, got)
		}
	}
}

func TestAssignBytes(t *testing.T) {
	t.Parallel()

	var dst []byte
	src := []byte("hello")

	Assign(&dst, src)
	if !bytes.Equal(dst, src) {
		t.Error("bytes were not equal!")
	}
}

func TestAssignScanValue(t *testing.T) {
	t.Parallel()

	var nsDst sql.NullString
	var nsSrc sql.NullString

	nsSrc.String = "hello"
	nsSrc.Valid = true

	Assign(&nsDst, nsSrc)

	if !nsDst.Valid {
		t.Error("n was still null")
	}
	if nsDst.String != "hello" {
		t.Error("assignment did not occur")
	}

	var niDst sql.NullInt64
	var niSrc sql.NullInt64

	niSrc.Valid = true
	niSrc.Int64 = 5

	Assign(&niDst, niSrc)

	if !niDst.Valid {
		t.Error("n was still null")
	}
	if niDst.Int64 != 5 {
		t.Error("assignment did not occur")
	}
}

func TestAssignScanNoValue(t *testing.T) {
	t.Parallel()

	var ns sql.NullString
	s := "hello"

	Assign(&ns, s)

	if !ns.Valid {
		t.Error("n was still null")
	}
	if ns.String != "hello" {
		t.Error("assignment did not occur")
	}

	var niDst sql.NullInt64
	i := 5

	Assign(&niDst, i)

	if !niDst.Valid {
		t.Error("n was still null")
	}
	if niDst.Int64 != 5 {
		t.Error("assignment did not occur")
	}
}

func TestAssignNoScanValue(t *testing.T) {
	t.Parallel()

	var ns sql.NullString
	var s string

	ns.String = "hello"
	ns.Valid = true
	Assign(&s, ns)

	if s != "hello" {
		t.Error("assignment did not occur")
	}

	var ni sql.NullInt64
	var i int

	ni.Int64 = 5
	ni.Valid = true
	Assign(&i, ni)

	if i != 5 {
		t.Error("assignment did not occur")
	}
}

func TestAssignNil(t *testing.T) {
	t.Parallel()

	var ns sql.NullString
	s := "hello"

	Assign(&s, ns)
	if s != "" {
		t.Errorf("should have assigned a zero value: %q", s)
	}
}

func TestAssignPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected a panic")
		}
	}()

	var aint, bint int
	Assign(&aint, bint)
}

type nullTime struct {
	Time  time.Time
	Valid bool
}

func (t *nullTime) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		t.Time = x
	case nil:
		t.Valid = false
		return nil
	default:
		err = fmt.Errorf("cannot scan type %T into nullTime: %v", value, value)
	}
	t.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (t nullTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

func TestMustTime(t *testing.T) {
	t.Parallel()

	var nt nullTime

	if !MustTime(nt).IsZero() {
		t.Error("should be zero")
	}

	now := time.Now()

	nt.Valid = true
	nt.Time = now

	if !MustTime(nt).Equal(now) {
		t.Error("time was wrong")
	}
}

func TestMustTimePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("it should have panic'd")
		}
	}()

	var ns sql.NullString
	ns.Valid = true
	ns.String = "hello"
	MustTime(ns)
}

func TestIsValuerNil(t *testing.T) {
	t.Parallel()

	var ns sql.NullString
	if !IsValuerNil(ns) {
		t.Error("it should be nil")
	}

	ns.Valid = true
	if IsValuerNil(ns) {
		t.Error("it should not be nil")
	}
}

func TestSetScanner(t *testing.T) {
	t.Parallel()

	var ns sql.NullString
	SetScanner(&ns, "hello")

	if !ns.Valid {
		t.Error("it should not be null")
	}
	if ns.String != "hello" {
		t.Error("it's value should have been hello")
	}
}

func TestSetScannerPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("it should have panic'd")
		}
	}()

	var ns nullTime
	SetScanner(&ns, "hello")
}

func TestUnTitleCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		In  string
		Out string
	}{
		{"HelloThere", "hello_there"},
		{"", ""},
		{"AA", "aa"},
		{"FunID", "fun_id"},
		{"UID", "uid"},
		{"GUID", "guid"},
		{"UID", "uid"},
		{"UUID", "uuid"},
		{"SSN", "ssn"},
		{"TZ", "tz"},
		{"ThingGUID", "thing_guid"},
		{"GUIDThing", "guid_thing"},
		{"ThingGUIDThing", "thing_guid_thing"},
		{"ID", "id"},
		{"GVZXC", "gvzxc"},
		{"IDTRGBID", "id_trgb_id"},
		{"ThingZXCStuffVXZ", "thing_zxc_stuff_vxz"},
		{"ZXCThingVXZStuff", "zxc_thing_vxz_stuff"},
		{"ZXCVDF9C9Hello9", "zxcvdf9_c9_hello9"},
		{"ID9UID911GUID9E9", "id9_uid911_guid9_e9"},
		{"ZXCVDF0C0Hello0", "zxcvdf0_c0_hello0"},
		{"ID0UID000GUID0E0", "id0_uid000_guid0_e0"},
		{"Ab5ZXC5D5", "ab5_zxc5_d5"},
		{"Identifier", "identifier"},
	}

	for i, test := range tests {
		if out := unTitleCase(test.In); out != test.Out {
			t.Errorf("[%d] (%s) Out was wrong: %q, want: %q", i, test.In, out, test.Out)
		}
	}
}
