package pgeo

import (
	"database/sql/driver"
)

// NullLseg allows line segment to be null
type NullLseg struct {
	Lseg
	Valid bool `json:"valid"`
}

// Value for database
func (l NullLseg) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}

	return valueLseg(l.Lseg)
}

// Scan from sql query
func (l *NullLseg) Scan(src interface{}) error {
	if src == nil {
		l.Lseg, l.Valid = NewLseg(Point{}, Point{}), false
		return nil
	}

	l.Valid = true
	return scanLseg(&l.Lseg, src)
}

// Randomize for sqlboiler
func (l *NullLseg) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		l.Valid = false
		return
	}

	l.Valid = true
	l.Lseg = randLseg(nextInt)
}
