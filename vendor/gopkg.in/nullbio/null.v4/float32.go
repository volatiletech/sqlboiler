package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"gopkg.in/nullbio/null.v4/convert"
)

// NullFloat32 is a replica of sql.NullFloat64 for float32 types.
type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Float32 is a nullable float32.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float32 struct {
	NullFloat32
}

// NewFloat32 creates a new Float32
func NewFloat32(f float32, valid bool) Float32 {
	return Float32{
		NullFloat32: NullFloat32{
			Float32: f,
			Valid:   valid,
		},
	}
}

// Float32From creates a new Float32 that will always be valid.
func Float32From(f float32) Float32 {
	return NewFloat32(f, true)
}

// Float32FromPtr creates a new Float32 that be null if f is nil.
func Float32FromPtr(f *float32) Float32 {
	if f == nil {
		return NewFloat32(0, false)
	}
	return NewFloat32(*f, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Float32.
// It also supports unmarshalling a sql.NullFloat32.
func (f *Float32) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case float64:
		f.Float32 = float32(x)
	case map[string]interface{}:
		err = json.Unmarshal(data, &f.NullFloat32)
	case nil:
		f.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Float32", reflect.TypeOf(v).Name())
	}
	f.Valid = err == nil
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float32 if the input is a blank or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float32) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseFloat(string(text), 32)
	f.Valid = err == nil
	if f.Valid {
		f.Float32 = float32(res)
	}
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this Float32 is null.
func (f Float32) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float32 is null.
func (f Float32) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatFloat(float64(f.Float32), 'f', -1, 32)), nil
}

// SetValid changes this Float32's value and also sets it to be non-null.
func (f *Float32) SetValid(n float32) {
	f.Float32 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float32's value, or a nil pointer if this Float32 is null.
func (f Float32) Ptr() *float32 {
	if !f.Valid {
		return nil
	}
	return &f.Float32
}

// IsZero returns true for invalid Float32s, for future omitempty support (Go 1.4?)
// A non-null Float32 with a 0 value will not be considered zero.
func (f Float32) IsZero() bool {
	return !f.Valid
}

// Scan implements the Scanner interface.
func (n *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		n.Float32, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Float32, value)
}

// Value implements the driver Valuer interface.
func (n NullFloat32) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return float64(n.Float32), nil
}
