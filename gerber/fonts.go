package gerber

// Font represents a webfont.
type Font struct {
	ID               string
	HorizAdvX        float64
	UnitsPerEm       float64
	Ascent           float64
	Descent          float64
	MissingHorizAdvX float64
	Glyphs           map[string]*Glyph
}

// Glyph represents an individual character of the webfont data.
type Glyph struct {
	HorizAdvX float64
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
	C byte      // C is the command.
	P []float64 // P are the parameters of the command.
}

// Fonts is a map of all the available fonts.
//
// The map is initialized at runtime by `init` functions
// in order to reduce the overall initial compile time
// of the package.
var Fonts = map[string]*Font{}
