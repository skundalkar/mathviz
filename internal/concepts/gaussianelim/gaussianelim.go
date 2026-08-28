// Package gaussianelim visualizes Gaussian elimination: solving a system of
// linear equations by row-reducing its augmented matrix to echelon form
// using only the three moves that never change the solution (swap two rows,
// scale a row, subtract a multiple of one row from another), then reading
// the system's rank — and whether it has one solution, none, or infinitely
// many — straight off the reduced result.
package gaussianelim

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
