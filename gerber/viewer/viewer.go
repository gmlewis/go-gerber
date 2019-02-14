// Package viewer views a Gerber design using Fyne.
package viewer

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/gmlewis/go-gerber/gerber"
)

type viewController struct {
	g         *gerber.Gerber
	lastW     int
	lastH     int
	scale     float64
	drawLayer []bool

	indexTopSilkscreen    int
	indexTopSolderMask    int
	indexTop              int
	indexLayer2           int
	indexLayer3           int
	indexLayer4           int
	indexLayer5           int
	indexBottom           int
	indexBottomSilkscreen int
	indexBottomSolderMask int
	indexDrill            int
	indexOutline          int
}

func initController(g *gerber.Gerber) *viewController {
	vc := &viewController{
		g:                     g,
		drawLayer:             make([]bool, len(g.Layers)),
		indexTopSilkscreen:    -1,
		indexTopSolderMask:    -1,
		indexTop:              -1,
		indexLayer2:           -1,
		indexLayer3:           -1,
		indexLayer4:           -1,
		indexLayer5:           -1,
		indexBottom:           -1,
		indexBottomSilkscreen: -1,
		indexBottomSolderMask: -1,
		indexDrill:            -1,
		indexOutline:          -1,
	}

	for i, layer := range g.Layers {
		vc.drawLayer[i] = true
		switch layer.Filename[len(layer.Filename)-3:] {
		case "gtl":
			vc.indexTop = i
		case "gts":
			vc.indexTopSolderMask = i
		case "gto":
			vc.indexTopSilkscreen = i
		case "gbl":
			vc.indexBottom = i
		case "gbs":
			vc.indexBottomSolderMask = i
		case "gbo":
			vc.indexBottomSilkscreen = i
		case "g2l":
			vc.indexLayer2 = i
		case "g3l":
			vc.indexLayer3 = i
		case "g4l":
			vc.indexLayer4 = i
		case "g5l":
			vc.indexLayer5 = i
		case "xln":
			vc.indexDrill = i
		case "gko":
			vc.indexOutline = i
		default:
			log.Fatalf("Unknown Gerber layer: %v", layer.Filename)
		}
	}
	return vc
}

func Gerber(g *gerber.Gerber) {
	a := app.New()

	vc := initController(g)
	c := canvas.NewRaster(vc.pixelFunc)
	c.SetMinSize(fyne.Size{Width: 800, Height: 800})

	layers := widget.NewVBox()
	addCheck := func(index int, label string) {
		if index >= 0 {
			check := widget.NewCheck(label, func(v bool) {
				vc.drawLayer[index] = v
				canvas.Refresh(c)
			})
			check.SetChecked(true)
			layers.Append(widget.NewHBox(check, layout.NewSpacer()))
		}
	}
	addCheck(vc.indexTopSilkscreen, "Top Silkscreen")
	addCheck(vc.indexTopSolderMask, "Top Solder Mask")
	addCheck(vc.indexTop, "Top")
	addCheck(vc.indexLayer2, "Layer 2")
	addCheck(vc.indexLayer3, "Layer 3")
	addCheck(vc.indexLayer4, "Layer 4")
	addCheck(vc.indexLayer5, "Layer 5")
	addCheck(vc.indexBottom, "Bottom")
	addCheck(vc.indexBottomSolderMask, "Bottom Solder Mask")
	addCheck(vc.indexBottomSilkscreen, "Bottom Silkscreen")
	addCheck(vc.indexDrill, "Drill")
	addCheck(vc.indexOutline, "Outline")
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

func (vc *viewController) pixelFunc(x, y, w, h int) color.Color {
	mbb := vc.g.MBB()
	if vc.lastW != w || vc.lastH != h {
		vc.lastW, vc.lastH = w, h
		vc.scale = (mbb.Max[0] - mbb.Min[0]) / float64(w)
		if s := (mbb.Max[1] - mbb.Min[1]) / float64(h); s < vc.scale {
			vc.scale = s
		}
		log.Printf("(%v,%v): mbb=%v, scale=%v", w, h, mbb, vc.scale)
	}

	ll := gerber.Pt{vc.scale*(float64(x)-0.5) + mbb.Min[0], vc.scale*(float64(h-y-1)-0.5) + mbb.Min[1]}
	ur := gerber.Pt{vc.scale*(float64(x)+0.5) + mbb.Min[0], vc.scale*(float64(h-y-1)+0.5) + mbb.Min[1]}
	bbox := &gerber.MBB{Min: ll, Max: ur}

	// Draw layers from bottom up
	var colors []color.Color

	renderLayer := func(index int, color color.Color) {
		if index < 0 || !vc.drawLayer[index] {
			return
		}
		if vc.g.Layers[index].IsDark(bbox) {
			colors = append(colors, color)
		}
	}
	renderLayer(vc.indexOutline, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	renderLayer(vc.indexDrill, color.RGBA{R: 50, G: 50, B: 50, A: 255})
	renderLayer(vc.indexBottomSilkscreen, color.RGBA{R: 250, G: 50, B: 250, A: 255})
	renderLayer(vc.indexBottomSolderMask, color.RGBA{R: 250, G: 50, B: 50, A: 255})
	renderLayer(vc.indexBottom, color.RGBA{R: 50, G: 50, B: 250, A: 255})
	renderLayer(vc.indexLayer5, color.RGBA{R: 50, G: 150, B: 250, A: 255})
	renderLayer(vc.indexLayer4, color.RGBA{R: 150, G: 50, B: 250, A: 255})
	renderLayer(vc.indexLayer3, color.RGBA{R: 50, G: 50, B: 150, A: 255})
	renderLayer(vc.indexLayer2, color.RGBA{R: 50, G: 250, B: 250, A: 255})
	renderLayer(vc.indexTop, color.RGBA{R: 250, G: 50, B: 250, A: 255})
	renderLayer(vc.indexTopSolderMask, color.RGBA{R: 250, G: 250, B: 50, A: 255})
	renderLayer(vc.indexTopSilkscreen, color.RGBA{R: 250, G: 150, B: 0, A: 255})

	if len(colors) == 0 {
		return color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	// blend colors?
	return colors[len(colors)-1]
}
