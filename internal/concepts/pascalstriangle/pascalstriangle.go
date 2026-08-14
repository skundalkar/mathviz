// Package pascalstriangle visualizes Pascal's triangle: each entry built
// as the sum of the two entries above it, and why that same number also
// answers "how many ways are there to choose k items from n" — the
// binomial coefficient "n choose k."
package pascalstriangle

import (
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

func render(p map[string]float64) string {
	c := viz.New(600, 440, 0, 1, 0, 1)
	return c.String()
}
