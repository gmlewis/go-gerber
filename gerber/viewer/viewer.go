// Package viewer views a Gerber design using Fyne.
package viewer

import (
	"image"
	"image/color"
	"log"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/fogleman/gg"
	"github.com/gmlewis/go-gerber/gerber"
)

type viewController struct {
	g         *gerber.Gerber
	mbb       gerber.MBB
	lastW     int
	lastH     int
	scale     float64
	drawLayer []bool
	app       fyne.App
	canvasObj fyne.CanvasObject
	img       *image.RGBA

	// These control the panning within the drawing area.
	xOffset int
	yOffset int

	indexDrill            int
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
	indexOutline          int
}

func initController(g *gerber.Gerber, app fyne.App) *viewController {
	vc := &viewController{
		g:                     g,
		app:                   app,
		mbb:                   g.MBB(),
		drawLayer:             make([]bool, len(g.Layers)),
		indexDrill:            -1,
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

	vc := initController(g, a)
	vc.img = image.NewRGBA(image.Rect(0, 0, 800, 800))
	vc.Resize(fyne.Size{Width: 800, Height: 800})
	c := canvas.NewRaster(vc.pixelFunc)
	c.SetMinSize(fyne.Size{Width: 800, Height: 800})
	vc.canvasObj = c

	layers := widget.NewVBox()
	addCheck := func(index int, label string) {
		if index >= 0 {
			check := widget.NewCheck(label, func(v bool) {
				vc.drawLayer[index] = v
				// widget.Refresh(vc)
				vc.Refresh()
				canvas.Refresh(c)
			})
			check.SetChecked(true)
			layers.Append(widget.NewHBox(check, layout.NewSpacer()))
		}
	}
	addCheck(vc.indexDrill, "Drill")
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
	addCheck(vc.indexOutline, "Outline")
	quit := widget.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Quit", func() { a.Quit() }),
	)

	w := a.NewWindow("Gerber viewer")
	w.Canvas().SetOnTypedRune(vc.OnTypedRune)
	w.Canvas().SetOnTypedKey(vc.OnTypedKey)
	w.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewBorderLayout(nil, quit, nil, layers),
			c,
			layers,
			quit,
		))

	w.ShowAndRun()
}

func (vc *viewController) OnTypedRune(key rune) {
	log.Printf("rune=%+q", key)
	switch key {
	case 'q': // TODO: Switch this to Alt-q when available.
		vc.app.Quit()
	case '-', '_':
		vc.zoom(0.25)
	case '+', '=':
		vc.zoom(-0.25)
	}
}

func (vc *viewController) OnTypedKey(event *fyne.KeyEvent) {
	if event == nil {
		return
	}
	log.Printf("event=%#v", *event)
	switch event.Name {
	case "Up":
		vc.pan(0, -vc.canvasObj.Size().Height/5)
	case "Down":
		vc.pan(0, vc.canvasObj.Size().Height/5)
	case "Left":
		vc.pan(vc.canvasObj.Size().Width/5, 0)
	case "Right":
		vc.pan(-vc.canvasObj.Size().Width/5, 0)
	}
}

func (vc *viewController) zoom(amount float64) {
	vc.scale = math.Exp2(amount) * vc.scale
	// widget.Refresh(vc)
	vc.Refresh()
	canvas.Refresh(vc.canvasObj)
}

func (vc *viewController) pan(dx, dy int) {
	vc.xOffset += dx
	vc.yOffset += dy
	// widget.Refresh(vc)
	vc.Refresh()
	canvas.Refresh(vc.canvasObj)
}

func (vc *viewController) ApplyTheme() {
	log.Printf("ApplyTheme")
}

func (vc *viewController) BackgroundColor() color.Color {
	return color.RGBA{R: 0, G: 0, B: 0, A: 255}
}

func (vc *viewController) CreateRenderer() fyne.WidgetRenderer {
	// log.Printf("CreateRenderer")
	return vc
}

func (vc *viewController) Hide() {
	// log.Printf("Hide")
	vc.canvasObj.Hide()
}

func (vc *viewController) Layout(size fyne.Size) {
	log.Printf("Layout")
}

func (vc *viewController) MinSize() fyne.Size {
	return vc.canvasObj.MinSize()
}

func (vc *viewController) Move(pos fyne.Position) {
	// log.Printf("Move")
	vc.canvasObj.Move(pos)
}

func (vc *viewController) Objects() []fyne.CanvasObject {
	// log.Printf("Objects")
	return []fyne.CanvasObject{vc.canvasObj}
}

func (vc *viewController) Position() fyne.Position {
	// log.Printf("Position")
	return vc.canvasObj.Position()
}

func (vc *viewController) Resize(size fyne.Size) {
	// log.Printf("Resize")
	if w, h := size.Width, size.Height; vc.lastW != w || vc.lastH != h {
		// vc.canvasObj.Resize(size)
		vc.lastW, vc.lastH = w, h
		vc.scale = (vc.mbb.Max[0] - vc.mbb.Min[0]) / float64(w)
		if s := (vc.mbb.Max[1] - vc.mbb.Min[1]) / float64(h); s > vc.scale {
			vc.scale = s
		}
		if w == h {
			vc.xOffset, vc.yOffset = 0, 0
		} else if w > h {
			vc.xOffset, vc.yOffset = (w-h)/2, 0
		} else {
			vc.xOffset, vc.yOffset = 0, (h-w)/2
		}
		log.Printf("(%v,%v): mbb=%v, scale=%v", w, h, vc.mbb, vc.scale)
		vc.img = image.NewRGBA(image.Rect(0, 0, w, h))
		vc.Refresh()
	}
}

func (vc *viewController) Show() {
	// log.Printf("Show")
	vc.canvasObj.Show()
}

func (vc *viewController) Size() fyne.Size {
	// log.Printf("Size")
	return vc.canvasObj.Size()
}

func (vc *viewController) Visible() bool {
	// log.Printf("Visible")
	return vc.canvasObj.Visible()
}

func (vc *viewController) Refresh() {
	const cs = 1.0 / float64(0xffff)
	ll := gerber.Pt{
		vc.scale*(float64(-vc.xOffset)) + vc.mbb.Min[0],
		vc.scale*(float64(-vc.yOffset)) + vc.mbb.Min[1],
	}
	ur := gerber.Pt{
		vc.scale*(float64(vc.lastW-1-vc.xOffset)) + vc.mbb.Min[0],
		vc.scale*(float64(vc.lastH-1-vc.yOffset)) + vc.mbb.Min[1],
	}
	bbox := &gerber.MBB{Min: ll, Max: ur}
	log.Printf("Refresh: MBB=%v", bbox)

	dc := gg.NewContextForImage(vc.img)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	renderLayer := func(index int, color color.Color) {
		if index < 0 || !vc.drawLayer[index] {
			return
		}
		r, g, b, a := color.RGBA()
		dc.SetRGBA(float64(r)*cs, float64(g)*cs, float64(b)*cs, float64(a)*cs)
		layer := vc.g.Layers[index]
		for _, p := range layer.Primitives {
			mbb := p.MBB()
			if !bbox.Contains(&mbb) {
				continue
			}
			// Render this primitive.
			switch v := p.(type) {
			case *gerber.CircleT:
				log.Printf("Render circle %v", mbb)
				x, y, r := 0.5*(mbb.Min[0]+mbb.Max[0]), 0.5*(mbb.Min[1]+mbb.Max[1]), 0.5*(mbb.Max[0]-mbb.Min[0])
				dc.DrawCircle(x/vc.scale-float64(vc.xOffset), y/vc.scale-float64(vc.yOffset), r/vc.scale)
				dc.Fill()
			default:
				log.Printf("%T not yet supported", v)
			}
		}
	}
	// Draw layers from bottom up
	renderLayer(vc.indexOutline, color.RGBA{R: 0, G: 255, B: 0, A: 255})
	renderLayer(vc.indexBottomSilkscreen, color.RGBA{R: 250, G: 50, B: 250, A: 255})
	renderLayer(vc.indexBottomSolderMask, color.RGBA{R: 250, G: 50, B: 50, A: 255})
	renderLayer(vc.indexBottom, color.RGBA{R: 50, G: 50, B: 250, A: 255})
	renderLayer(vc.indexLayer5, color.RGBA{R: 50, G: 150, B: 250, A: 255})
	renderLayer(vc.indexLayer4, color.RGBA{R: 150, G: 50, B: 250, A: 255})
	renderLayer(vc.indexLayer3, color.RGBA{R: 50, G: 50, B: 150, A: 255})
	renderLayer(vc.indexLayer2, color.RGBA{R: 50, G: 250, B: 250, A: 255})
	renderLayer(vc.indexTop, color.RGBA{R: 250, G: 50, B: 250, A: 255})
	renderLayer(vc.indexTopSolderMask, color.RGBA{R: 0, G: 150, B: 200, A: 255})
	renderLayer(vc.indexTopSilkscreen, color.RGBA{R: 250, G: 150, B: 0, A: 255})
	renderLayer(vc.indexDrill, color.RGBA{R: 200, G: 200, B: 200, A: 255})
	vc.img = dc.Image().(*image.RGBA)
}

func (vc *viewController) pixelFunc(x, y, w, h int) color.Color {
	vc.Resize(fyne.Size{Width: w, Height: h})
	return vc.img.At(x, y)
}
