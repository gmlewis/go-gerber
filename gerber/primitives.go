package gerber

import (
	"fmt"
	"io"
)

const (
	sf = 1e6 // scale factor
)

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
	WriteGerber(w io.Writer) error
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

func (p *ArcT) WriteGerber(w io.Writer) error {
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

func (p *CircleT) WriteGerber(w io.Writer) error {
	return nil
}

// LineT represents a line and satisfies the Primitive interface.
type LineT struct {
	x1, y1   float64
	x2, y2   float64
	aperture int
}

// Line returns a line primitive.
// All dimensions are in millimeters.
func Line(x1, y1, x2, y2 float64, shape Shape, layer int, component string, thickness float64) *LineT {
	return &LineT{
		x1: x1,
		y1: y1,
		x2: x2,
		y2: y2,
	}
}

func (p *LineT) WriteGerber(w io.Writer) error {
	fmt.Fprintf(w, "G54D%d*\n", p.aperture+12)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(p.x1)), int(0.5+sf*(p.y1)))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(p.x2)), int(0.5+sf*(p.y2)))
	return nil
}

// PolygonT represents a polygon and satisfies the Primitive interface.
type PolygonT struct {
	x, y   float64
	points []Pt
}

// Polygon returns a polygon primitive.
// All dimensions are in millimeters.
func Polygon(x, y float64, filled bool, points []Pt, layer int, component string, thickness float64) *PolygonT {
	return &PolygonT{
		x:      x,
		y:      y,
		points: points,
	}
}

func (p *PolygonT) WriteGerber(w io.Writer) error {
	io.WriteString(w, "G54D11*\n")
	io.WriteString(w, "G36*\n")
	for i, pt := range p.points {
		if i == 0 {
			fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(pt.X+p.x)), int(0.5+sf*(pt.Y+p.y)))
			continue
		}
		fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(pt.X+p.x)), int(0.5+sf*(pt.Y+p.y)))
	}
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(p.points[0].X+p.x)), int(0.5+sf*(p.points[0].Y+p.y)))
	io.WriteString(w, "G37*\n")
	return nil
}
