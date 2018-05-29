package pgeo

import (
	"database/sql/driver"

	"github.com/volatiletech/sqlboiler/randomize"
)

// NullBox allows a box to be null
type NullBox struct {
	Box
	Valid bool `json:"valid"`
}

// Value for the database
func (b NullBox) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}

	return valueBox(b.Box)
}

// Scan from sql query
func (b *NullBox) Scan(src interface{}) error {
	if src == nil {
		b.Box, b.Valid = NewBox(Point{}, Point{}), false
		return nil
	}

	b.Valid = true
	return scanBox(&b.Box, src)
}

// Randomize for sqlboiler
func (b *NullBox) Randomize(seed *randomize.Seed, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		b.Valid = false
		return
	}

	b.Valid = true
	b.Box = randBox(seed)
}
