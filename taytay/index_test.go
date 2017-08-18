package taytay

import (
	"image"
	"testing"
)

func bounds(width, height int) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	}
}

func TestClosest(t *testing.T) {

	taytays = []*TaylorSwift{
		{"1.jpg", bounds(100, 100)},
		{"2.jpg", bounds(200, 200)},
		{"3.jpg", bounds(300, 300)},
	}

	if ts := Closest(120, 120); ts.Filename != "2.jpg" {
		t.Errorf("expected 2.jpg, got: %v", ts.Filename)
	}

	if ts := Closest(199, 199); ts.Filename != "2.jpg" {
		t.Errorf("expected 2.jpg, got: %v", ts.Filename)
	}

	if ts := Closest(200, 200); ts.Filename != "2.jpg" {
		t.Errorf("expected 2.jpg, got: %v", ts.Filename)
	}

	if ts := Closest(201, 201); ts.Filename != "3.jpg" {
		t.Errorf("expected 3.jpg, got: %v", ts.Filename)
	}

	taytays = []*TaylorSwift{
		{"1.jpg", bounds(300, 600)},
		{"2.jpg", bounds(300, 1200)},
		{"3.jpg", bounds(300, 1800)},
		{"4.jpg", bounds(500, 1000)},
		{"5.jpg", bounds(500, 1500)},
		{"6.jpg", bounds(500, 2000)},
	}

	if ts := Closest(100, 230); ts.Filename != "1.jpg" {
		t.Errorf("expected 1.jpg, got: %v", ts.Filename)
	}

	if ts := Closest(400, 920); ts.Filename != "4.jpg" {
		t.Errorf("expected 4.jpg, got: %v", ts.Filename)
	}

}
