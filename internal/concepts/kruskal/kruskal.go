// Package kruskal visualizes Kruskal's algorithm: building the cheapest
// possible tree that still connects every node in a weighted graph, by
// sorting every road cheapest-first and greedily taking each one unless it
// would reconnect two cities that a cheaper road already connected --
// which would only create a wasteful cycle, never a shorter tree. It works
// on the exact same 5-city network dijkstras-algorithm walks through, but
// answers a different question: not "what's the cheapest route between two
// cities" but "what's the cheapest set of roads that keeps every city
// reachable from every other one."
package kruskal

import (
	"fmt"
	"sort"
	"strings"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "kruskals-mst",
		Seq:   100,
		Title: "Kruskal's algorithm (cheapest network connecting every node)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"dijkstras-algorithm answers 'what's the cheapest ROUTE between two " +
						"specific cities' -- it only ever needs some of a network's roads, the " +
						"ones on the one route it's tracing. Suppose the question is different: " +
						"you're planning which roads to actually build so every city in the " +
						"5-city network S,A,B,C,D is reachable from every other one, as cheaply as " +
						"possible overall -- not the cheapest route between any one pair, the " +
						"cheapest whole network. Gut instinct: build every road that's individually " +
						"cheap, in order -- S-B (cost 1) and A-C (cost 1) both look like obvious " +
						"first picks, and so does A-B (cost 2) right after. Keep blindly grabbing " +
						"the next-cheapest road and by the time C-D (cost 3) goes in, every city is " +
						"already connected -- so what happens when S-A (cost 4) comes up next? It's " +
						"cheaper than B-D (cost 8), so the same instinct says build it too. But S and " +
						"A are already both reachable from each other (through B and C) -- adding " +
						"S-A wouldn't connect anything new, it would just add a second, redundant " +
						"route between two cities that don't need one. Is there a way to grab roads " +
						"cheapest-first while automatically recognizing the moment a 'cheap' road " +
						"has stopped being useful?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the same 5-city network dijkstras-algorithm uses -- S-A(4), S-B(1), " +
						"A-B(2), A-C(1), B-C(5), B-D(8), C-D(3) -- and sort its 7 roads " +
						"cheapest-first: S-B(1), A-C(1), A-B(2), C-D(3), S-A(4), B-C(5), B-D(8). " +
						"Track which cities are already connected to which (a 'components' " +
						"grouping, starting with every city in its own group of one), and go " +
						"through the sorted list: if a road's two cities are still in different " +
						"groups, build it and merge the groups; if they're already in the same " +
						"group, skip it -- building it could only create a redundant loop.",
					"• S-B(1): S and B are in different groups -- build it. Groups: {S,B}, {A}, " +
						"{C}, {D}.",
					"• A-C(1): A and C are in different groups -- build it. Groups: {S,B}, " +
						"{A,C}, {D}.",
					"• A-B(2): A (in {A,C}) and B (in {S,B}) are in different groups -- build " +
						"it, merging them. Groups: {S,A,B,C}, {D}.",
					"• C-D(3): C (in {S,A,B,C}) and D (in {D}) are in different groups -- build " +
						"it. Groups: {S,A,B,C,D} -- every city is now connected.",
					"• S-A(4): S and A are already in the SAME group -- skip it. Building it " +
						"would only close a loop between two cities that can already reach each " +
						"other.",
					"• B-C(5): same group -- skip.",
					"• B-D(8): same group -- skip.",
					"4 roads got built (S-B, A-C, A-B, C-D) for a total cost of 1+1+2+3=7, " +
						"connecting all 5 cities -- and the 3 roads that got skipped (S-A, B-C, " +
						"B-D) are exactly the ones that would only have added a redundant loop, " +
						"never a new connection.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Cities sit at the same positions dijkstras-algorithm uses; roads are the " +
						"lines between them, labeled with cost. The step slider walks through the " +
						"roads in cheapest-first order one at a time: the road being considered " +
						"this step is drawn in orange, roads already built are green, roads " +
						"already skipped are a dashed red line marked ✕, and roads not reached yet " +
						"stay a plain gray line. Every city is colored by which group it currently " +
						"belongs to -- watch S and B share a color as soon as step 1 builds S-B, A " +
						"and C share a different color after step 2, and by step 4 every city " +
						"shares one color, the visual signal that the network is now fully " +
						"connected.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Build the cheapest possible network that still connects every node, " +
						"instead of just the cheapest route between one pair -- and know for " +
						"certain when a cheap-looking edge is actually redundant (S-A here, at " +
						"cost 4, cheaper than two of the roads that DID get built) rather than " +
						"having to guess or check by hand. dijkstras-algorithm tells you the best " +
						"way from S to any one city; kruskals-mst tells you the fewest, cheapest " +
						"roads needed to keep every city reachable from every other one at once.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Utility companies use minimum spanning trees to plan the cheapest layout " +
						"of power lines, water pipes, or fiber-optic cable that still reaches " +
						"every building in a service area. Network engineers use the same idea to " +
						"design a backbone that connects every office or data center for the least " +
						"total cable cost. It also shows up inside other algorithms -- some image " +
						"segmentation and clustering methods build a minimum spanning tree over " +
						"data points and then cut its most expensive edges to find natural " +
						"groupings.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: sort every edge cheapest-first, and add each one only if " +
						"its two endpoints are still in different components -- the union-find " +
						"structure that tracks components is what makes 'would this close a " +
						"loop' a cheap, certain check instead of a guess.",
					"Not like this: assuming the cheapest edges you accept will always be the " +
						"first ones in sorted order, with no skips in between. S-A(4) here is " +
						"cheaper than B-C(5) and B-D(8), yet it's the one that gets skipped, " +
						"because by the time the algorithm reaches it, S and A are already " +
						"connected through B and C -- being cheap doesn't matter once an edge's " +
						"two endpoints are already reachable from each other; only whether it " +
						"still connects something new does.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "step", Label: "Step (consider one edge at a time)", Min: 0, Max: 7, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// NodeNames labels the fixed 5-city example graph every Section walks
// through, in the same index order Edges uses: S=0, A=1, B=2, C=3, D=4 --
// the identical 5-city network and edge weights dijkstras-algorithm uses.
var NodeNames = []string{"S", "A", "B", "C", "D"}

// Edge is one undirected road between two cities.
type Edge struct {
	U, V   int
	Weight float64
}

// Edges is the fixed 5-city road network every Section walks through:
// S-A(4), S-B(1), A-B(2), A-C(1), B-C(5), B-D(8), C-D(3) -- the same
// network and weights dijkstras-algorithm uses, listed once each (an
// undirected road, not a directed edge in both directions).
var Edges = []Edge{
	{0, 1, 4}, // S-A
	{0, 2, 1}, // S-B
	{1, 2, 2}, // A-B
	{1, 3, 1}, // A-C
	{2, 3, 5}, // B-C
	{2, 4, 8}, // B-D
	{3, 4, 3}, // C-D
}

// TotalWeight sums the weight of every edge, e.g. to report an MST's total
// cost.
func TotalWeight(edges []Edge) float64 {
	sum := 0.0
	for _, e := range edges {
		sum += e.Weight
	}
	return sum
}

// UnionFind (a.k.a. disjoint-set) tracks which nodes currently belong to
// the same connected component, using union-by-rank and path-compressed
// Find so both operations run in near-constant time even on much larger
// graphs than this concept's toy example.
type UnionFind struct {
	parent []int
	rank   []int
}

// NewUnionFind returns a UnionFind over n nodes, each its own singleton
// component.
func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent: parent, rank: make([]int, n)}
}

// Find returns x's component root, compressing the path from x to the root
// along the way so future calls are faster.
func (uf *UnionFind) Find(x int) int {
	for uf.parent[x] != x {
		uf.parent[x] = uf.parent[uf.parent[x]] // path halving
		x = uf.parent[x]
	}
	return x
}

// Union merges the components containing a and b. It returns true if a and
// b were in different components (and are now merged into one) -- the
// signal Kruskal uses to accept an edge -- or false if they were already in
// the same component, meaning this edge would only close a cycle.
func (uf *UnionFind) Union(a, b int) bool {
	ra, rb := uf.Find(a), uf.Find(b)
	if ra == rb {
		return false
	}
	if uf.rank[ra] < uf.rank[rb] {
		ra, rb = rb, ra
	}
	uf.parent[rb] = ra
	if uf.rank[ra] == uf.rank[rb] {
		uf.rank[ra]++
	}
	return true
}

// Components returns, for every node 0..n-1, the current Find() root of the
// component it belongs to.
func (uf *UnionFind) Components() []int {
	out := make([]int, len(uf.parent))
	for i := range out {
		out[i] = uf.Find(i)
	}
	return out
}

// Step is one snapshot of Kruskal's algorithm: the state right after one
// edge from the weight-sorted list has been considered (or, for the first
// Step, before any edge has been considered at all). EdgeIndex is that
// edge's position in the weight-sorted list (-1 for the initial "start"
// step). Added reports whether the edge was accepted into the MST (true)
// or skipped because its endpoints were already connected (false). MST is
// every edge accepted so far, in the order accepted.
type Step struct {
	EdgeIndex   int
	Added       bool
	MST         []Edge
	Components  []int
	Description string
}

// SortedEdges returns a copy of edges sorted cheapest-first. Ties keep
// their original relative order (a stable sort), which is what makes the
// algorithm's trace deterministic and reproducible.
func SortedEdges(edges []Edge) []Edge {
	out := append([]Edge(nil), edges...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Weight < out[j].Weight })
	return out
}

// Kruskal builds a minimum spanning tree over n nodes from edges. It sorts
// every edge cheapest-first (SortedEdges) and considers them one at a
// time, accepting an edge into the MST exactly when its two endpoints are
// still in different union-find components (Union returns true) -- adding
// an edge between two nodes already connected by cheaper accepted edges
// could only create a cycle, never shorten the tree, so Union rejecting it
// is exactly the check that skips it.
//
// A real implementation would stop the moment the MST has n-1 edges, since
// every node is provably already connected by then and no further edge
// could ever be accepted. This walkthrough deliberately keeps considering
// every remaining edge instead, so the trace shows the skip-because-cycle
// case happening for real rather than only asserting that it would.
func Kruskal(edges []Edge, n int) []Step {
	sorted := SortedEdges(edges)
	uf := NewUnionFind(n)

	steps := []Step{{
		EdgeIndex: -1, Components: uf.Components(),
		Description: "Start: sort every road cheapest-first; every city is its own component",
	}}

	var mst []Edge
	for i, e := range sorted {
		added := uf.Union(e.U, e.V)
		if added {
			mst = append(mst, e)
		}
		steps = append(steps, Step{
			EdgeIndex: i, Added: added, MST: append([]Edge(nil), mst...),
			Components:  uf.Components(),
			Description: describe(e, added),
		})
	}
	return steps
}

func describe(e Edge, added bool) string {
	if added {
		return "Different components -- add this road to the network"
	}
	return "Already connected -- adding this road would only close a cycle, skip it"
}

// positions lays the 5-city example graph out on a unit square, in the same
// index order as NodeNames -- the identical layout dijkstras-algorithm uses,
// so the two concepts' pictures of "the same graph" line up.
var positions = [][2]float64{
	{0.08, 0.5}, // S, left
	{0.4, 0.18}, // A
	{0.4, 0.82}, // B
	{0.68, 0.5}, // C
	{0.92, 0.5}, // D, right
}

// groupPalette colors the (at most two, on this graph) currently-merged,
// multi-city components; a city still alone in its own singleton component
// is drawn in viz.Faint instead, since there's nothing yet to distinguish
// it from any other lone city.
var groupPalette = []string{viz.Accent, viz.Warm, viz.Good, viz.Bad}

// groupColors picks one fill color per node from components (a Find() root
// per node, as Step.Components stores it). Every node sharing a component
// gets the same color, chosen from groupPalette by that component's
// lowest-indexed member -- not by the union-find root value itself, which
// can change identity across steps as trees merge -- so a component's color
// stays stable across the trace instead of jumping around as it grows.
// Singleton (not-yet-merged) components are all drawn in viz.Faint.
func groupColors(components []int) []string {
	groups := map[int][]int{}
	for i, root := range components {
		groups[root] = append(groups[root], i)
	}
	colors := make([]string, len(components))
	for _, members := range groups {
		if len(members) < 2 {
			colors[members[0]] = viz.Faint
			continue
		}
		minIdx := members[0]
		for _, m := range members {
			if m < minIdx {
				minIdx = m
			}
		}
		col := groupPalette[minIdx%len(groupPalette)]
		for _, m := range members {
			colors[m] = col
		}
	}
	return colors
}

func render(p map[string]float64) string {
	steps := Kruskal(Edges, len(NodeNames))
	sorted := SortedEdges(Edges)
	maxStep := len(steps) - 1

	step := int(p["step"] + 0.5)
	if step < 0 {
		step = 0
	}
	if step > maxStep {
		step = maxStep
	}
	cur := steps[step]

	inMST := make(map[Edge]bool, len(cur.MST))
	for _, e := range cur.MST {
		inMST[e] = true
	}

	c := viz.New(700, 460, 0, 1, 0, 1)
	colors := groupColors(cur.Components)

	// Roads, drawn first so nodes and highlights sit on top. cur.EdgeIndex
	// is the edge this exact step just considered (-1 on the start step,
	// before anything has been looked at).
	for i, e := range sorted {
		color, width, dash := viz.Muted, 1.5, false
		switch {
		case i == cur.EdgeIndex:
			// The edge this step's Description is about -- orange whether
			// it ends up added or skipped, so it's easy to find on the page.
			color, width = viz.Warm, 3.5
		case inMST[e]:
			color, width = viz.Good, 3
		case i < cur.EdgeIndex:
			// Already considered and skipped in an earlier step.
			color, dash = viz.Bad, true
		default:
			// Not reached yet.
		}
		sx, sy := positions[e.U][0], positions[e.U][1]
		ex, ey := positions[e.V][0], positions[e.V][1]
		c.Path([][2]float64{{sx, sy}, {ex, ey}}, color, width)
		if dash {
			// viz.Path has no native dash option; an "x" marks a skipped
			// road instead, at its midpoint.
			mx, my := (sx+ex)/2, (sy+ey)/2
			c.Text(c.X(mx), c.Y(my)+4, "✕", 14, viz.Bad, "middle")
		}
		lx, ly := sx*0.6+ex*0.4, sy*0.6+ey*0.4
		c.Text(c.X(lx), c.Y(ly)-6, fmt.Sprintf("%.0f", e.Weight), 12, viz.Muted, "middle")
	}

	for i, pos := range positions {
		px, py := c.X(pos[0]), c.Y(pos[1])
		const side = 44.0
		c.Rect(px-side/2, py-side/2, side, side, colors[i], 0.85)
		c.Text(px, py+5, NodeNames[i], 15, "white", "middle")
	}

	if cur.EdgeIndex < 0 {
		c.Text(16, 24, fmt.Sprintf("Step %d/%d: %s", step, maxStep, cur.Description), 13, viz.Ink, "start")
	} else {
		e := sorted[cur.EdgeIndex]
		verdict := "ADD"
		if !cur.Added {
			verdict = "SKIP"
		}
		c.Text(16, 24, fmt.Sprintf("Step %d/%d: consider %s-%s (%.0f) -- %s: %s",
			step, maxStep, NodeNames[e.U], NodeNames[e.V], e.Weight, verdict, cur.Description), 13, viz.Ink, "start")
	}

	names := make([]string, len(cur.MST))
	for i, e := range cur.MST {
		names[i] = fmt.Sprintf("%s-%s", NodeNames[e.U], NodeNames[e.V])
	}
	c.Text(16, 44, fmt.Sprintf("Network so far: %s (total cost %.0f, %d/%d roads)",
		strings.Join(names, ", "), TotalWeight(cur.MST), len(cur.MST), len(NodeNames)-1), 13, viz.Accent, "start")

	c.Text(16, 440, "green=in the network  orange=considering now  red ✕=skipped (cycle)  gray=not reached yet",
		12, viz.Muted, "start")

	return c.String()
}
