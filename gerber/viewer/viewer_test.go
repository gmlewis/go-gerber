package viewer

import (
	"testing"

	"github.com/gmlewis/go-gerber/gerber"
)

func TestMBB(t *testing.T) {
	tests := []struct {
		name          string
		w, h          int
		scaleOverride float64
		vc            *viewController
		want          gerber.MBB
	}{
		{
			name: "unit scale",
			w:    2,
			h:    2,
			vc: &viewController{
				mbb:    gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
				center: gerber.Pt{0, 0},
			},
			want: gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
		},
		{
			name:          "normal window, max zoom",
			w:             801,
			h:             801,
			scaleOverride: 1.0,
			vc: &viewController{
				mbb:    gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
				center: gerber.Pt{0, 0},
			},
			want: gerber.MBB{Min: gerber.Pt{-400, -400}, Max: gerber.Pt{400, 400}},
		},
		{
			name: "normal window, fit design",
			w:    800,
			h:    800,
			vc: &viewController{
				mbb:    gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
				center: gerber.Pt{0, 0},
			},
			want: gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
		},
		{
			name: "wide window, fit design",
			w:    1601,
			h:    801,
			vc: &viewController{
				mbb:    gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
				center: gerber.Pt{0, 0},
			},
			want: gerber.MBB{Min: gerber.Pt{-2, -1}, Max: gerber.Pt{2, 1}},
		},
		{
			name: "tall window, fit design",
			w:    801,
			h:    1601,
			vc: &viewController{
				mbb:    gerber.MBB{Min: gerber.Pt{-1, -1}, Max: gerber.Pt{1, 1}},
				center: gerber.Pt{0, 0},
			},
			want: gerber.MBB{Min: gerber.Pt{-1, -2}, Max: gerber.Pt{1, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.vc.scaleToFit(tt.w, tt.h)
			if tt.scaleOverride != 0.0 {
				tt.vc.scale = tt.scaleOverride
			}
			got := tt.vc.MBB()
			if *got != tt.want {
				t.Errorf("MBB=%v, want %v", got, tt.want)
			}

			// Test xf,yf:
			xf := tt.vc.xf(&tt.want)
			yf := tt.vc.yf(&tt.want)
			llx, lly := xf(tt.want.Min[0]), yf(tt.want.Min[1])
			if want := 0.0; llx != want {
				t.Errorf("llx=%v, want %v", llx, want)
			}
			if want := float64(tt.vc.lastH - 1); lly != want {
				t.Errorf("lly=%v, want %v", lly, want)
			}

			cx, cy := xf(0.5*(tt.want.Max[0]+tt.want.Min[0])), yf(0.5*(tt.want.Max[1]+tt.want.Min[1]))
			if want := 0.5 * float64(tt.vc.lastW-1); cx != want {
				t.Errorf("cx=%v, want %v", cx, want)
			}
			if want := 0.5 * float64(tt.vc.lastH-1); cy != want {
				t.Errorf("cy=%v, want %v", cy, want)
			}

			urx, ury := xf(tt.want.Max[0]), yf(tt.want.Max[1])
			if want := float64(tt.vc.lastW - 1); urx != want {
				t.Errorf("urx=%v, want %v", urx, want)
			}
			if want := 0.0; ury != want {
				t.Errorf("ury=%v, want %v", ury, want)
			}
		})
	}
}
