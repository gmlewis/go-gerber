package gerber

import "io"

// Shape represents the type of shape the primitives use.
type Shape int

const (
	// RectShape uses rectangles for the primitive.
	RectShape Shape = iota
	// CircleShape uses circles for the primitive.
	CircleShape
)

// Primitive is a Gerber primitive.
type Primitive interface {
	Write(w io.Writer) error
}

// Pt represents a 2D Point.
type Pt struct {
	X, Y float64
}

// Point is a simple convenience function that keeps the code
// easy to read.
// All dimensions are in millimeters.
func Point(x, y float64) Pt {
	return Pt{X: x, Y: y}
}

// ArcT represents an arc and satisfies the Primitive interface.
type ArcT struct {
}

// Arc returns an arc primitive.
// All dimensions are in millimeters. Angles are in radians.
func Arc(
	x, y, radius float64,
	shape Shape,
	xScale, yScale, startAngle, endAngle float64,
	layer int,
	component string,
	thickness float64) *ArcT {
	return &ArcT{}
}

func (p *ArcT) Write(w io.Writer) error {
	return nil
}

// CircleT represents a circle and satisfies the Primitive interface.
type CircleT struct {
}

// Circle returns a circle primitive.
// All dimensions are in millimeters. Angles are in radians.
func Circle(x, y float64, layer int, component string, thickness float64) *CircleT {
	return &CircleT{}
}

func (p *CircleT) Write(w io.Writer) error {
	return nil
}

// LineT represents a line and satisfies the Primitive interface.
type LineT struct {
}

// Line returns a line primitive.
// All dimensions are in millimeters.
func Line(x1, y1, x2, y2 float64, shape Shape, layer int, component string, thickness float64) *LineT {
	return &LineT{}
}

func (p *LineT) Write(w io.Writer) error {
	return nil
}

// PolygonT represents a polygon and satisfies the Primitive interface.
type PolygonT struct {
}

// Polygon returns a polygon primitive.
// All dimensions are in millimeters.
func Polygon(x, y float64, filled bool, points []Pt, layer int, component string, thickness float64) *PolygonT {
	return &PolygonT{}
}

func (p *PolygonT) Write(w io.Writer) error {
	return nil
}
