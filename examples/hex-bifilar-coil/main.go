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
	step     = flag.Float64("step", 0.04, "Resolution (in radians) of the spiral")
	n        = flag.Int("n", 100, "Number of full winds in each spiral")
	gap      = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace    = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix   = flag.String("prefix", "hex-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	pts      = flag.Float64("pts", 18.0, "Font point size (72 pts = 1 inch = 25.4 mm)")
)

const (
	message = `With a trace and gap size of 0.15mm, this
hex bifilar coil should have a DC resistance
of approx. 1393.2Ω. Each spiral has 100 coils.`
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

	g := New(*prefix)

	s := newSpiral()

	startTopR, topSpiralR, endTopR := s.genSpiral(1.0, 0, 0)
	startTopL, topSpiralL, endTopL := s.genSpiral(1.0, math.Pi, 0)
	startBotR, botSpiralR, endBotR := s.genSpiral(-1.0, 0, 0)
	// startTopR := genPt(1.0, s.startAngle, 0, 0)
	// startTopL := genPt(1.0, s.startAngle, 0, math.Pi)
	log.Printf("startTopR=%v, startTopL=%v", startTopR, startTopL)
	log.Printf("endTopR=%v, endTopL=%v", endTopR, endTopL)

	startBotL, botSpiralL, endBotL := s.genSpiral(-1.0, math.Pi, 0 /**trace+padD */)
	log.Printf("startBotR=%v, startBotL=%v", startBotR, startBotL)
	log.Printf("endBotL=%v", endBotL)

	shiftAngle := math.Pi / 3.0
	startLayer2R, layer2SpiralR, endLayer2R := s.genSpiral(1.0, shiftAngle, 0)
	startLayer2L, layer2SpiralL, endLayer2L := s.genSpiral(1.0, math.Pi+shiftAngle, 0 /*-(*trace + padD)*/)
	startLayer3R, layer3SpiralR, endLayer3R := s.genSpiral(-1.0, shiftAngle, 0)
	startLayer3L, layer3SpiralL, endLayer3L := s.genSpiral(-1.0, math.Pi+shiftAngle, 0 /* *trace+padD*/)
	log.Printf("startLayer2R=%v, startLayer2L=%v", startLayer2R, startLayer2L)
	log.Printf("startLayer3R=%v, startLayer3L=%v", startLayer3R, startLayer3L)
	log.Printf("endLayer2L=%v", endLayer2L)
	log.Printf("endLayer3L=%v", endLayer3L)
	log.Printf("endLayer2R=%v", endLayer2R)
	log.Printf("endLayer3R=%v", endLayer3R)

	shiftAngle = -math.Pi / 3.0
	startLayer4R, layer4SpiralR, endLayer4R := s.genSpiral(1.0, shiftAngle, 0)
	startLayer4L, layer4SpiralL, endLayer4L := s.genSpiral(1.0, math.Pi+shiftAngle, 0 /*-(*trace + padD)*/)
	startLayer5R, layer5SpiralR, endLayer5R := s.genSpiral(-1.0, shiftAngle, 0)
	startLayer5L, layer5SpiralL, endLayer5L := s.genSpiral(-1.0, math.Pi+shiftAngle, 0 /**trace+padD*/)
	// startLayer4R := genPt(1.0, s.startAngle, 0, shiftAngle)
	// startLayer4L := genPt(1.0, s.startAngle, 0, math.Pi+shiftAngle)
	log.Printf("startLayer4R=%v, startLayer4L=%v", startLayer4R, startLayer4L)
	log.Printf("startLayer5R=%v, startLayer5L=%v", startLayer5R, startLayer5L)
	log.Printf("endLayer4R=%v, endLayer4L=%v", endLayer4R, endLayer4L)
	log.Printf("endLayer5L=%v", endLayer5L)
	log.Printf("%v, %v, %v, %v, %v, %v", len(topSpiralL), len(botSpiralL), len(layer2SpiralL), len(layer3SpiralL), len(layer4SpiralL), len(layer5SpiralL))

	viaPadD := 0.5
	padD := 2.0
	// viaOffset := math.Sqrt(0.5 * (*trace + viaPadD) * (*trace + viaPadD))
	hole2Offset := 0.5 * (*trace + viaPadD)
	// hole4PadOffset := 0.5 * (viaPadD + *trace)
	hole5PadOffset := 0.5 * (padD + *trace)
	innerHole1Y := 0.5 * (*trace + viaPadD) / math.Sin(math.Pi/6)
	innerHole6X := innerHole1Y * math.Cos(math.Pi/6)

	// Lower connecting trace between two spirals
	// hole1 := Point(viaOffset, 0)
	hole1 := Point(0, innerHole1Y)
	hole2 := Point(endTopL.X-hole2Offset, endTopL.Y)
	// Upper connecting trace for left spiral
	hole3 := Point(0, -innerHole1Y)
	// FIX hole4 := Point(endTopR.X+hole4PadOffset, *trace+padD)
	// Lower connecting trace for right spiral
	hole5 := Point(endTopR.X+hole5PadOffset, endTopR.Y)
	// Layer 2 and 3 inner connecting holes
	hole6 := Point(innerHole6X, 0.5*(*trace+viaPadD))
	hole7 := Point(-innerHole6X, -0.5*(*trace+viaPadD))
	// Layer 2 and 3 outer connecting hole
	// FIX hole8 := Point(endLayer2R.X, endLayer2R.Y+hole2Offset)
	// FIX THIS hole9 := Point(endLayer2L.X+hole5PadOffset, -(*trace + padD))
	// Layer 4 and 5 inner connecting holes
	hole10 := Point(innerHole6X, -0.5*(*trace+viaPadD))
	hole11 := Point(-innerHole6X, 0.5*(*trace+viaPadD))

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
		// Lower connecting trace between two spirals
		viaDrill(hole1),
		viaDrill(hole2),
		// Upper connecting trace for left spiral
		viaDrill(hole3),
		// 		viaDrill(hole4),
		// Lower connecting trace for right spiral
		contactDrill(hole5),
		// Layer 2 and 3 inner connecting holes
		viaDrill(hole6),
		viaDrill(hole7),
		// Layer 2 and 3 outer connecting hole
		//		viaDrill(hole8),
		//FIX THIS Circle(hole9.X, hole9.Y, drillD),
		// Layer 4 and 5 inner connecting holes
		viaDrill(hole10),
		viaDrill(hole11),
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
		// Lower connecting trace between two spirals
		viaPad(hole1),
		padLine(startTopL, hole1),
		viaPad(hole2),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		padLine(startTopR, hole3),
		//		viaPad(hole4),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		// Layer 2 and 3 outer connecting hole
		//		viaPad(hole8),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	layer2 := g.Layer2()
	layer2.Add(
		Polygon(0, 0, true, layer2SpiralR, 0.0),
		Polygon(0, 0, true, layer2SpiralL, 0.0),
		// Lower connecting trace between two spirals
		viaPad(hole1),
		viaPad(hole2),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		//		Line(startLayer2R.X, startLayer2R.Y, hole3.X, hole3.Y, RectShape, *trace),
		//		Line(startLayer2L.X, startLayer2L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//		Line(endLayer2R.X, endLayer2R.Y, hole8.X, hole8.Y, RectShape, *trace),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		//		Line(endLayer2L.X, endLayer2L.Y, hole9.X, hole9.Y, RectShape, *trace),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		padLine(startLayer2R, hole10),
		viaPad(hole11),
		padLine(startLayer2L, hole11),
	)

	layer4 := g.Layer4()
	layer4.Add(
		Polygon(0, 0, true, layer4SpiralR, 0.0),
		Polygon(0, 0, true, layer4SpiralL, 0.0),
		// Lower connecting trace between two spirals
		viaPad(hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		padLine(startLayer4L, hole6),
		viaPad(hole7),
		padLine(startLayer4R, hole7),
		//		Line(startLayer4R.X, startLayer4R.Y, hole3.X, hole3.Y, RectShape, *trace),
		//		Line(startLayer4L.X, startLayer4L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//		Line(endLayer4R.X, endLayer4R.Y, hole8.X, hole8.Y, RectShape, *trace),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		//		Line(endLayer4L.X, endLayer4L.Y, hole9.X, hole9.Y, RectShape, *trace),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		// Lower connecting trace between two spirals
		viaPad(hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(0, 0, true, botSpiralR, 0.0),
		Polygon(0, 0, true, botSpiralL, 0.0),
		//		Line(endTopL.X, endTopL.Y, hole2.X, hole2.Y, RectShape, *trace),
		// Lower connecting trace between two spirals
		viaPad(hole1),
		padLine(startBotL, hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		//		Line(endTopL.X, endTopL.Y, hole2.X, hole2.Y, RectShape, *trace),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		padLine(startBotR, hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		//		Line(endBotL.X, endBotL.Y, hole4.X, hole4.Y, RectShape, *trace),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		//		Line(startTopR.X, startTopR.Y, hole6.X, hole6.Y, RectShape, *trace),
		//		Line(startL.X, startL.Y, hole7.X, hole7.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	layer3 := g.Layer3()
	layer3.Add(
		Polygon(0, 0, true, layer3SpiralR, 0.0),
		Polygon(0, 0, true, layer3SpiralL, 0.0),
		// Lower connecting trace between two spirals
		viaPad(hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		//		Line(endLayer3L.X, endLayer3L.Y, hole4.X, hole4.Y, RectShape, *trace),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		padLine(startLayer3L, hole6),
		viaPad(hole7),
		padLine(startLayer3R, hole7),
		//		Line(startLayer2R.X, startLayer2R.Y, hole3.X, hole3.Y, RectShape, *trace),
		//		Line(startLayer2L.X, startLayer2L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//		Line(endLayer2R.X, endLayer2R.Y, hole8.X, hole8.Y, RectShape, *trace),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	layer5 := g.Layer5()
	layer5.Add(
		Polygon(0, 0, true, layer5SpiralR, 0.0),
		Polygon(0, 0, true, layer5SpiralL, 0.0),
		// Lower connecting trace between two spirals
		viaPad(hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		//		Circle(hole4.X, hole4.Y, viaPadD),
		//		Line(endLayer5L.X, endLayer5L.Y, hole4.X, hole4.Y, RectShape, *trace),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		//		Line(startLayer2R.X, startLayer2R.Y, hole3.X, hole3.Y, RectShape, *trace),
		//		Line(startLayer2L.X, startLayer2L.Y, hole1.X, hole1.Y, RectShape, *trace),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//		Line(endLayer2R.X, endLayer2R.Y, hole8.X, hole8.Y, RectShape, *trace),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		padLine(startLayer5R, hole10),
		viaPad(hole11),
		padLine(startLayer5L, hole11),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		// Lower connecting trace between two spirals
		viaPad(hole1),
		//		Circle(hole2.X, hole2.Y, viaPadD),
		// Upper connecting trace for left spiral
		viaPad(hole3),
		// Circle(hole4.X, hole4.Y, viaPadD),
		// Lower connecting trace for right spiral
		contactPad(hole5),
		// Layer 2 and 3 inner connecting holes
		viaPad(hole6),
		viaPad(hole7),
		// Layer 2 and 3 outer connecting hole
		//		Circle(hole8.X, hole8.Y, viaPadD),
		//FIX THIS Circle(hole9.X, hole9.Y, padD),
		// Layer 4 and 5 inner connecting holes
		viaPad(hole10),
		viaPad(hole11),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(0, 0, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		radius := endLayer3L.X
		labelSize := 4.0

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*radius, 1.0, message, *fontName, *pts, Center),
			Text(hole1.X, hole1.Y+viaPadD, 1.0, "TL/BL", *fontName, labelSize, BottomCenter),
			//			Text(hole2.X+viaPadD, hole2.Y, 1.0, "hole2", *fontName, labelSize, CenterLeft),
			Text(hole3.X, hole3.Y-viaPadD, 1.0, "TR/BR", *fontName, labelSize, TopCenter),
			//			Text(hole4.X-padD, hole4.Y, 1.0, "hole4", *fontName, labelSize, CenterRight),
			// Text(hole5.X-padD, hole5.Y, 1.0, "hole5", *fontName, labelSize, CenterRight),
			Text(hole6.X+viaPadD, hole6.Y-0.5*viaPadD, 1.0, "3L/4L", *fontName, labelSize, BottomLeft),
			Text(hole7.X-viaPadD, hole7.Y+0.5*viaPadD, 1.0, "3R/4R", *fontName, labelSize, TopRight),
			//			Text(hole8.X, hole8.Y-viaPadD, 1.0, "hole8", *fontName, labelSize, TopCenter),
			//FIX THIS Text(hole9.X-padD, hole9.Y, 1.0, "hole9", *fontName, labelSize, CenterRight),
			Text(hole10.X+viaPadD, hole10.Y+0.5*viaPadD, 1.0, "2R/5R", *fontName, labelSize, TopLeft),
			Text(hole11.X-viaPadD, hole11.Y-0.5*viaPadD, 1.0, "2L/5L", *fontName, labelSize, BottomRight),
			Text(-0.5*radius, -0.4*radius, 1.0, message2, *fontName, *pts, Center),
			Text(0.5*radius, -0.4*radius, 1.0, message3, *fontName, *pts, Center),

			// Debugging outer connections
			Text(endTopR.X, endTopR.Y, 1.0, "TR", *fontName, labelSize, TopCenter),
			Text(endTopL.X, endTopL.Y, 1.0, "TL", *fontName, labelSize, BottomCenter),
			Text(endBotR.X, endBotR.Y, 1.0, "BR", *fontName, labelSize, TopCenter),
			Text(endBotL.X, endBotL.Y, 1.0, "BL", *fontName, labelSize, BottomCenter),
			Text(endLayer2R.X, endLayer2R.Y, 1.0, "2R", *fontName, labelSize, TopCenter),
			Text(endLayer2L.X, endLayer2L.Y, 1.0, "2L", *fontName, labelSize, TopCenter),
			Text(endLayer3R.X, endLayer3R.Y, 1.0, "3R", *fontName, labelSize, BottomCenter),
			Text(endLayer3L.X, endLayer3L.Y, 1.0, "3L", *fontName, labelSize, BottomCenter),
			Text(endLayer4R.X, endLayer4R.Y, 1.0, "4R", *fontName, labelSize, BottomCenter),
			Text(endLayer4L.X, endLayer4L.Y, 1.0, "4L", *fontName, labelSize, TopCenter),
			Text(endLayer5R.X, endLayer5R.Y, 1.0, "5R", *fontName, labelSize, BottomCenter),
			Text(endLayer5L.X, endLayer5L.Y, 1.0, "5L", *fontName, labelSize, TopCenter),
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
