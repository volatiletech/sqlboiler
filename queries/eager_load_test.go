package queries

import (
	"fmt"
	"testing"

	"github.com/vattle/sqlboiler/boil"
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
		o.R.ChildOne = &testEagerChild{ID: 1}
	}

	testEagerCounters.ChildOne++
	fmt.Println("l! ChildOne")

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
			&testEagerChild{ID: 2},
			&testEagerChild{ID: 3},
		}
	}

	testEagerCounters.ChildMany++
	fmt.Println("l! ChildMany")

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
		o.R.NestedOne = &testEagerNested{ID: 6}
	}

	testEagerCounters.NestedOne++
	fmt.Println("l! NestedOne")

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
			&testEagerNested{ID: 6},
			&testEagerNested{ID: 7},
		}
	}

	testEagerCounters.NestedMany++
	fmt.Println("l! NestedMany")

	return nil
}

func TestEagerLoadFromOne(t *testing.T) {
	testEagerCounters.ChildOne = 0
	testEagerCounters.ChildMany = 0
	testEagerCounters.NestedOne = 0
	testEagerCounters.NestedMany = 0

	obj := &testEager{}

	toLoad := []string{"ChildOne", "ChildMany.NestedMany", "ChildMany.NestedOne"}
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
	if testEagerCounters.NestedMany != 1 {
		t.Error(testEagerCounters.NestedMany)
	}
	if testEagerCounters.NestedOne != 1 {
		t.Error(testEagerCounters.NestedOne)
	}
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

	toLoad := []string{"ChildOne", "ChildMany.NestedMany", "ChildMany.NestedOne"}
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
	if testEagerCounters.NestedMany != 1 {
		t.Error(testEagerCounters.NestedMany)
	}
	if testEagerCounters.NestedOne != 1 {
		t.Error(testEagerCounters.NestedOne)
	}
}
