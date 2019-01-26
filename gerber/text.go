package gerber

import (
	"fmt"
	"io"
	"log"
)

const (
	defaultFont = "ubuntumonoregular"
	fsf         = 400
)

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y float64
	s    string
	font string
}

// Text returns a text primitive.
// All dimensions are in millimeters.
func Text(x, y float64, s, font string) *TextT {
	return &TextT{
		x:    x,
		y:    y,
		s:    s,
		font: font,
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
			x, y = t.x, y-float64(f.Ascent-f.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * float64(f.HorizAdvX)
			continue
		}
		g, ok := f.Glyphs[string(c)]
		if !ok {
			log.Printf("Warning: missing glyph %+q: skipping", c)
			x += float64(f.HorizAdvX)
			continue
		}
		dx := g.WriteGerber(w, apertureIndex, x, y)
		if dx == 0 {
			dx = float64(f.HorizAdvX)
		}
		x += dx
	}
	return nil
}

// Aperture returns nil for TextT because it uses the default aperture.
func (t *TextT) Aperture() *Aperture {
	return nil
}

// WriteGerber writes the primitive to the Gerber file.
func (g *Glyph) WriteGerber(w io.Writer, apertureIndex int, x, y float64) float64 {
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
		case "M":
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = oX+ps.P[0], oY+ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case "m":
			if len(pts) > 0 {
				dumpPoly()
			}
			x, y = x+ps.P[0], y+ps.P[1]
			pts = []Pt{{X: x, Y: y}}
		case "L":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = oX+ps.P[i], oY+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "l":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+ps.P[i], y+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "H":
			for i := 0; i < len(ps.P); i++ {
				x = oX + ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "h":
			for i := 0; i < len(ps.P); i++ {
				x += ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "V":
			for i := 0; i < len(ps.P); i++ {
				y = oY + ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "v":
			for i := 0; i < len(ps.P); i++ {
				y += ps.P[i]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "C":
		case "c":
		case "S":
		case "s":
		case "Q":
		case "q":

		case "T":
		case "t":
			for i := 0; i < len(ps.P); i += 2 {
				x, y = x+ps.P[i], y+ps.P[i+1]
				pts = append(pts, Pt{X: x, Y: y})
			}
		case "A":
		case "a":
		case "Z", "z":
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
	return float64(g.HorizAdvX)
}
