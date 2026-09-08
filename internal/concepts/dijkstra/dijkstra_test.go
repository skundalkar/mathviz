package dijkstra

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

// TestDijkstraFinalDistances is pinned to the final-distance values worked
// step by step in LEARNINGS.md and the concept's own Sections.
func TestDijkstraFinalDistances(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	final := steps[len(steps)-1]
	want := []float64{0, 3, 1, 4, 7} // S, A, B, C, D
	for i, w := range want {
		if !almostEqual(final.Dist[i], w, 1e-9) {
			t.Errorf("Dist[%s] = %v, want %v", NodeNames[i], final.Dist[i], w)
		}
	}
}

// TestDijkstraFinalizationOrder is pinned to the specific order the
// Sections walk through: S, then B, then A, then C, then D.
func TestDijkstraFinalizationOrder(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	wantOrder := []int{-1, 0, 2, 1, 3, 4} // step0 has no Current
	if len(steps) != len(wantOrder) {
		t.Fatalf("len(steps) = %d, want %d", len(steps), len(wantOrder))
	}
	for i, want := range wantOrder {
		if steps[i].Current != want {
			t.Errorf("steps[%d].Current = %d, want %d", i, steps[i].Current, want)
		}
	}
}

// TestDijkstraRelaxationImprovesDistance checks the concept's specific
// claim that A's distance drops from 4 to 3 (via B) and C's drops from 6
// to 4 (via A) across the trace, not just that the final values are right.
func TestDijkstraRelaxationImprovesDistance(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	// Step 1: only S finalized -> A=4, C=inf still.
	if !almostEqual(steps[1].Dist[1], 4, 1e-9) {
		t.Errorf("after finalizing S, Dist[A] = %v, want 4", steps[1].Dist[1])
	}
	// Step 2: B finalized -> A improves to 3, C becomes 6.
	if !almostEqual(steps[2].Dist[1], 3, 1e-9) {
		t.Errorf("after finalizing B, Dist[A] = %v, want 3", steps[2].Dist[1])
	}
	if !almostEqual(steps[2].Dist[3], 6, 1e-9) {
		t.Errorf("after finalizing B, Dist[C] = %v, want 6", steps[2].Dist[3])
	}
	// Step 3: A finalized -> C improves to 4.
	if !almostEqual(steps[3].Dist[3], 4, 1e-9) {
		t.Errorf("after finalizing A, Dist[C] = %v, want 4", steps[3].Dist[3])
	}
}

// TestDijkstraStartStepAllInfiniteExceptSrc checks the initial snapshot.
func TestDijkstraStartStepAllInfiniteExceptSrc(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	start := steps[0]
	if start.Current != -1 {
		t.Errorf("start.Current = %d, want -1", start.Current)
	}
	if start.Dist[0] != 0 {
		t.Errorf("start.Dist[S] = %v, want 0", start.Dist[0])
	}
	for i := 1; i < len(start.Dist); i++ {
		if !math.IsInf(start.Dist[i], 1) {
			t.Errorf("start.Dist[%s] = %v, want +Inf", NodeNames[i], start.Dist[i])
		}
	}
}

// TestDijkstraUnreachableNodeNeverFinalized checks that a node with no path
// from src is simply skipped rather than causing a crash or a bogus finite
// distance.
func TestDijkstraUnreachableNodeNeverFinalized(t *testing.T) {
	g := [][]Edge{
		{{1, 1}}, // 0 -> 1
		{{0, 1}}, // 1 -> 0
		{},       // 2: isolated, unreachable from 0
	}
	steps := Dijkstra(g, 0)
	final := steps[len(steps)-1]
	if !math.IsInf(final.Dist[2], 1) {
		t.Errorf("Dist[2] = %v, want +Inf (unreachable)", final.Dist[2])
	}
	if final.Visited[2] {
		t.Error("Visited[2] = true, want false (never finalized)")
	}
}

// TestShortestPathKnownRoute is pinned to the S->B->A->C->D route worked in
// LEARNINGS.md and the concept's own Sections.
func TestShortestPathKnownRoute(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	final := steps[len(steps)-1]
	got := ShortestPath(final.Prev, 0, 4) // src=S, target=D
	want := []int{0, 2, 1, 3, 4}          // S, B, A, C, D
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ShortestPath(S, D) = %v, want %v", got, want)
	}
}

// TestShortestPathSourceToItself checks the degenerate single-node path.
func TestShortestPathSourceToItself(t *testing.T) {
	steps := Dijkstra(Graph, 0)
	final := steps[len(steps)-1]
	got := ShortestPath(final.Prev, 0, 0)
	want := []int{0}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ShortestPath(S, S) = %v, want %v", got, want)
	}
}

// TestShortestPathUnreachableReturnsNil checks that an unreached target
// (prev == -1) reports no path rather than a bogus partial one.
func TestShortestPathUnreachableReturnsNil(t *testing.T) {
	prev := []int{-1, -1, -1}
	if got := ShortestPath(prev, 0, 2); got != nil {
		t.Errorf("ShortestPath(unreachable) = %v, want nil", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("dijkstras-algorithm")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
