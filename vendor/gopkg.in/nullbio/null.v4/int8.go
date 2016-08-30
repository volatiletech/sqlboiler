package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullInt8 is a replica of sql.NullInt64 for int8 types.
type NullInt8 struct {
	Int8  int8
	Valid bool
}

// Int8 is an nullable int8.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int8 struct {
	NullInt8
}

// NewInt8 creates a new Int8
func NewInt8(i int8, valid bool) Int8 {
	return Int8{
		NullInt8: NullInt8{
			Int8:  i,
			Valid: valid,
		},
	}
}

// Int8From creates a new Int8 that will always be valid.
func Int8From(i int8) Int8 {
	return NewInt8(i, true)
}

// Int8FromPtr creates a new Int8 that be null if i is nil.
func Int8FromPtr(i *int8) Int8 {
	if i == nil {
		return NewInt8(0, false)
	}
	return NewInt8(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int8.
// It also supports unmarshalling a sql.NullInt8.
func (i *Int8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int8, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int8)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int8", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int8 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 8)
	i.Valid = err == nil
	if i.Valid {
		i.Int8 = int8(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int8 is null.
func (i Int8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int8 is null.
func (i Int8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int8), 10)), nil
}

// SetValid changes this Int8's value and also sets it to be non-null.
func (i *Int8) SetValid(n int8) {
	i.Int8 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int8's value, or a nil pointer if this Int8 is null.
func (i Int8) Ptr() *int8 {
	if !i.Valid {
		return nil
	}
	return &i.Int8
}

// IsZero returns true for invalid Int8's, for future omitempty support (Go 1.4?)
// A non-null Int8 with a 0 value will not be considered zero.
func (i Int8) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullInt8) Scan(value interface{}) error {
	if value == nil {
		n.Int8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int8, value)
}

// Value implements the driver Valuer interface.
func (n NullInt8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Int8), nil
}
