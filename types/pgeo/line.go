package pgeo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
)

// Line represents a infinite line with the linear equation Ax + By + C = 0, where A and B are not both zero.
type Line struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

// Value for database
func (l Line) Value() (driver.Value, error) {
	return valueLine(l)
}

// Scan from sql query
func (l *Line) Scan(src interface{}) error {
	return scanLine(l, src)
}

func valueLine(l Line) (driver.Value, error) {
	return fmt.Sprintf(`{%[1]v,%[2]v,%[3]v}`, l.A, l.B, l.C), nil
}

func scanLine(l *Line, src interface{}) error {
	if src == nil {
		*l = NewLine(0, 0, 0)
		return nil
	}

	val, err := iToS(src)
	if err != nil {
		return err
	}

	pdzs := regexp.MustCompile(`^\{(-?[0-9]+(?:\.[0-9]+)?),(-?[0-9]+(?:\.[0-9]+)?),(-?[0-9]+(?:\.[0-9]+)?)\}$`).FindStringSubmatch(val)
	if len(pdzs) != 4 {
		return errors.New("wrong line")
	}

	nums, err := parseNums(pdzs[1:4])
	if err != nil {
		return err
	}

	*l = NewLine(nums[0], nums[1], nums[2])

	return nil
}

func randLine(nextInt func() int64) Line {
	return Line{newRandNum(nextInt), newRandNum(nextInt), 0}
}

// Randomize for sqlboiler
func (l *Line) Randomize(nextInt func() int64, fieldType string, shouldBeNull bool) {
	*l = randLine(nextInt)
}
