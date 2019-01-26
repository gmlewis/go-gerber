// bifilar-coil creates Gerber files (and a bundled ZIP) representing
// a bifilar coil (https://en.wikipedia.org/wiki/Bifilar_coil) for
// manufacture on a printed circuit board (PCB).
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	. "github.com/gmlewis/go-gerber/gerber"
)

var (
	step     = flag.Float64("step", 0.01, "Resolution (in radians) of the spiral")
	n        = flag.Int("n", 100, "Number of full winds in each spiral")
	gap      = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace    = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix   = flag.String("prefix", "bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName = flag.String("font", "ubuntumonoregular", "Name of font to use for writing source on PCB (empty to not write)")
)

func main() {
	flag.Parse()

	g := New(*prefix)

	s := newSpiral()
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, s.size, s.size)

	spiralR := s.genSpiral(0)
	spiralL := s.genSpiral(math.Pi)
	startR := genPt(s.startAngle, 0, 0)
	endR := genPt(s.endAngle, 0, 0)
	startL := genPt(s.startAngle, 0, math.Pi)
	endL := genPt(s.endAngle, 0, math.Pi)

	padD := 0.5
	drillD := 0.25
	padOffset := 0.5 * (padD - *trace)

	// Lower connecting trace between two spirals
	hole1 := Point(startR.X, startR.Y+padOffset)
	hole2 := Point(endL.X-padOffset, endL.Y)
	// Upper connecting trace for left spiral
	hole3 := Point(startL.X, startL.Y-padOffset)
	hole4 := Point(endR.X+padOffset, startL.Y+padOffset)
	// Lower connecting trace for right spiral
	hole5 := Point(endR.X+padOffset, endR.Y)

	top := g.TopCopper()
	top.Add(
		Polygon(0, 0, true, spiralR, 0.0),
		Polygon(0, 0, true, spiralL, 0.0),
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, padD),
		Circle(hole2.X, hole2.Y, padD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, padD),
		Circle(hole4.X, hole4.Y, padD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, padD),
		Circle(hole2.X, hole2.Y, padD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, padD),
		Circle(hole4.X, hole4.Y, padD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, padD),
		Line(startR.X, startR.Y, endL.X, startR.Y, RectShape, *trace),
		Line(endL.X, startR.Y, endL.X, endL.Y, RectShape, *trace),
		Circle(hole2.X, hole2.Y, padD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, padD),
		Line(startL.X, startL.Y, endR.X+padOffset, startL.Y, RectShape, *trace),
		Circle(hole4.X, hole4.Y, padD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, padD),
		Circle(hole2.X, hole2.Y, padD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, padD),
		Circle(hole4.X, hole4.Y, padD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, padD),
	)

	drill := g.Drill()
	drill.Add(
		// Lower connecting trace between two spirals
		Circle(hole1.X, hole1.Y, drillD),
		Circle(hole2.X, hole2.Y, drillD),
		// Upper connecting trace for left spiral
		Circle(hole3.X, hole3.Y, drillD),
		Circle(hole4.X, hole4.Y, drillD),
		// Lower connecting trace for right spiral
		Circle(hole5.X, hole5.Y, drillD),
	)

	outline := g.Outline()
	outline.Add(
		Arc(0, 0, 0.5*s.size+padD, CircleShape, 1, 1, 0, 360, 0.1),
	)

	if *fontName != "" { // Write the code on the top and bottom silkscreens (split in half).
		buf, err := ioutil.ReadFile("main.go")
		if err != nil {
			log.Fatal(err)
		}
		lines := strings.Split(string(buf), "\n")
		half := len(lines) / 2
		x, y := -40000.0, 100000.0

		tss := g.TopSilkscreen()
		tss.Add(
			Text(x, y, 1.0, strings.Join(lines[:half], "\n"), *fontName),
		)

		bss := g.BottomSilkscreen()
		bss.Add(
			Text(-x, y, -1.0, strings.Join(lines[half:], "\n"), *fontName),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
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
	size := 2 * math.Abs(p1.X)
	p2 := genPt(endAngle, *trace*0.5, math.Pi)
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
