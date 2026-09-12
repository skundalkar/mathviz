package knapsack

import (
	"reflect"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

// TestKnapsackTableKnownValues is pinned to the exact table worked step by
// step in LEARNINGS.md and the concept's own Sections.
func TestKnapsackTableKnownValues(t *testing.T) {
	dp := KnapsackTable(Items, Capacity)
	want := [][]int{
		{0, 0, 0, 0, 0, 0},     // no items yet
		{0, 6, 6, 6, 6, 6},     // + A (w1,v6)
		{0, 6, 10, 16, 16, 16}, // + B (w2,v10)
		{0, 6, 10, 16, 18, 22}, // + C (w3,v12)
	}
	if !reflect.DeepEqual(dp, want) {
		t.Errorf("KnapsackTable = %v, want %v", dp, want)
	}
}

// TestKnapsackTableRowZeroIsAlwaysZero checks the base case: with no items
// to choose from, no capacity can ever produce a positive value.
func TestKnapsackTableRowZeroIsAlwaysZero(t *testing.T) {
	dp := KnapsackTable(Items, Capacity)
	for w, v := range dp[0] {
		if v != 0 {
			t.Errorf("dp[0][%d] = %d, want 0", w, v)
		}
	}
}

// TestBacktrackFindsOptimalItems is pinned to the concept's specific claim
// that the optimal packing is {B, C} (value 22), not {A, B} (value 16, what
// GreedyByRatio finds instead).
func TestBacktrackFindsOptimalItems(t *testing.T) {
	dp := KnapsackTable(Items, Capacity)
	chosen := Backtrack(dp, Items, Capacity)
	want := []int{1, 2} // B, C
	if !reflect.DeepEqual(chosen, want) {
		t.Errorf("Backtrack = %v, want %v", chosen, want)
	}
	if got := TotalValue(Items, chosen); got != 22 {
		t.Errorf("TotalValue(optimal) = %d, want 22", got)
	}
	if got := TotalWeight(Items, chosen); got != Capacity {
		t.Errorf("TotalWeight(optimal) = %d, want %d (fills the knapsack exactly)", got, Capacity)
	}
}

// TestGreedyByRatioFallsShortOfOptimal is pinned to the concept's central
// claim: ranking by value/weight ratio (A=6.0, B=5.0, C=4.0) picks A then
// B (value 16), missing the true optimum B+C (value 22) that Backtrack
// finds.
func TestGreedyByRatioFallsShortOfOptimal(t *testing.T) {
	greedy := GreedyByRatio(Items, Capacity)
	want := []int{0, 1} // A, B
	if !reflect.DeepEqual(greedy, want) {
		t.Errorf("GreedyByRatio = %v, want %v", greedy, want)
	}
	greedyValue := TotalValue(Items, greedy)
	if greedyValue != 16 {
		t.Errorf("TotalValue(greedy) = %d, want 16", greedyValue)
	}

	dp := KnapsackTable(Items, Capacity)
	optimalValue := TotalValue(Items, Backtrack(dp, Items, Capacity))
	if greedyValue >= optimalValue {
		t.Errorf("greedy value %d should be strictly worse than optimal %d", greedyValue, optimalValue)
	}
}

// TestKnapsackTableFinalCellMatchesDirectComputation is a redundant, more
// literal check on the single most important cell: dp[len(items)][Capacity]
// must equal the best value achievable, 22.
func TestKnapsackTableFinalCellMatchesDirectComputation(t *testing.T) {
	dp := KnapsackTable(Items, Capacity)
	got := dp[len(Items)][Capacity]
	if got != 22 {
		t.Errorf("dp[%d][%d] = %d, want 22", len(Items), Capacity, got)
	}
}

// TestKnapsackTableZeroCapacity checks the degenerate case: a knapsack with
// no room at all can never hold anything, regardless of how many items are
// available.
func TestKnapsackTableZeroCapacity(t *testing.T) {
	dp := KnapsackTable(Items, 0)
	for i, row := range dp {
		if row[0] != 0 {
			t.Errorf("dp[%d][0] = %d, want 0", i, row[0])
		}
	}
}

// TestBacktrackNoItemsChosenWhenNoneFit checks that an item heavier than
// the whole capacity is correctly never chosen.
func TestBacktrackNoItemsChosenWhenNoneFit(t *testing.T) {
	items := []Item{{"Heavy", 100, 999}}
	dp := KnapsackTable(items, 5)
	chosen := Backtrack(dp, items, 5)
	if len(chosen) != 0 {
		t.Errorf("Backtrack = %v, want empty (item doesn't fit)", chosen)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("dynamic-programming-knapsack")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
