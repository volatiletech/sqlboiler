package queries

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/strmangle"
)

type loadRelationshipState struct {
	exec   boil.Executor
	loaded map[string]struct{}
	toLoad []string
}

func (l loadRelationshipState) hasLoaded(depth int) bool {
	_, ok := l.loaded[l.buildKey(depth)]
	return ok
}

func (l loadRelationshipState) setLoaded(depth int) {
	l.loaded[l.buildKey(depth)] = struct{}{}
}

func (l loadRelationshipState) buildKey(depth int) string {
	buf := strmangle.GetBuffer()

	for i, piece := range l.toLoad[:depth+1] {
		if i != 0 {
			buf.WriteByte('.')
		}
		buf.WriteString(piece)
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}

// eagerLoad loads all of the model's relationships
//
// toLoad should look like:
// []string{"Relationship", "Relationship.NestedRelationship"} ... etc
func eagerLoad(exec boil.Executor, toLoad []string, obj interface{}, bkind bindKind) error {
	state := loadRelationshipState{
		exec:   exec,
		loaded: map[string]struct{}{},
	}
	for _, toLoad := range toLoad {
		state.toLoad = strings.Split(toLoad, ".")
		if err := state.loadRelationships(0, obj, bkind); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationships dynamically calls the template generated eager load
// functions of the form:
//
//   func (t *TableR) LoadRelationshipName(exec Executor, singular bool, obj interface{})
//
// The arguments to this function are:
//   - t is not considered here, and is always passed nil. The function exists on a loaded
//     struct to avoid a circular dependency with boil, and the receiver is ignored.
//   - exec is used to perform additional queries that might be required for loading the relationships.
//   - singular is passed in to identify whether or not this was a single object
//     or a slice that must be loaded into.
//   - obj is the object or slice of objects, always of the type *obj or *[]*obj as per bind.
func (l loadRelationshipState) loadRelationships(depth int, obj interface{}, bkind bindKind) error {
	typ := reflect.TypeOf(obj).Elem()
	if bkind == kindPtrSliceStruct {
		typ = typ.Elem().Elem()
	}

	loadingFrom := reflect.ValueOf(obj)
	if loadingFrom.IsNil() {
		return nil
	}
	loadingFrom = reflect.Indirect(loadingFrom)

	fmt.Println("load rels", typ.String(), l.toLoad[depth:])

	if !l.hasLoaded(depth) {
		fmt.Println("!loaded", l.toLoad[depth])
		current := l.toLoad[depth]
		ln, found := typ.FieldByName(loaderStructName)
		// It's possible a Loaders struct doesn't exist on the struct.
		if !found {
			return errors.Errorf("attempted to load %s but no L struct was found", current)
		}

		// Attempt to find the LoadRelationshipName function
		loadMethod, found := ln.Type.MethodByName(loadMethodPrefix + current)
		if !found {
			return errors.Errorf("could not find %s%s method for eager loading", loadMethodPrefix, current)
		}

		// Hack to allow nil executors
		execArg := reflect.ValueOf(l.exec)
		if !execArg.IsValid() {
			execArg = reflect.ValueOf((*sql.DB)(nil))
		}

		// Get a loader instance from anything we have, *struct, or *[]*struct
		val := reflect.Indirect(loadingFrom)
		if bkind == kindPtrSliceStruct {
			val = reflect.Indirect(val.Index(0))
		}

		methodArgs := []reflect.Value{
			val.FieldByName(loaderStructName),
			execArg,
			reflect.ValueOf(bkind == kindStruct),
			reflect.ValueOf(obj),
		}
		resp := loadMethod.Func.Call(methodArgs)
		if intf := resp[0].Interface(); intf != nil {
			return errors.Wrapf(intf.(error), "failed to eager load %s", current)
		}

		l.setLoaded(depth)
	} else {
		fmt.Println("!loading", l.toLoad[depth])
	}

	// Check if we can stop
	if depth+1 >= len(l.toLoad) {
		return nil
	}

	// If it's singular we can just immediately call without looping
	if bkind == kindStruct {
		return l.loadRelationshipsRecurse(depth, loadingFrom)
	}

	// Loop over all eager loaded objects
	ln := loadingFrom.Len()
	if ln == 0 {
		return nil
	}
	for i := 0; i < ln; i++ {
		iter := reflect.Indirect(loadingFrom.Index(i))
		if err := l.loadRelationshipsRecurse(depth, iter); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationshipsRecurse is a helper function for taking a reflect.Value and
// Basically calls loadRelationships with: obj.R.EagerLoadedObj, and whether it's a string or slice
func (l loadRelationshipState) loadRelationshipsRecurse(depth int, obj reflect.Value) error {
	// Get relationship struct, and if it's good to go, grab the value we just loaded.
	relationshipStruct := obj.FieldByName(relationshipStructName)
	if !relationshipStruct.IsValid() || relationshipStruct.IsNil() {
		return errors.Errorf("could not traverse into loaded %s relationship to load more things", l.toLoad[depth])
	}

	loadedObject := reflect.Indirect(relationshipStruct).FieldByName(l.toLoad[depth])
	fmt.Println("loadRecurse", l.toLoad[depth])

	// Pop one off the queue
	depth++

	bkind := kindStruct
	if reflect.Indirect(loadedObject).Kind() != reflect.Struct {
		bkind = kindPtrSliceStruct
		loadedObject = loadedObject.Addr()
	}
	return l.loadRelationships(depth, loadedObject.Interface(), bkind)
}
