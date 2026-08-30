// Package matinv visualizes the matrix inverse: the matrix A^-1 that
// undoes whatever A does, found by running Gauss-Jordan elimination (full
// row reduction, not just echelon form) on the augmented matrix [A | I]
// until the left half becomes I -- at which point the right half has
// become A^-1. The running example is A = [[4,7],[2,k]], invertible for
// every k except k=3.5, where `determinant`'s det(A)=4k-14 hits zero and
// the matrix collapses the plane onto a line with no way back.
package matinv

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "matrix-inverse",
		Seq:   81,
		Title: "Matrix inverse (undoing a transformation)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`gaussian-elimination` already showed how to solve Ax=b for one specific " +
						"b, by row-reducing the augmented matrix [A | b]. But what if you need to " +
						"solve Ax=b for the same A against ten different b's -- ten different " +
						"batches of sensor readings run through the same transformation, say? " +
						"Rerunning the full elimination procedure from scratch every single time " +
						"works, but it repeats identical work each round: A never changes, only b " +
						"does. Is there a single matrix, call it A^-1, you could compute once so " +
						"that x=A^-1*b directly for any b afterward, without redoing elimination " +
						"every time?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take A = [[4,7],[2,6]]. Augment it with the identity matrix and row-reduce " +
						"the whole thing at once, using the exact same three moves " +
						"`gaussian-elimination` used -- but this time pushing all the way to " +
						"reduced form (every pivot scaled to exactly 1, and cleared out both " +
						"above and below), not just echelon form:",
					"• [4,7 | 1,0] and [2,6 | 0,1]. Scale row 1 by 1/4: [1,1.75 | 0.25,0].",
					"• Eliminate column 1 from row 2 (subtract 2x row 1): [0,2.5 | -0.5,1].",
					"• Scale row 2 by 1/2.5: [0,1 | -0.2,0.4].",
					"• Eliminate column 2 from row 1 (subtract 1.75x row 2): [1,0 | 0.6,-0.7].",
					"The left half has become the identity, so the right half is the answer: " +
						"A^-1 = [[0.6,-0.7],[-0.2,0.4]]. Check it: 4(0.6)+7(-0.2)=1 and " +
						"2(0.6)+6(-0.2)=0, matching the identity's first column exactly.",
					"`determinant` already gave a one-number litmus test for whether this could " +
						"even work: det(A) = 4x6 - 7x2 = 10, nonzero, so an inverse exists -- and " +
						"for a 2x2 matrix it even hands you a shortcut formula, " +
						"A^-1 = (1/det(A)) x [[d,-b],[-c,a]] = (1/10)x[[6,-7],[-2,4]], the same " +
						"answer elimination found. If det(A) were 0 instead, elimination would hit " +
						"a column with no nonzero pivot available anywhere below it -- the same " +
						"stuck spot `gaussian-elimination` hit for a contradictory system -- and " +
						"there would be no A^-1 to find, because a zero-determinant matrix " +
						"collapses the whole plane onto a line, and no matrix can undo throwing " +
						"information away like that.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"k sets A's bottom-right entry (6 in the worked example above); step scrubs " +
						"through the Gauss-Jordan elimination one row operation at a time on the " +
						"augmented matrix [A | I], highlighting the pivot cell and printing the " +
						"operation that produced it. At the final step, with k away from 3.5, the " +
						"left half has become the identity and the right half is A^-1, read off " +
						"directly. Drag k down toward 3.5 and watch det(A)=4k-14 shrink toward " +
						"zero along with it -- at k=3.5 exactly, elimination can't find a pivot in " +
						"column 2 at all, the left half never reaches the identity, and the " +
						"readout reports 'no inverse exists' instead of a matrix.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Compute A^-1 once and reuse it: x=A^-1*b solves Ax=b for any new b with a " +
						"single matrix-vector multiply, instead of rerunning full elimination from " +
						"scratch every time the right-hand side changes. You can also now test " +
						"invertibility directly, before even attempting to solve anything: a " +
						"square matrix has an inverse exactly when its determinant is nonzero, the " +
						"same det(A)!=0 condition `determinant` already introduced for 'this " +
						"transformation doesn't collapse area to zero.'",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Computer graphics converts a point from world coordinates back into an " +
						"object's own local coordinates by multiplying by the object's placement " +
						"matrix's inverse -- undoing the exact rotation, scale, and translation " +
						"that placed it in the world. `linear-regression`'s closed-form least- " +
						"squares solution, beta = (X^T*X)^-1 * X^T*y, uses a matrix inverse " +
						"directly instead of iterating gradient-descent-style. Robotics' inverse " +
						"kinematics asks the same question in reverse: given where a robot arm's " +
						"hand ended up, what joint angles (transformed forward by known matrices) " +
						"produced it -- answered, when the transform is linear and square, by " +
						"inverting it.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'A has an inverse exactly when A is square and " +
						"det(A) != 0' -- both conditions matter, and elimination hitting a missing " +
						"pivot is just this same fact discovered mechanically instead of computed " +
						"up front.",
					"Not like this: assuming every matrix has an inverse the way every nonzero " +
						"number has a reciprocal. A non-square matrix has no inverse at all (that's " +
						"exactly the gap `singular-value-decomposition` fills, generalizing " +
						"'undo this transformation as best you can' to shapes where a true inverse " +
						"can't exist), and a square matrix with det(A)=0 has thrown away " +
						"information no matrix could recover. A second slip: computing A^-1 and " +
						"then multiplying it by b when you only ever need to solve Ax=b once -- " +
						"running elimination on [A | b] directly is both faster and more " +
						"numerically stable; the inverse earns its keep specifically when many " +
						"different b's will reuse the same A.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "k", Label: "k (A's bottom-right entry)", Min: 1, Max: 10, Step: 0.25, Def: 6},
			{Key: "step", Label: "Elimination step", Min: 0, Max: 4, Step: 1, Def: 4},
		},
		Render: render,
	})
}

// tol is the tolerance below which a value is treated as zero when
// deciding whether a pivot exists. Mirrors `gaussian-elimination`'s own
// tol.
const tol = 1e-9

// Matrix is a plain 2D grid of entries, row-major. Used both for a bare
// n x n matrix and for an augmented n x 2n matrix ([A | I] or, after
// elimination, [I | A^-1]).
type Matrix [][]float64

// cloneMatrix returns a deep copy so Step snapshots never alias each other
// or the caller's matrix.
func cloneMatrix(m Matrix) Matrix {
	out := make(Matrix, len(m))
	for i, row := range m {
		out[i] = append([]float64(nil), row...)
	}
	return out
}

// Identity returns the n x n identity matrix.
func Identity(n int) Matrix {
	m := make(Matrix, n)
	for i := range m {
		m[i] = make([]float64, n)
		m[i][i] = 1
	}
	return m
}

// Augment returns [a | Identity(n)], the starting point for Gauss-Jordan
// matrix inversion: n x 2n, a's own entries on the left, the identity on
// the right.
func Augment(a Matrix) Matrix {
	n := len(a)
	out := make(Matrix, n)
	for i, row := range a {
		out[i] = make([]float64, 2*n)
		copy(out[i], row)
		out[i][n+i] = 1
	}
	return out
}

// Step is one snapshot of the augmented matrix, right after one row
// operation (or, for the first Step, before any operation at all), plus a
// human-readable description of what produced it. PivotRow/PivotCol name
// the pivot the operation is working with; both are -1 for the initial
// "start" step.
type Step struct {
	Matrix      Matrix
	Description string
	PivotRow    int
	PivotCol    int
}

// GaussJordan runs full row reduction on an augmented n x 2n matrix (as
// produced by Augment) -- unlike `gaussian-elimination`'s Eliminate, which
// only reaches echelon form, this scales every pivot to exactly 1 and
// clears its column both above and below, so a fully successful run
// leaves the identity in the left n columns and the inverse in the right
// n columns. m is left untouched; every Step holds its own copy. If a
// column has no nonzero entry available in or below the current pivot row,
// that column is left without a pivot and reduction stops there -- the
// matrix has no inverse (see ExtractInverse).
func GaussJordan(m Matrix) []Step {
	if len(m) == 0 {
		return nil
	}
	cur := cloneMatrix(m)
	n := len(m)
	cols := len(m[0])

	steps := []Step{{Matrix: cloneMatrix(cur), Description: "Start", PivotRow: -1, PivotCol: -1}}

	for col := 0; col < n; col++ {
		if math.Abs(cur[col][col]) < tol {
			swapRow := -1
			for r := col + 1; r < n; r++ {
				if math.Abs(cur[r][col]) >= tol {
					swapRow = r
					break
				}
			}
			if swapRow == -1 {
				break // no pivot available anywhere in this column -- not invertible
			}
			cur[col], cur[swapRow] = cur[swapRow], cur[col]
			steps = append(steps, Step{
				Matrix:      cloneMatrix(cur),
				Description: fmt.Sprintf("Swap R%d <-> R%d", col+1, swapRow+1),
				PivotRow:    col,
				PivotCol:    col,
			})
		}

		pivotVal := cur[col][col]
		if math.Abs(pivotVal-1) >= tol {
			for c := 0; c < cols; c++ {
				cur[col][c] /= pivotVal
			}
			steps = append(steps, Step{
				Matrix:      cloneMatrix(cur),
				Description: fmt.Sprintf("R%d <- R%d / %.2f", col+1, col+1, pivotVal),
				PivotRow:    col,
				PivotCol:    col,
			})
		}

		for r := 0; r < n; r++ {
			if r == col || math.Abs(cur[r][col]) < tol {
				continue
			}
			factor := cur[r][col]
			for c := 0; c < cols; c++ {
				cur[r][c] -= factor * cur[col][c]
			}
			steps = append(steps, Step{
				Matrix:      cloneMatrix(cur),
				Description: fmt.Sprintf("R%d <- R%d - (%.2f)*R%d", r+1, r+1, factor, col+1),
				PivotRow:    col,
				PivotCol:    col,
			})
		}
	}
	return steps
}

// ExtractInverse reads the inverse off a fully Gauss-Jordan-reduced n x 2n
// matrix: ok is true exactly when the left n columns have become the
// identity, in which case the right n columns are A^-1.
func ExtractInverse(m Matrix, n int) (inv Matrix, ok bool) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			want := 0.0
			if i == j {
				want = 1
			}
			if math.Abs(m[i][j]-want) >= tol {
				return nil, false
			}
		}
	}
	inv = make(Matrix, n)
	for i := 0; i < n; i++ {
		inv[i] = append([]float64(nil), m[i][n:2*n]...)
	}
	return inv, true
}

// Determinant2x2 returns ad-bc for a 2x2 matrix [[a,b],[c,d]] -- the same
// litmus test `determinant` introduced: nonzero means an inverse exists,
// zero means the matrix collapses the plane onto a line with no way back.
func Determinant2x2(m Matrix) float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

// System builds the worked example's 2x2 matrix, A = [[4,7],[2,k]]. k is
// the concept's slider; det(A) = 4k-14 hits zero (singular, no inverse) at
// k=3.5.
func System(k float64) Matrix {
	return Matrix{
		{4, 7},
		{2, k},
	}
}

// Layout constants for the matrix-grid diagram, in pixels. Mirrors
// `gaussian-elimination`'s own grid layout, sized down for a 2x4 augmented
// matrix instead of 3x4.
const (
	gridX, gridY   = 40.0, 110.0
	cellW, cellH   = 90.0, 60.0
	rhsGap         = 24.0 // extra horizontal gap before the right-half (I / A^-1) columns
	stepListX      = 460.0
	stepListY      = 110.0
	stepLineHeight = 22.0
)

func render(p map[string]float64) string {
	k := p["k"]
	step := int(p["step"])

	a := System(k)
	steps := GaussJordan(Augment(a))
	if step < 0 {
		step = 0
	}
	if step > len(steps)-1 {
		step = len(steps) - 1
	}
	cur := steps[step]
	n := len(a)

	// A canvas we never call Axes()/Sample() on -- every draw call below is
	// in raw pixel space, since the diagram is a table, not a function plot.
	c := viz.New(760, 420, 0, 1, 0, 1)

	c.Text(20, 28, fmt.Sprintf("A = [[4,7],[2,%.2f]]    augmented with I, reducing [A | I] -> [I | A^-1]", k), 14, viz.Ink, "start")
	c.Text(20, 48, fmt.Sprintf("det(A) = 4k-14 = %.2f", Determinant2x2(a)), 13, viz.Muted, "start")
	c.Text(20, 74, fmt.Sprintf("Step %d/%d: %s", step, len(steps)-1, cur.Description), 13, viz.Accent, "start")

	for r, row := range cur.Matrix {
		rowIsPivot := r == cur.PivotRow
		for col, v := range row {
			x := gridX + float64(col)*cellW
			if col == n {
				x += rhsGap // visually separate the right half (I / A^-1)
			}
			y := gridY + float64(r)*cellH

			fill := viz.Faint
			if rowIsPivot {
				fill = viz.Accent
			}
			opacity := 0.35
			if rowIsPivot && col == cur.PivotCol {
				opacity = 0.7
			}
			c.Rect(x, y, cellW-4, cellH-4, fill, opacity)

			textColor := viz.Ink
			if rowIsPivot && col == cur.PivotCol {
				textColor = viz.Warm
			}
			c.Text(x+(cellW-4)/2, y+(cellH-4)/2+5, fmt.Sprintf("%.2f", v), 14, textColor, "middle")
		}
		// A "|" separating A's columns from the identity/inverse columns.
		sepX := gridX + float64(n)*cellW + rhsGap/2
		c.Text(sepX, gridY+float64(r)*cellH+(cellH-4)/2+5, "|", 16, viz.Muted, "middle")
	}

	shown := steps
	truncated := false
	const maxStepLines = 8
	if len(shown) > maxStepLines {
		shown = shown[:maxStepLines]
		truncated = true
	}
	for i, st := range shown {
		y := stepListY + float64(i)*stepLineHeight
		color := viz.Muted
		if i == step {
			color = viz.Accent
		}
		c.Text(stepListX, y, fmt.Sprintf("%d: %s", i, st.Description), 12, color, "start")
	}
	if truncated {
		c.Text(stepListX, stepListY+float64(len(shown))*stepLineHeight, "...", 12, viz.Muted, "start")
	}

	final := steps[len(steps)-1].Matrix
	readoutY := gridY + float64(n)*cellH + 30

	var verdict string
	if inv, ok := ExtractInverse(final, n); ok {
		verdict = fmt.Sprintf("A^-1 = [[%.2f,%.2f],[%.2f,%.2f]]", inv[0][0], inv[0][1], inv[1][0], inv[1][1])
	} else {
		verdict = "no inverse exists -- elimination found no pivot for this column (A is singular)"
	}
	c.Text(20, readoutY, verdict, 13, viz.Ink, "start")

	return c.String()
}
