package boil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/strmangle"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}
)

const (
	loadMethodPrefix       = "Load"
	relationshipStructName = "R"
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
	return BindFast(rows, obj, nil)
}

// BindFast uses a lookup table for column_name to ColumnName to avoid TitleCase.
func BindFast(rows *sql.Rows, obj interface{}, titleCases map[string]string) error {
	structType, sliceType, singular, err := bindChecks(obj)

	if err != nil {
		return err
	}

	return bind(rows, obj, structType, sliceType, singular, titleCases)
}

// Bind executes the query and inserts the
// result into the passed in object pointer
//
// See documentation for boil.Bind()
func (q *Query) Bind(obj interface{}) error {
	return q.BindFast(obj, nil)
}

// BindFast uses a lookup table for column_name to ColumnName to avoid TitleCase.
func (q *Query) BindFast(obj interface{}, titleCases map[string]string) error {
	structType, sliceType, singular, err := bindChecks(obj)
	if err != nil {
		return err
	}

	rows, err := ExecQueryAll(q)
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}
	defer rows.Close()

	if res := bind(rows, obj, structType, sliceType, singular, titleCases); res != nil {
		return res
	}

	for _, toLoad := range q.load {
		toLoadFragments := strings.Split(toLoad, ".")
		if err = loadRelationships(q.executor, toLoadFragments, obj, singular); err != nil {
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
//
// It takes list of nested relationships to load.
func loadRelationships(exec Executor, toLoad []string, obj interface{}, singular bool) error {
	typ := reflect.TypeOf(obj).Elem()
	if !singular {
		typ = typ.Elem().Elem()
	}

	current := toLoad[0]
	r, found := typ.FieldByName(relationshipStructName)
	// It's possible a Relationship struct doesn't exist on the struct.
	if !found {
		return errors.Errorf("attempted to load %s but no R struct was found", current)
	}

	// Attempt to find the LoadRelationshipName function
	loadMethod, found := r.Type.MethodByName(loadMethodPrefix + current)
	if !found {
		return errors.Errorf("could not find %s%s method for eager loading", loadMethodPrefix, current)
	}

	// Hack to allow nil executors
	execArg := reflect.ValueOf(exec)
	if !execArg.IsValid() {
		execArg = reflect.ValueOf((*sql.DB)(nil))
	}

	methodArgs := []reflect.Value{
		reflect.Indirect(reflect.New(r.Type)),
		execArg,
		reflect.ValueOf(singular),
		reflect.ValueOf(obj),
	}

	resp := loadMethod.Func.Call(methodArgs)
	if resp[0].Interface() != nil {
		return errors.Wrapf(resp[0].Interface().(error), "failed to eager load %s", current)
	}

	// Pull one off the queue, continue if there's still some to go
	toLoad = toLoad[1:]
	if len(toLoad) == 0 {
		return nil
	}

	loadedObject := reflect.ValueOf(obj)
	// If we eagerly loaded nothing
	if loadedObject.IsNil() {
		return nil
	}
	loadedObject = reflect.Indirect(loadedObject)

	// If it's singular we can just immediately call without looping
	if singular {
		return loadRelationshipsRecurse(exec, current, toLoad, singular, loadedObject)
	}

	// Loop over all eager loaded objects
	ln := loadedObject.Len()
	if ln == 0 {
		return nil
	}
	for i := 0; i < ln; i++ {
		iter := loadedObject.Index(i).Elem()
		if err := loadRelationshipsRecurse(exec, current, toLoad, singular, iter); err != nil {
			return err
		}
	}

	return nil
}

// loadRelationshipsRecurse is a helper function for taking a reflect.Value and
// Basically calls loadRelationships with: obj.R.EagerLoadedObj, and whether it's a string or slice
func loadRelationshipsRecurse(exec Executor, current string, toLoad []string, singular bool, obj reflect.Value) error {
	r := obj.FieldByName(relationshipStructName)
	if !r.IsValid() || r.IsNil() {
		return errors.Errorf("could not traverse into loaded %s relationship to load more things", current)
	}
	newObj := reflect.Indirect(r).FieldByName(current)
	singular = reflect.Indirect(newObj).Kind() == reflect.Struct
	if !singular {
		newObj = newObj.Addr()
	}
	return loadRelationships(exec, toLoad, newObj.Interface(), singular)
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

var mut sync.RWMutex
var mappingThings = map[string][]uint64{}

func bind(rows *sql.Rows, obj interface{}, structType, sliceType reflect.Type, singular bool, titleCases map[string]string) error {
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

	buf := strmangle.GetBuffer()

	buf.WriteString(structType.String())
	for _, s := range cols {
		buf.WriteString(s)
	}
	mapKey := buf.String()
	strmangle.PutBuffer(buf)

	mut.RLock()
	mapping, ok = mappingThings[mapKey]
	mut.RUnlock()

	if !ok {
		mapping, err = bindMapping(structType, titleCases, cols)
		if err != nil {
			return err
		}

		mut.Lock()
		mappingThings[mapKey] = mapping
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

func bindMapping(typ reflect.Type, titleCases map[string]string, cols []string) ([]uint64, error) {
	ptrs := make([]uint64, len(cols))
	mapping := makeStructMapping(typ, titleCases)

ColLoop:
	for i, c := range cols {
		names := strings.Split(c, ".")
		for j, n := range names {
			t, ok := titleCases[n]
			if ok {
				names[j] = t
				continue
			}
			names[j] = strmangle.TitleCase(n)
		}
		name := strings.Join(names, ".")

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
// for things on
func ptrsFromMapping(val reflect.Value, mapping []uint64) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m).Interface()
	}
	return ptrs
}

var sentinel = uint64(255)

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

func makeStructMapping(typ reflect.Type, titleCases map[string]string) map[string]uint64 {
	fieldMaps := make(map[string]uint64)
	makeStructMappingHelper(typ, "", 0, 0, fieldMaps, titleCases)
	return fieldMaps
}

func makeStructMappingHelper(typ reflect.Type, prefix string, current uint64, depth uint, fieldMaps map[string]uint64, titleCases map[string]string) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)

		tag, recurse := getBoilTag(f, titleCases)
		if len(tag) == 0 {
			tag = f.Name
		} else if tag[0] == '-' {
			continue
		}

		if len(prefix) != 0 {
			tag = fmt.Sprintf("%s.%s", prefix, tag)
		}

		if recurse {
			makeStructMappingHelper(f.Type, tag, current|uint64(i)<<depth, depth+8, fieldMaps, titleCases)
			continue
		}

		fieldMaps[tag] = current | (sentinel << (depth + 8)) | (uint64(i) << depth)
	}
}

func bin64(i uint64) string {
	str := strconv.FormatUint(i, 2)
	pad := 64 - len(str)
	if pad > 0 {
		str = strings.Repeat("0", pad) + str
	}

	var newStr string
	for i := 0; i < len(str); i += 8 {
		if i != 0 {
			newStr += " "
		}
		newStr += str[i : i+8]
	}

	return newStr
}

func findField(names []string, titleCases map[string]string, v reflect.Value) (interface{}, bool) {
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

	var name string
	var ok bool
	name, ok = titleCases[names[0]]
	if !ok {
		name = strmangle.TitleCase(names[0])
	}
	typ := v.Type()

	n := typ.NumField()
	for i := 0; i < n; i++ {
		f := typ.Field(i)
		fieldName, recurse := getBoilTag(f, titleCases)

		if fieldName == "-" {
			continue
		}

		if recurse {
			if fieldName == name {
				names = names[1:]
			}
			if ptr, ok := findField(names, titleCases, v.Field(i)); ok {
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

func getBoilTag(field reflect.StructField, titleCases map[string]string) (name string, recurse bool) {
	tag := field.Tag.Get("boil")
	name = field.Name

	if len(tag) == 0 {
		return name, false
	}

	var ok bool
	ind := strings.IndexByte(tag, ',')
	if ind == -1 {
		name, ok = titleCases[tag]
		if !ok {
			name = strmangle.TitleCase(tag)
		}
		return name, false
	} else if ind == 0 {
		return name, true
	}

	nameFragment := tag[:ind]
	name, ok = titleCases[nameFragment]
	if !ok {
		name = strmangle.TitleCase(nameFragment)
	}
	return name, true
}

// GetStructValues returns the values (as interface) of the matching columns in obj
func GetStructValues(obj interface{}, titleCases map[string]string, columns ...string) []interface{} {
	ret := make([]interface{}, len(columns))
	val := reflect.Indirect(reflect.ValueOf(obj))

	for i, c := range columns {
		var fieldName string
		if titleCases == nil {
			fieldName = strmangle.TitleCase(c)
		} else {
			fieldName = titleCases[c]
		}
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("unable to find field with name: %s\n%#v", fieldName, obj))
		}
		ret[i] = field.Interface()
	}

	return ret
}

// GetSliceValues returns the values (as interface) of the matching columns in obj.
func GetSliceValues(slice []interface{}, titleCases map[string]string, columns ...string) []interface{} {
	ret := make([]interface{}, len(slice)*len(columns))

	for i, obj := range slice {
		val := reflect.Indirect(reflect.ValueOf(obj))
		for j, c := range columns {
			var fieldName string
			if titleCases == nil {
				fieldName = strmangle.TitleCase(c)
			} else {
				fieldName = titleCases[c]
			}

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
func GetStructPointers(obj interface{}, titleCases map[string]string, columns ...string) []interface{} {
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
			var fieldName string
			if titleCases == nil {
				fieldName = strmangle.TitleCase(columns[i])
			} else {
				fieldName = titleCases[columns[i]]
			}

			return v.FieldByName(fieldName)
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
