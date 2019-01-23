package gerber

import (
	"io"
)

// Layer represents a printed circuit board layer.
type Layer struct {
	// Filename is the filename of the Gerber layer.
	Filename string
	// Primitives represents the collection of primitives.
	Primitives []Primitive

	g *Gerber // Root Gerber object.
}

// Add adds primitives to a layer.
func (l *Layer) Add(primitives ...Primitive) {
	l.Primitives = append(l.Primitives, primitives...)
}

// WriteGerber writes a layer to its corresponding Gerber layer file.
func (l *Layer) WriteGerber(w io.Writer) error {
	io.WriteString(w, "%FSLAX36Y36*%\n")
	io.WriteString(w, "%MOMM*%\n")
	io.WriteString(w, "%LPD*%\n")

	for _, p := range l.Primitives {
		p.WriteGerber(w)
	}
	return nil
}

// TopCopper adds a top copper layer to the design
// and returns the layer.
func (g *Gerber) TopCopper() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".gtl",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}

// TopSolderMask adds a top solder mask layer to the design
// and returns the layer.
func (g *Gerber) TopSolderMask() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".gts",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}

// BottomCopper adds a bottom copper layer to the design
// and returns the layer.
func (g *Gerber) BottomCopper() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".gbl",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}

// BottomSolderMask adds a bottom solder mask layer to the design
// and returns the layer.
func (g *Gerber) BottomSolderMask() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".gbs",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}

// Drill adds a drill layer to the design
// and returns the layer.
func (g *Gerber) Drill() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".xln",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}

// Outline adds an outline layer to the design
// and returns the layer.
func (g *Gerber) Outline() *Layer {
	layer := &Layer{
		Filename: g.FilenamePrefix + ".gko",
	}
	g.Layers = append(g.Layers, layer)
	return layer
}
