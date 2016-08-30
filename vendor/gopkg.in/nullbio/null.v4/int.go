package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullInt is a replica of sql.NullInt64 for int types.
type NullInt struct {
	Int   int
	Valid bool
}

// Int is an nullable int.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	NullInt
}

// NewInt creates a new Int
func NewInt(i int, valid bool) Int {
	return Int{
		NullInt: NullInt{
			Int:   i,
			Valid: valid,
		},
	}
}

// IntFrom creates a new Int that will always be valid.
func IntFrom(i int) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr(i *int) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Int.
// It also supports unmarshalling a sql.NullInt.
func (i *Int) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Int", reflect.TypeOf(v).Name())
	}
	i.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 0)
	i.Valid = err == nil
	if i.Valid {
		i.Int = int(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Int is null.
func (i Int) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(int64(i.Int), 10)), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int) {
	i.Int = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int {
	if !i.Valid {
		return nil
	}
	return &i.Int
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}

// Scan implements the Scanner interface.
func (n *NullInt) Scan(value interface{}) error {
	if value == nil {
		n.Int, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int, value)
}

// Value implements the driver Valuer interface.
func (n NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return int64(n.Int), nil
}
