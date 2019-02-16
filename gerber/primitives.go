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
	Center     Pt
	Radius     float64
	Shape      Shape
	XScale     float64
	YScale     float64
	StartAngle float64
	EndAngle   float64
	Thickness  float64
	mbb        *MBB // cached minimum bounding box
}

// Arc returns an arc primitive.
// All dimensions are in millimeters.
// Angles are specified in degrees (and stored as radians).
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
		Center:     center,
		Radius:     radius,
		Shape:      shape,
		XScale:     math.Abs(xScale),
		YScale:     math.Abs(yScale),
		StartAngle: math.Pi * startAngle / 180.0,
		EndAngle:   math.Pi * endAngle / 180.0,
		Thickness:  thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (a *ArcT) WriteGerber(w io.Writer, apertureIndex int) error {
	delta := a.EndAngle - a.StartAngle
	length := delta * a.Radius
	// Resolution of segments is 0.1mm
	segments := int(0.5+length*10.0) + 1
	delta /= float64(segments)

	angle := float64(a.StartAngle)
	for i := 0; i < segments; i++ {
		x1 := a.Center[0] + a.XScale*math.Cos(angle)*a.Radius
		y1 := a.Center[1] + a.YScale*math.Sin(angle)*a.Radius

		angle += delta

		x2 := a.Center[0] + a.XScale*math.Cos(angle)*a.Radius
		y2 := a.Center[1] + a.YScale*math.Sin(angle)*a.Radius

		line := Line(x1, y1, x2, y2, a.Shape, a.Thickness)
		line.WriteGerber(w, apertureIndex)
	}
	return nil
}

// Aperture returns the primitive's desired aperture.
func (a *ArcT) Aperture() *Aperture {
	return &Aperture{
		Shape: a.Shape,
		Size:  a.Thickness,
	}
}

func (a *ArcT) MBB() MBB {
	if a.mbb != nil {
		return *a.mbb
	}

	delta := a.EndAngle - a.StartAngle
	length := delta * a.Radius
	// Resolution of segments is 0.1mm
	segments := int(0.5+length*10.0) + 1
	delta /= float64(segments)

	angle := float64(a.StartAngle)
	for i := 0; i < segments; i++ {
		x1 := a.Center[0] + a.XScale*math.Cos(angle)*a.Radius
		y1 := a.Center[1] + a.YScale*math.Sin(angle)*a.Radius

		angle += delta

		x2 := a.Center[0] + a.XScale*math.Cos(angle)*a.Radius
		y2 := a.Center[1] + a.YScale*math.Sin(angle)*a.Radius

		line := Line(x1, y1, x2, y2, a.Shape, a.Thickness)
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
	P1, P2    Pt
	Shape     Shape
	Thickness float64
	mbb       *MBB // cached minimum bounding box
}

// Line returns a line primitive.
// All dimensions are in millimeters.
func Line(x1, y1, x2, y2 float64, shape Shape, thickness float64) *LineT {
	p1 := Pt{x1, y1}
	p2 := Pt{x2, y2}
	return &LineT{
		P1:        p1,
		P2:        p2,
		Shape:     shape,
		Thickness: thickness,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (l *LineT) WriteGerber(w io.Writer, apertureIndex int) error {
	fmt.Fprintf(w, "G54D%d*\n", apertureIndex)
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(l.P1[0])), int(0.5+sf*(l.P1[1])))
	fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(l.P2[0])), int(0.5+sf*(l.P2[1])))
	return nil
}

// Aperture returns the primitive's desired aperture.
func (l *LineT) Aperture() *Aperture {
	return &Aperture{
		Shape: l.Shape,
		Size:  l.Thickness,
	}
}

func (l *LineT) MBB() MBB {
	if l.mbb != nil {
		return *l.mbb
	}
	l.mbb = &MBB{Min: l.P1, Max: l.P1}
	l.mbb.Join(&MBB{Min: l.P2, Max: l.P2})
	l.mbb.Min[0] -= 0.5 * l.Thickness
	l.mbb.Min[1] -= 0.5 * l.Thickness
	l.mbb.Max[0] += 0.5 * l.Thickness
	l.mbb.Max[1] += 0.5 * l.Thickness
	return *l.mbb
}

// PolygonT represents a polygon and satisfies the Primitive interface.
type PolygonT struct {
	Offset Pt
	Points []Pt
	mbb    *MBB // cached minimum bounding box
}

// Polygon returns a polygon primitive.
// All dimensions are in millimeters.
func Polygon(offset Pt, filled bool, points []Pt, thickness float64) *PolygonT {
	return &PolygonT{
		Offset: offset,
		Points: points,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (p *PolygonT) WriteGerber(w io.Writer, apertureIndex int) error {
	io.WriteString(w, "G54D11*\n")
	io.WriteString(w, "G36*\n")
	for i, pt := range p.Points {
		if i == 0 {
			fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(pt[0]+p.Offset[0])), int(0.5+sf*(pt[1]+p.Offset[1])))
			continue
		}
		fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*(pt[0]+p.Offset[0])), int(0.5+sf*(pt[1]+p.Offset[1])))
	}
	fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*(p.Points[0][0]+p.Offset[0])), int(0.5+sf*(p.Points[0][1]+p.Offset[1])))
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
	for i, pt := range p.Points {
		newPt := Pt{pt[0] + p.Offset[0], pt[1] + p.Offset[1]}
		v := &MBB{Min: newPt, Max: newPt}
		if i == 0 {
			p.mbb = v
			continue
		}
		p.mbb.Join(v)
	}

	return *p.mbb
}
