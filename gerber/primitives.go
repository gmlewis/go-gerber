package gerber

import (
	"fmt"
	"io"
	"math"

	"github.com/gmlewis/go3d/float64/vec2"
)

const (
	sf     = 1e6 // scale factor
	maxPts = 10000
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
	// MBB returns the minimum bounding box in millimeters.
	MBB() MBB
}

// Aperture represents the nature of the primitive
// and satisfies the Primitive interface.
type Aperture struct {
	Shape Shape
	Size  float64
}

func (a *Aperture) MBB() MBB { return MBB{} }

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
type Pt = vec2.T

// MBB represents a minimum bounding box.
type MBB = vec2.Rect

// Point is a simple convenience function that keeps the code easy to read.
// All dimensions are in millimeters.
func Point(x, y float64) Pt {
	return Pt{x, y}
}

// ArcT represents an arc and satisfies the Primitive interface.
type ArcT struct {
	center     Pt
	radius     float64
	shape      Shape
	xScale     float64
	yScale     float64
	startAngle float64
	endAngle   float64
	thickness  float64
	mbb        *MBB // cached minimum bounding box
}

// Arc returns an arc primitive.
// All dimensions are in millimeters. Angles are in degrees.
func Arc(
	center Pt,
	radius float64,
	shape Shape,
	xScale, yScale, startAngle, endAngle float64,
	thickness float64) *ArcT {
	if startAngle > endAngle {
		startAngle, endAngle = endAngle, startAngle
	}
	return &ArcT{
		center:     center,
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
		x1 := a.center[0] + a.xScale*math.Cos(angle)*a.radius
		y1 := a.center[1] + a.yScale*math.Sin(angle)*a.radius

		angle += delta

		x2 := a.center[0] + a.xScale*math.Cos(angle)*a.radius
		y2 := a.center[1] + a.yScale*math.Sin(angle)*a.radius

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

func (a *ArcT) MBB() MBB {
	if a.mbb != nil {
		return *a.mbb
	}

	delta := a.endAngle - a.startAngle
	length := delta * a.radius
	// Resolution of segments is 0.1mm
	segments := int(0.5+length*10.0) + 1
	delta /= float64(segments)

	angle := float64(a.startAngle)
	for i := 0; i < segments; i++ {
		x1 := a.center[0] + a.xScale*math.Cos(angle)*a.radius
		y1 := a.center[1] + a.yScale*math.Sin(angle)*a.radius

		angle += delta

		x2 := a.center[0] + a.xScale*math.Cos(angle)*a.radius
		y2 := a.center[1] + a.yScale*math.Sin(angle)*a.radius

		line := Line(x1, y1, x2, y2, a.shape, a.thickness)
		mbb := line.MBB()
		if a.mbb == nil {
			a.mbb = &mbb
		} else {
			a.mbb.Join(&mbb)
		}
	}

	return *a.mbb
}

// CircleT represents a circle and satisfies the Primitive interface.
type CircleT struct {
	pt        Pt
	thickness float64
	mbb       *MBB // cached minimum bounding box
}

// Circle returns a circle primitive.
// All dimensions are in millimeters.
func Circle(center Pt, thickness float64) *CircleT {
	return &CircleT{
		pt:        center,
		thickness: thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (c *CircleT) WriteGerber(w io.Writer, apertureIndex int) error {
	fmt.Fprintf(w, "G54D%d*\n", apertureIndex)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(c.pt[0])), int(0.5+sf*(c.pt[1])))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(c.pt[0])), int(0.5+sf*(c.pt[1])))
	return nil
}

// Aperture returns the primitive's desired aperture.
func (c *CircleT) Aperture() *Aperture {
	return &Aperture{
		Shape: CircleShape,
		Size:  c.thickness,
	}
}

func (c *CircleT) MBB() MBB {
	if c.mbb != nil {
		return *c.mbb
	}
	r := 0.5 * c.thickness
	ll := Pt{c.pt[0] - r, c.pt[1] - r}
	ur := Pt{c.pt[0] + r, c.pt[1] + r}
	c.mbb = &MBB{Min: ll, Max: ur}
	return *c.mbb
}

// LineT represents a line and satisfies the Primitive interface.
type LineT struct {
	p1, p2    Pt
	shape     Shape
	thickness float64
	mbb       *MBB    // cached minimum bounding box
	length    float64 // cached length of the line
}

// Line returns a line primitive.
// All dimensions are in millimeters.
func Line(x1, y1, x2, y2 float64, shape Shape, thickness float64) *LineT {
	p1 := Pt{x1, y1}
	p2 := Pt{x2, y2}
	v := vec2.Sub(&p1, &p2)
	length := v.Length()
	return &LineT{
		p1:        p1,
		p2:        p2,
		shape:     shape,
		thickness: thickness,
		length:    length,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (l *LineT) WriteGerber(w io.Writer, apertureIndex int) error {
	fmt.Fprintf(w, "G54D%d*\n", apertureIndex)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(l.p1[0])), int(0.5+sf*(l.p1[1])))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(l.p2[0])), int(0.5+sf*(l.p2[1])))
	return nil
}

// Aperture returns the primitive's desired aperture.
func (l *LineT) Aperture() *Aperture {
	return &Aperture{
		Shape: l.shape,
		Size:  l.thickness,
	}
}

func (l *LineT) MBB() MBB {
	if l.mbb != nil {
		return *l.mbb
	}
	l.mbb = &MBB{Min: l.p1, Max: l.p1}
	l.mbb.Join(&MBB{Min: l.p2, Max: l.p2})
	l.mbb.Min[0] -= 0.5 * l.thickness
	l.mbb.Min[1] -= 0.5 * l.thickness
	l.mbb.Max[0] += 0.5 * l.thickness
	l.mbb.Max[1] += 0.5 * l.thickness
	return *l.mbb
}

// PolygonT represents a polygon and satisfies the Primitive interface.
type PolygonT struct {
	offset Pt
	points []Pt
	mbb    *MBB // cached minimum bounding box
}

// Polygon returns a polygon primitive.
// All dimensions are in millimeters.
func Polygon(offset Pt, filled bool, points []Pt, thickness float64) *PolygonT {
	return &PolygonT{
		offset: offset,
		points: points,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (p *PolygonT) WriteGerber(w io.Writer, apertureIndex int) error {
	io.WriteString(w, "G54D11*\n")
	io.WriteString(w, "G36*\n")
	for i, pt := range p.points {
		if i == 0 {
			fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(pt[0]+p.offset[0])), int(0.5+sf*(pt[1]+p.offset[1])))
			continue
		}
		fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(pt[0]+p.offset[0])), int(0.5+sf*(pt[1]+p.offset[1])))
	}
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(p.points[0][0]+p.offset[0])), int(0.5+sf*(p.points[0][1]+p.offset[1])))
	io.WriteString(w, "G37*\n")
	return nil
}

// Aperture returns nil for PolygonT because it uses the default aperture.
func (p *PolygonT) Aperture() *Aperture {
	return nil
}

func (p *PolygonT) MBB() MBB {
	if p.mbb != nil {
		return *p.mbb
	}
	for i, pt := range p.points {
		newPt := Pt{pt[0] + p.offset[0], pt[1] + p.offset[1]}
		v := &MBB{Min: newPt, Max: newPt}
		if i == 0 {
			p.mbb = v
			continue
		}
		p.mbb.Join(v)
	}

	return *p.mbb
}
