package boil

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/strmangle"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}
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

	if len(q.load) == 0 {
		return nil
	}

	return q.loadRelationships(obj, singular)
}

// loadRelationships dynamically calls the template generated eager load
// functions of the form:
//
//   func (t *TableLoaded) LoadRelationshipName(exec Executor, singular bool, obj interface{})
//
// The arguments to this function are:
//   - t is not considered here, and is always passed nil. The function exists on a loaded
//     struct to avoid a circular dependency with boil, and the receiver is ignored.
//   - exec is used to perform additional queries that might be required for loading the relationships.
//   - singular is passed in to identify whether or not this was a single object
//     or a slice that must be loaded into.
//   - obj is the object or slice of objects, always of the type *obj or *[]*obj as per bind.
func (q *Query) loadRelationships(obj interface{}, singular bool) error {
	typ := reflect.TypeOf(obj).Elem()
	if !singular {
		typ = typ.Elem()
	}

	rel, found := typ.FieldByName("Loaded")
	// If the users object has no loaded struct, it must be
	// a custom object and we should not attempt to load any relationships.
	if !found {
		return errors.New("load query mod was used but bound struct contained no Loaded field")
	}

	for _, relationship := range q.load {
		// Attempt to find the LoadRelationshipName function
		loadMethod, found := rel.Type.MethodByName("Load" + relationship)
		if !found {
			return errors.Errorf("could not find Load%s method for eager loading", relationship)
		}

		execArg := reflect.ValueOf(q.executor)
		if !execArg.IsValid() {
			execArg = reflect.ValueOf((*sql.DB)(nil))
		}

		methodArgs := []reflect.Value{
			reflect.Indirect(reflect.New(rel.Type)),
			execArg,
			reflect.ValueOf(singular),
			reflect.ValueOf(obj),
		}

		resp := loadMethod.Func.Call(methodArgs)
		if resp[0].Interface() != nil {
			return resp[0].Interface().(error)
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

	foundOne := false
	for rows.Next() {
		foundOne = true
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

	if singular && !foundOne {
		return sql.ErrNoRows
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
		if !field.IsValid() {
			panic(fmt.Sprintf("unable to find field with name: %s\n%#v", strmangle.TitleCase(c), obj))
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
			field := val.FieldByName(strmangle.TitleCase(c))
			if !field.IsValid() {
				panic(fmt.Sprintf("unable to find field with name: %s\n%#v", strmangle.TitleCase(c), obj))
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
