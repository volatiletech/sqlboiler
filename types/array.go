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
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/lib/pq/oid"
	"github.com/volatiletech/randomize"
)

type parameterStatus struct {
	// server version in the same format as server_version_num, or 0 if
	// unavailable
	serverVersion int

	// the current location based on the TimeZone value of the session, if
	// available
	currentLocation *time.Location
}

func errorf(s string, args ...interface{}) {
	panic(fmt.Errorf("pq: %s", fmt.Sprintf(s, args...)))
}

func encode(parameterStatus *parameterStatus, x interface{}, pgtypOid oid.Oid) []byte {
	switch v := x.(type) {
	case int64:
		return strconv.AppendInt(nil, v, 10)
	case float64:
		return strconv.AppendFloat(nil, v, 'f', -1, 64)
	case []byte:
		if pgtypOid == oid.T_bytea {
			return encodeBytea(parameterStatus.serverVersion, v)
		}

		return v
	case string:
		if pgtypOid == oid.T_bytea {
			return encodeBytea(parameterStatus.serverVersion, []byte(v))
		}

		return []byte(v)
	case bool:
		return strconv.AppendBool(nil, v)
	case time.Time:
		return formatTs(v)

	default:
		errorf("encode: unknown type for %T", v)
	}

	panic("not reached")
}

// Parse a bytea value received from the server.  Both "hex" and the legacy
// "escape" format are supported.
func parseBytea(s []byte) (result []byte, err error) {
	if len(s) >= 2 && bytes.Equal(s[:2], []byte("\\x")) {
		// bytea_output = hex
		s = s[2:] // trim off leading "\\x"
		result = make([]byte, hex.DecodedLen(len(s)))
		_, err := hex.Decode(result, s)
		if err != nil {
			return nil, err
		}
	} else {
		// bytea_output = escape
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

func encodeBytea(serverVersion int, v []byte) (result []byte) {
	if serverVersion >= 90000 {
		// Use the hex format if we know that the server supports it
		result = make([]byte, 2+hex.EncodedLen(len(v)))
		result[0] = '\\'
		result[1] = 'x'
		hex.Encode(result[2:], v)
	} else {
		// .. or resort to "escape"
		for _, b := range v {
			if b == '\\' {
				result = append(result, '\\', '\\')
			} else if b < 0x20 || b > 0x7e {
				result = append(result, []byte(fmt.Sprintf("\\%03o", b))...)
			} else {
				result = append(result, b)
			}
		}
	}

	return result
}

var errInvalidTimestamp = errors.New("invalid timestamp")

type timestampParser struct {
	err error
}

func (p *timestampParser) expect(str string, char byte, pos int) {
	if p.err != nil {
		return
	}
	if pos+1 > len(str) {
		p.err = errInvalidTimestamp
		return
	}
	if c := str[pos]; c != char && p.err == nil {
		p.err = fmt.Errorf("expected '%v' at position %v; got '%v'", char, pos, c)
	}
}

func (p *timestampParser) mustAtoi(str string, begin int, end int) int {
	if p.err != nil {
		return 0
	}
	if begin < 0 || end < 0 || begin > end || end > len(str) {
		p.err = errInvalidTimestamp
		return 0
	}
	result, err := strconv.Atoi(str[begin:end])
	if err != nil {
		if p.err == nil {
			p.err = fmt.Errorf("expected number; got '%v'", str)
		}
		return 0
	}
	return result
}

// The location cache caches the time zones typically used by the client.
type locationCache struct {
	cache map[int]*time.Location
	lock  sync.Mutex
}

// All connections share the same list of timezones. Benchmarking shows that
// about 5% speed could be gained by putting the cache in the connection and
// losing the mutex, at the cost of a small amount of memory and a somewhat
// significant increase in code complexity.
var globalLocationCache = newLocationCache()

func newLocationCache() *locationCache {
	return &locationCache{cache: make(map[int]*time.Location)}
}

// Returns the cached timezone for the specified offset, creating and caching
// it if necessary.
func (c *locationCache) getLocation(offset int) *time.Location {
	c.lock.Lock()
	defer c.lock.Unlock()

	location, ok := c.cache[offset]
	if !ok {
		location = time.FixedZone("", offset)
		c.cache[offset] = location
	}

	return location
}

var infinityTsEnabled = false
var infinityTsNegative time.Time
var infinityTsPositive time.Time

const (
	infinityTsEnabledAlready        = "pq: infinity timestamp enabled already"
	infinityTsNegativeMustBeSmaller = "pq: infinity timestamp: negative value must be smaller (before) than positive"
)

// EnableInfinityTs controls the handling of Postgres' "-infinity" and
// "infinity" "timestamp"s.
//
// If EnableInfinityTs is not called, "-infinity" and "infinity" will return
// []byte("-infinity") and []byte("infinity") respectively, and potentially
// cause error "sql: Scan error on column index 0: unsupported driver -> Scan
// pair: []uint8 -> *time.Time", when scanning into a time.Time value.
//
// Once EnableInfinityTs has been called, all connections created using this
// driver will decode Postgres' "-infinity" and "infinity" for "timestamp",
// "timestamp with time zone" and "date" types to the predefined minimum and
// maximum times, respectively.  When encoding time.Time values, any time which
// equals or precedes the predefined minimum time will be encoded to
// "-infinity".  Any values at or past the maximum time will similarly be
// encoded to "infinity".
//
// If EnableInfinityTs is called with negative >= positive, it will panic.
// Calling EnableInfinityTs after a connection has been established results in
// undefined behavior.  If EnableInfinityTs is called more than once, it will
// panic.
func EnableInfinityTs(negative time.Time, positive time.Time) {
	if infinityTsEnabled {
		panic(infinityTsEnabledAlready)
	}
	if !negative.Before(positive) {
		panic(infinityTsNegativeMustBeSmaller)
	}
	infinityTsEnabled = true
	infinityTsNegative = negative
	infinityTsPositive = positive
}

/*
 * Testing might want to toggle infinityTsEnabled
 */
func disableInfinityTs() {
	infinityTsEnabled = false
}

// This is a time function specific to the Postgres default DateStyle
// setting ("ISO, MDY"), the only one we currently support. This
// accounts for the discrepancies between the parsing available with
// time.Parse and the Postgres date formatting quirks.
func parseTs(currentLocation *time.Location, str string) interface{} {
	switch str {
	case "-infinity":
		if infinityTsEnabled {
			return infinityTsNegative
		}
		return []byte(str)
	case "infinity":
		if infinityTsEnabled {
			return infinityTsPositive
		}
		return []byte(str)
	}
	t, err := ParseTimestamp(currentLocation, str)
	if err != nil {
		panic(err)
	}
	return t
}

// ParseTimestamp parses Postgres' text format. It returns a time.Time in
// currentLocation iff that time's offset agrees with the offset sent from the
// Postgres server. Otherwise, ParseTimestamp returns a time.Time with the
// fixed offset offset provided by the Postgres server.
func ParseTimestamp(currentLocation *time.Location, str string) (time.Time, error) {
	p := timestampParser{}

	monSep := strings.IndexRune(str, '-')
	// this is Gregorian year, not ISO Year
	// In Gregorian system, the year 1 BC is followed by AD 1
	year := p.mustAtoi(str, 0, monSep)
	daySep := monSep + 3
	month := p.mustAtoi(str, monSep+1, daySep)
	p.expect(str, '-', daySep)
	timeSep := daySep + 3
	day := p.mustAtoi(str, daySep+1, timeSep)

	var hour, minute, second int
	if len(str) > monSep+len("01-01")+1 {
		p.expect(str, ' ', timeSep)
		minSep := timeSep + 3
		p.expect(str, ':', minSep)
		hour = p.mustAtoi(str, timeSep+1, minSep)
		secSep := minSep + 3
		p.expect(str, ':', secSep)
		minute = p.mustAtoi(str, minSep+1, secSep)
		secEnd := secSep + 3
		second = p.mustAtoi(str, secSep+1, secEnd)
	}
	remainderIdx := monSep + len("01-01 00:00:00") + 1
	// Three optional (but ordered) sections follow: the
	// fractional seconds, the time zone offset, and the BC
	// designation. We set them up here and adjust the other
	// offsets if the preceding sections exist.

	nanoSec := 0
	tzOff := 0

	if remainderIdx < len(str) && str[remainderIdx] == '.' {
		fracStart := remainderIdx + 1
		fracOff := strings.IndexAny(str[fracStart:], "-+ ")
		if fracOff < 0 {
			fracOff = len(str) - fracStart
		}
		fracSec := p.mustAtoi(str, fracStart, fracStart+fracOff)
		nanoSec = fracSec * (1000000000 / int(math.Pow(10, float64(fracOff))))

		remainderIdx += fracOff + 1
	}
	if tzStart := remainderIdx; tzStart < len(str) && (str[tzStart] == '-' || str[tzStart] == '+') {
		// time zone separator is always '-' or '+' (UTC is +00)
		var tzSign int
		switch c := str[tzStart]; c {
		case '-':
			tzSign = -1
		case '+':
			tzSign = +1
		default:
			return time.Time{}, fmt.Errorf("expected '-' or '+' at position %v; got %v", tzStart, c)
		}
		tzHours := p.mustAtoi(str, tzStart+1, tzStart+3)
		remainderIdx += 3
		var tzMin, tzSec int
		if remainderIdx < len(str) && str[remainderIdx] == ':' {
			tzMin = p.mustAtoi(str, remainderIdx+1, remainderIdx+3)
			remainderIdx += 3
		}
		if remainderIdx < len(str) && str[remainderIdx] == ':' {
			tzSec = p.mustAtoi(str, remainderIdx+1, remainderIdx+3)
			remainderIdx += 3
		}
		tzOff = tzSign * ((tzHours * 60 * 60) + (tzMin * 60) + tzSec)
	}
	var isoYear int
	if remainderIdx+3 <= len(str) && str[remainderIdx:remainderIdx+3] == " BC" {
		isoYear = 1 - year
		remainderIdx += 3
	} else {
		isoYear = year
	}
	if remainderIdx < len(str) {
		return time.Time{}, fmt.Errorf("expected end of input, got %v", str[remainderIdx:])
	}
	t := time.Date(isoYear, time.Month(month), day,
		hour, minute, second, nanoSec,
		globalLocationCache.getLocation(tzOff))

	if currentLocation != nil {
		// Set the location of the returned Time based on the session's
		// TimeZone value, but only if the local time zone database agrees with
		// the remote database on the offset.
		lt := t.In(currentLocation)
		_, newOff := lt.Zone()
		if newOff == tzOff {
			t = lt
		}
	}

	return t, p.err
}

// formatTs formats t into a format postgres understands.
func formatTs(t time.Time) []byte {
	if infinityTsEnabled {
		// t <= -infinity : ! (t > -infinity)
		if !t.After(infinityTsNegative) {
			return []byte("-infinity")
		}
		// t >= infinity : ! (!t < infinity)
		if !t.Before(infinityTsPositive) {
			return []byte("infinity")
		}
	}
	return FormatTimestamp(t)
}

// FormatTimestamp formats t into Postgres' text format for timestamps.
func FormatTimestamp(t time.Time) []byte {
	// Need to send dates before 0001 A.D. with " BC" suffix, instead of the
	// minus sign preferred by Go.
	// Beware, "0000" in ISO is "1 BC", "-0001" is "2 BC" and so on
	bc := false
	if t.Year() <= 0 {
		// flip year sign, and add 1, e.g: "0" will be "1", and "-10" will be "11"
		t = t.AddDate((-t.Year())*2+1, 0, 0)
		bc = true
	}
	b := []byte(t.Format("2006-01-02 15:04:05.999999999Z07:00"))

	_, offset := t.Zone()
	offset %= 60
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

var typeByteSlice = reflect.TypeOf([]byte{})
var typeDriverValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var typeSQLScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

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
	}

	return GenericArray{a}
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
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to BoolArray", src)
}

func (a *BoolArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BoolArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BoolArray, len(elems))
		for i, v := range elems {
			if len(v) < 1 {
				return fmt.Errorf("boil: could not parse boolean array index %d: invalid boolean %q", i, v)
			}
			switch v[:1][0] {
			case 't', 'T':
				b[i] = true
			case 'f', 'F':
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

// Randomize for sqlboiler
func (a *BoolArray) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*a = BoolArray{nextInt()%2 == 0, nextInt()%2 == 0, nextInt()%2 == 0}
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
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to BytesArray", src)
}

func (a *BytesArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "BytesArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(BytesArray, len(elems))
		for i, v := range elems {
			b[i], err = parseBytea(v)
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

// Randomize for sqlboiler
func (a *BytesArray) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*a = BytesArray{randomize.ByteSlice(nextInt, 4), randomize.ByteSlice(nextInt, 4), randomize.ByteSlice(nextInt, 4)}
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
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to Float64Array", src)
}

func (a *Float64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Float64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
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

// Randomize for sqlboiler
func (a *Float64Array) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*a = Float64Array{float64(nextInt()), float64(nextInt())}
}

// GenericArray implements the driver.Valuer and sql.Scanner interfaces for
// an array or slice of any dimension.
type GenericArray struct{ A interface{} }

func (GenericArray) evaluateDestination(rt reflect.Type) (reflect.Type, func([]byte, reflect.Value) error, string) {
	var assign func([]byte, reflect.Value) error
	var del = ","

	// TODO calculate the assign function for other types
	// TODO repeat this section on the element type of arrays or slices (multidimensional)
	{
		if reflect.PtrTo(rt).Implements(typeSQLScanner) {
			// dest is always addressable because it is an element of a slice.
			assign = func(src []byte, dest reflect.Value) (err error) {
				ss := dest.Addr().Interface().(sql.Scanner)
				if src == nil {
					err = ss.Scan(nil)
				} else {
					err = ss.Scan(src)
				}
				return
			}
			goto FoundType
		}

		assign = func([]byte, reflect.Value) error {
			return fmt.Errorf("boil: scanning to %s is not implemented; only sql.Scanner", rt)
		}
	}

FoundType:

	if ad, ok := reflect.Zero(rt).Interface().(ArrayDelimiter); ok {
		del = ad.ArrayDelimiter()
	}

	return rt, assign, del
}

// Scan implements the sql.Scanner interface.
func (a GenericArray) Scan(src interface{}) error {
	dpv := reflect.ValueOf(a.A)
	switch {
	case dpv.Kind() != reflect.Ptr:
		return fmt.Errorf("boil: destination %T is not a pointer to array or slice", a.A)
	case dpv.IsNil():
		return fmt.Errorf("boil: destination %T is nil", a.A)
	}

	dv := dpv.Elem()
	switch dv.Kind() {
	case reflect.Slice:
	case reflect.Array:
	default:
		return fmt.Errorf("boil: destination %T is not a pointer to array or slice", a.A)
	}

	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src, dv)
	case string:
		return a.scanBytes([]byte(src), dv)
	case nil:
		if dv.Kind() == reflect.Slice {
			dv.Set(reflect.Zero(dv.Type()))
			return nil
		}
	}

	return fmt.Errorf("boil: cannot convert %T to %s", src, dv.Type())
}

func (a GenericArray) scanBytes(src []byte, dv reflect.Value) error {
	dtype, assign, del := a.evaluateDestination(dv.Type().Elem())
	dims, elems, err := parseArray(src, []byte(del))
	if err != nil {
		return err
	}

	// TODO allow multidimensional

	if len(dims) > 1 {
		return fmt.Errorf("boil: scanning from multidimensional ARRAY%s is not implemented",
			strings.ReplaceAll(fmt.Sprint(dims), " ", "]["))
	}

	// Treat a zero-dimensional array like an array with a single dimension of zero.
	if len(dims) == 0 {
		dims = append(dims, 0)
	}

	for i, rt := 0, dv.Type(); i < len(dims); i, rt = i+1, rt.Elem() {
		switch rt.Kind() {
		case reflect.Slice:
		case reflect.Array:
			if rt.Len() != dims[i] {
				return fmt.Errorf("boil: cannot convert ARRAY%s to %s",
					strings.ReplaceAll(fmt.Sprint(dims), " ", "]["), dv.Type())
			}
		default:
			// TODO handle multidimensional
		}
	}

	values := reflect.MakeSlice(reflect.SliceOf(dtype), len(elems), len(elems))
	for i, e := range elems {
		if err := assign(e, values.Index(i)); err != nil {
			return fmt.Errorf("boil: parsing array element index %d: %v", i, err)
		}
	}

	// TODO handle multidimensional

	switch dv.Kind() {
	case reflect.Slice:
		dv.Set(values.Slice(0, dims[0]))
	case reflect.Array:
		for i := 0; i < dims[0]; i++ {
			dv.Index(i).Set(values.Index(i))
		}
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (a GenericArray) Value() (driver.Value, error) {
	if a.A == nil {
		return nil, nil
	}

	rv := reflect.ValueOf(a.A)

	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			return nil, nil
		}
	case reflect.Array:
	default:
		return nil, fmt.Errorf("boil: Unable to convert %T to array", a.A)
	}

	if n := rv.Len(); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 0, 1+2*n)

		b, _, err := appendArray(b, rv, n)
		return string(b), err
	}

	return "{}", nil
}

// Int64Array represents a one-dimensional array of the PostgreSQL integer types.
type Int64Array []int64

// Scan implements the sql.Scanner interface.
func (a *Int64Array) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to Int64Array", src)
}

func (a *Int64Array) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "Int64Array")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
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

// Randomize for sqlboiler
func (a *Int64Array) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*a = Int64Array{int64(nextInt()), int64(nextInt())}
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
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to StringArray", src)
}

func (a *StringArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "StringArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
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

// Randomize for sqlboiler
func (a *StringArray) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	strs := make([]string, 2)
	fieldType = strings.TrimPrefix(fieldType, "ARRAY")

	for i := range strs {
		val, ok := randomize.FormattedString(nextInt, fieldType)
		if ok {
			strs[i] = val
			continue
		}

		strs[i] = randomize.Str(nextInt, 1)
	}

	*a = strs
}

// DecimalArray represents a one-dimensional array of the decimal type.
type DecimalArray []Decimal

// Scan implements the sql.Scanner interface.
func (a *DecimalArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("boil: cannot convert %T to DecimalArray", src)
}

func (a *DecimalArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "DecimalArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(DecimalArray, len(elems))
		for i, v := range elems {
			var success bool
			b[i].Big, success = new(decimal.Big).SetString(string(v))
			if !success {
				return fmt.Errorf("boil: parsing decimal element index as decimal %d: %s", i, v)
			}
		}
		*a = b
	}
	return nil
}

// Value implements the driver.Valuer interface.
func (a DecimalArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	} else if len(a) == 0 {
		return "{}", nil
	}

	strs := make([]string, len(a))
	for i, d := range a {
		strs[i] = d.String()
	}

	return "{" + strings.Join(strs, ",") + "}", nil
}

// Randomize for sqlboiler
func (a *DecimalArray) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	d1, d2 := NewDecimal(new(decimal.Big)), NewDecimal(new(decimal.Big))
	d1.SetString(fmt.Sprintf("%d.%d", nextInt()%10, nextInt()%10))
	d2.SetString(fmt.Sprintf("%d.%d", nextInt()%10, nextInt()%10))
	*a = DecimalArray{d1, d2}
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
	return append(b, encode(nil, v, 0)...), nil
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
			if depth == len(dims) {
				break Element
			}
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
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
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
		return nil, fmt.Errorf("boil: cannot convert ARRAY%s to %s", strings.ReplaceAll(fmt.Sprint(dims), " ", "]["), typ)
	}
	return elems, err
}
