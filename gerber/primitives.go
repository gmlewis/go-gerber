package gerber

import (
	"fmt"
	"io"
	"math"
)

const (
	sf = 1e6 // scale factor
)

// Shape represents the type of shape the apertures use.
type Shape string

const (
	// RectShape uses rectangles for the aperture.
	RectShape Shape = "R"
	// CircleShape uses circles for the aperture.
	CircleShape Shape = "C"
)

// Primitive is a Gerber primitive.
type Primitive interface {
	WriteGerber(w io.Writer, apertureIndex int) error
	Aperture() *Aperture
}

// Aperture represents the nature of the primitive
// and satisfies the Primitive interface.
type Aperture struct {
	Shape Shape
	Size  float64
}

// WriteGerber writes the aperture to the Gerber file.
func (a *Aperture) WriteGerber(w io.Writer, apertureIndex int) error {
	if a.Shape == CircleShape {
		fmt.Fprintf(w, "%%ADD%vC,%0.5f*%%\n", apertureIndex, a.Size)
		return nil
	}
	fmt.Fprintf(w, "%%ADD%vR,%0.5fX%0.5f*%%\n", apertureIndex, a.Size, a.Size)
	return nil
}

// Aperture is defined to implement the Primitive interface.
func (a *Aperture) Aperture() *Aperture {
	return a
}

// ID returns a unique ID for the Aperture.
func (a *Aperture) ID() string {
	if a == nil {
		return "default"
	}
	return fmt.Sprintf("%v%0.5f", a.Shape, sf*a.Size)
}

// Pt represents a 2D Point.
type Pt struct {
	X, Y float64
}

// Point is a simple convenience function that keeps the code easy to read.
// All dimensions are in millimeters.
func Point(x, y float64) Pt {
	return Pt{X: x, Y: y}
}

// ArcT represents an arc and satisfies the Primitive interface.
type ArcT struct {
	x          float64
	y          float64
	radius     float64
	shape      Shape
	xScale     float64
	yScale     float64
	startAngle float64
	endAngle   float64
	thickness  float64
}

// Arc returns an arc primitive.
// All dimensions are in millimeters. Angles are in degrees.
func Arc(
	x, y, radius float64,
	shape Shape,
	xScale, yScale, startAngle, endAngle float64,
	thickness float64) *ArcT {
	if startAngle > endAngle {
		startAngle, endAngle = endAngle, startAngle
	}
	return &ArcT{
		x:          x,
		y:          y,
		radius:     radius,
		shape:      shape,
		xScale:     math.Abs(xScale),
		yScale:     math.Abs(yScale),
		startAngle: math.Pi * startAngle / 180.0,
		endAngle:   math.Pi * endAngle / 180.0,
		thickness:  thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (a *ArcT) WriteGerber(w io.Writer, apertureIndex int) error {
	delta := a.endAngle - a.startAngle
	length := delta * a.radius
	// Resolution of segments is 0.1mm
	segments := int(0.5+length*10.0) + 1
	delta /= float64(segments)

	angle := float64(a.startAngle)
	for i := 0; i < segments; i++ {
		x1 := a.x + a.xScale*math.Cos(angle)*a.radius
		y1 := a.y + a.yScale*math.Sin(angle)*a.radius

		angle += delta

		x2 := a.x + a.xScale*math.Cos(angle)*a.radius
		y2 := a.y + a.yScale*math.Sin(angle)*a.radius

		line := Line(x1, y1, x2, y2, a.shape, a.thickness)
		line.WriteGerber(w, apertureIndex)
	}
	return nil
}

// Aperture returns the primitive's desired aperture.
func (a *ArcT) Aperture() *Aperture {
	return &Aperture{
		Shape: a.shape,
		Size:  a.thickness,
	}
}

// CircleT represents a circle and satisfies the Primitive interface.
type CircleT struct {
	x, y      float64
	thickness float64
}

// Circle returns a circle primitive.
// All dimensions are in millimeters.
func Circle(x, y float64, thickness float64) *CircleT {
	return &CircleT{
		x:         x,
		y:         y,
		thickness: thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (c *CircleT) WriteGerber(w io.Writer, apertureIndex int) error {
	fmt.Fprintf(w, "G54D%d*\n", apertureIndex)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(c.x)), int(0.5+sf*(c.y)))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(c.x)), int(0.5+sf*(c.y)))
	return nil
}

// Aperture returns the primitive's desired aperture.
func (c *CircleT) Aperture() *Aperture {
	return &Aperture{
		Shape: CircleShape,
		Size:  c.thickness,
	}
}

// LineT represents a line and satisfies the Primitive interface.
type LineT struct {
	x1, y1    float64
	x2, y2    float64
	shape     Shape
	thickness float64
}

// Line returns a line primitive.
// All dimensions are in millimeters.
func Line(x1, y1, x2, y2 float64, shape Shape, thickness float64) *LineT {
	return &LineT{
		x1:        x1,
		y1:        y1,
		x2:        x2,
		y2:        y2,
		shape:     shape,
		thickness: thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (l *LineT) WriteGerber(w io.Writer, apertureIndex int) error {
	fmt.Fprintf(w, "G54D%d*\n", apertureIndex)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(l.x1)), int(0.5+sf*(l.y1)))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(l.x2)), int(0.5+sf*(l.y2)))
	return nil
}

// Aperture returns the primitive's desired aperture.
func (l *LineT) Aperture() *Aperture {
	return &Aperture{
		Shape: l.shape,
		Size:  l.thickness,
	}
}

// PolygonT represents a polygon and satisfies the Primitive interface.
type PolygonT struct {
	x, y   float64
	points []Pt
}

// Polygon returns a polygon primitive.
// All dimensions are in millimeters.
func Polygon(x, y float64, filled bool, points []Pt, thickness float64) *PolygonT {
	return &PolygonT{
		x:      x,
		y:      y,
		points: points,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (p *PolygonT) WriteGerber(w io.Writer, apertureIndex int) error {
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

// Aperture returns nil for PolygonT because it uses the default aperture.
func (p *PolygonT) Aperture() *Aperture {
	return nil
}
