package pgeo

import (
	"database/sql/driver"
)

// Point is the fundamental two-dimensional building block for geometric types.
// X and Y are the respective coordinates, as floating-point numbers
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Value representation for database
func (p Point) Value() (driver.Value, error) {
	return valuePoint(p)
}

// Scan from query
func (p *Point) Scan(src interface{}) error {
	return scanPoint(p, src)
}

func valuePoint(p Point) (driver.Value, error) {
	return formatPoint(p), nil
}

func scanPoint(p *Point, src interface{}) error {
	if src == nil {
		*p = NewPoint(0, 0)
		return nil
	}

	val, err := iToS(src)
	if err != nil {
		return err
	}

	*p, err = parsePoint(val)
	if err != nil {
		return err
	}

	return nil

}

func randPoint(nextInt func() int64) Point {
	return Point{newRandNum(nextInt), newRandNum(nextInt)}
}

func randPoints(nextInt func() int64, n int) []Point {
	var points = []Point{}
	if n <= 0 {
		return points
	}

	for i := 0; i < n; i++ {
		points = append(points, randPoint(nextInt))
	}

	return points
}

// Randomize for sqlboiler
func (p *Point) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*p = randPoint(nextInt)
}
