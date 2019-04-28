// bifilar-with-capacitor creates Gerber files (and a bundled ZIP) representing
// a bifilar coil (https://en.wikipedia.org/wiki/Bifilar_coil) with the
// ability to connect the two windings with a capacitor (or not) for
// manufacture on a printed circuit board (PCB).
//
// This designs differs from the others in that a single coil is devoted
// to one layer and the board itself is the dielectric between the top and
// bottom coils.
//
// This design requires that an external wire connect the inner terminal
// of one coil to the outer terminal of the other coil (with the added
// benefit of being able to insert a tuning capacitor in between the
// two coils.)
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
	width      = flag.Float64("width", 100.0, "Width of PCB")
	height     = flag.Float64("height", 100.0, "Height of PCB")
	step       = flag.Float64("step", 0.01, "Resolution (in radians) of the spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 2.0, "Width of traces in mm")
	prefix     = flag.String("prefix", "bifilar-cap", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	messageFmt = `This is a single (2-layer)
bifilar coil with one coil per layer.
Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`

	padD   = 2
	padR   = padD / 2
	drillD = padD / 2
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

	startR, endR, spiralR := s.genSpiral(0)
	startL, endL, spiralL := s.genSpiral(math.Pi)

	centerHole := Point(0.5**width, 0.5**height)
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, padD)
	}

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, spiralR, 0.0),
		Circle(startR, padD),
		Circle(centerHole, padD),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		Circle(startR, padD),
		Circle(centerHole, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, spiralL, 0.0),
		Circle(startR, padD),
		Circle(centerHole, padD),
		padLine(startL, centerHole),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		Circle(startR, padD),
		Circle(centerHole, padD),
	)

	drill := g.Drill()
	drill.Add(
		Circle(startR, drillD),
		Circle(centerHole, drillD),
		Circle(endR, drillD),
		Circle(endL, drillD),
	)

	outline := g.Outline()
	border := []Pt{{0, 0}, {*width, 0}, {*width, *height}, {0, *height}}
	outline.Add(
		Line(border[0][0], border[0][1], border[1][0], border[1][1], CircleShape, 0.1),
		Line(border[1][0], border[1][1], border[2][0], border[2][1], CircleShape, 0.1),
		Line(border[2][0], border[2][1], border[3][0], border[3][1], CircleShape, 0.1),
		Line(border[3][0], border[3][1], border[0][0], border[0][1], CircleShape, 0.1),
	)

	if *fontName != "" {
		// pts := 36.0 * r / 139.18 // determined emperically
		// labelSize := pts
		// message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		// tss := g.TopSilkscreen()
		// tss.Add(
		// 	Text(0, 0.5*r, 1.0, message, *fontName, pts, &Center),
		// 	Text(hole2[0]+0.5*padD, hole2[1], 1.0, "TR/TL", *fontName, labelSize, &CenterLeft),
		// 	Text(hole4[0]-0.5*padD, hole4[1], 1.0, "TL", *fontName, labelSize, &CenterRight),
		// 	Text(hole5[0]-0.6*padD, hole5[1], 1.0, "TR", *fontName, labelSize, &CenterRight),
		// )
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g, true)
	}
}

func genPt(angle, rOffset, angleOffset float64) Pt {
	r := (*trace + *gap) * angle / (2.0 * math.Pi)
	x := 0.5**width + (r+rOffset)*math.Cos(angle+angleOffset)
	y := 0.5**height + (r+rOffset)*math.Sin(angle+angleOffset)
	return Point(x, y)
}

type spiral struct {
	startAngle float64
	endAngle   float64
	size       float64
}

func newSpiral() *spiral {
	startAngle := 2.5 * math.Pi
	n := math.Floor(0.5**width/(*trace+*gap)) - 0.375
	endAngle := 2.0 * math.Pi * n
	log.Printf("n=%v, start=%v, end=%v", n, genPt(startAngle, *trace*0.5, 0.0), genPt(endAngle, *trace*0.5, 0.0))

	p1 := genPt(endAngle, *trace*0.5, 0)
	size := 2 * math.Abs(p1[0])
	p2 := genPt(endAngle, *trace*0.5, math.Pi)
	xl := 2 * math.Abs(p2[0])
	if xl > size {
		size = xl
	}
	log.Printf("startAngle=%v, endAngle=%v, size=%v", startAngle, endAngle, size)
	return &spiral{
		startAngle: startAngle,
		endAngle:   endAngle,
		size:       size,
	}
}

func (s *spiral) genSpiral(angleOffset float64) (startPt, endPt Pt, pts []Pt) {
	start := s.startAngle + angleOffset
	end := s.endAngle + angleOffset
	halfTW := *trace * 0.5

	steps := int(0.5 + (end-start) / *step)
	for i := 0; i < steps; i++ {
		angle := start + *step*float64(i)
		pts = append(pts, genPt(angle, halfTW, 0.0))
	}
	pts = append(pts, genPt(end, halfTW, 0.0))
	pts = append(pts, genPt(end, -halfTW, 0.0))
	for i := steps - 1; i >= 0; i-- {
		angle := start + *step*float64(i)
		pts = append(pts, genPt(angle, -halfTW, 0.0))
	}
	pts = append(pts, pts[0])
	return genPt(start, 0.0, 0.0), genPt(end, 0.0, 0.0), pts
}
