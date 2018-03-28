package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

//Polygons are represented by lists of points (the vertexes of the polygon).
type Polygon []Point

func (p Polygon) Value() (driver.Value, error) {
	return valuePolygon(p)
}

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
