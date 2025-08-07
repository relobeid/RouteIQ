package sim_test

import (
	"testing"

	sim "routeiq/internal/sim"
)

func TestLightCycle_Timings(t *testing.T) {
	l := sim.NewLightCycle()
	for i := 0; i < 30; i++ { l.Tick() }
	if l.State() != "yellow" {
		t.Fatalf("expected yellow after 30s green, got %s", l.State())
	}
	for i := 0; i < 5; i++ { l.Tick() }
	if l.State() != "red" {
		t.Fatalf("expected red after yellow, got %s", l.State())
	}
	for i := 0; i < 25; i++ { l.Tick() }
	if l.State() != "green" {
		t.Fatalf("expected green after red, got %s", l.State())
	}
}

func TestLightCycle_ElapsedAndFullCycle(t *testing.T) {
	l := sim.NewLightCycle()
	if l.State() != "green" || l.Elapsed() != 0 {
		t.Fatalf("expected initial state green, elapsed 0")
	}
	for i := 0; i < 30; i++ { l.Tick() }
	if l.State() != "yellow" || l.Elapsed() != 0 {
		t.Fatalf("expected yellow with elapsed reset, got state=%s elapsed=%d", l.State(), l.Elapsed())
	}
	for i := 0; i < 5; i++ { l.Tick() }
	if l.State() != "red" || l.Elapsed() != 0 {
		t.Fatalf("expected red with elapsed reset, got state=%s elapsed=%d", l.State(), l.Elapsed())
	}
	for i := 0; i < 25; i++ { l.Tick() }
	if l.State() != "green" || l.Elapsed() != 0 {
		t.Fatalf("expected green with elapsed reset, got state=%s elapsed=%d", l.State(), l.Elapsed())
	}
	for i := 0; i < 10; i++ { l.Tick() }
	if l.State() != "green" || l.Elapsed() != 10 {
		t.Fatalf("expected green elapsed=10, got state=%s elapsed=%d", l.State(), l.Elapsed())
	}
}
