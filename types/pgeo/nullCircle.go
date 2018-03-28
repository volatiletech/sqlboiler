package pgeo

import (
	"database/sql/driver"
)

type NullCircle struct {
	Circle
	Valid bool `json:"valid"`
}

func (c NullCircle) Value() (driver.Value, error) {
	if !c.Valid {
		return nil, nil
	}

	return valueCircle(c.Circle)
}

func (c *NullCircle) Scan(src interface{}) error {
	if src == nil {
		c.Circle, c.Valid = NewCircle(Point{}, 0), false
		return nil
	}

	c.Valid = true
	return scanCircle(&c.Circle, src)
}
