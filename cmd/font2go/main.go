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
	"strings"
	"text/template"
)

const (
	prefix = "font-"
)

var (
	outTemp = template.Must(template.New("out").Funcs(funcMap).Parse(goTemplate))
	funcMap = template.FuncMap{
		"floats":  floats,
		"orEmpty": orEmpty,
		"utf8":    utf8Escape,
	}
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		log.Printf("Processing file %q ...", arg)

		fontData := &FontData{}
		if buf, err := ioutil.ReadFile(arg); err != nil {
			log.Fatal(err)
		} else {
			if err := xml.Unmarshal(buf, fontData); err != nil {
				log.Fatal(err)
			}
		}

		fontData.Font.ID = strings.ToLower(fontData.Font.ID)

		for _, g := range fontData.Font.Glyphs {
			g.ParsePath()
		}

		var buf bytes.Buffer
		if err := outTemp.Execute(&buf, fontData.Font); err != nil {
			log.Fatal(err)
		}

		filename := fmt.Sprintf("%v%v.go", prefix, fontData.Font.ID)
		fmtBuf, err := format.Source(buf.Bytes())
		if err != nil {
			ioutil.WriteFile(filename, buf.Bytes(), 0644) // Dump the unformatted output.
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(filename, fmtBuf, 0644); err != nil {
			log.Fatal(err)
		}
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

func init() {
  Fonts["{{ .ID }}"] = {{ .ID }}Font
}

var {{ .ID }}Font = &Font{
	// ID: "{{ .ID }}",
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
				{ C: '{{ .Command }}'{{ if .Parameters }}, P: {{ .Parameters | floats }}{{ end }} },{{ end }}
			},
		},{{ end }}{{ end }}
	},
}
`
