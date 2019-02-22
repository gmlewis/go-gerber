// tri-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// three bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil) (one on
// each layer of a three-layer PCB) for manufacture on a printed circuit
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
	step       = flag.Float64("step", 0.0171, "Resolution (in radians) of the spiral")
	n          = flag.Int("n", 100, "Number of full winds in each spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "tri-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	pts        = flag.Float64("pts", 18.0, "Font point size (72 pts = 1 inch = 25.4 mm)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	messageFmt = `Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`
	message2 = `Top layer: TR ⇨ BR
Bottom layer: BR ⇨ TL
Top layer: TL ⇨ BL
Bottom layer: BL ⇨ 2L
Layer 2: 2L ⇨ 2R`
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

	g := New(*prefix)

	s := newSpiral()

	topSpiralR, endR := s.genSpiral(1.0, 0, 0)
	topSpiralL, endL := s.genSpiral(1.0, math.Pi, 0)
	botSpiralR, _ := s.genSpiral(-1.0, 0, 0)
	startR := genPt(1.0, s.startAngle, 0, 0)
	startL := genPt(1.0, s.startAngle, 0, math.Pi)

	viaPadD := 0.5
	viaDrillD := 0.25

	padD := 2.0
	drillD := 1.0
	botSpiralL, botEndL := s.genSpiral(-1.0, math.Pi, *trace+padD)

	shiftAngle := 0.5 * math.Pi
	layer2SpiralR, layer2EndR := s.genSpiral(-1.0, shiftAngle, 0)
	layer2SpiralL, layer2EndL := s.genSpiral(-1.0, math.Pi+shiftAngle, *trace+padD)
	startLayer2R := genPt(1.0, s.startAngle, 0, shiftAngle)
	startLayer2L := genPt(1.0, s.startAngle, 0, math.Pi+shiftAngle)

	viaOffset := math.Sqrt(0.5 * (*trace + viaPadD) * (*trace + viaPadD))
	hole2Offset := 0.5 * (*trace + padD)
	hole4PadOffset := 0.5 * (padD + *trace)
	hole5PadOffset := 0.5 * (padD + *trace)

	// Lower connecting trace between two spirals
	// hole1 := Point(startR[0], startR[1]+viaPadOffset)
	hole1 := Point(viaOffset, 0)
	hole2 := Point(endL[0]-hole2Offset, endL[1])
	// Upper connecting trace for left spiral
	// hole3 := Point(startL[0], startL[1]-viaPadOffset)
	hole3 := Point(-viaOffset, 0)
	hole4 := Point(endR[0]+hole4PadOffset, *trace+padD)
	// Lower connecting trace for right spiral
	hole5 := Point(endR[0]+hole5PadOffset, endR[1])
	// Layer 2 inner connecting holes
	hole6 := Point(0, viaOffset)
	hole7 := Point(0, -viaOffset)
	// Layer 2 outer connecting hole
	hole8 := Point(layer2EndR[0], layer2EndR[1]+hole2Offset)

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, topSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, topSpiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, padD),
		Line(endL[0], endL[1], hole2[0], hole2[1], RectShape, *trace),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
		Line(endR[0], endR[1], hole5[0], hole5[1], RectShape, *trace),
		// Layer 2 inner connecting holes
		Circle(hole6, viaPadD),
		Circle(hole7, viaPadD),
		Line(startR[0], startR[1], hole6[0], hole6[1], RectShape, *trace),
		Line(startL[0], startL[1], hole7[0], hole7[1], RectShape, *trace),
		// Layer 2 outer connecting hole
		Circle(hole8, padD),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, padD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
		// Layer 2 inner connecting holes
		Circle(hole6, viaPadD),
		Circle(hole7, viaPadD),
		// Layer 2 outer connecting hole
		Circle(hole8, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, botSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, botSpiralL, 0.0),
		Line(endL[0], endL[1], hole2[0], hole2[1], RectShape, *trace),
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, padD),
		Line(endL[0], endL[1], hole2[0], hole2[1], RectShape, *trace),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		Line(botEndL[0], botEndL[1], hole4[0], hole4[1], RectShape, *trace),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
		// Layer 2 inner connecting holes
		Circle(hole6, viaPadD),
		Circle(hole7, viaPadD),
		Line(startR[0], startR[1], hole6[0], hole6[1], RectShape, *trace),
		Line(startL[0], startL[1], hole7[0], hole7[1], RectShape, *trace),
		// Layer 2 outer connecting hole
		Circle(hole8, padD),
	)

	layer2 := g.LayerN(2)
	layer2.Add(
		Polygon(Pt{0, 0}, true, layer2SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer2SpiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, padD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		Line(layer2EndL[0], layer2EndL[1], hole4[0], hole4[1], RectShape, *trace),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
		// Layer 2 inner connecting holes
		Circle(hole6, viaPadD),
		Circle(hole7, viaPadD),
		Circle(hole1, viaPadD),
		Circle(hole3, viaPadD),
		Line(startLayer2R[0], startLayer2R[1], hole3[0], hole3[1], RectShape, *trace),
		Line(startLayer2L[0], startLayer2L[1], hole1[0], hole1[1], RectShape, *trace),
		Line(hole1[0], hole1[1], hole3[0], hole3[1], RectShape, *trace),
		// Layer 2 outer connecting hole
		Circle(hole8, padD),
		Line(layer2EndR[0], layer2EndR[1], hole8[0], hole8[1], RectShape, *trace),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, padD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
		// Layer 2 inner connecting holes
		Circle(hole6, viaPadD),
		Circle(hole7, viaPadD),
		// Layer 2 outer connecting hole
		Circle(hole8, padD),
	)

	drill := g.Drill()
	drill.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaDrillD),
		Circle(hole2, drillD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaDrillD),
		Circle(hole4, drillD),
		// Lower connecting trace for right spiral
		Circle(hole5, drillD),
		// Layer 2 inner connecting holes
		Circle(hole6, viaDrillD),
		Circle(hole7, viaDrillD),
		// Layer 2 outer connecting hole
		Circle(hole8, drillD),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		labelSize := pts
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.5*r, 1.0, message, *fontName, pts, &Center),
			Text(hole1[0]+viaPadD, hole1[1], 1.0, "2R", *fontName, labelSize, &CenterLeft),
			Text(hole2[0]+viaPadD, hole2[1], 1.0, "BR/TL", *fontName, labelSize, &CenterLeft),
			Text(hole3[0]-viaPadD, hole3[1], 1.0, "2L", *fontName, labelSize, &CenterRight),
			Text(hole4[0]-padD, hole4[1], 1.0, "BL/2L", *fontName, labelSize, &CenterRight),
			Text(hole5[0]-padD, hole5[1], 1.0, "TR", *fontName, labelSize, &CenterRight),
			Text(hole6[0], hole6[1]+viaPadD, 1.0, "TR/BR", *fontName, labelSize, &BottomCenter),
			Text(hole7[0], hole7[1]-viaPadD, 1.0, "TL/BL", *fontName, labelSize, &TopCenter),
			Text(hole8[0], hole8[1]-viaPadD, 1.0, "2R", *fontName, labelSize, &TopCenter),
			Text(0, -0.5*r, 1.0, message2, *fontName, pts, &Center),
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
	startAngle := 2.5 * math.Pi
	endAngle := 2*math.Pi + float64(*n)*2.0*math.Pi
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

func (s *spiral) genSpiral(xScale, offset, trimY float64) (pts []Pt, endPt Pt) {
	halfTW := *trace * 0.5
	endAngle := s.endAngle
	if trimY < 0 { // Only for layer2SpiralL - extend another Pi/2
		endAngle += 0.5 * math.Pi
	}
	steps := int(0.5 + (endAngle-s.startAngle) / *step)
	for i := 0; i < steps; i++ {
		angle := s.startAngle + *step*float64(i)
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
	} else if trimY < 0 { // Only for layer2SpiralL
		trimYsteps++
		for {
			if pts[len(pts)-trimYsteps][1] < trimY {
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
	return pts, endPt
}
