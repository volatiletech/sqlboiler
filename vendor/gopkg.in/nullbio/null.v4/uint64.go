package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullUint64 is a replica of sql.NullInt64 for uint64 types.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Uint64 is an nullable uint64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint64 struct {
	NullUint64
}

// NewUint64 creates a new Uint64
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		NullUint64: NullUint64{
			Uint64: i,
			Valid:  valid,
		},
	}
}

// Uint64From creates a new Uint64 that will always be valid.
func Uint64From(i uint64) Uint64 {
	return NewUint64(i, true)
}

// Uint64FromPtr creates a new Uint64 that be null if i is nil.
func Uint64FromPtr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}
	return NewUint64(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint64.
// It also supports unmarshalling a sql.NullUint64.
func (i *Uint64) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint64, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint64)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint64)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint64", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint64 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint64) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 64)
	i.Valid = err == nil
	if i.Valid {
		i.Uint64 = uint64(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint64 is null.
func (i Uint64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(i.Uint64, 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint64 is null.
func (i Uint64) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(i.Uint64, 10)), nil
}

// SetValid changes this Uint64's value and also sets it to be non-null.
func (i *Uint64) SetValid(n uint64) {
	i.Uint64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint64's value, or a nil pointer if this Uint64 is null.
func (i Uint64) Ptr() *uint64 {
	if !i.Valid {
		return nil
	}
	return &i.Uint64
}

// IsZero returns true for invalid Uint64's, for future omitempty support (Go 1.4?)
// A non-null Uint64 with a 0 value will not be considered zero.
func (i Uint64) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullUint64) Scan(value interface{}) error {
	if value == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint64, value)
}

// Value implements the driver Valuer interface.
func (n NullUint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint64), nil
}
