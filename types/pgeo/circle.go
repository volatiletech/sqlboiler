package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/volatiletech/sqlboiler/randomize"
)

// Circle is represented by a center point and radius.
type Circle struct {
	Point
	Radius float64 `json:"radius"`
}

// Value for the database
func (c Circle) Value() (driver.Value, error) {
	return valueCircle(c)
}

// Scan from sql query
func (c *Circle) Scan(src interface{}) error {
	return scanCircle(c, src)
}

func valueCircle(c Circle) (driver.Value, error) {
	return fmt.Sprintf(`<%s,%v>`, formatPoint(c.Point), c.Radius), nil
}

func scanCircle(c *Circle, src interface{}) error {
	if src == nil {
		*c = NewCircle(Point{}, 0)
		return nil
	}

	val, err := iToS(src)
	if err != nil {
		return err
	}

	points, err := parsePoints(val)
	if err != nil {
		return err
	}

	pdzs := strings.Split(val, "),")

	if len(points) != 1 || len(pdzs) != 2 {
		return errors.New("wrong circle")
	}

	r, err := strconv.ParseFloat(strings.Trim(pdzs[1], ">"), 64)
	if err != nil {
		return err
	}

	*c = NewCircle(points[0], r)

	return nil
}

func randCircle(seed *randomize.Seed) Circle {
	return Circle{randPoint(seed), newRandNum(seed)}
}

// Randomize for sqlboiler
func (c *Circle) Randomize(seed *randomize.Seed, fieldType string, shouldBeNull bool) {
	*c = randCircle(seed)
}
