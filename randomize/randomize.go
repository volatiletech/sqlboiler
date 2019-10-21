// Package randomize has helpers for randomization of structs and fields
package randomize

import (
	"math"
	"reflect"
	"regexp"
	"sort"
	"sync/atomic"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Randomizer allows a field to be randomized
type Randomizer interface {
	// Randomize should panic if there's no ability to randomize with the current parameters.
	//
	// nextInt can be called to create "random" sequential integers. This is done to avoid collisions in unique columns
	// for the tests.
	//
	// fieldType is used in the cases where the actual type (string, null string etc.) can actually be multiple
	// types of things that have specific randomization requirements, like a uuid for example is a normal null.String
	// but when randomizing that null string it must create a valid uuid or the database will reject it.
	//
	// shouldBeNull is a suggestion that a field should be null in this instance. The randomize implementation
	// can ignore this if the field cannot be null either because the type doesn't support it or there
	// is no ability for a field of this type to be null.
	Randomize(nextInt func() int64, fieldType string, shouldBeNull bool)
}

var (
	typeTime     = reflect.TypeOf(time.Time{})
	rgxValidTime = regexp.MustCompile(`[2-9]+`)
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

// NextInt retrives an integer in order
func (s *Seed) NextInt() int64 {
	return atomic.AddInt64((*int64)(s), 1)
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
			if strmangle.TitleCase(v) == fieldTyp.Name || v == fieldTyp.Tag.Get("boil") {
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

// randomizeField changes the value at field to a "randomized" value.
//
// If canBeNull is false:
//  The value will always be a non-null and non-zero value.
//
// If canBeNull is true:
//  The value has the possibility of being null or a non-zero value at random.
func randomizeField(s *Seed, field reflect.Value, fieldType string, canBeNull bool) error {
	kind := field.Kind()
	typ := field.Type()

	var shouldBeNull bool
	// Check the regular columns, these can be set or not set
	// depending on the canBeNull flag.
	// if canBeNull is false, then never return null values.
	if canBeNull {
		// 1 in 3 chance of being null or zero value
		shouldBeNull = s.NextInt()%3 == 0
	}

	// The struct and it's fields should always be addressable
	ptrToField := field.Addr()
	if r, ok := ptrToField.Interface().(Randomizer); ok {
		r.Randomize(s.NextInt, fieldType, shouldBeNull)
		return nil
	}

	var value interface{}

	if kind == reflect.Struct {
		if shouldBeNull {
			value = getStructNullValue(s, fieldType, typ)
		} else {
			value = getStructRandValue(s, fieldType, typ)
		}
	} else {
		// only get zero values for non byte slices
		// to stop mysql from being a jerk
		if shouldBeNull && kind != reflect.Slice {
			value = getVariableZeroValue(s, fieldType, kind, typ)
		} else {
			value = getVariableRandValue(s, fieldType, kind, typ)
		}
	}

	if value == nil {
		return errors.Errorf("unsupported type: %s", typ.String())
	}

	newValue := reflect.ValueOf(value)
	if reflect.TypeOf(value) != typ {
		newValue = newValue.Convert(typ)
	}

	field.Set(newValue)

	return nil
}

// getStructNullValue for the matching type.
func getStructNullValue(s *Seed, fieldType string, typ reflect.Type) interface{} {
	if typ == typeTime {
		// MySQL does not support 0 value time.Time, so use rand
		return Date(s.NextInt)
	}

	return nil
}

// getStructRandValue returns a "random" value for the matching type.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getStructRandValue(s *Seed, fieldType string, typ reflect.Type) interface{} {
	if typ == typeTime {
		return Date(s.NextInt)
	}

	return nil
}

// getVariableZeroValue for the matching type.
func getVariableZeroValue(s *Seed, fieldType string, kind reflect.Kind, typ reflect.Type) interface{} {
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
		// Some of these formatted strings cannot tolerate 0 values, so
		// we ignore the request for a null value.
		str, ok := FormattedString(s.NextInt, fieldType)
		if ok {
			return str
		}
		return ""
	case reflect.Slice:
		return []byte{}
	}

	return nil
}

// getVariableRandValue returns a "random" value for the matching kind.
// The randomness is really an incrementation of the global seed,
// this is done to avoid duplicate key violations.
func getVariableRandValue(s *Seed, fieldType string, kind reflect.Kind, typ reflect.Type) interface{} {
	switch kind {
	case reflect.Float32:
		return float32(float32(s.NextInt()%10)/10.0 + float32(s.NextInt()%10))
	case reflect.Float64:
		return float64(float64(s.NextInt()%10)/10.0 + float64(s.NextInt()%10))
	case reflect.Int:
		return int(s.NextInt())
	case reflect.Int8:
		return int8(s.NextInt() % math.MaxInt8)
	case reflect.Int16:
		return int16(s.NextInt() % math.MaxInt16)
	case reflect.Int32:
		val, ok := MediumInt(s.NextInt, fieldType)
		if ok {
			return val
		}
		return int32(s.NextInt() % math.MaxInt32)
	case reflect.Int64:
		return int64(s.NextInt())
	case reflect.Uint:
		return uint(s.NextInt())
	case reflect.Uint8:
		return uint8(s.NextInt() % math.MaxUint8)
	case reflect.Uint16:
		return uint16(s.NextInt() % math.MaxUint16)
	case reflect.Uint32:
		val, ok := MediumUint(s.NextInt, fieldType)
		if ok {
			return val
		}
		return uint32(s.NextInt() % math.MaxUint32)
	case reflect.Uint64:
		return uint64(s.NextInt())
	case reflect.Bool:
		return true
	case reflect.String:
		str, ok := FormattedString(s.NextInt, fieldType)
		if ok {
			return str
		}
		return Str(s.NextInt, 1)
	case reflect.Slice:
		sliceVal := typ.Elem()
		if sliceVal.Kind() != reflect.Uint8 {
			return errors.Errorf("unsupported slice type: %T, was expecting byte slice.", typ.String())
		}
		return ByteSlice(s.NextInt, 1)
	}

	return nil
}
