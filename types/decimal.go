package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/ericlagergren/decimal"
)

var (
	// DecimalContext is a global context that will be used when creating
	// decimals. It should be set once before any sqlboiler and then
	// assumed to be read-only after sqlboiler's first use.
	DecimalContext decimal.Context

	nullBytes = []byte("null")
)

var (
	_ driver.Valuer = Decimal{}
	_ driver.Valuer = NullDecimal{}
	_ sql.Scanner   = &Decimal{}
	_ sql.Scanner   = &NullDecimal{}
)

// Decimal is a DECIMAL in sql. Its zero value is valid for use with both
// Value and Scan.
//
// Although decimal can represent NaN and Infinity it will return an error
// if an attempt to store these values in the database is made.
//
// Because it cannot be nil, when Big is nil Value() will return "0"
// It will error if an attempt to Scan() a "null" value into it.
type Decimal struct {
	*decimal.Big
}

// NullDecimal is the same as Decimal, but allows the Big pointer to be nil.
// See documentation for Decimal for more details.
//
// When going into a database, if Big is nil it's value will be "null".
type NullDecimal struct {
	*decimal.Big
}

// NewDecimal creates a new decimal from a decimal
func NewDecimal(d *decimal.Big) Decimal {
	return Decimal{Big: d}
}

// NewNullDecimal creates a new null decimal from a decimal
func NewNullDecimal(d *decimal.Big) NullDecimal {
	return NullDecimal{Big: d}
}

// Value implements driver.Valuer.
func (d Decimal) Value() (driver.Value, error) {
	return decimalValue(d.Big, false)
}

// Scan implements sql.Scanner.
func (d *Decimal) Scan(val interface{}) error {
	newD, err := decimalScan(d.Big, val, false)
	if err != nil {
		return err
	}

	d.Big = newD
	return nil
}

// UnmarshalJSON allows marshalling JSON into a null pointer
func (d *Decimal) UnmarshalJSON(data []byte) error {
	if d.Big == nil {
		d.Big = new(decimal.Big)
	}

	return d.Big.UnmarshalJSON(data)
}

// MarshalText marshals a decimal value
func (d Decimal) MarshalText() ([]byte, error) {
	if d.Big == nil {
		return nullBytes, nil
	}

	return d.Big.MarshalText()
}

// UnmarshalText allows marshalling text into a null pointer
func (d *Decimal) UnmarshalText(data []byte) error {
	if d.Big == nil {
		d.Big = new(decimal.Big)
	}

	return d.Big.UnmarshalText(data)
}

// Randomize implements sqlboiler's randomize interface
func (d *Decimal) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	d.Big = randomDecimal(nextInt, fieldType, false)
}

// Value implements driver.Valuer.
func (n NullDecimal) Value() (driver.Value, error) {
	return decimalValue(n.Big, true)
}

// Scan implements sql.Scanner.
func (n *NullDecimal) Scan(val interface{}) error {
	newD, err := decimalScan(n.Big, val, true)
	if err != nil {
		return err
	}

	n.Big = newD
	return nil
}

// UnmarshalJSON allows marshalling JSON into a null pointer
func (n *NullDecimal) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		if n != nil {
			n.Big = nil
		}
		return nil
	}

	if n.Big == nil {
		n.Big = decimal.WithContext(DecimalContext)
	}

	return n.Big.UnmarshalJSON(data)
}

// MarshalText marshals a decimal value
func (n NullDecimal) MarshalText() ([]byte, error) {
	if n.Big == nil {
		return nullBytes, nil
	}

	return n.Big.MarshalText()
}

// UnmarshalText allows marshalling text into a null pointer
func (n *NullDecimal) UnmarshalText(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		if n != nil {
			n.Big = nil
		}
		return nil
	}

	if n.Big == nil {
		n.Big = decimal.WithContext(DecimalContext)
	}

	return n.Big.UnmarshalText(data)
}

// String impl
func (n NullDecimal) String() string {
	if n.Big == nil {
		return "nil"
	}
	return n.Big.String()
}

func (n NullDecimal) Format(f fmt.State, verb rune) {
	if n.Big == nil {
		fmt.Fprint(f, "nil")
		return
	}
	n.Big.Format(f, verb)
}

// MarshalJSON marshals a decimal value
func (n NullDecimal) MarshalJSON() ([]byte, error) {
	if n.Big == nil {
		return nullBytes, nil
	}

	return n.Big.MarshalText()
}

// IsZero implements qmhelper.Nullable
func (n NullDecimal) IsZero() bool {
	return n.Big == nil
}

// Randomize implements sqlboiler's randomize interface
func (n *NullDecimal) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	n.Big = randomDecimal(nextInt, fieldType, shouldBeNull)
}

func randomDecimal(nextInt func() int64, fieldType string, shouldBeNull bool) *decimal.Big {
	if shouldBeNull {
		return nil
	}

	randVal := fmt.Sprintf("%d.%d", nextInt()%10, nextInt()%10)
	random, success := decimal.WithContext(DecimalContext).SetString(randVal)
	if !success {
		panic("randVal could not be turned into a decimal")
	}

	return random
}

func decimalValue(d *decimal.Big, canNull bool) (driver.Value, error) {
	if d == nil {
		if canNull {
			return nil, nil
		}

		return "0", nil
	}

	if d.IsNaN(0) {
		return nil, errors.New("refusing to allow NaN into database")
	}
	if d.IsInf(0) {
		return nil, errors.New("refusing to allow infinity into database")
	}

	return d.String(), nil
}

func decimalScan(d *decimal.Big, val interface{}, canNull bool) (*decimal.Big, error) {
	if val == nil {
		if !canNull {
			return nil, errors.New("null cannot be scanned into decimal")
		}

		return nil, nil
	}

	switch t := val.(type) {
	case float64:
		if d == nil {
			d = decimal.WithContext(DecimalContext)
		}
		d.SetFloat64(t)
		return d, nil
	case int64:
		return decimal.WithContext(DecimalContext).SetMantScale(t, 0), nil
	case string:
		if d == nil {
			d = decimal.WithContext(DecimalContext)
		}
		if _, ok := d.SetString(t); !ok {
			if err := d.Context.Err(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("invalid decimal syntax: %q", t)
		}
		return d, nil
	case []byte:
		if d == nil {
			d = decimal.WithContext(DecimalContext)
		}
		if err := d.UnmarshalText(t); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, fmt.Errorf("cannot scan decimal value: %#v", val)
	}
}
