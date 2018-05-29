package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/randomize"
)

// Lseg is a line segment and is represented by pairs of points that are the endpoints of the segment.
type Lseg [2]Point

// Value for the database
func (l Lseg) Value() (driver.Value, error) {
	return valueLseg(l)
}

// Scan from sql query
func (l *Lseg) Scan(src interface{}) error {
	return scanLseg(l, src)
}

func valueLseg(l Lseg) (driver.Value, error) {
	return fmt.Sprintf(`[%s]`, formatPoints(l[:])), nil
}

func scanLseg(l *Lseg, src interface{}) error {
	if src == nil {
		*l = NewLseg(Point{}, Point{})
		return nil
	}

	points, err := parsePointsSrc(src)
	if err != nil {
		return err
	}

	if len(points) != 2 {
		return errors.New("wrong lseg")
	}

	*l = NewLseg(points[0], points[1])

	return nil
}

func randLseg(seed *randomize.Seed) Lseg {
	return Lseg([2]Point{randPoint(seed), randPoint(seed)})
}

// Randomize for sqlboiler
func (l *Lseg) Randomize(seed *randomize.Seed, fieldType string, shouldBeNull bool) {
	*l = randLseg(seed)
}
