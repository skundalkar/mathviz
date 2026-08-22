// Package chisquared visualizes the chi-squared test of independence on a
// 2x2 contingency table: comparing counted outcomes in two groups against
// what you'd expect if the grouping had nothing to do with the outcome, and
// converting the size of that deviation into a p-value.
package chisquared

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "chi-squared-test",
		Seq:   58,
		Title: "Chi-squared test of independence",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You test two versions of a signup page: version A converts 15 of 50 " +
						"visitors (30%), version B converts 25 of 50 (50%). That's a 20-point gap " +
						"in this one batch of visitors — but `p-value` already taught the core " +
						"lesson for a single observed number: any batch would show SOME random " +
						"gap even if both pages were secretly equally good. Is a 20-point gap in " +
						"conversion counts the kind of thing pure luck produces at this sample " +
						"size, or is it large enough that 'both pages are equally good' stops " +
						"being a believable explanation? `p-value` handled one measured statistic " +
						"against a null distribution; chi-squared testing does the same job for " +
						"CATEGORY COUNTS specifically — comparing what you actually counted in " +
						"each group against what you'd expect if the two variables (which page, " +
						"whether they converted) were completely unrelated.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"With 50 visitors in each group and the rates above: observed counts are A: " +
						"15 converted, 35 didn't; B: 25 converted, 25 didn't. If conversion really " +
						"had nothing to do with which page a visitor saw, the OVERALL conversion " +
						"rate — (15+25)/100 = 40% — should apply equally to both groups: expected " +
						"conversions = 50×40% = 20 for each group, expected non-conversions = " +
						"50×60% = 30 for each group.",
					"Compare observed to expected, cell by cell, with Σ(observed−expected)² ÷ " +
						"expected: (15−20)²/20 + (35−30)²/30 + (25−20)²/20 + (25−30)²/30 = 1.25 " +
						"+ 0.833 + 1.25 + 0.833 = 4.17. That total is the chi-squared statistic — " +
						"bigger means the counts deviate more from what independence predicts.",
					"Converting 4.17 into a p-value uses the chi-squared distribution the same " +
						"way `p-value` used a null curve: for a 2×2 table (1 degree of freedom), " +
						"the tail probability has a closed form, P(χ²≥x) = 1 − erf(√(x/2)). " +
						"Plugging in 4.17 gives p ≈ 0.041 — under 5%, the conventional cutoff, so " +
						"'both pages convert equally' starts to look like a poor explanation for " +
						"what was counted.",
					"Now try the exact same 30%-vs-50% gap with only 10 visitors per group " +
						"instead of 50: observed becomes A: 3 converted, 7 didn't; B: 5 converted, " +
						"5 didn't. Same percentages, same 20-point gap — but chi² drops all the " +
						"way to 0.83, and p jumps to about 0.36. The identical-looking gap is no " +
						"longer surprising, because 10 visitors per group is little enough that a " +
						"20-point swing happens by chance all the time.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The n, group-A-rate, and group-B-rate sliders build the observed 2×2 table " +
						"reported at the top — exactly the arithmetic from the previous section, " +
						"recomputed live. Below it, the curve plots the tail probability " +
						"P(χ²≥x) for every possible statistic x from 0 up — a smooth, " +
						"always-decreasing curve from 1 down to 0, since a bigger deviation from " +
						"independence is always less and less likely under pure chance. The " +
						"orange vertical line marks the current chi² value; where it crosses the " +
						"curve IS the p-value, marked with a small dashed readout back to the " +
						"y-axis. The two dashed grey reference lines mark the conventional 0.05 " +
						"and 0.01 significance thresholds (χ²=3.841 and 6.635) so you can see " +
						"directly whether the orange line has crossed into 'usually considered " +
						"significant' territory.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Decide whether two categorical groupings are actually related, or whether " +
						"an apparent gap in counts is small-sample noise, using one number " +
						"instead of eyeballing raw counts — and see directly how much SAMPLE " +
						"SIZE alone changes that verdict, holding the underlying rates fixed. " +
						"That's the same lesson `p-value` and `confidence-interval` taught for a " +
						"single measured statistic, extended to whole categories of counts at " +
						"once — the same chi-squared machinery, with more rows and columns, is " +
						"what powers tests over any number of categories, not just a 2×2 table.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A/B tests on web pages and ads, clinical trials checking whether a " +
						"treatment's outcome rates differ from a control group's, market research " +
						"checking whether product preference differs by age bracket — anywhere " +
						"the question is 'do these category counts line up with what I'd expect " +
						"if there were no relationship.' Geneticists use exactly this test " +
						"(checking observed offspring counts against Mendelian ratios) to see " +
						"whether real data matches a genetic model's predicted proportions.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the 20-point gap gives p≈0.04 at n=50 per group, below " +
						"the conventional 0.05 threshold, so it's unlikely to be pure chance at " +
						"this sample size' — the verdict is always relative to the sample size " +
						"actually collected.",
					"Not like this: 'a big percentage gap always means a real effect.' The exact " +
						"same 30%-vs-50% gap was NOT significant (p≈0.36) at n=10 per group; a " +
						"p-value is a statement about how surprising the counted numbers are, not " +
						"a direct measurement of how large or important the underlying effect is.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Sample size per group (n)", Min: 5, Max: 200, Step: 1, Def: 50},
			{Key: "pA", Label: "Group A success rate", Min: 0, Max: 1, Step: 0.02, Def: 0.30},
			{Key: "pB", Label: "Group B success rate", Min: 0, Max: 1, Step: 0.02, Def: 0.50},
		},
		Render: render,
	})
}

// Table2x2 is a 2x2 contingency table: two rows (groups) by two columns
// (outcomes).
type Table2x2 [2][2]float64

// RowTotal sums one row of the table.
func RowTotal(t Table2x2, row int) float64 {
	return t[row][0] + t[row][1]
}

// ColTotal sums one column of the table.
func ColTotal(t Table2x2, col int) float64 {
	return t[0][col] + t[1][col]
}

// GrandTotal sums every cell of the table.
func GrandTotal(t Table2x2) float64 {
	return RowTotal(t, 0) + RowTotal(t, 1)
}

// Expected returns the table's expected counts under independence: each
// cell is (its row total × its column total) ÷ the grand total -- the count
// that row/column combination would get if the row and column variables
// were completely unrelated, while keeping the same row and column totals
// as the observed table.
func Expected(t Table2x2) Table2x2 {
	n := GrandTotal(t)
	var e Table2x2
	if n == 0 {
		return e
	}
	for r := 0; r < 2; r++ {
		for c := 0; c < 2; c++ {
			e[r][c] = RowTotal(t, r) * ColTotal(t, c) / n
		}
	}
	return e
}

// ChiSquareStat is Pearson's chi-squared statistic Σ(observed−expected)²÷
// expected, summed over every cell. A cell with 0 expected count is skipped
// (it can only happen when a whole row or column is entirely empty).
func ChiSquareStat(observed, expected Table2x2) float64 {
	sum := 0.0
	for r := 0; r < 2; r++ {
		for c := 0; c < 2; c++ {
			if expected[r][c] == 0 {
				continue
			}
			d := observed[r][c] - expected[r][c]
			sum += d * d / expected[r][c]
		}
	}
	return sum
}

// ChiSquareCDFDF1 is the CDF of the chi-squared distribution with 1 degree
// of freedom, P(χ²≤x). A chi-squared(1) variable is the square of a
// standard normal, which gives this a closed form via the error function
// instead of needing a general incomplete-gamma implementation.
func ChiSquareCDFDF1(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Erf(math.Sqrt(x / 2))
}

// PValueDF1 is the upper-tail probability P(χ²≥x) for 1 degree of freedom --
// the p-value of an observed chi-squared statistic from a 2x2 table.
func PValueDF1(chi2 float64) float64 {
	return 1 - ChiSquareCDFDF1(chi2)
}

// ObservedTable builds the concept's 2x2 observed table from the slider
// values: n visitors in each of two groups, with pA and pB the fraction of
// each group that "succeeded" (rounded to the nearest whole count, since a
// real count can't be fractional).
func ObservedTable(n int, pA, pB float64) Table2x2 {
	oA := math.Round(pA * float64(n))
	oB := math.Round(pB * float64(n))
	return Table2x2{
		{oA, float64(n) - oA},
		{oB, float64(n) - oB},
	}
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 460, 0, 1, 0, 1).String()
}
