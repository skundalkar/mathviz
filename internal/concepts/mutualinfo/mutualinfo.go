// Package mutualinfo visualizes mutual information: how many bits of
// uncertainty about one variable get resolved by learning another, built
// from `entropy`'s bits-of-surprise and `kl-divergence`'s "how far is this
// distribution from that one." Unlike `correlation`, it works for
// categorical variables and non-linear relationships, and it's exactly the
// quantity `decision-trees` calls "information gain."
package mutualinfo

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "mutual-information",
		Seq:   67,
		Title: "Mutual information (how much one variable reveals about another)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`correlation` already measures whether two numeric variables move together " +
						"— but it only catches a straight-line relationship. Does someone carrying " +
						"an umbrella tell you anything about whether it's raining? 'Rain' and " +
						"'umbrella' aren't numbers on a line you can compute a slope through, so " +
						"correlation doesn't even apply. And even for numeric variables, " +
						"correlation can completely miss a real, strong relationship that isn't a " +
						"straight line (Y=X² has correlation near 0 despite Y being entirely " +
						"determined by X). `entropy` measures how uncertain one variable is by " +
						"itself. Is there a single number — working for categories, curves, " +
						"anything — for exactly how much learning one variable shrinks your " +
						"uncertainty about the other?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take P(rain)=0.3, and someone's umbrella habit: P(umbrella | rain)=0.9, " +
						"P(umbrella | no rain)=0.1. Multiply out the joint probabilities of every " +
						"rain/umbrella combination:",
					"• P(rain, umbrella) = 0.3×0.9 = 0.27",
					"• P(rain, no umbrella) = 0.3×0.1 = 0.03",
					"• P(no rain, umbrella) = 0.7×0.1 = 0.07",
					"• P(no rain, no umbrella) = 0.7×0.9 = 0.63",
					"Add up the umbrella column: P(umbrella) = 0.27+0.07 = 0.34. Now three " +
						"entropies, all in bits, all from `entropy`'s same -Σp·log2(p): H(rain) = " +
						"0.881 (rain alone), H(umbrella) = 0.925 (umbrella alone), and H(rain, " +
						"umbrella) = 1.350 (both together, computed the same way over all four " +
						"joint probabilities above). If rain and umbrella-carrying were totally " +
						"unrelated, knowing both would cost exactly as many bits as knowing each " +
						"separately: H(rain,umbrella) would equal H(rain)+H(umbrella) = 1.806. It " +
						"doesn't — 1.350 is less. That gap is mutual information: I(X;Y) = " +
						"H(X)+H(Y)-H(X,Y) = 0.881+0.925-1.350 = 0.456 bits, shared between the two " +
						"variables rather than paid for twice. `kl-divergence` gives a second way " +
						"to compute the identical number: I(X;Y) = KL(P(X,Y) ‖ P(X)·P(Y)) — how " +
						"far the real joint distribution sits from the 'as if independent' " +
						"distribution you'd get by just multiplying the marginals. Both routes " +
						"give 0.456 bits, because they're the same quantity derived two ways.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Two mosaic plots, side by side. Each splits a square into a rain column " +
						"(width P(rain)) and a no-rain column, then splits each column vertically " +
						"into umbrella (top) and no-umbrella (bottom) at that column's own " +
						"conditional probability. The left mosaic uses the real conditionals — " +
						"P(umbrella|rain)=0.9 splits the rain column high, P(umbrella|no " +
						"rain)=0.1 splits the no-rain column low, so the boundary jumps sharply " +
						"between columns. The right mosaic draws the 'as if independent' version, " +
						"splitting both columns at the same marginal P(umbrella)=0.34 — a flat, " +
						"unbroken line. The orange reference line marks that same flat height on " +
						"both plots. How far the left mosaic's boundary strays from that line is " +
						"mutual information made visible: drag P(umbrella|rain) and P(umbrella|no " +
						"rain) apart and the jump grows along with I(X;Y); set them equal and the " +
						"left mosaic snaps flat, matching the right one exactly, at I(X;Y)=0.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Measure dependency between any two variables — categorical, numeric, or a " +
						"mix, linear or not — with one number that's always ≥0 and hits exactly 0 " +
						"only when the variables are truly independent. `decision-trees` already " +
						"used this exact quantity without naming it: 'information gain,' the " +
						"criterion each split is chosen to maximize, is I(feature; label) — this " +
						"concept is what was happening under the hood there.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every decision tree split (`decision-trees`) picks the feature and threshold " +
						"that maximizes mutual information with the label, under the name " +
						"'information gain.' Feature selection in machine learning more broadly: " +
						"ranking candidate input variables by mutual information with the target " +
						"catches useful non-linear predictors that a correlation-based ranking " +
						"would score near zero and discard. Biology uses it to find genes whose " +
						"activity moves together in patterns too complex for a straight-line " +
						"correlation to detect; linguistics uses it to find which words tend to " +
						"co-occur more than chance would predict.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'I(X;Y) = H(X)+H(Y)-H(X,Y) = 0.456 bits — symmetric " +
						"(I(X;Y)=I(Y;X), unlike KL divergence's order-sensitive KL(P‖Q)), always " +
						"≥0, and exactly 0 only when the variables are independent.'",
					"Not like this: treating 'correlation near 0' as 'no relationship at all.' Y=X² " +
						"is a textbook case — Y is completely determined by X, yet the correlation " +
						"coefficient comes out near 0 because the relationship curves instead of " +
						"running in a straight line; mutual information between X and Y stays high " +
						"in that same case, because it doesn't assume any particular shape of " +
						"relationship. Also don't treat mutual information as bounded by 1 the way " +
						"a correlation coefficient's magnitude is — I(X;Y) can never exceed " +
						"min(H(X),H(Y)), which depends entirely on how uncertain the two variables " +
						"are to begin with.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "px", Label: "P(rain)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.3},
			{Key: "py1", Label: "P(umbrella | rain)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.9},
			{Key: "py0", Label: "P(umbrella | no rain)", Min: 0.01, Max: 0.99, Step: 0.01, Def: 0.1},
		},
		Render: render,
	})
}

// BinaryEntropy returns the Shannon entropy, in bits, of a two-outcome
// distribution with P(true)=p: H(p) = -p*log2(p) - (1-p)*log2(1-p).
// Identical to entropy.BinaryEntropy -- each concept package is
// self-contained (see BUILD_CYCLE.md).
func BinaryEntropy(p float64) float64 {
	term := func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return -x * math.Log2(x)
	}
	return term(p) + term(1-p)
}

// Joint is the 2x2 joint probability table for two binary variables X
// (e.g. "did it rain") and Y (e.g. "did they carry an umbrella"): P11 =
// P(X=1,Y=1), P10 = P(X=1,Y=0), P01 = P(X=0,Y=1), P00 = P(X=0,Y=0). The
// four entries always sum to 1.
type Joint struct{ P11, P10, P01, P00 float64 }

// NewJoint builds the joint distribution from P(X=1)=px and the two
// conditionals P(Y=1|X=1)=py1, P(Y=1|X=0)=py0.
func NewJoint(px, py1, py0 float64) Joint {
	return Joint{
		P11: px * py1,
		P10: px * (1 - py1),
		P01: (1 - px) * py0,
		P00: (1 - px) * (1 - py0),
	}
}

// MarginalX returns P(X=1).
func (j Joint) MarginalX() float64 { return j.P11 + j.P10 }

// MarginalY returns P(Y=1).
func (j Joint) MarginalY() float64 { return j.P11 + j.P01 }

// JointEntropy returns H(X,Y) in bits: -Σ p(x,y)*log2(p(x,y)) over all
// four joint outcomes (a p==0 cell contributes 0, the same convention
// BinaryEntropy uses).
func (j Joint) JointEntropy() float64 {
	term := func(x float64) float64 {
		if x <= 0 {
			return 0
		}
		return -x * math.Log2(x)
	}
	return term(j.P11) + term(j.P10) + term(j.P01) + term(j.P00)
}

// MutualInformation returns I(X;Y) in bits, computed as
// H(X)+H(Y)-H(X,Y): the bits saved by describing X and Y together instead
// of paying for each variable's uncertainty separately. Never negative
// (H(X,Y) can never exceed H(X)+H(Y)), and exactly 0 only when X and Y are
// independent.
func (j Joint) MutualInformation() float64 {
	return BinaryEntropy(j.MarginalX()) + BinaryEntropy(j.MarginalY()) - j.JointEntropy()
}

// MutualInformationViaKL computes I(X;Y) the other way described in
// LEARNINGS.md: as KL(P(X,Y) ‖ P(X)*P(Y)), summing
// p(x,y)*log2(p(x,y)/(p(x)*p(y))) over all four cells, where p(x)*p(y) is
// the joint distribution X and Y would have if they were independent. This
// should agree with MutualInformation up to floating-point error -- see
// TestMutualInformationMatchesKLFormula.
func (j Joint) MutualInformationViaKL() float64 {
	px, py := j.MarginalX(), j.MarginalY()
	term := func(pxy, qxy float64) float64 {
		if pxy <= 0 {
			return 0
		}
		return pxy * math.Log2(pxy/qxy)
	}
	return term(j.P11, px*py) + term(j.P10, px*(1-py)) +
		term(j.P01, (1-px)*py) + term(j.P00, (1-px)*(1-py))
}

func clamp01(x float64) float64 {
	if x < 0.01 {
		return 0.01
	}
	if x > 0.99 {
		return 0.99
	}
	return x
}

// Layout constants for the two side-by-side mosaic panels, in pixels.
const (
	panelW, panelH = 260.0, 260.0
	panel1X        = 40.0
	panelGap       = 60.0
	panelY         = 160.0
)

func render(p map[string]float64) string {
	px, py1, py0 := clamp01(p["px"]), clamp01(p["py1"]), clamp01(p["py0"])

	j := NewJoint(px, py1, py0)
	py := j.MarginalY()
	hx, hy := BinaryEntropy(j.MarginalX()), BinaryEntropy(py)
	hxy := j.JointEntropy()
	mi := j.MutualInformation()

	c := viz.New(680, 480, 0, 1, 0, 1)

	panel2X := panel1X + panelW + panelGap
	drawMosaic(c, panel1X, panelY, py1, py0, px)
	drawMosaic(c, panel2X, panelY, py, py, px) // "as if independent": both columns split at the marginal

	c.Text(panel1X+panelW/2, panelY-14, "actual joint P(X,Y)", 13, viz.Ink, "middle")
	c.Text(panel2X+panelW/2, panelY-14, "if X, Y independent", 13, viz.Ink, "middle")

	// Orange reference line at the marginal P(Y=1), spanning both panels --
	// any vertical jump between this line and the left mosaic's own column
	// boundaries is mutual information made visible.
	refY := panelY + py*panelH
	c.Rect(panel1X, refY-0.75, panel2X+panelW-panel1X, 1.5, viz.Warm, 0.9)

	c.Text(20, 24, fmt.Sprintf("P(rain)=%.2f   P(umbrella|rain)=%.2f   P(umbrella|no rain)=%.2f", px, py1, py0),
		14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("H(X)=%.3f bits   H(Y)=%.3f bits   H(X,Y)=%.3f bits", hx, hy, hxy),
		13, viz.Muted, "start")
	c.Text(20, 62, fmt.Sprintf("I(X;Y) = H(X)+H(Y)-H(X,Y) = %.3f bits", mi), 15, viz.Good, "start")
	c.Text(20, 80, "blue = P(umbrella)   orange line = marginal P(umbrella), same on both panels",
		12, viz.Muted, "start")

	return c.String()
}

// drawMosaic draws one mosaic plot with its top-left corner at (x0,y0): an
// X=1 column of width px*panelW on the left and an X=0 column of width
// (1-px)*panelW on the right, each split vertically into a Y=1 region
// (top, height = its own column's conditional probability * panelH) and a
// Y=0 region (bottom). Passing the same value for yGivenX1 and yGivenX0
// draws the "as if independent" baseline, where both columns split at
// the same height.
func drawMosaic(c *viz.Canvas, x0, y0, yGivenX1, yGivenX0, px float64) {
	x1w := px * panelW
	x0w := panelW - x1w

	y1h := yGivenX1 * panelH
	c.Rect(x0, y0, x1w, y1h, viz.Accent, 0.8)
	c.Rect(x0, y0+y1h, x1w, panelH-y1h, viz.Faint, 1)

	y0h := yGivenX0 * panelH
	c.Rect(x0+x1w, y0, x0w, y0h, viz.Accent, 0.8)
	c.Rect(x0+x1w, y0+y0h, x0w, panelH-y0h, viz.Faint, 1)

	c.Text(x0+x1w/2, y0+panelH+18, "X=1 (rain)", 12, viz.Muted, "middle")
	c.Text(x0+x1w+x0w/2, y0+panelH+18, "X=0 (no rain)", 12, viz.Muted, "middle")
}
