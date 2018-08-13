package pgeo

import (
	"database/sql/driver"
)

// NullLine allows line to be null
type NullLine struct {
	Line
	Valid bool `json:"valid"`
}

// Value for database
func (l NullLine) Value() (driver.Value, error) {
	if !l.Valid {
		return nil, nil
	}

	return valueLine(l.Line)
}

// Scan from sql query
func (l *NullLine) Scan(src interface{}) error {
	if src == nil {
		l.Line, l.Valid = NewLine(0, 0, 0), false
		return nil
	}

	l.Valid = true
	return scanLine(&l.Line, src)
}

// Randomize for sqlboiler
func (l *NullLine) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		l.Valid = false
		return
	}

	l.Valid = true
	l.Line = randLine(nextInt)
}
