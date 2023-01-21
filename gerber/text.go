package gerber

import (
	"fmt"
	"io"
	"log"

	"github.com/gmlewis/go-fonts/fonts"
)

const (
	mmPerPt = 25.4 / 72.0

	XLeft   = fonts.XLeft
	XCenter = fonts.XCenter
	XRight  = fonts.XRight
	YBottom = fonts.YBottom
	YCenter = fonts.YCenter
	YTop    = fonts.YTop
)

var (
	BottomLeft   = fonts.BottomLeft
	BottomCenter = fonts.BottomCenter
	BottomRight  = fonts.BottomRight
	CenterLeft   = fonts.CenterLeft
	Center       = fonts.Center
	CenterRight  = fonts.CenterRight
	TopLeft      = fonts.TopLeft
	TopCenter    = fonts.TopCenter
	TopRight     = fonts.TopRight
)

// TextOpts provides options for positioning (aligning) the text based on
// its minimum bounding box.
type TextOpts = fonts.TextOpts

// TextT represents text and satisfies the Primitive interface.
type TextT struct {
	x, y     float64 // in mm
	xScale   float64
	opts     *TextOpts
	message  string
	fontName string
	pts      float64
	Render   *fonts.Render
}

func verifyOrSubstituteFont(fontName string) string {
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
	return fontName
}

// TextBox returns a text primitive where the text fills the MBB and aligns
// the text according to TextOpts.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func TextBox(mbb MBB, xScale float64, message, fontName string, opts *TextOpts) *TextT {
	fontName = verifyOrSubstituteFont(fontName)

	x, y, pts, err := fonts.FillBox(mbb, xScale, 1.0, message, fontName, opts)
	if err != nil {
		log.Fatalf("TextBox(message=%q,fontName=%q): %v", message, fontName, err)
	}
	return &TextT{
		x:        x,
		y:        y,
		xScale:   xScale,
		opts:     opts,
		message:  message,
		fontName: fontName,
		pts:      pts / mmPerPt,
	}
}

// Text returns a text primitive.
// All dimensions are in millimeters.
// xScale is 1.0 for top silkscreen and -1.0 for bottom silkscreen.
func Text(x, y, xScale float64, message, fontName string, pts float64, opts *TextOpts) *TextT {
	fontName = verifyOrSubstituteFont(fontName)

	return &TextT{
		x:        x,
		y:        y,
		xScale:   xScale,
		opts:     opts,
		message:  message,
		fontName: fontName,
		pts:      pts,
	}
}

func (t *TextT) renderText() error {
	if t.Render == nil {
		yScale := t.pts * mmPerPt
		xScale := t.xScale * yScale
		var err error
		if t.Render, err = fonts.Text(t.x, t.y, xScale, yScale, t.message, t.fontName, t.opts); err != nil {
			return err
		}
	}
	return nil
}

func (t *TextT) MBB() MBB {
	if err := t.renderText(); err != nil {
		log.Fatalf("MBB (message=%q,fontName=%q): %v", t.message, t.fontName, err)
	}
	return t.Render.MBB
}

// Width returns the width of the text in millimeters.
func (t *TextT) Width() float64 {
	if err := t.renderText(); err != nil {
		log.Fatalf("Width (message=%q,fontName=%q): %v", t.message, t.fontName, err)
	}
	width := t.Render.MBB.Max[0] - t.Render.MBB.Min[0]
	return width
}

// Height returns the height of the text in millimeters
func (t *TextT) Height() float64 {
	if err := t.renderText(); err != nil {
		log.Fatalf("Height (message=%q,fontName=%q): %v", t.message, t.fontName, err)
	}
	height := t.Render.MBB.Max[1] - t.Render.MBB.Min[1]
	return height
}

// WriteGerber writes the primitive to the Gerber file.
func (t *TextT) WriteGerber(w io.Writer, apertureIndex int) error {
	if err := t.renderText(); err != nil {
		return err
	}

	currentDark := true
	for _, poly := range t.Render.Polygons {
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
				fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*pt[0]), int(0.5+sf*pt[1]))
				continue
			}
			fmt.Fprintf(w, "X%06dY%06dD01*\n", int(0.5+sf*pt[0]), int(0.5+sf*pt[1]))
		}
		fmt.Fprintf(w, "X%06dY%06dD02*\n", int(0.5+sf*poly.Pts[0][0]), int(0.5+sf*poly.Pts[0][1]))
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
