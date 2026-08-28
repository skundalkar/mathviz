// Package gaussianelim visualizes Gaussian elimination: solving a system of
// linear equations by row-reducing its augmented matrix to echelon form
// using only the three moves that never change the solution (swap two rows,
// scale a row, subtract a multiple of one row from another), then reading
// the system's rank — and whether it has one solution, none, or infinitely
// many — straight off the reduced result.
package gaussianelim

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "gaussian-elimination",
		Seq:   74,
		Title: "Gaussian elimination (row-reducing to solve a linear system)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You already know how to solve two equations in two unknowns by " +
						"substitution: solve one equation for x, plug that into the other, solve " +
						"for y, then plug back in for x. That works fine by hand for two " +
						"unknowns. But a system with three equations and three unknowns doesn't " +
						"substitute nearly as cleanly — every substitution drags the other " +
						"variables along with it, and by the time you're solving for the third " +
						"variable the algebra is already a tangle. And even after slogging " +
						"through it, substitution never tells you what to do if the equations " +
						"aren't actually independent — if one of them turns out to be built out " +
						"of the other two, so there's no single point where all three meet. Is " +
						"there a systematic procedure that scales past two variables, and that " +
						"reveals when a system doesn't pin down one unique answer?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the system 2x+y−z=8, −3x−y+2z=−11, −2x+y+2z=−3. Write it as an " +
						"augmented matrix — one row per equation, one column per variable plus " +
						"one for the right-hand side — and eliminate one variable at a time using " +
						"only moves that don't change the solution: swap two rows, scale a row, " +
						"or subtract a multiple of one row from another (each move just recombines " +
						"equations that were already true).",
					"• Eliminate x from row 2: row 2's x-coefficient is −3, row 1's is 2, so " +
						"subtract (−3/2)×row 1 from row 2 — i.e. add 1.5×row 1 — leaving row 2 = " +
						"[0, 0.5, 0.5 | 1].",
					"• Eliminate x from row 3: subtract (−2/2)=−1 times row 1 — i.e. add row 1 — " +
						"leaving row 3 = [0, 2, 1 | 5].",
					"• Eliminate y from row 3 using row 2's pivot (0.5): subtract (2/0.5)=4 " +
						"times row 2 from row 3, leaving row 3 = [0, 0, −1 | 1].",
					"The matrix is now in echelon form — each row's first nonzero entry sits " +
						"strictly right of the row above it. Back-substitute from the bottom: row " +
						"3 says −z=1 so z=−1; row 2 says 0.5y+0.5(−1)=1 so y=3; row 1 says " +
						"2x+3−(−1)=8 so x=2. Three elimination moves and one back-substitution " +
						"pass solved three equations at once — no juggling which variable to " +
						"substitute into which equation next.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"e sets row 3's z-coefficient (2 in the worked example above); step scrubs " +
						"through the elimination one row operation at a time, highlighting the " +
						"pivot row and printing the operation that produced it. Watch what happens " +
						"as you raise e toward 3: row 3's elimination result shrinks toward an " +
						"all-zero coefficient row. At e=3 exactly, row 3 becomes [0,0,0 | 1] — " +
						"'0 = 1', a flat contradiction. The rank of the coefficient part has " +
						"dropped from 3 to 2, but the augmented row still carries a nonzero " +
						"right-hand side, so the system has gone from one unique solution to no " +
						"solution at all, and the readout at the final step says so directly.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Solve systems with any number of variables by the same mechanical " +
						"procedure — swap, scale, subtract — instead of an ad hoc substitution " +
						"chase that gets harder every time a variable is added. And you can now " +
						"diagnose a system instead of just solving it: comparing the rank of the " +
						"coefficient columns to the rank of the full augmented matrix tells you, " +
						"before you even try to back-substitute, whether you're looking at exactly " +
						"one solution (both ranks equal the number of variables), infinitely many " +
						"(both ranks equal but less than the number of variables), or none (the " +
						"ranks disagree, exactly what e=3 above produces).",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Circuit analysis writes one equation per loop from Kirchhoff's voltage and " +
						"current laws and row-reduces them to find every current at once. " +
						"Structural engineers set up a system of force-balance equations at every " +
						"joint of a bridge truss and solve them the same way. Balancing a chemical " +
						"equation is a small linear system in the atom counts. And " +
						"`linear-regression`'s least-squares line is found, under the hood, by " +
						"row-reducing the 'normal equations' — the same elimination machinery, " +
						"just on a matrix built from the data instead of physical constraints.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'row 3 became [0,0,0 | 1], so the coefficient rank is 2 " +
						"but the augmented rank is 3 — the equations contradict each other and " +
						"there's no solution,' naming both ranks rather than just eyeballing the " +
						"zero row.",
					"Not like this: forgetting to apply a row operation to the right-hand-side " +
						"column too — subtracting 4×row 2 from row 3's coefficients but leaving " +
						"row 3's b-value untouched silently produces a wrong system. A second slip: " +
						"seeing an all-zero coefficient row and immediately declaring 'no " +
						"solution' without checking its right-hand side — an all-zero row with a " +
						"zero right-hand side too means infinitely many solutions, not zero.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "e", Label: "e (row 3's z-coefficient)", Min: -1, Max: 5, Step: 0.5, Def: 2},
			{Key: "step", Label: "Elimination step", Min: 0, Max: 3, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// tol is the tolerance below which a value is treated as zero when deciding
// whether a pivot exists or a row is "all zero".
const tol = 1e-9

// Matrix is an augmented matrix: each row holds the coefficient columns
// followed by the right-hand-side value as its last entry.
type Matrix [][]float64

// cloneMatrix returns a deep copy so Step snapshots don't alias each other
// or the caller's matrix.
func cloneMatrix(m Matrix) Matrix {
	out := make(Matrix, len(m))
	for i, row := range m {
		out[i] = append([]float64(nil), row...)
	}
	return out
}

// Step is one snapshot in the elimination process: the matrix state right
// after one row operation (or, for the first Step, before any operation at
// all), plus a human-readable description of what produced it.
// PivotRow/PivotCol name the pivot the operation is working with; both are
// -1 for the initial "start" step, which performed no operation.
type Step struct {
	Matrix      Matrix
	Description string
	PivotRow    int
	PivotCol    int
}

// Eliminate runs forward Gaussian elimination on m (which is left
// untouched — every Step holds its own copy) and returns every intermediate
// state, starting with the unmodified matrix. Each pivot column is handled
// in order: if the current pivot row has a ~zero entry there, the first row
// below it with a nonzero entry is swapped in; then every row below the
// pivot has a multiple of the pivot row subtracted from it to zero out that
// column. A column with no nonzero entry available anywhere below the pivot
// row is left alone (it has no pivot) and elimination moves on to the next
// column without advancing the pivot row -- this is what a rank-deficient
// matrix looks like.
func Eliminate(m Matrix) []Step {
	if len(m) == 0 {
		return nil
	}
	cur := cloneMatrix(m)
	cols := len(m[0])
	rows := len(m)

	steps := []Step{{Matrix: cloneMatrix(cur), Description: "Start", PivotRow: -1, PivotCol: -1}}

	pivotRow := 0
	for col := 0; col < cols-1 && pivotRow < rows; col++ {
		if math.Abs(cur[pivotRow][col]) < tol {
			swapRow := -1
			for r := pivotRow + 1; r < rows; r++ {
				if math.Abs(cur[r][col]) >= tol {
					swapRow = r
					break
				}
			}
			if swapRow == -1 {
				continue // no pivot available in this column; try the next one
			}
			cur[pivotRow], cur[swapRow] = cur[swapRow], cur[pivotRow]
			steps = append(steps, Step{
				Matrix:      cloneMatrix(cur),
				Description: fmt.Sprintf("Swap R%d <-> R%d", pivotRow+1, swapRow+1),
				PivotRow:    pivotRow,
				PivotCol:    col,
			})
		}

		for r := pivotRow + 1; r < rows; r++ {
			if math.Abs(cur[r][col]) < tol {
				continue
			}
			factor := cur[r][col] / cur[pivotRow][col]
			for c := col; c < cols; c++ {
				cur[r][c] -= factor * cur[pivotRow][c]
			}
			steps = append(steps, Step{
				Matrix:      cloneMatrix(cur),
				Description: fmt.Sprintf("R%d <- R%d - (%.2f)*R%d", r+1, r+1, factor, pivotRow+1),
				PivotRow:    pivotRow,
				PivotCol:    col,
			})
		}
		pivotRow++
	}
	return steps
}

// Rank counts how many rows of m have at least one entry larger than tol
// among its first cols columns. Called with cols = number of variables it
// reports the coefficient matrix's rank; called with cols = the full width
// of an augmented matrix it reports the augmented matrix's rank. Comparing
// the two on the same (already eliminated) matrix is the standard test for
// how many solutions a system has.
func Rank(m Matrix, cols int) int {
	rank := 0
	for _, row := range m {
		nonzero := false
		for c := 0; c < cols && c < len(row); c++ {
			if math.Abs(row[c]) >= tol {
				nonzero = true
				break
			}
		}
		if nonzero {
			rank++
		}
	}
	return rank
}

// BackSubstitute solves a square echelon-form augmented matrix (n rows, n+1
// columns, the last column being the right-hand side) for its n unknowns,
// working from the last row up. It reports ok=false the moment it finds a
// ~zero pivot, since that means the system doesn't have exactly one
// solution (use Rank to tell "infinitely many" apart from "none").
func BackSubstitute(m Matrix) (solution []float64, ok bool) {
	n := len(m)
	for _, row := range m {
		if len(row) != n+1 {
			return nil, false // not a square coefficient matrix
		}
	}
	sol := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		if math.Abs(m[i][i]) < tol {
			return nil, false
		}
		sum := m[i][n]
		for j := i + 1; j < n; j++ {
			sum -= m[i][j] * sol[j]
		}
		sol[i] = sum / m[i][i]
	}
	return sol, true
}

// system builds the worked example's augmented matrix: 2x+y-z=8,
// -3x-y+2z=-11, -2x+y+e*z=-3. e is row 3's z-coefficient -- the concept's
// slider. At e=2 the system has the unique solution (2,3,-1); at e=3 the
// coefficient matrix becomes singular and the system is inconsistent.
func system(e float64) Matrix {
	return Matrix{
		{2, 1, -1, 8},
		{-3, -1, 2, -11},
		{-2, 1, e, -3},
	}
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
