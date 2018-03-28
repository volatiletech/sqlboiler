package pgeo

import (
	"database/sql/driver"
)

type NullPoint struct {
	Point
	Valid bool `json:"valid"`
}

func (p NullPoint) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}

	return valuePoint(p.Point)
}

func (p *NullPoint) Scan(src interface{}) error {
	if src == nil {
		p.Point, p.Valid = NewPoint(0, 0), false
		return nil
	}

	p.Valid = true
	return scanPoint(&p.Point, src)
}
