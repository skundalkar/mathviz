// Package matinv visualizes the matrix inverse: the matrix A^-1 that
// undoes whatever A does, found by running Gauss-Jordan elimination (full
// row reduction, not just echelon form) on the augmented matrix [A | I]
// until the left half becomes I -- at which point the right half has
// become A^-1. The running example is A = [[4,7],[2,k]], invertible for
// every k except k=3.5, where `determinant`'s det(A)=4k-14 hits zero and
// the matrix collapses the plane onto a line with no way back.
package matinv

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(760, 420, 0, 1, 0, 1).String()
}
