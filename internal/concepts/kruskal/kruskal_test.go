package kruskal

import (
	"reflect"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

// TestSortedEdgesOrder is pinned to the specific weight-sorted order worked
// through in LEARNINGS.md and the concept's own Sections, including the
// stable tie-break between the two weight-1 edges (S-B before A-C, matching
// their original order in Edges).
func TestSortedEdgesOrder(t *testing.T) {
	sorted := SortedEdges(Edges)
	want := []Edge{
		{0, 2, 1}, // S-B
		{1, 3, 1}, // A-C
		{1, 2, 2}, // A-B
		{3, 4, 3}, // C-D
		{0, 1, 4}, // S-A
		{2, 3, 5}, // B-C
		{2, 4, 8}, // B-D
	}
	if !reflect.DeepEqual(sorted, want) {
		t.Errorf("SortedEdges = %v, want %v", sorted, want)
	}
}

// TestUnionFindMergesAndDetectsCycles walks through a few unions by hand to
// check both the true (merged) and false (already connected) return paths.
func TestUnionFindMergesAndDetectsCycles(t *testing.T) {
	uf := NewUnionFind(4)
	if uf.Find(0) == uf.Find(1) {
		t.Fatal("0 and 1 start in different components")
	}
	if !uf.Union(0, 1) {
		t.Error("Union(0,1) on fresh singletons should return true")
	}
	if uf.Find(0) != uf.Find(1) {
		t.Error("0 and 1 should share a component after Union")
	}
	if uf.Union(0, 1) {
		t.Error("Union(0,1) again should return false -- already connected")
	}
	if !uf.Union(1, 2) {
		t.Error("Union(1,2) should return true -- 2 was still its own component")
	}
	if uf.Find(0) != uf.Find(2) {
		t.Error("0 and 2 should now share a component via 1")
	}
	if uf.Find(3) == uf.Find(0) {
		t.Error("3 was never unioned, should still be its own component")
	}
}

// TestKruskalAcceptsFourEdges is pinned to the exact trace worked in
// LEARNINGS.md: the first four sorted edges (S-B, A-C, A-B, C-D) are each
// accepted, giving a spanning tree of 4 edges over 5 nodes.
func TestKruskalAcceptsFourEdges(t *testing.T) {
	steps := Kruskal(Edges, len(NodeNames))
	final := steps[len(steps)-1]
	want := []Edge{
		{0, 2, 1}, // S-B
		{1, 3, 1}, // A-C
		{1, 2, 2}, // A-B
		{3, 4, 3}, // C-D
	}
	if !reflect.DeepEqual(final.MST, want) {
		t.Errorf("final MST = %v, want %v", final.MST, want)
	}
}

// TestKruskalSkipsRemainingThreeEdges is pinned to the concept's specific
// claim that S-A, B-C, and B-D are each skipped as cycle-forming once the
// tree is already fully connected.
func TestKruskalSkipsRemainingThreeEdges(t *testing.T) {
	steps := Kruskal(Edges, len(NodeNames))
	// steps[0] is the start step; steps[1..7] correspond to the 7
	// weight-sorted edges from TestSortedEdgesOrder.
	wantAdded := []bool{true, true, true, true, false, false, false}
	if len(steps)-1 != len(wantAdded) {
		t.Fatalf("len(steps)-1 = %d, want %d", len(steps)-1, len(wantAdded))
	}
	for i, want := range wantAdded {
		if steps[i+1].Added != want {
			t.Errorf("steps[%d].Added = %v, want %v (edge %v)", i+1, steps[i+1].Added, want, SortedEdges(Edges)[i])
		}
	}
}

// TestKruskalTotalWeight is pinned to the MST weight worked in
// LEARNINGS.md: 1 + 1 + 2 + 3 = 7.
func TestKruskalTotalWeight(t *testing.T) {
	steps := Kruskal(Edges, len(NodeNames))
	final := steps[len(steps)-1]
	if got := TotalWeight(final.MST); got != 7 {
		t.Errorf("TotalWeight(final.MST) = %v, want 7", got)
	}
}

// TestKruskalFinalComponentsAllConnected checks that every node ends up in
// the same component once the spanning tree is complete.
func TestKruskalFinalComponentsAllConnected(t *testing.T) {
	steps := Kruskal(Edges, len(NodeNames))
	final := steps[len(steps)-1]
	root := final.Components[0]
	for i, c := range final.Components {
		if c != root {
			t.Errorf("Components[%d] = %d, want %d (everyone connected)", i, c, root)
		}
	}
}

// TestKruskalStartStepHasNoEdgesYet checks the initial snapshot: every node
// is its own singleton component and no edge has been considered.
func TestKruskalStartStepHasNoEdgesYet(t *testing.T) {
	steps := Kruskal(Edges, len(NodeNames))
	start := steps[0]
	if start.EdgeIndex != -1 {
		t.Errorf("start.EdgeIndex = %d, want -1", start.EdgeIndex)
	}
	if len(start.MST) != 0 {
		t.Errorf("start.MST = %v, want empty", start.MST)
	}
	for i, c := range start.Components {
		if c != i {
			t.Errorf("start.Components[%d] = %d, want %d (singleton)", i, c, i)
		}
	}
}

// TestKruskalDisconnectedGraphLeavesGaps checks that a graph which can't be
// fully connected simply ends up with fewer than n-1 accepted edges,
// instead of crashing or fabricating a connection that doesn't exist.
func TestKruskalDisconnectedGraphLeavesGaps(t *testing.T) {
	edges := []Edge{
		{0, 1, 1},
		// node 2 has no edge at all -- can never join the tree
	}
	steps := Kruskal(edges, 3)
	final := steps[len(steps)-1]
	if len(final.MST) != 1 {
		t.Errorf("len(final.MST) = %d, want 1 (only 0-1 connectable)", len(final.MST))
	}
	if final.Components[2] == final.Components[0] {
		t.Error("node 2 has no edge, should never share a component with 0")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("kruskals-mst")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
