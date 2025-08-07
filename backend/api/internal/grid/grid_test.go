package grid

import "testing"

func TestNewGrid_Builds400Cells(t *testing.T) {
	g := NewGrid(20, 20)
	if g.NumCells() != 400 {
		t.Fatalf("expected 400 cells, got %d", g.NumCells())
	}
}

func TestIntersections_Placement(t *testing.T) {
	g := NewGrid(20, 20)
	expected := map[[2]int]bool{
		{5, 5}: true,
		{5, 15}: true,
		{15, 5}: true,
		{15, 15}: true,
	}
	got := map[[2]int]bool{}
	for _, it := range g.Intersections() {
		got[[2]int{it.X, it.Y}] = true
		if it.LightState != LightRed {
			t.Fatalf("expected default light state red, got %s", it.LightState)
		}
	}
	for k := range expected {
		if !got[k] {
			t.Fatalf("missing intersection at %v", k)
		}
	}
}

func TestToID_FromID_RoundTrip(t *testing.T) {
	g := NewGrid(20, 20)
	cases := [][2]int{
		{0, 0}, {19, 0}, {0, 19}, {19, 19}, // corners
		{10, 0}, {0, 10}, {19, 10}, {10, 19}, // edges
		{5, 5}, {7, 13}, {12, 4}, {15, 15}, // random
	}
	for _, c := range cases {
		x, y := c[0], c[1]
		id := g.ToID(x, y)
		if id == -1 {
			t.Fatalf("unexpected -1 id for valid (%d,%d)", x, y)
		}
		rx, ry := g.FromID(id)
		if rx != x || ry != y {
			t.Fatalf("roundtrip mismatch: (%d,%d)->%d->(%d,%d)", x, y, id, rx, ry)
		}
	}
}

func TestBoundsChecks(t *testing.T) {
	g := NewGrid(20, 20)
	invalid := [][2]int{{-1, 0}, {0, -1}, {20, 0}, {0, 20}, {20, 20}}
	for _, p := range invalid {
		if g.IsValid(p[0], p[1]) {
			t.Fatalf("expected invalid for %v", p)
		}
		if id := g.ToID(p[0], p[1]); id != -1 {
			t.Fatalf("expected -1 id for %v, got %d", p, id)
		}
	}
	if x, y := g.FromID(-1); x != -1 || y != -1 {
		t.Fatalf("FromID(-1) should return (-1,-1), got (%d,%d)", x, y)
	}
	if x, y := g.FromID(400); x != -1 || y != -1 { // 0..399 are valid
		t.Fatalf("FromID(400) should return (-1,-1), got (%d,%d)", x, y)
	}
}
