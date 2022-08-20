package queries

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/friendsofgo/errors"
	"github.com/razor-1/sqlboiler/v4/boil"
	"github.com/volatiletech/strmangle"
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
	boilTag                = "boil"
	nullJoinTagValue       = "nulljoin"
)

// BindP executes the query and inserts the
// result into the passed in object pointer.
// It panics on error.
// Also see documentation for Bind() and Query.Bind()
func (q *Query) BindP(ctx context.Context, exec boil.Executor, obj interface{}) {
	if err := q.Bind(ctx, exec, obj); err != nil {
		panic(boil.WrapErr(err))
	}
}

// BindG executes the query and inserts
// the result into the passed in object pointer.
// It uses the global executor.
// Also see documentation for Bind() and Query.Bind()
func (q *Query) BindG(ctx context.Context, obj interface{}) error {
	return q.Bind(ctx, boil.GetDB(), obj)
}

// Bind inserts the rows into the passed in object pointer, because the caller
// owns the rows it is imperative to note that the caller MUST both close the
// rows and check for errors on the rows.
//
// If you neglect closing the rows your application may have a memory leak
// if the rows are not implicitly closed by iteration alone.
// If you neglect checking the rows.Err silent failures may occur in your
// application.
//
// Valid types to bind to are: *Struct, []*Struct, and []Struct. Keep in mind
// if you use []Struct that Bind will be doing copy-by-value as a method
// of keeping heap memory usage low which means if your Struct contains
// reference types/pointers you will see incorrect results, do not use
// []Struct with a Struct with reference types.
//
// Bind rules:
//   - Struct tags control bind, in the form of: `boil:"name,bind"`
//   - If "name" is omitted the sql column names that come back are TitleCased
//     and matched against the field name.
//   - If the "name" part of the struct tag is specified, the given name will
//     be used instead of the struct field name for binding.
//   - If the "name" of the struct tag is "-", this field will not be bound to.
//   - If the ",bind" option is specified on a struct field and that field
//     is a struct itself, it will be recursed into to look for fields for
//     binding.
//   - If one or more boil struct tags are duplicated and there are multiple
//     matching columns for those tags the behaviour of Bind will be undefined
//     for those fields with duplicated struct tags.
//
// Example usage:
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
//   models.Users(
//     qm.InnerJoin("users as friend on users.friend_id = friend.id")
//   ).Bind(&joinStruct)
//
// For custom objects that want to use eager loading, please see the
// loadRelationships function.
func Bind(rows *sql.Rows, obj interface{}) error {
	structType, _, singular, err := bindChecks(obj)
	if err != nil {
		return err
	}
	_, err = bind(rows, obj, structType, singular)
	return err
}

// Bind executes the query and inserts the
// result into the passed in object pointer.
//
// If Context is non-nil it will upgrade the
// Executor to a ContextExecutor and query with the passed context.
// If Context is non-nil, any eager loading that's done must also
// be using load* methods that support context as the first parameter.
//
// Also see documentation for Bind()
func (q *Query) Bind(ctx context.Context, exec boil.Executor, obj interface{}) error {
	structType, _, bkind, err := bindChecks(obj)
	if err != nil {
		return err
	}
	currentObj := obj

	if q.loadJoin {
		currentObj, err = eagerLoadJoin(q, obj, bkind)
		if err != nil {
			return err
		}
		// using a new bind type, so re-populate
		structType, _, bkind, err = bindChecks(currentObj)
		if err != nil {
			return err
		}
	}

	var rows *sql.Rows
	if ctx != nil {
		rows, err = q.QueryContext(ctx, exec.(boil.ContextExecutor))
	} else {
		rows, err = q.Query(exec)
	}
	if err != nil {
		return errors.Wrap(err, "bind failed to execute query")
	}
	var nullTables [][]string
	if nullTables, err = bind(rows, currentObj, structType, bkind); err != nil {
		if innerErr := rows.Close(); innerErr != nil {
			return errors.Wrapf(err, "error on rows.Close after bind error: %+v", innerErr)
		}

		return err
	}
	if err = rows.Close(); err != nil {
		return errors.Wrap(err, "failed to clean up rows in bind")
	}
	if err = rows.Err(); err != nil {
		return errors.Wrap(err, "error from rows in bind")
	}

	if len(q.load) != 0 && !q.loadJoin {
		return eagerLoad(ctx, exec, q.load, q.loadMods, obj, bkind)
	}

	if q.loadJoin {
		// we have to take the data that was bound into the load join struct and get it into
		// the expected struct type
		objVal := reflect.Indirect(reflect.ValueOf(obj))
		currentObjVal := reflect.Indirect(reflect.ValueOf(currentObj))
		if bkind == kindPtrSliceStruct {
			objType := reflect.TypeOf(obj).Elem().Elem().Elem()
			for i := 0; i < currentObjVal.Len(); i++ {
				newObj := reflect.Indirect(reflect.New(objType))
				err = processLoadJoinStruct(reflect.Indirect(currentObjVal.Index(i)), newObj, nullTables[i])
				if err != nil {
					return err
				}
				objVal.Set(reflect.Append(objVal, newObj.Addr()))
			}
		} else if bkind == kindStruct {
			// currentObj will have a field with the same type as the type of objVal
			err = processLoadJoinStruct(currentObjVal, objVal, nullTables[0])
			if err != nil {
				return err
			}
		} else {
			return errors.New("unsupported bindKind for loadJoin: " + string(bkind))
		}
	}

	return nil
}

// processLoadJoinStruct uses reflection to populate a struct of a type from a generated table
// from a struct of the load join type
func processLoadJoinStruct(currentObjVal, objVal reflect.Value, nullTables []string) error {
	numFields := currentObjVal.NumField()
	for i := 0; i < numFields; i++ {
		coField := currentObjVal.Field(i)
		if coField.Type() == objVal.Type() {
			objVal.Set(coField)
			break
		}
	}

	nilRTypes := make(map[string]bool)
	if len(nullTables) > 0 {
		// we need to see which types in the R struct should be nil pointers (nullable join was null)
		joinRespType := currentObjVal.Type()
		numFields = joinRespType.NumField()
		for i := 0; i < numFields; i++ {
			f := joinRespType.Field(i)
			table, _, nulljoin := parseBoilTag(f)
			if nulljoin && inSlice(table, nullTables) {
				nilRTypes[f.Type.Name()] = true
			}
		}
	}

	// take the loadjoin struct and map it into the original obj struct (populate *R)
	relationshipStruct := objVal.FieldByName(relationshipStructName)
	if !relationshipStruct.IsValid() {
		return errors.New("cannot find relationship struct")
	}
	rStruct := reflect.Indirect(reflect.New(relationshipStruct.Type().Elem()))
	for i := 0; i < rStruct.NumField(); i++ {
		f := rStruct.Field(i)
		if f.Kind() == reflect.Ptr {
			baseType := f.Type().Elem()
			fieldTypeName := f.Type().Elem().Name()
			if _, ok := nilRTypes[fieldTypeName]; ok {
				f.Set(reflect.Zero(f.Type()))
				continue
			}

			// currentObj is where the data is. we need to get the value of the same type from currentObj
			// and set it in our new rStruct
			numFields = currentObjVal.NumField()
			for j := 0; j < numFields; j++ {
				cf := currentObjVal.Field(j)
				if cf.Type() == baseType {
					f.Set(cf.Addr())
					break
				}
			}
		}
	}
	// set the R struct to be our new value
	relationshipStruct.Set(rStruct.Addr())
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

//find all items in slice that start with "needle."
func sliceTableIndex(slice []string, needle string) (map[int]bool, error) {
	indexes := make(map[int]bool)
	for idx, val := range slice {
		if strings.HasPrefix(val, needle+".") {
			indexes[idx] = true
		}
	}

	var err error
	if len(indexes) == 0 {
		err = errors.New("needle not found in slice")
	}

	return indexes, err
}

func parseBoilTag(field reflect.StructField) (table string, bind, nulljoin bool) {
	tag := field.Tag.Get(boilTag)
	if len(tag) == 0 {
		return
	}
	tagItems := strings.Split(tag, ",")
	if len(tagItems) < 2 {
		return
	}
	table = tagItems[0]
	bind = tagItems[1] == "bind"
	if len(tagItems) > 2 {
		nulljoin = tagItems[2] == nullJoinTagValue
	}

	return
}

func getNullableColumnIndexes(cols []string, structType reflect.Type) map[int]bool {
	indexes := make(map[int]bool)

	numFields := structType.NumField()
	for i := 0; i < numFields; i++ {
		f := structType.Field(i)
		table, _, nulljoin := parseBoilTag(f)
		if nulljoin && len(table) > 0 {
			// this is a hit. add all the indexes to the map we will return
			colIndexes, err := sliceTableIndex(cols, table)
			if err == nil {
				for colIndex, _ := range colIndexes {
					indexes[colIndex] = true
				}
			}
		}
	}

	return indexes
}

func bind(rows *sql.Rows, obj interface{}, structType reflect.Type, bkind bindKind) ([][]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "bind failed to get column names")
	}

	var ptrSlice reflect.Value
	switch bkind {
	case kindSliceStruct, kindPtrSliceStruct:
		ptrSlice = reflect.Indirect(reflect.ValueOf(obj))
	}

	mapping, err := getMappingCache(structType).mapping(cols)
	if err != nil {
		return nil, err
	}
	/*
		var strMapping map[string]uint64
		var sok bool
		var mapping []uint64
		var ok bool

		typKey := makeTypeKey(structType)
		colsKey := makeColsKey(structType, cols)



		mut.RLock()
		mapping, ok = bindingMaps[colsKey]
		if !ok {
			if strMapping, sok = structMaps[typKey]; !sok {
				strMapping = MakeStructMapping(structType)
			}
		}
		mut.RUnlock()

		if !ok {
			mapping, err = BindMapping(structType, strMapping, cols)
			if err != nil {
				return nil, err
			}

			mut.Lock()
			if !sok {
				structMaps[typKey] = strMapping
			}
			bindingMaps[colsKey] = mapping
			mut.Unlock()
		}
	*/

	//see if this has the "nulljoin" option on the tag
	nullableIndexes := getNullableColumnIndexes(cols, structType)

	foundOne := false
	nullTables := make([][]string, 0)

Rows:
	for rows.Next() {
		foundOne = true
		var newStruct reflect.Value
		var pointers []interface{}
		nullCols := make([]string, 0)
		var oneStruct reflect.Value
		if bkind == kindSliceStruct {
			oneStruct = reflect.Indirect(reflect.New(structType))
		}

		switch bkind {
		case kindStruct:
			pointers = PtrsFromMapping(reflect.Indirect(reflect.ValueOf(obj)), mapping)
		case kindSliceStruct:
			pointers = PtrsFromMapping(oneStruct, mapping)
		case kindPtrSliceStruct:
			newStruct = reflect.New(structType)
			pointers = PtrsFromMapping(reflect.Indirect(newStruct), mapping)
		}

		if len(nullableIndexes) > 0 {
			nullablePointers := make([]interface{}, len(pointers))
			for i := 0; i < len(nullablePointers); i++ {
				if _, ok := nullableIndexes[i]; ok {
					nullablePointers[i] = new(sql.RawBytes)
				} else {
					nullablePointers[i] = pointers[i]
				}
			}
			// preflight Scan to check for nulls
			if err := rows.Scan(nullablePointers...); err != nil {
				return nil, errors.Wrap(err, "failed to bind nullable pointers to obj")
			}

			// see if anything that could be null actually was
			for nullIndex, _ := range nullableIndexes {
				rb := nullablePointers[nullIndex].(*sql.RawBytes)
				// it's not a nil pointer, but a pointer to an empty byte slice:
				if *rb == nil {
					// change pointers at this index to allow a null value in the real Scan below
					pointers[nullIndex] = new(sql.RawBytes)
					nullCols = append(nullCols, cols[nullIndex])
				}
			}
		}

		//for each row, store the result of findNullTables into nullTables
		//this is used later to set the corresponding .R pointer to nil
		nullTables = append(nullTables, findNullTables(cols, nullCols))

		if err := rows.Scan(pointers...); err != nil {
			return nil, errors.Wrap(err, "failed to bind pointers to obj")
		}

		switch bkind {
		case kindStruct:
			break Rows
		case kindSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, oneStruct))
		case kindPtrSliceStruct:
			ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
		}
	}

	if bkind == kindStruct && !foundOne {
		return nil, sql.ErrNoRows
	}

	return nullTables, nil
}

// findNullTables: if every value for a given table was nil, then add it to our list of null tables
// this is used to set table.R.ForeignTable = nil instead of having it populate with zero values
// both string slices passed to this function will contain entries like "table.column"
func findNullTables(allColumns, nullColumns []string) []string {
	nullTables := make([]string, 0)

	allColumnsMap := tableColumnMap(allColumns)
	nullColumnsMap := tableColumnMap(nullColumns)

	for table, columns := range nullColumnsMap {
		// if all the columns in the null map are the same as all the columns in the table, then it's a null table
		if len(columns) == len(allColumnsMap[table]) {
			equal := true
			for _, nullCol := range columns {
				if !inSlice(nullCol, allColumnsMap[table]) {
					equal = false
				}
			}
			if equal {
				nullTables = append(nullTables, table)
			}
		}
	}

	return nullTables
}

// see if needle is in haystack
func inSlice(needle string, haystack []string) bool {
	for _, haystackEntry := range haystack {
		if needle == haystackEntry {
			return true
		}
	}

	return false
}

// make a map of table to columns
func tableColumnMap(columns []string) map[string][]string {
	tables := make(map[string][]string)
	for _, col := range columns {
		colsplit := strings.Split(col, ".")
		if len(colsplit) != 2 {
			continue
		}
		if _, ok := tables[colsplit[0]]; !ok {
			tables[colsplit[0]] = []string{colsplit[1]}
		} else {
			tables[colsplit[0]] = append(tables[colsplit[0]], colsplit[1])
		}
	}

	return tables
}

// BindMapping creates a mapping that helps look up the pointer for the
// column given.
func BindMapping(typ reflect.Type, mapping map[string]uint64, cols []string) ([]uint64, error) {
	ptrs := make([]uint64, len(cols))

ColLoop:
	for i, c := range cols {
		ptrMap, ok := mapping[c]
		if ok {
			ptrs[i] = ptrMap
			continue
		}

		suffix := "." + c
		for maybeMatch, mapping := range mapping {
			if strings.HasSuffix(maybeMatch, suffix) {
				ptrs[i] = mapping
				continue ColLoop
			}
		}
		// if c doesn't exist in the model, the pointer will be the zero value in the ptrs array and it's value will be thrown away
		continue
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
	if mapping == 0 {
		var ignored interface{}
		return reflect.ValueOf(&ignored)
	}
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
			tag = unTitleCase(f.Name)
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
	tag := field.Tag.Get(boilTag)

	if len(tag) == 0 {
		return "", false
	}

	ind := strings.IndexByte(tag, ',')
	if ind == -1 {
		return tag, false
	} else if ind == 0 {
		return "", true
	}

	nameFragment := tag[:ind]
	return nameFragment, true
}

var (
	mappingCachesMu sync.Mutex
	mappingCaches   = make(map[reflect.Type]*mappingCache)
)

func getMappingCache(typ reflect.Type) *mappingCache {
	mappingCachesMu.Lock()
	defer mappingCachesMu.Unlock()

	cache := mappingCaches[typ]
	if cache != nil {
		return cache
	}

	cache = newMappingCache(typ)
	mappingCaches[typ] = cache

	return cache
}

type mappingCache struct {
	typ reflect.Type

	mu          sync.Mutex
	structMap   map[string]uint64
	colMappings map[string][]uint64
}

func newMappingCache(typ reflect.Type) *mappingCache {
	return &mappingCache{
		typ:         typ,
		structMap:   MakeStructMapping(typ),
		colMappings: make(map[string][]uint64),
	}
}

func (b *mappingCache) mapping(cols []string) ([]uint64, error) {
	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	for _, s := range cols {
		buf.WriteString(s)
		buf.WriteByte(0)
	}

	key := buf.Bytes()

	b.mu.Lock()
	defer b.mu.Unlock()

	mapping := b.colMappings[string(key)]
	if mapping != nil {
		return mapping, nil
	}

	mapping, err := BindMapping(b.typ, b.structMap, cols)
	if err != nil {
		return nil, err
	}

	b.colMappings[string(key)] = mapping

	return mapping, nil
}

func makeTypeKey(typ reflect.Type) string {
	buf := strmangle.GetBuffer()
	buf.WriteString(typ.String())
	for i, n := 0, typ.NumField(); i < n; i++ {
		field := typ.Field(i)
		buf.WriteString(field.Name)
		buf.WriteString(field.Type.String())
	}
	hash := buf.String()
	strmangle.PutBuffer(buf)

	return hash
}

func makeColsKey(typ reflect.Type, cols []string) string {
	buf := strmangle.GetBuffer()
	buf.WriteString(typ.String())
	for _, s := range cols {
		buf.WriteString(s)
	}
	hash := buf.String()
	strmangle.PutBuffer(buf)

	return hash
}

// Equal is different to reflect.DeepEqual in that it's both less efficient
// less magical, and dosen't concern itself with a wide variety of types that could
// be present but it does use the driver.Valuer interface since many types that will
// go through database things will use these.
//
// We're focused on basic types + []byte. Since we're really only interested in things
// that are typically used for primary keys in a database.
//
// Choosing not to use the DefaultParameterConverter here because sqlboiler doesn't generate
// pointer columns.
func Equal(a, b interface{}) bool {
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	// Here we make a fast-path for bytes, because it's the most likely thing
	// this method will be called with.
	if ab, ok := a.([]byte); ok {
		if bb, ok := b.([]byte); ok {
			return bytes.Equal(ab, bb)
		}
	}

	var err error
	// If either is a sql.Scanner, pull the primitive value out before we get into type checking
	// since we can't compare complex types anyway.
	if v, ok := a.(driver.Valuer); ok {
		a, err = v.Value()
		if err != nil {
			panic(fmt.Sprintf("while comparing values, although 'a' implemented driver.Valuer, an error occured when calling it: %+v", err))
		}
	}
	if v, ok := b.(driver.Valuer); ok {
		b, err = v.Value()
		if err != nil {
			panic(fmt.Sprintf("while comparing values, although 'b' implemented driver.Valuer, an error occured when calling it: %+v", err))
		}
	}

	// Do nil checks again, since a Null type could have returned nil
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	// If either is string and another is numeric, try to parse string as numeric
	if as, ok := a.(string); ok && isNumeric(b) {
		a = parseNumeric(as, reflect.TypeOf(b))
	}
	if bs, ok := b.(string); ok && isNumeric(a) {
		b = parseNumeric(bs, reflect.TypeOf(a))
	}

	a = upgradeNumericTypes(a)
	b = upgradeNumericTypes(b)

	if at, bt := reflect.TypeOf(a), reflect.TypeOf(b); at != bt {
		panic(fmt.Sprintf("primitive type of a (%s) was not the same primitive type as b (%s)", at.String(), bt.String()))
	}

	switch t := a.(type) {
	case int64, float64, bool, string:
		return a == b
	case []byte:
		return bytes.Equal(t, b.([]byte))
	case time.Time:
		return t.Equal(b.(time.Time))
	}

	return false
}

// isNumeric tests if i is a numeric value.
func isNumeric(i interface{}) bool {
	switch i.(type) {
	case int,
		int8,
		int16,
		int32,
		int64,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		float32,
		float64:
		return true
	}
	return false
}

// parseNumeric tries to parse s as t.
// t must be a numeric type.
func parseNumeric(s string, t reflect.Type) interface{} {
	var (
		res interface{}
		err error
	)
	switch t.Kind() {
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		res, err = strconv.ParseInt(s, 0, t.Bits())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		res, err = strconv.ParseUint(s, 0, t.Bits())
	case reflect.Float32,
		reflect.Float64:
		res, err = strconv.ParseFloat(s, t.Bits())
	}
	if err != nil {
		panic(fmt.Sprintf("tries to parse %q as %s but got error: %+v", s, t.String(), err))
	}
	return res
}

// Assign assigns a value to another using reflection.
// Dst must be a pointer.
func Assign(dst, src interface{}) {
	// Fast path for []byte since it's one of the
	// most frequent other "ids" we'll be assigning.
	if db, ok := dst.(*[]byte); ok {
		if sb, ok := src.([]byte); ok {
			*db = make([]byte, len(sb))
			copy(*db, sb)
			return
		}
	}

	scan, isDstScanner := dst.(sql.Scanner)
	val, isSrcValuer := src.(driver.Valuer)

	switch {
	case isDstScanner && isSrcValuer:
		val, err := val.Value()
		if err != nil {
			panic(fmt.Sprintf("tried to call value on %T but got err: %+v", src, err))
		}

		err = scan.Scan(val)
		if err != nil {
			panic(fmt.Sprintf("tried to call Scan on %T with %#v but got err: %+v", dst, val, err))
		}

	case isDstScanner && !isSrcValuer:
		// Compress any lower width integer types
		src = upgradeNumericTypes(src)

		if err := scan.Scan(src); err != nil {
			panic(fmt.Sprintf("tried to call Scan on %T with %#v but got err: %+v", dst, src, err))
		}

	case !isDstScanner && isSrcValuer:
		val, err := val.Value()
		if err != nil {
			panic(fmt.Sprintf("tried to call value on %T but got err: %+v", src, err))
		}

		assignValue(dst, val)

	default:
		// We should always be comparing primitives with each other with == in templates
		// so this method should never be called for say: string, string, or int, int
		panic("this case should have been handled by something other than this method")
	}
}

func upgradeNumericTypes(i interface{}) interface{} {
	switch t := i.(type) {
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case float32:
		return float64(t)
	default:
		return i
	}
}

// This whole function makes assumptions that whatever type
// dst is, will be compatible with whatever came out of the Valuer.
// We handle the types that driver.Value could possibly be.
func assignValue(dst interface{}, val driver.Value) {
	dstType := reflect.TypeOf(dst).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	if val == nil {
		dstVal.Set(reflect.Zero(dstType))
		return
	}

	v := reflect.ValueOf(val)

	switch dstType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dstVal.SetInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dstVal.SetInt(int64(v.Uint()))
	case reflect.Bool:
		dstVal.SetBool(v.Bool())
	case reflect.String:
		dstVal.SetString(v.String())
	case reflect.Float32, reflect.Float64:
		dstVal.SetFloat(v.Float())
	case reflect.Slice:
		// Assume []byte
		db, sb := dst.(*[]byte), val.([]byte)
		*db = make([]byte, len(sb))
		copy(*db, sb)
	case reflect.Struct:
		// Assume time.Time
		dstVal.Set(v)
	}
}

// MustTime retrieves a time value from a valuer.
func MustTime(val driver.Valuer) time.Time {
	v, err := val.Value()
	if err != nil {
		panic(fmt.Sprintf("attempted to call value on %T to get time but got an error: %+v", val, err))
	}

	if v == nil {
		return time.Time{}
	}

	return v.(time.Time)
}

// IsValuerNil returns true if the valuer's value is null.
func IsValuerNil(val driver.Valuer) bool {
	v, err := val.Value()
	if err != nil {
		panic(fmt.Sprintf("attempted to call value on %T but got an error: %+v", val, err))
	}

	return v == nil
}

// IsNil is a more generic version of IsValuerNil, will check to make sure it's
// not a valuer first.
func IsNil(val interface{}) bool {
	if val == nil {
		return true
	}

	valuer, ok := val.(driver.Valuer)
	if ok {
		return IsValuerNil(valuer)
	}

	return reflect.ValueOf(val).IsNil()
}

// SetScanner attempts to set a scannable value on a scanner.
func SetScanner(scanner sql.Scanner, val driver.Value) {
	if err := scanner.Scan(val); err != nil {
		panic(fmt.Sprintf("attempted to call Scan on %T with %#v but got an error: %+v", scanner, val, err))
	}
}

// These are sorted by size so that the biggest thing
// gets replaced first (think guid/id). This list is copied
// from strmangle.uppercaseWords and should hopefully be kept
// in sync.
var specialWordReplacer = strings.NewReplacer(
	"ASCII", "Ascii",
	"GUID", "Guid",
	"JSON", "Json",
	"UUID", "Uuid",
	"UTF8", "Utf8",
	"ACL", "Acl",
	"API", "Api",
	"CPU", "Cpu",
	"EOF", "Eof",
	"RAM", "Ram",
	"SLA", "Sla",
	"UDP", "Udp",
	"UID", "Uid",
	"URI", "Uri",
	"URL", "Url",
	"ID", "Id",
	"IP", "Ip",
	"UI", "Ui",
)

// unTitleCase attempts to undo a title-cased string.
//
// DO NOT USE THIS METHOD IF YOU CAN AVOID IT
//
// Normally this would be easy but we have to deal with uppercased words
// of varying lengths. We almost never use this function so it
// can be as badly performing as we want. If people don't want to incur
// it's cost they should be able to use the `boil` struct tag to avoid it.
//
// We did not put this in strmangle because we don't want it being part
// of any public API as it's loaded with corner cases and sad performance.
func unTitleCase(n string) string {
	if len(n) == 0 {
		return ""
	}

	// Make our words no longer special case
	n = specialWordReplacer.Replace(n)

	buf := strmangle.GetBuffer()

	first := true

	writeIt := func(s string) {
		if first {
			first = false
		} else {
			buf.WriteByte('_')
		}
		buf.WriteString(strings.ToLower(s))
	}

	lastUp := true
	start := 0
	for i, r := range n {
		currentUp := unicode.IsUpper(r)
		isDigit := unicode.IsDigit(r)

		if !isDigit && !lastUp && currentUp {
			fragment := n[start:i]
			writeIt(fragment)
			start = i
		}

		if !isDigit && lastUp && !currentUp && i-1-start > 1 {
			fragment := n[start : i-1]
			writeIt(fragment)
			start = i - 1
		}

		lastUp = currentUp
	}

	remaining := n[start:]
	if len(remaining) > 0 {
		writeIt(remaining)
	}

	ret := buf.String()
	strmangle.PutBuffer(buf)
	return ret
}
