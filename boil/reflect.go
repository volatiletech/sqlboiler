package boil

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nullbio/sqlboiler/strmangle"
	"github.com/pkg/errors"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}
)

// BindP executes the query and inserts the
// result into the passed in object pointer.
// It panics on error.
func (q *Query) BindP(obj interface{}) {
	if err := q.Bind(obj); err != nil {
		panic(WrapErr(err))
	}
}

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// Bind rules:
// - Struct tags control bind, in the form of: `boil:"name,bind"`
// - If the "name" part of the struct tag is specified, that will be used
// for binding instead of the snake_cased field name.
// - If the ",bind" option is specified on an struct field, it will be recursed
// into to look for fields for binding.
// - If the name of the struct tag is "-", this field will not be bound to
//
// Example Query:
//
//   type JoinStruct struct {
//     User1      *models.User `boil:"user,bind"`
//     User2      *models.User `boil:"friend,bind"`
//     // RandomData will not be recursed into to look for fields to
//     // bind and will not be bound to because of the - for the name.
//     RandomData myStruct     `boil:"-"`
//     // Date will not be recursed into to look for fields to bind because
//     // it does not specify ,bind in the struct tag. But it can be bound to
//     // as it does not specify a - for the name.
//     Date       time.Time
//   }
//
//   models.Users(qm.InnerJoin("users as friend on users.friend_id = friend.id")).Bind(&joinStruct)
func (q *Query) Bind(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

	var structType reflect.Type
	var sliceType reflect.Type
	var singular bool

	for i := 0; i < len(bindAccepts); i++ {
		exp := bindAccepts[i]

		if i != 0 {
			typ = typ.Elem()
			kind = typ.Kind()
		}

		if kind != exp {
			if exp == reflect.Slice || kind == reflect.Struct {
				structType = typ
				singular = true
				break
			}

			return errors.Errorf("obj type should be *[]*Type or *Type but was %q", reflect.TypeOf(obj).String())
		}

		switch kind {
		case reflect.Struct:
			structType = typ
		case reflect.Slice:
			sliceType = typ
		}
	}

	return bind(q, obj, structType, sliceType, singular)
}

func bind(q *Query, obj interface{}, structType, sliceType reflect.Type, singular bool) error {
	rows, err := ExecQueryAll(q)
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}

	cols, err := rows.Columns()
	if err != nil {
		return errors.Wrap(err, "bind failed to get column names")
	}

	var ptrSlice reflect.Value
	if !singular {
		ptrSlice = reflect.Indirect(reflect.ValueOf(obj))
	}

	for rows.Next() {
		var newStruct reflect.Value
		var pointers []interface{}

		if singular {
			pointers, err = bindPtrs(obj, cols...)
		} else {
			newStruct = reflect.New(structType)
			pointers, err = bindPtrs(newStruct.Interface(), cols...)
		}
		if err != nil {
			return err
		}

		if err := rows.Scan(pointers...); err != nil {
			return errors.Wrap(err, "failed to bind pointers to obj")
		}

		if !singular {
			ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
		}
	}

	return nil
}

func bindPtrs(obj interface{}, cols ...string) ([]interface{}, error) {
	v := reflect.ValueOf(obj)
	ptrs := make([]interface{}, len(cols))

	for i, c := range cols {
		names := strings.Split(c, ".")

		ptr, ok := findField(names, v)
		if !ok {
			return nil, errors.Errorf("bindPtrs failed to find field %s", c)
		}

		ptrs[i] = ptr
	}

	return ptrs, nil
}

func findField(names []string, v reflect.Value) (interface{}, bool) {
	if !v.IsValid() || len(names) == 0 {
		return nil, false
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, false
		}
		v = reflect.Indirect(v)
	}

	if v.Kind() != reflect.Struct {
		return nil, false
	}

	name := strmangle.TitleCase(names[0])
	typ := v.Type()

	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		fieldName, recurse := getBoilTag(f)

		if fieldName == "-" {
			continue
		}

		if recurse {
			if fieldName == name {
				names = names[1:]
			}
			if ptr, ok := findField(names, v.Field(i)); ok {
				return ptr, ok
			}
		}

		if fieldName != name || len(names) > 1 {
			continue
		}

		fieldVal := v.Field(i)
		if fieldVal.Kind() != reflect.Ptr {
			return fieldVal.Addr().Interface(), true
		}
		return fieldVal.Interface(), true
	}

	return nil, false
}

func getBoilTag(field reflect.StructField) (name string, recurse bool) {
	tag := field.Tag.Get("boil")

	if len(tag) != 0 {
		tagTokens := strings.Split(tag, ",")
		name = strmangle.TitleCase(tagTokens[0])
		recurse = len(tagTokens) > 1 && tagTokens[1] == "bind"
	}

	if len(name) == 0 {
		name = field.Name
	}

	return name, recurse
}

// GetStructValues returns the values (as interface) of the matching columns in obj
func GetStructValues(obj interface{}, columns ...string) []interface{} {
	ret := make([]interface{}, len(columns))
	val := reflect.Indirect(reflect.ValueOf(obj))

	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		ret[i] = field.Interface()
	}

	return ret
}

// GetStructPointers returns a slice of pointers to the matching columns in obj
func GetStructPointers(obj interface{}, columns ...string) []interface{} {
	val := reflect.ValueOf(obj).Elem()

	var ln int
	var getField func(reflect.Value, int) reflect.Value

	if len(columns) == 0 {
		ln = val.NumField()
		getField = func(v reflect.Value, i int) reflect.Value {
			return v.Field(i)
		}
	} else {
		ln = len(columns)
		getField = func(v reflect.Value, i int) reflect.Value {
			return v.FieldByName(strmangle.TitleCase(columns[i]))
		}
	}

	ret := make([]interface{}, ln)
	for i := 0; i < ln; i++ {
		field := getField(val, i)

		if !field.IsValid() {
			// Although this breaks the abstraction of getField above - we know that v.Field(i) can't actually
			// produce an Invalid value, so we make a hopefully safe assumption here.
			panic(fmt.Sprintf("Could not find field on struct %T for field %s", obj, strmangle.TitleCase(columns[i])))
		}

		ret[i] = field.Addr().Interface()
	}

	return ret
}
