// Package gerber writes Gerber RS274X files (for PCBs).
package gerber

import (
	"archive/zip"
	"os"
	"sync"
)

// Gerber represents the layers needed to build a PCB.
type Gerber struct {
	// FilenamePrefix is the filename prefix for the Gerber design files.
	FilenamePrefix string
	// Layers represents the layers making up the Gerber design.
	Layers []*Layer

	mu  sync.Mutex // protects mbb against multiple requests
	mbb *MBB       // cached minimum bounding box
}

// New returns a new Gerber design.
// filenamePrefix is the base filename for all gerber files (e.g. "bifilar-coil").
func New(filenamePrefix string) *Gerber {
	return &Gerber{
		FilenamePrefix: filenamePrefix,
	}
}

// WriteGerber writes all the Gerber layers to their respective files
// then zips them all together into a ZIP file with the same prefix
// for sending to PCB manufacturers.
func (g *Gerber) WriteGerber() error {
	zf, err := os.Create(g.FilenamePrefix + ".zip")
	if err != nil {
		return err
	}
	zw := zip.NewWriter(zf)
	for _, layer := range g.Layers {
		f, err := zw.Create(layer.Filename)
		if err != nil {
			return err
		}
		if err := layer.WriteGerber(f); err != nil {
			return err
		}
		w, err := os.Create(layer.Filename)
		if err != nil {
			return err
		}
		if err := layer.WriteGerber(w); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	}
	return zw.Close()
}

// MBB returns the minimum bounding box of the design in millimeters.
func (g *Gerber) MBB() MBB {
	g.mu.Lock()
	defer g.mu.Unlock() // Only calculate MBB once.
	if g.mbb != nil {
		return *g.mbb
	}

	var finalMu sync.Mutex
	var wg sync.WaitGroup
	for _, p := range g.Layers {
		wg.Add(1)
		go func(p *Layer) {
			v := p.MBB()
			finalMu.Lock()
			if g.mbb == nil {
				g.mbb = &v
			} else {
				g.mbb.Join(&v)
			}
			finalMu.Unlock()
			wg.Done()
		}(p)
	}
	wg.Wait()
	return *g.mbb
}
