package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullUint32 is a replica of sql.NullInt64 for uint32 types.
type NullUint32 struct {
	Uint32 uint32
	Valid  bool
}

// Uint32 is an nullable uint32.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint32 struct {
	NullUint32
}

// NewUint32 creates a new Uint32
func NewUint32(i uint32, valid bool) Uint32 {
	return Uint32{
		NullUint32: NullUint32{
			Uint32: i,
			Valid:  valid,
		},
	}
}

// Uint32From creates a new Uint32 that will always be valid.
func Uint32From(i uint32) Uint32 {
	return NewUint32(i, true)
}

// Uint32FromPtr creates a new Uint32 that be null if i is nil.
func Uint32FromPtr(i *uint32) Uint32 {
	if i == nil {
		return NewUint32(0, false)
	}
	return NewUint32(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint32.
// It also supports unmarshalling a sql.NullUint32.
func (i *Uint32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint32, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint32)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint32)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint32", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 32)
	i.Valid = err == nil
	if i.Valid {
		i.Uint32 = uint32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint32 is null.
func (i Uint32) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint32 is null.
func (i Uint32) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint32), 10)), nil
}

// SetValid changes this Uint32's value and also sets it to be non-null.
func (i *Uint32) SetValid(n uint32) {
	i.Uint32 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint32's value, or a nil pointer if this Uint32 is null.
func (i Uint32) Ptr() *uint32 {
	if !i.Valid {
		return nil
	}
	return &i.Uint32
}

// IsZero returns true for invalid Uint32's, for future omitempty support (Go 1.4?)
// A non-null Uint32 with a 0 value will not be considered zero.
func (i Uint32) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullUint32) Scan(value interface{}) error {
	if value == nil {
		n.Uint32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint32, value)
}

// Value implements the driver Valuer interface.
func (n NullUint32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint32), nil
}
