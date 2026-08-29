// Package svm visualizes the support vector machine: among the infinitely
// many straight lines that separate two linearly-separable classes, the one
// that maximizes the gap (margin) to the nearest point on each side. The
// running example is four points -- two per class -- where the two closest
// points across classes turn out to be the only ones that matter: the
// "support vectors" that pin down the boundary.
package svm

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "support-vector-machine",
		Seq:   77,
		Title: "Support vector machine (the maximum-margin boundary)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`logistic-regression` fit one specific S-curve by minimizing log-loss over " +
						"every point, and `decision-trees` picked one specific split by scoring " +
						"information gain -- both land on a single answer because their scoring " +
						"rule prefers it. But once two classes are cleanly separable, there isn't " +
						"just one line that gets 100% right: rotate or shift a separating line a " +
						"little and it can still classify every training point correctly. Two " +
						"students working the same four-point dataset by hand could draw two " +
						"different 'perfect' lines and both be technically right. If accuracy on " +
						"the training points can't tell those lines apart, what can?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take four points: class +1 at (3,3) and (5,4); class -1 at (1,1) and " +
						"(-1,0). Instead of scoring a line by how many points it gets right (all " +
						"of these lines get all four right), score it by its margin: the distance " +
						"from the line to the single nearest point of either class, doubled to " +
						"cover both sides. A line that grazes past a point has almost no buffer " +
						"before a slightly different point of that class would land on the wrong " +
						"side; a line with room on both sides is safer.",
					"The two closest points across the two classes are (3,3) and (1,1) -- " +
						"distance sqrt((3-1)²+(3-1)²) = sqrt(8) ≈ 2.83 apart. The maximum-margin " +
						"line is their perpendicular bisector: the direction straight from (1,1) " +
						"to (3,3) is (2,2), so the line's normal vector (after scaling to length " +
						"1) is n ≈ (0.707, 0.707), passing through their midpoint (2,2). Checking " +
						"the other two points against this line: (5,4) sits 3.54 units out and " +
						"(-1,0) sits 3.54 units out on the correct sides -- both comfortably " +
						"farther from the line than (3,3) and (1,1) are (1.41 units each). Moving " +
						"(5,4) or (-1,0) anywhere in that safe zone never changes the boundary at " +
						"all; only (3,3) and (1,1) -- the 'support vectors' -- decide where the " +
						"line sits, which is exactly the pattern `lagrange-multipliers` pointed at: " +
						"the points with a nonzero multiplier at the solution are the ones that " +
						"actually constrain it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"theta and offset draw your own candidate line (blue) by its angle and " +
						"distance from the origin; the green line and shaded band are the true " +
						"maximum-margin solution from the worked example above, with (3,3) and " +
						"(1,1) circled as its support vectors. The readout reports your line's " +
						"achieved margin -- or flags it as not separating the classes at all, if " +
						"you rotate it too far -- next to the best possible margin (2.83) for " +
						"comparison.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Pick, out of infinitely many training-accuracy-perfect lines, provably the " +
						"one that leaves the most breathing room before a new point near the " +
						"boundary would get misclassified -- and know in advance which handful of " +
						"points (the support vectors) that decision actually depends on, so " +
						"collecting more of the other points wouldn't change the answer at all.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Text classification (spam vs. not-spam from word-count features) and " +
						"bioinformatics (classifying tumor samples from a handful of gene-" +
						"expression readings) were classic SVM strongholds, especially when there " +
						"are more features than examples -- a regime where margin, not just " +
						"training accuracy, matters a lot for how well the boundary generalizes. " +
						"The 'kernel trick' (not shown here) lets the same max-margin idea draw " +
						"curved boundaries by measuring distance in a reshaped space instead of " +
						"the raw one.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the boundary is decided by the closest points across " +
						"classes, so only those points are support vectors' -- moving any other " +
						"point, as long as it stays on its own side of the margin, leaves the " +
						"line exactly where it was.",
					"Not like this: assuming the line that looks centered by eye, or the one a " +
						"different training method (like logistic regression's log-loss fit) " +
						"happens to produce, is automatically the max-margin one -- two lines can " +
						"both classify every point correctly while having very different margins, " +
						"and only measuring the distance to the nearest point tells them apart.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "theta", Label: "your line's angle (theta)", Min: 1, Max: 179, Step: 1, Def: 80, Unit: "°"},
			{Key: "offset", Label: "your line's offset from origin (d)", Min: -2, Max: 5, Step: 0.1, Def: 1.5},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(560, 420, -3, 6, -3, 6).String()
}
