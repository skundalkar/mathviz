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
					"Say you're packing a knapsack with capacity 5, choosing among three " +
						"items: A (weight 1, value 6), B (weight 2, value 10), C (weight 3, " +
						"value 12). Each item is all-or-nothing -- you either take the whole " +
						"thing or leave it, no splitting an item in half. Gut instinct: rank items " +
						"by value per unit of weight (A: 6/1=6.0, B: 10/2=5.0, C: 12/3=4.0) and " +
						"grab the richest ones first, as long as they still fit. That takes A " +
						"(weight 1, 4 left), then B (weight 2, 2 left) -- C needs weight 3 and " +
						"only 2 remains, so it's skipped. Total value: 6+10=16, with 2 units of " +
						"capacity left unused. That instinct is wrong: leaving A out entirely and " +
						"taking B and C instead uses the full capacity (2+3=5) for value " +
						"10+12=22 -- six points better, just by NOT taking the single best-ratio " +
						"item. With only 3 items you could check every subset by hand (2^3=8 of " +
						"them) to find that out, but that approach doubles in cost with every " +
						"single item added -- 20 items would mean checking over a million " +
						"subsets. Is there a way to find the guaranteed-best combination without " +
						"either a flawed shortcut or checking every possibility one by one?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Build a table: one row per item added so far (row 0 = no items " +
						"available yet), one column per capacity from 0 up to 5. Each cell " +
						"dp[i][w] answers a smaller version of the same question: 'using only " +
						"the first i items, what's the most value you can pack into capacity " +
						"w?' Row 0 is all zeros -- with no items to choose from, nothing can be " +
						"packed no matter the capacity. Every later cell reuses the row above it " +
						"(an already-solved smaller subproblem) instead of starting over: dp[i][w] " +
						"is the better of leaving item i out (dp[i-1][w], unchanged from the row " +
						"above) or, if it fits, taking it (dp[i-1][w-weight]+value, the best the " +
						"REMAINING capacity could do without this item, plus this item's value).",
					"• Row 1 (+A, weight 1, value 6): for every capacity 1 through 5, taking A " +
						"(dp[0][w-1]+6 = 0+6 = 6) beats leaving it out (dp[0][w]=0), so " +
						"dp[1][w]=6 for w=1..5, and dp[1][0]=0 (no room even for A).",
					"• Row 2 (+B, weight 2, value 10): dp[2][2] = max(dp[1][2]=6, " +
						"dp[1][0]+10=10) = 10 -- take B. dp[2][3] = max(dp[1][3]=6, " +
						"dp[1][1]+10=16) = 16 -- take B AND still have A's row-1 value banked at " +
						"capacity 1. dp[2][5] = max(dp[1][5]=6, dp[1][3]+10=16) = 16.",
					"• Row 3 (+C, weight 3, value 12): dp[3][4] = max(dp[2][4]=16, " +
						"dp[2][1]+12=12) = 16 -- leave C out here, it doesn't help yet. dp[3][5] = " +
						"max(dp[2][5]=16, dp[2][2]+12=10+12=22) = 22 -- taking C, on top of " +
						"whatever capacity 2 could best achieve without it (dp[2][2]=10, which is " +
						"B alone), wins.",
					"The bottom-right cell, dp[3][5]=22, is the answer. Walking back from " +
						"there -- did dp[3][5] change from dp[2][5]? Yes (22 vs 16), so C is in; " +
						"drop to capacity 5-3=2. Did dp[2][2] change from dp[1][2]? Yes (10 vs 6), " +
						"so B is in; drop to capacity 2-2=0. Did dp[1][0] change from dp[0][0]? No " +
						"(0 vs 0), so A is out -- recovers the exact combination, {B, C}, not just " +
						"its value.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The grid is the dp table: rows are 'items considered so far', columns " +
						"are capacity 0 through 5. The step slider fills the table left-to-right, " +
						"top-to-bottom, one cell at a time -- orange is the cell just filled, gray " +
						"is already filled, and the description line above explains that cell's " +
						"take-it-or-leave-it decision in the same terms as section 2. Once every " +
						"cell is filled, the backtrack chain -- the specific cells that explain " +
						"where the final answer came from -- lights up green, and a line below " +
						"compares that optimal combination to what ranking by value/weight ratio " +
						"would have picked instead.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Find the provably best combination of all-or-nothing choices under a " +
						"shared limit -- guaranteed optimal, not just a plausible-looking " +
						"shortcut -- by solving it bottom-up from tiny subproblems instead of " +
						"checking every combination by brute force. The table has only " +
						"(items+1)×(capacity+1) cells, so it grows by simple multiplication as " +
						"more items are added, not by doubling the way checking every subset " +
						"does -- the difference between a method that still works with hundreds " +
						"of items and one that stops being practical past about 25.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A shipping company deciding which packages fit on a delivery truck's " +
						"remaining weight allowance, to maximize total delivered value. A cloud " +
						"provider deciding which jobs to run on a server with a fixed memory " +
						"budget. A student with 5 hours left before an exam deciding which " +
						"subset of practice topics (each taking a known amount of time and worth " +
						"a known number of expected points) to study -- exactly this problem, at " +
						"a scale small enough to solve by hand. More generally, dynamic " +
						"programming (reusing solved subproblems instead of recomputing them) is " +
						"the same idea behind edit-distance spell checkers, DNA sequence " +
						"alignment, and shortest-path routing over a schedule of connections.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: dp[i][w] is always the better of 'leave item i out' " +
						"and 'take item i and add its value to whatever the remaining capacity " +
						"already achieved without it' -- both options reuse an already-solved, " +
						"smaller subproblem from the row above; neither one is solved from " +
						"scratch.",
					"Not like this: ranking items by value-per-weight and taking the richest " +
						"ones first. That greedy shortcut is exactly right for the fractional " +
						"version of this problem (where you're allowed to take, say, half of an " +
						"item), but for the ALL-OR-NOTHING version it can miss the true optimum, " +
						"as A/B/C above shows -- the single best-ratio item (A) isn't even in the " +
						"optimal combination.",
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
