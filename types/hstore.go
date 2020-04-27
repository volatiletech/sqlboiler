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
	"database/sql"
	"database/sql/driver"
	"strings"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/randomize"
)

// HStore is a wrapper for transferring HStore values back and forth easily.
type HStore map[string]null.String

// escapes and quotes hstore keys/values
// s should be a sql.NullString or string
func hQuote(s interface{}) string {
	var str string
	switch v := s.(type) {
	case null.String:
		if !v.Valid {
			return "NULL"
		}
		str = v.String
	case sql.NullString:
		if !v.Valid {
			return "NULL"
		}
		str = v.String
	case string:
		str = v
	default:
		panic("not a string or sql.NullString")
	}

	str = strings.ReplaceAll(str, "\\", "\\\\")
	return `"` + strings.ReplaceAll(str, "\"", "\\\"") + `"`
}

// Scan implements the Scanner interface.
//
// Note h is reallocated before the scan to clear existing values. If the
// hstore column's database value is NULL, then h is set to nil instead.
func (h *HStore) Scan(value interface{}) error {
	if value == nil {
		h = nil
		return nil
	}
	*h = make(map[string]null.String)
	var b byte
	pair := [][]byte{{}, {}}
	pi := 0
	inQuote := false
	didQuote := false
	sawSlash := false
	bindex := 0
	for bindex, b = range value.([]byte) {
		if sawSlash {
			pair[pi] = append(pair[pi], b)
			sawSlash = false
			continue
		}

		switch b {
		case '\\':
			sawSlash = true
			continue
		case '"':
			inQuote = !inQuote
			if !didQuote {
				didQuote = true
			}
			continue
		default:
			if !inQuote {
				switch b {
				case ' ', '\t', '\n', '\r':
					continue
				case '=':
					continue
				case '>':
					pi = 1
					didQuote = false
					continue
				case ',':
					s := string(pair[1])
					if !didQuote && len(s) == 4 && strings.EqualFold(s, "null") {
						(*h)[string(pair[0])] = null.String{String: "", Valid: false}
					} else {
						(*h)[string(pair[0])] = null.String{String: string(pair[1]), Valid: true}
					}
					pair[0] = []byte{}
					pair[1] = []byte{}
					pi = 0
					continue
				}
			}
		}
		pair[pi] = append(pair[pi], b)
	}
	if bindex > 0 {
		s := string(pair[1])
		if !didQuote && len(s) == 4 && strings.EqualFold(s, "null") {
			(*h)[string(pair[0])] = null.String{String: "", Valid: false}
		} else {
			(*h)[string(pair[0])] = null.String{String: string(pair[1]), Valid: true}
		}
	}
	return nil
}

// Value implements the driver Valuer interface. Note if h is nil, the
// database column value will be set to NULL.
func (h HStore) Value() (driver.Value, error) {
	if h == nil {
		return nil, nil
	}
	parts := []string{}
	for key, val := range h {
		thispart := hQuote(key) + "=>" + hQuote(val)
		parts = append(parts, thispart)
	}
	return []byte(strings.Join(parts, ",")), nil
}

// Randomize for sqlboiler
func (h *HStore) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		*h = nil
		return
	}

	*h = make(map[string]null.String)
	(*h)[randomize.Str(nextInt, 3)] = null.String{String: randomize.Str(nextInt, 3), Valid: nextInt()%3 == 0}
	(*h)[randomize.Str(nextInt, 3)] = null.String{String: randomize.Str(nextInt, 3), Valid: nextInt()%3 == 0}
}
