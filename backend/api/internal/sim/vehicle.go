package sim

import "time"

// Vehicle represents a simulated vehicle in the grid.
type Vehicle struct {
	ID        string
	X         int
	Y         int
	Speed     float64 // cells per second
	DestX     int
	DestY     int
	CreatedAt time.Time
}
