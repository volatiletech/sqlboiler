// Package pgeo implements geometric types for Postgres
//
// Geometryc types:
// https://www.postgresql.org/docs/current/static/datatype-geometric.html
//
// Description:
// https://github.com/saulortega/pgeo
package pgeo

func NewPoint(X, Y float64) Point {
	return Point{X, Y}
}

func NewLine(A, B, C float64) Line {
	return Line{A, B, C}
}

func NewLseg(A, B Point) Lseg {
	return Lseg([2]Point{A, B})
}

func NewBox(A, B Point) Box {
	return Box([2]Point{A, B})
}

func NewPath(P []Point, C bool) Path {
	return Path{P, C}
}

func NewPolygon(P []Point) Polygon {
	return Polygon(P)
}

func NewCircle(P Point, R float64) Circle {
	return Circle{P, R}
}

func NewNullPoint(P Point, v bool) NullPoint {
	return NullPoint{P, v}
}

func NewNullLine(L Line, v bool) NullLine {
	return NullLine{L, v}
}

func NewNullLseg(L Lseg, v bool) NullLseg {
	return NullLseg{L, v}
}

func NewNullBox(B Box, v bool) NullBox {
	return NullBox{B, v}
}

func NewNullPath(P Path, v bool) NullPath {
	return NullPath{P, v}
}

func NewNullPolygon(P Polygon, v bool) NullPolygon {
	return NullPolygon{P, v}
}

func NewNullCircle(C Circle, v bool) NullCircle {
	return NullCircle{C, v}
}

func NewRandPoint() Point {
	return Point{newRandNum(), newRandNum()}
}

func NewRandLine() Line {
	return Line{newRandNum(), newRandNum(), 0}
}

func NewRandLseg() Lseg {
	return Lseg([2]Point{NewRandPoint(), NewRandPoint()})
}

func NewRandBox() Box {
	return Box([2]Point{NewRandPoint(), NewRandPoint()})
}

func NewRandPath() Path {
	return Path{RandPoints(3), newRandNum() < 40}
}

func NewRandPolygon() Polygon {
	return Polygon(RandPoints(3))
}

func NewRandCircle() Circle {
	return Circle{NewRandPoint(), newRandNum()}
}

func NewZeroPoint() Point {
	return Point{0, 0}
}

func RandPoints(n int) []Point {
	var points = []Point{}
	if n <= 0 {
		return points
	}

	for i := 0; i < n; i++ {
		points = append(points, NewRandPoint())
	}

	return points
}
