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
//
// Copyright 2019 Glenn M. Lewis. All Rights Reserved.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"strings"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
	"github.com/gmlewis/go-gerber/gerber/viewer"
)

var (
	width  = flag.Float64("width", 400.0, "Width of PCB")
	height = flag.Float64("height", 400.0, "Height of PCB")
	step   = flag.Float64("step", 0.01, "Resolution (in radians) of the spiral")
	gap    = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	padGap = flag.Float64("pad_gap", 0.2, "Gap between pads in mm")
	trace  = flag.Float64("trace", 2.0, "Width of traces in mm")
	prefix = flag.String("prefix", "bifilar-cap",
		"Filename prefix for all Gerber files and zip")
	fontName = flag.String("font", "freeserif",
		"Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
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

	startT, endT, spiralT := s.genSpiral(0, 0)
	startB, endB, spiralB := s.genSpiral(0.5*math.Pi, 0.5)

	centerHole := Point(0.5**width, 0.5**height)
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, padD)
	}
	outerContact := func(pt Pt) Pt {
		x := pt[0] - 0.5**width
		y := pt[1] - 0.5**height
		r := math.Sqrt(x*x+y*y) + *trace + *gap
		angle := math.Atan2(y, x)
		return Point(0.5**width+r*math.Cos(angle),
			0.5**height+r*math.Sin(angle))
	}

	topOuter := outerContact(endT)
	botOuter := outerContact(endB)

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, spiralT, 0.0),
		Circle(startT, padD),
		Circle(centerHole, padD),
		Circle(topOuter, padD),
		Circle(botOuter, padD),
		padLine(topOuter, endT),
	)

	topMask := g.TopSolderMask()
	topMask.Add(
		Circle(startT, padD),
		Circle(centerHole, padD),
		Circle(topOuter, padD),
		Circle(botOuter, padD),
	)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, spiralB, 0.0),
		Circle(startT, padD),
		Circle(centerHole, padD),
		padLine(startB, centerHole),
		Circle(topOuter, padD),
		Circle(botOuter, padD),
		padLine(botOuter, endB),
	)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		Circle(startT, padD),
		Circle(centerHole, padD),
		Circle(topOuter, padD),
		Circle(botOuter, padD),
	)

	drill := g.Drill()
	drill.Add(
		Circle(startT, drillD),
		Circle(centerHole, drillD),
		Circle(topOuter, drillD),
		Circle(botOuter, drillD),
	)

	outline := g.Outline()
	border := []Pt{{0, 0}, {*width, 0}, {*width, *height}, {0, *height}}
	outline.Add(
		Line(border[0][0], border[0][1], border[1][0], border[1][1], CircleShape, 0.1),
		Line(border[1][0], border[1][1], border[2][0], border[2][1], CircleShape, 0.1),
		Line(border[2][0], border[2][1], border[3][0], border[3][1], CircleShape, 0.1),
		Line(border[3][0], border[3][1], border[0][0], border[0][1], CircleShape, 0.1),
	)

	// Now populate the board with breadboard points...
	d := *trace + *padGap
	for y := 0.75 * *trace; y <= *height-0.5**trace; y += d {
		ry := y - 0.5**height
		n := -1
		created := map[int]bool{}
		for x := 0.75 * *trace; x <= *height-0.5**trace; x += d {
			n++
			rx := x - 0.5**width
			r := math.Sqrt(rx*rx+ry*ry) - *padGap
			if r-*trace-4**padGap <= 0.5**width {
				continue
			}
			created[n] = true
			c := Circle(Point(x, y), padD)
			top.Add(c)
			topMask.Add(c)
			bottom.Add(c)
			bottomMask.Add(c)
			drill.Add(Circle(Point(x, y), drillD))
			if created[n-1] && n%2 == 1 {
				line := padLine(Point(x-d, y), Point(x, y))
				top.Add(line)
				bottom.Add(line)
			}
		}
	}

	if *fontName != "" {
		buf, err := ioutil.ReadFile("main.go")
		if err != nil {
			log.Fatalf("ReadFile: %v", err)
		}
		lines := strings.Split(string(buf), "\n")
		quarter := len(lines) / 4
		if quarter*4 < len(lines) {
			quarter++
		}
		t1 := strings.Join(lines[0:quarter], "\n")
		t2 := strings.Join(lines[quarter:2*quarter], "\n")
		t3 := strings.Join(lines[2*quarter:3*quarter], "\n")
		t4 := strings.Join(lines[3*quarter:], "\n")

		const margin = 3
		mbbL := MBB{Min: Pt{margin, margin},
			Max: Pt{0.5**width - margin, *height - margin}}
		mbbR := MBB{Min: Pt{0.5**width + margin, margin},
			Max: Pt{*width - margin, *height - margin}}
		tss := g.TopSilkscreen()
		tss.Add(
			TextBox(mbbL, 1.0, t1, *fontName, &Center),
			TextBox(mbbR, 1.0, t2, *fontName, &Center),
		)

		bss := g.BottomSilkscreen()
		bss.Add(
			TextBox(mbbR, -1.0, t3, *fontName, &Center),
			TextBox(mbbL, -1.0, t4, *fontName, &Center),
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

func (s *spiral) genSpiral(
	startAngleOffset, endAngleOffset float64) (startPt, endPt Pt, pts []Pt) {
	start := s.startAngle + startAngleOffset
	end := s.endAngle + endAngleOffset
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
