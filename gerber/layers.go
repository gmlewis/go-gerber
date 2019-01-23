package gerber

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

// Write writes a layer to its corresponding layer file.
func (l *Layer) Write() error {
	return nil
}

// TopCopper adds a top copper layer to the design
// and returns the layer.
func (g *Gerber) TopCopper() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".gtl",
	}
}

// TopSolderMask adds a top solder mask layer to the design
// and returns the layer.
func (g *Gerber) TopSolderMask() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".gts",
	}
}

// BottomCopper adds a bottom copper layer to the design
// and returns the layer.
func (g *Gerber) BottomCopper() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".gbl",
	}
}

// BottomSolderMask adds a bottom solder mask layer to the design
// and returns the layer.
func (g *Gerber) BottomSolderMask() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".gbs",
	}
}

// Drill adds a drill layer to the design
// and returns the layer.
func (g *Gerber) Drill() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".xln",
	}
}

// Outline adds an outline layer to the design
// and returns the layer.
func (g *Gerber) Outline() *Layer {
	return &Layer{
		Filename: g.FilenamePrefix + ".gko",
	}
}
