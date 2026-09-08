// Package kerneltrick visualizes the kernel trick: when no straight-line
// boundary can separate two classes in their original representation,
// mapping every point into a higher-dimensional feature space can make them
// linearly separable there -- and a kernel function lets support-vector-
// machine's maximum-margin machinery work in that mapped space using only a
// direct formula on the original points, without ever building the mapped
// vectors themselves.
package kerneltrick

import (
	"math"
	"sort"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "kernel-trick",
		Seq:   97,
		Title: "The kernel trick (making the unseparable separable)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"support-vector-machine found the max-margin line that separates two " +
						"classes -- but that whole procedure quietly assumes a straight line " +
						"CAN separate them. What if it can't? Take seven points on a number " +
						"line: x = -3, -2 belong to class B; x = -1, 0, 1 belong to class A; x = " +
						"2, 3 belong to class B again. Class A sits in the middle, class B sits " +
						"on both ends. No single split point on this line can put all of A on " +
						"one side and all of B on the other -- any threshold you pick either cuts " +
						"class A in half or lets some of class B slip onto A's side. " +
						"support-vector-machine's machinery has nothing to grab onto here. Is " +
						"there a way to reuse that same max-margin idea on data that genuinely " +
						"isn't a straight-line problem in the form it's given to you?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Map every point x to two coordinates instead of one: φ(x) = (x, x²) -- " +
						"keep the original value, and add its square as a second dimension. Look " +
						"at just the new second coordinate for each point:",
					"• Class A: x=-1 → x²=1. x=0 → x²=0. x=1 → x²=1.",
					"• Class B: x=-2 → x²=4. x=2 → x²=4. x=-3 → x²=9. x=3 → x²=9.",
					"Every class-A point now has x² ≤ 1, and every class-B point has x² ≥ 4 -- " +
						"a gap opened up in this second coordinate that never existed in x alone. " +
						"Exactly like support-vector-machine's own worked example, take the two " +
						"closest opposite-class points in the mapped space (x²=1 from A, x²=4 " +
						"from B) and split the difference: threshold = (1+4)/2 = 2.5. Classify " +
						"any point by 'is x² below 2.5?' -- a single straight (horizontal) line in " +
						"the 2D mapped space, even though no single point could do it back in 1D.",
					"Now the trick: an SVM's math only ever needs the DOT PRODUCT of two mapped " +
						"points, never the mapped vectors on their own. φ(a)·φ(b) = a·b + a²·b² " +
						"-- and that right-hand side is a formula in a and b directly, no need to " +
						"ever build the pair (a, a²) as an actual vector first. Call that formula " +
						"the kernel, K(a,b) = a·b + a²·b². For a=2, b=-1: K(2,-1) = 2·(-1) + " +
						"2²·(-1)² = -2 + 4 = 2 -- and computing φ(2)·φ(-1) the long way, (2,4)·" +
						"(-1,1) = -2+4 = 2, gives the exact same number.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The view slider switches between two pictures of the same seven points. " +
						"At view=0 they sit on the plain 1D number line, colored by class, and no " +
						"vertical split line can cleanly separate the colors. At view=1 the same " +
						"points are plotted at (x, x²) in 2D; the horizontal margin line at " +
						"y=2.5 now cleanly separates class A (below) from class B (above), and " +
						"the two support points on each side of the gap (x²=1 and x²=4) are " +
						"highlighted. The query slider adds an eighth, movable point: its class " +
						"is decided purely by K(query, support point) arithmetic, shown in the " +
						"readout, and drawn in whichever picture is active.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Classify data that genuinely isn't a straight-line problem in the form " +
						"it's given to you, by reusing support-vector-machine's exact max-margin " +
						"machinery in a mapped space instead of inventing a new algorithm for " +
						"curved boundaries. And because the classifier only ever needs a kernel's " +
						"dot-product number -- never the mapped vector itself -- this scales to " +
						"feature spaces far too large (or, for some kernels, infinite-dimensional) " +
						"to ever build explicitly; you just need a formula for what their dot " +
						"product would have been.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Handwriting and image classifiers used polynomial and RBF-kernel SVMs " +
						"heavily before deep learning became the default, precisely because raw " +
						"pixel brightness is almost never linearly separable by digit or object " +
						"class. Bioinformatics tools classify protein or gene sequences with " +
						"string kernels that compare sequences directly, without ever constructing " +
						"an explicit numeric feature vector for each one. Any 'similarity function " +
						"between two things' that's cheap to compute directly, even when the " +
						"implied feature space would be unwieldy to build, is a candidate kernel.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a kernel is a direct formula for the dot product AFTER " +
						"mapping -- K(a,b) = φ(a)·φ(b) -- computed straight from a and b, with φ " +
						"never actually built as a vector.",
					"Not like this: assuming any similarity-looking function can be used as a " +
						"kernel. A valid kernel has to correspond to an honest dot product in SOME " +
						"feature space (Mercer's condition) -- pick an arbitrary formula that " +
						"doesn't, and there's no consistent φ backing it up, so the max-margin math " +
						"built on top of it can silently stop making sense.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "view", Label: "View (0=original 1D line, 1=mapped 2D space)", Min: 0, Max: 1, Step: 1, Def: 0},
			{Key: "query", Label: "Query point x", Min: -3.5, Max: 3.5, Step: 0.25, Def: 1.5},
		},
		Render: render,
	})
}

// Point is one labeled 1D data point: X is its position on the number
// line, Label is +1 for class A or -1 for class B.
type Point struct {
	X     float64
	Label int
}

// Points is the fixed 7-point example every Section walks through: class A
// (label +1) sits near zero, class B (label -1) sits on both ends, so no
// single threshold on X alone separates them.
var Points = []Point{
	{-3, -1},
	{-2, -1},
	{-1, 1},
	{0, 1},
	{1, 1},
	{2, -1},
	{3, -1},
}

// Phi is the explicit feature map used throughout: keep the original value
// and add its square as a second coordinate.
func Phi(x float64) (float64, float64) {
	return x, x * x
}

// Kernel computes K(a, b) = Phi(a)·Phi(b) -- the dot product the mapped
// space's classifier actually needs -- directly from a and b, without ever
// building the two Phi vectors. See PhiDot for the same number computed the
// long way; TestKernelMatchesExplicitPhiDot checks they always agree.
func Kernel(a, b float64) float64 {
	return a*b + a*a*b*b
}

// PhiDot computes Phi(a)·Phi(b) the long way -- explicitly mapping both
// points first, then dotting the resulting vectors. Exists purely so tests
// (and the concept's Sections) can show Kernel gives the identical answer
// without doing that work.
func PhiDot(a, b float64) float64 {
	a1, a2 := Phi(a)
	b1, b2 := Phi(b)
	return a1*b1 + a2*b2
}

// MarginThreshold finds the maximum-margin split in the mapped space's
// second coordinate (x²), the same "split the difference between the
// nearest opposite-class points" rule support-vector-machine uses: the
// midpoint between the largest x² among class-A points and the smallest x²
// among class-B points. It also returns those two bracketing values, the
// mapped-space support points.
func MarginThreshold(points []Point) (threshold, innerMax, outerMin float64) {
	innerMax = math.Inf(-1)
	outerMin = math.Inf(1)
	for _, pt := range points {
		_, y := Phi(pt.X)
		if pt.Label == 1 && y > innerMax {
			innerMax = y
		}
		if pt.Label == -1 && y < outerMin {
			outerMin = y
		}
	}
	return (innerMax + outerMin) / 2, innerMax, outerMin
}

// Classify decides a point's class in the mapped space: +1 if its mapped
// x² sits below threshold, -1 otherwise.
func Classify(x, threshold float64) int {
	_, y := Phi(x)
	if y < threshold {
		return 1
	}
	return -1
}

// SeparableByThreshold reports whether ANY single threshold on vals (with
// "below the threshold is one class, above is the other", in either
// direction) classifies every corresponding label correctly -- the general
// test for "is this 1D representation linearly separable". Candidates are
// every midpoint between consecutive sorted values, plus one threshold
// below all of them and one above, which is exhaustive: a threshold
// anywhere between two candidate points behaves identically to the
// midpoint between them.
func SeparableByThreshold(vals []float64, labels []int) bool {
	n := len(vals)
	if n == 0 {
		return true
	}
	sorted := append([]float64(nil), vals...)
	sort.Float64s(sorted)

	candidates := make([]float64, 0, n+1)
	candidates = append(candidates, sorted[0]-1)
	for i := 0; i+1 < n; i++ {
		candidates = append(candidates, (sorted[i]+sorted[i+1])/2)
	}
	candidates = append(candidates, sorted[n-1]+1)

	for _, t := range candidates {
		belowIsPositive, aboveIsPositive := true, true
		for i, v := range vals {
			pred := -1
			if v < t {
				pred = 1
			}
			if pred != labels[i] {
				belowIsPositive = false
			}
			if -pred != labels[i] {
				aboveIsPositive = false
			}
		}
		if belowIsPositive || aboveIsPositive {
			return true
		}
	}
	return false
}

func render(p map[string]float64) string {
	c := viz.New(680, 460, 0, 1, 0, 1)
	return c.String()
}
