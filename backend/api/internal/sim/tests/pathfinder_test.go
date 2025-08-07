package sim_test

import (
	"testing"
	"time"

	sim "routeiq/internal/sim"
)

func TestPathFinder_BasicPath(t *testing.T) {
	pf := sim.NewPathFinder(20,20, nil)
	p := pf.Path(0,0, 19,19)
	if len(p) == 0 {
		t.Fatalf("expected path, got none")
	}
	if len(p) < 39 {
		t.Fatalf("expected path length >=39, got %d", len(p))
	}
}

func TestPathFinder_WithBlockersGap(t *testing.T) {
	blocked := make(map[[2]int]bool)
	for y:=0; y<20; y++ { if y!=10 { blocked[[2]int{10,y}] = true } }
	pf := sim.NewPathFinder(20,20, blocked)
	p := pf.Path(0,10, 19,10)
	if len(p) == 0 {
		t.Fatalf("expected path through the gap, got none")
	}
	seen := false
	for _, pt := range p {
		if pt.X==10 && pt.Y==10 { seen = true; break }
	}
	if !seen {
		t.Fatalf("expected path to cross (10,10) gap")
	}
}

func TestPathFinder_PerfUnder50ms(t *testing.T) {
	pf := sim.NewPathFinder(20,20, nil)
	start := time.Now()
	_ = pf.Path(0,0, 19,19)
	if time.Since(start) > 50*time.Millisecond {
		t.Fatalf("expected under 50ms, took %s", time.Since(start))
	}
}
