// Package decisiontree visualizes how a decision tree picks where to split:
// scoring every candidate threshold on a single feature by information gain
// (built from `entropy`'s bits-of-surprise measure) rather than by where the
// feature's values happen to be numerically centered.
package decisiontree

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "decision-trees",
		Seq:   56,
		Title: "Decision trees (splitting by information gain)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`entropy` measured how surprising a set of outcomes is, in bits. Here are " +
						"10 students' hours studied and whether they passed: 5 passed, 5 didn't — " +
						"exactly the maximum-uncertainty case entropy warns about, a coin flip. You " +
						"have one number per student (hours studied) and want a simple rule: 'if " +
						"hours studied is above some threshold, predict pass.' Gut instinct: just " +
						"split at the middle of the range — hours run from 1 to 5.5, so try 3.25. " +
						"Is the numeric midpoint actually the best place to draw that line, and how " +
						"would you even measure 'best'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Score a candidate split the same way entropy scored a whole distribution: " +
						"how much less uncertain are the two resulting groups, on average, than the " +
						"whole group was before you split? Start with the whole group: 5 of 10 " +
						"passed, so parent entropy = 1.0 bit — total uncertainty, a coin flip.",
					"Try the 'obvious' midpoint, threshold = 3.25 hours: the low-hours group " +
						"{1, 1.5, 2, 2.5, 3} has 1 pass out of 5 (entropy ≈ 0.72), and the " +
						"high-hours group {3.5, 4, 4.5, 5, 5.5} has 4 pass out of 5 (entropy ≈ " +
						"0.72 too) — both sides are still fairly mixed. Weighted average child " +
						"entropy = 0.5×0.72 + 0.5×0.72 ≈ 0.72, so information gain = 1.0 − 0.72 = " +
						"0.28 bits. The midpoint only shaves off about a quarter of the original " +
						"uncertainty.",
					"Now try threshold = 2.75 instead: the low side {1, 1.5, 2, 2.5} has 0 passes " +
						"out of 4 — perfectly pure, entropy = 0 — and the high side {3, 3.5, 4, " +
						"4.5, 5, 5.5} has 5 passes out of 6 (entropy ≈ 0.65). Weighted average = " +
						"0.4×0 + 0.6×0.65 ≈ 0.39, so gain = 1.0 − 0.39 = 0.61 bits — more than " +
						"double the midpoint's gain, because the low side came out completely " +
						"pure. Scanning every candidate split confirms 2.75 (tied with 3.75) is " +
						"the actual best split, not 3.25 — the midpoint just happened to land " +
						"inside the messiest part of the data.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"All 10 students plotted by hours studied, colored by pass/fail. The " +
						"threshold slider draws a vertical split line; the readout below reports " +
						"exactly the numbers worked through above for whichever threshold you pick " +
						"— left/right group sizes, each side's pass fraction and entropy, and the " +
						"resulting information gain — plus, for comparison, the best gain " +
						"achievable by any threshold on this data. Drag the slider off 3.25 toward " +
						"2.75 or 3.75 and watch the gain jump from 0.28 up to 0.61.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Instead of eyeballing where to draw a split, you can score every candidate " +
						"threshold with one number — information gain — and mechanically pick the " +
						"one that reduces uncertainty the most. That's the exact rule real " +
						"decision-tree algorithms (ID3, C4.5, CART) use at every branch: try every " +
						"feature and every candidate threshold, split on whichever gives the " +
						"highest information gain, then repeat inside each resulting group.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Spam filters splitting on word counts, medical triage splitting on a lab " +
						"value, credit scoring splitting on an income-to-debt ratio — anywhere you " +
						"want a human-readable 'if X is above/below this line' rule instead of an " +
						"opaque score. Random forests and gradient-boosted trees are, underneath, " +
						"large committees of exactly this kind of tree, each one grown by repeating " +
						"this same highest-information-gain split over and over.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the split at 2.75 hours gives more than double the " +
						"midpoint's information gain, because it produces a perfectly pure " +
						"low-hours group' — anchoring the claim on the measured entropy reduction.",
					"Not like this: 'the best threshold is always the middle or median of the " +
						"data' — the best split is decided by how cleanly it separates the labels, " +
						"not by which value is numerically in the center of the feature's range; " +
						"here the truly best splits (2.75 and 3.75) sit off-center, on either side " +
						"of the messiest two data points, not on top of them.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "threshold", Label: "Split threshold (hours studied)", Min: 1, Max: 5.5, Step: 0.05, Def: 3.25},
		},
		Render: render,
	})
}

// Hours and Pass are the 10 students used throughout this concept: hours
// studied and whether each one passed (1) or failed (0). Two points break
// the otherwise-clean trend on purpose — hours=3 passed early, hours=3.5
// failed late — so the "obviously best" midpoint split isn't actually best.
var (
	Hours = []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5}
	Pass  = []int{0, 0, 0, 0, 1, 0, 1, 1, 1, 1}
)

// Entropy is the binary entropy, in bits, of a class with fraction p
// positive. It's 0 when p is 0 or 1 (no uncertainty at all) and peaks at 1
// bit when p=0.5 (maximum uncertainty, a coin flip).
func Entropy(p float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	return -p*math.Log2(p) - (1-p)*math.Log2(1-p)
}

// classEntropy is Entropy applied to a slice of 0/1 labels, using the
// fraction of 1s as p.
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
// entropies. It's always >= 0 (a split can never increase uncertainty on
// average) and equals the parent entropy exactly when the split perfectly
// separates the two classes (both children pure, entropy 0 each).
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

// bestSplit scans every midpoint between consecutive Hours values and
// returns the threshold with the highest information gain (the first one
// found, in case of a tie).
func bestSplit() (threshold, gain float64) {
	best := -1.0
	bestT := 0.0
	for i := 0; i < len(Hours)-1; i++ {
		th := (Hours[i] + Hours[i+1]) / 2
		left, right := Split(Hours, Pass, th)
		g := InfoGain(Pass, left, right)
		if g > best {
			best = g
			bestT = th
		}
	}
	return bestT, best
}

func countAndFrac(labels []int) (pos int, frac float64) {
	for _, l := range labels {
		pos += l
	}
	if len(labels) > 0 {
		frac = float64(pos) / float64(len(labels))
	}
	return
}

func render(p map[string]float64) string {
	threshold := p["threshold"]

	left, right := Split(Hours, Pass, threshold)
	parent := classEntropy(Pass)
	leftPos, leftFrac := countAndFrac(left)
	rightPos, rightFrac := countAndFrac(right)
	leftE, rightE := classEntropy(left), classEntropy(right)
	gain := InfoGain(Pass, left, right)
	bestT, bestGain := bestSplit()

	c := viz.New(680, 400, 0, 6, 0, 1)
	c.Axes()
	for x := 1.0; x <= 5.5; x += 0.5 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	c.Text(10, c.Y(0.72), "pass", 11, viz.Good, "start")
	c.Text(10, c.Y(0.32), "fail", 11, viz.Bad, "start")

	for i, h := range Hours {
		y, color := 0.3, viz.Bad
		if Pass[i] == 1 {
			y, color = 0.7, viz.Good
		}
		px, py := c.X(h), c.Y(y)
		c.Rect(px-5, py-5, 10, 10, color, 0.9)
	}

	c.VLine(threshold, viz.Ink, false)

	c.Text(20, 24, fmt.Sprintf("threshold = %.2f hours   parent entropy = %.3f bits (5/10 passed)", threshold, parent), 13, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("left:  n=%d, %d passed (p=%.2f), entropy=%.3f", len(left), leftPos, leftFrac, leftE), 13, viz.Muted, "start")
	c.Text(20, 62, fmt.Sprintf("right: n=%d, %d passed (p=%.2f), entropy=%.3f", len(right), rightPos, rightFrac, rightE), 13, viz.Muted, "start")
	c.Text(20, 88, fmt.Sprintf("information gain = %.3f bits", gain), 15, viz.Accent, "start")
	c.Text(20, 108, fmt.Sprintf("best possible split: threshold=%.2f hours, gain=%.3f bits", bestT, bestGain), 12, viz.Warm, "start")

	return c.String()
}
