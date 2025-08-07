package sim

import (
	"math/rand/v2"
	"sync"
	"time"

	"github.com/google/uuid"
)

// VehicleManager manages the lifecycle of vehicles in-memory.
type VehicleManager struct {
	mu       sync.RWMutex
	vehicles map[string]*Vehicle
}

func NewVehicleManager() *VehicleManager {
	return &VehicleManager{vehicles: make(map[string]*Vehicle)}
}

// Spawn creates n vehicles at random positions within bounds [0,width) x [0,height).
func (m *VehicleManager) Spawn(n, width, height int) []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		id := uuid.New().String()
		x := rand.IntN(width)
		y := rand.IntN(height)
		dx := rand.IntN(width)
		dy := rand.IntN(height)
		v := &Vehicle{
			ID:        id,
			X:         x,
			Y:         y,
			Speed:     1.0,
			DestX:     dx,
			DestY:     dy,
			CreatedAt: time.Now(),
		}
		m.vehicles[id] = v
		ids = append(ids, id)
	}
	return ids
}

// Despawn removes the vehicles with the provided IDs.
func (m *VehicleManager) Despawn(ids ...string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	removed := 0
	for _, id := range ids {
		if _, ok := m.vehicles[id]; ok {
			delete(m.vehicles, id)
			removed++
		}
	}
	return removed
}

// Get returns the vehicle by ID.
func (m *VehicleManager) Get(id string) (*Vehicle, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.vehicles[id]
	return v, ok
}

// Count returns the number of active vehicles.
func (m *VehicleManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.vehicles)
}

// List returns a snapshot of vehicles.
func (m *VehicleManager) List() []*Vehicle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*Vehicle, 0, len(m.vehicles))
	for _, v := range m.vehicles {
		out = append(out, v)
	}
	return out
}
