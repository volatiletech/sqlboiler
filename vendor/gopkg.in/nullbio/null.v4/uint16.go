package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullUint16 is a replica of sql.NullInt64 for uint16 types.
type NullUint16 struct {
	Uint16 uint16
	Valid  bool
}

// Uint16 is an nullable uint16.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint16 struct {
	NullUint16
}

// NewUint16 creates a new Uint16
func NewUint16(i uint16, valid bool) Uint16 {
	return Uint16{
		NullUint16: NullUint16{
			Uint16: i,
			Valid:  valid,
		},
	}
}

// Uint16From creates a new Uint16 that will always be valid.
func Uint16From(i uint16) Uint16 {
	return NewUint16(i, true)
}

// Uint16FromPtr creates a new Uint16 that be null if i is nil.
func Uint16FromPtr(i *uint16) Uint16 {
	if i == nil {
		return NewUint16(0, false)
	}
	return NewUint16(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint16.
// It also supports unmarshalling a sql.NullUint16.
func (i *Uint16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint16, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint16)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint16)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint16", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint16 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 16)
	i.Valid = err == nil
	if i.Valid {
		i.Uint16 = uint16(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint16 is null.
func (i Uint16) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint16), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint16 is null.
func (i Uint16) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint16), 10)), nil
}

// SetValid changes this Uint16's value and also sets it to be non-null.
func (i *Uint16) SetValid(n uint16) {
	i.Uint16 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint16's value, or a nil pointer if this Uint16 is null.
func (i Uint16) Ptr() *uint16 {
	if !i.Valid {
		return nil
	}
	return &i.Uint16
}

// IsZero returns true for invalid Uint16's, for future omitempty support (Go 1.4?)
// A non-null Uint16 with a 0 value will not be considered zero.
func (i Uint16) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullUint16) Scan(value interface{}) error {
	if value == nil {
		n.Uint16, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint16, value)
}

// Value implements the driver Valuer interface.
func (n NullUint16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint16), nil
}
