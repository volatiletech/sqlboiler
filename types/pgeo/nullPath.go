package pgeo

import (
	"database/sql/driver"

	"github.com/volatiletech/sqlboiler/randomize"
)

// NullPath allows path to be null
type NullPath struct {
	Path
	Valid bool `json:"valid"`
}

// Value for database
func (p NullPath) Value() (driver.Value, error) {
	if !p.Valid {
		return nil, nil
	}

	return valuePath(p.Path)
}

// Scan from sql query
func (p *NullPath) Scan(src interface{}) error {
	if src == nil {
		p.Path, p.Valid = NewPath([]Point{Point{}, Point{}}, false), false
		return nil
	}

	p.Valid = true
	return scanPath(&p.Path, src)
}

// Randomize for sqlboiler
func (p *NullPath) Randomize(seed *randomize.Seed, fieldType string, shouldBeNull bool) {
	if shouldBeNull {
		p.Valid = false
		return
	}

	p.Valid = true
	p.Path = randPath(seed)
}
