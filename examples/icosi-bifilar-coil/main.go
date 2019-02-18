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
	innerHole6R := innerViaPts[19]
	innerHole6L := innerViaPts[9]
	innerHole7R := innerViaPts[11]
	innerHole7L := innerViaPts[1]
	innerHole8R := innerViaPts[15]
	innerHole8L := innerViaPts[5]
	innerHole9R := innerViaPts[15]
	innerHole9L := innerViaPts[5]
	innerHole10R := innerViaPts[0]
	innerHole10L := innerViaPts[10]
	innerHole11R := innerViaPts[10]
	innerHole11L := innerViaPts[0]
	innerHole12R := innerViaPts[14]
	innerHole12L := innerViaPts[4]
	innerHole13R := innerViaPts[16]
	innerHole13L := innerViaPts[6]
	innerHole14R := innerViaPts[1]
	innerHole14L := innerViaPts[11]
	innerHole15R := innerViaPts[9]
	innerHole15L := innerViaPts[19]
	innerHole16R := innerViaPts[13]
	innerHole16L := innerViaPts[3]
	innerHole17R := innerViaPts[17]
	innerHole17L := innerViaPts[7]
	innerHole18R := innerViaPts[12]
	innerHole18L := innerViaPts[2]
	innerHole19R := innerViaPts[18]
	innerHole19L := innerViaPts[8]

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
	outerHoleTR := outerViaPts[0]
	outerHoleTL := outerViaPts[10]
	outerHoleBR := outerViaPts[10]
	outerHoleBL := outerViaPts[20]
	outerHole2R := outerViaPts[1]
	outerHole2L := outerViaPts[11]
	outerHole3R := outerViaPts[9]
	outerHole3L := outerViaPts[19]
	outerHole4R := outerViaPts[19]
	outerHole4L := outerViaPts[9]
	outerHole5R := outerViaPts[11]
	outerHole5L := outerViaPts[1]
	outerHole6R := outerViaPts[2]
	outerHole6L := outerViaPts[12]
	outerHole7R := outerViaPts[8]
	outerHole7L := outerViaPts[18]
	outerHole8R := outerViaPts[18]
	outerHole8L := outerViaPts[8]
	outerHole9R := outerViaPts[12]
	outerHole9L := outerViaPts[2]
	outerHole10R := outerViaPts[3]
	outerHole10L := outerViaPts[13]
	outerHole11R := outerViaPts[7]
	outerHole11L := outerViaPts[17]
	outerHole12R := outerViaPts[17]
	outerHole12L := outerViaPts[7]
	outerHole13R := outerViaPts[13]
	outerHole13L := outerViaPts[3]
	outerHole14R := outerViaPts[4]
	outerHole14L := outerViaPts[14]
	outerHole15R := outerViaPts[6]
	outerHole15L := outerViaPts[16]
	outerHole16R := outerViaPts[16]
	outerHole16L := outerViaPts[6]
	outerHole17R := outerViaPts[14]
	outerHole17L := outerViaPts[4]
	outerHole18R := outerViaPts[15]
	outerHole18L := outerViaPts[5]
	outerHole19R := outerViaPts[15]
	outerHole19L := outerViaPts[5]

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
		padLine(startTopR, innerHoleTR),
		padLine(startTopL, innerHoleTL),
		padLine(endTopR, outerHoleTR),
		padLine(endTopL, outerHoleTL),
	)
	addVias(top)

	topMask := g.TopSolderMask()
	addVias(topMask)

	layer2 := g.LayerN(2)
	layer2.Add(
		Polygon(Pt{0, 0}, true, layer2SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer2SpiralL, 0.0),
		padLine(startLayer2R, innerHole2R),
		padLine(startLayer2L, innerHole2L),
		padLine(endLayer2R, outerHole2R),
		padLine(endLayer2L, outerHole2L),
	)
	addVias(layer2)

	layer4 := g.LayerN(4)
	layer4.Add(
		Polygon(Pt{0, 0}, true, layer4SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer4SpiralL, 0.0),
		padLine(startLayer4R, innerHole4R),
		padLine(startLayer4L, innerHole4L),
		padLine(endLayer4R, outerHole4R),
		padLine(endLayer4L, outerHole4L),
	)
	addVias(layer4)

	layer6 := g.LayerN(6)
	layer6.Add(
		Polygon(Pt{0, 0}, true, layer6SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer6SpiralL, 0.0),
		padLine(startLayer6R, innerHole6R),
		padLine(startLayer6L, innerHole6L),
		padLine(endLayer6R, outerHole6R),
		padLine(endLayer6L, outerHole6L),
	)
	addVias(layer6)

	layer8 := g.LayerN(8)
	layer8.Add(
		Polygon(Pt{0, 0}, true, layer8SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer8SpiralL, 0.0),
		padLine(startLayer8R, innerHole8R),
		padLine(startLayer8L, innerHole8L),
		padLine(endLayer8R, outerHole8R),
		padLine(endLayer8L, outerHole8L),
	)
	addVias(layer8)

	layer10 := g.LayerN(10)
	layer10.Add(
		Polygon(Pt{0, 0}, true, layer10SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer10SpiralL, 0.0),
		padLine(startLayer10R, innerHole10R),
		padLine(startLayer10L, innerHole10L),
		padLine(endLayer10R, outerHole10R),
		padLine(endLayer10L, outerHole10L),
	)
	addVias(layer10)

	layer12 := g.LayerN(12)
	layer12.Add(
		Polygon(Pt{0, 0}, true, layer12SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer12SpiralL, 0.0),
		padLine(startLayer12R, innerHole12R),
		padLine(startLayer12L, innerHole12L),
		padLine(endLayer12R, outerHole12R),
		padLine(endLayer12L, outerHole12L),
	)
	addVias(layer12)

	layer14 := g.LayerN(14)
	layer14.Add(
		Polygon(Pt{0, 0}, true, layer14SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer14SpiralL, 0.0),
		padLine(startLayer14R, innerHole14R),
		padLine(startLayer14L, innerHole14L),
		padLine(endLayer14R, outerHole14R),
		padLine(endLayer14L, outerHole14L),
	)
	addVias(layer14)

	layer16 := g.LayerN(16)
	layer16.Add(
		Polygon(Pt{0, 0}, true, layer16SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer16SpiralL, 0.0),
		padLine(startLayer16R, innerHole16R),
		padLine(startLayer16L, innerHole16L),
		padLine(endLayer16R, outerHole16R),
		padLine(endLayer16L, outerHole16L),
	)
	addVias(layer16)

	layer18 := g.LayerN(18)
	layer18.Add(
		Polygon(Pt{0, 0}, true, layer18SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer18SpiralL, 0.0),
		padLine(startLayer18R, innerHole18R),
		padLine(startLayer18L, innerHole18L),
		padLine(endLayer18R, outerHole18R),
		padLine(endLayer18L, outerHole18L),
	)
	addVias(layer18)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, true, botSpiralR, 0.0),
		Polygon(Pt{0, 0}, true, botSpiralL, 0.0),
		padLine(startBotR, innerHoleBR),
		padLine(startBotL, innerHoleBL),
		padLine(endBotR, outerHoleBR),
		padLine(endBotL, outerHoleBL),
	)
	addVias(bottom)

	bottomMask := g.BottomSolderMask()
	addVias(bottomMask)

	layer3 := g.LayerN(3)
	layer3.Add(
		Polygon(Pt{0, 0}, true, layer3SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer3SpiralL, 0.0),
		padLine(startLayer3R, innerHole3R),
		padLine(startLayer3L, innerHole3L),
		padLine(endLayer3R, outerHole3R),
		padLine(endLayer3L, outerHole3L),
	)
	addVias(layer3)

	layer5 := g.LayerN(5)
	layer5.Add(
		Polygon(Pt{0, 0}, true, layer5SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer5SpiralL, 0.0),
		padLine(startLayer5R, innerHole5R),
		padLine(startLayer5L, innerHole5L),
		padLine(endLayer5R, outerHole5R),
		padLine(endLayer5L, outerHole5L),
	)
	addVias(layer5)

	layer7 := g.LayerN(7)
	layer7.Add(
		Polygon(Pt{0, 0}, true, layer7SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer7SpiralL, 0.0),
		padLine(startLayer7R, innerHole7R),
		padLine(startLayer7L, innerHole7L),
		padLine(endLayer7R, outerHole7R),
		padLine(endLayer7L, outerHole7L),
	)
	addVias(layer7)

	layer9 := g.LayerN(9)
	layer9.Add(
		Polygon(Pt{0, 0}, true, layer9SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer9SpiralL, 0.0),
		padLine(startLayer9R, innerHole9R),
		padLine(startLayer9L, innerHole9L),
		padLine(endLayer9R, outerHole9R),
		padLine(endLayer9L, outerHole9L),
	)
	addVias(layer9)

	layer11 := g.LayerN(11)
	layer11.Add(
		Polygon(Pt{0, 0}, true, layer11SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer11SpiralL, 0.0),
		padLine(startLayer11R, innerHole11R),
		padLine(startLayer11L, innerHole11L),
		padLine(endLayer11R, outerHole11R),
		padLine(endLayer11L, outerHole11L),
	)
	addVias(layer11)

	layer13 := g.LayerN(13)
	layer13.Add(
		Polygon(Pt{0, 0}, true, layer13SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer13SpiralL, 0.0),
		padLine(startLayer13R, innerHole13R),
		padLine(startLayer13L, innerHole13L),
		padLine(endLayer13R, outerHole13R),
		padLine(endLayer13L, outerHole13L),
	)
	addVias(layer13)

	layer15 := g.LayerN(15)
	layer15.Add(
		Polygon(Pt{0, 0}, true, layer15SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer15SpiralL, 0.0),
		padLine(startLayer15R, innerHole15R),
		padLine(startLayer15L, innerHole15L),
		padLine(endLayer15R, outerHole15R),
		padLine(endLayer15L, outerHole15L),
	)
	addVias(layer15)

	layer17 := g.LayerN(17)
	layer17.Add(
		Polygon(Pt{0, 0}, true, layer17SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer17SpiralL, 0.0),
		padLine(startLayer17R, innerHole17R),
		padLine(startLayer17L, innerHole17L),
		padLine(endLayer17R, outerHole17R),
		padLine(endLayer17L, outerHole17L),
	)
	addVias(layer17)

	layer19 := g.LayerN(19)
	layer19.Add(
		Polygon(Pt{0, 0}, true, layer19SpiralR, 0.0),
		Polygon(Pt{0, 0}, true, layer19SpiralL, 0.0),
		padLine(startLayer19R, innerHole19R),
		padLine(startLayer19L, innerHole19L),
		padLine(endLayer19R, outerHole19R),
		padLine(endLayer19L, outerHole19L),
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

		innerLabel := func(label string, num float64) *TextT {
			r := innerR - viaPadD
			x := r * math.Cos(num*angleDelta)
			y := r * math.Sin(num*angleDelta)
			return Text(x, y, 1.0, label, *fontName, labelSize, &Center)
		}
		innerLabel2 := func(label string, num float64) *TextT {
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
			innerLabel("6R", 19),
			innerLabel("6L", 9),
			innerLabel("7R", 11),
			innerLabel("7L", 1),
			innerLabel("8R", 15),
			innerLabel("8L", 5),
			innerLabel2("9R", 15),
			innerLabel2("9L", 5),
			innerLabel("10R", 0),
			innerLabel("10L", 10),
			innerLabel2("11R", 10),
			innerLabel2("11L", 0),
			innerLabel2("12R", 14),
			innerLabel2("12L", 4),
			innerLabel2("13R", 16),
			innerLabel2("13L", 6),
			innerLabel2("14R", 1),
			innerLabel2("14L", 11),
			innerLabel2("15R", 9),
			innerLabel2("15L", 19),
			innerLabel2("16R", 13),
			innerLabel2("16L", 3),
			innerLabel2("17R", 17),
			innerLabel2("17L", 7),
			innerLabel2("18R", 12),
			innerLabel2("18L", 2),
			innerLabel2("19R", 18),
			innerLabel2("19L", 8),

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
