// font2go reads one or more standard SVG webfont file(s) and writes Go file(s)
// used to render them to a Gerber layer.
package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"sort"
	"strings"
	"text/template"
)

var (
	filename = flag.String("out", "fonts.go", "Output filename for Go fonts file")

	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(goTemplate))
	funcMap = template.FuncMap{
		"floats":  floats,
		"orEmpty": orEmpty,
		"utf8":    utf8Escape,
	}
)

func main() {
	flag.Parse()

	var fonts []*Font
	for _, arg := range flag.Args() {
		log.Printf("Processing file %q ...", arg)

		buf, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		fontData := &FontData{}
		if err := xml.Unmarshal(buf, fontData); err != nil {
			log.Fatal(err)
		}

		fontData.Font.ID = strings.ToLower(fontData.Font.ID)

		for _, g := range fontData.Font.Glyphs {
			g.ParsePath()
		}

		fonts = append(fonts, fontData.Font)
	}

	sort.Slice(fonts, func(a, b int) bool { return fonts[a].ID < fonts[b].ID })

	var buf bytes.Buffer
	if err := outTemp.Execute(&buf, fonts); err != nil {
		log.Fatal(err)
	}

	fmtBuf, err := format.Source(buf.Bytes())
	if err != nil {
		ioutil.WriteFile(*filename, buf.Bytes(), 0644) // Dump the unformatted output.
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(*filename, fmtBuf, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

func utf8Escape(s *string) string {
	if s == nil || *s == "" {
		return `""`
	}
	return fmt.Sprintf("%+q", *s)
}

func orEmpty(s *string) string {
	if s == nil || *s == "" {
		return `""`
	}
	return fmt.Sprintf("%q", *s)
}

func floats(f []float64) string {
	return fmt.Sprintf("%#v", f)
}

var goTemplate = `// Auto-generated - DO NOT EDIT!

package gerber

// Font represents a webfont.
type Font struct {
	ID           string
	HorizAdvX    int
	UnitsPerEm   int
	Ascent       int
	Descent      int
	MissingHorizAdvX int
	Glyphs       map[string]*Glyph
}

// Glyph represents an individual character of the webfont data.
type Glyph struct {
	HorizAdvX int
	Unicode   string
	GerberLP  string
	PathSteps []*PathStep
}

// PathStep represents a single path step.
//
// There are 20 possible commands, broken up into 6 types,
// with each command having an "absolute" (upper case) and
// a "relative" (lower case) version.
//
// MoveTo: M, m
// LineTo: L, l, H, h, V, v
// Cubic Bézier Curve: C, c, S, s
// Quadratic Bézier Curve: Q, q, T, t
// Elliptical Arc Curve: A, a
// ClosePath: Z, z
type PathStep struct {
	C    string // C is the command.
	P []float64 // P are the parameters of the command.
}

var fonts = map[string]*Font{ {{ range . }}
	"{{ .ID }}": {
		ID: "{{ .ID }}",
		HorizAdvX:  {{ .HorizAdvX }},
		UnitsPerEm: {{ .FontFace.UnitsPerEm }},
		Ascent:     {{ .FontFace.Ascent }},
		Descent:    {{ .FontFace.Descent }},
		MissingHorizAdvX: {{ .MissingGlyph.HorizAdvX }},
		Glyphs: map[string]*Glyph{ {{ range .Glyphs }}{{ if .Unicode }}
			{{ .Unicode | utf8 }}: {
				HorizAdvX: {{ .HorizAdvX }},
				Unicode: {{ .Unicode | utf8 }},
				GerberLP: {{ .GerberLP | orEmpty }},
				PathSteps: []*PathStep{ {{ range .PathSteps }}
					{ C: "{{ .Command }}"{{ if .Parameters }}, P: {{ .Parameters | floats }}{{ end }} },{{ end }}
				},
			},{{ end }}{{ end }}
		},
	},{{ end }}
}
`
