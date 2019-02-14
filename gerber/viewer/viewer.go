// Package viewer views a Gerber design using Fyne.
package viewer

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/gmlewis/go-gerber/gerber"
)

type viewController struct {
	g *gerber.Gerber
}

func Gerber(g *gerber.Gerber) {
	a := app.New()

	vc := &viewController{g: g}
	c := canvas.NewRaster(vc.pixelFunc)
	c.SetMinSize(fyne.Size{Width: 800, Height: 800})

	layers := widget.NewVBox(
		widget.NewHBox(
			widget.NewCheck("Top Silkscreen", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Top", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Layer 2", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Layer 3", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Layer 4", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Layer 5", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Bottom", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Bottom Silkscreen", func(v bool) {}),
			layout.NewSpacer(),
		),
		widget.NewHBox(
			widget.NewCheck("Drill", func(v bool) {}),
			layout.NewSpacer(),
		),
	)
	quit := widget.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Quit", func() { a.Quit() }),
	)

	w := a.NewWindow("Gerber viewer")
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, quit, nil, layers),
			c,
			layers,
			quit,
		))

	w.ShowAndRun()
}

func (v *viewController) pixelFunc(x, y, w, h int) color.Color {
	return color.RGBA{R: 0, G: 255, B: 0, A: 255}
}
