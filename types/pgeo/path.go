package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
)

// Path is represented by lists of connected points.
// Paths can be open, where the first and last points in the list are considered not connected,
// or closed, where the first and last points are considered connected.
type Path struct {
	Points []Point
	Closed bool `json:"closed"`
}

// Value for database
func (p Path) Value() (driver.Value, error) {
	return valuePath(p)
}

// Scan from sql query
func (p *Path) Scan(src interface{}) error {
	return scanPath(p, src)
}

func valuePath(p Path) (driver.Value, error) {
	var val string
	if p.Closed {
		val = fmt.Sprintf(`(%s)`, formatPoints(p.Points))
	} else {
		val = fmt.Sprintf(`[%s]`, formatPoints(p.Points))
	}
	return val, nil
}

func scanPath(p *Path, src interface{}) error {
	if src == nil {
		return nil
	}

	val, err := iToS(src)
	if err != nil {
		return err
	}

	(*p).Points, err = parsePoints(val)
	if err != nil {
		return err
	}

	if len((*p).Points) < 2 {
		return errors.New("wrong path")
	}

	(*p).Closed = regexp.MustCompile(`^\(\(`).MatchString(val)

	return nil
}

func randPath(nextInt func() int64) Path {
	return Path{randPoints(nextInt, 3), newRandNum(nextInt) < 40}
}

// Randomize for sqlboiler
func (p *Path) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*p = randPath(nextInt)
}
