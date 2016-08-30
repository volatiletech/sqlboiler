package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullUint is a replica of sql.NullInt64 for uint types.
type NullUint struct {
	Uint  uint
	Valid bool
}

// Uint is an nullable uint.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Uint struct {
	NullUint
}

// NewUint creates a new Uint
func NewUint(i uint, valid bool) Uint {
	return Uint{
		NullUint: NullUint{
			Uint:  i,
			Valid: valid,
		},
	}
}

// UintFrom creates a new Uint that will always be valid.
func UintFrom(i uint) Uint {
	return NewUint(i, true)
}

// UintFromPtr creates a new Uint that be null if i is nil.
func UintFromPtr(i *uint) Uint {
	if i == nil {
		return NewUint(0, false)
	}
	return NewUint(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Uint.
// It also supports unmarshalling a sql.NullUint.
func (i *Uint) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to uint, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Uint)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullUint)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Uint", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Uint if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Uint) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseUint(string(text), 10, 0)
	i.Valid = err == nil
	if i.Valid {
		i.Uint = uint(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Uint is null.
func (i Uint) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Uint is null.
func (i Uint) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatUint(uint64(i.Uint), 10)), nil
}

// SetValid changes this Uint's value and also sets it to be non-null.
func (i *Uint) SetValid(n uint) {
	i.Uint = n
	i.Valid = true
}

// Ptr returns a pointer to this Uint's value, or a nil pointer if this Uint is null.
func (i Uint) Ptr() *uint {
	if !i.Valid {
		return nil
	}
	return &i.Uint
}

// IsZero returns true for invalid Uints, for future omitempty support (Go 1.4?)
// A non-null Uint with a 0 value will not be considered zero.
func (i Uint) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullUint) Scan(value interface{}) error {
	if value == nil {
		n.Uint, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Uint, value)
}

// Value implements the driver Valuer interface.
func (n NullUint) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Uint), nil
}
