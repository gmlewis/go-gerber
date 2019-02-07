// hex-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// six bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil) (one on
// each layer of a six-layer PCB) for manufacture on a printed circuit
// board (PCB).
package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
)

var (
	step     = flag.Float64("step", 0.02, "Resolution (in radians) of the spiral")
	n        = flag.Int("n", 100, "Number of full winds in each spiral")
	gap      = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace    = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix   = flag.String("prefix", "hex-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
)

const (
	messageFmt = `Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`
	message2 = `3L ⇨ 4L
4L ⇨ BL
BL ⇨ TL
TL ⇨ 5L
5L ⇨ 2L
2L ⇨ 3R`
	message3 = `3R ⇨ 4R
4R ⇨ BR
BR ⇨ TR
TR ⇨ 5R
5R ⇨ 2R`
)

func main() {
	flag.Parse()

	if *n < 12 {
		flag.Usage()
		log.Fatal("N must be >= 12.")
	}

	g := New(*prefix)

	s := newSpiral()

	startTopR, topSpiralR, endTopR := s.genSpiral(1.0, 0, 0)
	startTopL, topSpiralL, endTopL := s.genSpiral(1.0, math.Pi, 0)
	startBotR, botSpiralR, endBotR := s.genSpiral(-1.0, 0, 0)

	startBotL, botSpiralL, endBotL := s.genSpiral(-1.0, math.Pi, 0)

	shiftAngle := math.Pi / 3.0
	padD := 2.0
	startLayer2R, layer2SpiralR, endLayer2R := s.genSpiral(1.0, shiftAngle, 0)
	startLayer2L, layer2SpiralL, endLayer2L := s.genSpiral(1.0, math.Pi+shiftAngle, 0)
	startLayer3R, layer3SpiralR, endLayer3R := s.genSpiral(-1.0, shiftAngle, 0)
	startLayer3L, layer3SpiralL, endLayer3L := s.genSpiral(-1.0, math.Pi+shiftAngle, *trace+padD)

	shiftAngle = -math.Pi / 3.0
	startLayer4R, layer4SpiralR, endLayer4R := s.genSpiral(1.0, shiftAngle, 0)
	startLayer4L, layer4SpiralL, endLayer4L := s.genSpiral(1.0, math.Pi+shiftAngle, 0)
	startLayer5R, layer5SpiralR, endLayer5R := s.genSpiral(-1.0, shiftAngle, 0)
	startLayer5L, layer5SpiralL, endLayer5L := s.genSpiral(-1.0, math.Pi+shiftAngle, 0)

	viaPadD := 0.5
	innerHole1Y := 0.5 * (*trace + viaPadD) / math.Sin(math.Pi/6)
	innerHole6X := innerHole1Y * math.Cos(math.Pi/6)

	hole1 := Point(0, innerHole1Y)
	hole3 := Point(0, -innerHole1Y)
	hole6 := Point(innerHole6X, 0.5*(*trace+viaPadD))
	hole7 := Point(-innerHole6X, -0.5*(*trace+viaPadD))
	hole10 := Point(innerHole6X, -0.5*(*trace+viaPadD))
	hole11 := Point(-innerHole6X, 0.5*(*trace+viaPadD))

	outerViaPt := func(pt Pt, angle float64) Pt {
		r := *trace*1.5 + 0.5*viaPadD
		dx := r * math.Cos(angle)
		dy := r * math.Sin(angle)
		return Point(pt.X+dx, pt.Y+dy)
	}

	holeBL4L := outerViaPt(endBotL, math.Pi/3.0)
	holeTL5L := outerViaPt(endTopL, 2.0*math.Pi/3.0)
	hole2L3R := outerViaPt(endLayer2L, math.Pi)
	holeBR4R := outerViaPt(endBotR, 4.0*math.Pi/3.0)
	holeTR5R := outerViaPt(endTopR, 5.0*math.Pi/3.0)

	outerContactPt := func(pt Pt, angle float64) Pt {
		r := *trace*1.5 + 0.5*padD
		dx := r * math.Cos(angle)
		dy := r * math.Sin(angle)
		return Point(pt.X+dx, pt.Y+dy)
	}

	hole2R := outerContactPt(endLayer2R, 0)
	hole3L := outerContactPt(endLayer3L, 0)

	viaDrill := func(pt Pt) *CircleT {
		const viaDrillD = 0.25
		return Circle(pt.X, pt.Y, viaDrillD)
	}
	contactDrill := func(pt Pt) *CircleT {
		const drillD = 1.0
		return Circle(pt.X, pt.Y, drillD)
	}

	drill := g.Drill()
	drill.Add(
		viaDrill(hole1),

		viaDrill(hole3),

		viaDrill(hole6),
		viaDrill(hole7),

		viaDrill(hole10),
		viaDrill(hole11),

		viaDrill(holeBL4L),
		viaDrill(holeTL5L),
		viaDrill(hole2L3R),
		viaDrill(holeBR4R),
		viaDrill(holeTR5R),
		contactDrill(hole2R),
		contactDrill(hole3L),
	)

	viaPad := func(pt Pt) *CircleT {
		return Circle(pt.X, pt.Y, viaPadD)
	}
	contactPad := func(pt Pt) *CircleT {
		return Circle(pt.X, pt.Y, padD)
	}
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1.X, pt1.Y, pt2.X, pt2.Y, CircleShape, *trace)
	}

	top := g.TopCopper()
	top.Add(
		Polygon(0, 0, true, topSpiralR, 0.0),
		Polygon(0, 0, true, topSpiralL, 0.0),

		viaPad(hole1),
		padLine(startTopL, hole1),

		viaPad(hole3),
		padLine(startTopR, hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		padLine(endTopL, holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		padLine(endTopR, holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	layer2 := g.Layer2()
	layer2.Add(
		Polygon(0, 0, true, layer2SpiralR, 0.0),
		Polygon(0, 0, true, layer2SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		padLine(startLayer2R, hole10),
		viaPad(hole11),
		padLine(startLayer2L, hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		padLine(endLayer2L, hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		padLine(endLayer2R, hole2R),
		contactPad(hole3L),
	)

	layer4 := g.Layer4()
	layer4.Add(
		Polygon(0, 0, true, layer4SpiralR, 0.0),
		Polygon(0, 0, true, layer4SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		padLine(startLayer4L, hole6),
		viaPad(hole7),
		padLine(startLayer4R, hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		padLine(endLayer4L, holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		padLine(endLayer4R, holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(0, 0, true, botSpiralR, 0.0),
		Polygon(0, 0, true, botSpiralL, 0.0),

		viaPad(hole1),
		padLine(startBotL, hole1),

		viaPad(hole3),
		padLine(startBotR, hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		padLine(endBotL, holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		padLine(endBotR, holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	layer3 := g.Layer3()
	layer3.Add(
		Polygon(0, 0, true, layer3SpiralR, 0.0),
		Polygon(0, 0, true, layer3SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		padLine(startLayer3L, hole6),
		viaPad(hole7),
		padLine(startLayer3R, hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		padLine(endLayer3R, hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
		padLine(endLayer3L, hole3L),
	)

	layer5 := g.Layer5()
	layer5.Add(
		Polygon(0, 0, true, layer5SpiralR, 0.0),
		Polygon(0, 0, true, layer5SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		padLine(startLayer5R, hole10),
		viaPad(hole11),
		padLine(startLayer5L, hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		padLine(endLayer5L, holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		padLine(endLayer5R, holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		viaPad(holeBL4L),
		viaPad(holeTL5L),
		viaPad(hole2L3R),
		viaPad(holeBR4R),
		viaPad(holeTR5R),
		contactPad(hole2R),
		contactPad(hole3L),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(0, 0, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		labelSize := pts * 4.0 / 18.0
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*r, 1.0, message, *fontName, pts, Center),
			Text(hole1.X, hole1.Y+viaPadD, 1.0, "TL/BL", *fontName, labelSize, BottomCenter),
			Text(hole3.X, hole3.Y-viaPadD, 1.0, "TR/BR", *fontName, labelSize, TopCenter),
			Text(hole6.X+viaPadD, hole6.Y-0.5*viaPadD, 1.0, "3L/4L", *fontName, labelSize, BottomLeft),
			Text(hole7.X-viaPadD, hole7.Y+0.5*viaPadD, 1.0, "3R/4R", *fontName, labelSize, TopRight),
			Text(hole10.X+viaPadD, hole10.Y+0.5*viaPadD, 1.0, "2R/5R", *fontName, labelSize, TopLeft),
			Text(hole11.X-viaPadD, hole11.Y-0.5*viaPadD, 1.0, "2L/5L", *fontName, labelSize, BottomRight),
			Text(-0.5*r, -0.4*r, 1.0, message2, *fontName, pts, Center),
			Text(0.5*r, -0.4*r, 1.0, message3, *fontName, pts, Center),

			// Outer connections
			Text(holeTR5R.X, holeTR5R.Y+viaPadD, 1.0, "5R", *fontName, labelSize, BottomCenter),
			Text(holeTR5R.X, holeTR5R.Y-viaPadD, 1.0, "TR", *fontName, labelSize, TopCenter),

			Text(holeTL5L.X, holeTL5L.Y+viaPadD, 1.0, "TL", *fontName, labelSize, BottomCenter),
			Text(holeTL5L.X, holeTL5L.Y-viaPadD, 1.0, "5L", *fontName, labelSize, TopCenter),

			Text(holeBR4R.X, holeBR4R.Y+viaPadD, 1.0, "4R", *fontName, labelSize, BottomCenter),
			Text(holeBR4R.X, holeBR4R.Y-viaPadD, 1.0, "BR", *fontName, labelSize, TopCenter),

			Text(holeBL4L.X, holeBL4L.Y+viaPadD, 1.0, "BL", *fontName, labelSize, BottomCenter),
			Text(holeBL4L.X, holeBL4L.Y-viaPadD, 1.0, "4L", *fontName, labelSize, TopCenter),

			Text(hole2L3R.X, hole2L3R.Y+viaPadD, 1.0, "3R", *fontName, labelSize, BottomCenter),
			Text(hole2L3R.X, hole2L3R.Y-viaPadD, 1.0, "2L", *fontName, labelSize, TopCenter),

			Text(endLayer3L.X-0.5*padD, endLayer3L.Y, 1.0, "3L", *fontName, pts, BottomRight),
			Text(endLayer2R.X-0.5*padD, endLayer2R.Y, 1.0, "2R", *fontName, pts, TopRight),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

func genPt(xScale, angle, halfTW, offset float64) Pt {
	r := (angle + *trace + *gap) / (3 * math.Pi)
	x := (r + halfTW) * math.Cos(angle+offset)
	y := (r + halfTW) * math.Sin(angle+offset)
	return Point(x*xScale, y)
}

type spiral struct {
	startAngle float64
	endAngle   float64
	size       float64
}

func newSpiral() *spiral {
	startAngle := 3.5 * math.Pi
	endAngle := 2.0*math.Pi + float64(*n)*2.0*math.Pi
	p1 := genPt(1.0, endAngle, *trace*0.5, 0)
	size := 2 * math.Abs(p1.X)
	p2 := genPt(1.0, endAngle, *trace*0.5, math.Pi)
	xl := 2 * math.Abs(p2.X)
	if xl > size {
		size = xl
	}
	return &spiral{
		startAngle: startAngle,
		endAngle:   endAngle,
		size:       size,
	}
}

func (s *spiral) genSpiral(xScale, offset, trimY float64) (startPt Pt, pts []Pt, endPt Pt) {
	halfTW := *trace * 0.5
	endAngle := s.endAngle - math.Pi/3.0
	if trimY < 0 { // Only for layer2SpiralL - extend another Pi/2
		endAngle += 0.5 * math.Pi
	}
	steps := int(0.5 + (endAngle-s.startAngle) / *step)
	for i := 0; i < steps; i++ {
		angle := s.startAngle + *step*float64(i)
		if i == 0 {
			startPt = genPt(xScale, angle, 0, offset)
		}
		pts = append(pts, genPt(xScale, angle, halfTW, offset))
	}
	var trimYsteps int
	if trimY > 0 {
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps].Y > trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{X: pts[lastStep].X, Y: trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{X: eX.X, Y: trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{X: nX.X, Y: trimY})
	} else if trimY < 0 { // Only for layer2SpiralL
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps].Y < trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{X: pts[lastStep].X, Y: trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{X: eX.X, Y: trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{X: nX.X, Y: trimY})
	} else {
		pts = append(pts, genPt(xScale, endAngle, halfTW, offset))
		endPt = genPt(xScale, endAngle, 0, offset)
		pts = append(pts, genPt(xScale, endAngle, -halfTW, offset))
	}
	for i := steps - 1 - trimYsteps; i >= 0; i-- {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(xScale, angle, -halfTW, offset))
	}
	pts = append(pts, pts[0])
	return startPt, pts, endPt
}
