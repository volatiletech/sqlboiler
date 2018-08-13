// Package pgeo implements geometric types for Postgres
//
// Geometric types:
// https://www.postgresql.org/docs/current/static/datatype-geometric.html
//
// Description:
// https://github.com/saulortega/pgeo
package pgeo

// NewPoint creates a point
func NewPoint(X, Y float64) Point {
	return Point{X, Y}
}

// NewLine creates a line
func NewLine(A, B, C float64) Line {
	return Line{A, B, C}
}

// NewLseg creates a line segment
func NewLseg(A, B Point) Lseg {
	return Lseg([2]Point{A, B})
}

// NewBox creates a box
func NewBox(A, B Point) Box {
	return Box([2]Point{A, B})
}

// NewPath creates a path
func NewPath(P []Point, C bool) Path {
	return Path{P, C}
}

// NewPolygon creates a polygon
func NewPolygon(P []Point) Polygon {
	return Polygon(P)
}

// NewCircle creates a circle from a radius and a point
func NewCircle(P Point, R float64) Circle {
	return Circle{P, R}
}

// NewNullPoint creates a point which can be null
func NewNullPoint(P Point, v bool) NullPoint {
	return NullPoint{P, v}
}

// NewNullLine creates a line which can be null
func NewNullLine(L Line, v bool) NullLine {
	return NullLine{L, v}
}

// NewNullLseg creates a line segment which can be null
func NewNullLseg(L Lseg, v bool) NullLseg {
	return NullLseg{L, v}
}

// NewNullBox creates a box which can be null
func NewNullBox(B Box, v bool) NullBox {
	return NullBox{B, v}
}

// NewNullPath creates a path which can be null
func NewNullPath(P Path, v bool) NullPath {
	return NullPath{P, v}
}

// NewNullPolygon creates a polygon which can be null
func NewNullPolygon(P Polygon, v bool) NullPolygon {
	return NullPolygon{P, v}
}

// NewNullCircle creates a circle which can be null
func NewNullCircle(C Circle, v bool) NullCircle {
	return NullCircle{C, v}
}
