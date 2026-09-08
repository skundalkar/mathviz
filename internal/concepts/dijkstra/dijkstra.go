// Package dijkstra visualizes Dijkstra's algorithm: finding the cheapest
// route from one source node to every other node in a weighted graph by
// repeatedly locking in the not-yet-finalized node with the smallest known
// distance, then using it to shorten ("relax") its neighbors' distances --
// the same graph-of-nodes idea pagerank walked, but now every edge carries
// a cost instead of counting equally, so the algorithm has to keep
// reconsidering a node's distance until it's actually finalized.
package dijkstra

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "dijkstras-algorithm",
		Seq:   96,
		Title: "Dijkstra's algorithm (cheapest route through a weighted graph)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"pagerank walked a graph's link structure to find each page's " +
						"steady-state importance, but every edge there counted equally — a link " +
						"was a link, with no notion of 'distance' or 'cost'. Real networks aren't " +
						"like that: a road network's edges have different lengths, a flight " +
						"network's edges have different prices. Say you're standing at city S and " +
						"want the cheapest way to reach every other city. Gut instinct: the road " +
						"S→A, cost 4, is the only edge that touches A directly from S, so surely " +
						"it's A's cheapest route. That instinct is wrong — going the 'long way' " +
						"through B first (S→B costs 1, then B→A costs 2, total 3) is actually " +
						"cheaper than the direct road. A plan that locks in the first route it " +
						"finds to each city, without ever reconsidering it once a shorter detour " +
						"turns up, will happily report the wrong answer. You need a procedure that " +
						"only finalizes a city's cheapest route once it's actually certain — " +
						"nothing left unexplored could possibly beat it.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the 5-city network S, A, B, C, D with roads S-A(4), S-B(1), A-B(2), " +
						"A-C(1), B-C(5), B-D(8), C-D(3). Track a tentative distance for every " +
						"city (0 for S, infinity for the rest) and repeat: finalize whichever " +
						"not-yet-finalized city currently has the smallest tentative distance, " +
						"then 'relax' every one of its neighbors — if reaching a neighbor through " +
						"the city just finalized beats its current tentative distance, lower it.",
					"• Finalize S (dist 0). Relax A: 0+4=4. Relax B: 0+1=1.",
					"• Smallest unfinalized is B (dist 1). Finalize B. Relax A: 1+2=3 < 4 — " +
						"lower it. Relax C: 1+5=6. Relax D: 1+8=9.",
					"• Smallest unfinalized is A (dist 3, exactly the S→B→A route the gut " +
						"instinct in section 1 missed). Finalize A. Relax C: 3+1=4 < 6 — lower it " +
						"again.",
					"• Smallest unfinalized is C (dist 4). Finalize C. Relax D: 4+3=7 < 9 — " +
						"lower it.",
					"• Smallest unfinalized is D (dist 7). Finalize D — nothing left to relax.",
					"Final distances from S: S=0, B=1, A=3, C=4, D=7 — and D's true cheapest " +
						"route, S→B→A→C→D, only came together after three separate rounds of " +
						"relaxation lowered it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each city is a box on the map; roads are the lines between them, labeled " +
						"with cost. The step slider scrubs through the trace above one " +
						"finalization at a time: green boxes are finalized, orange is the city " +
						"just finalized this step, blue is reachable-but-not-yet-finalized (its " +
						"number is only a tentative distance, still able to drop), and gray hasn't " +
						"been reached at all yet. The number under each box is its current " +
						"distance from S. The target slider picks a destination city; once that " +
						"city has a known distance, the road segments of its current best route " +
						"back to S are drawn in a heavier highlight color — watch that highlighted " +
						"route re-route itself through C once C gets finalized.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Find the guaranteed-cheapest route (and its exact cost) from one source to " +
						"every other node in a weighted graph, in a single pass — not just a " +
						"distance number, either: by remembering which finalized neighbor produced " +
						"each city's current-best distance, you can walk that chain of neighbors " +
						"backward from any target all the way to S and read off the actual turn-by" +
						"-turn route, not only how long it is.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"GPS and mapping apps run this (or a close, faster variant) to turn a road " +
						"network into turn-by-turn directions. Internet routers use the same idea " +
						"(OSPF's link-state routing) to find the lowest-latency path for a packet " +
						"across a network of routers. Flight-search and package-delivery routing " +
						"tools use it with ticket price, or delivery time, as the edge weight " +
						"instead of road distance.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: always finalize the cheapest not-yet-finalized city next, " +
						"and every time you finalize one, relax its neighbors — a city's distance " +
						"can keep dropping right up until the moment it's finalized, exactly what " +
						"happened to A (4→3) and C (6→4) above.",
					"Not like this: assuming this works when a road can have a negative cost. " +
						"Dijkstra's algorithm never revisits a city once it's finalized, trusting " +
						"that nothing unexplored could possibly beat its current distance — a move " +
						"that only holds when every edge weight is non-negative. A single " +
						"negative-weight edge discovered later could undercut an already-finalized " +
						"distance, and Dijkstra would never notice; graphs with negative weights " +
						"need a different algorithm (Bellman-Ford) built to handle that case.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "step", Label: "Step (finalize one city at a time)", Min: 0, Max: 5, Step: 1, Def: 0},
			{Key: "target", Label: "Highlight route to city (0=S,1=A,2=B,3=C,4=D)", Min: 0, Max: 4, Step: 1, Def: 4},
		},
		Render: render,
	})
}

// NodeNames labels the fixed 5-city example graph every Section walks
// through, in the same index order Graph and positions use: S=0, A=1, B=2,
// C=3, D=4.
var NodeNames = []string{"S", "A", "B", "C", "D"}

// Edge is one directed edge in an adjacency list: To is the destination
// node's index, Weight its non-negative cost.
type Edge struct {
	To     int
	Weight float64
}

// Graph is the fixed 5-city road network every Section walks through:
// S-A(4), S-B(1), A-B(2), A-C(1), B-C(5), B-D(8), C-D(3), stored as an
// undirected graph (each road listed as a directed edge in both
// directions).
var Graph = [][]Edge{
	{{1, 4}, {2, 1}},                 // S -> A(4), B(1)
	{{0, 4}, {2, 2}, {3, 1}},         // A -> S(4), B(2), C(1)
	{{0, 1}, {1, 2}, {3, 5}, {4, 8}}, // B -> S(1), A(2), C(5), D(8)
	{{1, 1}, {2, 5}, {4, 3}},         // C -> A(1), B(5), D(3)
	{{2, 8}, {3, 3}},                 // D -> B(8), C(3)
}

// Step is one snapshot of Dijkstra's algorithm: the tentative distance and
// current-best predecessor for every node right after one finalization (or,
// for the first Step, before any finalization at all). Current is the node
// finalized this step (-1 for the initial "start" step). Relaxed lists
// every neighbor whose distance was lowered as a result.
type Step struct {
	Dist        []float64
	Prev        []int
	Visited     []bool
	Current     int
	Relaxed     []int
	Description string
}

func cloneFloats(v []float64) []float64 { return append([]float64(nil), v...) }
func cloneInts(v []int) []int           { return append([]int(nil), v...) }
func cloneBools(v []bool) []bool        { return append([]bool(nil), v...) }

// Dijkstra runs Dijkstra's algorithm on g from src and returns every
// intermediate state, starting with the unexplored start (dist[src]=0,
// everything else infinity, nothing finalized). Each round finalizes the
// not-yet-finalized node with the smallest tentative distance and relaxes
// every one of its neighbors -- lowering a neighbor's distance, and
// recording the finalized node as that neighbor's current-best predecessor,
// whenever routing through the just-finalized node beats what was known
// before. A node with no path from src at all is simply never finalized;
// Dijkstra stops once nothing reachable remains unfinalized.
func Dijkstra(g [][]Edge, src int) []Step {
	n := len(g)
	dist := make([]float64, n)
	prev := make([]int, n)
	visited := make([]bool, n)
	for i := range dist {
		dist[i] = math.Inf(1)
		prev[i] = -1
	}
	dist[src] = 0

	steps := []Step{{
		Dist: cloneFloats(dist), Prev: cloneInts(prev), Visited: cloneBools(visited),
		Current: -1, Description: "Start: dist[src]=0, every other city unreached",
	}}

	for round := 0; round < n; round++ {
		u := -1
		best := math.Inf(1)
		for v := 0; v < n; v++ {
			if !visited[v] && dist[v] < best {
				best = dist[v]
				u = v
			}
		}
		if u == -1 {
			break // everything still unvisited is unreachable from src
		}
		visited[u] = true

		var relaxed []int
		for _, e := range g[u] {
			if visited[e.To] {
				continue
			}
			nd := dist[u] + e.Weight
			if nd < dist[e.To] {
				dist[e.To] = nd
				prev[e.To] = u
				relaxed = append(relaxed, e.To)
			}
		}

		steps = append(steps, Step{
			Dist: cloneFloats(dist), Prev: cloneInts(prev), Visited: cloneBools(visited),
			Current: u, Relaxed: relaxed,
			Description: "Finalize the cheapest unfinalized city and relax its neighbors",
		})
	}
	return steps
}

// ShortestPath walks prev backward from target to src and returns the
// route as a slice of node indices from src to target (inclusive). It
// returns nil if target hasn't been reached (prev[target] == -1 and
// target != src).
func ShortestPath(prev []int, src, target int) []int {
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
	c := viz.New(700, 460, 0, 1, 0, 1)
	return c.String()
}
