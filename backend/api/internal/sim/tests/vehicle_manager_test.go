package sim_test

import (
	"testing"

	sim "routeiq/internal/sim"
)

func TestVehicleManager_SpawnAndUniqueness(t *testing.T) {
	m := sim.NewVehicleManager()
	ids := m.Spawn(150, 20, 20)
	if len(ids) != 150 {
		t.Fatalf("expected 150 ids, got %d", len(ids))
	}
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate id generated: %s", id)
		}
		seen[id] = true
		if v, ok := m.Get(id); !ok || v == nil {
			t.Fatalf("vehicle not retrievable: %s", id)
		}
	}
	if m.Count() != 150 {
		t.Fatalf("expected count 150, got %d", m.Count())
	}
}

func TestVehicleManager_Despawn(t *testing.T) {
	m := sim.NewVehicleManager()
	ids := m.Spawn(10, 20, 20)
	removed := m.Despawn(ids[:5]...)
	if removed != 5 {
		t.Fatalf("expected removed 5, got %d", removed)
	}
	if m.Count() != 5 {
		t.Fatalf("expected count 5, got %d", m.Count())
	}
}

func TestVehicleManager_ListAndGetNotFound(t *testing.T) {
	m := sim.NewVehicleManager()
	if _, ok := m.Get("nope"); ok {
		t.Fatalf("expected not found for unknown id")
	}
	m.Spawn(3, 2, 2)
	list := m.List()
	if len(list) != 3 {
		t.Fatalf("expected list len 3, got %d", len(list))
	}
}

func TestVehicleManager_DespawnEdgeCases(t *testing.T) {
	m := sim.NewVehicleManager()
	ids := m.Spawn(2, 2, 2)
	removed := m.Despawn(ids[0], "unknown", ids[0])
	if removed != 1 {
		t.Fatalf("expected to remove exactly 1, got %d", removed)
	}
	if m.Count() != 1 {
		t.Fatalf("expected count 1, got %d", m.Count())
	}
	removed = m.Despawn(ids[1])
	if removed != 1 || m.Count() != 0 {
		t.Fatalf("expected final remove=1 and count=0, got removed=%d count=%d", removed, m.Count())
	}
	if m.Despawn("still-unknown") != 0 {
		t.Fatalf("expected remove=0 on empty manager")
	}
}

func TestVehicleManager_SpawnZero(t *testing.T) {
	m := sim.NewVehicleManager()
	ids := m.Spawn(0, 10, 10)
	if len(ids) != 0 || m.Count() != 0 {
		t.Fatalf("spawn(0) should not change state")
	}
}
