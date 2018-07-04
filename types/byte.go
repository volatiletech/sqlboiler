package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Byte is an alias for byte.
// Byte implements Marshal and Unmarshal.
type Byte byte

// String output your byte.
func (b Byte) String() string {
	return string(b)
}

// UnmarshalJSON sets *b to a copy of data.
func (b *Byte) UnmarshalJSON(data []byte) error {
	if b == nil {
		return errors.New("json: unmarshal json on nil pointer to byte")
	}

	var x string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if len(x) > 1 {
		return errors.New("json: cannot convert to byte, text len is greater than one")
	}

	*b = Byte(x[0])
	return nil
}

// MarshalJSON returns the JSON encoding of b.
func (b Byte) MarshalJSON() ([]byte, error) {
	return []byte{'"', byte(b), '"'}, nil
}

// Value returns b as a driver.Value.
func (b Byte) Value() (driver.Value, error) {
	return []byte{byte(b)}, nil
}

// Scan stores the src in *b.
func (b *Byte) Scan(src interface{}) error {
	switch src.(type) {
	case uint8:
		*b = Byte(src.(uint8))
	case string:
		*b = Byte(src.(string)[0])
	case []byte:
		*b = Byte(src.([]byte)[0])
	default:
		return errors.New("incompatible type for byte")
	}

	return nil
}

// Randomize for sqlboiler
func (b *Byte) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		*b = Byte(65) // Can't deal with a true 0-value
	}

	*b = Byte(nextInt()%60 + 65) // Can't deal with non-ascii characters in some databases
}
