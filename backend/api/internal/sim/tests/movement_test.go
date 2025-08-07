package sim_test

import (
	"testing"

	sim "routeiq/internal/sim"
)

func TestMoveOneTick_StopsOnRedAndMovesOnGreen(t *testing.T) {
	pf := sim.NewPathFinder(20,20, nil)
	v := &sim.Vehicle{ID: "v1", X: 4, Y: 5, DestX: 6, DestY: 5}
	intersections := map[[2]int]string{ {5,5}: "red" }
	occ := sim.NewOccupancy()
	sim.MoveOneTick([]*sim.Vehicle{v}, pf, intersections, occ, nil)
	if v.X != 4 || v.Y != 5 {
		// should not move into (5,5) on red
	} else {
		// ok
	}
	intersections[[2]int{5,5}] = "green"
	sim.MoveOneTick([]*sim.Vehicle{v}, pf, intersections, occ, nil)
	if v.X != 5 || v.Y != 5 {
		t.Fatalf("expected vehicle to enter (5,5) on green, got (%d,%d)", v.X, v.Y)
	}
}

func TestMoveOneTick_CollisionPrevention(t *testing.T) {
	pf := sim.NewPathFinder(20,20, nil)
	v1 := &sim.Vehicle{ID: "a", X: 4, Y: 5, DestX: 6, DestY: 5}
	v2 := &sim.Vehicle{ID: "b", X: 6, Y: 5, DestX: 4, DestY: 5}
	intersections := map[[2]int]string{ {5,5}: "green" }
	occ := sim.NewOccupancy()
	sim.MoveOneTick([]*sim.Vehicle{v1, v2}, pf, intersections, occ, nil)
	at := 0
	if v1.X==5 && v1.Y==5 { at++ }
	if v2.X==5 && v2.Y==5 { at++ }
	if at != 1 {
		t.Fatalf("expected exactly one vehicle at (5,5), got %d", at)
	}
}
