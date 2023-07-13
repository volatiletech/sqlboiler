package pgeo

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func iToS(src interface{}) (string, error) {
	var val string
	var err error

	switch src.(type) {
	case string:
		val = src.(string)
	case []byte:
		val = string(src.([]byte))
	default:
		err = fmt.Errorf("incompatible type %v", reflect.ValueOf(src).Kind().String())
	}

	return val, err
}

func parseNums(s []string) ([]float64, error) {
	var pts = []float64{}
	for _, p := range s {
		pt, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return pts, err
		}

		pts = append(pts, pt)
	}

	return pts, nil
}

func formatPoint(point Point) string {
	return fmt.Sprintf(`(%v,%v)`, point.X, point.Y)
}

func formatPoints(points []Point) string {
	var pts = []string{}
	for _, p := range points {
		pts = append(pts, formatPoint(p))
	}
	return strings.Join(pts, ",")
}

var parsePointRegexp = regexp.MustCompile(`^\(([0-9\.Ee-]+),([0-9\.Ee-]+)\)$`)

func parsePoint(pt string) (Point, error) {
	var point = Point{}
	var err error

	pdzs := parsePointRegexp.FindStringSubmatch(pt)
	if len(pdzs) != 3 {
		return point, errors.New("wrong point")
	}

	nums, err := parseNums(pdzs[1:3])
	if err != nil {
		return point, err
	}

	point.X = nums[0]
	point.Y = nums[1]

	return point, nil
}

var parsePointsRegexp = regexp.MustCompile(`\(([0-9\.Ee-]+),([0-9\.Ee-]+)\)`)

func parsePoints(pts string) ([]Point, error) {
	var points = []Point{}

	pdzs := parsePointsRegexp.FindAllString(pts, -1)
	for _, pt := range pdzs {
		point, err := parsePoint(pt)
		if err != nil {
			return points, err
		}

		points = append(points, point)
	}

	return points, nil
}

func parsePointsSrc(src interface{}) ([]Point, error) {
	val, err := iToS(src)
	if err != nil {
		return []Point{}, err
	}

	return parsePoints(val)
}

func newRandNum(nextInt func() int64) float64 {
	return float64(nextInt())
}
