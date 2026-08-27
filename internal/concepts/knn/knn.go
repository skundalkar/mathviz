// Package knn visualizes k-nearest-neighbors classification: instead of
// fitting an explicit decision boundary in advance (the single threshold
// decision-trees picks, or the coefficients logistic-regression fits),
// classify a new point by simply looking up its k closest labeled examples
// and taking a majority vote among them.
package knn

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "k-nearest-neighbors",
		Seq:   72,
		Title: "k-nearest neighbors (classify by asking your neighbors)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"decision-trees classifies a new student by comparing one number — hours " +
						"studied — against a single threshold. That works when one feature really " +
						"does separate the classes with one clean cut. But what if two features " +
						"matter together — say, hours studied AND hours slept — and the boundary " +
						"between pass and fail isn't a straight cut on either one alone, just a " +
						"rough 'similar situations tend to turn out similarly'? Fitting an explicit " +
						"formula for that boundary in advance is one option (that's what " +
						"logistic-regression does) — but is there a way to classify a new case " +
						"just by looking at what happened to the most similar past cases, without " +
						"fitting any formula for the boundary at all?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take 10 past students, each plotted by (hours studied, hours slept), " +
						"labeled pass or fail: 5 failed — (1,2), (2,1), (2,4), (3,2), and an " +
						"outlier (8,2) who studied a lot but slept almost none and still failed. " +
						"5 passed — (7,8), (9,7), (8,9), (9,9), and an outlier (2,8) who barely " +
						"studied but slept a lot and still passed. A new student sits at (4,6). " +
						"Measure the plain straight-line (Euclidean) distance from (4,6) to every " +
						"one of the 10, and sort nearest first:",
					"• Closest: (2,4) fail at 2.828, and (2,8) pass, also at 2.828 — an exact " +
						"tie. With only the single nearest neighbor to ask (k=1), a tie has to be " +
						"broken somehow; breaking it toward whichever point came first in the data " +
						"gives k=1 a verdict of FAIL.",
					"• Next closest: (7,8) pass at 3.606, then (3,2) fail at 4.123, then (1,2) " +
						"fail and (8,9) pass, both at 5.0.",
					"• k=1: 1 fail vote → FAIL. k=3: 1 fail, 2 pass → PASS. k=5: 3 fail, 2 pass " +
						"→ FAIL. k=7: 3 fail, 4 pass → PASS. k=9: 5 fail, 4 pass → FAIL.",
					"The verdict flips four times as k grows from 1 to 9 — not because the " +
						"student changed, but because each new k pulls in one more vote from " +
						"farther away, and this particular query sits almost exactly on the " +
						"boundary between the two clusters.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Red squares are past fails, green squares are past passes, sized the same " +
						"for every point. The orange diamond is the query point — drag its qx/qy " +
						"sliders to move it anywhere on the plane. A thin line connects the query " +
						"to each of its current k nearest neighbors (set by the k slider), and the " +
						"verdict banner reads off the majority vote among exactly those connected " +
						"points, live. Leave the query at its default (4,6) and step k from 1 to 9 " +
						"to watch which lines light up change, and the verdict flip back and forth " +
						"exactly as section 2 walked through.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Classify a new point using two (or more) features at once, with a decision " +
						"boundary that can bend around clusters however they're actually shaped — " +
						"no formula to fit, no threshold to choose in advance. The tradeoff is " +
						"right there in the flipping verdict above: a small k reacts to the single " +
						"nearest point, including a lucky or unlucky outlier (high variance, low " +
						"bias); a large k averages over more of the space and is steadier, but can " +
						"blur past a real local boundary or, taken to its limit (k = every point), " +
						"collapse to always predicting whichever class has more total examples, " +
						"ignoring location entirely.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Recommendation systems suggest items liked by 'users similar to you.' Image " +
						"and handwriting recognition can classify a new image by finding the most " +
						"visually similar labeled images. Anomaly detection flags a data point as " +
						"suspicious when its nearest neighbors are all unusually far away. It's " +
						"also a common baseline in medical diagnosis support tools — classify a " +
						"new patient's risk by looking at outcomes for the most similar past " +
						"patients on record.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: k-nearest-neighbors classifies by majority vote among the " +
						"k closest labeled points, measured by distance in feature space — there's " +
						"no fitted formula, only the training data itself and a distance rule. Not " +
						"like this: assuming a bigger k is always 'more accurate' — push k too far " +
						"(all the way to every training point, as section 2's k=... all case shows) " +
						"and the model stops looking at the query's location at all, just returning " +
						"the overall class balance. Also not like this: forgetting that features on " +
						"very different scales (say, dollars vs. years) will silently let the " +
						"large-scale feature dominate the distance calculation — a proper " +
						"implementation typically standardizes features first (see z-score) so " +
						"'closeness' means the same thing along every axis.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "qx", Label: "Query hours studied", Min: 0, Max: 10, Step: 0.5, Def: 4},
			{Key: "qy", Label: "Query hours slept", Min: 0, Max: 10, Step: 0.5, Def: 6},
			{Key: "k", Label: "k (neighbors)", Min: 1, Max: 9, Step: 2, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	return viz.New(680, 460, 0, 10, 0, 10).String()
}
