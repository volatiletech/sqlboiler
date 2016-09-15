package queries

import (
	"database/sql"
	"reflect"

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
//
// It takes list of nested relationships to load.
func (l loadRelationshipState) loadRelationships(depth int, obj interface{}, bkind bindKind) error {
	typ := reflect.TypeOf(obj).Elem()
	if bkind == kindPtrSliceStruct {
		typ = typ.Elem().Elem()
	}

	if !l.hasLoaded(depth) {
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

		val := reflect.ValueOf(obj).Elem()
		if bkind == kindPtrSliceStruct {
			val = val.Index(0).Elem()
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
	}

	// Pull one off the queue, continue if there's still some to go
	depth++
	if depth >= len(l.toLoad) {
		return nil
	}

	loadedObject := reflect.ValueOf(obj)
	// If we eagerly loaded nothing
	if loadedObject.IsNil() {
		return nil
	}
	loadedObject = reflect.Indirect(loadedObject)

	// If it's singular we can just immediately call without looping
	if bkind == kindStruct {
		return l.loadRelationshipsRecurse(depth, loadedObject)
	}

	// Loop over all eager loaded objects
	ln := loadedObject.Len()
	if ln == 0 {
		return nil
	}
	for i := 0; i < ln; i++ {
		iter := loadedObject.Index(i).Elem()
		if err := l.loadRelationshipsRecurse(depth, iter); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationshipsRecurse is a helper function for taking a reflect.Value and
// Basically calls loadRelationships with: obj.R.EagerLoadedObj, and whether it's a string or slice
func (l loadRelationshipState) loadRelationshipsRecurse(depth int, obj reflect.Value) error {
	r := obj.FieldByName(relationshipStructName)
	if !r.IsValid() || r.IsNil() {
		return errors.Errorf("could not traverse into loaded %s relationship to load more things", l.toLoad[depth])
	}
	newObj := reflect.Indirect(r).FieldByName(l.toLoad[depth])
	bkind := kindStruct
	if reflect.Indirect(newObj).Kind() != reflect.Struct {
		bkind = kindPtrSliceStruct
		newObj = newObj.Addr()
	}
	return l.loadRelationships(depth, newObj.Interface(), bkind)
}
