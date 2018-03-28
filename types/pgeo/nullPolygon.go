package pgeo

import (
	"database/sql/driver"
)

type NullPolygon struct {
	Polygon
	Valid bool `json:"valid"`
}

func (p NullPolygon) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}

	return valuePolygon(p.Polygon)
}

func (p *NullPolygon) Scan(src interface{}) error {
	if src == nil {
		p.Polygon, p.Valid = NewPolygon([]Point{Point{}, Point{}, Point{}, Point{}}), false
		return nil
	}

	p.Valid = true
	return scanPolygon(&p.Polygon, src)
}
