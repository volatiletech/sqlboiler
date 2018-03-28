package pgeo

import (
	"database/sql/driver"
)

type NullLseg struct {
	Lseg
	Valid bool `json:"valid"`
}

func (l NullLseg) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}

	return valueLseg(l.Lseg)
}

func (l *NullLseg) Scan(src interface{}) error {
	if src == nil {
		l.Lseg, l.Valid = NewLseg(Point{}, Point{}), false
		return nil
	}

	l.Valid = true
	return scanLseg(&l.Lseg, src)
}
