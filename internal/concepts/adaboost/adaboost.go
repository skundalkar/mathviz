// Package adaboost visualizes AdaBoost: instead of scoring a single
// `decision-trees`-style split once and living with whatever it gets wrong,
// fit a sequence of one-split stumps, each on a version of the data where
// every point misclassified so far has been reweighted heavier -- so each
// new stump is aimed specifically at whatever the stumps before it still
// get wrong -- then combine every stump's vote, weighted by how accurate it
// was, into one stronger classifier.
package adaboost

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "adaboost",
		Seq:   98,
		Title: "AdaBoost (reweighting the hardest examples each round)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`decision-trees` scored every candidate split by information gain and " +
						"picked the single best one -- but a single stump is only ever one split: " +
						"it draws one line and lives with whatever it gets wrong forever. Take " +
						"eight points on a line, alternating in pairs by class: x=1,2 and x=5,6 " +
						"are class +1; x=3,4 and x=7,8 are class -1. Gut instinct: reuse " +
						"`decision-trees`'s recipe, find the single best split, and use it. The " +
						"best single split here (at x=2.5, predicting +1 to its left and -1 to " +
						"its right) only gets 6 of the 8 points right -- it separates x=1,2 from " +
						"the rest, but nothing about that same line can also carve out x=5,6 " +
						"(also +1) from x=3,4,7,8 (-1). No single stump, however cleverly placed, " +
						"can ever separate this alternating pattern perfectly. Is there a way to " +
						"combine several deliberately imperfect, one-line classifiers into " +
						"something that gets all eight right?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Start every point with equal weight, 1/8 = 0.125. Fit a stump the same " +
						"way `decision-trees` does, but score each candidate split by weighted " +
						"error instead of a plain count. Round 1's best split is threshold=2.5 " +
						"(+1 for x≤2.5, -1 for x>2.5); it's wrong on x=5 and x=6 (both true +1, " +
						"predicted -1), so weighted error e1 = 2×0.125 = 0.25. That error earns " +
						"the stump a vote weight, alpha1 = 0.5·ln((1-e1)/e1) = 0.5·ln(3) ≈ 0.549 " +
						"-- the smaller a stump's error, the larger its vote. Then reweight: every " +
						"point this stump got wrong grows heavier, every point it got right " +
						"shrinks. Worked out exactly, a misclassified point's new weight is " +
						"old/(2·e1) = 0.125/0.5 = 0.25, and a correctly-classified point's new " +
						"weight is old/(2·(1-e1)) = 0.125/1.5 ≈ 0.083 -- x=5 and x=6 now carry " +
						"0.25 each, double their starting weight, while the other six shrink to " +
						"about 0.083.",
					"Round 2 fits a new stump to those reweighted points, and with x=5,x=6 now " +
						"dominant, the best split shifts to threshold=6.5 (+1 for x≤6.5, -1 for " +
						"x>6.5) -- it fixes x=5,x=6 but now misses x=3,x=4 instead (weighted " +
						"error e2 = 2×0.083 ≈ 0.167, alpha2 = 0.5·ln(5) ≈ 0.805). Reweighting " +
						"again pushes x=3,x=4 up to 0.25 each. Round 3 fits a third stump to " +
						"that: threshold=4.5 (-1 for x≤4.5, +1 for x>4.5), which by itself is " +
						"wrong on x=1,x=2,x=7,x=8 (e3=0.2, alpha3 = 0.5·ln(4) ≈ 0.693) -- taken " +
						"alone, this third stump is only 60% right.",
					"None of the three stumps is individually correct on all 8 points -- the " +
						"best any one manages alone is 75%. But sum each stump's vote weighted " +
						"by its own alpha (alpha1·h1(x) + alpha2·h2(x) + alpha3·h3(x)) and take " +
						"the sign of the total: every single one of the 8 points comes out " +
						"correct. Three specifically, differently wrong classifiers cancel out " +
						"each other's blind spots.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The rounds slider adds stumps one at a time (0 = no classifier yet, every " +
						"weight still uniform). Each active round's threshold is a dashed " +
						"vertical line, colored to match that round's readout line below. Every " +
						"point is a small square (accent = true class +1, warm = true class -1), " +
						"sized by its current weight -- watch x=5,x=6 balloon after round 1, then " +
						"shrink back down as x=3,x=4 balloon after round 2 -- ringed green if the " +
						"combined vote of every active round gets that point right, red if it's " +
						"still wrong. At rounds=1, two red rings remain (x=5,x=6); at rounds=2, " +
						"two different points turn red (x=3,x=4); at rounds=3, every ring turns " +
						"green.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Build a strong classifier out of several classifiers that are each " +
						"individually mediocre -- none of the three stumps above tops 75% alone " +
						"-- by having every new one specifically target whatever the ensemble so " +
						"far still gets wrong, and letting the more accurate rounds (smaller " +
						"error → larger alpha) carry more of the final vote than the shakier " +
						"ones. `decision-trees` showed how to score and pick one split; AdaBoost " +
						"turns 'fit one split' into a loop that keeps refitting a split to " +
						"whatever the vote so far still leaves unsolved.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Viola-Jones (2001), the first face detector fast enough to run in real " +
						"time, was exactly this: thousands of simple 'is this patch of pixels " +
						"brighter than that one' stumps, boosted together into one fast, accurate " +
						"detector. AdaBoost also shows up in spam filtering, combining many weak " +
						"keyword and header rules, and was one of the first boosting algorithms " +
						"to move from theory into everyday production use, years before " +
						"`gradient-boosting`'s more general residual-fitting version took over on " +
						"tabular data.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'round 2's stump is fit on reweighted data where x=5 and " +
						"x=6 now carry double the vote of the other points, not on the original " +
						"evenly-weighted set' -- every round after the first sees a different " +
						"weighted distribution than the round before it.",
					"Not like this: assuming AdaBoost keeps only the latest or the " +
						"best-performing stump and throws the rest away. Every stump built along " +
						"the way keeps voting forever, each at its own fixed alpha -- round 1's " +
						"stump still votes on x=1 even after round 3 is added, it just carries " +
						"less weight than a lower-error round would.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "rounds", Label: "Boosting rounds", Min: 0, Max: 3, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 460, 0, 1, 0, 1)
	c.Axes()
	return c.String()
}
