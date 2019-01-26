package gerber

import (
	"fmt"
	"io"
	"log"

	"github.com/ungerik/go3d/float64/hermit2"
	"github.com/ungerik/go3d/float64/vec2"
)

const (
	defaultFont = "ubuntumonoregular"
	fsf         = 600
	resolution  = 1000
)

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y, xScale float64
	s            string
	font         string
}

// Text returns a text primitive.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func Text(x, y, xScale float64, s, font string) *TextT {
	return &TextT{
		x:      x,
		y:      y,
		xScale: xScale,
		s:      s,
		font:   font,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (t *TextT) WriteGerber(w io.Writer, apertureIndex int) error {
	f, ok := fonts[t.font]
	if !ok {
		log.Printf("Could not find font %q: using %q instead", t.font, defaultFont)
		f, ok = fonts[defaultFont]
		if !ok {
			log.Fatalf("Could not find default font %q", defaultFont)
		}
	}

	x, y := t.x, t.y
	for _, c := range t.s {
		if c == rune('\n') {
			x, y = t.x, y-(f.Ascent-f.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * t.xScale * f.HorizAdvX
			continue
		}
		g, ok := f.Glyphs[string(c)]
		if !ok {
			log.Printf("Warning: missing glyph %+q: skipping", c)
			x += t.xScale * f.HorizAdvX
			continue
		}
		dx := g.WriteGerber(w, apertureIndex, x, y, t.xScale)
		if dx == 0 {
			dx = f.HorizAdvX
		}
		x += dx * t.xScale
	}
	return nil
}

// Aperture returns nil for TextT because it uses the default aperture.
func (t *TextT) Aperture() *Aperture {
	return nil
}

// WriteGerber writes the primitive to the Gerber file.
func (g *Glyph) WriteGerber(w io.Writer, apertureIndex int, x, y, xScale float64) float64 {
	oX, oY := x, y // origin for this glyph
	var pts []Pt   // Current polygon

	dumpPoly := func() {
		io.WriteString(w, "G54D11*\n")
		io.WriteString(w, "G36*\n")
		for i, pt := range pts {
			if i == 0 {
				fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+fsf*pt.X), int(0.5+fsf*pt.Y))
				continue
			}
			fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+fsf*pt.X), int(0.5+fsf*pt.Y))
		}
		fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+fsf*pts[0].X), int(0.5+fsf*pts[0].Y))
		io.WriteString(w, "G37*\n")
		pts = []Pt{}
	}

	for _, ps := range g.PathSteps {
		switch ps.C {
		case 'M':
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = oX+xScale*ps.P[0], oY+ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case 'm':
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = x+xScale*ps.P[0], y+ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case 'L':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = oX+xScale*ps.P[i], oY+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'l':
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+xScale*ps.P[i], y+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'H':
			for i := 0; i < len(ps.P); i++ {
				x = oX + xScale*ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'h':
			for i := 0; i < len(ps.P); i++ {
				x += xScale * ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'V':
			for i := 0; i < len(ps.P); i++ {
				y = oY + ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'v':
			for i := 0; i < len(ps.P); i++ {
				y += ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case 'C':
			log.Printf("C: %#v", ps.P)
		case 'c':
			log.Printf("c: %#v", ps.P)
		// case 'S':
		// case 's':
		// case 'Q':
		case 'q':
			for i := 0; i < len(ps.P); i += 4 {
				dx1, dy1, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3]
				h := &hermit2.T{
					A: hermit2.PointTangent{Point: vec2.T{x, y}, Tangent: vec2.T{dx1, dy1}},
					B: hermit2.PointTangent{Point: vec2.T{x + dx, y + dy}, Tangent: vec2.T{-dx1, -dy1}},
				}
				length := h.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < 3 {
					steps = 3
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := h.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'T':
		case 't':
			log.Printf("t: %#v", ps.P)
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+xScale*ps.P[i], y+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		// case 'A':
		// case 'a':
		case 'Z', 'z':
			if len(pts) > 0 {
				dumpPoly()
			}
		default:
			log.Fatalf("Unknown path command %q", ps.C)
		}
	}
	if len(pts) > 0 {
		dumpPoly()
	}
	return g.HorizAdvX
}
