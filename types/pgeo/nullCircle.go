package pgeo

import (
	"database/sql/driver"
)

// NullCircle allows circle to be null
type NullCircle struct {
	Circle
	Valid bool `json:"valid"`
}

// Value for database
func (c NullCircle) Value() (driver.Value, error) {
	if !c.Valid {
		return nil, nil
	}

	return valueCircle(c.Circle)
}

// Scan from sql query
func (c *NullCircle) Scan(src interface{}) error {
	if src == nil {
		c.Circle, c.Valid = NewCircle(Point{}, 0), false
		return nil
	}

	c.Valid = true
	return scanCircle(&c.Circle, src)
}

// Randomize for sqlboiler
func (c *NullCircle) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		c.Valid = false
		return
	}

	c.Valid = true
	c.Circle = randCircle(nextInt)
}
