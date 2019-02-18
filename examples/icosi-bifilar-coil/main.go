// icosi-bifilar-coil creates Gerber files (and a bundled ZIP) representing
// 20 bifilar coils (https://en.wikipedia.org/wiki/Bifilar_coil) (one on
// each layer of a 20-layer PCB) for manufacture on a printed circuit
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
	n          = flag.Int("n", 12, "Number of full winds in each spiral")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "icosi-bifilar-coil", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	ncoils     = 20
	angleDelta = 2.0 * math.Pi / ncoils

	messageFmt = `Trace size = %0.2fmm.
Gap size = %0.2fmm.
Each spiral has %v coils.`

// 	message2 = `3L ⇨ 4L
// 4L ⇨ BL
// BL ⇨ TL
// TL ⇨ 5L
// 5L ⇨ 2L
// 2L ⇨ 3R`
// 	message3 = `3R ⇨ 4R
// 4R ⇨ BR
// BR ⇨ TR
// TR ⇨ 5R
// 5R ⇨ 2R`
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
	log.Printf("%v, %v, %v", startTopR, len(topSpiralR), endTopR)
	startTopL, topSpiralL, endTopL := s.genSpiral(1.0, math.Pi, 0)
	log.Printf("%v, %v, %v", startTopL, len(topSpiralL), endTopL)
	startBotR, botSpiralR, endBotR := s.genSpiral(-1.0, 0, 0)
	log.Printf("%v, %v, %v", startBotR, len(botSpiralR), endBotR)
	startBotL, botSpiralL, endBotL := s.genSpiral(-1.0, math.Pi, 0)
	log.Printf("%v, %v, %v", startBotL, len(botSpiralL), endBotL)

	padD := 2.0
	startLayer2R, layer2SpiralR, endLayer2R := s.genSpiral(1.0, angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer2R, len(layer2SpiralR), endLayer2R)
	startLayer2L, layer2SpiralL, endLayer2L := s.genSpiral(1.0, math.Pi+angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer2L, len(layer2SpiralL), endLayer2L)
	startLayer3R, layer3SpiralR, endLayer3R := s.genSpiral(-1.0, angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer3R, len(layer3SpiralR), endLayer3R)
	startLayer3L, layer3SpiralL, endLayer3L := s.genSpiral(-1.0, math.Pi+angleDelta, 0) // *trace+padD)
	log.Printf("%v, %v, %v", startLayer3L, len(layer3SpiralL), endLayer3L)

	startLayer4R, layer4SpiralR, endLayer4R := s.genSpiral(1.0, -angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer4R, len(layer4SpiralR), endLayer4R)
	startLayer4L, layer4SpiralL, endLayer4L := s.genSpiral(1.0, math.Pi-angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer4L, len(layer4SpiralL), endLayer4L)
	startLayer5R, layer5SpiralR, endLayer5R := s.genSpiral(-1.0, -angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer5R, len(layer5SpiralR), endLayer5R)
	startLayer5L, layer5SpiralL, endLayer5L := s.genSpiral(-1.0, math.Pi-angleDelta, 0)
	log.Printf("%v, %v, %v", startLayer5L, len(layer5SpiralL), endLayer5L)

	viaPadD := 0.5
	innerR := (*gap + viaPadD) / math.Sin(angleDelta)
	// minStartAngle := (innerR + *gap + 0.5**trace + 0.5*viaPadD) * (3 * math.Pi)
	// log.Printf("innerR=%v, minStartAngle=%v/Pi=%v", innerR, minStartAngle, minStartAngle/math.Pi)
	var innerViaPts []Pt
	for i := 0; i < ncoils; i++ {
		x := innerR * math.Cos(float64(i)*angleDelta)
		y := innerR * math.Sin(float64(i)*angleDelta)
		innerViaPts = append(innerViaPts, Pt{x, y})
	}
	innerHoleTR := innerViaPts[17]
	innerHoleTL := innerViaPts[7]
	innerHoleBR := innerViaPts[13]
	innerHoleBL := innerViaPts[3]
	innerHole2R := innerViaPts[18]
	innerHole2L := innerViaPts[8]
	innerHole3R := innerViaPts[12]
	innerHole3L := innerViaPts[2]
	innerHole4R := innerViaPts[16]
	innerHole4L := innerViaPts[6]
	innerHole5R := innerViaPts[14]
	innerHole5L := innerViaPts[4]

	// innerHole1Y := 0.5 * (*trace + viaPadD) / math.Sin(math.Pi/ncoils)
	// innerHole6X := innerHole1Y * math.Cos(math.Pi/ncoils)
	// hole1 := Point(0, innerHole1Y)
	// hole3 := Point(0, -innerHole1Y)
	// hole6 := Point(innerHole6X, 0.5*(*trace+viaPadD))
	// hole7 := Point(-innerHole6X, -0.5*(*trace+viaPadD))
	// hole10 := Point(innerHole6X, -0.5*(*trace+viaPadD))
	// hole11 := Point(-innerHole6X, 0.5*(*trace+viaPadD))

	// outerViaPt := func(pt Pt, angle float64) Pt {
	// 	r := *trace*1.5 + 0.5*viaPadD
	// 	dx := r * math.Cos(angle)
	// 	dy := r * math.Sin(angle)
	// 	return Point(pt[0]+dx, pt[1]+dy)
	// }
	// log.Printf("%v", outerViaPt)
	outerR := (2.0*math.Pi + float64(*n)*2.0*math.Pi + *trace + *gap) / (3.0 * math.Pi)
	var outerViaPts []Pt
	for i := 0; i < ncoils; i++ {
		r := outerR + 0.5**trace + *gap + 0.5*viaPadD
		x := r * math.Cos(float64(i)*angleDelta)
		y := r * math.Sin(float64(i)*angleDelta)
		outerViaPts = append(outerViaPts, Pt{x, y})
	}
	outerHoleTR := outerViaPts[0]
	outerHoleTL := outerViaPts[10]
	outerHoleBR := outerViaPts[10]
	outerHoleBL := outerViaPts[0]
	outerHole2R := outerViaPts[1]
	outerHole2L := outerViaPts[11]
	outerHole3R := outerViaPts[9]
	outerHole3L := outerViaPts[19]
	outerHole4R := outerViaPts[19]
	outerHole4L := outerViaPts[9]
	outerHole5R := outerViaPts[11]
	outerHole5L := outerViaPts[1]

	// holeBL4L := outerViaPt(endBotL, math.Pi/3.0)
	// holeTL5L := outerViaPt(endTopL, 2.0*math.Pi/3.0)
	// hole2L3R := outerViaPt(endLayer2L, math.Pi)
	// holeBR4R := outerViaPt(endBotR, 4.0*math.Pi/3.0)
	// holeTR5R := outerViaPt(endTopR, 5.0*math.Pi/3.0)

	outerContactPt := func(pt Pt, angle float64) Pt {
		r := *trace*1.5 + 0.5*padD
		dx := r * math.Cos(angle)
		dy := r * math.Sin(angle)
		return Point(pt[0]+dx, pt[1]+dy)
	}
	log.Printf("%v", outerContactPt)

	// hole2R := outerContactPt(endLayer2R, 0)
	// hole3L := outerContactPt(endLayer3L, 0)

	viaDrill := func(pt Pt) *CircleT {
		const viaDrillD = 0.25
		return Circle(pt, viaDrillD)
	}
	contactDrill := func(pt Pt) *CircleT {
		const drillD = 1.0
		return Circle(pt, drillD)
	}
	log.Printf("%v", contactDrill)

	drill := g.Drill()
	// drill.Add(
	//   contactDrill(hole2R),
	//   contactDrill(hole3L),
	// )
	for _, pt := range innerViaPts {
		drill.Add(viaDrill(pt))
	}
	for _, pt := range outerViaPts {
		drill.Add(viaDrill(pt))
	}

	viaPad := func(pt Pt) *CircleT {
		return Circle(pt, viaPadD)
	}
	contactPad := func(pt Pt) *CircleT {
		return Circle(pt, padD)
	}
	log.Printf("%v", contactPad)
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, *trace)
	}
	log.Printf("%v", padLine)

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, topSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, topSpiralL, 0.0),
		padLine(startTopR, innerHoleTR),
		padLine(startTopL, innerHoleTL),
		padLine(endTopR, outerHoleTR),
		padLine(endTopL, outerHoleTL),
	//
	// viaPad(hole1),
	// padLine(startTopL, hole1),
	//
	// viaPad(hole3),
	// padLine(startTopR, hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// padLine(endTopL, holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// padLine(endTopR, holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		top.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		top.Add(viaPad(pt))
	}

	layer2 := g.Layer2()
	layer2.Add(
		Polygon(Pt{0, 0}, true, layer2SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer2SpiralL, 0.0),
		padLine(startLayer2R, innerHole2R),
		padLine(startLayer2L, innerHole2L),
		padLine(endLayer2R, outerHole2R),
		padLine(endLayer2L, outerHole2L),
	//
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// padLine(startLayer2R, hole10),
	// viaPad(hole11),
	// padLine(startLayer2L, hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// padLine(endLayer2L, hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// padLine(endLayer2R, hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		layer2.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		layer2.Add(viaPad(pt))
	}

	layer4 := g.Layer4()
	layer4.Add(
		Polygon(Pt{0, 0}, true, layer4SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer4SpiralL, 0.0),
		padLine(startLayer4R, innerHole4R),
		padLine(startLayer4L, innerHole4L),
		padLine(endLayer4R, outerHole4R),
		padLine(endLayer4L, outerHole4L),
	//
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// padLine(startLayer4L, hole6),
	// viaPad(hole7),
	// padLine(startLayer4R, hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// padLine(endLayer4L, holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// padLine(endLayer4R, holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		layer4.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		layer4.Add(viaPad(pt))
	}

	topMask := g.TopSolderMask()
	topMask.Add(
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		topMask.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		topMask.Add(viaPad(pt))
	}

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, botSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, botSpiralL, 0.0),
		padLine(startBotR, innerHoleBR),
		padLine(startBotL, innerHoleBL),
		padLine(endBotR, outerHoleBR),
		padLine(endBotL, outerHoleBL),
	//
	// viaPad(hole1),
	// padLine(startBotL, hole1),
	//
	// viaPad(hole3),
	// padLine(startBotR, hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// padLine(endBotL, holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// padLine(endBotR, holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		bottom.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		bottom.Add(viaPad(pt))
	}

	layer3 := g.Layer3()
	layer3.Add(
		Polygon(Pt{0, 0}, true, layer3SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer3SpiralL, 0.0),
		padLine(startLayer3R, innerHole3R),
		padLine(startLayer3L, innerHole3L),
		padLine(endLayer3R, outerHole3R),
		padLine(endLayer3L, outerHole3L),
	//
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// padLine(startLayer3L, hole6),
	// viaPad(hole7),
	// padLine(startLayer3R, hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// padLine(endLayer3R, hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	// padLine(endLayer3L, hole3L),
	)
	for _, pt := range innerViaPts {
		layer3.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		layer3.Add(viaPad(pt))
	}

	layer5 := g.Layer5()
	layer5.Add(
		Polygon(Pt{0, 0}, true, layer5SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer5SpiralL, 0.0),
		padLine(startLayer5R, innerHole5R),
		padLine(startLayer5L, innerHole5L),
		padLine(endLayer5R, outerHole5R),
		padLine(endLayer5L, outerHole5L),
	//
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// padLine(startLayer5R, hole10),
	// viaPad(hole11),
	// padLine(startLayer5L, hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// padLine(endLayer5L, holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// padLine(endLayer5R, holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		layer5.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		layer5.Add(viaPad(pt))
	}

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
	// viaPad(hole1),
	//
	// viaPad(hole3),
	//
	// viaPad(hole6),
	// viaPad(hole7),
	//
	// viaPad(hole10),
	// viaPad(hole11),
	//
	// viaPad(holeBL4L),
	// viaPad(holeTL5L),
	// viaPad(hole2L3R),
	// viaPad(holeBR4R),
	// viaPad(holeTR5R),
	// contactPad(hole2R),
	// contactPad(hole3L),
	)
	for _, pt := range innerViaPts {
		bottomMask.Add(viaPad(pt))
	}
	for _, pt := range outerViaPts {
		bottomMask.Add(viaPad(pt))
	}

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		labelSize := pts * 4.0 / 18.0
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		innerLabel := func(label string, num int) *TextT {
			r := innerR - viaPadD
			x := r * math.Cos(float64(num)*angleDelta)
			y := r * math.Sin(float64(num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		outerLabel := func(label string, num int) *TextT {
			r := outerR + 0.5**trace + *gap + 1.5*viaPadD
			x := r * math.Cos(float64(num)*angleDelta)
			y := r * math.Sin(float64(num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		outerLabel2 := func(label string, num int) *TextT {
			r := outerR + 0.5**trace + *gap + 2.5*viaPadD
			x := r * math.Cos(float64(num)*angleDelta)
			y := r * math.Sin(float64(num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*r, 1.0, message, *fontName, pts, &Center),
			innerLabel("TR", 17),
			innerLabel("TL", 7),
			innerLabel("BR", 13),
			innerLabel("BL", 3),
			innerLabel("2R", 18),
			innerLabel("2L", 8),
			innerLabel("3R", 12),
			innerLabel("3L", 2),
			innerLabel("4R", 16),
			innerLabel("4L", 6),
			innerLabel("5R", 14),
			innerLabel("5L", 4),

			outerLabel("TR", 0),
			outerLabel("TL", 10),
			outerLabel2("BR", 10),
			outerLabel2("BL", 0),
			outerLabel("2R", 1),
			outerLabel("2L", 11),
			outerLabel("3R", 9),
			outerLabel("3L", 19),
			outerLabel2("4R", 19),
			outerLabel2("4L", 9),
			outerLabel2("5R", 11),
			outerLabel2("5L", 1),
			// Text(-0.5*r, -0.4*r, 1.0, message2, *fontName, pts, &Center),
			// Text(0.5*r, -0.4*r, 1.0, message3, *fontName, pts, &Center),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g)
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
	startAngle := 7.71 * math.Pi
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
	endAngle := s.endAngle // - math.Pi/3.0
	if trimY < 0 {         // Only for layer2SpiralL - extend another Pi/2
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
	return startPt, pts, endPt
}
