// penta-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// five bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil) (one on
// each layer of a five-layer PCB) for manufacture on a printed circuit
// board (PCB).
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
	"github.com/gmlewis/go-gerber/gerber/viewer"
)

var (
	step       = flag.Float64("step", 0.02, "Resolution (in radians) of the spiral")
	n          = flag.Int("n", 100, "Number of full winds in each spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "penta-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	messageFmt = `Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`
	message2 = `2R ⇨ 4R
4R ⇨ TR
TR ⇨ BR
BR ⇨ 3R
3R ⇨ 3L`
	message3 = `3L ⇨ 2L
2L ⇨ 4L
4L ⇨ TL
TL ⇨ BL`
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	shiftAngle = -math.Pi / 3.0
	startLayer3R, layer3SpiralR, endLayer3R := s.genSpiral(1.0, shiftAngle, 0)
	startLayer3L, layer3SpiralL, endLayer3L := s.genSpiral(1.0, math.Pi+shiftAngle, -0.1)
	startLayer4R, layer4SpiralR, endLayer4R := s.genSpiral(-1.0, shiftAngle, 0)
	startLayer4L, layer4SpiralL, endLayer4L := s.genSpiral(-1.0, math.Pi+shiftAngle, 0)

	viaPadD := 0.5
	innerHole1Y := 0.5 * (*trace + viaPadD) / math.Sin(math.Pi/6)
	innerHole6X := innerHole1Y * math.Cos(math.Pi/6)

	hole1 := Point(0, innerHole1Y)
	hole3 := Point(0, -innerHole1Y)
	hole6 := Point(innerHole6X, 0.5*(*trace+viaPadD))
	hole7 := Point(-innerHole6X, -0.5*(*trace+viaPadD))
	hole10 := Point(innerHole6X, -0.5*(*trace+viaPadD))
	hole11 := Point(-innerHole6X, 0.5*(*trace+viaPadD))

	outerContactPt := func(pt Pt, angle float64) Pt {
		r := *trace*1.5 + 0.5*padD
		dx := r * math.Cos(angle)
		dy := r * math.Sin(angle)
		return Point(pt[0]+dx, pt[1]+dy)
	}
	outerContactPt2 := func(pt Pt, angle float64) Pt {
		r := *trace*2.5 + 0.5*padD
		dx := r * math.Cos(angle)
		dy := r * math.Sin(angle)
		return Point(pt[0]+dx, pt[1]+dy)
	}

	holeBL := outerContactPt(endBotL, math.Pi/3.0)
	holeTL4L := outerContactPt2(endTopL, 2.0*math.Pi/3.0)
	hole2L3L := outerContactPt(endLayer2L, math.Pi)
	holeBR3R := outerContactPt(endBotR, 4.0*math.Pi/3.0)
	holeTR4R := outerContactPt(endTopR, 5.0*math.Pi/3.0)

	hole2R := outerContactPt(endLayer2R, 0)

	viaDrill := func(pt Pt) *CircleT {
		const viaDrillD = 0.25
		return Circle(pt, viaDrillD)
	}
	contactDrill := func(pt Pt) *CircleT {
		const drillD = 1.0
		return Circle(pt, drillD)
	}

	drill := g.Drill()
	drill.Add(
		viaDrill(hole1),

		viaDrill(hole3),

		viaDrill(hole6),
		viaDrill(hole7),

		viaDrill(hole10),
		viaDrill(hole11),

		contactDrill(holeBL),
		contactDrill(holeTL4L),
		contactDrill(hole2L3L),
		contactDrill(holeBR3R),
		contactDrill(holeTR4R),
		contactDrill(hole2R),
	)

	viaPad := func(pt Pt) *CircleT {
		return Circle(pt, viaPadD)
	}
	contactPad := func(pt Pt) *CircleT {
		return Circle(pt, padD)
	}
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, *trace)
	}

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, topSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, topSpiralL, 0.0),

		viaPad(hole1),
		padLine(startTopL, hole1),

		viaPad(hole3),
		padLine(startTopR, hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		padLine(endTopL, holeTL4L),
		contactPad(hole2L3L),
		contactPad(holeBR3R),
		contactPad(holeTR4R),
		padLine(endTopR, holeTR4R),
		contactPad(hole2R),
	)

	layer2 := g.LayerN(2)
	layer2.Add(
		Polygon(Pt{0, 0}, true, layer2SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer2SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		padLine(startLayer2R, hole10),
		viaPad(hole11),
		padLine(startLayer2L, hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		contactPad(hole2L3L),
		padLine(endLayer2L, hole2L3L),
		contactPad(holeBR3R),
		contactPad(holeTR4R),
		contactPad(hole2R),
		padLine(endLayer2R, hole2R),
	)

	layer3 := g.LayerN(3)
	layer3.Add(
		Polygon(Pt{0, 0}, true, layer3SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer3SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		padLine(startLayer3L, hole6),
		viaPad(hole7),
		padLine(startLayer3R, hole7),
		padLine(hole6, hole7),

		viaPad(hole10),
		viaPad(hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		contactPad(hole2L3L),
		padLine(endLayer3L, hole2L3L),
		contactPad(holeBR3R),
		padLine(endLayer3R, holeBR3R),
		contactPad(holeTR4R),
		contactPad(hole2R),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		contactPad(hole2L3L),
		contactPad(holeBR3R),
		contactPad(holeTR4R),
		contactPad(hole2R),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, botSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, botSpiralL, 0.0),

		viaPad(hole1),
		padLine(startBotL, hole1),

		viaPad(hole3),
		padLine(startBotR, hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		contactPad(holeBL),
		padLine(endBotL, holeBL),
		contactPad(holeTL4L),
		contactPad(hole2L3L),
		contactPad(holeBR3R),
		padLine(endBotR, holeBR3R),
		contactPad(holeTR4R),
		contactPad(hole2R),
	)

	layer4 := g.LayerN(4)
	layer4.Add(
		Polygon(Pt{0, 0}, true, layer4SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer4SpiralL, 0.0),

		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		padLine(startLayer4R, hole10),
		viaPad(hole11),
		padLine(startLayer4L, hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		padLine(endLayer4L, holeTL4L),
		contactPad(hole2L3L),
		contactPad(holeBR3R),
		contactPad(holeTR4R),
		padLine(endLayer4R, holeTR4R),
		contactPad(hole2R),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		viaPad(hole1),

		viaPad(hole3),

		viaPad(hole6),
		viaPad(hole7),

		viaPad(hole10),
		viaPad(hole11),

		contactPad(holeBL),
		contactPad(holeTL4L),
		contactPad(hole2L3L),
		contactPad(holeBR3R),
		contactPad(holeTR4R),
		contactPad(hole2R),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		labelSize := pts * 12.0 / 18.0
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		outerLabel := func(pt Pt, label string) *TextT {
			r := math.Sqrt(pt[0]*pt[0] + pt[1]*pt[1])
			angle := 0.15 + math.Atan2(pt[1], pt[0])
			x := r * math.Cos(angle)
			y := r * math.Sin(angle)
			return Text(x, y, 1.0, label, *fontName, pts, &Center)
		}
		outerLabel2 := func(pt Pt, label string) *TextT {
			r := math.Sqrt(pt[0]*pt[0] + pt[1]*pt[1])
			angle := -0.15 + math.Atan2(pt[1], pt[0])
			x := r * math.Cos(angle)
			y := r * math.Sin(angle)
			return Text(x, y, 1.0, label, *fontName, pts, &Center)
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.5*r, 1.0, message, *fontName, pts, &Center),
			Text(hole1[0], hole1[1]+viaPadD, 1.0, "TL/BL", *fontName, labelSize, &BottomCenter),
			Text(hole3[0], hole3[1]-viaPadD, 1.0, "TR/BR", *fontName, labelSize, &TopCenter),
			Text(hole6[0]+0.5*viaPadD, hole6[1], 1.0, "3L", *fontName, labelSize, &BottomLeft),
			Text(hole7[0]-0.5*viaPadD, hole7[1], 1.0, "3R", *fontName, labelSize, &TopRight),
			Text(hole10[0]+0.5*viaPadD, hole10[1], 1.0, "2R/4R", *fontName, labelSize, &TopLeft),
			Text(hole11[0]-0.5*viaPadD, hole11[1], 1.0, "2L/4L", *fontName, labelSize, &BottomRight),
			Text(-0.5*r, -0.3*r, 1.0, message2, *fontName, pts, &Center),
			Text(0.5*r, -0.3*r, 1.0, message3, *fontName, pts, &Center),

			// Outer connections
			outerLabel(holeTR4R, "4R"),
			outerLabel2(holeTR4R, "TR"),

			outerLabel(holeTL4L, "TL"),
			outerLabel2(holeTL4L, "4L"),

			outerLabel(holeBR3R, "3R"),
			outerLabel2(holeBR3R, "BR"),

			outerLabel(holeBL, "BL"),

			outerLabel(hole2L3L, "3L"),
			outerLabel2(hole2L3L, "2L"),

			outerLabel2(hole2R, "2R"),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g, true)
	}
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
	size := 2 * math.Abs(p1[0])
	p2 := genPt(1.0, endAngle, *trace*0.5, math.Pi)
	xl := 2 * math.Abs(p2[0])
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
	if trimY < 0 { // Only for layer3SpiralL - extend another 2*Pi/3
		endAngle += 2.0 * math.Pi / 3.0
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
			if pts[len(pts)-trimYsteps][1] > trimY {
				break
			}
			trimYsteps++
		}
		lastStep := len(pts) - trimYsteps
		trimYsteps--
		pts = pts[0 : lastStep+1]
		pts = append(pts, Pt{pts[lastStep][0], trimY})
		angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		eX := genPt(xScale, angle, 0, offset)
		endPt = Pt{eX[0], trimY}
		nX := genPt(xScale, angle, -halfTW, offset)
		pts = append(pts, Pt{nX[0], trimY})
		// } else if trimY < 0 { // Only for layer2SpiralL
		// 	trimYsteps++
		// 	for {
		// 		if pts[len(pts)-trimYsteps][1] < trimY {
		// 			break
		// 		}
		// 		trimYsteps++
		// 	}
		// 	lastStep := len(pts) - trimYsteps
		// 	trimYsteps--
		// 	pts = pts[0 : lastStep+1]
		// 	pts = append(pts, Pt{pts[lastStep][0], trimY})
		// 	angle := s.startAngle + *step*float64(steps-1-trimYsteps)
		// 	eX := genPt(xScale, angle, 0, offset)
		// 	endPt = Pt{eX[0], trimY}
		// 	nX := genPt(xScale, angle, -halfTW, offset)
		// 	pts = append(pts, Pt{nX[0], trimY})
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
