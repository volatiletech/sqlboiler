package boil

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/nullbio/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v4"
)

var (
	typeNullFloat32 = reflect.TypeOf(null.Float32{})
	typeNullFloat64 = reflect.TypeOf(null.Float64{})
	typeNullInt     = reflect.TypeOf(null.Int{})
	typeNullInt8    = reflect.TypeOf(null.Int8{})
	typeNullInt16   = reflect.TypeOf(null.Int16{})
	typeNullInt32   = reflect.TypeOf(null.Int32{})
	typeNullInt64   = reflect.TypeOf(null.Int64{})
	typeNullUint    = reflect.TypeOf(null.Uint{})
	typeNullUint8   = reflect.TypeOf(null.Uint8{})
	typeNullUint16  = reflect.TypeOf(null.Uint16{})
	typeNullUint32  = reflect.TypeOf(null.Uint32{})
	typeNullUint64  = reflect.TypeOf(null.Uint64{})
	typeNullString  = reflect.TypeOf(null.String{})
	typeNullBool    = reflect.TypeOf(null.Bool{})
	typeNullTime    = reflect.TypeOf(null.Time{})
	typeTime        = reflect.TypeOf(time.Time{})
)

// Bind executes the query and inserts the
// result into the passed in object pointer
func (q *Query) Bind(obj interface{}) error {
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()

	if kind != reflect.Ptr {
		return fmt.Errorf("Bind not given a pointer to a slice or struct: %s", typ.String())
	}

	typ = typ.Elem()
	kind = typ.Kind()

	if kind == reflect.Struct {
		row := ExecQueryOne(q)
		err := BindOne(row, q.selectCols, obj)
		if err != nil {
			return fmt.Errorf("Failed to execute Bind query for %s: %s", q.table, err)
		}
	} else if kind == reflect.Slice {
		rows, err := ExecQueryAll(q)
		if err != nil {
			return fmt.Errorf("Failed to execute Bind query for %s: %s", q.table, err)
		}
		err = BindAll(rows, q.selectCols, obj)
		if err != nil {
			return fmt.Errorf("Failed to Bind results to object provided for %s: %s", q.table, err)
		}
	} else {
		return fmt.Errorf("Bind given a pointer to a non-slice or non-struct: %s", typ.String())
	}

	return nil
}

// BindOne inserts the returned row columns into the
// passed in object pointer
func BindOne(row *sql.Row, selectCols []string, obj interface{}) error {
	kind := reflect.ValueOf(obj).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("BindOne given a non-pointer type")
	}

	pointers := GetStructPointers(obj, selectCols...)
	if err := row.Scan(pointers...); err != nil {
		return fmt.Errorf("Unable to scan into pointers: %s", err)
	}

	return nil
}

// BindAll inserts the returned rows columns into the
// passed in slice of object pointers
func BindAll(rows *sql.Rows, selectCols []string, obj interface{}) error {
	ptrSlice := reflect.ValueOf(obj)
	typ := ptrSlice.Type()
	ptrSlice = ptrSlice.Elem()
	kind := typ.Kind()

	var structTyp reflect.Type

	for i, exp := range []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct} {
		if i != 0 {
			typ = typ.Elem()
			kind = typ.Kind()
		}

		if kind != exp {
			return fmt.Errorf("[%d] BindAll object type should be *[]*Type but was: %s", i, ptrSlice.Type().String())
		}

		if kind == reflect.Struct {
			structTyp = typ
		}
	}

	for rows.Next() {
		newStruct := reflect.New(structTyp)
		pointers := GetStructPointers(newStruct.Interface(), selectCols...)
		if err := rows.Scan(pointers...); err != nil {
			return fmt.Errorf("Unable to scan into pointers: %s", err)
		}

		ptrSlice.Set(reflect.Append(ptrSlice, newStruct))
	}

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
	var ret []interface{}

	if len(columns) == 0 {
		fieldsLen := val.NumField()
		ret = make([]interface{}, fieldsLen)
		for i := 0; i < fieldsLen; i++ {
			ret[i] = val.Field(i).Addr().Interface()
		}
		return ret
	}

	ret = make([]interface{}, len(columns))

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

// RandomizeSlice takes a pointer to a slice of pointers to objects
// and fills the pointed to objects with random data.
// It will ignore the fields in the blacklist.
func RandomizeSlice(obj interface{}, blacklist ...string) error {
	ptrSlice := reflect.ValueOf(obj)
	typ := ptrSlice.Type()
	ptrSlice = ptrSlice.Elem()
	kind := typ.Kind()

	var structTyp reflect.Type

	for i, exp := range []reflect.Kind{reflect.Ptr, reflect.Slice, reflect.Ptr, reflect.Struct} {
		if i != 0 {
			typ = typ.Elem()
			kind = typ.Kind()
		}

		if kind != exp {
			return fmt.Errorf("[%d] RandomizeSlice object type should be *[]*Type but was: %s", i, ptrSlice.Type().String())
		}

		if kind == reflect.Struct {
			structTyp = typ
		}
	}

	for i := 0; i < ptrSlice.Len(); i++ {
		o := ptrSlice.Index(i)
		o.Set(reflect.New(structTyp))
		if err := RandomizeStruct(o.Interface(), blacklist...); err != nil {
			return err
		}
	}

	return nil
}

// RandomizeStruct takes an object and fills it with random data.
// It will ignore the fields in the blacklist.
func RandomizeStruct(str interface{}, blacklist ...string) error {
	// Don't modify blacklist
	copyBlacklist := make([]string, len(blacklist))
	copy(copyBlacklist, blacklist)
	blacklist = copyBlacklist

	sort.Strings(blacklist)

	// Check if it's pointer
	value := reflect.ValueOf(str)
	kind := value.Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("Outer element should be a pointer, given a non-pointer: %T", str)
	}

	// Check if it's a struct
	value = value.Elem()
	kind = value.Kind()
	if kind != reflect.Struct {
		return fmt.Errorf("Inner element should be a struct, given a non-struct: %T", str)
	}

	typ := value.Type()
	nFields := value.NumField()

	// Iterate through fields, randomizing
	for i := 0; i < nFields; i++ {
		fieldVal := value.Field(i)
		fieldTyp := typ.Field(i)

		found := sort.Search(len(blacklist), func(i int) bool {
			return blacklist[i] == fieldTyp.Name
		})
		if found != len(blacklist) {
			continue
		}

		if err := randomizeField(fieldVal); err != nil {
			return err
		}
	}

	return nil
}

// randDate generates a random time.Time between 1850 and 2050.
// Only the Day/Month/Year columns are set so that Dates and DateTimes do
// not cause mismatches in the test data comparisons.
func randDate() time.Time {
	t := time.Date(
		1850+rand.Intn(200),
		time.Month(1+rand.Intn(12)),
		1+rand.Intn(25),
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return t
}

func randomizeField(field reflect.Value) error {
	kind := field.Kind()
	typ := field.Type()

	var newVal interface{}

	if kind == reflect.Struct {
		b := rand.Intn(2) == 1
		switch typ {
		case typeNullBool:
			if b {
				newVal = null.NewBool(rand.Intn(2) == 1, b)
			} else {
				newVal = null.NewBool(false, false)
			}
		case typeNullString:
			if b {
				newVal = null.NewString(randStr(1), b)
			} else {
				newVal = null.NewString("", false)
			}
		case typeNullTime:
			if b {
				newVal = null.NewTime(randDate(), b)
			} else {
				newVal = null.NewTime(time.Time{}, false)
			}
		case typeTime:
			newVal = randDate()
		case typeNullFloat32:
			if b {
				newVal = null.NewFloat32(float32(rand.Intn(9))/10.0+float32(rand.Intn(9)), b)
			} else {
				newVal = null.NewFloat32(0.0, false)
			}
		case typeNullFloat64:
			if b {
				newVal = null.NewFloat64(float64(rand.Intn(9))/10.0+float64(rand.Intn(9)), b)
			} else {
				newVal = null.NewFloat64(0.0, false)
			}
		case typeNullInt:
			if b {
				newVal = null.NewInt(rand.Int(), b)
			} else {
				newVal = null.NewInt(0, false)
			}
		case typeNullInt8:
			if b {
				newVal = null.NewInt8(int8(rand.Intn(int(math.MaxInt8))), b)
			} else {
				newVal = null.NewInt8(0, false)
			}
		case typeNullInt16:
			if b {
				newVal = null.NewInt16(int16(rand.Intn(int(math.MaxInt16))), b)
			} else {
				newVal = null.NewInt16(0, false)
			}
		case typeNullInt32:
			if b {
				newVal = null.NewInt32(rand.Int31(), b)
			} else {
				newVal = null.NewInt32(0, false)
			}
		case typeNullInt64:
			if b {
				newVal = null.NewInt64(rand.Int63(), b)
			} else {
				newVal = null.NewInt64(0, false)
			}
		case typeNullUint:
			if b {
				newVal = null.NewUint(uint(rand.Int()), b)
			} else {
				newVal = null.NewUint(0, false)
			}
		case typeNullUint8:
			if b {
				newVal = null.NewUint8(uint8(rand.Intn(int(math.MaxInt8))), b)
			} else {
				newVal = null.NewUint8(0, false)
			}
		case typeNullUint16:
			if b {
				newVal = null.NewUint16(uint16(rand.Intn(int(math.MaxInt16))), b)
			} else {
				newVal = null.NewUint16(0, false)
			}
		case typeNullUint32:
			if b {
				newVal = null.NewUint32(uint32(rand.Int31()), b)
			} else {
				newVal = null.NewUint32(0, false)
			}
		case typeNullUint64:
			if b {
				newVal = null.NewUint64(uint64(rand.Int63()), b)
			} else {
				newVal = null.NewUint64(0, false)
			}
		}
	} else {
		switch kind {
		case reflect.Float32:
			newVal = float32(rand.Intn(9))/10.0 + float32(rand.Intn(9))
		case reflect.Float64:
			newVal = float64(rand.Intn(9))/10.0 + float64(rand.Intn(9))
		case reflect.Int:
			newVal = rand.Int()
		case reflect.Int8:
			newVal = int8(rand.Intn(int(math.MaxInt8)))
		case reflect.Int16:
			newVal = int16(rand.Intn(int(math.MaxInt16)))
		case reflect.Int32:
			newVal = rand.Int31()
		case reflect.Int64:
			newVal = rand.Int63()
		case reflect.Uint:
			newVal = uint(rand.Int())
		case reflect.Uint8:
			newVal = uint8(rand.Intn(int(math.MaxInt8)))
		case reflect.Uint16:
			newVal = uint16(rand.Intn(int(math.MaxInt16)))
		case reflect.Uint32:
			newVal = uint32(rand.Int31())
		case reflect.Uint64:
			newVal = uint64(rand.Int63())
		case reflect.Bool:
			var b bool
			if rand.Intn(2) == 1 {
				b = true
			}
			newVal = b
		case reflect.String:
			newVal = randStr(1)
		case reflect.Slice:
			sliceVal := typ.Elem()
			if sliceVal.Kind() != reflect.Uint8 {
				return fmt.Errorf("unsupported slice type: %T", typ.String())
			}
			newVal = randByteSlice(5 + rand.Intn(20))
		default:
			return fmt.Errorf("unsupported type: %T", typ.String())
		}
	}

	field.Set(reflect.ValueOf(newVal))

	return nil
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr(ln int) string {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(alphabet[rand.Intn(len(alphabet))])
	}

	return string(str)
}

func randByteSlice(ln int) []byte {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(rand.Intn(256))
	}

	return str
}
