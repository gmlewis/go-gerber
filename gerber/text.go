package gerber

import (
	"fmt"
	"io"
	"log"

	"github.com/gmlewis/go-fonts/fonts"
)

const (
	mmPerPt = 25.4 / 72.0
)

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y, xScale float64
	message      string
	fontName     string
	pts          float64
}

// Text returns a text primitive.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func Text(x, y, xScale float64, message, fontName string, pts float64) *TextT {
	if len(fonts.Fonts) == 0 {
		log.Fatal("No fonts available")
	}

	if _, ok := fonts.Fonts[fontName]; !ok {
		var name string
		for name = range fonts.Fonts {
			break
		}
		log.Printf("Could not find font %q: using %q instead", fontName, name)
		fontName = name
	}

	return &TextT{
		x:        sf * x,
		y:        sf * y,
		xScale:   xScale,
		message:  message,
		fontName: fontName,
		pts:      pts,
	}
}

// WriteGerber writes the primitive to the Gerber file.
func (t *TextT) WriteGerber(w io.Writer, apertureIndex int) error {
	yScale := sf * t.pts * mmPerPt
	xScale := t.xScale * yScale

	render, err := fonts.Text(t.x, t.y, xScale, yScale, t.message, t.fontName)
	if err != nil {
		return err
	}

	currentDark := true
	for _, poly := range render.Polygons {
		if poly.Dark && !currentDark {
			io.WriteString(w, "%LPD*%\n")
			currentDark = true
		} else if !poly.Dark && currentDark {
			io.WriteString(w, "%LPC*%\n")
			currentDark = false
		}

		io.WriteString(w, "G54D11*\n")
		io.WriteString(w, "G36*\n")
		for i, pt := range poly.Pts {
			if i == 0 {
				fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+pt.X), int(0.5+pt.Y))
				continue
			}
			fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+pt.X), int(0.5+pt.Y))
		}
		fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+poly.Pts[0].X), int(0.5+poly.Pts[0].Y))
		io.WriteString(w, "G37*\n")
	}

	if !currentDark {
		io.WriteString(w, "%LPD*%\n")
	}
	return nil
}

// Aperture returns nil for TextT because it uses the default aperture.
func (t *TextT) Aperture() *Aperture {
	return nil
}
