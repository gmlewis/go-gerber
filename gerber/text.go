package gerber

import (
	"fmt"
	"io"
	"log"

	"github.com/gmlewis/go-fonts/fonts"
)

const (
	mmPerPt = 25.4 / 72.0

	XLeft   = 0
	XCenter = 0.5
	XRight  = 1
	YBottom = 0
	YCenter = 0.5
	YTop    = 1
)

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y, xScale   float64 // x, y in Gerber units (nm)
	xAlign, yAlign float64
	message        string
	fontName       string
	pts            float64
	render         *fonts.Render
}

// TextOpts provides options for positioning (aligning) the text based on
// its minimum bounding box.
type TextOpts struct {
	// XAlign represents the horizontal alignment of the text.
	// 0=x origin at left (the default), 1=x origin at right, 0.5=center.
	// XLeft, XCenter, and XRight are defined for convenience and
	// readability of the code.
	XAlign float64
	// YAlign represents the vertical alignment of the text.
	// 0=y origin at bottom (the default), 1=y origin at top, 0.5=center.
	// YBottom, YCenter, and YTop are defined for convenience and
	// readbility of the code.
	YAlign float64
}

// Text returns a text primitive.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func Text(x, y, xScale float64, message, fontName string, pts float64, opts *TextOpts) *TextT {
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

	var xAlign float64
	var yAlign float64
	if opts != nil {
		xAlign = opts.XAlign
		yAlign = opts.YAlign
	}
	return &TextT{
		x:        sf * x,
		y:        sf * y,
		xScale:   xScale,
		xAlign:   xAlign,
		yAlign:   yAlign,
		message:  message,
		fontName: fontName,
		pts:      pts,
	}
}

func (t *TextT) renderText() error {
	if t.render == nil {
		yScale := sf * t.pts * mmPerPt
		xScale := t.xScale * yScale
		// Get the MBB.
		render, err := fonts.Text(t.x, t.y, xScale, yScale, t.message, t.fontName)
		if err != nil {
			return err
		}
		// Re-render with MBB info.
		width := (render.Xmax - render.Xmin)
		height := (render.Ymax - render.Ymin)
		xError := render.Xmin - t.x
		yError := render.Ymin - t.y
		x := t.x - t.xAlign*width - xError
		y := t.y - t.yAlign*height - yError
		if render, err = fonts.Text(x, y, xScale, yScale, t.message, t.fontName); err != nil {
			return err
		}
		// log.Printf("t.message=%q t.x,t.y=(%.2f, %.2f), MBB=(%.2f,%.2f)-(%.2f,%.2f)", t.message, t.x, t.y, render.Xmin, render.Ymin, render.Xmax, render.Ymax)
		t.render = render
	}
	return nil
}

// Width returns the width of the text in millimeters.
func (t *TextT) Width() float64 {
	if err := t.renderText(); err != nil {
		log.Fatal(err)
	}
	width := t.render.Xmax - t.render.Xmin
	return width / sf
}

// Height returns the height of the text in millimeters
func (t *TextT) Height() float64 {
	if err := t.renderText(); err != nil {
		log.Fatal(err)
	}
	height := t.render.Ymax - t.render.Ymin
	return height / sf
}

// WriteGerber writes the primitive to the Gerber file.
func (t *TextT) WriteGerber(w io.Writer, apertureIndex int) error {
	if err := t.renderText(); err != nil {
		return err
	}

	currentDark := true
	for _, poly := range t.render.Polygons {
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
