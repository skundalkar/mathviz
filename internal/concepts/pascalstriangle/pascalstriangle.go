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
				Body:    []string{"placeholder"},
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
