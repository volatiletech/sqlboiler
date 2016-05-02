package boil

import (
	"fmt"
	"reflect"

	"github.com/pobri19/sqlboiler/strmangle"
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

// GetStructValues returns the values (as interface) of the matching columns in obj
func GetStructValues(obj interface{}, columns ...string) []interface{} {
	ret := make([]interface{}, len(columns))
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		ret[i] = field.Interface()
	}

	return ret
}

// GetStructPointers returns a slice of pointers to the matching columns in obj
func GetStructPointers(obj interface{}, columns ...string) []interface{} {
	val := reflect.ValueOf(obj).Elem()
	ret := make([]interface{}, len(columns))

	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		if !field.IsValid() {
			panic(fmt.Sprintf("Could not find field on struct %T for field %s", obj, strmangle.TitleCase(c)))
		}

		field = field.Addr()
		ret[i] = field.Interface()
	}

	return ret
}
