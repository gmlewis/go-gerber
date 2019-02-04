package gerber

import (
	"math"
	"testing"

	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
)

func TestText(t *testing.T) {
	const (
		message    = "012"
		fontName   = "freeserif"
		wantWidth  = 36.8554
		wantHeight = 17.526
		pts        = 72
		eps        = 1e-6
		sf         = 1e6
	)

	tests := []struct {
		name     string
		x, y     float64
		opts     *TextOpts
		wantXmin float64
		wantYmin float64
		wantXmax float64
		wantYmax float64
	}{
		{
			name:     "XLeft,YBottom",
			wantXmin: 0,
			wantYmin: 0,
			wantXmax: wantWidth * sf,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XLeft,YBottom w/ offset",
			x:        10,
			y:        20,
			wantXmin: 10 * sf,
			wantYmin: 20 * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XCenter,YBottom",
			opts:     &TextOpts{XAlign: XCenter, YAlign: YBottom},
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: 0,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XCenter,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XCenter, YAlign: YBottom},
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: 20 * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XRight,YBottom",
			opts:     &TextOpts{XAlign: XRight, YAlign: YBottom},
			wantXmin: -wantWidth * sf,
			wantYmin: 0,
			wantXmax: 0,
			wantYmax: wantHeight * sf,
		},
		{
			name:     "XRight,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XRight, YAlign: YBottom},
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: 20 * sf,
			wantXmax: 10 * sf,
			wantYmax: (20 + wantHeight) * sf,
		},
		{
			name:     "XLeft,YCenter",
			opts:     &TextOpts{XAlign: XLeft, YAlign: YCenter},
			wantXmin: 0,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: wantWidth * sf,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XLeft,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XLeft, YAlign: YCenter},
			wantXmin: 10 * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XCenter,YCenter",
			opts:     &TextOpts{XAlign: XCenter, YAlign: YCenter},
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XCenter,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XCenter, YAlign: YCenter},
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XRight,YCenter",
			opts:     &TextOpts{XAlign: XRight, YAlign: YCenter},
			wantXmin: -wantWidth * sf,
			wantYmin: -0.5 * wantHeight * sf,
			wantXmax: 0,
			wantYmax: 0.5 * wantHeight * sf,
		},
		{
			name:     "XRight,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XRight, YAlign: YCenter},
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: (20 - 0.5*wantHeight) * sf,
			wantXmax: 10 * sf,
			wantYmax: (20 + 0.5*wantHeight) * sf,
		},
		{
			name:     "XLeft,YTop",
			opts:     &TextOpts{XAlign: XLeft, YAlign: YTop},
			wantXmin: 0,
			wantYmin: -wantHeight * sf,
			wantXmax: wantWidth * sf,
			wantYmax: 0,
		},
		{
			name:     "XLeft,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XLeft, YAlign: YTop},
			wantXmin: 10 * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: (10 + wantWidth) * sf,
			wantYmax: 20 * sf,
		},
		{
			name:     "XCenter,YTop",
			opts:     &TextOpts{XAlign: XCenter, YAlign: YTop},
			wantXmin: -0.5 * wantWidth * sf,
			wantYmin: -wantHeight * sf,
			wantXmax: 0.5 * wantWidth * sf,
			wantYmax: 0,
		},
		{
			name:     "XCenter,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XCenter, YAlign: YTop},
			wantXmin: (10 - 0.5*wantWidth) * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: (10 + 0.5*wantWidth) * sf,
			wantYmax: 20 * sf,
		},
		{
			name:     "XRight,YTop",
			opts:     &TextOpts{XAlign: XRight, YAlign: YTop},
			wantXmin: -wantWidth * sf,
			wantYmin: -wantHeight * sf,
			wantXmax: 0,
			wantYmax: 0,
		},
		{
			name:     "XRight,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     &TextOpts{XAlign: XRight, YAlign: YTop},
			wantXmin: (10 - wantWidth) * sf,
			wantYmin: (20 - wantHeight) * sf,
			wantXmax: 10 * sf,
			wantYmax: 20 * sf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := Text(tt.x, tt.y, 1, message, fontName, pts, tt.opts)
			gotWidth := text.Width()
			if math.Abs(gotWidth-wantWidth) > eps {
				t.Errorf("width = %v, want %v", gotWidth, wantWidth)
			}
			gotHeight := text.Height()
			if math.Abs(gotHeight-wantHeight) > eps {
				t.Errorf("height = %v, want %v", gotHeight, wantHeight)
			}
			render := text.render
			if math.Abs(render.Xmin-tt.wantXmin) > eps {
				t.Errorf("Xmin = %v, want %v", render.Xmin, tt.wantXmin)
			}
			if math.Abs(render.Ymin-tt.wantYmin) > eps {
				t.Errorf("Ymin = %v, want %v", render.Ymin, tt.wantYmin)
			}
			if math.Abs(render.Xmax-tt.wantXmax) > eps {
				t.Errorf("Xmax = %v, want %v", render.Xmax, tt.wantXmax)
			}
			if math.Abs(render.Ymax-tt.wantYmax) > eps {
				t.Errorf("Ymax = %v, want %v", render.Ymax, tt.wantYmax)
			}
		})
	}
}
