package boil

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
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

	rgxValidTime = regexp.MustCompile(`[2-9]+`)
)

type seed int

var sd = new(seed)

func (s *seed) nextInt() int {
	nextInt := int(*s)
	*s++
	return nextInt
}

// IsZeroValue checks if the variables with matching columns in obj
// are or are not zero values, depending on whether shouldZero is true or false
func IsZeroValue(obj interface{}, shouldZero bool, columns ...string) []error {
	val := reflect.Indirect(reflect.ValueOf(obj))

	var errs []error
	for _, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		if !field.IsValid() {
			panic(fmt.Sprintf("Unable to find variable with column name %s", c))
		}

		zv := reflect.Zero(field.Type())
		if shouldZero && !reflect.DeepEqual(field.Interface(), zv.Interface()) {
			errs = append(errs, fmt.Errorf("Column with name %s is not zero value: %#v, %#v", c, field.Interface(), zv.Interface()))
		} else if !shouldZero && reflect.DeepEqual(field.Interface(), zv.Interface()) {
			errs = append(errs, fmt.Errorf("Column with name %s is zero value: %#v, %#v", c, field.Interface(), zv.Interface()))
		}
	}

	return errs
}

// IsValueMatch checks whether the variables in obj with matching column names
// match the values in the values slice.
func IsValueMatch(obj interface{}, columns []string, values []interface{}) []error {
	val := reflect.Indirect(reflect.ValueOf(obj))

	var errs []error
	for i, c := range columns {
		field := val.FieldByName(strmangle.TitleCase(c))
		if !field.IsValid() {
			panic(fmt.Sprintf("Unable to find variable with column name %s", c))
		}

		typ := field.Type().String()
		if typ == "time.Time" || typ == "null.Time" {
			var timeField reflect.Value
			var valTimeStr string
			if typ == "time.Time" {
				valTimeStr = values[i].(time.Time).String()
				timeField = field
			} else {
				valTimeStr = values[i].(null.Time).Time.String()
				timeField = field.FieldByName("Time")
				validField := field.FieldByName("Valid")
				if validField.Interface() != values[i].(null.Time).Valid {
					errs = append(errs, fmt.Errorf("Null.Time column with name %s Valid field does not match: %v ≠ %v", c, values[i].(null.Time).Valid, validField.Interface()))
				}
			}

			if (rgxValidTime.MatchString(valTimeStr) && timeField.Interface() == reflect.Zero(timeField.Type()).Interface()) ||
				(!rgxValidTime.MatchString(valTimeStr) && timeField.Interface() != reflect.Zero(timeField.Type()).Interface()) {
				errs = append(errs, fmt.Errorf("Time column with name %s Time field does not match: %v ≠ %v", c, values[i], timeField.Interface()))
			}

			continue
		}

		if !reflect.DeepEqual(field.Interface(), values[i]) {
			errs = append(errs, fmt.Errorf("Column with name %s does not match value: %#v ≠ %#v", c, values[i], field.Interface()))
		}
	}

	return errs
}

// RandomizeSlice takes a pointer to a slice of pointers to objects
// and fills the pointed to objects with random data.
// It will ignore the fields in the blacklist.
func RandomizeSlice(obj interface{}, colTypes map[string]string, includeInvalid bool, blacklist ...string) error {
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
		if err := RandomizeStruct(o.Interface(), colTypes, includeInvalid, blacklist...); err != nil {
			return err
		}
	}

	return nil
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

		fieldDBType := colTypes[typ.Field(i).Name]
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

func randomizeField(field reflect.Value, fieldType string, includeInvalid bool) error {
	kind := field.Kind()
	typ := field.Type()

	var newVal interface{}

	if kind == reflect.Struct {
		var b bool
		if includeInvalid {
			b = rand.Intn(2) == 1
		} else {
			b = true
		}
		switch typ {
		case typeNullBool:
			if b {
				newVal = null.NewBool(sd.nextInt()%2 == 0, b)
			} else {
				newVal = null.NewBool(false, false)
			}
		case typeNullString:
			if b {
				if fieldType == "interval" {
					newVal = null.NewString(strconv.Itoa((sd.nextInt()%26)+2)+" days", b)
				} else {
					newVal = null.NewString(randStr(1, sd.nextInt()), b)
				}
			} else {
				newVal = null.NewString("", false)
			}
		case typeNullTime:
			if b {
				newVal = null.NewTime(randDate(sd.nextInt()), b)
			} else {
				newVal = null.NewTime(time.Time{}, false)
			}
		case typeTime:
			newVal = randDate(sd.nextInt())
		case typeNullFloat32:
			if b {
				newVal = null.NewFloat32(float32(sd.nextInt()%10)/10.0+float32(sd.nextInt()%10), b)
			} else {
				newVal = null.NewFloat32(0.0, false)
			}
		case typeNullFloat64:
			if b {
				newVal = null.NewFloat64(float64(sd.nextInt()%10)/10.0+float64(sd.nextInt()%10), b)
			} else {
				newVal = null.NewFloat64(0.0, false)
			}
		case typeNullInt:
			if b {
				newVal = null.NewInt(sd.nextInt(), b)
			} else {
				newVal = null.NewInt(0, false)
			}
		case typeNullInt8:
			if b {
				newVal = null.NewInt8(int8(sd.nextInt()), b)
			} else {
				newVal = null.NewInt8(0, false)
			}
		case typeNullInt16:
			if b {
				newVal = null.NewInt16(int16(sd.nextInt()), b)
			} else {
				newVal = null.NewInt16(0, false)
			}
		case typeNullInt32:
			if b {
				newVal = null.NewInt32(int32(sd.nextInt()), b)
			} else {
				newVal = null.NewInt32(0, false)
			}
		case typeNullInt64:
			if b {
				newVal = null.NewInt64(int64(sd.nextInt()), b)
			} else {
				newVal = null.NewInt64(0, false)
			}
		case typeNullUint:
			if b {
				newVal = null.NewUint(uint(sd.nextInt()), b)
			} else {
				newVal = null.NewUint(0, false)
			}
		case typeNullUint8:
			if b {
				newVal = null.NewUint8(uint8(sd.nextInt()), b)
			} else {
				newVal = null.NewUint8(0, false)
			}
		case typeNullUint16:
			if b {
				newVal = null.NewUint16(uint16(sd.nextInt()), b)
			} else {
				newVal = null.NewUint16(0, false)
			}
		case typeNullUint32:
			if b {
				newVal = null.NewUint32(uint32(sd.nextInt()), b)
			} else {
				newVal = null.NewUint32(0, false)
			}
		case typeNullUint64:
			if b {
				newVal = null.NewUint64(uint64(sd.nextInt()), b)
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
			var b bool
			if includeInvalid {
				b = sd.nextInt()%2 == 0
			} else {
				b = true
			}
			newVal = b
		case reflect.String:
			if fieldType == "interval" {
				newVal = strconv.Itoa((sd.nextInt()%26)+2) + " days"
			} else {
				newVal = randStr(1, sd.nextInt())
			}
		case reflect.Slice:
			sliceVal := typ.Elem()
			if sliceVal.Kind() != reflect.Uint8 {
				return fmt.Errorf("unsupported slice type: %T", typ.String())
			}
			newVal = randByteSlice(5+rand.Intn(20), sd.nextInt())
		default:
			return fmt.Errorf("unsupported type: %T", typ.String())
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
