package boil

import (
	"fmt"
	"reflect"
)

func (q *Query) Bind(obj interface{}) error {
	return nil
}

func checkType(obj interface{}) (reflect.Type, bool, error) {
	val := reflect.ValueOf(obj)
	typ := val.Type()
	kind := val.Kind()

	if kind != reflect.Ptr {
		return nil, false, fmt.Errorf("Bind must be given pointers to structs but got type: %s, kind: %s", typ.String(), kind)
	}

	typ = typ.Elem()
	kind = typ.Kind()
	isSlice := false

	switch kind {
	case reflect.Slice:
		typ = typ.Elem()
		kind = typ.Kind()
		isSlice = true
	case reflect.Struct:
		return typ, isSlice, nil
	default:
		return nil, false, fmt.Errorf("Bind was given an invalid object must be []*T or *T but got type: %s, kind: %s", typ.String(), kind)
	}

	if kind != reflect.Ptr {
		return nil, false, fmt.Errorf("Bind must be given pointers to structs but got type: %s, kind: %s", typ.String(), kind)
	}

	typ = typ.Elem()
	kind = typ.Kind()

	if kind != reflect.Struct {
		return nil, false, fmt.Errorf("Bind must be a struct but got type: %s, kind: %s", typ.String(), kind)
	}

	return typ, isSlice, nil
}
