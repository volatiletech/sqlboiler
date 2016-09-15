// Copyright (c) 2011-2013, 'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany. MIT license.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software
// is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var typeByteSlice = reflect.TypeOf([]byte{})
var typeDriverValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var typeSQLScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func encode(x interface{}) []byte {
	switch v := x.(type) {
	case int64:
		return strconv.AppendInt(nil, v, 10)
	case float64:
		return strconv.AppendFloat(nil, v, 'f', -1, 64)
	case []byte:
		return encodeBytes(v)
	case string:
		return []byte(v)
	case bool:
		return strconv.AppendBool(nil, v)
	case time.Time:
		return formatTimestamp(v)

	default:
		panic(fmt.Errorf("encode: unknown type for %T", v))
	}
}

// FormatTimestamp formats t into Postgres' text format for timestamps.
func formatTimestamp(t time.Time) []byte {
	// Need to send dates before 0001 A.D. with " BC" suffix, instead of the
	// minus sign preferred by Go.
	// Beware, "0000" in ISO is "1 BC", "-0001" is "2 BC" and so on
	bc := false
	if t.Year() <= 0 {
		// flip year sign, and add 1, e.g: "0" will be "1", and "-10" will be "11"
		t = t.AddDate((-t.Year())*2+1, 0, 0)
		bc = true
	}
	b := []byte(t.Format(time.RFC3339Nano))

	_, offset := t.Zone()
	offset = offset % 60
	if offset != 0 {
		// RFC3339Nano already printed the minus sign
		if offset < 0 {
			offset = -offset
		}

		b = append(b, ':')
		if offset < 10 {
			b = append(b, '0')
		}
		b = strconv.AppendInt(b, int64(offset), 10)
	}

	if bc {
		b = append(b, " BC"...)
	}
	return b
}

func encodeBytes(v []byte) (result []byte) {
	for _, b := range v {
		if b == '\\' {
			result = append(result, '\\', '\\')
		} else if b < 0x20 || b > 0x7e {
			result = append(result, []byte(fmt.Sprintf("\\%03o", b))...)
		} else {
			result = append(result, b)
		}
	}

	return result
}

// Parse a bytea value received from the server.  Both "hex" and the legacy
// "escape" format are supported.
func parseBytes(s []byte) (result []byte, err error) {
	if len(s) >= 2 && bytes.Equal(s[:2], []byte("\\x")) {
		// bytea_output = hex
		s = s[2:] // trim off leading "\\x"
		result = make([]byte, hex.DecodedLen(len(s)))
		_, err := hex.Decode(result, s)
		if err != nil {
			return nil, err
		}
	} else {
		for len(s) > 0 {
			if s[0] == '\\' {
				// escaped '\\'
				if len(s) >= 2 && s[1] == '\\' {
					result = append(result, '\\')
					s = s[2:]
					continue
				}

				// '\\' followed by an octal number
				if len(s) < 4 {
					return nil, fmt.Errorf("invalid bytea sequence %v", s)
				}
				r, err := strconv.ParseInt(string(s[1:4]), 8, 9)
				if err != nil {
					return nil, fmt.Errorf("could not parse bytea value: %s", err.Error())
				}
				result = append(result, byte(r))
				s = s[4:]
			} else {
				// We hit an unescaped, raw byte.  Try to read in as many as
				// possible in one go.
				i := bytes.IndexByte(s, '\\')
				if i == -1 {
					result = append(result, s...)
					break
				}
				result = append(result, s[:i]...)
				s = s[i:]
			}
		}
	}

	return result, nil
}

// Array returns the optimal driver.Valuer and sql.Scanner for an array or
// slice of any dimension.
//
// For example:
//  db.Query(`SELECT * FROM t WHERE id = ANY($1)`, pq.Array([]int{235, 401}))
//
//  var x []sql.NullInt64
//  db.QueryRow('SELECT ARRAY[235, 401]').Scan(pq.Array(&x))
//
// Scanning multi-dimensional arrays is not supported.  Arrays where the lower
// bound is not one (such as `[0:0]={1}') are not supported.
func Array(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	switch a := a.(type) {
	case []bool:
		return (*BoolArray)(&a)
	case []float64:
		return (*Float64Array)(&a)
	case []int64:
		return (*Int64Array)(&a)
	case []string:
		return (*StringArray)(&a)

	case *[]bool:
		return (*BoolArray)(a)
	case *[]float64:
		return (*Float64Array)(a)
	case *[]int64:
		return (*Int64Array)(a)
	case *[]string:
		return (*StringArray)(a)

	default:
		panic(fmt.Sprintf("boil: invalid type received %T", a))
	}
}

// ArrayDelimiter may be optionally implemented by driver.Valuer or sql.Scanner
// to override the array delimiter used by GenericArray.
type ArrayDelimiter interface {
	// ArrayDelimiter returns the delimiter character(s) for this element's type.
	ArrayDelimiter() string
}

// BoolArray represents a one-dimensional array of the PostgreSQL boolean type.
type BoolArray []bool

// Scan implements the sql.Scanner interface.
func (a *BoolArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	}

	return fmt.Errorf("boil: cannot convert %T to BoolArray", src)
}

func (a *BoolArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BoolArray")
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolArray, len(elems))
		for i, v := range elems {
			if len(v) != 1 {
				return fmt.Errorf("boil: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
			switch v[0] {
			case 't':
				b[i] = true
			case 'f':
				b[i] = false
			default:
				return fmt.Errorf("boil: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a BoolArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be exactly two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1+2*n)

		for i := 0; i < n; i++ {
			b[2*i] = ','
			if a[i] {
				b[1+2*i] = 't'
			} else {
				b[1+2*i] = 'f'
			}
		}

		b[0] = '{'
		b[2*n] = '}'

		return string(b), nil
	}

	return "{}", nil
}

// BytesArray represents a one-dimensional array of the PostgreSQL bytea type.
type BytesArray [][]byte

// Scan implements the sql.Scanner interface.
func (a *BytesArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	}

	return fmt.Errorf("boil: cannot convert %T to BytesArray", src)
}

func (a *BytesArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BytesArray")
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BytesArray, len(elems))
		for i, v := range elems {
			b[i], err = parseBytes(v)
			if err != nil {
				return fmt.Errorf("could not parse bytea array index %d: %s", i, err.Error())
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface. It uses the "hex" format which
// is only supported on PostgreSQL 9.0 or newer.
func (a BytesArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// 3*N bytes of hex formatting, and N-1 bytes of delimiters.
		size := 1 + 6*n
		for _, x := range a {
			size += hex.EncodedLen(len(x))
		}

		b := make([]byte, size)

		for i, s := 0, b; i < n; i++ {
			o := copy(s, `,"\\x`)
			o += hex.Encode(s[o:], a[i])
			s[o] = '"'
			s = s[o+1:]
		}

		b[0] = '{'
		b[size-1] = '}'

		return string(b), nil
	}

	return "{}", nil
}

// Float64Array represents a one-dimensional array of the PostgreSQL double
// precision type.
type Float64Array []float64

// Scan implements the sql.Scanner interface.
func (a *Float64Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	}

	return fmt.Errorf("boil: cannot convert %T to Float64Array", src)
}

func (a *Float64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float64Array")
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Float64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseFloat(string(v), 64); err != nil {
				return fmt.Errorf("boil: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Float64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendFloat(b, a[0], 'f', -1, 64)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, a[i], 'f', -1, 64)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

type Int64Array []int64

// Scan implements the sql.Scanner interface.
func (a *Int64Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	}

	return fmt.Errorf("boil: cannot convert %T to Int64Array", src)
}

func (a *Int64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Int64Array")
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(Int64Array, len(elems))
		for i, v := range elems {
			if b[i], err = strconv.ParseInt(string(v), 10, 64); err != nil {
				return fmt.Errorf("boil: parsing array element index %d: %v", i, err)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a Int64Array) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '{'

		b = strconv.AppendInt(b, a[0], 10)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendInt(b, a[i], 10)
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// StringArray represents a one-dimensional array of the PostgreSQL character types.
type StringArray []string

// Scan implements the sql.Scanner interface.
func (a *StringArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	}

	return fmt.Errorf("boil: cannot convert %T to StringArray", src)
}

func (a *StringArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "StringArray")
	if err != nil {
		return err
	}
	if len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(StringArray, len(elems))
		for i, v := range elems {
			if b[i] = string(v); v == nil {
				return fmt.Errorf("boil: parsing array element index %d: cannot convert nil to string", i)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}

	if n := len(a); n > 0 {
		// There will be at least two curly brackets, 2*N bytes of quotes,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+3*n)
		b[0] = '{'

		b = appendArrayQuotedBytes(b, []byte(a[0]))
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = appendArrayQuotedBytes(b, []byte(a[i]))
		}

		return string(append(b, '}')), nil
	}

	return "{}", nil
}

// appendArray appends rv to the buffer, returning the extended buffer and
// the delimiter used between elements.
//
// It panics when n <= 0 or rv's Kind is not reflect.Array nor reflect.Slice.
func appendArray(b []byte, rv reflect.Value, n int) ([]byte, string, error) {
	var del string
	var err error

	b = append(b, '{')

	if b, del, err = appendArrayElement(b, rv.Index(0)); err != nil {
		return b, del, err
	}

	for i := 1; i < n; i++ {
		b = append(b, del...)
		if b, del, err = appendArrayElement(b, rv.Index(i)); err != nil {
			return b, del, err
		}
	}

	return append(b, '}'), del, nil
}

// appendArrayElement appends rv to the buffer, returning the extended buffer
// and the delimiter to use before the next element.
//
// When rv's Kind is neither reflect.Array nor reflect.Slice, it is converted
// using driver.DefaultParameterConverter and the resulting []byte or string
// is double-quoted.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func appendArrayElement(b []byte, rv reflect.Value) ([]byte, string, error) {
	if k := rv.Kind(); k == reflect.Array || k == reflect.Slice {
		if t := rv.Type(); t != typeByteSlice && !t.Implements(typeDriverValuer) {
			if n := rv.Len(); n > 0 {
				return appendArray(b, rv, n)
			}

			return b, "", nil
		}
	}

	var del = ","
	var err error
	var iv interface{} = rv.Interface()

	if ad, ok := iv.(ArrayDelimiter); ok {
		del = ad.ArrayDelimiter()
	}

	if iv, err = driver.DefaultParameterConverter.ConvertValue(iv); err != nil {
		return b, del, err
	}

	switch v := iv.(type) {
	case nil:
		return append(b, "NULL"...), del, nil
	case []byte:
		return appendArrayQuotedBytes(b, v), del, nil
	case string:
		return appendArrayQuotedBytes(b, []byte(v)), del, nil
	}

	b, err = appendValue(b, iv)
	return b, del, err
}

func appendArrayQuotedBytes(b, v []byte) []byte {
	b = append(b, '"')
	for {
		i := bytes.IndexAny(v, `"\`)
		if i < 0 {
			b = append(b, v...)
			break
		}
		if i > 0 {
			b = append(b, v[:i]...)
		}
		b = append(b, '\\', v[i])
		v = v[i+1:]
	}
	return append(b, '"')
}

func appendValue(b []byte, v driver.Value) ([]byte, error) {
	return append(b, encode(v)...), nil
}

// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("boil: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("boil: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("boil: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("boil: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("boil: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("boil: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}

func scanLinearArray(src, del []byte, typ string) (elems [][]byte, err error) {
	dims, elems, err := parseArray(src, del)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 {
		return nil, fmt.Errorf("boil: cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}
	return elems, err
}
