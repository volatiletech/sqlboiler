package queries

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/strmangle"
)

var (
	bindAccepts = []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct}

	mut         sync.RWMutex
	bindingMaps = make(map[string][]uint64)
	structMaps  = make(map[string]map[string]uint64)
)

// Identifies what kind of object we're binding to
type bindKind int

const (
	kindStruct bindKind = iota
	kindSliceStruct
	kindPtrSliceStruct
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
		panic(boil.WrapErr(err))
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
	structType, sliceType, bkind, err := bindChecks(obj)
	if err != nil {
		return err
	}

	rows, err := q.Query()
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}
	defer rows.Close()
	if res := bind(rows, obj, structType, sliceType, bkind); res != nil {
		return res
	}

	if len(q.load) != 0 {
		return eagerLoad(q.executor, q.load, obj, bkind)
	}

	return nil
}

// bindChecks resolves information about the bind target, and errors if it's not an object
// we can bind to.
func bindChecks(obj interface{}) (structType reflect.Type, sliceType reflect.Type, bkind bindKind, err error) {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

	setErr := func() {
		err = errors.Errorf("obj type should be *Type, *[]Type, or *[]*Type but was %q", reflect.TypeOf(obj).String())
	}

	for i := 0; ; i++ {
		switch i {
		case 0:
			if kind != reflect.Ptr {
				setErr()
				return
			}
		case 1:
			switch kind {
			case reflect.Struct:
				structType = typ
				bkind = kindStruct
				return
			case reflect.Slice:
				sliceType = typ
			default:
				setErr()
				return
			}
		case 2:
			switch kind {
			case reflect.Struct:
				structType = typ
				bkind = kindSliceStruct
				return
			case reflect.Ptr:
			default:
				setErr()
				return
			}
		case 3:
			if kind != reflect.Struct {
				setErr()
				return
			}
			structType = typ
			bkind = kindPtrSliceStruct
			return
		}

		typ = typ.Elem()
		kind = typ.Kind()
	}
}

func bind(rows *sql.Rows, obj interface{}, structType, sliceType reflect.Type, bkind bindKind) error {
	cols, err := rows.Columns()
	if err != nil {
		return errors.Wrap(err, "bind failed to get column names")
	}

	var ptrSlice reflect.Value
	switch bkind {
	case kindSliceStruct, kindPtrSliceStruct:
		ptrSlice = reflect.Indirect(reflect.ValueOf(obj))
	}

	var strMapping map[string]uint64
	var sok bool
	var mapping []uint64
	var ok bool

	typStr := structType.String()

	mapKey := makeCacheKey(typStr, cols)
	mut.RLock()
	mapping, ok = bindingMaps[mapKey]
	if !ok {
		if strMapping, sok = structMaps[typStr]; !sok {
			strMapping = MakeStructMapping(structType)
		}
	}
	mut.RUnlock()

	if !ok {
		mapping, err = BindMapping(structType, strMapping, cols)
		if err != nil {
			return err
		}

		mut.Lock()
		if !sok {
			structMaps[typStr] = strMapping
		}
		bindingMaps[mapKey] = mapping
		mut.Unlock()
	}

	var oneStruct reflect.Value
	if bkind == kindSliceStruct {
		oneStruct = reflect.Indirect(reflect.New(structType))
	}

	foundOne := false
	for rows.Next() {
		foundOne = true
		var newStruct reflect.Value
		var pointers []interface{}

		switch bkind {
		case kindStruct:
			pointers = PtrsFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mapping)
		case kindSliceStruct:
			pointers = PtrsFromMapping(oneStruct, mapping)
		case kindPtrSliceStruct:
			newStruct = reflect.New(structType)
			pointers = PtrsFromMapping(reflect.Indirect(newStruct), mapping)
		}
		if err != nil {
			return err
		}

		if err := rows.Scan(pointers...); err != nil {
			return errors.Wrap(err, "failed to bind pointers to obj")
		}

		switch bkind {
		case kindSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, oneStruct))
		case kindPtrSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
		}
	}

	if bkind == kindStruct && !foundOne {
		return sql.ErrNoRows
	}

	return nil
}

// BindMapping creates a mapping that helps look up the pointer for the
// column given.
func BindMapping(typ reflect.Type, mapping map[string]uint64, cols []string) ([]uint64, error) {
	ptrs := make([]uint64, len(cols))

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

// PtrsFromMapping expects to be passed an addressable struct and a mapping
// of where to find things. It pulls the pointers out referred to by the mapping.
func PtrsFromMapping(val reflect.Value, mapping []uint64) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m, true).Interface()
	}
	return ptrs
}

// ValuesFromMapping expects to be passed an addressable struct and a mapping
// of where to find things. It pulls the pointers out referred to by the mapping.
func ValuesFromMapping(val reflect.Value, mapping []uint64) []interface{} {
	ptrs := make([]interface{}, len(mapping))
	for i, m := range mapping {
		ptrs[i] = ptrFromMapping(val, m, false).Interface()
	}
	return ptrs
}

// ptrFromMapping expects to be passed an addressable struct that it's looking
// for things on.
func ptrFromMapping(val reflect.Value, mapping uint64, addressOf bool) reflect.Value {
	for i := 0; i < 8; i++ {
		v := (mapping >> uint(i*8)) & sentinel

		if v == sentinel {
			if addressOf && val.Kind() != reflect.Ptr {
				return val.Addr()
			} else if !addressOf && val.Kind() == reflect.Ptr {
				return reflect.Indirect(val)
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

// MakeStructMapping creates a map of the struct to be able to quickly look
// up its pointers and values by name.
func MakeStructMapping(typ reflect.Type) map[string]uint64 {
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
