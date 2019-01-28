package gerber

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gmlewis/go3d/float64/bezier2"
	"github.com/gmlewis/go3d/float64/qbezier2"
	"github.com/gmlewis/go3d/float64/vec2"
)

const (
	mmPerPt    = 25.4 / 72.0
	resolution = 0.1 // mm
	minSteps   = 4
	maxSteps   = 100
)

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y, xScale float64
	s            string
	font         *Font
	pts          float64
}

// Text returns a text primitive.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func Text(x, y, xScale float64, s, fontName string, pts float64) *TextT {
	if len(Fonts) == 0 {
		log.Fatal("No fonts available")
	}

	font, ok := Fonts[fontName]
	if !ok {
		var name string
		for name, font = range Fonts {
			break
		}
		log.Printf("Could not find font %q: using %q instead", fontName, name)
	}

	return &TextT{
		x:      x,
		y:      y,
		xScale: xScale,
		s:      s,
		font:   font,
		pts:    pts,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (t *TextT) WriteGerber(w io.Writer, apertureIndex int) error {
	x, y := t.x, t.y
	for _, c := range t.s {
		if c == rune('\n') {
			x, y = t.x, y-(t.font.Ascent-t.font.Descent)
			continue
		}
		if c == rune('\t') {
			x += 2.0 * t.xScale * t.font.HorizAdvX
			continue
		}
		g, ok := t.font.Glyphs[string(c)]
		if !ok {
			log.Printf("Warning: missing glyph %+q: skipping", c)
			x += t.xScale * t.font.HorizAdvX
			continue
		}
		dx := g.WriteGerber(w, apertureIndex, t, x, y)
		if dx == 0 {
			dx = t.font.HorizAdvX
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
func (g *Glyph) WriteGerber(w io.Writer, apertureIndex int, t *TextT, x, y float64) float64 {
	xScale := t.xScale
	oX, oY := x, y         // origin for this glyph
	var pts []Pt           // Current polygon
	currentPolarity := "d" // d=dark, c=clear
	var curveNum int

	fsf := sf * t.pts * mmPerPt / t.font.HorizAdvX

	dumpPoly := func() {
		if g.GerberLP != "" && curveNum < len(g.GerberLP) {
			polarity := g.GerberLP[curveNum : curveNum+1]
			if polarity != currentPolarity {
				fmt.Fprintf(w, "%%LP%v*%%\n", strings.ToUpper(polarity))
				currentPolarity = polarity
			}
			// } else if g.GerberLP == "" && curveNum > 0 && currentPolarity == "d" {
			// 	io.WriteString(w, "%LPC%*\n")
			// 	currentPolarity = "c"
		}

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

	var lastQ *qbezier2.T
	var lastCommand byte
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
			for i := 0; i < len(ps.P); i += 6 {
				x1, y1, x2, y2, ex, ey := oX+xScale*ps.P[i], oY+ps.P[i+1], oX+xScale*ps.P[i+2], oY+ps.P[i+3], oX+xScale*ps.P[i+4], oY+ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x1, y1},
					P2: vec2.T{x2, y2},
					P3: vec2.T{ex, ey},
				}
				// lastQ = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > maxSteps {
					steps = maxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = ex, ey
			}
		case 'c':
			for i := 0; i < len(ps.P); i += 6 {
				dx1, dy1, dx2, dy2, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3], xScale*ps.P[i+4], ps.P[i+5]
				b := &bezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx2, y + dy2},
					P3: vec2.T{x + dx, y + dy},
				}
				// lastQ = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > maxSteps {
					steps = maxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'S':
		// case 's':
		// case 'Q':
		case 'q':
			for i := 0; i < len(ps.P); i += 4 {
				dx1, dy1, dx, dy := xScale*ps.P[i], ps.P[i+1], xScale*ps.P[i+2], ps.P[i+3]
				b := &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				lastQ = b
				length := b.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > maxSteps {
					steps = maxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := b.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'T':
		case 't':
			for i := 0; i < len(ps.P); i += 2 {
				dx, dy := xScale*ps.P[i], ps.P[i+1]
				dx1, dy1 := 0.0, 0.0
				if lastQ != nil && (lastCommand == 'q' || lastCommand == 't') {
					dx1, dy1 = lastQ.P2[0]-lastQ.P1[0], lastQ.P2[1]-lastQ.P1[1]
				}
				lastQ = &qbezier2.T{
					P0: vec2.T{x, y},
					P1: vec2.T{x + dx1, y + dy1},
					P2: vec2.T{x + dx, y + dy},
				}
				lastCommand = ps.C
				length := lastQ.Length(1)
				steps := int(0.5 + length/resolution)
				if steps < minSteps {
					steps = minSteps
				}
				if steps > maxSteps {
					steps = maxSteps
				}
				for j := 1; j <= steps; j++ {
					t := float64(j) / float64(steps)
					p := lastQ.Point(t)
					pts = append(pts, Pt{X: p[0], Y: p[1]})
				}
				x, y = x+dx, y+dy
			}
		// case 'A':
		// case 'a':
		case 'Z', 'z':
			if len(pts) > 0 {
				pts = append(pts, pts[0]) // Close the path.
				dumpPoly()
			}
			curveNum++
		default:
			log.Fatalf("Unsupported path command %q", ps.C)
		}
		lastCommand = ps.C
	}
	if len(pts) > 0 {
		dumpPoly()
	}

	// Restore dark polarity for the rest of the Gerber layer.
	if currentPolarity != "d" {
		io.WriteString(w, "%LPD%*\n")
	}

	return g.HorizAdvX
}
