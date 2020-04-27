package queries

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/strmangle"
)

type loadRelationshipState struct {
	ctx    context.Context
	exec   boil.Executor
	loaded map[string]struct{}
	toLoad []string
	mods   map[string]Applicator
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
// obj should be one of:
// *[]*struct or *struct
// bkind should reflect what kind of thing it is above
func eagerLoad(ctx context.Context, exec boil.Executor, toLoad []string, mods map[string]Applicator, obj interface{}, bkind bindKind) error {
	state := loadRelationshipState{
		ctx:    ctx, // defiant to the end, I know this is frowned upon
		exec:   exec,
		loaded: map[string]struct{}{},
		mods:   mods,
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
//   - bkind is passed in to identify whether or not this was a single object
//     or a slice that must be loaded into.
//   - obj is the object or slice of objects, always of the type *obj or *[]*obj as per bind.
//
// We start with a normal select before eager loading anything: select * from a;
// Then we start eager loading things, it can be represented by a DAG
//          a1, a2           select id, a_id from b where id in (a1, a2)
//         / |    \
//        b1 b2    b3        select id, b_id from c where id in (b2, b3, b4)
//       /   | \     \
//      c1  c2 c3    c4
//
// That's to say that we descend the graph of relationships, and at each level
// we gather all the things up we want to load into, load them, and then move
// to the next level of the graph.
func (l loadRelationshipState) loadRelationships(depth int, obj interface{}, bkind bindKind) error {
	typ := reflect.TypeOf(obj).Elem()
	if bkind == kindPtrSliceStruct {
		typ = typ.Elem().Elem()
	}

	loadingFrom := reflect.ValueOf(obj)
	if loadingFrom.IsNil() {
		return nil
	}

	if !l.hasLoaded(depth) {
		if err := l.callLoadFunction(depth, loadingFrom, typ, bkind); err != nil {
			return err
		}
	}

	// Check if we can stop
	if depth+1 >= len(l.toLoad) {
		return nil
	}

	// *[]*struct -> []*struct
	// *struct -> struct
	loadingFrom = reflect.Indirect(loadingFrom)

	// If it's singular we can just immediately call without looping
	if bkind == kindStruct {
		return l.loadRelationshipsRecurse(depth, loadingFrom)
	}

	// If we were an empty slice to begin with, bail, probably a useless check
	if loadingFrom.Len() == 0 {
		return nil
	}

	// Collect eagerly loaded things to send into next eager load call
	slice, nextBKind, err := collectLoaded(l.toLoad[depth], loadingFrom)
	if err != nil {
		return err
	}

	// If we could collect nothing we're done
	if slice.Len() == 0 {
		return nil
	}

	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)

	return l.loadRelationships(depth+1, ptr.Interface(), nextBKind)
}

// callLoadFunction finds the loader struct, finds the method that we need
// to call and calls it.
func (l loadRelationshipState) callLoadFunction(depth int, loadingFrom reflect.Value, typ reflect.Type, bkind bindKind) error {
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

	ctxArg := reflect.ValueOf(l.ctx)
	// Hack to allow nil executors
	execArg := reflect.ValueOf(l.exec)
	if !execArg.IsValid() {
		execArg = reflect.ValueOf((*sql.DB)(nil))
	}

	// Get a loader instance from anything we have, *struct, or *[]*struct
	val := reflect.Indirect(loadingFrom)
	if bkind == kindPtrSliceStruct {
		if val.Len() == 0 {
			return nil
		}
		val = val.Index(0)
		if val.IsNil() {
			return nil
		}
		val = reflect.Indirect(val)
	}

	methodArgs := make([]reflect.Value, 0, 5)
	methodArgs = append(methodArgs, val.FieldByName(loaderStructName))
	if ctxArg.IsValid() {
		methodArgs = append(methodArgs, ctxArg)
	}
	methodArgs = append(methodArgs, execArg, reflect.ValueOf(bkind == kindStruct), loadingFrom)
	if mods, ok := l.mods[l.buildKey(depth)]; ok {
		methodArgs = append(methodArgs, reflect.ValueOf(mods))
	} else {
		methodArgs = append(methodArgs, applicatorSentinelVal)
	}

	ret := loadMethod.Func.Call(methodArgs)
	if intf := ret[0].Interface(); intf != nil {
		return errors.Wrapf(intf.(error), "failed to eager load %s", current)
	}

	l.setLoaded(depth)
	return nil
}

// loadRelationshipsRecurse is a helper function for taking a reflect.Value and
// Basically calls loadRelationships with: obj.R.EagerLoadedObj
// Called with an obj of *struct
func (l loadRelationshipState) loadRelationshipsRecurse(depth int, obj reflect.Value) error {
	key := l.toLoad[depth]
	r, err := findRelationshipStruct(obj)
	if err != nil {
		return errors.Wrapf(err, "failed to append loaded %s", key)
	}

	loadedObject := reflect.Indirect(r).FieldByName(key)
	if loadedObject.IsNil() {
		return nil
	}

	bkind := kindStruct
	if derefed := reflect.Indirect(loadedObject); derefed.Kind() != reflect.Struct {
		bkind = kindPtrSliceStruct

		// Convert away any helper slice types
		// elemType is *elem (from []*elem or helperSliceType)
		// sliceType is *[]*elem
		elemType := derefed.Type().Elem()
		sliceType := reflect.PtrTo(reflect.SliceOf(elemType))

		loadedObject = loadedObject.Addr().Convert(sliceType)
	}
	return l.loadRelationships(depth+1, loadedObject.Interface(), bkind)
}

// collectLoaded traverses the next level of the graph and picks up all
// the values that we need for the next eager load query.
//
// For example when loadingFrom is [parent1, parent2]
//
//   parent1 -> child1
//          \-> child2
//   parent2 -> child3
//
// This should return [child1, child2, child3]
func collectLoaded(key string, loadingFrom reflect.Value) (reflect.Value, bindKind, error) {
	// Pull the first one so we can get the types out of it in order to
	// create the proper type of slice.
	current := reflect.Indirect(loadingFrom.Index(0))
	lnFrom := loadingFrom.Len()

	r, err := findRelationshipStruct(current)
	if err != nil {
		return reflect.Value{}, 0, errors.Wrapf(err, "failed to collect loaded %s", key)
	}

	loadedObject := reflect.Indirect(r).FieldByName(key)
	loadedType := loadedObject.Type() // Should be *obj or []*obj

	bkind := kindPtrSliceStruct
	if loadedType.Elem().Kind() == reflect.Struct {
		bkind = kindStruct
		loadedType = reflect.SliceOf(loadedType)
	} else {
		// Ensure that we get rid of all the helper "XSlice" types
		loadedType = reflect.SliceOf(loadedType.Elem())
	}

	collection := reflect.MakeSlice(loadedType, 0, 0)

	i := 0
	for {
		switch bkind {
		case kindStruct:
			if !loadedObject.IsNil() {
				collection = reflect.Append(collection, loadedObject)
			}
		case kindPtrSliceStruct:
			collection = reflect.AppendSlice(collection, loadedObject)
		}

		i++
		if i >= lnFrom {
			break
		}

		current = reflect.Indirect(loadingFrom.Index(i))
		r, err = findRelationshipStruct(current)
		if err != nil {
			return reflect.Value{}, 0, errors.Wrapf(err, "failed to collect loaded %s", key)
		}

		loadedObject = reflect.Indirect(r).FieldByName(key)
	}

	return collection, kindPtrSliceStruct, nil
}

func findRelationshipStruct(obj reflect.Value) (reflect.Value, error) {
	relationshipStruct := obj.FieldByName(relationshipStructName)
	if !relationshipStruct.IsValid() {
		return reflect.Value{}, errors.New("relationship struct was invalid")
	} else if relationshipStruct.IsNil() {
		return reflect.Value{}, errors.New("relationship struct was nil")
	}

	return relationshipStruct, nil
}

var (
	applicatorSentinel    Applicator
	applicatorSentinelVal = reflect.ValueOf(&applicatorSentinel).Elem()
)
