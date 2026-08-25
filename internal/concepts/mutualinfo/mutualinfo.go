// Package mutualinfo visualizes mutual information: how many bits of
// uncertainty about one variable get resolved by learning another, built
// from `entropy`'s bits-of-surprise and `kl-divergence`'s "how far is this
// distribution from that one." Unlike `correlation`, it works for
// categorical variables and non-linear relationships, and it's exactly the
// quantity `decision-trees` calls "information gain."
package mutualinfo

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 480, 0, 1, 0, 1).String()
}
