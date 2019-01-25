package gerber

import "testing"

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

func TestCircleT_Primitive(t *testing.T) {
	var p Primitive = &CircleT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("CircleT does not implement the Primitive interface")
	}
}

func TestLineT_Primitive(t *testing.T) {
	var p Primitive = &LineT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("LineT does not implement the Primitive interface")
	}
}

func TestPolygonT_Primitive(t *testing.T) {
	var p Primitive = &PolygonT{}
	if p == nil {
		// In actuality, this test won't compile if it isn't a Primitive.
		t.Errorf("PolygonT does not implement the Primitive interface")
	}
}
