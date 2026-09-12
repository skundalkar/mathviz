// Package astar visualizes A* search: Dijkstra's algorithm with one change
// -- instead of always expanding whichever not-yet-closed node has the
// smallest cost-so-far, it expands whichever has the smallest cost-so-far
// PLUS a heuristic estimate of the remaining distance to a fixed goal. Both
// are built from the exact same relaxation rule dijkstras-algorithm used;
// the only difference is what the search sorts by, and dijkstras-algorithm
// is exactly the special case where that heuristic is zero everywhere.
package astar

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "a-star-search",
		Seq:   99,
		Title: "A* search (Dijkstra with a distance-to-goal hint)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "heuristic", Label: "Use heuristic (0=plain Dijkstra, 1=A*)", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "step", Label: "Step (expand one node at a time)", Min: 0, Max: 6, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// NodeNames labels the fixed 6-node example graph every Section walks
// through, in the same index order Pos and Graph use: S=0, A=1, B=2, C=3,
// D=4, G=5.
var NodeNames = []string{"S", "A", "B", "C", "D", "G"}

// Pos gives each node's 2D coordinates. Both Graph's edge weights and
// Heuristic are plain Euclidean distances derived from these same
// coordinates -- which is exactly what guarantees the heuristic never
// overestimates the true remaining distance: the straight line from any
// node to the goal can never be longer than a route that detours through a
// third point (the triangle inequality).
var Pos = [][2]float64{
	{0, 0},  // S
	{2, 3},  // A
	{4, 0},  // B
	{6, 3},  // C
	{8, 0},  // D
	{10, 0}, // G
}

// Goal is the fixed destination every search in this concept runs toward.
const Goal = 5 // G

// Edge is one directed edge in an adjacency list: To is the destination
// node's index, Weight its non-negative cost.
type Edge struct {
	To     int
	Weight float64
}

func dist(a, b [2]float64) float64 {
	dx, dy := a[0]-b[0], a[1]-b[1]
	return math.Sqrt(dx*dx + dy*dy)
}

// roadPairs lists each undirected road once, as (from, to) node indices;
// buildGraph mirrors every pair into both directions with its Euclidean
// weight.
var roadPairs = [][2]int{
	{0, 1}, // S-A
	{0, 2}, // S-B
	{1, 2}, // A-B
	{1, 3}, // A-C
	{2, 3}, // B-C
	{2, 4}, // B-D
	{3, 4}, // C-D
	{3, 5}, // C-G
	{4, 5}, // D-G
}

func buildGraph() [][]Edge {
	g := make([][]Edge, len(Pos))
	for _, pr := range roadPairs {
		w := dist(Pos[pr[0]], Pos[pr[1]])
		g[pr[0]] = append(g[pr[0]], Edge{pr[1], w})
		g[pr[1]] = append(g[pr[1]], Edge{pr[0], w})
	}
	return g
}

// Graph is the fixed 6-node road network every Section walks through: S-A,
// S-B, A-B, A-C, B-C, B-D, C-D, C-G, D-G, each edge weighted by the
// Euclidean distance between its endpoints' Pos.
var Graph = buildGraph()

// Heuristic estimates the remaining distance from node n to goal as the
// straight-line (Euclidean) distance between their coordinates. Because
// every edge weight in Graph is also a Euclidean distance between the same
// coordinates, the triangle inequality guarantees this never overestimates
// the true cheapest remaining route -- the "admissible" property A* needs
// in order to still guarantee it finds the optimal path.
func Heuristic(n, goal int) float64 {
	return dist(Pos[n], Pos[goal])
}

// Step is one snapshot of the search: the cost-so-far (G) and current-best
// predecessor for every node right after one node is expanded (or, for the
// first Step, before anything has been expanded at all). Current is the
// node expanded this step (-1 for the initial "start" step); Relaxed lists
// every neighbor whose cost-so-far was lowered as a result.
type Step struct {
	G           []float64
	Prev        []int
	Closed      []bool
	Current     int
	Relaxed     []int
	Description string
}

func cloneFloats(v []float64) []float64 { return append([]float64(nil), v...) }
func cloneInts(v []int) []int           { return append([]int(nil), v...) }
func cloneBools(v []bool) []bool        { return append([]bool(nil), v...) }

// Search runs a best-first search from src toward goal, expanding whichever
// not-yet-closed node has the smallest g(n) (plus h(n,goal) when
// useHeuristic is true) each round. useHeuristic=false is plain Dijkstra;
// useHeuristic=true is A*. It stops the moment goal itself is expanded:
// because Heuristic never overestimates, that already is goal's cheapest
// possible cost, so there is nothing left to gain by continuing.
func Search(g [][]Edge, src, goal int, useHeuristic bool) []Step {
	n := len(g)
	gscore := make([]float64, n)
	prev := make([]int, n)
	closed := make([]bool, n)
	open := make([]bool, n)
	for i := range gscore {
		gscore[i] = math.Inf(1)
		prev[i] = -1
	}
	gscore[src] = 0
	open[src] = true

	steps := []Step{{
		G: cloneFloats(gscore), Prev: cloneInts(prev), Closed: cloneBools(closed),
		Current: -1, Description: "Start: g[src]=0, every other node unreached",
	}}

	for {
		u := -1
		bestF := math.Inf(1)
		for v := 0; v < n; v++ {
			if !open[v] || closed[v] {
				continue
			}
			f := gscore[v]
			if useHeuristic {
				f += Heuristic(v, goal)
			}
			if f < bestF {
				bestF = f
				u = v
			}
		}
		if u == -1 {
			break // open exhausted -- goal unreachable
		}
		closed[u] = true

		var relaxed []int
		for _, e := range g[u] {
			if closed[e.To] {
				continue
			}
			nd := gscore[u] + e.Weight
			if nd < gscore[e.To] {
				gscore[e.To] = nd
				prev[e.To] = u
				open[e.To] = true
				relaxed = append(relaxed, e.To)
			}
		}

		steps = append(steps, Step{
			G: cloneFloats(gscore), Prev: cloneInts(prev), Closed: cloneBools(closed),
			Current: u, Relaxed: relaxed,
			Description: "Expand the open node with the smallest priority and relax its neighbors",
		})

		if u == goal {
			break // goal's cost is now final -- nothing left to gain
		}
	}
	return steps
}

// PathTo walks prev backward from target to src and returns the route as a
// slice of node indices from src to target (inclusive). It returns nil if
// target hasn't been reached (prev[target] == -1 and target != src).
func PathTo(prev []int, src, target int) []int {
	if target != src && prev[target] == -1 {
		return nil
	}
	path := []int{target}
	for v := target; v != src; {
		v = prev[v]
		path = append(path, v)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
