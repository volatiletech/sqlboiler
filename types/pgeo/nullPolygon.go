package pgeo

import (
	"database/sql/driver"
)

// NullPolygon allows polygon to be null
type NullPolygon struct {
	Polygon
	Valid bool `json:"valid"`
}

// Value for database
func (p NullPolygon) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}

	return valuePolygon(p.Polygon)
}

// Scan from sql query
func (p *NullPolygon) Scan(src interface{}) error {
	if src == nil {
		p.Polygon, p.Valid = NewPolygon([]Point{Point{}, Point{}, Point{}, Point{}}), false
		return nil
	}

	p.Valid = true
	return scanPolygon(&p.Polygon, src)
}

// Randomize for sqlboiler
func (p *NullPolygon) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		p.Valid = false
		return
	}

	p.Valid = true
	p.Polygon = randPolygon(nextInt)
}
