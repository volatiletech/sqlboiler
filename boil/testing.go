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
func RandomizeStruct(str interface{}, colTypes map[string]string, includeInvalid bool, blacklist ...string) error {
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
		if err := randomizeField(fieldVal, fieldDBType, includeInvalid); err != nil {
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

func randomizeField(field reflect.Value, fieldType string, canBeNull bool) error {
	kind := field.Kind()
	typ := field.Type()

	var newVal interface{}

	if kind == reflect.Struct {
		var notNull bool
		if canBeNull {
			notNull = rand.Intn(2) == 1
		} else {
			notNull = false
		}

		switch typ {
		case typeNullBool:
			if notNull {
				newVal = null.NewBool(sd.nextInt()%2 == 0, true)
			} else {
				newVal = null.NewBool(false, false)
			}
		case typeNullString:
			if fieldType == "uuid" {
				newVal = null.NewString(uuid.NewV4().String(), true)
			} else if notNull {
				switch fieldType {
				case "interval":
					newVal = null.NewString(strconv.Itoa((sd.nextInt()%26)+2)+" days", true)
				default:
					newVal = null.NewString(randStr(1, sd.nextInt()), true)
				}
			} else {
				newVal = null.NewString("", false)
			}
		case typeNullTime:
			if notNull {
				newVal = null.NewTime(randDate(sd.nextInt()), true)
			} else {
				newVal = null.NewTime(time.Time{}, false)
			}
		case typeTime:
			newVal = randDate(sd.nextInt())
		case typeNullFloat32:
			if notNull {
				newVal = null.NewFloat32(float32(sd.nextInt()%10)/10.0+float32(sd.nextInt()%10), true)
			} else {
				newVal = null.NewFloat32(0.0, false)
			}
		case typeNullFloat64:
			if notNull {
				newVal = null.NewFloat64(float64(sd.nextInt()%10)/10.0+float64(sd.nextInt()%10), true)
			} else {
				newVal = null.NewFloat64(0.0, false)
			}
		case typeNullInt:
			if notNull {
				newVal = null.NewInt(sd.nextInt(), true)
			} else {
				newVal = null.NewInt(0, false)
			}
		case typeNullInt8:
			if notNull {
				newVal = null.NewInt8(int8(sd.nextInt()), true)
			} else {
				newVal = null.NewInt8(0, false)
			}
		case typeNullInt16:
			if notNull {
				newVal = null.NewInt16(int16(sd.nextInt()), true)
			} else {
				newVal = null.NewInt16(0, false)
			}
		case typeNullInt32:
			if notNull {
				newVal = null.NewInt32(int32(sd.nextInt()), true)
			} else {
				newVal = null.NewInt32(0, false)
			}
		case typeNullInt64:
			if notNull {
				newVal = null.NewInt64(int64(sd.nextInt()), true)
			} else {
				newVal = null.NewInt64(0, false)
			}
		case typeNullUint:
			if notNull {
				newVal = null.NewUint(uint(sd.nextInt()), true)
			} else {
				newVal = null.NewUint(0, false)
			}
		case typeNullUint8:
			if notNull {
				newVal = null.NewUint8(uint8(sd.nextInt()), true)
			} else {
				newVal = null.NewUint8(0, false)
			}
		case typeNullUint16:
			if notNull {
				newVal = null.NewUint16(uint16(sd.nextInt()), true)
			} else {
				newVal = null.NewUint16(0, false)
			}
		case typeNullUint32:
			if notNull {
				newVal = null.NewUint32(uint32(sd.nextInt()), true)
			} else {
				newVal = null.NewUint32(0, false)
			}
		case typeNullUint64:
			if notNull {
				newVal = null.NewUint64(uint64(sd.nextInt()), true)
			} else {
				newVal = null.NewUint64(0, false)
			}
		}
	} else {
		switch kind {
		case reflect.Float32:
			newVal = float32(float32(sd.nextInt()%10)/10.0 + float32(sd.nextInt()%10))
		case reflect.Float64:
			newVal = float64(float64(sd.nextInt()%10)/10.0 + float64(sd.nextInt()%10))
		case reflect.Int:
			newVal = sd.nextInt()
		case reflect.Int8:
			newVal = int8(sd.nextInt())
		case reflect.Int16:
			newVal = int16(sd.nextInt())
		case reflect.Int32:
			newVal = int32(sd.nextInt())
		case reflect.Int64:
			newVal = int64(sd.nextInt())
		case reflect.Uint:
			newVal = uint(sd.nextInt())
		case reflect.Uint8:
			newVal = uint8(sd.nextInt())
		case reflect.Uint16:
			newVal = uint16(sd.nextInt())
		case reflect.Uint32:
			newVal = uint32(sd.nextInt())
		case reflect.Uint64:
			newVal = uint64(sd.nextInt())
		case reflect.Bool:
			newVal = sd.nextInt()%2 == 0
		case reflect.String:
			switch fieldType {
			case "interval":
				newVal = strconv.Itoa((sd.nextInt()%26)+2) + " days"
			case "uuid":
				newVal = uuid.NewV4().String()
			default:
				newVal = randStr(1, sd.nextInt())
			}
		case reflect.Slice:
			sliceVal := typ.Elem()
			if sliceVal.Kind() != reflect.Uint8 {
				return errors.Errorf("unsupported slice type: %T", typ.String())
			}
			newVal = randByteSlice(5+rand.Intn(20), sd.nextInt())
		default:
			return errors.Errorf("unsupported type: %T", typ.String())
		}
	}

	field.Set(reflect.ValueOf(newVal))

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
