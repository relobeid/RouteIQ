package sim

import (
	"container/heap"
	"math"
)

type point struct{ X, Y int }

type PathFinder struct{
	Width int
	Height int
	Blocked map[[2]int]bool // blocked cells
}

func NewPathFinder(width, height int, blocked map[[2]int]bool) *PathFinder {
	if blocked == nil { blocked = make(map[[2]int]bool) }
	return &PathFinder{Width: width, Height: height, Blocked: blocked}
}

func (p *PathFinder) inBounds(x,y int) bool {
	return x>=0 && x<p.Width && y>=0 && y<p.Height
}

func (p *PathFinder) isBlocked(x,y int) bool {
	return p.Blocked[[2]int{x,y}]
}

func manhattan(a, b point) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

type node struct {
	pt point
	g  int // cost from start
	h  int // heuristic to goal
	f  int // g+h
	idx int // heap idx
}

type nodePQ []*node

func (pq nodePQ) Len() int { return len(pq) }
func (pq nodePQ) Less(i, j int) bool { return pq[i].f < pq[j].f }
func (pq nodePQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i]; pq[i].idx=i; pq[j].idx=j }
func (pq *nodePQ) Push(x any) { n := x.(*node); n.idx = len(*pq); *pq = append(*pq, n) }
func (pq *nodePQ) Pop() any { old := *pq; n := len(old); x := old[n-1]; *pq = old[:n-1]; return x }

// Path returns a sequence of points from start to goal, inclusive. Empty if no path.
func (p *PathFinder) Path(sx, sy, gx, gy int) []point {
	start := point{sx,sy}
	goal := point{gx,gy}
	if !p.inBounds(sx,sy) || !p.inBounds(gx,gy) || p.isBlocked(gx,gy) { return nil }
	if sx==gx && sy==gy { return []point{start} }

	came := make(map[point]point)
	gscore := make(map[point]int)
	gscore[start] = 0

	open := &nodePQ{}
	heap.Init(open)
	heap.Push(open, &node{pt:start, g:0, h:manhattan(start,goal), f:manhattan(start,goal)})
	inOpen := map[point]*node{start: (*open)[0]}
	closed := make(map[point]bool)

	neighbors := func(c point) []point {
		cand := []point{{c.X+1,c.Y},{c.X-1,c.Y},{c.X,c.Y+1},{c.X,c.Y-1}}
		out := make([]point,0,4)
		for _, n := range cand {
			if p.inBounds(n.X,n.Y) && !p.isBlocked(n.X,n.Y) { out = append(out, n) }
		}
		return out
	}

	for open.Len()>0 {
		cur := heap.Pop(open).(*node)
		delete(inOpen, cur.pt)
		if cur.pt == goal { // reconstruct
			var path []point
			u := goal
			for u != start {
				path = append(path, u)
				u = came[u]
			}
			path = append(path, start)
			// reverse
			for i,j := 0, len(path)-1; i<j; i,j = i+1, j-1 { path[i],path[j] = path[j],path[i] }
			return path
		}
		closed[cur.pt] = true

		for _, nb := range neighbors(cur.pt) {
			if closed[nb] { continue }
			tentative := cur.g + 1
			if g, ok := gscore[nb]; !ok || tentative < g {
				came[nb] = cur.pt
				gscore[nb] = tentative
				h := manhattan(nb, goal)
				if on, ok := inOpen[nb]; ok {
					on.g = tentative; on.h=h; on.f = tentative + h
					heap.Fix(open, on.idx)
				} else {
					n := &node{pt: nb, g: tentative, h: h, f: tentative + h}
					heap.Push(open, n)
					inOpen[nb] = n
				}
			}
		}
	}
	return nil
}
