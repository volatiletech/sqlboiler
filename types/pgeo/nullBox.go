package pgeo

import (
	"database/sql/driver"
)

type NullBox struct {
	Box
	Valid bool `json:"valid"`
}

func (b NullBox) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}

	return valueBox(b.Box)
}

func (b *NullBox) Scan(src interface{}) error {
	if src == nil {
		b.Box, b.Valid = NewBox(Point{}, Point{}), false
		return nil
	}

	b.Valid = true
	return scanBox(&b.Box, src)
}
