package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// Box is represented by pairs of points that are opposite corners of the box.
type Box [2]Point

// Value for the database
func (b Box) Value() (driver.Value, error) {
	return valueBox(b)
}

// Scan from sql query
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

func randBox(nextInt func() int64) Box {
	return Box([2]Point{randPoint(nextInt), randPoint(nextInt)})
}

// Randomize for sqlboiler
func (b *Box) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*b = randBox(nextInt)
}
