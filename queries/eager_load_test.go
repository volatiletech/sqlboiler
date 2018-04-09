package queries

import (
	"fmt"
	"testing"

	"github.com/ann-kilzer/sqlboiler/boil"
)

var testEagerCounters struct {
	ChildOne   int
	ChildMany  int
	NestedOne  int
	NestedMany int
}

type testEager struct {
	ID int
	R  *testEagerR
	L  testEagerL
}

type testEagerR struct {
	ChildOne  *testEagerChild
	ChildMany []*testEagerChild
	ZeroOne   *testEagerZero
	ZeroMany  []*testEagerZero
}
type testEagerL struct {
}

type testEagerChild struct {
	ID int
	R  *testEagerChildR
	L  testEagerChildL
}
type testEagerChildR struct {
	NestedOne  *testEagerNested
	NestedMany []*testEagerNested
}
type testEagerChildL struct {
}

type testEagerNested struct {
	ID int
	R  *testEagerNestedR
	L  testEagerNestedL
}
type testEagerNestedR struct {
}
type testEagerNestedL struct {
}

type testEagerZero struct {
	ID int
	R  *testEagerZeroR
	L  testEagerZeroL
}
type testEagerZeroR struct {
	NestedOne  *testEagerNested
	NestedMany []*testEagerNested
}
type testEagerZeroL struct {
}

func (testEagerL) LoadChildOne(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEager
	if singular {
		toSetOn = []*testEager{obj.(*testEager)}
	} else {
		toSetOn = *obj.(*[]*testEager)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerR{}
		}
		o.R.ChildOne = &testEagerChild{ID: 11}
	}

	testEagerCounters.ChildOne++

	return nil
}

func (testEagerL) LoadChildMany(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEager
	if singular {
		toSetOn = []*testEager{obj.(*testEager)}
	} else {
		toSetOn = *obj.(*[]*testEager)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerR{}
		}
		o.R.ChildMany = []*testEagerChild{
			&testEagerChild{ID: 12},
			&testEagerChild{ID: 13},
		}
	}

	testEagerCounters.ChildMany++

	return nil
}

func (testEagerChildL) LoadNestedOne(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEagerChild
	if singular {
		toSetOn = []*testEagerChild{obj.(*testEagerChild)}
	} else {
		toSetOn = *obj.(*[]*testEagerChild)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerChildR{}
		}
		o.R.NestedOne = &testEagerNested{ID: 21}
	}

	testEagerCounters.NestedOne++

	return nil
}

func (testEagerChildL) LoadNestedMany(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEagerChild
	if singular {
		toSetOn = []*testEagerChild{obj.(*testEagerChild)}
	} else {
		toSetOn = *obj.(*[]*testEagerChild)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerChildR{}
		}
		o.R.NestedMany = []*testEagerNested{
			&testEagerNested{ID: 22},
			&testEagerNested{ID: 23},
		}
	}

	testEagerCounters.NestedMany++

	return nil
}

func (testEagerL) LoadZeroOne(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEager
	if singular {
		toSetOn = []*testEager{obj.(*testEager)}
	} else {
		toSetOn = *obj.(*[]*testEager)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerR{}
		}
	}

	return nil
}

func (testEagerL) LoadZeroMany(_ boil.Executor, singular bool, obj interface{}) error {
	var toSetOn []*testEager
	if singular {
		toSetOn = []*testEager{obj.(*testEager)}
	} else {
		toSetOn = *obj.(*[]*testEager)
	}

	for _, o := range toSetOn {
		if o.R == nil {
			o.R = &testEagerR{}
		}
		o.R.ZeroMany = []*testEagerZero{}
	}
	return nil
}

func (testEagerZeroL) LoadNestedOne(_ boil.Executor, singular bool, obj interface{}) error {
	return nil
}

func (testEagerZeroL) LoadNestedMany(_ boil.Executor, singular bool, obj interface{}) error {
	return nil
}

func TestEagerLoadFromOne(t *testing.T) {
	testEagerCounters.ChildOne = 0
	testEagerCounters.ChildMany = 0
	testEagerCounters.NestedOne = 0
	testEagerCounters.NestedMany = 0

	obj := &testEager{}

	toLoad := []string{"ChildOne.NestedMany", "ChildOne.NestedOne", "ChildMany.NestedMany", "ChildMany.NestedOne"}
	err := eagerLoad(nil, toLoad, obj, kindStruct)
	if err != nil {
		t.Fatal(err)
	}

	if testEagerCounters.ChildMany != 1 {
		t.Error(testEagerCounters.ChildMany)
	}
	if testEagerCounters.ChildOne != 1 {
		t.Error(testEagerCounters.ChildOne)
	}
	if testEagerCounters.NestedMany != 2 {
		t.Error(testEagerCounters.NestedMany)
	}
	if testEagerCounters.NestedOne != 2 {
		t.Error(testEagerCounters.NestedOne)
	}

	checkChildOne(obj.R.ChildOne)
	checkChildMany(obj.R.ChildMany)

	checkNestedOne(obj.R.ChildOne.R.NestedOne)
	checkNestedOne(obj.R.ChildMany[0].R.NestedOne)
	checkNestedOne(obj.R.ChildMany[1].R.NestedOne)

	checkNestedMany(obj.R.ChildOne.R.NestedMany)
	checkNestedMany(obj.R.ChildMany[0].R.NestedMany)
	checkNestedMany(obj.R.ChildMany[1].R.NestedMany)
}

func TestEagerLoadFromMany(t *testing.T) {
	testEagerCounters.ChildOne = 0
	testEagerCounters.ChildMany = 0
	testEagerCounters.NestedOne = 0
	testEagerCounters.NestedMany = 0

	slice := []*testEager{
		{ID: -1},
		{ID: -2},
	}

	toLoad := []string{"ChildOne.NestedMany", "ChildOne.NestedOne", "ChildMany.NestedMany", "ChildMany.NestedOne"}
	err := eagerLoad(nil, toLoad, &slice, kindPtrSliceStruct)
	if err != nil {
		t.Fatal(err)
	}

	if testEagerCounters.ChildMany != 1 {
		t.Error(testEagerCounters.ChildMany)
	}
	if testEagerCounters.ChildOne != 1 {
		t.Error(testEagerCounters.ChildOne)
	}
	if testEagerCounters.NestedMany != 2 {
		t.Error(testEagerCounters.NestedMany)
	}
	if testEagerCounters.NestedOne != 2 {
		t.Error(testEagerCounters.NestedOne)
	}

	checkChildOne(slice[0].R.ChildOne)
	checkChildOne(slice[1].R.ChildOne)
	checkChildMany(slice[0].R.ChildMany)
	checkChildMany(slice[1].R.ChildMany)

	checkNestedOne(slice[0].R.ChildOne.R.NestedOne)
	checkNestedOne(slice[0].R.ChildMany[0].R.NestedOne)
	checkNestedOne(slice[0].R.ChildMany[1].R.NestedOne)
	checkNestedOne(slice[1].R.ChildOne.R.NestedOne)
	checkNestedOne(slice[1].R.ChildMany[0].R.NestedOne)
	checkNestedOne(slice[1].R.ChildMany[1].R.NestedOne)

	checkNestedMany(slice[0].R.ChildOne.R.NestedMany)
	checkNestedMany(slice[0].R.ChildMany[0].R.NestedMany)
	checkNestedMany(slice[0].R.ChildMany[1].R.NestedMany)
	checkNestedMany(slice[1].R.ChildOne.R.NestedMany)
	checkNestedMany(slice[1].R.ChildMany[0].R.NestedMany)
	checkNestedMany(slice[1].R.ChildMany[1].R.NestedMany)
}

func TestEagerLoadZeroParents(t *testing.T) {
	t.Parallel()

	obj := &testEager{}

	toLoad := []string{"ZeroMany.NestedMany", "ZeroOne.NestedOne", "ZeroMany.NestedMany", "ZeroOne.NestedOne"}
	err := eagerLoad(nil, toLoad, obj, kindStruct)
	if err != nil {
		t.Fatal(err)
	}

	if len(obj.R.ZeroMany) != 0 {
		t.Error("should have loaded nothing")
	}
	if obj.R.ZeroOne != nil {
		t.Error("should have loaded nothing")
	}
}

func TestEagerLoadZeroParentsMany(t *testing.T) {
	t.Parallel()

	obj := []*testEager{
		&testEager{},
		&testEager{},
	}

	toLoad := []string{"ZeroMany.NestedMany", "ZeroOne.NestedOne", "ZeroMany.NestedMany", "ZeroOne.NestedOne"}
	err := eagerLoad(nil, toLoad, &obj, kindPtrSliceStruct)
	if err != nil {
		t.Fatal(err)
	}

	if len(obj[0].R.ZeroMany) != 0 {
		t.Error("should have loaded nothing")
	}
	if obj[0].R.ZeroOne != nil {
		t.Error("should have loaded nothing")
	}
	if len(obj[1].R.ZeroMany) != 0 {
		t.Error("should have loaded nothing")
	}
	if obj[1].R.ZeroOne != nil {
		t.Error("should have loaded nothing")
	}
}

func checkChildOne(c *testEagerChild) {
	if c == nil {
		panic("c was nil")
	}

	if c.ID != 11 {
		panic(fmt.Sprintf("ChildOne id was not loaded correctly: %d", c.ID))
	}
}

func checkChildMany(cs []*testEagerChild) {
	if len(cs) != 2 {
		panic("cs len was not 2")
	}

	if cs[0].ID != 12 {
		panic(fmt.Sprintf("cs[0] had wrong id: %d", cs[0].ID))
	}
	if cs[1].ID != 13 {
		panic(fmt.Sprintf("cs[1] had wrong id: %d", cs[1].ID))
	}
}

func checkNestedOne(n *testEagerNested) {
	if n == nil {
		panic("n was nil")
	}

	if n.ID != 21 {
		panic(fmt.Sprintf("NestedOne id was not loaded correctly: %d", n.ID))
	}
}

func checkNestedMany(ns []*testEagerNested) {
	if len(ns) != 2 {
		panic("ns len was not 2")
	}

	if ns[0].ID != 22 {
		panic(fmt.Sprintf("ns[0] had wrong id: %d", ns[0].ID))
	}
	if ns[1].ID != 23 {
		panic(fmt.Sprintf("ns[1] had wrong id: %d", ns[1].ID))
	}
}
