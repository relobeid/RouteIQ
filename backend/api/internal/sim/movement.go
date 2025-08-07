package sim

// Occupancy tracks reserved cells per tick to avoid collisions.
type Occupancy struct{ cells map[[2]int]bool }
func NewOccupancy() *Occupancy { return &Occupancy{cells: make(map[[2]int]bool)} }
func (o *Occupancy) Reset() { o.cells = make(map[[2]int]bool) }
func (o *Occupancy) TryReserve(x,y int) bool {
	k := [2]int{x,y}
	if o.cells[k] { return false }
	o.cells[k] = true
	return true
}

// MoveOneTick moves vehicles by at most one cell along their path, respecting intersection lights and collisions.
// intersections: map of intersection coords -> current light state ("red"/"yellow"/"green").
func MoveOneTick(vehicles []*Vehicle, pf *PathFinder, intersections map[[2]int]string, occ *Occupancy, dests map[string]point) {
	if occ == nil { occ = NewOccupancy() } else { occ.Reset() }
	for _, v := range vehicles {
		d, ok := dests[v.ID]
		if !ok { d = point{v.DestX, v.DestY} }
		path := pf.Path(v.X, v.Y, d.X, d.Y)
		if len(path) <= 1 { // at destination or no path
			continue
		}
		next := path[1]
		// stop at red or yellow when entering an intersection cell
		if state, isIntersection := intersections[[2]int{next.X, next.Y}]; isIntersection {
			if state == "red" || state == "yellow" {
				continue
			}
		}
		// collision prevention: reserve next cell
		if !occ.TryReserve(next.X, next.Y) {
			continue
		}
		v.X, v.Y = next.X, next.Y
	}
}
