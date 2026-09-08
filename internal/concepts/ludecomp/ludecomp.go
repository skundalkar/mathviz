// Package ludecomp visualizes LU decomposition: factoring a square matrix A
// into a unit-lower-triangular L and an upper-triangular U such that A =
// L*U, by recording the multipliers gaussian-elimination normally discards.
// Once L and U are known, solving A x = b for any b becomes two cheap
// triangular substitutions (forward through L, back through U) instead of
// repeating the full O(n^3) elimination from scratch for every new b.
package ludecomp

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "lu-decomposition",
		Seq:   95,
		Title: "LU decomposition (factor once, solve many)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"gaussian-elimination solves A x = b by bolting b onto A as one augmented " +
						"matrix and row-reducing the whole thing together. That works, but notice " +
						"what it throws away: every row operation elimination performs — which " +
						"multiple of which row to subtract from which other row — depends only on " +
						"A, never on b. Now suppose you need to solve the same system again with a " +
						"different b: a bridge's force-balance equations under a second load " +
						"scenario, or a circuit's equations for a second input signal. " +
						"gaussian-elimination has no memory of the work it already did — it " +
						"augments the new b onto A and repeats the entire elimination from " +
						"scratch, redoing the O(n^3) row-reduction work even though A itself, and " +
						"therefore every multiplier elimination computes, hasn't changed at all. " +
						"Is there a way to do A's elimination work exactly once, and reuse it " +
						"cheaply for every new b?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take A = [[4,3,2],[2,3,1],[1,1,2]]. Run the exact same forward elimination " +
						"gaussian-elimination does, but this time keep the multiplier used at each " +
						"step instead of discarding it once the row is updated:",
					"• Zero A's column 0 below row 0: row 1's multiplier is 2/4=0.5, row 2's is " +
						"1/4=0.25. Store these in a matrix L at positions L[1][0] and L[2][0]; " +
						"apply them to update rows 1 and 2, same as gaussian-elimination always did.",
					"• Zero column 1 below row 1 (now [0, 1.5, 0]): row 2's multiplier is " +
						"0.25/1.5=1/6. Store it at L[2][1]; apply it, leaving row 2 = [0, 0, 1.5].",
					"• L gets 1s on its diagonal (a row is never eliminated using itself), and U " +
						"is exactly the echelon form elimination produced: U = [[4,3,2],[0,1.5,0]," +
						"[0,0,1.5]], L = [[1,0,0],[0.5,1,0],[0.25,1/6,1]].",
					"Multiply L by U back out and you get A again exactly — L*U = A. That's the " +
						"whole factorization: L records *how* elimination combined the rows, U " +
						"records *what* it produced. Now for any b, substitute A x = (L U) x = b, " +
						"and let y = U x: first solve L y = b for y (forward substitution, top " +
						"row down — cheap, since L is triangular and each row only needs the y " +
						"values already found above it), then solve U x = y for x (back " +
						"substitution, bottom row up). For b = [9, 7, 6]: forward gives y = " +
						"[9, 2.5, 3.333]; back-substituting through U gives x ≈ [−0.111, 1.667, " +
						"2.222] — two O(n^2) triangular solves, no re-elimination of A at all.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The b slider picks which right-hand side to solve for; the step slider " +
						"scrubs through the process in order: A sits unfactored, then L and U " +
						"appear (the one-time O(n^3) factoring work), then y fills in row by row " +
						"as forward substitution walks down through L, then x fills in row by row " +
						"as back substitution walks back up through U. Switching b never changes " +
						"L or U — only the forward/back substitution numbers change — which is the " +
						"whole point made visible.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Solve A x = b for as many different b vectors as you need, each one costing " +
						"only two cheap triangular substitutions instead of a full re-elimination. " +
						"Solving with b set to each of the identity's columns (e0, e1, e2) one at " +
						"a time, and stacking the resulting x's side by side, gives you the columns " +
						"of A⁻¹ — the same inverse matrix-inverse computes by row-reducing [A|I], " +
						"just reusing one L/U factorization instead of row-reducing from scratch. " +
						"And because L's diagonal is all 1s, det(A) = det(U) = the product of U's " +
						"diagonal entries — 4×1.5×1.5 = 9 here — a byproduct of the factorization " +
						"determinant would otherwise compute separately.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Structural engineers solving the same bridge or truss equations under many " +
						"different load scenarios factor the stiffness matrix once and re-solve " +
						"for each load with cheap substitutions. Circuit simulators do the same " +
						"across many input signals. Numerical libraries like LAPACK use LU " +
						"decomposition as the default way to solve linear systems and compute " +
						"determinants and inverses, precisely because factoring once and reusing " +
						"it is so much cheaper than re-deriving from scratch every time.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'factor A into L and U once, then solve L y = b and " +
						"U x = y for each new b' — the factorization is reused, only the two " +
						"triangular solves repeat.",
					"Not like this: assuming this simple no-pivoting version (called Doolittle's " +
						"method) always works. If elimination ever lands on a ~zero pivot — U's " +
						"diagonal entry at that step — it can't proceed without first swapping in " +
						"a row that has a nonzero entry there, exactly the row-swap move " +
						"gaussian-elimination already uses. The fix, used by real numerical " +
						"libraries, is 'LU with partial pivoting' (PA = LU for some permutation " +
						"matrix P) rather than assuming a plain factorization always exists.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "b", Label: "Right-hand side b (0=[9,7,6], 1..3=unit vectors e0..e2)", Min: 0, Max: 3, Step: 1, Def: 0},
			{Key: "step", Label: "Solve step", Min: 0, Max: 7, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// tol is the tolerance below which a pivot is treated as zero -- this
// simple (no-pivoting) Doolittle factorization can't proceed past a pivot
// that small, exactly the failure mode the concept's "common mistake"
// section calls out.
const tol = 1e-9

// A is the fixed 3x3 example every Section walks through: it has three
// clean, nonzero pivots under plain Doolittle elimination, so no row
// swapping is ever needed.
var A = Matrix{
	{4, 3, 2},
	{2, 3, 1},
	{1, 1, 2},
}

// Matrix is a square matrix stored row-major, dense.
type Matrix [][]float64

// cloneMatrix returns a deep copy so Decompose never mutates its input.
func cloneMatrix(m Matrix) Matrix {
	out := make(Matrix, len(m))
	for i, row := range m {
		out[i] = append([]float64(nil), row...)
	}
	return out
}

// identity returns the n x n identity matrix.
func identity(n int) Matrix {
	m := make(Matrix, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

// Decompose factors square matrix a into a unit-lower-triangular L and an
// upper-triangular U such that L*U = a, using Doolittle's method: run
// ordinary forward Gaussian elimination, but instead of discarding the
// multiplier used to zero out each entry, record it into L. U ends up
// holding exactly the echelon form elimination produces. It reports
// ok=false the moment a pivot (U's diagonal entry at that step) would be
// ~zero, since this simple version has no way to swap in a better row --
// see gaussianelim.Eliminate for the row-swap gaussian-elimination itself
// relies on to handle that case.
func Decompose(a Matrix) (L, U Matrix, ok bool) {
	n := len(a)
	L = identity(n)
	U = make(Matrix, n)
	for i := range U {
		U[i] = make([]float64, n)
	}

	work := cloneMatrix(a)
	for i := 0; i < n; i++ {
		for k := i; k < n; k++ {
			U[i][k] = work[i][k]
		}
		if math.Abs(U[i][i]) < tol {
			return nil, nil, false
		}
		for r := i + 1; r < n; r++ {
			factor := work[r][i] / U[i][i]
			L[r][i] = factor
			for k := i; k < n; k++ {
				work[r][k] -= factor * U[i][k]
			}
		}
	}
	return L, U, true
}

// ForwardSubstitute solves L y = b for y, where L is unit lower triangular
// (as Decompose produces): each row's diagonal entry is always 1, so no
// division is needed, only subtracting off the contribution of the y
// values already found above it.
func ForwardSubstitute(L Matrix, b []float64) []float64 {
	n := len(L)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		sum := b[i]
		for j := 0; j < i; j++ {
			sum -= L[i][j] * y[j]
		}
		y[i] = sum
	}
	return y
}

// BackSubstitute solves U x = y for x, where U is upper triangular (as
// Decompose produces): works from the last row up, subtracting off the
// contribution of the x values already found below it before dividing by
// the diagonal pivot.
func BackSubstitute(U Matrix, y []float64) []float64 {
	n := len(U)
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sum := y[i]
		for j := i + 1; j < n; j++ {
			sum -= U[i][j] * x[j]
		}
		x[i] = sum / U[i][i]
	}
	return x
}

// Solve factors a once via Decompose and solves A x = b via forward then
// back substitution, reporting ok=false if a couldn't be factored.
func Solve(a Matrix, b []float64) (x []float64, ok bool) {
	L, U, ok := Decompose(a)
	if !ok {
		return nil, false
	}
	y := ForwardSubstitute(L, b)
	return BackSubstitute(U, y), true
}

// Determinant computes det(a) as the product of U's diagonal entries after
// LU decomposition -- valid because det(L)=1 (unit diagonal, triangular)
// and det of a triangular matrix is the product of its diagonal, so
// det(a) = det(L)*det(U) = det(U). Reports ok=false if a couldn't be
// factored.
func Determinant(a Matrix) (det float64, ok bool) {
	_, U, ok := Decompose(a)
	if !ok {
		return 0, false
	}
	det = 1
	for i := range U {
		det *= U[i][i]
	}
	return det, true
}

// bVecs holds the four right-hand sides the b slider picks between: a
// general system, then the identity's three columns, whose solutions are
// A⁻¹'s columns (see "What can you do now that you couldn't before?").
var bVecs = [][]float64{
	{9, 7, 6},
	{1, 0, 0},
	{0, 1, 0},
	{0, 0, 1},
}

var bNames = []string{
	"b = [9, 7, 6]",
	"b = e0 = [1, 0, 0]  (solving recovers column 0 of A⁻¹)",
	"b = e1 = [0, 1, 0]  (solving recovers column 1 of A⁻¹)",
	"b = e2 = [0, 0, 1]  (solving recovers column 2 of A⁻¹)",
}

// Layout constants for render's matrix grids and vector columns, in pixels.
const (
	cellW, cellH = 58.0, 38.0
	gridAx       = 24.0
	gridLx       = gridAx + 3*cellW + 40
	gridUx       = gridLx + 3*cellW + 40
	gridY        = 100.0

	vecY    = gridY + 3*cellH + 56
	vecBx   = 24.0
	vecYx   = vecBx + 190
	vecXx   = vecYx + 190
	vecRowH = 26.0
)

// drawMatrix draws an n x n matrix as a grid of cells at (x0, gridY). When
// revealed is false every cell shows "?" instead of its value, dimmed --
// used for L and U before the factorization step.
func drawMatrix(c *viz.Canvas, m Matrix, x0 float64, revealed bool, label string) {
	c.Text(x0, gridY-14, label, 13, viz.Muted, "start")
	for r, row := range m {
		for col, v := range row {
			x := x0 + float64(col)*cellW
			y := gridY + float64(r)*cellH
			opacity := 0.3
			if revealed {
				opacity = 0.6
			}
			c.Rect(x, y, cellW-4, cellH-4, viz.Faint, opacity)
			text, color := "?", viz.Muted
			if revealed {
				text, color = fmt.Sprintf("%.2f", v), viz.Ink
			}
			c.Text(x+(cellW-4)/2, y+(cellH-4)/2+5, text, 12, color, "middle")
		}
	}
}

// drawVector draws one labeled column of an n-entry vector, one row per
// entry, revealing only the rows revealedUpTo allows and highlighting
// activeRow (the entry the current step is computing) in a different color
// -- revealedUpTo(i) true means row i already shows its value.
func drawVector(c *viz.Canvas, vals []float64, x0 float64, symbol string, revealedUpTo func(i int) bool, activeRow int) {
	c.Text(x0, vecY-14, symbol, 13, viz.Muted, "start")
	for i, v := range vals {
		y := vecY + float64(i)*vecRowH
		text, color := "?", viz.Muted
		switch {
		case i == activeRow:
			color = viz.Warm
			if revealedUpTo(i) {
				text = fmt.Sprintf("%.3f", v)
			}
		case revealedUpTo(i):
			text, color = fmt.Sprintf("%.3f", v), viz.Ink
		}
		c.Text(x0, y, fmt.Sprintf("%s[%d] = %s", symbol, i, text), 13, color, "start")
	}
}

func render(p map[string]float64) string {
	n := len(A)
	bIdx := int(p["b"] + 0.5)
	if bIdx < 0 {
		bIdx = 0
	}
	if bIdx > len(bVecs)-1 {
		bIdx = len(bVecs) - 1
	}
	b := bVecs[bIdx]

	// Total steps: 1 (start) + 1 (factored) + n (forward) + n (back).
	maxStep := 1 + 2*n
	step := int(p["step"] + 0.5)
	if step < 0 {
		step = 0
	}
	if step > maxStep {
		step = maxStep
	}

	L, U, _ := Decompose(A) // A is fixed and known to factor cleanly.
	y := ForwardSubstitute(L, b)
	x := BackSubstitute(U, y)
	det, _ := Determinant(A)

	factored := step >= 1
	forwardRevealed := step - 1
	if forwardRevealed < 0 {
		forwardRevealed = 0
	}
	if forwardRevealed > n {
		forwardRevealed = n
	}
	backRevealed := step - 1 - n
	if backRevealed < 0 {
		backRevealed = 0
	}
	if backRevealed > n {
		backRevealed = n
	}
	forwardActive, backActive := -1, -1
	if step >= 2 && step <= 1+n {
		forwardActive = step - 2
	}
	if step >= 2+n && step <= 1+2*n {
		backActive = n - 1 - (step - (2 + n))
	}

	c := viz.New(760, 460, 0, 1, 0, 1)

	c.Text(16, 26, fmt.Sprintf("A = [[4,3,2],[2,3,1],[1,1,2]]      %s", bNames[bIdx]), 13, viz.Ink, "start")

	var stepDesc string
	switch {
	case step == 0:
		stepDesc = "Start: A given, not yet factored -- every b would require full elimination again."
	case step == 1:
		stepDesc = fmt.Sprintf("Factored: A = L*U (Doolittle elimination; multipliers recorded into L). det(A) = product(diag(U)) = %.2f*%.2f*%.2f = %.2f.",
			U[0][0], U[1][1], U[2][2], det)
	case forwardActive >= 0:
		row := forwardActive
		terms := ""
		for j := 0; j < row; j++ {
			terms += fmt.Sprintf(" - L[%d][%d]*y[%d](%.2f*%.3f)", row, j, j, L[row][j], y[j])
		}
		stepDesc = fmt.Sprintf("Forward substitution: y[%d] = b[%d]%s = %.3f", row, row, terms, y[row])
	case backActive >= 0:
		row := backActive
		terms := ""
		for j := row + 1; j < n; j++ {
			terms += fmt.Sprintf(" - U[%d][%d]*x[%d](%.2f*%.3f)", row, j, j, U[row][j], x[j])
		}
		stepDesc = fmt.Sprintf("Back substitution: x[%d] = (y[%d]%s) / U[%d][%d](%.2f) = %.3f",
			row, row, terms, row, row, U[row][row], x[row])
	}
	c.Text(16, 48, stepDesc, 13, viz.Accent, "start")

	drawMatrix(c, A, gridAx, true, "A")
	drawMatrix(c, L, gridLx, factored, "L")
	drawMatrix(c, U, gridUx, factored, "U")

	drawVector(c, b, vecBx, "b", func(i int) bool { return true }, -1)
	drawVector(c, y, vecYx, "y", func(i int) bool { return i < forwardRevealed }, forwardActive)
	drawVector(c, x, vecXx, "x", func(i int) bool { return i >= n-backRevealed }, backActive)

	return c.String()
}
