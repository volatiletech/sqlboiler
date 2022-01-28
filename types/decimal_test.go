package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
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

	zero := Decimal{}
	if _, err := zero.Value(); err != nil {
		t.Error("zero value should not error")
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

	zero := NullDecimal{}
	if _, err := zero.Value(); err != nil {
		t.Error("zero value should not error")
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

func TestDecimal_JSON(t *testing.T) {
	t.Parallel()

	s := struct {
		D Decimal `json:"d"`
	}{}

	err := json.Unmarshal([]byte(`{"d":"54.45"}`), &s)
	if err != nil {
		t.Error(err)
	}

	want, _ := new(decimal.Big).SetString("54.45")
	if s.D.Cmp(want) != 0 {
		t.Error("D was wrong:", s.D)
	}
}

func TestNullDecimal_JSON(t *testing.T) {
	t.Parallel()

	s := struct {
		N NullDecimal `json:"n"`
	}{}

	err := json.Unmarshal([]byte(`{"n":"54.45"}`), &s)
	if err != nil {
		t.Error(err)
	}

	want, _ := new(decimal.Big).SetString("54.45")
	if s.N.Cmp(want) != 0 {
		fmt.Println(want, s.N)
		t.Error("N was wrong:", s.N)
	}
}

func TestNullDecimal_JSONNil(t *testing.T) {
	t.Parallel()

	var n NullDecimal
	b, _ := json.Marshal(n)
	if string(b) != `null` {
		t.Errorf("want: null, got: %s", b)
	}

	n2 := new(NullDecimal)
	b, _ = json.Marshal(n2)
	if string(b) != `null` {
		t.Errorf("want: null, got: %s", b)
	}
}

func TestNullDecimal_IsZero(t *testing.T) {
	t.Parallel()

	var nullable qmhelper.Nullable = NullDecimal{}

	if !nullable.IsZero() {
		t.Error("it should be zero")
	}

	nullable = NullDecimal{Big: new(decimal.Big)}
	if nullable.IsZero() {
		t.Error("it should not be zero")
	}
}
