package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Polygon is represented by lists of points (the vertexes of the polygon).
type Polygon []Point

// Value for database
func (p Polygon) Value() (driver.Value, error) {
	return valuePolygon(p)
}

// Scan from sql query
func (p *Polygon) Scan(src interface{}) error {
	return scanPolygon(p, src)
}

func valuePolygon(p Polygon) (driver.Value, error) {
	return fmt.Sprintf(`(%s)`, formatPoints(p[:])), nil
}

func scanPolygon(p *Polygon, src interface{}) error {
	if src == nil {
		return nil
	}

	var err error
	*p, err = parsePointsSrc(src)
	if err != nil {
		return err
	}

	if len(*p) <= 2 {
		return errors.New("wrong polygon")
	}

	return nil
}

func randPolygon(nextInt func() int64) Polygon {
	return Polygon(randPoints(nextInt, 3))
}

// Randomize for sqlboiler
func (p *Polygon) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*p = randPolygon(nextInt)
}
