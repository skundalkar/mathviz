package astar

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// TestHeuristicMatchesEuclideanDistance is pinned to the coordinates worked
// through in LEARNINGS.md and the concept's own Sections.
func TestHeuristicMatchesEuclideanDistance(t *testing.T) {
	cases := []struct {
		n    int
		want float64
	}{
		{0, 10},            // S: straight line to G is 10
		{1, math.Sqrt(73)}, // A: sqrt((10-2)^2+(0-3)^2) = sqrt(73)
		{2, 6},             // B
		{3, 5},             // C: a 3-4-5 triangle to G
		{4, 2},             // D
		{5, 0},             // G to itself
	}
	for _, c := range cases {
		if got := Heuristic(c.n, Goal); !almostEqual(got, c.want, 1e-9) {
			t.Errorf("Heuristic(%s, G) = %v, want %v", NodeNames[c.n], got, c.want)
		}
	}
}

// TestHeuristicNeverOverestimates checks A*'s admissibility requirement
// directly: for every node, the straight-line distance to the goal must
// never exceed that node's true cheapest graph distance to the goal (found
// here by plain Dijkstra, i.e. Search with useHeuristic=false).
func TestHeuristicNeverOverestimates(t *testing.T) {
	for n := range NodeNames {
		steps := Search(Graph, n, Goal, false)
		trueDist := steps[len(steps)-1].G[Goal]
		h := Heuristic(n, Goal)
		if h > trueDist+1e-9 {
			t.Errorf("Heuristic(%s, G) = %v overestimates true distance %v", NodeNames[n], h, trueDist)
		}
	}
}

// TestSearchAStarExpansionOrder is pinned to the trace worked step by step
// in LEARNINGS.md and the concept's own Sections: A* expands only S, B, D,
// G, in that order, and stops the moment G is expanded.
func TestSearchAStarExpansionOrder(t *testing.T) {
	steps := Search(Graph, 0, Goal, true)
	wantOrder := []int{-1, 0, 2, 4, 5} // step0 has no Current
	if len(steps) != len(wantOrder) {
		t.Fatalf("len(steps) = %d, want %d", len(steps), len(wantOrder))
	}
	for i, want := range wantOrder {
		if steps[i].Current != want {
			t.Errorf("steps[%d].Current = %d, want %d", i, steps[i].Current, want)
		}
	}
}

// TestSearchAStarNeverExpandsAorC checks the concept's specific claim that
// A* never opens the A or C branch on this graph.
func TestSearchAStarNeverExpandsAorC(t *testing.T) {
	steps := Search(Graph, 0, Goal, true)
	final := steps[len(steps)-1]
	if final.Closed[1] {
		t.Error("A* closed A, want it never expanded")
	}
	if final.Closed[3] {
		t.Error("A* closed C, want it never expanded")
	}
}

// TestSearchDijkstraExpandsAllSix is pinned to the comparison trace: plain
// Dijkstra (useHeuristic=false) has to expand every node, S through G, six
// expansions instead of A*'s four, before it reaches the same goal.
func TestSearchDijkstraExpandsAllSix(t *testing.T) {
	steps := Search(Graph, 0, Goal, false)
	wantOrder := []int{-1, 0, 1, 2, 3, 4, 5}
	if len(steps) != len(wantOrder) {
		t.Fatalf("len(steps) = %d, want %d", len(steps), len(wantOrder))
	}
	for i, want := range wantOrder {
		if steps[i].Current != want {
			t.Errorf("steps[%d].Current = %d, want %d", i, steps[i].Current, want)
		}
	}
}

// TestSearchBothAlgorithmsAgreeOnCost checks that A* and plain Dijkstra,
// despite expanding a different number of nodes, arrive at the exact same
// optimal cost for the goal.
func TestSearchBothAlgorithmsAgreeOnCost(t *testing.T) {
	astarFinal := Search(Graph, 0, Goal, true)
	dijkstraFinal := Search(Graph, 0, Goal, false)
	got := astarFinal[len(astarFinal)-1].G[Goal]
	want := dijkstraFinal[len(dijkstraFinal)-1].G[Goal]
	if !almostEqual(got, want, 1e-9) {
		t.Errorf("A* cost = %v, Dijkstra cost = %v, want equal", got, want)
	}
	if !almostEqual(got, 10, 1e-9) {
		t.Errorf("goal cost = %v, want 10 (4.0 + 4.0 + 2.0 along S->B->D->G)", got)
	}
}

// TestPathToKnownRoute is pinned to the S->B->D->G route worked in
// LEARNINGS.md and the concept's own Sections.
func TestPathToKnownRoute(t *testing.T) {
	steps := Search(Graph, 0, Goal, true)
	final := steps[len(steps)-1]
	got := PathTo(final.Prev, 0, Goal)
	want := []int{0, 2, 4, 5} // S, B, D, G
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PathTo(S, G) = %v, want %v", got, want)
	}
}

// TestPathToSourceToItself checks the degenerate single-node path.
func TestPathToSourceToItself(t *testing.T) {
	steps := Search(Graph, 0, Goal, true)
	final := steps[len(steps)-1]
	got := PathTo(final.Prev, 0, 0)
	want := []int{0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PathTo(S, S) = %v, want %v", got, want)
	}
}

// TestPathToUnreachableReturnsNil checks that an unreached target
// (prev == -1) reports no path rather than a bogus partial one.
func TestPathToUnreachableReturnsNil(t *testing.T) {
	prev := []int{-1, -1, -1}
	if got := PathTo(prev, 0, 2); got != nil {
		t.Errorf("PathTo(unreachable) = %v, want nil", got)
	}
}

// TestSearchUnreachableNodeNeverClosed checks that a node with no path from
// src is simply skipped rather than causing a crash or a bogus finite cost.
func TestSearchUnreachableNodeNeverClosed(t *testing.T) {
	g := [][]Edge{
		{{1, 1}}, // 0 -> 1
		{{0, 1}}, // 1 -> 0
		{},       // 2: isolated, unreachable from 0
	}
	steps := Search(g, 0, 1, true)
	final := steps[len(steps)-1]
	if !final.Closed[1] {
		t.Error("goal 1 is reachable from 0, want it closed")
	}
	if final.Closed[2] {
		t.Error("node 2 is unreachable from 0, want it never closed")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("a-star-search")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
