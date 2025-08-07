package grid

// LightState represents the state of a traffic light at an intersection.
type LightState string

const (
	LightRed    LightState = "red"
	LightYellow LightState = "yellow"
	LightGreen  LightState = "green"
)

// Intersection represents a major intersection location in the grid.
type Intersection struct {
	X          int
	Y          int
	LightState LightState
}

// Cell represents a single cell in the traffic grid.
type Cell struct {
	X int
	Y int
}

// Grid represents the city grid and intersection layout.
type Grid struct {
	Width         int
	Height        int
	cells         []Cell
	intersections []Intersection
}

// NewGrid constructs a grid of Width x Height and seeds four major intersections.
func NewGrid(width, height int) *Grid {
	g := &Grid{Width: width, Height: height}
	g.buildCells()
	g.seedDefaultIntersections()
	return g
}

func (g *Grid) buildCells() {
	total := g.Width * g.Height
	g.cells = make([]Cell, 0, total)
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.cells = append(g.cells, Cell{X: x, Y: y})
		}
	}
}

func (g *Grid) seedDefaultIntersections() {
	// Default four major intersections for a 20x20 grid
	defaults := [][2]int{
		{5, 5}, {5, 15}, {15, 5}, {15, 15},
	}
	for _, p := range defaults {
		x, y := p[0], p[1]
		if g.IsValid(x, y) {
			g.intersections = append(g.intersections, Intersection{X: x, Y: y, LightState: LightRed})
		}
	}
}

// NumCells returns the total number of cells in the grid.
func (g *Grid) NumCells() int { return g.Width * g.Height }

// IsValid reports whether (x,y) is inside the grid bounds.
func (g *Grid) IsValid(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

// ToID maps (x,y) to a linear cell id. Returns -1 if out of bounds.
func (g *Grid) ToID(x, y int) int {
	if !g.IsValid(x, y) {
		return -1
	}
	return y*g.Width + x
}

// FromID maps a linear cell id to (x,y). Returns (-1,-1) if id is invalid.
func (g *Grid) FromID(id int) (int, int) {
	if id < 0 || id >= g.NumCells() {
		return -1, -1
	}
	x := id % g.Width
	y := id / g.Width
	return x, y
}

// Intersections returns a copy of the seeded intersections for safety.
func (g *Grid) Intersections() []Intersection {
	out := make([]Intersection, len(g.intersections))
	copy(out, g.intersections)
	return out
}
