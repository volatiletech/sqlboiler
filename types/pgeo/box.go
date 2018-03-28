package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

//Boxes are represented by pairs of points that are opposite corners of the box.
type Box [2]Point

func (b Box) Value() (driver.Value, error) {
	return valueBox(b)
}

func (b *Box) Scan(src interface{}) error {
	return scanBox(b, src)
}

func valueBox(b Box) (driver.Value, error) {
	return fmt.Sprintf(`(%s)`, formatPoints(b[:])), nil
}

func scanBox(b *Box, src interface{}) error {
	if src == nil {
		*b = NewBox(Point{}, Point{})
		return nil
	}

	points, err := parsePointsSrc(src)
	if err != nil {
		return err
	}

	if len(points) != 2 {
		return errors.New("wrong box")
	}

	*b = NewBox(points[0], points[1])

	return nil
}
