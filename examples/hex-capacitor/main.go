// hex-capacitor creates Gerber files (and a bundled ZIP) representing
// six-layer capactors for manufacture on a printed circuit
// board (PCB).
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
	. "github.com/gmlewis/go-gerber/gerber"
	"github.com/gmlewis/go-gerber/gerber/viewer"
)

var (
	width      = flag.Float64("width", 100.0, "Width of PCB")
	height     = flag.Float64("height", 100.0, "Height of PCB")
	gap        = flag.Float64("gap", 0.15, "Gap between traces in mm (6mil = 0.15mm)")
	trace      = flag.Float64("trace", 0.15, "Width of traces in mm")
	prefix     = flag.String("prefix", "hex-capacitor", "Filename prefix for all Gerber files and zip")
	fontName   = flag.String("font", "freeserif", "Name of font to use for writing source on PCB (empty to not write)")
	view       = flag.Bool("view", false, "View the resulting design using Fyne")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

const (
	messageFmt = `This is a hex (4-layer)
capacitor. Size=%0.2fmm x %0.2fmm.
Trace size = %0.2fmm.
Gap size = %0.2fmm.`
	padD = 4
	padR = padD / 2
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

	railL, railR := genRails()

	pad1 := Point(padR, padR)
	pad2 := Point(*width-padR, *height-padR)

	contactDrill := func(pt Pt) *CircleT {
		const drillD = 0.5 * padD
		return Circle(pt, drillD)
	}

	drill := g.Drill()
	drill.Add(
		contactDrill(pad1),
		contactDrill(pad2),
	)

	contactPad := func(pt Pt) *CircleT {
		return Circle(pt, padD)
	}
	padLine := func(pt1, pt2 Pt) *LineT {
		return Line(pt1[0], pt1[1], pt2[0], pt2[1], CircleShape, padD)
	}

	leftLines := genLeftLines()
	rightLines := genRightLines()

	top := g.TopCopper()
	top.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{padR, 2 * padD}, pad1),
		padLine(Pt{*width - padR, *height - 2*padD}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	top.Add(leftLines...)
	top.Add(rightLines...)

	topMask := g.TopSolderMask()
	topMask.Add(
		contactPad(pad1),
		contactPad(pad2),
	)

	layer2 := g.LayerN(2)
	layer2.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{2 * padD, padR}, pad1),
		padLine(Pt{*width - 2*padD, *height - padR}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	layer2.Add(leftLines...)
	layer2.Add(rightLines...)

	layer4 := g.LayerN(4)
	layer4.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{2 * padD, padR}, pad1),
		padLine(Pt{*width - 2*padD, *height - padR}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	layer4.Add(leftLines...)
	layer4.Add(rightLines...)

	bottom := g.BottomCopper()
	bottom.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{2 * padD, padR}, pad1),
		padLine(Pt{*width - 2*padD, *height - padR}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	bottom.Add(leftLines...)
	bottom.Add(rightLines...)

	bottomMask := g.BottomSolderMask()
	bottomMask.Add(
		contactPad(pad1),
		contactPad(pad2),
	)

	layer3 := g.LayerN(3)
	layer3.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{padR, 2 * padD}, pad1),
		padLine(Pt{*width - padR, *height - 2*padD}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	layer3.Add(leftLines...)
	layer3.Add(rightLines...)

	layer5 := g.LayerN(5)
	layer5.Add(
		Polygon(Pt{0, 0}, false, railL, 0.0),
		Polygon(Pt{0, 0}, false, railR, 0.0),

		padLine(Pt{padR, 2 * padD}, pad1),
		padLine(Pt{*width - padR, *height - 2*padD}, pad2),
		contactPad(pad1),
		contactPad(pad2),
	)
	layer5.Add(leftLines...)
	layer5.Add(rightLines...)

	outline := g.Outline()
	border := []Pt{{0, 0}, {*width, 0}, {*width, *height}, {0, *height}}
	outline.Add(
		Polygon(Pt{0, 0}, true, border, 0.0),
	)

	// if *fontName != "" {
	// 	tss := g.TopSilkscreen()
	// 	tsAdd(
	// 	)
	// }

	if err := g.WriteGerber(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")

	if *view {
		viewer.Gerber(g, true)
	}
}

func genLeftLines() (lines []Primitive) {
	dx := *trace + *gap
	topY := *height - padD - *gap - 0.5**trace

	for x := padD + *gap + 0.5**trace; x < *width-padD-*gap; x += 2 * dx {
		lines = append(lines,
			Line(x, padR, x, topY, CircleShape, *trace),
		)
	}

	return lines
}

func genRightLines() (lines []Primitive) {
	dx := *trace + *gap
	botY := padD + *gap + 0.5**trace

	for x := padD + 2**gap + 1.5**trace; x < *width-padD-*gap; x += 2 * dx {
		lines = append(lines,
			Line(x, botY, x, *height-padR, CircleShape, *trace),
		)
	}

	return lines
}

func genRails() (railL, railR []Pt) {

	// Create the power bus lines.
	railL = append(railL,
		Pt{0, padD + *gap},
		Pt{0, *height},
		Pt{*width - padD - *gap, *height},
		Pt{*width - padD - *gap, *height - padD},
		Pt{padD, *height - padD},
		Pt{padD, padD + *gap},
		Pt{0, padD + *gap},
	)

	railR = append(railR,
		Pt{*width, *height - padD - *gap},
		Pt{*width, 0},
		Pt{padD + *gap, 0},
		Pt{padD + *gap, padD},
		Pt{*width - padD, padD},
		Pt{*width - padD, *height - padD - *gap},
		Pt{*width, *height - padD - *gap},
	)

	return railL, railR
}
