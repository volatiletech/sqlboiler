package boil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/strmangle"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}
)

var (
	mut         sync.RWMutex
	bindingMaps = make(map[string][]uint64)
)

const (
	loadMethodPrefix       = "Load"
	relationshipStructName = "R"
	loaderStructName       = "L"
	sentinel               = uint64(255)
)

// BindP executes the query and inserts the
// result into the passed in object pointer.
// It panics on error. See boil.Bind() documentation.
func (q *Query) BindP(obj interface{}) {
	if err := q.Bind(obj); err != nil {
		panic(WrapErr(err))
	}
}

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// Bind rules:
//   - Struct tags control bind, in the form of: `boil:"name,bind"`
//   - If "name" is omitted the sql column names that come back are TitleCased
//     and matched against the field name.
//   - If the "name" part of the struct tag is specified, the given name will
//     be used instead of the struct field name for binding.
//   - If the "name" of the struct tag is "-", this field will not be bound to.
//   - If the ",bind" option is specified on a struct field and that field
//     is a struct itself, it will be recursed into to look for fields for binding.
//
// Example Query:
//
//   type JoinStruct struct {
//     // User1 can have it's struct fields bound to since it specifies
//     // ,bind in the struct tag, it will look specifically for
//     // fields that are prefixed with "user." returning from the query.
//     // For example "user.id" column name will bind to User1.ID
//     User1      *models.User `boil:"user,bind"`
//     // User2 will follow the same rules as noted above except it will use
//     // "friend." as the prefix it's looking for.
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
//
// For custom objects that want to use eager loading, please see the
// loadRelationships function.
func Bind(rows *sql.Rows, obj interface{}) error {
	structType, sliceType, singular, err := bindChecks(obj)
	if err != nil {
		return err
	}

	return bind(rows, obj, structType, sliceType, singular)
}

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// See documentation for boil.Bind()
func (q *Query) Bind(obj interface{}) error {
	structType, sliceType, singular, err := bindChecks(obj)
	if err != nil {
		return err
	}

	rows, err := ExecQueryAll(q)
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}
	defer rows.Close()
	if res := bind(rows, obj, structType, sliceType, singular); res != nil {
		return res
	}

	state := loadRelationshipState{
		exec:   q.executor,
		loaded: map[string]struct{}{},
	}
	for _, toLoad := range q.load {
		state.toLoad = strings.Split(toLoad, ".")
		if err = state.loadRelationships(0, obj, singular); err != nil {
			return err
		}
	}
	return nil
}

// bindChecks resolves information about the bind target, and errors if it's not an object
// we can bind to.
func bindChecks(obj interface{}) (structType reflect.Type, sliceType reflect.Type, singular bool, err error) {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

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

			return nil, nil, false, errors.Errorf("obj type should be *[]*Type or *Type but was %q", reflect.TypeOf(obj).String())
		}

		switch kind {
		case reflect.Struct:
			structType = typ
		case reflect.Slice:
			sliceType = typ
		}
	}

	return structType, sliceType, singular, nil
}

func bind(rows *sql.Rows, obj interface{}, structType, sliceType reflect.Type, singular bool) error {
	cols, err := rows.Columns()
	if err != nil {
		return errors.Wrap(err, "bind failed to get column names")
	}

	var ptrSlice reflect.Value
	if !singular {
		ptrSlice = reflect.Indirect(reflect.ValueOf(obj))
	}

	var mapping []uint64
	var ok bool

	mapKey := makeCacheKey(structType.String(), cols)
	mut.RLock()
	mapping, ok = bindingMaps[mapKey]
	mut.RUnlock()

	if !ok {
		mapping, err = bindMapping(structType, cols)
		if err != nil {
			return err
		}

		mut.Lock()
		bindingMaps[mapKey] = mapping
		mut.Unlock()
	}

	foundOne := false
	for rows.Next() {
		foundOne = true
		var newStruct reflect.Value
		var pointers []interface{}

		if singular {
			pointers = ptrsFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mapping)
		} else {
			newStruct = reflect.New(structType)
			pointers = ptrsFromMapping(reflect.Indirect(newStruct), mapping)
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

	if singular && !foundOne {
		return sql.ErrNoRows
	}

	return nil
}

// bindMapping creates a mapping that helps look up the pointer for the
// column given.
func bindMapping(typ reflect.Type, cols []string) ([]uint64, error) {
	ptrs := make([]uint64, len(cols))
	mapping := makeStructMapping(typ)

ColLoop:
	for i, c := range cols {
		name := strmangle.TitleCaseIdentifier(c)
		ptrMap, ok := mapping[name]
		if ok {
			ptrs[i] = ptrMap
			continue
		}

		suffix := "." + name
		for maybeMatch, mapping := range mapping {
			if strings.HasSuffix(maybeMatch, suffix) {
				ptrs[i] = mapping
				continue ColLoop
			}
		}

		return nil, errors.Errorf("could not find struct field name in mapping: %s", name)
	}

	return ptrs, nil
}

// ptrsFromMapping expects to be passed an addressable struct that it's looking
// for things on.
func ptrsFromMapping(val reflect.Value, mapping []uint64) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m).Interface()
	}
	return ptrs
}

// ptrFromMapping expects to be passed an addressable struct that it's looking
// for things on.
func ptrFromMapping(val reflect.Value, mapping uint64) reflect.Value {
	for i := 0; i < 8; i++ {
		v := (mapping >> uint(i*8)) & sentinel

		if v == sentinel {
			if val.Kind() != reflect.Ptr {
				return val.Addr()
			}
			return val
		}

		val = val.Field(int(v))
		if val.Kind() == reflect.Ptr {
			val = reflect.Indirect(val)
		}
	}

	panic("could not find pointer from mapping")
}

func makeStructMapping(typ reflect.Type) map[string]uint64 {
	fieldMaps := make(map[string]uint64)
	makeStructMappingHelper(typ, "", 0, 0, fieldMaps)
	return fieldMaps
}

func makeStructMappingHelper(typ reflect.Type, prefix string, current uint64, depth uint, fieldMaps map[string]uint64) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)

		tag, recurse := getBoilTag(f)
		if len(tag) == 0 {
			tag = f.Name
		} else if tag[0] == '-' {
			continue
		}

		if len(prefix) != 0 {
			tag = fmt.Sprintf("%s.%s", prefix, tag)
		}

		if recurse {
			makeStructMappingHelper(f.Type, tag, current|uint64(i)<<depth, depth+8, fieldMaps)
			continue
		}

		fieldMaps[tag] = current | (sentinel << (depth + 8)) | (uint64(i) << depth)
	}
}

func getBoilTag(field reflect.StructField) (name string, recurse bool) {
	tag := field.Tag.Get("boil")
	name = field.Name

	if len(tag) == 0 {
		return name, false
	}

	ind := strings.IndexByte(tag, ',')
	if ind == -1 {
		return strmangle.TitleCase(tag), false
	} else if ind == 0 {
		return name, true
	}

	nameFragment := tag[:ind]
	return strmangle.TitleCase(nameFragment), true
}

func makeCacheKey(typ string, cols []string) string {
	buf := strmangle.GetBuffer()
	buf.WriteString(typ)
	for _, s := range cols {
		buf.WriteString(s)
	}
	mapKey := buf.String()
	strmangle.PutBuffer(buf)

	return mapKey
}

// GetStructValues returns the values (as interface) of the matching columns in obj
func GetStructValues(obj interface{}, columns ...string) []interface{} {
	ret := make([]interface{}, len(columns))
	val := reflect.Indirect(reflect.ValueOf(obj))

	for i, c := range columns {
		fieldName := strmangle.TitleCase(c)
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("unable to find field with name: %s\n%#v", fieldName, obj))
		}
		ret[i] = field.Interface()
	}

	return ret
}

// GetSliceValues returns the values (as interface) of the matching columns in obj.
func GetSliceValues(slice []interface{}, columns ...string) []interface{} {
	ret := make([]interface{}, len(slice)*len(columns))

	for i, obj := range slice {
		val := reflect.Indirect(reflect.ValueOf(obj))
		for j, c := range columns {
			fieldName := strmangle.TitleCase(c)
			field := val.FieldByName(fieldName)
			if !field.IsValid() {
				panic(fmt.Sprintf("unable to find field with name: %s\n%#v", fieldName, obj))
			}
			ret[i*len(columns)+j] = field.Interface()
		}
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
