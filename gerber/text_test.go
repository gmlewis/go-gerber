package gerber

import (
	"math"
	"testing"

	"github.com/gmlewis/go-fonts/fonts"
	_ "github.com/gmlewis/go-fonts/fonts/freeserif"
)

func TestTextT_Primitive(t *testing.T) {
	var p Primitive = &TextT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("TextT does not implement the Primitive interface")
	}
}

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
		opts     TextOpts
		wantXmin float64
		wantYmin float64
		wantXmax float64
		wantYmax float64
	}{
		{
			name:     "XLeft,YBottom",
			wantXmin: 0,
			wantYmin: 0,
			wantXmax: wantWidth,
			wantYmax: wantHeight,
		},
		{
			name:     "XLeft,YBottom w/ offset",
			x:        10,
			y:        20,
			wantXmin: 10,
			wantYmin: 20,
			wantXmax: (10 + wantWidth),
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XCenter,YBottom",
			opts:     BottomCenter,
			wantXmin: -0.5 * wantWidth,
			wantYmin: 0,
			wantXmax: 0.5 * wantWidth,
			wantYmax: wantHeight,
		},
		{
			name:     "XCenter,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     BottomCenter,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: 20,
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XRight,YBottom",
			opts:     BottomRight,
			wantXmin: -wantWidth,
			wantYmin: 0,
			wantXmax: 0,
			wantYmax: wantHeight,
		},
		{
			name:     "XRight,YBottom w/ offset",
			x:        10,
			y:        20,
			opts:     BottomRight,
			wantXmin: (10 - wantWidth),
			wantYmin: 20,
			wantXmax: 10,
			wantYmax: (20 + wantHeight),
		},
		{
			name:     "XLeft,YCenter",
			opts:     CenterLeft,
			wantXmin: 0,
			wantYmin: -0.5 * wantHeight,
			wantXmax: wantWidth,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XLeft,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     CenterLeft,
			wantXmin: 10,
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: (10 + wantWidth),
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XCenter,YCenter",
			opts:     Center,
			wantXmin: -0.5 * wantWidth,
			wantYmin: -0.5 * wantHeight,
			wantXmax: 0.5 * wantWidth,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XCenter,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     Center,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XRight,YCenter",
			opts:     CenterRight,
			wantXmin: -wantWidth,
			wantYmin: -0.5 * wantHeight,
			wantXmax: 0,
			wantYmax: 0.5 * wantHeight,
		},
		{
			name:     "XRight,YCenter w/ offset",
			x:        10,
			y:        20,
			opts:     CenterRight,
			wantXmin: (10 - wantWidth),
			wantYmin: (20 - 0.5*wantHeight),
			wantXmax: 10,
			wantYmax: (20 + 0.5*wantHeight),
		},
		{
			name:     "XLeft,YTop",
			opts:     TopLeft,
			wantXmin: 0,
			wantYmin: -wantHeight,
			wantXmax: wantWidth,
			wantYmax: 0,
		},
		{
			name:     "XLeft,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopLeft,
			wantXmin: 10,
			wantYmin: (20 - wantHeight),
			wantXmax: (10 + wantWidth),
			wantYmax: 20,
		},
		{
			name:     "XCenter,YTop",
			opts:     TopCenter,
			wantXmin: -0.5 * wantWidth,
			wantYmin: -wantHeight,
			wantXmax: 0.5 * wantWidth,
			wantYmax: 0,
		},
		{
			name:     "XCenter,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopCenter,
			wantXmin: (10 - 0.5*wantWidth),
			wantYmin: (20 - wantHeight),
			wantXmax: (10 + 0.5*wantWidth),
			wantYmax: 20,
		},
		{
			name:     "XRight,YTop",
			opts:     TopRight,
			wantXmin: -wantWidth,
			wantYmin: -wantHeight,
			wantXmax: 0,
			wantYmax: 0,
		},
		{
			name:     "XRight,YTop w/ offset",
			x:        10,
			y:        20,
			opts:     TopRight,
			wantXmin: (10 - wantWidth),
			wantYmin: (20 - wantHeight),
			wantXmax: 10,
			wantYmax: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := Text(tt.x, tt.y, 1, message, fontName, pts, &tt.opts)
			gotWidth := text.Width()
			if math.Abs(gotWidth-wantWidth) > eps {
				t.Errorf("width = %v, want %v", gotWidth, wantWidth)
			}
			gotHeight := text.Height()
			if math.Abs(gotHeight-wantHeight) > eps {
				t.Errorf("height = %v, want %v", gotHeight, wantHeight)
			}
			mbb := text.MBB()
			if math.Abs(mbb.Min[0]-tt.wantXmin) > eps {
				t.Errorf("Xmin = %v, want %v", mbb.Min[0], tt.wantXmin)
			}
			if math.Abs(mbb.Min[1]-tt.wantYmin) > eps {
				t.Errorf("Ymin = %v, want %v", mbb.Min[1], tt.wantYmin)
			}
			if math.Abs(mbb.Max[0]-tt.wantXmax) > eps {
				t.Errorf("Xmax = %v, want %v", mbb.Max[0], tt.wantXmax)
			}
			if math.Abs(mbb.Max[1]-tt.wantYmax) > eps {
				t.Errorf("Ymax = %v, want %v", mbb.Max[1], tt.wantYmax)
			}
		})
	}
}

func TestText_Empty(t *testing.T) {
	// See: https://github.com/gmlewis/go-gerber/issues/8
	g := New("textbug")
	g.TopSilkscreen().Add(
		Text(
			25, 25,
			1.0,
			"", // should not cause panic
			"freeserif",
			12,
			&TextOpts{XAlign: fonts.XCenter, YAlign: fonts.YCenter},
		),
	)
	g.MBB() // should not panic
}
