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
	nlayers    = 20
	angleDelta = 2.0 * math.Pi / nlayers

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

	if *n < 20 {
		flag.Usage()
		log.Fatal("N must be >= 20.")
	}

	g := New(*prefix)

	s := newSpiral()

	padD := 2.0
	startTopR, topSpiralR, endTopR := s.genSpiral(1, 0, 0)
	startTopL, topSpiralL, endTopL := s.genSpiral(1, math.Pi, 0)
	startBotR, botSpiralR, endBotR := s.genSpiral(-1, 0, 0)
	startBotL, botSpiralL, endBotL := s.genSpiral(-1, math.Pi, 0) // *trace+padD)

	startLayerNR := map[int]Pt{}
	startLayerNL := map[int]Pt{}
	endLayerNR := map[int]Pt{}
	endLayerNL := map[int]Pt{}
	layerNSpiralR := map[int][]Pt{}
	layerNSpiralL := map[int][]Pt{}
	for n := 2; n < nlayers; n += 4 {
		af := float64((n + 2) / 4)
		startLayerNR[n], layerNSpiralR[n], endLayerNR[n] = s.genSpiral(1, af*angleDelta, 0)
		startLayerNL[n], layerNSpiralL[n], endLayerNL[n] = s.genSpiral(1, math.Pi+af*angleDelta, 0)
		startLayerNR[n+1], layerNSpiralR[n+1], endLayerNR[n+1] = s.genSpiral(-1, af*angleDelta, 0)
		startLayerNL[n+1], layerNSpiralL[n+1], endLayerNL[n+1] = s.genSpiral(-1, math.Pi+af*angleDelta, 0)

		if n+2 < nlayers {
			startLayerNR[n+2], layerNSpiralR[n+2], endLayerNR[n+2] = s.genSpiral(1, -af*angleDelta, 0)
			startLayerNL[n+2], layerNSpiralL[n+2], endLayerNL[n+2] = s.genSpiral(1, math.Pi-af*angleDelta, 0)
			startLayerNR[n+3], layerNSpiralR[n+3], endLayerNR[n+3] = s.genSpiral(-1, -af*angleDelta, 0)
			trimY := 0.0
			if n+3 == 5 { // 5L
				trimY = *trace + padD
			}
			startLayerNL[n+3], layerNSpiralL[n+3], endLayerNL[n+3] = s.genSpiral(-1, math.Pi-af*angleDelta, trimY)
		}
	}

	viaPadD := 0.5
	innerR := (*gap + viaPadD) / math.Sin(angleDelta)
	// minStartAngle := (innerR + *gap + 0.5**trace + 0.5*viaPadD) * (3 * math.Pi)
	// log.Printf("innerR=%v, minStartAngle=%v/Pi=%v", innerR, minStartAngle, minStartAngle/math.Pi)
	var innerViaPts []Pt
	for i := 0; i < nlayers; i++ {
		x := innerR * math.Cos(float64(i)*angleDelta)
		y := innerR * math.Sin(float64(i)*angleDelta)
		innerViaPts = append(innerViaPts, Pt{x, y})
	}
	innerHole := map[string]int{
		"TR": 17, "TL": 7, "BR": 13, "BL": 3,
		"2R": 18, "2L": 8, "3R": 12, "3L": 2,
		"4R": 16, "4L": 6, "5R": 14, "5L": 4,
		"6R": 19, "6L": 9, "7R": 11, "7L": 1,
		"8R": 15, "8L": 5, "9R": 15, "9L": 5,
		"10R": 0, "10L": 10, "11R": 10, "11L": 0,
		"12R": 14, "12L": 4, "13R": 16, "13L": 6,
		"14R": 1, "14L": 11, "15R": 9, "15L": 19,
		"16R": 13, "16L": 3, "17R": 17, "17L": 7,
		"18R": 2, "18L": 12, "19R": 8, "19L": 18,
	}

	outerR := (2.0*math.Pi + float64(*n)*2.0*math.Pi + *trace + *gap) / (3.0 * math.Pi)
	outerContactPt := func(n float64) Pt {
		r := outerR + 0.5**trace + *gap + 0.5*padD
		x := r * math.Cos(n*angleDelta)
		y := r * math.Sin(n*angleDelta)
		return Pt{x, y}
	}

	var outerViaPts []Pt
	for i := 0; i < nlayers; i++ {
		pt := outerContactPt(float64(i) - 0.5)
		outerViaPts = append(outerViaPts, pt)
	}
	outerViaPts = append(outerViaPts, outerContactPt(0.0))
	outerHole := map[string]int{
		"TR": 0, "TL": 10, "BR": 9, "BL": 19,
		"2R": 1, "2L": 11, "3R": 8, "3L": 18,
		"4R": 19, "4L": 9, "5R": 10, "5L": 20,
		"6R": 2, "6L": 12, "7R": 7, "7L": 17,
		"8R": 18, "8L": 8, "9R": 11, "9L": 1,
		"10R": 3, "10L": 13, "11R": 6, "11L": 16,
		"12R": 17, "12L": 7, "13R": 12, "13L": 2,
		"14R": 4, "14L": 14, "15R": 5, "15L": 15,
		"16R": 16, "16L": 6, "17R": 13, "17L": 3,
		"18R": 5, "18L": 15, "19R": 4, "19L": 14,
	}

	drill := g.Drill()
	for _, pt := range innerViaPts {
		const viaDrillD = 0.25
		drill.Add(Circle(pt, viaDrillD))
	}
	for _, pt := range outerViaPts {
		const drillD = 1.0
		drill.Add(Circle(pt, drillD))
	}

	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, *trace)
	}
	addVias := func(layer *Layer) {
		for _, pt := range innerViaPts {
			layer.Add(Circle(pt, viaPadD))
		}
		for _, pt := range outerViaPts {
			layer.Add(Circle(pt, padD))
		}
	}

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, true, topSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, topSpiralL, 0.0),
		padLine(startTopR, innerViaPts[innerHole["TR"]]),
		padLine(startTopL, innerViaPts[innerHole["TL"]]),
		padLine(endTopR, outerViaPts[outerHole["TR"]]),
		padLine(endTopL, outerViaPts[outerHole["TL"]]),
	)
	addVias(top)

	topMask := g.TopSolderMask()
	addVias(topMask)

	for n := 2; n < nlayers; n++ {
		nr := fmt.Sprintf("%vR", n)
		nl := fmt.Sprintf("%vL", n)
		layer := g.LayerN(n)
		layer.Add(
			Polygon(Pt{0, 0}, true, layerNSpiralR[n], 0.0),
			Polygon(Pt{0, 0}, true, layerNSpiralL[n], 0.0),
			padLine(startLayerNR[n], innerViaPts[innerHole[nr]]),
			padLine(startLayerNL[n], innerViaPts[innerHole[nl]]),
			padLine(endLayerNR[n], outerViaPts[outerHole[nr]]),
			padLine(endLayerNL[n], outerViaPts[outerHole[nl]]),
		)
		addVias(layer)
	}

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, botSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, botSpiralL, 0.0),
		padLine(startBotR, innerViaPts[innerHole["BR"]]),
		padLine(startBotL, innerViaPts[innerHole["BL"]]),
		padLine(endBotR, outerViaPts[outerHole["BR"]]),
		padLine(endBotL, outerViaPts[outerHole["BL"]]),
	)
	addVias(bottom)

	bottomMask := g.BottomSolderMask()
	addVias(bottomMask)

	outline := g.Outline()
	r := 0.5*s.size + padD + *trace
	outline.Add(
		Arc(Pt{0, 0}, r, CircleShape, 1, 1, 0, 360, 0.1),
	)
	fmt.Printf("n=%v: (%.2f,%.2f)\n", *n, 2*r, 2*r)

	if *fontName != "" {
		pts := 36.0 * r / 139.18 // determined emperically
		labelSize := pts * 4.0 / 18.0
		outerLabelSize := 32.0 * r / 139.18 // determined emperically
		message := fmt.Sprintf(messageFmt, *trace, *gap, *n)

		innerLabel := func(label string) *TextT {
			num := float64(innerHole[label])
			r := innerR - viaPadD
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		innerLabel2 := func(label string) *TextT {
			num := float64(innerHole[label])
			r := innerR + viaPadD
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		outerLabel := func(label string) *TextT {
			num := float64(outerHole[label])
			if outerHole[label] != 20 {
				num -= 0.5
			}
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((0.3+num)*angleDelta)
			y := r * math.Sin((0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}
		outerLabel2 := func(label string) *TextT {
			num := float64(outerHole[label])
			if outerHole[label] == 20 {
				num = 0.3
			} else {
				num -= 0.5
			}
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((-0.3+num)*angleDelta)
			y := r * math.Sin((-0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*r, 1.0, message, *fontName, pts, &Center),
			innerLabel("TR"), innerLabel("TL"), innerLabel("BR"), innerLabel("BL"),
			innerLabel("2R"), innerLabel("2L"), innerLabel("3R"), innerLabel("3L"),
			innerLabel("4R"), innerLabel("4L"), innerLabel("5R"), innerLabel("5L"),
			innerLabel("6R"), innerLabel("6L"), innerLabel("7R"), innerLabel("7L"),
			innerLabel("8R"), innerLabel("8L"), innerLabel2("9R"), innerLabel2("9L"),
			innerLabel("10R"), innerLabel("10L"), innerLabel2("11R"), innerLabel2("11L"),
			innerLabel2("12R"), innerLabel2("12L"), innerLabel2("13R"), innerLabel2("13L"),
			innerLabel2("14R"), innerLabel2("14L"), innerLabel2("15R"), innerLabel2("15L"),
			innerLabel2("16R"), innerLabel2("16L"), innerLabel2("17R"), innerLabel2("17L"),
			innerLabel2("18R"), innerLabel2("18L"), innerLabel2("19R"), innerLabel2("19L"),

			outerLabel2("TR"), outerLabel("TL"), outerLabel2("BR"), outerLabel2("BL"),
			outerLabel("2R"), outerLabel("2L"), outerLabel("3R"), outerLabel("3L"),
			outerLabel2("4R"), outerLabel2("4L"), outerLabel2("5R"), outerLabel2("5L"),
			outerLabel("6R"), outerLabel("6L"), outerLabel("7R"), outerLabel("7L"),
			outerLabel2("8R"), outerLabel2("8L"), outerLabel2("9R"), outerLabel2("9L"),
			outerLabel("10R"), outerLabel("10L"), outerLabel("11R"), outerLabel("11L"),
			outerLabel2("12R"), outerLabel2("12L"), outerLabel2("13R"), outerLabel2("13L"),
			outerLabel("14R"), outerLabel("14L"), outerLabel("15R"), outerLabel("15L"),
			outerLabel2("16R"), outerLabel2("16L"), outerLabel2("17R"), outerLabel2("17L"),
			outerLabel("18R"), outerLabel("18L"), outerLabel2("19R"), outerLabel2("19L"),

			// Text(-0.5*r, -0.4*r, 1.0, message2, *fontName, pts, &Center),
			// Text(0.5*r, -0.4*r, 1.0, message3, *fontName, pts, &Center),
		)
	}

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g, false)
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
	endAngle := 2.0*math.Pi + float64(*n)*2.0*math.Pi + 0.15*math.Pi
	p1 := genPt(1.0, endAngle, *trace*0.5, 0)
	size := 2 * p1.Length()
	p2 := genPt(1.0, endAngle, *trace*0.5, math.Pi)
	if v := 2 * p2.Length(); v > size {
		size = v
	}
	return &spiral{
		startAngle: startAngle,
		endAngle:   endAngle,
		size:       size,
	}
}

func (s *spiral) genSpiral(xScale, offset, trimY float64) (startPt Pt, pts []Pt, endPt Pt) {
	halfTW := *trace * 0.5
	endAngle := s.endAngle - 4.0*math.Pi/nlayers
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
