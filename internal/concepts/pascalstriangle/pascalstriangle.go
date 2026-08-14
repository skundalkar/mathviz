// Package pascalstriangle visualizes Pascal's triangle: each entry built
// as the sum of the two entries above it, and why that same number also
// answers "how many ways are there to choose k items from n" — the
// binomial coefficient "n choose k."
package pascalstriangle

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "pascals-triangle",
		Seq:   34,
		Title: "Pascal's triangle (binomial coefficients, row by row)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You're picking a 2-person subcommittee out of a 5-person team, or figuring " +
						"out how many different 3-topping pizzas you can order from 8 available " +
						"toppings. Gut instinct: 'easy, 5 people times 4 remaining choices, that's " +
						"20 ways to pick 2.' That overcounts — picking Alice then Bob lands you the " +
						"same 2-person subcommittee as picking Bob then Alice, but the instinctive " +
						"count treats them as two different outcomes, so the real number of " +
						"distinct groups is smaller than 20. Listing every group by hand to avoid " +
						"that trap works for 2 of 5 people, gets tedious for 3 of 8 toppings, and " +
						"is hopeless for, say, 5 cards out of a 52-card deck. Is there a fast, " +
						"unambiguous way to count 'how many distinct groups of k can I make from n " +
						"things' — one you can look up or build up once, instead of listing or " +
						"re-deriving it by hand every time?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Build a triangle of numbers from the top, one row at a time. Row 0 is just a " +
						"single 1 (there's exactly one way to choose nothing from nothing). Every " +
						"row after that starts and ends with 1 — there's always exactly one way to " +
						"choose none of n things, and exactly one way to choose all of them — and " +
						"every entry in between is simply the sum of the two entries diagonally " +
						"above it.",
					"• Row 0 = [1]. Row 1 = [1, 1]. Row 2 = [1, 2, 1] — the middle 2 is 1+1 from " +
						"row 1. Row 3 = [1, 3, 3, 1] — each 3 is 1+2. Row 4 = [1, 4, 6, 4, 1] — the " +
						"6 is 3+3 from row 3.",
					"• Row 5 = [1, 5, 10, 10, 5, 1] — the first 10 (position k=2) is row 4's 4 " +
						"plus row 4's 6: 4 + 6 = 10.",
					"That 10 also has a second, completely independent way to arrive at the same " +
						"number: the combinatorial formula for 'n choose k,' n!/(k!(n-k)!) = " +
						"5!/(2!×3!) = 120/(2×6) = 10 — exact same answer, reached by pure " +
						"multiplication and division instead of repeated addition. That's not a " +
						"coincidence: every entry the triangle produces by addition always equals " +
						"its row and position's 'n choose k' value, proven equal, not just " +
						"observed to often match.",
					"Apply it: 8 toppings, choose 3 — row 8 of the triangle is " +
						"[1, 8, 28, 56, 70, 56, 28, 8, 1], and position k=3 reads 56. There are 56 " +
						"distinct 3-topping pizzas possible from 8 toppings.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Every row 0 through 8 drawn as a triangle of numbers. The Row n and Position " +
						"k sliders pick one entry, highlighted in blue. Whenever that entry isn't " +
						"on the edge of its row, its two parents one row up are highlighted in " +
						"orange and connected to it with two orange lines, tracing out exactly the " +
						"addition that produced it — move either slider and watch which two " +
						"numbers feed into the newly highlighted one. Below the triangle, the same " +
						"entry is broken down three ways at once: the addition (parent + parent = " +
						"child), the factorial formula, and a plain-English 'ways to choose' " +
						"sentence — all three landing on the same number, since they're three views " +
						"of one fact rather than three separate facts.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Instantly count how many distinct groups of k you can form from n things, " +
						"for numbers far too large to list by hand — using the fast recursive " +
						"picture if you already have the row above, or the direct formula if you " +
						"don't — and trust the two always agree, because they're proven equal, not " +
						"just usually consistent. 8 toppings choose 3 = 56 possible pizzas; 52 " +
						"cards choose 5 = 2,598,960 possible poker hands — numbers nobody is " +
						"realistically listing out one at a time by hand.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Counting problems generally — how many ways to pick a starting five from a " +
						"twelve-player roster, how many different lottery ticket combinations " +
						"exist, how many possible poker hands a deck can deal. It's also the exact " +
						"set of coefficients that appear when you expand (a+b)ⁿ in algebra (row n " +
						"of the triangle gives the coefficients of (a+b)ⁿ's expanded terms), and it " +
						"underlies the binomial probability distribution — the math behind 'what " +
						"are the odds of exactly 6 heads in 10 coin flips' or reading results out " +
						"of an A/B test. 'Pascal's triangle' is also just the name most people " +
						"already know this shape by from a math class, whether or not they " +
						"remember why it works.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'C(n,k) counts groups where order doesn't matter — picking " +
						"Alice then Bob is the same group as picking Bob then Alice.' Not like " +
						"this: confusing it with counting ordered arrangements (permutations), " +
						"which is a bigger number — 5×4=20 ordered pairs from a 5-person team, " +
						"versus C(5,2)=10 unordered pairs, exactly double, because each unordered " +
						"pair corresponds to 2 possible orderings. Also not like this: treating the " +
						"addition rule and the factorial formula as two competing methods that " +
						"might disagree — they're proven to always produce the identical number by " +
						"two different routes, so if a hand calculation ever gives different " +
						"answers from the two methods, the arithmetic has a mistake in it " +
						"somewhere, not the underlying rule.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Row n", Min: 0, Max: 8, Step: 1, Def: 5},
			{Key: "k", Label: "Position k (choose k of n)", Min: 0, Max: 8, Step: 1, Def: 2},
		},
		Render: render,
	})
}

// MaxRows caps how many rows the picture ever draws (rows 0 through
// MaxRows-1), so the triangle and its Params stay in sync.
const MaxRows = 9

// Factorial returns n! (1 for n <= 1). Negative n returns 0 — undefined.
func Factorial(n int) int64 {
	if n < 0 {
		return 0
	}
	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

// BinomialCoefficient returns "n choose k" — the number of ways to pick an
// unordered group of k items out of n — via the combinatorial formula
// n!/(k!(n-k)!). Returns 0 for an out-of-range k (k<0 or k>n).
func BinomialCoefficient(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return Factorial(n) / (Factorial(k) * Factorial(n-k))
}

// PascalTriangle builds the first `rows` rows (row 0 through row
// rows-1) the way the picture does: not from the factorial formula, but
// from Pascal's own addition rule. Every row starts and ends with 1 (there
// is exactly one way to choose none of n items, and exactly one way to
// choose all of them); every entry in between is the sum of the two
// entries above it. rows < 1 is treated as 1.
func PascalTriangle(rows int) [][]int64 {
	if rows < 1 {
		rows = 1
	}
	t := make([][]int64, rows)
	for i := 0; i < rows; i++ {
		t[i] = make([]int64, i+1)
		t[i][0], t[i][i] = 1, 1
		for j := 1; j < i; j++ {
			t[i][j] = t[i-1][j-1] + t[i-1][j]
		}
	}
	return t
}

const (
	cellW, cellH     = 54.0, 42.0
	marginT          = 84.0
	canvasW, canvasH = 560.0, 560.0
	triangleCenterX  = canvasW / 2
	// triangleBottom is the pixel y just past row MaxRows-1's cells (see
	// cellCenter: row MaxRows-1's center is marginT+(MaxRows-1)*cellH+cellH/2,
	// and each cell extends 18px past its center), leaving the formula text
	// below room to sit without overlapping the grid.
	triangleBottom = marginT + (MaxRows-1)*cellH + cellH/2 + 18
)

func render(p map[string]float64) string {
	n := int(p["n"])
	if n < 0 {
		n = 0
	}
	if n > MaxRows-1 {
		n = MaxRows - 1
	}
	k := int(p["k"])
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}

	tri := PascalTriangle(MaxRows)
	target := tri[n][k]

	// The grid is drawn entirely in pixel space via Rect/Text, so the
	// data-space range passed to New is unused -- any placeholder works.
	c := viz.New(canvasW, canvasH, 0, 1, 0, 1)

	cellCenter := func(row, col int) (float64, float64) {
		rowWidth := float64(row+1) * cellW
		xStart := triangleCenterX - rowWidth/2
		return xStart + float64(col)*cellW + cellW/2, marginT + float64(row)*cellH + cellH/2
	}

	// The edges of every row (k=0 or k=n) have nothing to sum -- they're 1
	// by definition, not by addition. Interior entries always have exactly
	// two real parents, one row up.
	interior := n > 0 && k > 0 && k < n

	// Connecting lines from parents down to the target, drawn first so the
	// cells sit on top of them.
	if interior {
		px, py := cellCenter(n-1, k-1)
		tx, ty := cellCenter(n, k)
		c.Path([][2]float64{{px, py}, {tx, ty}}, viz.Warm, 2)
		px2, py2 := cellCenter(n-1, k)
		c.Path([][2]float64{{px2, py2}, {tx, ty}}, viz.Warm, 2)
	}

	for row := 0; row < MaxRows; row++ {
		for col := 0; col <= row; col++ {
			x, y := cellCenter(row, col)
			fill, opacity, textColor := viz.Faint, 1.0, viz.Ink
			switch {
			case row == n && col == k:
				fill, opacity, textColor = viz.Accent, 0.9, "white"
			case interior && row == n-1 && (col == k-1 || col == k):
				fill, opacity, textColor = viz.Warm, 0.35, viz.Ink
			}
			const r = 18.0
			c.Rect(x-r, y-r, 2*r, 2*r, fill, opacity)
			c.Text(x, y+5, fmt.Sprintf("%d", tri[row][col]), 13, textColor, "middle")
		}
	}

	c.Text(20, 24, fmt.Sprintf("Row %d, position %d: C(%d,%d) = %d", n, k, n, k, target), 15, viz.Ink, "start")

	y := triangleBottom + 25.0
	if interior {
		left, right := tri[n-1][k-1], tri[n-1][k]
		c.Text(20, y, fmt.Sprintf("Addition rule: C(%d,%d) = C(%d,%d) + C(%d,%d) = %d + %d = %d",
			n, k, n-1, k-1, n-1, k, left, right, target), 13, viz.Warm, "start")
	} else {
		c.Text(20, y, "Edge entries are always 1 — there's exactly one way to choose none of n items, or all of them.",
			13, viz.Warm, "start")
	}
	y += 22
	c.Text(20, y, fmt.Sprintf("Combinatorial formula: %d! / (%d! × %d!) = %d / (%d × %d) = %d",
		n, k, n-k, Factorial(n), Factorial(k), Factorial(n-k), target), 13, viz.Muted, "start")
	y += 22
	c.Text(20, y, fmt.Sprintf("In words: there are %d different ways to choose %d items out of %d.", target, k, n),
		13, viz.Ink, "start")

	return c.String()
}
