package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullUint8 is a replica of sql.NullInt64 for uint8 types.
type NullUint8 struct {
	Uint8 uint8
	Valid bool
}

// Uint8 is an nullable uint8.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint8 struct {
	NullUint8
}

// NewUint8 creates a new Uint8
func NewUint8(i uint8, valid bool) Uint8 {
	return Uint8{
		NullUint8: NullUint8{
			Uint8: i,
			Valid: valid,
		},
	}
}

// Uint8From creates a new Uint8 that will always be valid.
func Uint8From(i uint8) Uint8 {
	return NewUint8(i, true)
}

// Uint8FromPtr creates a new Uint8 that be null if i is nil.
func Uint8FromPtr(i *uint8) Uint8 {
	if i == nil {
		return NewUint8(0, false)
	}
	return NewUint8(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint8.
// It also supports unmarshalling a sql.NullUint8.
func (i *Uint8) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint8, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint8)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint8)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint8", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint8 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint8) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 8)
	i.Valid = err == nil
	if i.Valid {
		i.Uint8 = uint8(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint8 is null.
func (i Uint8) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint8), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint8 is null.
func (i Uint8) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint8), 10)), nil
}

// SetValid changes this Uint8's value and also sets it to be non-null.
func (i *Uint8) SetValid(n uint8) {
	i.Uint8 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint8's value, or a nil pointer if this Uint8 is null.
func (i Uint8) Ptr() *uint8 {
	if !i.Valid {
		return nil
	}
	return &i.Uint8
}

// IsZero returns true for invalid Uint8's, for future omitempty support (Go 1.4?)
// A non-null Uint8 with a 0 value will not be considered zero.
func (i Uint8) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullUint8) Scan(value interface{}) error {
	if value == nil {
		n.Uint8, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint8, value)
}

// Value implements the driver Valuer interface.
func (n NullUint8) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint8), nil
}
