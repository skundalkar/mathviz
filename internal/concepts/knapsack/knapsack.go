// Package knapsack visualizes the 0/1 knapsack problem solved by dynamic
// programming: instead of checking every possible subset of items (2^n of
// them), build up a table of the best value achievable for every smaller
// item-count/capacity combination, then reuse those already-solved
// subproblems to answer the full problem in one pass.
package knapsack

import (
	"fmt"
	"sort"
	"strings"

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

// backtrackPath returns the sequence of table cells Backtrack visits while
// walking from the finished answer back to the base case, purely so render
// can highlight that exact chain once the table is complete.
func backtrackPath(dp [][]int, items []Item) [][2]int {
	w := Capacity
	path := [][2]int{{len(items), w}}
	for i := len(items); i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			w -= items[i-1].Weight
		}
		path = append(path, [2]int{i - 1, w})
	}
	return path
}

func render(p map[string]float64) string {
	dp := KnapsackTable(Items, Capacity)
	n, capW := len(Items), Capacity
	maxStep := n * (capW + 1)

	step := int(p["step"] + 0.5)
	if step < 0 {
		step = 0
	}
	if step > maxStep {
		step = maxStep
	}

	onPath := make(map[[2]int]bool)
	if step == maxStep {
		for _, cell := range backtrackPath(dp, Items) {
			onPath[cell] = true
		}
	}

	c := viz.New(700, 460, 0, 1, 0, 1)

	const originX, originY, cellW, cellH, gap = 170.0, 78.0, 78.0, 46.0, 6.0

	// Column headers (capacity 0..Capacity) and row headers (item names).
	for w := 0; w <= capW; w++ {
		x := originX + float64(w)*(cellW+gap) + cellW/2
		c.Text(x, originY-14, fmt.Sprintf("w=%d", w), 12, viz.Muted, "middle")
	}
	c.Text(originX-16, originY+cellH/2+5, "no items", 12, viz.Muted, "end")
	for i, it := range Items {
		y := originY + float64(i+1)*(cellH+gap) + cellH/2 + 5
		c.Text(originX-16, y, fmt.Sprintf("+%s (w%d,v%d)", it.Name, it.Weight, it.Value), 12, viz.Muted, "end")
	}

	for i := 0; i <= n; i++ {
		for w := 0; w <= capW; w++ {
			filled := i == 0
			flatIdx := -1
			if i > 0 {
				flatIdx = (i-1)*(capW+1) + w
				filled = flatIdx < step
			}
			isCurrent := flatIdx == step-1

			x := originX + float64(w)*(cellW+gap)
			y := originY + float64(i)*(cellH+gap)

			color, opacity := "white", 1.0
			switch {
			case isCurrent:
				color, opacity = viz.Warm, 0.85
			case onPath[[2]int{i, w}]:
				color, opacity = viz.Good, 0.55
			case filled:
				color, opacity = viz.Faint, 1.0
			default:
				opacity = 0
			}
			c.Rect(x, y, cellW-2, cellH-2, color, opacity)
			if filled || isCurrent {
				c.Text(x+cellW/2-1, y+cellH/2+5, fmt.Sprintf("%d", dp[i][w]), 14, viz.Ink, "middle")
			}
		}
	}

	if step == 0 {
		c.Text(16, 24, "Step 0: no items considered yet -- every cell in row 0 is 0", 13, viz.Ink, "start")
	} else if step < maxStep {
		i := (step-1)/(capW+1) + 1
		w := (step - 1) % (capW + 1)
		it := Items[i-1]
		took := dp[i][w] != dp[i-1][w]
		verdict := "leave it out (doesn't improve on the row above)"
		if took {
			verdict = fmt.Sprintf("take it: dp[%d][%d]+%d beats the row above", i-1, w-it.Weight, it.Value)
		}
		c.Text(16, 24, fmt.Sprintf("Step %d/%d: dp[%d][%d], considering %s (w%d,v%d) -- %s",
			step, maxStep, i, w, it.Name, it.Weight, it.Value, verdict), 13, viz.Ink, "start")
	} else {
		chosen := Backtrack(dp, Items, Capacity)
		names := make([]string, len(chosen))
		for i, idx := range chosen {
			names[i] = Items[idx].Name
		}
		greedy := GreedyByRatio(Items, Capacity)
		greedyNames := make([]string, len(greedy))
		for i, idx := range greedy {
			greedyNames[i] = Items[idx].Name
		}
		c.Text(16, 24, fmt.Sprintf("Table complete. Optimal: {%s} = value %d (weight %d/%d) -- highlighted in green",
			strings.Join(names, ","), TotalValue(Items, chosen), TotalWeight(Items, chosen), Capacity), 13, viz.Accent, "start")
		c.Text(16, 44, fmt.Sprintf("Greedy by value/weight would have picked {%s} = value %d instead",
			strings.Join(greedyNames, ","), TotalValue(Items, greedy)), 13, viz.Bad, "start")
	}

	c.Text(16, 440, "orange=cell just filled  gray=already filled  green=on the optimal backtrack path",
		12, viz.Muted, "start")

	return c.String()
}
