package boil

import "testing"

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
			[]*testNestedSlice{{ID: 5}},
		}
	case *[]*testNestedSlice:
		for _, r := range *x {
			r.R = &testNestedRSlice{
				[]*testNestedSlice{{ID: 5}},
			}
		}
	}
	loadFunctionNestedCalled++
	return nil
}

func testFakeState(toLoad ...string) loadRelationshipState {
	return loadRelationshipState{
		loaded: map[string]struct{}{},
		toLoad: toLoad,
	}
}

func TestLoadRelationshipsSlice(t *testing.T) {
	// t.Parallel() Function uses globals
	loadFunctionCalled = false

	testSlice := []*struct {
		ID int
		R  *testRStruct
		L  testLStruct
	}{{}}

	if err := testFakeState("TestOne").loadRelationships(0, &testSlice, kindPtrSliceStruct); err != nil {
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

	if err := testFakeState("TestOne").loadRelationships(0, &testSingular, kindStruct); err != nil {
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
	if err := testFakeState("ToEagerLoad", "ToEagerLoad", "ToEagerLoad").loadRelationships(0, &testSlice, kindPtrSliceStruct); err != nil {
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
	if err := testFakeState("ToEagerLoad", "ToEagerLoad", "ToEagerLoad").loadRelationships(0, &testSliceSlice, kindPtrSliceStruct); err != nil {
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
	if err := testFakeState("ToEagerLoad", "ToEagerLoad", "ToEagerLoad").loadRelationships(0, &testSingular, kindStruct); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}

	testSingularSlice := testNestedSlice{
		ID: 3,
	}
	loadFunctionNestedCalled = 0
	if err := testFakeState("ToEagerLoad", "ToEagerLoad", "ToEagerLoad").loadRelationships(0, &testSingularSlice, kindStruct); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 3 {
		t.Error("Load function was called:", loadFunctionNestedCalled, "times")
	}
}

func TestLoadRelationshipsNoReload(t *testing.T) {
	// t.Parallel() Function uses globals
	testSingular := testNestedStruct{
		ID: 3,
		R: &testNestedRStruct{
			&testNestedStruct{},
		},
	}

	loadFunctionNestedCalled = 0
	state := loadRelationshipState{
		loaded: map[string]struct{}{
			"ToEagerLoad":             {},
			"ToEagerLoad.ToEagerLoad": {},
		},
		toLoad: []string{"ToEagerLoad", "ToEagerLoad"},
	}

	if err := state.loadRelationships(0, &testSingular, kindStruct); err != nil {
		t.Error(err)
	}
	if loadFunctionNestedCalled != 0 {
		t.Error("didn't want this called")
	}
}
