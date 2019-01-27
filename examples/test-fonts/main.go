// test-fonts creates Gerber files (and a bundled ZIP) with the
// given message in all fonts in order to test out the font rendering.
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	. "github.com/gmlewis/go-gerber/gerber"
)

var (
	all = flag.Bool("all", false, "All renders all glyphs and overrides -msg")
	msg = flag.String("msg", "M", "Message to write to Gerber file silkscreen")
	pts = flag.Float64("pts", 72.0, "Point size to render font")
	x   = flag.Float64("x", 0, "X starting position of font")
	y   = flag.Float64("y", 0, "Y starting position of font")
)

const (
	allLineLength = 20
)

func main() {
	flag.Parse()

	for name, font := range Fonts {
		g := New(name)

		message := *msg
		if *all {
			var glyphs []string
			for g := range font.Glyphs {
				glyphs = append(glyphs, g)
			}
			sort.Strings(glyphs)
			var lines []string
			for len(glyphs) > 0 {
				end := allLineLength
				if end > len(glyphs) {
					end = len(glyphs)
				}
				lines = append(lines, fmt.Sprintf("%v", strings.Join(glyphs[0:end], "")))
				glyphs = glyphs[end:]
			}
			message = strings.Join(lines, "\n")
		}

		tss := g.TopSilkscreen()
		tss.Add(
			Text(*x, *y, 1.0, message, name, *pts),
		)

		if err := g.WriteGerber(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done.")
}
