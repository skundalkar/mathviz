// Package knapsack visualizes the 0/1 knapsack problem solved by dynamic
// programming: instead of checking every possible subset of items (2^n of
// them), build up a table of the best value achievable for every smaller
// item-count/capacity combination, then reuse those already-solved
// subproblems to answer the full problem in one pass.
package knapsack

import (
	"sort"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "dynamic-programming-knapsack",
		Seq:   101,
		Title: "0/1 knapsack (dynamic programming)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "step", Label: "Step (fill the table one cell at a time)", Min: 0, Max: 18, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// Item is one candidate to pack: a name, an integer Weight (how much
// capacity it uses), and an integer Value (how much it's worth).
type Item struct {
	Name   string
	Weight int
	Value  int
}

// Items is the fixed 3-item example every Section walks through -- picked
// specifically because ranking by value-per-weight (the natural greedy
// instinct) does NOT find the best combination here, only dynamic
// programming's table does. See "Why would you need this?".
var Items = []Item{
	{"A", 1, 6},  // ratio 6.0
	{"B", 2, 10}, // ratio 5.0
	{"C", 3, 12}, // ratio 4.0
}

// Capacity is the fixed knapsack capacity every Section walks through.
const Capacity = 5

// KnapsackTable builds the classic 0/1 knapsack DP table: dp[i][w] is the
// best total value achievable using only the first i items (in the order
// given) with total weight at most w. Row 0 is all zeros (no items to
// choose from yet, so nothing to pack). Each later row reuses the row
// above it -- the already-solved "first i-1 items" subproblem -- either
// leaving the i-th item out entirely (dp[i-1][w], unchanged) or, if it
// fits (its Weight <= w), putting it in and adding its Value to whatever
// the remaining capacity could best achieve without it
// (dp[i-1][w-Weight]+Value). dp[i][w] is simply whichever of those two
// options is larger.
func KnapsackTable(items []Item, capacity int) [][]int {
	n := len(items)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}
	for i := 1; i <= n; i++ {
		it := items[i-1]
		for w := 0; w <= capacity; w++ {
			dp[i][w] = dp[i-1][w] // option 1: leave item i-1 out
			if it.Weight <= w {
				if withIt := dp[i-1][w-it.Weight] + it.Value; withIt > dp[i][w] {
					dp[i][w] = withIt // option 2: take it, beats leaving it out
				}
			}
		}
	}
	return dp
}

// Backtrack walks a finished KnapsackTable backward from the bottom-right
// corner to recover WHICH items the optimal dp[len(items)][capacity] value
// actually uses: at each row i, dp[i][w] came from taking item i-1 exactly
// when it differs from dp[i-1][w] (leaving it out would have scored less),
// in which case the search continues from a correspondingly smaller
// capacity. It returns the chosen items' indices in their original order.
func Backtrack(dp [][]int, items []Item, capacity int) []int {
	w := capacity
	var chosen []int
	for i := len(items); i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			chosen = append(chosen, i-1)
			w -= items[i-1].Weight
		}
	}
	for a, b := 0, len(chosen)-1; a < b; a, b = a+1, b-1 {
		chosen[a], chosen[b] = chosen[b], chosen[a]
	}
	return chosen
}

// GreedyByRatio picks items by value-per-weight ratio, richest first,
// taking each one that still fits -- the instinctive approach Section 1
// shows falling short of the true optimum. It returns the chosen items'
// indices in their original order.
func GreedyByRatio(items []Item, capacity int) []int {
	idx := make([]int, len(items))
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(a, b int) bool {
		ra := float64(items[idx[a]].Value) / float64(items[idx[a]].Weight)
		rb := float64(items[idx[b]].Value) / float64(items[idx[b]].Weight)
		return ra > rb
	})
	remaining := capacity
	var chosen []int
	for _, i := range idx {
		if items[i].Weight <= remaining {
			chosen = append(chosen, i)
			remaining -= items[i].Weight
		}
	}
	sort.Ints(chosen)
	return chosen
}

// TotalValue sums the Value of the items at the given indices.
func TotalValue(items []Item, chosen []int) int {
	sum := 0
	for _, i := range chosen {
		sum += items[i].Value
	}
	return sum
}

// TotalWeight sums the Weight of the items at the given indices.
func TotalWeight(items []Item, chosen []int) int {
	sum := 0
	for _, i := range chosen {
		sum += items[i].Weight
	}
	return sum
}

func render(p map[string]float64) string {
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
