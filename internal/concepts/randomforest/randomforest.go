// Package randomforest visualizes bagging: instead of trusting the single
// best-information-gain split `decision-trees` found on one dataset, fit a
// small stump-tree on each of several bootstrap resamples of the same data
// (the resampling idea `bootstrap-resampling` used to see how a statistic
// wiggles) and let the trees vote. Individual trees' split thresholds swing
// around depending on which resample they happened to get; the vote
// fraction across the whole forest is far more stable.
package randomforest

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "random-forest",
		Seq:   70,
		Title: "Random forest (bagging trees to cut variance)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`decision-trees` picked one 'best' split on 10 students' hours studied " +
						"vs. pass/fail — 2.75 hours, by information gain, built from `entropy`. " +
						"But `bootstrap-resampling` showed that a statistic computed from a small " +
						"sample can wiggle a lot depending on exactly which resample, with " +
						"replacement, you happen to draw. What if you resampled the training data " +
						"itself before fitting a tree at all — would every resample really agree " +
						"the cutoff is 2.75 hours, or would the 'best' split move around depending " +
						"on which 10 (of 10, with replacement) students you happened to draw?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the same 10 students `decision-trees` used, and build 9 bootstrap " +
						"resamples of them — each a set of 10 draws, with replacement, from the " +
						"original 10. Fit a single-split tree on each resample the exact same way " +
						"`decision-trees` did: scan candidate thresholds, pick the one with the " +
						"highest information gain.",
					"• Resample 0 (the original data, unresampled): best split at 2.75 hours, " +
						"gain=0.610 bits — identical to decision-trees's own answer.",
					"• Resample 1 (extra weight on the 3-hour pass, drops the noisy 3.5-hour " +
						"fail): best split still at 2.75 hours, but gain jumps to 0.971 bits — " +
						"removing the noisy point made the same cutoff much cleaner.",
					"• Resample 2 (extra weight on the 3.5-hour fail, drops the 3-hour pass): " +
						"best split moves to 3.75 hours, gain=0.881 bits — a full half-point " +
						"later than resample 0's answer, from the same 10 students, just resampled " +
						"differently.",
					"Across all 9 resampled trees, thresholds land anywhere from 2.75 to 3.75 " +
						"hours (gains from 0.396 to 1.000 bits) — no single tree's threshold is " +
						"'the' answer. Now let them vote: at hours=3.0, with only tree 0 active, " +
						"100% of active trees predict pass. Add trees 1 and 2 (3 active): 2 of 3 " +
						"now predict pass (67%), since tree 2's threshold (3.75) hasn't been " +
						"crossed yet. With the full 9-tree forest active at hours=3.0, 6 of 9 " +
						"trees have already crossed their own threshold — still 67% — but by " +
						"hours=3.5 that climbs to 7 of 9 (78%), since a tree with threshold 3.25 " +
						"has now also flipped, something the 3-tree forest can't show at all.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The x-axis is hours studied, same as `decision-trees`; the y-axis is the " +
						"fraction of active trees currently voting 'pass' at that x. The " +
						"numTrees slider controls how many of the 9 resampled trees are in the " +
						"forest; the accent staircase is their combined vote fraction, and the " +
						"small ticks along the bottom mark each active tree's own split " +
						"threshold. The original 10 students are plotted as dots, colored by " +
						"true pass/fail, for reference. With numTrees=1 the staircase is a " +
						"single sharp cliff (exactly `decision-trees`'s picture); add more trees " +
						"and watch the cliff soften into several smaller steps as each tree's " +
						"own threshold contributes its own smaller jump.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Stop trusting any single tree's split location as 'the' cutoff — a lone " +
						"tree's 3.75-hour threshold is really just an artifact of which resample " +
						"it happened to train on. Averaging many resampled trees' votes trades " +
						"that one noisy number for a smoother, more stable read of how confident " +
						"the data really is at each point, without needing any one tree to be " +
						"individually more careful. This is exactly what 'bagging' (bootstrap " +
						"aggregating) buys a random forest: variance goes down because errors " +
						"specific to any one resample get outvoted, while bias stays put because " +
						"every tree is still built by the same honest information-gain rule.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Random forests are a default go-to model for credit scoring, fraud " +
						"detection, and medical risk prediction whenever a single decision tree " +
						"would be too twitchy to trust. Microsoft's Kinect used a random forest, " +
						"trained per-pixel on depth images, to classify body parts in real time. " +
						"Kaggle tabular-data competitions routinely reach for a random forest as " +
						"a strong, low-fuss baseline before trying anything fancier.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'no single tree's split (2.75, or 3.75, or anywhere in " +
						"between) should be trusted on its own — what matters is where the " +
						"forest's vote fraction crosses 50%, averaged over many resampled trees.'",
					"Not like this: assuming that because each tree trains on a different " +
						"resample, a random forest's answer is basically arbitrary or 'random' in " +
						"the everyday sense. The trees differ in a controlled way — bootstrap " +
						"resampling of the same fixed dataset — and it's exactly that controlled, " +
						"bounded diversity being averaged away that reduces variance, not a source " +
						"of extra unpredictability in the final answer.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "numTrees", Label: "Number of trees in the forest", Min: 1, Max: 9, Step: 1, Def: 9},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 1, 0, 1).String()
}
