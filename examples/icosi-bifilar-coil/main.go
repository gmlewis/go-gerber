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

	padD := 2.0
	startTopR, topSpiralR, endTopR := s.genSpiral(1, 0, 0)
	startTopL, topSpiralL, endTopL := s.genSpiral(1, math.Pi, 0)
	startBotR, botSpiralR, endBotR := s.genSpiral(-1, 0, 0)
	startBotL, botSpiralL, endBotL := s.genSpiral(-1, math.Pi, *trace+padD)

	startLayer2R, layer2SpiralR, endLayer2R := s.genSpiral(1, angleDelta, 0)
	startLayer2L, layer2SpiralL, endLayer2L := s.genSpiral(1, math.Pi+angleDelta, 0)
	startLayer3R, layer3SpiralR, endLayer3R := s.genSpiral(-1, angleDelta, 0)
	startLayer3L, layer3SpiralL, endLayer3L := s.genSpiral(-1, math.Pi+angleDelta, 0)

	startLayer4R, layer4SpiralR, endLayer4R := s.genSpiral(1, -angleDelta, 0)
	startLayer4L, layer4SpiralL, endLayer4L := s.genSpiral(1, math.Pi-angleDelta, 0)
	startLayer5R, layer5SpiralR, endLayer5R := s.genSpiral(-1, -angleDelta, 0)
	startLayer5L, layer5SpiralL, endLayer5L := s.genSpiral(-1, math.Pi-angleDelta, 0)

	startLayer6R, layer6SpiralR, endLayer6R := s.genSpiral(1, 2*angleDelta, 0)
	startLayer6L, layer6SpiralL, endLayer6L := s.genSpiral(1, math.Pi+2*angleDelta, 0)
	startLayer7R, layer7SpiralR, endLayer7R := s.genSpiral(-1, 2*angleDelta, 0)
	startLayer7L, layer7SpiralL, endLayer7L := s.genSpiral(-1, math.Pi+2*angleDelta, 0)

	startLayer8R, layer8SpiralR, endLayer8R := s.genSpiral(1, -2*angleDelta, 0)
	startLayer8L, layer8SpiralL, endLayer8L := s.genSpiral(1, math.Pi-2*angleDelta, 0)
	startLayer9R, layer9SpiralR, endLayer9R := s.genSpiral(-1, -2*angleDelta, 0)
	startLayer9L, layer9SpiralL, endLayer9L := s.genSpiral(-1, math.Pi-2*angleDelta, 0)

	startLayer10R, layer10SpiralR, endLayer10R := s.genSpiral(1, 3*angleDelta, 0)
	startLayer10L, layer10SpiralL, endLayer10L := s.genSpiral(1, math.Pi+3*angleDelta, 0)
	startLayer11R, layer11SpiralR, endLayer11R := s.genSpiral(-1, 3*angleDelta, 0)
	startLayer11L, layer11SpiralL, endLayer11L := s.genSpiral(-1, math.Pi+3*angleDelta, 0)

	startLayer12R, layer12SpiralR, endLayer12R := s.genSpiral(1, -3*angleDelta, 0)
	startLayer12L, layer12SpiralL, endLayer12L := s.genSpiral(1, math.Pi-3*angleDelta, 0)
	startLayer13R, layer13SpiralR, endLayer13R := s.genSpiral(-1, -3*angleDelta, 0)
	startLayer13L, layer13SpiralL, endLayer13L := s.genSpiral(-1, math.Pi-3*angleDelta, 0)

	startLayer14R, layer14SpiralR, endLayer14R := s.genSpiral(1, 4*angleDelta, 0)
	startLayer14L, layer14SpiralL, endLayer14L := s.genSpiral(1, math.Pi+4*angleDelta, 0)
	startLayer15R, layer15SpiralR, endLayer15R := s.genSpiral(-1, 4*angleDelta, 0)
	startLayer15L, layer15SpiralL, endLayer15L := s.genSpiral(-1, math.Pi+4*angleDelta, 0)

	startLayer16R, layer16SpiralR, endLayer16R := s.genSpiral(1, -4*angleDelta, 0)
	startLayer16L, layer16SpiralL, endLayer16L := s.genSpiral(1, math.Pi-4*angleDelta, 0)
	startLayer17R, layer17SpiralR, endLayer17R := s.genSpiral(-1, -4*angleDelta, 0)
	startLayer17L, layer17SpiralL, endLayer17L := s.genSpiral(-1, math.Pi-4*angleDelta, 0)

	startLayer18R, layer18SpiralR, endLayer18R := s.genSpiral(1, -5*angleDelta, 0)
	startLayer18L, layer18SpiralL, endLayer18L := s.genSpiral(1, math.Pi-5*angleDelta, 0)
	startLayer19R, layer19SpiralR, endLayer19R := s.genSpiral(-1, -5*angleDelta, 0)
	startLayer19L, layer19SpiralL, endLayer19L := s.genSpiral(-1, math.Pi-5*angleDelta, 0)

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
	innerHole := map[string]int{
		"TR":  17,
		"TL":  7,
		"BR":  13,
		"BL":  3,
		"2R":  18,
		"2L":  8,
		"3R":  12,
		"3L":  2,
		"4R":  16,
		"4L":  6,
		"5R":  14,
		"5L":  4,
		"6R":  19,
		"6L":  9,
		"7R":  11,
		"7L":  1,
		"8R":  15,
		"8L":  5,
		"9R":  15,
		"9L":  5,
		"10R": 0,
		"10L": 10,
		"11R": 10,
		"11L": 0,
		"12R": 14,
		"12L": 4,
		"13R": 16,
		"13L": 6,
		"14R": 1,
		"14L": 11,
		"15R": 9,
		"15L": 19,
		"16R": 13,
		"16L": 3,
		"17R": 17,
		"17L": 7,
		"18R": 12,
		"18L": 2,
		"19R": 18,
		"19L": 8,
	}

	outerR := (2.0*math.Pi + float64(*n)*2.0*math.Pi + *trace + *gap) / (3.0 * math.Pi)
	outerContactPt := func(n float64) Pt {
		r := outerR + 0.5**trace + *gap + 0.5*padD
		x := r * math.Cos(n*angleDelta)
		y := r * math.Sin(n*angleDelta)
		return Pt{x, y}
	}

	var outerViaPts []Pt
	for i := 0; i < ncoils; i++ {
		pt := outerContactPt(float64(i))
		outerViaPts = append(outerViaPts, pt)
	}
	outerViaPts = append(outerViaPts, outerContactPt(0.5))
	outerHole := map[string]int{
		"TR":  0,
		"TL":  10,
		"BR":  10,
		"BL":  20,
		"2R":  1,
		"2L":  11,
		"3R":  9,
		"3L":  19,
		"4R":  19,
		"4L":  9,
		"5R":  11,
		"5L":  1,
		"6R":  2,
		"6L":  12,
		"7R":  8,
		"7L":  18,
		"8R":  18,
		"8L":  8,
		"9R":  12,
		"9L":  2,
		"10R": 3,
		"10L": 13,
		"11R": 7,
		"11L": 17,
		"12R": 17,
		"12L": 7,
		"13R": 13,
		"13L": 3,
		"14R": 4,
		"14L": 14,
		"15R": 6,
		"15L": 16,
		"16R": 16,
		"16L": 6,
		"17R": 14,
		"17L": 4,
		"18R": 15,
		"18L": 5,
		"19R": 15,
		"19L": 5,
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

	layer2 := g.LayerN(2)
	layer2.Add(
		Polygon(Pt{0, 0}, true, layer2SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer2SpiralL, 0.0),
		padLine(startLayer2R, innerViaPts[innerHole["2R"]]),
		padLine(startLayer2L, innerViaPts[innerHole["2L"]]),
		padLine(endLayer2R, outerViaPts[outerHole["2R"]]),
		padLine(endLayer2L, outerViaPts[outerHole["2L"]]),
	)
	addVias(layer2)

	layer4 := g.LayerN(4)
	layer4.Add(
		Polygon(Pt{0, 0}, true, layer4SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer4SpiralL, 0.0),
		padLine(startLayer4R, innerViaPts[innerHole["4R"]]),
		padLine(startLayer4L, innerViaPts[innerHole["4L"]]),
		padLine(endLayer4R, outerViaPts[outerHole["4R"]]),
		padLine(endLayer4L, outerViaPts[outerHole["4L"]]),
	)
	addVias(layer4)

	layer6 := g.LayerN(6)
	layer6.Add(
		Polygon(Pt{0, 0}, true, layer6SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer6SpiralL, 0.0),
		padLine(startLayer6R, innerViaPts[innerHole["6R"]]),
		padLine(startLayer6L, innerViaPts[innerHole["6L"]]),
		padLine(endLayer6R, outerViaPts[outerHole["6R"]]),
		padLine(endLayer6L, outerViaPts[outerHole["6L"]]),
	)
	addVias(layer6)

	layer8 := g.LayerN(8)
	layer8.Add(
		Polygon(Pt{0, 0}, true, layer8SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer8SpiralL, 0.0),
		padLine(startLayer8R, innerViaPts[innerHole["8R"]]),
		padLine(startLayer8L, innerViaPts[innerHole["8L"]]),
		padLine(endLayer8R, outerViaPts[outerHole["8R"]]),
		padLine(endLayer8L, outerViaPts[outerHole["8L"]]),
	)
	addVias(layer8)

	layer10 := g.LayerN(10)
	layer10.Add(
		Polygon(Pt{0, 0}, true, layer10SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer10SpiralL, 0.0),
		padLine(startLayer10R, innerViaPts[innerHole["10R"]]),
		padLine(startLayer10L, innerViaPts[innerHole["10L"]]),
		padLine(endLayer10R, outerViaPts[outerHole["10R"]]),
		padLine(endLayer10L, outerViaPts[outerHole["10L"]]),
	)
	addVias(layer10)

	layer12 := g.LayerN(12)
	layer12.Add(
		Polygon(Pt{0, 0}, true, layer12SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer12SpiralL, 0.0),
		padLine(startLayer12R, innerViaPts[innerHole["12R"]]),
		padLine(startLayer12L, innerViaPts[innerHole["12L"]]),
		padLine(endLayer12R, outerViaPts[outerHole["12R"]]),
		padLine(endLayer12L, outerViaPts[outerHole["12L"]]),
	)
	addVias(layer12)

	layer14 := g.LayerN(14)
	layer14.Add(
		Polygon(Pt{0, 0}, true, layer14SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer14SpiralL, 0.0),
		padLine(startLayer14R, innerViaPts[innerHole["14R"]]),
		padLine(startLayer14L, innerViaPts[innerHole["14L"]]),
		padLine(endLayer14R, outerViaPts[outerHole["14R"]]),
		padLine(endLayer14L, outerViaPts[outerHole["14L"]]),
	)
	addVias(layer14)

	layer16 := g.LayerN(16)
	layer16.Add(
		Polygon(Pt{0, 0}, true, layer16SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer16SpiralL, 0.0),
		padLine(startLayer16R, innerViaPts[innerHole["16R"]]),
		padLine(startLayer16L, innerViaPts[innerHole["16L"]]),
		padLine(endLayer16R, outerViaPts[outerHole["16R"]]),
		padLine(endLayer16L, outerViaPts[outerHole["16L"]]),
	)
	addVias(layer16)

	layer18 := g.LayerN(18)
	layer18.Add(
		Polygon(Pt{0, 0}, true, layer18SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer18SpiralL, 0.0),
		padLine(startLayer18R, innerViaPts[innerHole["18R"]]),
		padLine(startLayer18L, innerViaPts[innerHole["18L"]]),
		padLine(endLayer18R, outerViaPts[outerHole["18R"]]),
		padLine(endLayer18L, outerViaPts[outerHole["18L"]]),
	)
	addVias(layer18)

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

	layer3 := g.LayerN(3)
	layer3.Add(
		Polygon(Pt{0, 0}, true, layer3SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer3SpiralL, 0.0),
		padLine(startLayer3R, innerViaPts[innerHole["3R"]]),
		padLine(startLayer3L, innerViaPts[innerHole["3L"]]),
		padLine(endLayer3R, outerViaPts[outerHole["3R"]]),
		padLine(endLayer3L, outerViaPts[outerHole["3L"]]),
	)
	addVias(layer3)

	layer5 := g.LayerN(5)
	layer5.Add(
		Polygon(Pt{0, 0}, true, layer5SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer5SpiralL, 0.0),
		padLine(startLayer5R, innerViaPts[innerHole["5R"]]),
		padLine(startLayer5L, innerViaPts[innerHole["5L"]]),
		padLine(endLayer5R, outerViaPts[outerHole["5R"]]),
		padLine(endLayer5L, outerViaPts[outerHole["5L"]]),
	)
	addVias(layer5)

	layer7 := g.LayerN(7)
	layer7.Add(
		Polygon(Pt{0, 0}, true, layer7SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer7SpiralL, 0.0),
		padLine(startLayer7R, innerViaPts[innerHole["7R"]]),
		padLine(startLayer7L, innerViaPts[innerHole["7L"]]),
		padLine(endLayer7R, outerViaPts[outerHole["7R"]]),
		padLine(endLayer7L, outerViaPts[outerHole["7L"]]),
	)
	addVias(layer7)

	layer9 := g.LayerN(9)
	layer9.Add(
		Polygon(Pt{0, 0}, true, layer9SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer9SpiralL, 0.0),
		padLine(startLayer9R, innerViaPts[innerHole["9R"]]),
		padLine(startLayer9L, innerViaPts[innerHole["9L"]]),
		padLine(endLayer9R, outerViaPts[outerHole["9R"]]),
		padLine(endLayer9L, outerViaPts[outerHole["9L"]]),
	)
	addVias(layer9)

	layer11 := g.LayerN(11)
	layer11.Add(
		Polygon(Pt{0, 0}, true, layer11SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer11SpiralL, 0.0),
		padLine(startLayer11R, innerViaPts[innerHole["11R"]]),
		padLine(startLayer11L, innerViaPts[innerHole["11L"]]),
		padLine(endLayer11R, outerViaPts[outerHole["11R"]]),
		padLine(endLayer11L, outerViaPts[outerHole["11L"]]),
	)
	addVias(layer11)

	layer13 := g.LayerN(13)
	layer13.Add(
		Polygon(Pt{0, 0}, true, layer13SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer13SpiralL, 0.0),
		padLine(startLayer13R, innerViaPts[innerHole["13R"]]),
		padLine(startLayer13L, innerViaPts[innerHole["13L"]]),
		padLine(endLayer13R, outerViaPts[outerHole["13R"]]),
		padLine(endLayer13L, outerViaPts[outerHole["13L"]]),
	)
	addVias(layer13)

	layer15 := g.LayerN(15)
	layer15.Add(
		Polygon(Pt{0, 0}, true, layer15SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer15SpiralL, 0.0),
		padLine(startLayer15R, innerViaPts[innerHole["15R"]]),
		padLine(startLayer15L, innerViaPts[innerHole["15L"]]),
		padLine(endLayer15R, outerViaPts[outerHole["15R"]]),
		padLine(endLayer15L, outerViaPts[outerHole["15L"]]),
	)
	addVias(layer15)

	layer17 := g.LayerN(17)
	layer17.Add(
		Polygon(Pt{0, 0}, true, layer17SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer17SpiralL, 0.0),
		padLine(startLayer17R, innerViaPts[innerHole["17R"]]),
		padLine(startLayer17L, innerViaPts[innerHole["17L"]]),
		padLine(endLayer17R, outerViaPts[outerHole["17R"]]),
		padLine(endLayer17L, outerViaPts[outerHole["17L"]]),
	)
	addVias(layer17)

	layer19 := g.LayerN(19)
	layer19.Add(
		Polygon(Pt{0, 0}, true, layer19SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer19SpiralL, 0.0),
		padLine(startLayer19R, innerViaPts[innerHole["19R"]]),
		padLine(startLayer19L, innerViaPts[innerHole["19L"]]),
		padLine(endLayer19R, outerViaPts[outerHole["19R"]]),
		padLine(endLayer19L, outerViaPts[outerHole["19L"]]),
	)
	addVias(layer19)

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
		outerLabel := func(label string, num float64) *TextT {
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((0.3+num)*angleDelta)
			y := r * math.Sin((0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}
		outerLabel2 := func(label string, num float64) *TextT {
			r := outerR + 0.5**trace + *gap + 0.5*padD
			x := r * math.Cos((-0.3+num)*angleDelta)
			y := r * math.Sin((-0.3+num)*angleDelta)
			return Text(x, y, 1.0, label, *fontName, outerLabelSize, &Center)
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(0, 0.3*r, 1.0, message, *fontName, pts, &Center),
			innerLabel("TR"),
			innerLabel("TL"),
			innerLabel("BR"),
			innerLabel("BL"),
			innerLabel("2R"),
			innerLabel("2L"),
			innerLabel("3R"),
			innerLabel("3L"),
			innerLabel("4R"),
			innerLabel("4L"),
			innerLabel("5R"),
			innerLabel("5L"),
			innerLabel("6R"),
			innerLabel("6L"),
			innerLabel("7R"),
			innerLabel("7L"),
			innerLabel("8R"),
			innerLabel("8L"),
			innerLabel2("9R"),
			innerLabel2("9L"),
			innerLabel("10R"),
			innerLabel("10L"),
			innerLabel2("11R"),
			innerLabel2("11L"),
			innerLabel2("12R"),
			innerLabel2("12L"),
			innerLabel2("13R"),
			innerLabel2("13L"),
			innerLabel2("14R"),
			innerLabel2("14L"),
			innerLabel2("15R"),
			innerLabel2("15L"),
			innerLabel2("16R"),
			innerLabel2("16L"),
			innerLabel2("17R"),
			innerLabel2("17L"),
			innerLabel2("18R"),
			innerLabel2("18L"),
			innerLabel2("19R"),
			innerLabel2("19L"),

			outerLabel2("TR", 0),
			outerLabel("TL", 10),
			outerLabel2("BR", 10),
			outerLabel2("BL", 0.5),
			outerLabel("2R", 1),
			outerLabel("2L", 11),
			outerLabel("3R", 9),
			outerLabel("3L", 19),
			outerLabel2("4R", 19),
			outerLabel2("4L", 9),
			outerLabel2("5R", 11),
			outerLabel2("5L", 1),
			outerLabel("6R", 2),
			outerLabel("6L", 12),
			outerLabel("7R", 8),
			outerLabel("7L", 18),
			outerLabel2("8R", 18),
			outerLabel2("8L", 8),
			outerLabel2("9R", 12),
			outerLabel2("9L", 2),
			outerLabel("10R", 3),
			outerLabel("10L", 13),
			outerLabel("11R", 7),
			outerLabel("11L", 17),
			outerLabel2("12R", 17),
			outerLabel2("12L", 7),
			outerLabel2("13R", 13),
			outerLabel2("13L", 3),
			outerLabel("14R", 4),
			outerLabel("14L", 14),
			outerLabel("15R", 6),
			outerLabel("15L", 16),
			outerLabel2("16R", 16),
			outerLabel2("16L", 6),
			outerLabel2("17R", 14),
			outerLabel2("17L", 4),
			outerLabel("18R", 15),
			outerLabel("18L", 5),
			outerLabel2("19R", 15),
			outerLabel2("19L", 5),

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
