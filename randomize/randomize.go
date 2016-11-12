// Package randomize has helpers for randomization of structs and fields
package randomize

import (
	"database/sql"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	null "gopkg.in/nullbio/null.v6"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/vattle/sqlboiler/strmangle"
	"github.com/vattle/sqlboiler/types"
)

var (
	typeNullFloat32  = reflect.TypeOf(null.Float32{})
	typeNullFloat64  = reflect.TypeOf(null.Float64{})
	typeNullInt      = reflect.TypeOf(null.Int{})
	typeNullInt8     = reflect.TypeOf(null.Int8{})
	typeNullInt16    = reflect.TypeOf(null.Int16{})
	typeNullInt32    = reflect.TypeOf(null.Int32{})
	typeNullInt64    = reflect.TypeOf(null.Int64{})
	typeNullUint     = reflect.TypeOf(null.Uint{})
	typeNullUint8    = reflect.TypeOf(null.Uint8{})
	typeNullUint16   = reflect.TypeOf(null.Uint16{})
	typeNullUint32   = reflect.TypeOf(null.Uint32{})
	typeNullUint64   = reflect.TypeOf(null.Uint64{})
	typeNullString   = reflect.TypeOf(null.String{})
	typeNullByte     = reflect.TypeOf(null.Byte{})
	typeNullBool     = reflect.TypeOf(null.Bool{})
	typeNullTime     = reflect.TypeOf(null.Time{})
	typeNullBytes    = reflect.TypeOf(null.Bytes{})
	typeNullJSON     = reflect.TypeOf(null.JSON{})
	typeTime         = reflect.TypeOf(time.Time{})
	typeJSON         = reflect.TypeOf(types.JSON{})
	typeInt64Array   = reflect.TypeOf(types.Int64Array{})
	typeBytesArray   = reflect.TypeOf(types.BytesArray{})
	typeBoolArray    = reflect.TypeOf(types.BoolArray{})
	typeFloat64Array = reflect.TypeOf(types.Float64Array{})
	typeStringArray  = reflect.TypeOf(types.StringArray{})
	typeHStore       = reflect.TypeOf(types.HStore{})
	rgxValidTime     = regexp.MustCompile(`[2-9]+`)

	validatedTypes = []string{
		"inet", "line", "uuid", "interval", "mediumint",
		"json", "jsonb", "box", "cidr", "circle",
		"lseg", "macaddr", "path", "pg_lsn", "point",
		"polygon", "txid_snapshot", "money", "hstore",
	}
)

// Seed is an atomic counter for pseudo-randomization structs. Using full
// randomization leads to collisions in a domain where uniqueness is an
// important factor.
type Seed int64

// NewSeed creates a new seed for pseudo-randomization.
func NewSeed() *Seed {
	s := new(int64)
	*s = time.Now().Unix()
	return (*Seed)(s)
}

func (s *Seed) nextInt() int {
	return int(atomic.AddInt64((*int64)(s), 1))
}

// Struct gets its fields filled with random data based on the seed.
// It will ignore the fields in the blacklist.
// It will ignore fields that have the struct tag boil:"-"
func Struct(s *Seed, str interface{}, colTypes map[string]string, canBeNull bool, blacklist ...string) error {
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

		if fieldTyp.Tag.Get("boil") == "-" {
			continue
		}

		fieldDBType := colTypes[fieldTyp.Name]
		if err := randomizeField(s, fieldVal, fieldDBType, canBeNull); err != nil {
			return err
		}
	}

	return nil
}

// randDate generates a random time.Time between 1850 and 2050.
// Only the Day/Month/Year columns are set so that Dates and DateTimes do
// not cause mismatches in the test data comparisons.
func randDate(s *Seed) time.Time {
	t := time.Date(
		1972+s.nextInt()%60,
		time.Month(1+(s.nextInt()%12)),
		1+(s.nextInt()%25),
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
func randomizeField(s *Seed, field reflect.Value, fieldType string, canBeNull bool) error {

	kind := field.Kind()
	typ := field.Type()

	if strings.HasPrefix(fieldType, "enum") {
		enum, err := randEnumValue(fieldType)
		if err != nil {
			return err
		}

		if kind == reflect.Struct {
			val := null.NewString(enum, rand.Intn(1) == 0)
			field.Set(reflect.ValueOf(val))
		} else {
			field.Set(reflect.ValueOf(enum))
		}

		return nil
	}

	var value interface{}
	var isNull bool

	// Validated columns always need to be set regardless of canBeNull,
	// and they have to adhere to a strict value format.
	foundValidated := strmangle.SetInclude(fieldType, validatedTypes)

	if foundValidated {
		if kind == reflect.Struct {
			switch typ {
			case typeNullString:
				if fieldType == "interval" {
					value = null.NewString(strconv.Itoa((s.nextInt()%26)+2)+" days", true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "uuid" {
					value = null.NewString(uuid.NewV4().String(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "box" || fieldType == "line" || fieldType == "lseg" ||
					fieldType == "path" || fieldType == "polygon" {
					value = null.NewString(randBox(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "cidr" || fieldType == "inet" {
					value = null.NewString(randNetAddr(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "macaddr" {
					value = null.NewString(randMacAddr(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "circle" {
					value = null.NewString(randCircle(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "pg_lsn" {
					value = null.NewString(randLsn(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "point" {
					value = null.NewString(randPoint(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "txid_snapshot" {
					value = null.NewString(randTxID(), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "money" {
					value = null.NewString(randMoney(s), true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
			case typeNullInt32:
				if fieldType == "mediumint" {
					// 8388607 is the max for 3 byte int
					value = null.NewInt32(int32(s.nextInt())%8388607, true)
					field.Set(reflect.ValueOf(value))
					return nil
				}
			case typeNullJSON:
				value = null.NewJSON([]byte(fmt.Sprintf(`"%s"`, randStr(s, 1))), true)
				field.Set(reflect.ValueOf(value))
				return nil
			case typeHStore:
				value := types.HStore{}
				value[randStr(s, 3)] = sql.NullString{String: randStr(s, 3), Valid: s.nextInt()%3 == 0}
				value[randStr(s, 3)] = sql.NullString{String: randStr(s, 3), Valid: s.nextInt()%3 == 0}
				field.Set(reflect.ValueOf(value))
				return nil
			}

		} else {
			switch kind {
			case reflect.String:
				if fieldType == "interval" {
					value = strconv.Itoa((s.nextInt()%26)+2) + " days"
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "uuid" {
					value = uuid.NewV4().String()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "box" || fieldType == "line" || fieldType == "lseg" ||
					fieldType == "path" || fieldType == "polygon" {
					value = randBox()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "cidr" || fieldType == "inet" {
					value = randNetAddr()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "macaddr" {
					value = randMacAddr()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "circle" {
					value = randCircle()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "pg_lsn" {
					value = randLsn()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "point" {
					value = randPoint()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "txid_snapshot" {
					value = randTxID()
					field.Set(reflect.ValueOf(value))
					return nil
				}
				if fieldType == "money" {
					value = randMoney(s)
					field.Set(reflect.ValueOf(value))
					return nil
				}
			case reflect.Int32:
				if fieldType == "mediumint" {
					// 8388607 is the max for 3 byte int
					value = int32(s.nextInt()) % 8388607
					field.Set(reflect.ValueOf(value))
					return nil
				}
			}
			switch typ {
			case typeJSON:
				value = []byte(fmt.Sprintf(`"%s"`, randStr(s, 1)))
				field.Set(reflect.ValueOf(value))
				return nil
			case typeHStore:
				value := types.HStore{}
				value[randStr(s, 3)] = sql.NullString{String: randStr(s, 3), Valid: s.nextInt()%3 == 0}
				value[randStr(s, 3)] = sql.NullString{String: randStr(s, 3), Valid: s.nextInt()%3 == 0}
				field.Set(reflect.ValueOf(value))
				return nil
			}
		}
	}

	// Check the regular columns, these can be set or not set
	// depending on the canBeNull flag.
	if canBeNull {
		// 1 in 3 chance of being null or zero value
		isNull = s.nextInt()%3 == 0
	} else {
		// if canBeNull is false, then never return null values.
		isNull = false
	}

	// If it's a Postgres array, treat it like one
	if strings.HasPrefix(fieldType, "ARRAY") {
		value = getArrayRandValue(s, typ, fieldType)
		// Retrieve the value to be returned
	} else if kind == reflect.Struct {
		if isNull {
			value = getStructNullValue(s, typ)
		} else {
			value = getStructRandValue(s, typ)
		}
	} else {
		// only get zero values for non byte slices
		// to stop mysql from being a jerk
		if isNull && kind != reflect.Slice {
			value = getVariableZeroValue(s, kind, typ)
		} else {
			value = getVariableRandValue(s, kind, typ)
		}
	}

	if value == nil {
		return errors.Errorf("unsupported type: %s", typ.String())
	}

	field.Set(reflect.ValueOf(value))

	return nil
}

func getArrayRandValue(s *Seed, typ reflect.Type, fieldType string) interface{} {
	fieldType = strings.TrimLeft(fieldType, "ARRAY")
	switch typ {
	case typeInt64Array:
		return types.Int64Array{int64(s.nextInt()), int64(s.nextInt())}
	case typeFloat64Array:
		return types.Float64Array{float64(s.nextInt()), float64(s.nextInt())}
	case typeBoolArray:
		return types.BoolArray{s.nextInt()%2 == 0, s.nextInt()%2 == 0, s.nextInt()%2 == 0}
	case typeStringArray:
		if fieldType == "interval" {
			value := strconv.Itoa((s.nextInt()%26)+2) + " days"
			return types.StringArray{value, value}
		}
		if fieldType == "uuid" {
			value := uuid.NewV4().String()
			return types.StringArray{value, value}
		}
		if fieldType == "box" || fieldType == "line" || fieldType == "lseg" ||
			fieldType == "path" || fieldType == "polygon" {
			value := randBox()
			return types.StringArray{value, value}
		}
		if fieldType == "cidr" || fieldType == "inet" {
			value := randNetAddr()
			return types.StringArray{value, value}
		}
		if fieldType == "macaddr" {
			value := randMacAddr()
			return types.StringArray{value, value}
		}
		if fieldType == "circle" {
			value := randCircle()
			return types.StringArray{value, value}
		}
		if fieldType == "pg_lsn" {
			value := randLsn()
			return types.StringArray{value, value}
		}
		if fieldType == "point" {
			value := randPoint()
			return types.StringArray{value, value}
		}
		if fieldType == "txid_snapshot" {
			value := randTxID()
			return types.StringArray{value, value}
		}
		if fieldType == "money" {
			value := randMoney(s)
			return types.StringArray{value, value}
		}
		if fieldType == "json" || fieldType == "jsonb" {
			value := []byte(fmt.Sprintf(`"%s"`, randStr(s, 1)))
			return types.StringArray{string(value)}
		}
		return types.StringArray{randStr(s, 4), randStr(s, 4), randStr(s, 4)}
	case typeBytesArray:
		return types.BytesArray{randByteSlice(s, 4), randByteSlice(s, 4), randByteSlice(s, 4)}
	}

	return nil
}

// getStructNullValue for the matching type.
func getStructNullValue(s *Seed, typ reflect.Type) interface{} {
	switch typ {
	case typeTime:
		// MySQL does not support 0 value time.Time, so use rand
		return randDate(s)
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
	case typeNullBytes:
		return null.NewBytes(nil, false)
	case typeNullByte:
		return null.NewByte(byte(0), false)
	}

	return nil
}

// getStructRandValue returns a "random" value for the matching type.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getStructRandValue(s *Seed, typ reflect.Type) interface{} {
	switch typ {
	case typeTime:
		return randDate(s)
	case typeNullBool:
		return null.NewBool(s.nextInt()%2 == 0, true)
	case typeNullString:
		return null.NewString(randStr(s, 1), true)
	case typeNullTime:
		return null.NewTime(randDate(s), true)
	case typeNullFloat32:
		return null.NewFloat32(float32(s.nextInt()%10)/10.0+float32(s.nextInt()%10), true)
	case typeNullFloat64:
		return null.NewFloat64(float64(s.nextInt()%10)/10.0+float64(s.nextInt()%10), true)
	case typeNullInt:
		return null.NewInt(int(int32(s.nextInt())), true)
	case typeNullInt8:
		return null.NewInt8(int8(s.nextInt()), true)
	case typeNullInt16:
		return null.NewInt16(int16(s.nextInt()), true)
	case typeNullInt32:
		return null.NewInt32(int32(s.nextInt()), true)
	case typeNullInt64:
		return null.NewInt64(int64(s.nextInt()), true)
	case typeNullUint:
		return null.NewUint(uint(s.nextInt()), true)
	case typeNullUint8:
		return null.NewUint8(uint8(s.nextInt()), true)
	case typeNullUint16:
		return null.NewUint16(uint16(s.nextInt()), true)
	case typeNullUint32:
		return null.NewUint32(uint32(s.nextInt()), true)
	case typeNullUint64:
		return null.NewUint64(uint64(s.nextInt()), true)
	case typeNullBytes:
		return null.NewBytes(randByteSlice(s, 1), true)
	case typeNullByte:
		return null.NewByte(byte(rand.Intn(125-65)+65), true)
	}

	return nil
}

// getVariableZeroValue for the matching type.
func getVariableZeroValue(s *Seed, kind reflect.Kind, typ reflect.Type) interface{} {
	switch typ.String() {
	case "types.Byte":
		// Decimal 65 is 'A'. 0 is not a valid UTF8, so cannot use a zero value here.
		return types.Byte(65)
	}

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
		return []byte{}
	}

	return nil
}

// getVariableRandValue returns a "random" value for the matching kind.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getVariableRandValue(s *Seed, kind reflect.Kind, typ reflect.Type) interface{} {
	switch typ.String() {
	case "types.Byte":
		return types.Byte(rand.Intn(125-65) + 65)
	}

	switch kind {
	case reflect.Float32:
		return float32(float32(s.nextInt()%10)/10.0 + float32(s.nextInt()%10))
	case reflect.Float64:
		return float64(float64(s.nextInt()%10)/10.0 + float64(s.nextInt()%10))
	case reflect.Int:
		return s.nextInt()
	case reflect.Int8:
		return int8(s.nextInt())
	case reflect.Int16:
		return int16(s.nextInt())
	case reflect.Int32:
		return int32(s.nextInt())
	case reflect.Int64:
		return int64(s.nextInt())
	case reflect.Uint:
		return uint(s.nextInt())
	case reflect.Uint8:
		return uint8(s.nextInt())
	case reflect.Uint16:
		return uint16(s.nextInt())
	case reflect.Uint32:
		return uint32(s.nextInt())
	case reflect.Uint64:
		return uint64(s.nextInt())
	case reflect.Bool:
		return true
	case reflect.String:
		return randStr(s, 1)
	case reflect.Slice:
		sliceVal := typ.Elem()
		if sliceVal.Kind() != reflect.Uint8 {
			return errors.Errorf("unsupported slice type: %T, was expecting byte slice.", typ.String())
		}
		return randByteSlice(s, 1)
	}

	return nil
}

func randEnumValue(enum string) (string, error) {
	vals := strmangle.ParseEnumVals(enum)
	if vals == nil || len(vals) == 0 {
		return "", fmt.Errorf("unable to parse enum string: %s", enum)
	}

	return vals[rand.Intn(len(vals)-1)], nil
}
