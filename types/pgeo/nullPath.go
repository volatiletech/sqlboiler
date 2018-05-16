package pgeo

import (
	"database/sql/driver"
)

type NullPath struct {
	Path
	Valid bool `json:"valid"`
}

func (p NullPath) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}

	return valuePath(p.Path)
}

func (p *NullPath) Scan(src interface{}) error {
	if src == nil {
		p.Path, p.Valid = NewPath([]Point{Point{}, Point{}}, false), false
		return nil
	}

	p.Valid = true
	return scanPath(&p.Path, src)
}
