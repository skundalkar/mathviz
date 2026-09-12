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
	"sort"

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
					"Placeholder -- filled in once the math and picture exist.",
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

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
