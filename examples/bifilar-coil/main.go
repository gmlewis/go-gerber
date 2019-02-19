// bifilar-coil creates Gerber files (and a bundled ZIP) representing
// a bifilar coil (https://en.wikipedia.org/wiki/Bifilar_coil) for
// manufacture on a printed circuit board (PCB).
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
	step       = flag.Float64("step", 0.01, "Resolution (in radians) of the spiral")
	n          = flag.Int("n", 100, "Number of full winds in each spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	message = `With a trace and gap size of 0.15mm, this
bifilar coil has a DC resistance of 232.2Ω.
Each spiral has 100 coils.`
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

	spiralR := s.genSpiral(0)
	spiralL := s.genSpiral(math.Pi)
	startR := genPt(s.startAngle, 0, 0)
	endR := genPt(s.endAngle, 0, 0)
	startL := genPt(s.startAngle, 0, math.Pi)
	endL := genPt(s.endAngle, 0, math.Pi)

	viaPadD := 0.5
	viaDrillD := 0.25
	viaPadOffset := 0.5 * (viaPadD - *trace)

	padD := 2.0
	drillD := 1.0
	padOffset := 0.5 * (padD - *trace)

	// Lower connecting trace between two spirals
	hole1 := Point(startR[0], startR[1]+viaPadOffset)
	hole2 := Point(endL[0]-viaPadOffset, endL[1])
	// Upper connecting trace for left spiral
	hole3 := Point(startL[0], startL[1]-viaPadOffset)
	hole4 := Point(endR[0]+padOffset, startL[1]+2*padOffset)
	// Lower connecting trace for right spiral
	hole5 := Point(endR[0]+padOffset, endR[1])

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, spiralR, 0.0),
		Polygon(Pt{0, 0}, true, spiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Line(startR[0], startR[1], endL[0], startR[1], RectShape, *trace),
		Line(endL[0], startR[1], endL[0], endL[1], RectShape, *trace),
		Circle(hole2, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Line(startL[0], startL[1], startL[0], startL[1]+padOffset, RectShape, *trace),
		Line(startL[0], startL[1]+padOffset, endR[0]+padOffset, startL[1]+padOffset, RectShape, *trace),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaPadD),
		Circle(hole2, viaPadD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaPadD),
		Circle(hole4, padD),
		// Lower connecting trace for right spiral
		Circle(hole5, padD),
	)

	drill := g.Drill()
	drill.Add(
		// Lower connecting trace between two spirals
		Circle(hole1, viaDrillD),
		Circle(hole2, viaDrillD),
		// Upper connecting trace for left spiral
		Circle(hole3, viaDrillD),
		Circle(hole4, drillD),
		// Lower connecting trace for right spiral
		Circle(hole5, drillD),
	)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		y := 0.3 * r
		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, y, 1.0, message, *fontName, pts, &Center),
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

func genPt(angle, halfTW, offset float64) Pt {
	r := (angle + *trace + *gap) / (3 * math.Pi)
	x := (r + halfTW) * math.Cos(angle+offset)
	y := (r + halfTW) * math.Sin(angle+offset)
	return Point(x, y)
}

type spiral struct {
	startAngle float64
	endAngle   float64
	size       float64
}

func newSpiral() *spiral {
	startAngle := 1.5 * math.Pi
	endAngle := 2*math.Pi + float64(*n)*2.0*math.Pi
	p1 := genPt(endAngle, *trace*0.5, 0)
	size := 2 * math.Abs(p1[0])
	p2 := genPt(endAngle, *trace*0.5, math.Pi)
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

func (s *spiral) genSpiral(offset float64) []Pt {
	halfTW := *trace * 0.5
	var pts []Pt
	steps := int(0.5 + (s.endAngle-s.startAngle) / *step)
	for i := 0; i < steps; i++ {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(angle, halfTW, offset))
	}
	pts = append(pts, genPt(s.endAngle, halfTW, offset))
	pts = append(pts, genPt(s.endAngle, -halfTW, offset))
	for i := steps - 1; i >= 0; i-- {
		angle := s.startAngle + *step*float64(i)
		pts = append(pts, genPt(angle, -halfTW, offset))
	}
	pts = append(pts, pts[0])
	return pts
}
