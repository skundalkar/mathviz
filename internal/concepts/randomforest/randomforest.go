// Package randomforest visualizes bagging: instead of trusting the single
// best-information-gain split `decision-trees` found on one dataset, fit a
// small stump-tree on each of several bootstrap resamples of the same data
// (the resampling idea `bootstrap-resampling` used to see how a statistic
// wiggles) and let the trees vote. Individual trees' split thresholds swing
// around depending on which resample they happened to get; the vote
// fraction across the whole forest is far more stable.
package randomforest

import (
	"fmt"
	"math"
	"sort"

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

// Hours and Pass are the same 10 students `decision-trees` used: hours
// studied and whether each one passed (1) or failed (0). Kept as its own
// copy here -- each concept package is self-contained (see BUILD_CYCLE.md).
var (
	Hours = []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5}
	Pass  = []int{0, 0, 0, 0, 1, 0, 1, 1, 1, 1}
)

// treeSamples is 9 fixed, hand-picked bootstrap resamples of Hours/Pass --
// each a list of 10 indices into Hours/Pass, drawn with replacement. Kept as
// literals (like anova's baseOffsets) rather than generated with math/rand
// at Render time, since Render must be pure: same input, same output, no
// randomness. Sample 0 is the original, unresampled data -- exactly what
// decision-trees fit its single split on -- so the forest's first tree
// always matches decision-trees's own answer.
var treeSamples = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, // original data, unresampled
	{0, 1, 2, 3, 4, 4, 4, 6, 7, 9}, // triples the 3h pass, drops the 3.5h fail
	{0, 1, 2, 3, 5, 5, 5, 6, 7, 9}, // triples the 3.5h fail, drops the 3h pass
	{0, 0, 1, 2, 3, 6, 7, 8, 9, 9},
	{1, 2, 3, 4, 4, 5, 5, 6, 8, 9},
	{0, 1, 3, 4, 6, 7, 7, 8, 9, 9},
	{0, 2, 2, 4, 5, 6, 7, 8, 9, 9},
	{0, 1, 2, 3, 4, 5, 6, 7, 7, 8},
	{1, 1, 2, 3, 4, 5, 6, 8, 9, 9},
}

// Entropy is the binary entropy, in bits, of a class with fraction p
// positive -- 0 when p is 0 or 1, peaking at 1 bit when p=0.5.
func Entropy(p float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	return -p*math.Log2(p) - (1-p)*math.Log2(1-p)
}

func classEntropy(labels []int) float64 {
	if len(labels) == 0 {
		return 0
	}
	pos := 0
	for _, l := range labels {
		pos += l
	}
	return Entropy(float64(pos) / float64(len(labels)))
}

// Split partitions labels by whether each example's feature value is at or
// below threshold (left) or above it (right). hours and labels must be the
// same length and index-aligned.
func Split(hours []float64, labels []int, threshold float64) (left, right []int) {
	for i, h := range hours {
		if h <= threshold {
			left = append(left, labels[i])
		} else {
			right = append(right, labels[i])
		}
	}
	return
}

// InfoGain is the information gain of splitting labels into left and right:
// the parent's entropy minus the size-weighted average of the two children's
// entropies.
func InfoGain(labels, left, right []int) float64 {
	parent := classEntropy(labels)
	n := float64(len(labels))
	if n == 0 {
		return 0
	}
	wLeft := float64(len(left)) / n
	wRight := float64(len(right)) / n
	return parent - (wLeft*classEntropy(left) + wRight*classEntropy(right))
}

// majority returns 1 if labels has at least as many 1s as 0s, else 0 --
// what a leaf predicts for any example that lands in it.
func majority(labels []int) int {
	pos := 0
	for _, l := range labels {
		pos += l
	}
	if pos*2 >= len(labels) {
		return 1
	}
	return 0
}

// Tree is one single-split ("stump") decision tree: a threshold on hours,
// and the majority label each side predicts.
type Tree struct {
	Threshold             float64
	LeftLabel, RightLabel int
	Gain                  float64
}

// BuildTree fits a single-split tree on hours/labels by scanning every
// midpoint between consecutive distinct hours values (as decision-trees
// does) and keeping the threshold with the highest information gain.
func BuildTree(hours []float64, labels []int) Tree {
	uniq := append([]float64(nil), hours...)
	sort.Float64s(uniq)
	dedup := uniq[:0]
	for i, h := range uniq {
		if i == 0 || h != dedup[len(dedup)-1] {
			dedup = append(dedup, h)
		}
	}

	best := -1.0
	bestT := 0.0
	for i := 0; i < len(dedup)-1; i++ {
		th := (dedup[i] + dedup[i+1]) / 2
		left, right := Split(hours, labels, th)
		g := InfoGain(labels, left, right)
		if g > best {
			best = g
			bestT = th
		}
	}
	left, right := Split(hours, labels, bestT)
	return Tree{Threshold: bestT, LeftLabel: majority(left), RightLabel: majority(right), Gain: best}
}

// Predict returns the tree's predicted label for a given hours value.
func Predict(t Tree, hours float64) int {
	if hours <= t.Threshold {
		return t.LeftLabel
	}
	return t.RightLabel
}

// Forest builds the first n trees (1..len(treeSamples)) from treeSamples,
// each fit on its own bootstrap resample of Hours/Pass. n is clamped to
// [1, len(treeSamples)].
func Forest(n int) []Tree {
	if n < 1 {
		n = 1
	}
	if n > len(treeSamples) {
		n = len(treeSamples)
	}
	trees := make([]Tree, n)
	for i := 0; i < n; i++ {
		idx := treeSamples[i]
		hs := make([]float64, len(idx))
		ls := make([]int, len(idx))
		for j, ix := range idx {
			hs[j] = Hours[ix]
			ls[j] = Pass[ix]
		}
		trees[i] = BuildTree(hs, ls)
	}
	return trees
}

// VoteFraction is the fraction of trees predicting "pass" (1) at hours.
func VoteFraction(trees []Tree, hours float64) float64 {
	if len(trees) == 0 {
		return 0
	}
	votes := 0
	for _, t := range trees {
		votes += Predict(t, hours)
	}
	return float64(votes) / float64(len(trees))
}

func render(p map[string]float64) string {
	numTrees := int(p["numTrees"])
	trees := Forest(numTrees)

	// y-range: [0,1] holds the vote-fraction staircase; the strip below 0
	// holds a data rug (true labels) and each active tree's own threshold
	// tick, so both stay visible without overlapping the curve.
	c := viz.New(680, 440, 1, 5.5, -0.22, 1.08)
	c.PadT = 90
	c.PadB = 40
	c.Axes()
	for x := 1.0; x <= 5.5; x += 0.5 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	curve := viz.Sample(1, 5.5, 400, func(x float64) float64 { return VoteFraction(trees, x) })
	c.Path(curve, viz.Accent, 2.5)

	// Each active tree's own split threshold, as a small tick below the axis.
	for _, tr := range trees {
		px, py := c.X(tr.Threshold), c.Y(-0.18)
		c.Rect(px-2, py-5, 4, 10, viz.Warm, 0.8)
	}

	// The original 10 students, colored by true pass/fail, as a rug plot.
	for i, h := range Hours {
		color := viz.Bad
		if Pass[i] == 1 {
			color = viz.Good
		}
		px, py := c.X(h), c.Y(-0.08)
		c.Rect(px-4, py-4, 8, 8, color, 0.85)
	}

	minT, maxT := trees[0].Threshold, trees[0].Threshold
	for _, tr := range trees {
		if tr.Threshold < minT {
			minT = tr.Threshold
		}
		if tr.Threshold > maxT {
			maxT = tr.Threshold
		}
	}

	c.Text(20, 24, fmt.Sprintf("%d tree(s) active   individual thresholds range %.2f-%.2f hours", numTrees, minT, maxT), 13, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("forest vote-fraction(pass) at hours=3.0: %.0f%%    at hours=3.5: %.0f%%",
		VoteFraction(trees, 3.0)*100, VoteFraction(trees, 3.5)*100), 13, viz.Muted, "start")
	c.Text(20, 64, "accent staircase = fraction of active trees voting pass   orange ticks = each tree's own split", 12, viz.Muted, "start")

	return c.String()
}
