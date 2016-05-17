package boil

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"time"

	"github.com/guregu/null"
	"github.com/pobri19/sqlboiler/strmangle"
)

var (
	typeNullInt    = reflect.TypeOf(null.Int{})
	typeNullFloat  = reflect.TypeOf(null.Float{})
	typeNullString = reflect.TypeOf(null.String{})
	typeNullBool   = reflect.TypeOf(null.Bool{})
	typeNullTime   = reflect.TypeOf(null.Time{})
	typeTime       = reflect.TypeOf(time.Time{})
)

// Bind executes the query and inserts the
// result into the passed in object pointer
func (q *Query) Bind(obj interface{}) error {
	return nil
}

// BindOne inserts the returned row columns into the
// passed in object pointer
func BindOne(row *sql.Row, obj interface{}) error {
	return nil
}

// BindAll inserts the returned rows columns into the
// passed in slice of object pointers
func BindAll(rows *sql.Rows, obj interface{}) error {
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
		return fmt.Errorf("can only randomize pointers to structs, given: %T", str)
	}

	// Check if it's a struct
	value = value.Elem()
	kind = value.Kind()
	if kind != reflect.Struct {
		return fmt.Errorf("can only randomize pointers to structs, given: %T", str)
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

func randomizeField(field reflect.Value) error {
	kind := field.Kind()
	typ := field.Type()

	var newVal interface{}

	if kind == reflect.Struct {
		switch typ {
		case typeNullInt:
			newVal = null.NewInt(rand.Int63(), rand.Intn(2) == 1)
		case typeNullFloat:
			newVal = null.NewFloat(rand.Float64(), rand.Intn(2) == 1)
		case typeNullBool:
			newVal = null.NewBool(rand.Intn(2) == 1, rand.Intn(2) == 1)
		case typeNullString:
			newVal = null.NewString(randStr(5+rand.Intn(25)), rand.Intn(2) == 1)
		case typeNullTime:
			newVal = null.NewTime(time.Unix(rand.Int63(), 0), rand.Intn(2) == 1)
		case typeTime:
			newVal = time.Unix(rand.Int63(), 0)
		}
	} else {
		switch kind {
		case reflect.Int:
			newVal = rand.Int()
		case reflect.Int64:
			newVal = rand.Int63()
		case reflect.Float64:
			newVal = rand.Float64()
		case reflect.Bool:
			var b bool
			if rand.Intn(2) == 1 {
				b = true
			}
			newVal = b
		case reflect.String:
			newVal = randStr(5 + rand.Intn(20))
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
