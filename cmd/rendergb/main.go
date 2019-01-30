// rendergb renders Gerber files to a PNG image in Go.
package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

var (
	width  = flag.Int("width", 800, "Image width")
	height = flag.Int("height", 800, "Image height")
	out    = flag.String("out", "out.png", "Output image filename")

	xyRE = regexp.MustCompile(`^X(-?\d+)Y(-?\d+)D(\d+)\*\s*$`)
)

func main() {
	flag.Parse()

	// First, get the minimum bounding box of all Gerber files.
	var mbb *mbbT
	for _, arg := range flag.Args() {
		log.Printf("Processing %v ...", arg)
		mbb = updateMBB(arg, mbb)
	}
	if mbb == nil {
		log.Fatal("No polygons found.")
	}

	xs := float64(*width) / (mbb.xmax - mbb.xmin)
	ys := float64(*height) / (mbb.ymax - mbb.ymin)
	mbb.scale = xs
	log.Printf("xs=%v, ys=%v", xs, ys)
	if ys < mbb.scale {
		mbb.scale = ys
		*width = int(0.5 + mbb.scale*(mbb.xmax-mbb.xmin))
	} else {
		*height = int(0.5 + mbb.scale*(mbb.ymax-mbb.ymin))
	}
	log.Printf("MBB=%#v, scale=%v", *mbb, mbb.scale)

	dc := gg.NewContext(*width, *height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	currentColor := func() {
		dc.SetRGB(1, 1, 1)
	}
	for _, arg := range flag.Args() {
		log.Printf("Processing %v ...", arg)
		switch strings.ToLower(arg[len(arg)-4:]) {
		case ".gto":
			currentColor = func() {
				dc.SetRGB(1, 1, 1)
			}
		default:
			currentColor = func() {
				dc.SetRGB(1, 0, 0)
			}
		}
		currentColor()
		renderLayer(arg, dc, mbb, currentColor)
	}
	dc.SavePNG(*out)
}

type mbbT struct {
	xmin, xmax float64
	ymin, ymax float64
	scale      float64
}

func renderLayer(filename string, dc *gg.Context, mbb *mbbT, currentColor func()) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		s, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if strings.HasPrefix(s, "%LPD") {
			currentColor()
		} else if strings.HasPrefix(s, "%LPC") {
			dc.SetRGB(0, 0, 0)
		}

		m := xyRE.FindStringSubmatch(s)
		if len(m) == 4 {
			x, y, d := atof(m[1]), atof(m[2]), atoi(m[3])
			if d == 2 {
				dc.Fill()
				dc.MoveTo(mbb.scale*(x-mbb.xmin), float64(*height)-mbb.scale*(y-mbb.ymin))
			} else {
				dc.LineTo(mbb.scale*(x-mbb.xmin), float64(*height)-mbb.scale*(y-mbb.ymin))
			}
		}

		if err == io.EOF {
			break
		}
	}
}

func updateMBB(filename string, mbb *mbbT) *mbbT {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		s, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if len(s) > 0 && s[0] == 'X' {
			mbb = updateMBBLine(s, mbb)
		}
		if err == io.EOF {
			break
		}
	}
	return mbb
}

func updateMBBLine(line string, mbb *mbbT) *mbbT {
	m := xyRE.FindStringSubmatch(line)
	if len(m) != 4 {
		return mbb
	}
	x, y := atof(m[1]), atof(m[2])
	if mbb == nil {
		return &mbbT{
			xmin: x,
			xmax: x,
			ymin: y,
			ymax: y,
		}
	}
	if x < mbb.xmin {
		mbb.xmin = x
	}
	if y < mbb.ymin {
		mbb.ymin = y
	}
	if x > mbb.xmax {
		mbb.xmax = x
	}
	if y > mbb.ymax {
		mbb.ymax = y
	}
	return mbb
}

func atoi(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("unable to parse %q as int64", s)
	}
	return v
}

func atof(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("unable to parse %q as float64", s)
	}
	return v
}
