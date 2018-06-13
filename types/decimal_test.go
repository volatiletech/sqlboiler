package types

import (
	"testing"

	"github.com/ericlagergren/decimal"
)

func TestDecimal_Value(t *testing.T) {
	t.Parallel()

	tests := []string{
		"3.14",
		"0",
		"43.4292",
	}

	for i, test := range tests {
		d := Decimal{Big: new(decimal.Big)}
		d.Big, _ = d.SetString(test)

		val, err := d.Value()
		if err != nil {
			t.Errorf("%d) %+v", i, err)
		}

		s, ok := val.(string)
		if !ok {
			t.Errorf("%d) wrong type returned", i)
		}

		if s != test {
			t.Errorf("%d) want: %s, got: %s", i, test, s)
		}
	}

	infinity := Decimal{Big: new(decimal.Big).SetInf(true)}
	if _, err := infinity.Value(); err == nil {
		t.Error("infinity should not be passed into the database")
	}
	nan := Decimal{Big: new(decimal.Big).SetNaN(true)}
	if _, err := nan.Value(); err == nil {
		t.Error("nan should not be passed into the database")
	}
}

func TestDecimal_Scan(t *testing.T) {
	t.Parallel()

	tests := []string{
		"3.14",
		"0",
		"43.4292",
	}

	for i, test := range tests {
		var d Decimal
		if err := d.Scan(test); err != nil {
			t.Error(err)
		}

		if got := d.String(); got != test {
			t.Errorf("%d) want: %s, got: %s", i, test, got)
		}
	}

	var d Decimal
	if err := d.Scan(nil); err == nil {
		t.Error("it should disallow scanning from a null value")
	}
}

func TestNullDecimal_Value(t *testing.T) {
	t.Parallel()

	tests := []string{
		"3.14",
		"0",
		"43.4292",
	}

	for i, test := range tests {
		d := NullDecimal{Big: new(decimal.Big)}
		d.Big, _ = d.SetString(test)

		val, err := d.Value()
		if err != nil {
			t.Errorf("%d) %+v", i, err)
		}

		s, ok := val.(string)
		if !ok {
			t.Errorf("%d) wrong type returned", i)
		}

		if s != test {
			t.Errorf("%d) want: %s, got: %s", i, test, s)
		}
	}

	infinity := NullDecimal{Big: new(decimal.Big).SetInf(true)}
	if _, err := infinity.Value(); err == nil {
		t.Error("infinity should not be passed into the database")
	}
	nan := NullDecimal{Big: new(decimal.Big).SetNaN(true)}
	if _, err := nan.Value(); err == nil {
		t.Error("nan should not be passed into the database")
	}
}

func TestNullDecimal_Scan(t *testing.T) {
	t.Parallel()

	tests := []string{
		"3.14",
		"0",
		"43.4292",
	}

	for i, test := range tests {
		var d NullDecimal
		if err := d.Scan(test); err != nil {
			t.Error(err)
		}

		if got := d.String(); got != test {
			t.Errorf("%d) want: %s, got: %s", i, test, got)
		}
	}

	var d NullDecimal
	if err := d.Scan(nil); err != nil {
		t.Error(err)
	}
	if d.Big != nil {
		t.Error("it should have been nil")
	}
}
