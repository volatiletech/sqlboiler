package boil

import (
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"sync/atomic"
	"time"

	null "gopkg.in/nullbio/null.v4"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/vattle/sqlboiler/strmangle"
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

	rgxValidTime = regexp.MustCompile(`[2-9]+`)
)

type seed int64

var sd = new(seed)

func (s *seed) nextInt() int {
	return int(atomic.AddInt64((*int64)(s), 1))
}

// RandomizeStruct takes an object and fills it with random data.
// It will ignore the fields in the blacklist.
func RandomizeStruct(str interface{}, colTypes map[string]string, canBeNull bool, blacklist ...string) error {
	// Don't modify blacklist
	copyBlacklist := make([]string, len(blacklist))
	copy(copyBlacklist, blacklist)
	blacklist = copyBlacklist

	sort.Strings(blacklist)

	// Check if it's pointer
	value := reflect.ValueOf(str)
	kind := value.Kind()
	if kind != reflect.Ptr {
		return errors.Errorf("Outer element should be a pointer, given a non-pointer: %T", str)
	}

	// Check if it's a struct
	value = value.Elem()
	kind = value.Kind()
	if kind != reflect.Struct {
		return errors.Errorf("Inner element should be a struct, given a non-struct: %T", str)
	}

	typ := value.Type()
	nFields := value.NumField()

	// Iterate through fields, randomizing
	for i := 0; i < nFields; i++ {
		fieldVal := value.Field(i)
		fieldTyp := typ.Field(i)

		var found bool
		for _, v := range blacklist {
			if strmangle.TitleCase(v) == fieldTyp.Name {
				found = true
				break
			}
		}

		if found {
			continue
		}

		fieldDBType := colTypes[fieldTyp.Name]
		if err := randomizeField(fieldVal, fieldDBType, canBeNull); err != nil {
			return err
		}
	}

	return nil
}

// randDate generates a random time.Time between 1850 and 2050.
// Only the Day/Month/Year columns are set so that Dates and DateTimes do
// not cause mismatches in the test data comparisons.
func randDate(sd int) time.Time {
	t := time.Date(
		1850+rand.Intn(sd),
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

// randomizeField changes the value at field to a "randomized" value.
//
// If canBeNull is false:
//  The value will always be a non-null and non-zero value.

// If canBeNull is true:
//  The value has the possibility of being null or non-zero at random.
func randomizeField(field reflect.Value, fieldType string, canBeNull bool) error {
	kind := field.Kind()
	typ := field.Type()

	var value interface{}
	var isNull bool

	// Validated columns always need to be set regardless of canBeNull,
	// and they have to adhere to a strict value format.
	validatedTypes := []string{"uuid", "interval"}
	var foundValidated bool

	for _, validated := range validatedTypes {
		if fieldType == validated {
			foundValidated = true
			break
		}
	}

	if foundValidated {
		if kind == reflect.Struct {
			switch typ {
			case typeNullString:
				if fieldType == "interval" {
					value = null.NewString(strconv.Itoa((sd.nextInt()%26)+2)+" days", true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "uuid" {
					value = null.NewString(uuid.NewV4().String(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
			}
		} else {
			switch kind {
			case reflect.String:
				if fieldType == "interval" {
					value = strconv.Itoa((sd.nextInt()%26)+2) + " days"
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "uuid" {
					value = uuid.NewV4().String()
					field.Set(reflect.ValueOf(value))
					return nil
				}
			}
		}
	}

	// Check the regular columns, these can be set or not set
	// depending on the canBeNull flag.
	if canBeNull {
		// 1 in 3 chance of being null or zero value
		isNull = rand.Intn(3) == 1
	} else {
		// if canBeNull is false, then never return null values.
		isNull = false
	}

	// Retrieve the value to be returned
	if kind == reflect.Struct {
		if isNull {
			value = getStructNullValue(typ)
		} else {
			value = getStructRandValue(typ)
		}
	} else {
		if isNull {
			value = getVariableNullValue(kind)
		} else {
			value = getVariableRandValue(kind, typ)
		}
	}

	if value == nil {
		return errors.Errorf("unsupported type: %T", typ.String())
	}

	field.Set(reflect.ValueOf(value))
	return nil
}

// getStructNullValue for the matching type.
func getStructNullValue(typ reflect.Type) interface{} {
	switch typ {
	case typeTime:
		return time.Time{}
	case typeNullBool:
		return null.NewBool(false, false)
	case typeNullString:
		return null.NewString("", false)
	case typeNullTime:
		return null.NewTime(time.Time{}, false)
	case typeNullFloat32:
		return null.NewFloat32(0.0, false)
	case typeNullFloat64:
		return null.NewFloat64(0.0, false)
	case typeNullInt:
		return null.NewInt(0, false)
	case typeNullInt8:
		return null.NewInt8(0, false)
	case typeNullInt16:
		return null.NewInt16(0, false)
	case typeNullInt32:
		return null.NewInt32(0, false)
	case typeNullInt64:
		return null.NewInt64(0, false)
	case typeNullUint:
		return null.NewUint(0, false)
	case typeNullUint8:
		return null.NewUint8(0, false)
	case typeNullUint16:
		return null.NewUint16(0, false)
	case typeNullUint32:
		return null.NewUint32(0, false)
	case typeNullUint64:
		return null.NewUint64(0, false)
	}

	return nil
}

// getStructRandValue returns a "random" value for the matching type.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getStructRandValue(typ reflect.Type) interface{} {
	switch typ {
	case typeTime:
		return randDate(sd.nextInt())
	case typeNullBool:
		return null.NewBool(sd.nextInt()%2 == 0, true)
	case typeNullString:
		return null.NewString(randStr(1, sd.nextInt()), true)
	case typeNullTime:
		return null.NewTime(randDate(sd.nextInt()), true)
	case typeNullFloat32:
		return null.NewFloat32(float32(sd.nextInt()%10)/10.0+float32(sd.nextInt()%10), true)
	case typeNullFloat64:
		return null.NewFloat64(float64(sd.nextInt()%10)/10.0+float64(sd.nextInt()%10), true)
	case typeNullInt:
		return null.NewInt(sd.nextInt(), true)
	case typeNullInt8:
		return null.NewInt8(int8(sd.nextInt()), true)
	case typeNullInt16:
		return null.NewInt16(int16(sd.nextInt()), true)
	case typeNullInt32:
		return null.NewInt32(int32(sd.nextInt()), true)
	case typeNullInt64:
		return null.NewInt64(int64(sd.nextInt()), true)
	case typeNullUint:
		return null.NewUint(uint(sd.nextInt()), true)
	case typeNullUint8:
		return null.NewUint8(uint8(sd.nextInt()), true)
	case typeNullUint16:
		return null.NewUint16(uint16(sd.nextInt()), true)
	case typeNullUint32:
		return null.NewUint32(uint32(sd.nextInt()), true)
	case typeNullUint64:
		return null.NewUint64(uint64(sd.nextInt()), true)
	}

	return nil
}

// getVariableNullValue for the matching type.
func getVariableNullValue(kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Float32:
		return float32(0)
	case reflect.Float64:
		return float64(0)
	case reflect.Int:
		return int(0)
	case reflect.Int8:
		return int8(0)
	case reflect.Int16:
		return int16(0)
	case reflect.Int32:
		return int32(0)
	case reflect.Int64:
		return int64(0)
	case reflect.Uint:
		return uint(0)
	case reflect.Uint8:
		return uint8(0)
	case reflect.Uint16:
		return uint16(0)
	case reflect.Uint32:
		return uint32(0)
	case reflect.Uint64:
		return uint64(0)
	case reflect.Bool:
		return false
	case reflect.String:
		return ""
	case reflect.Slice:
		return []byte(nil)
	}

	return nil
}

// getVariableRandValue returns a "random" value for the matching kind.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getVariableRandValue(kind reflect.Kind, typ reflect.Type) interface{} {
	switch kind {
	case reflect.Float32:
		return float32(float32(sd.nextInt()%10)/10.0 + float32(sd.nextInt()%10))
	case reflect.Float64:
		return float64(float64(sd.nextInt()%10)/10.0 + float64(sd.nextInt()%10))
	case reflect.Int:
		return sd.nextInt()
	case reflect.Int8:
		return int8(sd.nextInt())
	case reflect.Int16:
		return int16(sd.nextInt())
	case reflect.Int32:
		return int32(sd.nextInt())
	case reflect.Int64:
		return int64(sd.nextInt())
	case reflect.Uint:
		return uint(sd.nextInt())
	case reflect.Uint8:
		return uint8(sd.nextInt())
	case reflect.Uint16:
		return uint16(sd.nextInt())
	case reflect.Uint32:
		return uint32(sd.nextInt())
	case reflect.Uint64:
		return uint64(sd.nextInt())
	case reflect.Bool:
		return true
	case reflect.String:
		return randStr(1, sd.nextInt())
	case reflect.Slice:
		sliceVal := typ.Elem()
		if sliceVal.Kind() != reflect.Uint8 {
			return errors.Errorf("unsupported slice type: %T, was expecting byte slice.", typ.String())
		}
		return randByteSlice(5+rand.Intn(20), sd.nextInt())
	}

	return nil
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStr(ln int, s int) string {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(alphabet[s%len(alphabet)])
	}

	return string(str)
}

func randByteSlice(ln int, s int) []byte {
	str := make([]byte, ln)
	for i := 0; i < ln; i++ {
		str[i] = byte(s % 256)
	}

	return str
}
