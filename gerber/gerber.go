// Package gerber writes Gerber RS274X files (for PCBs).
package gerber

// Gerber represents the layers needed to build a PCB.
type Gerber struct {
	// FilenamePrefix is the filename prefix for the Gerber design files.
	FilenamePrefix string
	// Layers represents the layers making up the Gerber design.
	Layers []*Layer
}

// New returns a new Gerber design.
// filenamePrefix is the base filename for all gerber files (e.g. "bifilar-coil").
func New(filenamePrefix string) *Gerber {
	return &Gerber{
		FilenamePrefix: filenamePrefix,
	}
}
