package pgeo

import (
	"database/sql/driver"
)

type NullLine struct {
	Line
	Valid bool `json:"valid"`
}

func (l NullLine) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}

	return valueLine(l.Line)
}

func (l *NullLine) Scan(src interface{}) error {
	if src == nil {
		l.Line, l.Valid = NewLine(0, 0, 0), false
		return nil
	}

	l.Valid = true
	return scanLine(&l.Line, src)
}
