package gerber

import (
	"math"
	"testing"
)

func TestAperture_Primitive(t *testing.T) {
	var p Primitive = &Aperture{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("Aperture does not implement the Primitive interface")
	}
}

func TestArcT_Primitive(t *testing.T) {
	var p Primitive = &ArcT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("ArcT does not implement the Primitive interface")
	}
}

func TestArcT_MBB(t *testing.T) {
	const eps = 1e-3
	tests := []struct {
		name string
		p    *ArcT
		want MBB
	}{
		{
			name: "unit circle",
			p:    Arc(Pt{0, 0}, 1, CircleShape, 1, 1, 0, 360, 0),
			want: MBB{Min: Pt{-1, -1}, Max: Pt{1, 1}},
		},
		{
			name: "unit circle w/ thickness",
			p:    Arc(Pt{0, 0}, 10, CircleShape, 1, 1, 0, 360, 2),
			want: MBB{Min: Pt{-11, -11}, Max: Pt{11, 11}},
		},
		{
			name: "first quadrant arc",
			p:    Arc(Pt{0, 0}, 10, CircleShape, 1, 1, 0, 90, 2),
			want: MBB{Min: Pt{-1, -1}, Max: Pt{11, 11}},
		},
		{
			name: "second quadrant arc",
			p:    Arc(Pt{0, 0}, 10, CircleShape, 1, 1, 90, 180, 2),
			want: MBB{Min: Pt{-11, -1}, Max: Pt{1, 11}},
		},
		{
			name: "third quadrant arc",
			p:    Arc(Pt{0, 0}, 10, CircleShape, 1, 1, 180, 270, 2),
			want: MBB{Min: Pt{-11, -11}, Max: Pt{1, 1}},
		},
		{
			name: "fourth quadrant arc",
			p:    Arc(Pt{0, 0}, 10, CircleShape, 1, 1, 270, 360, 2),
			want: MBB{Min: Pt{-1, -11}, Max: Pt{11, 1}},
		},
		{
			name: "first quadrant arc w/ offset",
			p:    Arc(Pt{10, 20}, 10, CircleShape, 1, 1, 0, 90, 2),
			want: MBB{Min: Pt{9, 19}, Max: Pt{21, 31}},
		},
		{
			name: "second quadrant arc w/ offset",
			p:    Arc(Pt{10, 20}, 10, CircleShape, 1, 1, 90, 180, 2),
			want: MBB{Min: Pt{-1, 19}, Max: Pt{11, 31}},
		},
		{
			name: "third quadrant arc w/ offset",
			p:    Arc(Pt{10, 20}, 10, CircleShape, 1, 1, 180, 270, 2),
			want: MBB{Min: Pt{-1, 9}, Max: Pt{11, 21}},
		},
		{
			name: "fourth quadrant arc w/ offset",
			p:    Arc(Pt{10, 20}, 10, CircleShape, 1, 1, 270, 360, 2),
			want: MBB{Min: Pt{9, 9}, Max: Pt{21, 21}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.MBB()
			if math.Abs(got.Min[0]-tt.want.Min[0]) > eps {
				t.Errorf("Min[0]=%v, want %v", got.Min[0], tt.want.Min[0])
			}
			if math.Abs(got.Min[1]-tt.want.Min[1]) > eps {
				t.Errorf("Min[1]=%v, want %v", got.Min[1], tt.want.Min[1])
			}
			if math.Abs(got.Max[0]-tt.want.Max[0]) > eps {
				t.Errorf("Max[0]=%v, want %v", got.Max[0], tt.want.Max[0])
			}
			if math.Abs(got.Max[1]-tt.want.Max[1]) > eps {
				t.Errorf("Max[1]=%v, want %v", got.Max[1], tt.want.Max[1])
			}
		})
	}
}

func TestCircleT_Primitive(t *testing.T) {
	var p Primitive = &CircleT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("CircleT does not implement the Primitive interface")
	}
}

func TestCircleT_MBB(t *testing.T) {
	const eps = 1e-12
	tests := []struct {
		name string
		p    *CircleT
		want MBB
	}{
		{
			name: "unit circle",
			p:    Circle(Pt{0, 0}, 1),
			want: MBB{Min: Pt{-0.5, -0.5}, Max: Pt{0.5, 0.5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.MBB()
			if math.Abs(got.Min[0]-tt.want.Min[0]) > eps {
				t.Errorf("Min[0]=%v, want %v", got.Min[0], tt.want.Min[0])
			}
			if math.Abs(got.Min[1]-tt.want.Min[1]) > eps {
				t.Errorf("Min[1]=%v, want %v", got.Min[1], tt.want.Min[1])
			}
			if math.Abs(got.Max[0]-tt.want.Max[0]) > eps {
				t.Errorf("Max[0]=%v, want %v", got.Max[0], tt.want.Max[0])
			}
			if math.Abs(got.Max[1]-tt.want.Max[1]) > eps {
				t.Errorf("Max[1]=%v, want %v", got.Max[1], tt.want.Max[1])
			}
		})
	}
}

func TestLineT_Primitive(t *testing.T) {
	var p Primitive = &LineT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("LineT does not implement the Primitive interface")
	}
}

func TestLineT_MBB(t *testing.T) {
	const eps = 1e-12
	tests := []struct {
		name string
		p    *LineT
		want MBB
	}{
		{
			name: "horizontal line",
			p:    Line(0, 0, 1, 0, CircleShape, 1),
			want: MBB{Min: Pt{-0.5, -0.5}, Max: Pt{1.5, 0.5}},
		},
		{
			name: "vertical line",
			p:    Line(0, 0, 0, 1, CircleShape, 1),
			want: MBB{Min: Pt{-0.5, -0.5}, Max: Pt{0.5, 1.5}},
		},
		{
			name: "diagonal line",
			p:    Line(0, 0, 1, 1, CircleShape, 1),
			want: MBB{Min: Pt{-0.5, -0.5}, Max: Pt{1.5, 1.5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.MBB()
			if math.Abs(got.Min[0]-tt.want.Min[0]) > eps {
				t.Errorf("Min[0]=%v, want %v", got.Min[0], tt.want.Min[0])
			}
			if math.Abs(got.Min[1]-tt.want.Min[1]) > eps {
				t.Errorf("Min[1]=%v, want %v", got.Min[1], tt.want.Min[1])
			}
			if math.Abs(got.Max[0]-tt.want.Max[0]) > eps {
				t.Errorf("Max[0]=%v, want %v", got.Max[0], tt.want.Max[0])
			}
			if math.Abs(got.Max[1]-tt.want.Max[1]) > eps {
				t.Errorf("Max[1]=%v, want %v", got.Max[1], tt.want.Max[1])
			}
		})
	}
}

func TestPolygonT_Primitive(t *testing.T) {
	var p Primitive = &PolygonT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("PolygonT does not implement the Primitive interface")
	}
}

func TestPolygonT_MBB(t *testing.T) {
	const eps = 1e-12
	tests := []struct {
		name string
		p    *PolygonT
		want MBB
	}{
		{
			name: "unit box",
			p:    Polygon(Pt{0, 0}, true, []Pt{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, 0),
			want: MBB{Min: Pt{0, 0}, Max: Pt{1, 1}},
		},
		{
			name: "shifted unit box",
			p:    Polygon(Pt{10, 20}, true, []Pt{{0, 0}, {1, 0}, {1, 1}, {0, 1}}, 0),
			want: MBB{Min: Pt{10, 20}, Max: Pt{11, 21}},
		},
		{
			name: "centered box",
			p:    Polygon(Pt{0, 0}, true, []Pt{{-1, -1}, {1, -1}, {1, 1}, {-1, 1}}, 0),
			want: MBB{Min: Pt{-1, -1}, Max: Pt{1, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.MBB()
			if math.Abs(got.Min[0]-tt.want.Min[0]) > eps {
				t.Errorf("Min[0]=%v, want %v", got.Min[0], tt.want.Min[0])
			}
			if math.Abs(got.Min[1]-tt.want.Min[1]) > eps {
				t.Errorf("Min[1]=%v, want %v", got.Min[1], tt.want.Min[1])
			}
			if math.Abs(got.Max[0]-tt.want.Max[0]) > eps {
				t.Errorf("Max[0]=%v, want %v", got.Max[0], tt.want.Max[0])
			}
			if math.Abs(got.Max[1]-tt.want.Max[1]) > eps {
				t.Errorf("Max[1]=%v, want %v", got.Max[1], tt.want.Max[1])
			}
		})
	}
}
