// test-fonts creates Gerber files (and a bundled ZIP) with the
// given message in all fonts in order to test out the font rendering.
package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/gmlewis/go-gerber/gerber"
)

var (
	msg = flag.String("msg", "M", "Message to write to Gerber file silkscreen")
	pts = flag.Float64("pts", 72.0, "Point size to render font")
	x   = flag.Float64("x", 0, "X starting position of font")
	y   = flag.Float64("y", 0, "Y starting position of font")
)

func main() {
	flag.Parse()

	for name := range Fonts {
		g := New(name)
		tss := g.TopSilkscreen()
		tss.Add(
			Text(*x, *y, 1.0, *msg, name, *pts),
		)

		if err := g.WriteGerber(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Done.")
}
