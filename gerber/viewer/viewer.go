// Package viewer views a Gerber design using Fyne.
package viewer

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"

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
	center    gerber.Pt
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

	// mu protects Refresh from being hit multiple times concurrently.
	mu sync.Mutex
}

func initController(g *gerber.Gerber, app fyne.App) *viewController {
	mbb := g.MBB()
	vc := &viewController{
		g:                     g,
		app:                   app,
		mbb:                   mbb,
		center:                gerber.Pt{0.5 * (mbb.Max[0] + mbb.Min[0]), 0.5 * (mbb.Max[1] + mbb.Min[1])},
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
	vc.scaleToFit(800, 800)
	vc.img = image.NewRGBA(image.Rect(0, 0, 800, 800))
	c := canvas.NewRaster(vc.pixelFunc)
	c.SetMinSize(fyne.Size{Width: 800, Height: 800})
	vc.canvasObj = c

	layers := widget.NewVBox()
	addCheck := func(index int, label string) {
		if index >= 0 {
			check := widget.NewCheck(label, func(v bool) {
				vc.drawLayer[index] = v
				vc.Refresh(nil)
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
	switch key {
	case 'q', 'Q': // TODO: Switch this to Alt-q when available.
		vc.app.Quit()
	case '-', '_':
		vc.zoom(-0.25)
	case '+', '=':
		vc.zoom(0.25)
	case 'f', 'F':
		vc.xOffset, vc.yOffset = 0, 0
		vc.scaleToFit(vc.lastW, vc.lastH)
		vc.Refresh(nil)
		canvas.Refresh(vc.canvasObj)
	default:
		log.Printf("Unhandled rune=%+q", key)
	}
}

func (vc *viewController) OnTypedKey(event *fyne.KeyEvent) {
	if event == nil {
		return
	}
	switch event.Name {
	case "Up":
		vc.pan(0, -vc.canvasObj.Size().Height/5)
	case "Down":
		vc.pan(0, vc.canvasObj.Size().Height/5)
	case "Left":
		vc.pan(vc.canvasObj.Size().Width/5, 0)
	case "Right":
		vc.pan(-vc.canvasObj.Size().Width/5, 0)
	default:
		log.Printf("Unhandled event=%#v", *event)
	}
}

func (vc *viewController) zoom(amount float64) {
	vc.scale = math.Exp2(amount) * vc.scale
	vc.Refresh(nil)
	canvas.Refresh(vc.canvasObj)
}

func (vc *viewController) pan(dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}
	vc.xOffset += dx
	vc.yOffset += dy
	b := vc.img.Bounds()
	p := image.Pt(-dx, dy)
	draw.Draw(vc.img, b, vc.img, b.Min.Add(p), draw.Src)
	dirtyRect := image.Rectangle{}
	switch {
	case dx > 0 && dy == 0: // pan left (image moves right)
		dirtyRect = b.Intersect(image.Rect(b.Min.X, b.Min.Y, b.Min.X+dx, b.Max.Y))
	case dx < 0 && dy == 0: // pan right (image moves left)
		dirtyRect = b.Intersect(image.Rect(b.Max.X+dx, b.Min.Y, b.Max.X, b.Max.Y))
	case dx == 0 && dy > 0: // pan down (image moves up)
		dirtyRect = b.Intersect(image.Rect(b.Min.X, b.Max.Y-dy, b.Max.X, b.Max.Y))
	case dx == 0 && dy < 0: // pan up (image moves down)
		dirtyRect = b.Intersect(image.Rect(b.Min.X, b.Min.Y, b.Max.X, b.Min.Y-dy))
	default:
		log.Printf("Unhandled pan (%v,%v)", dx, dy)
	}
	vc.Refresh(&dirtyRect)
	canvas.Refresh(vc.canvasObj)
}

func (vc *viewController) scaleToFit(w, h int) {
	vc.lastW, vc.lastH = w, h
	vc.scale = float64(w-1) / (vc.mbb.Max[0] - vc.mbb.Min[0])
	if s := float64(h-1) / (vc.mbb.Max[1] - vc.mbb.Min[1]); s < vc.scale {
		vc.scale = s
	}
	log.Printf("(%v,%v): mbb=%v, scale=%v", w, h, vc.mbb, vc.scale)
}

func (vc *viewController) Resize(w, h int) {
	if vc.lastW != w || vc.lastH != h {
		vc.lastW, vc.lastH = w, h
		vc.img = image.NewRGBA(image.Rect(0, 0, w, h))
		vc.Refresh(nil)
	}
}

func (vc *viewController) MBB() *gerber.MBB {
	// a vc.Offset of (0,0) means to center the design.
	xOffset, yOffset := float64(-vc.xOffset)/vc.scale, float64(-vc.yOffset)/vc.scale
	halfWidth, halfHeight := 0.5*float64(vc.lastW-1)/vc.scale, 0.5*float64(vc.lastH-1)/vc.scale
	ll := gerber.Pt{
		vc.center[0] + xOffset - halfWidth,
		vc.center[1] + yOffset - halfHeight,
	}
	ur := gerber.Pt{
		vc.center[0] + xOffset + halfWidth,
		vc.center[1] + yOffset + halfHeight,
	}
	return &gerber.MBB{Min: ll, Max: ur}
}

func (vc *viewController) xf(bbox *gerber.MBB) func(x float64) float64 {
	return func(x float64) float64 {
		return vc.scale * (x - bbox.Min[0])
	}
}

func (vc *viewController) yf(bbox *gerber.MBB) func(y float64) float64 {
	return func(y float64) float64 {
		return vc.scale * (bbox.Max[1] - y)
	}
}

func (vc *viewController) Refresh(dirtyRect *image.Rectangle) {
	const cs = 1.0 / float64(0xffff)
	bbox := vc.MBB()

	var dc *gg.Context
	if dirtyRect != nil {
		sx, sy := 1.0/float64(vc.lastW-1), 1.0/float64(vc.lastH-1)
		fx1, fx2 := sx*float64(dirtyRect.Min.X), sx*float64(dirtyRect.Max.X)
		fy1, fy2 := sy*float64(vc.lastH-1-dirtyRect.Min.Y), sy*float64(vc.lastH-1-dirtyRect.Max.Y)
		wx, wy := bbox.Max[0]-bbox.Min[0], bbox.Max[1]-bbox.Min[1]
		log.Printf("f1=(%v,%v), f2=(%v,%v), w=(%v,%v)", fx1, fy1, fx2, fy2, wx, wy)
		dc = gg.NewContext(dirtyRect.Max.X-dirtyRect.Min.X, dirtyRect.Max.Y-dirtyRect.Min.Y)
		ll := gerber.Pt{
			bbox.Min[0] + fx1*wx,
			bbox.Min[1] + fy1*wy,
		}
		ur := gerber.Pt{
			bbox.Min[0] + fx2*wx,
			bbox.Min[1] + fy2*wy,
		}
		log.Printf("dirtyRect=%v, old bbox=%v, new bbox=%v", *dirtyRect, *bbox, gerber.MBB{Min: ll, Max: ur})
		bbox = &gerber.MBB{Min: ll, Max: ur}
	} else {
		dc = gg.NewContextForImage(vc.img)
	}

	// log.Printf("Refresh: MBB=%v", bbox)
	xf := vc.xf(bbox)
	yf := vc.yf(bbox)

	dc.SetRGB(0, 0, 0)
	dc.Clear()
	renderLayer := func(index int, color color.Color) {
		if index < 0 || !vc.drawLayer[index] {
			return
		}
		r, g, b, a := color.RGBA()
		fr, fg, fb, fa := float64(r)*cs, float64(g)*cs, float64(b)*cs, float64(a)*cs
		foreground := func(ctx *gg.Context) {
			ctx.SetRGBA(fr, fg, fb, fa)
		}
		foreground(dc)
		layer := vc.g.Layers[index]
		for _, p := range layer.Primitives {
			mbb := p.MBB()
			if !bbox.Intersects(&mbb) {
				continue
			}
			// Render this primitive.
			switch v := p.(type) {
			case *gerber.ArcT:
				// TODO: account for line shape.
				dc.SetLineWidth(v.Thickness * vc.scale)
				delta := v.EndAngle - v.StartAngle
				length := delta * v.Radius
				// Resolution of segments is 0.1mm
				segments := int(0.5+length*10.0) + 1
				delta /= float64(segments)

				angle := float64(v.StartAngle)
				for i := 0; i < segments; i++ {
					x1 := v.Center[0] + v.XScale*math.Cos(angle)*v.Radius
					y1 := v.Center[1] + v.YScale*math.Sin(angle)*v.Radius

					angle += delta

					x2 := v.Center[0] + v.XScale*math.Cos(angle)*v.Radius
					y2 := v.Center[1] + v.YScale*math.Sin(angle)*v.Radius

					dc.DrawLine(xf(x1), yf(y1), xf(x2), yf(y2))
				}
				dc.Stroke()
			case *gerber.CircleT:
				x, y, r := 0.5*(mbb.Min[0]+mbb.Max[0]), 0.5*(mbb.Min[1]+mbb.Max[1]), 0.5*(mbb.Max[0]-mbb.Min[0])
				dc.DrawCircle(xf(x), yf(y), r*vc.scale)
				dc.Fill()
			case *gerber.LineT:
				// TODO: account for line shape.
				dc.SetLineWidth(v.Thickness * vc.scale)
				dc.DrawLine(xf(v.P1[0]), yf(v.P1[1]), xf(v.P2[0]), yf(v.P2[1]))
				dc.Stroke()
			case *gerber.TextT:
				// Render text into new context, the copy foreground pixels only.
				bnds := vc.img.Bounds()
				nc := gg.NewContext(bnds.Max.X, bnds.Max.Y)
				// nc := gg.NewContextForImage(vc.img)
				for _, poly := range v.Render.Polygons {
					if poly.Dark {
						foreground(nc)
					} else {
						nc.SetRGB(0, 0, 0)
					}
					for i, pt := range poly.Pts {
						if i == 0 {
							nc.MoveTo(xf(pt[0]), yf(pt[1]))
						} else {
							nc.LineTo(xf(pt[0]), yf(pt[1]))
						}
					}
					nc.Fill()
				}
				llx, lly := int(xf(mbb.Min[0])), int(yf(mbb.Max[1]))
				urx, ury := int(0.5+xf(mbb.Max[0])), int(0.5+yf(mbb.Min[1]))
				// log.Printf("ll=(%v,%v), ur=(%v,%v)", llx, lly, urx, ury)
				img := nc.Image()
				foreground(dc)
				for y := lly; y <= ury; y++ {
					for x := llx; x <= urx; x++ {
						c := img.At(x, y)
						cr, cg, cb, _ := c.RGBA()
						if cr == 0 && cg == 0 && cb == 0 {
							continue
						}
						dc.SetPixel(x, y)
					}
				}
			case *gerber.PolygonT:
				for i, pt := range v.Points {
					p := gerber.Pt{pt[0] + v.Offset[0], pt[1] + v.Offset[1]}
					if i == 0 {
						dc.MoveTo(xf(p[0]), yf(p[1]))
					} else {
						dc.LineTo(xf(p[0]), yf(p[1]))
					}
				}
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
	if dirtyRect != nil {
		draw.Draw(vc.img, *dirtyRect, dc.Image(), image.Point{0, 0}, draw.Src)
	} else {
		vc.img = dc.Image().(*image.RGBA)
	}
}

func (vc *viewController) pixelFunc(x, y, w, h int) color.Color {
	if vc.lastW != w || vc.lastH != h {
		vc.mu.Lock()
		vc.Resize(w, h)
		vc.mu.Unlock()
	}
	return vc.img.At(x, y)
}
