// Package adaboost visualizes AdaBoost: instead of scoring a single
// `decision-trees`-style split once and living with whatever it gets wrong,
// fit a sequence of one-split stumps, each on a version of the data where
// every point misclassified so far has been reweighted heavier -- so each
// new stump is aimed specifically at whatever the stumps before it still
// get wrong -- then combine every stump's vote, weighted by how accurate it
// was, into one stronger classifier.
package adaboost

import (
	"fmt"
	"math"
	"sort"

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

// Xs and Ys are the eight points every Section walks through: labels
// alternate in pairs by class -- x=1,2 and x=5,6 are +1; x=3,4 and x=7,8
// are -1 -- so no single stump can separate them, forcing AdaBoost to
// actually combine more than one.
var (
	Xs = []float64{1, 2, 3, 4, 5, 6, 7, 8}
	Ys = []int{1, 1, -1, -1, 1, 1, -1, -1}
)

// Stump is one weighted decision stump: predict LeftLabel for x<=Threshold,
// RightLabel otherwise. Unlike a single "flip both sides" polarity, each
// side's label is chosen independently by whichever class carries more of
// the weighted vote on that side, so a stump can predict +1 on both sides,
// -1 on both sides, or either order.
type Stump struct {
	Threshold             float64
	LeftLabel, RightLabel int
}

// StumpPredict returns a stump's predicted label (+1 or -1) for one x.
func StumpPredict(s Stump, x float64) int {
	if x <= s.Threshold {
		return s.LeftLabel
	}
	return s.RightLabel
}

func predictAll(s Stump, xs []float64) []int {
	preds := make([]int, len(xs))
	for i, x := range xs {
		preds[i] = StumpPredict(s, x)
	}
	return preds
}

// weightedLabel returns the weighted-majority label (+1 or -1) on each side
// of threshold th: whichever label carries more total weight among the
// points on that side. Ties resolve to +1.
func weightedLabel(xs []float64, ys []int, weights []float64, th float64) (left, right int) {
	var leftPos, leftNeg, rightPos, rightNeg float64
	for i, x := range xs {
		pos, neg := &rightPos, &rightNeg
		if x <= th {
			pos, neg = &leftPos, &leftNeg
		}
		if ys[i] == 1 {
			*pos += weights[i]
		} else {
			*neg += weights[i]
		}
	}
	left, right = 1, 1
	if leftNeg > leftPos {
		left = -1
	}
	if rightNeg > rightPos {
		right = -1
	}
	return
}

// WeightedError returns the total weight of the points where predicted
// disagrees with actual. actual, predicted, and weights must be the same
// length and index-aligned.
func WeightedError(actual, predicted []int, weights []float64) float64 {
	var err float64
	for i := range actual {
		if actual[i] != predicted[i] {
			err += weights[i]
		}
	}
	return err
}

// FitWeightedStump scans every midpoint between consecutive distinct sorted
// x values and returns the threshold (with each side's weighted-majority
// label) that minimizes total weighted classification error, plus that
// error. xs, ys, and weights must be the same length and index-aligned.
func FitWeightedStump(xs []float64, ys []int, weights []float64) (Stump, float64) {
	uniq := append([]float64(nil), xs...)
	sort.Float64s(uniq)
	dedup := uniq[:0]
	for i, x := range uniq {
		if i == 0 || x != dedup[len(dedup)-1] {
			dedup = append(dedup, x)
		}
	}

	best := Stump{Threshold: dedup[len(dedup)-1], LeftLabel: 1, RightLabel: 1}
	bestErr := math.Inf(1)
	for i := 0; i < len(dedup)-1; i++ {
		th := (dedup[i] + dedup[i+1]) / 2
		left, right := weightedLabel(xs, ys, weights, th)
		st := Stump{Threshold: th, LeftLabel: left, RightLabel: right}
		err := WeightedError(ys, predictAll(st, xs), weights)
		if err < bestErr {
			bestErr = err
			best = st
		}
	}
	return best, bestErr
}

// Alpha is AdaBoost's per-round classifier vote weight: 0.5*ln((1-err)/err).
// It's positive whenever err<0.5 (better than a coin flip) and grows
// without bound as err shrinks toward 0 -- a near-perfect stump dominates
// the combined vote.
func Alpha(err float64) float64 {
	if err <= 0 {
		err = 1e-10
	}
	if err >= 1 {
		err = 1 - 1e-10
	}
	return 0.5 * math.Log((1-err)/err)
}

// UpdateWeights reweights each point by exp(-alpha*actual*predicted): a
// point the current stump got right (actual*predicted=+1) shrinks by
// exp(-alpha), a point it got wrong (actual*predicted=-1) grows by
// exp(alpha), then the whole set is renormalized back to sum to 1 so it
// stays a valid distribution for the next round's fit.
func UpdateWeights(actual, predicted []int, weights []float64, alpha float64) []float64 {
	out := make([]float64, len(weights))
	var sum float64
	for i := range weights {
		out[i] = weights[i] * math.Exp(-alpha*float64(actual[i]*predicted[i]))
		sum += out[i]
	}
	for i := range out {
		out[i] /= sum
	}
	return out
}

// Round is the record of one boosting round: the stump fit to that round's
// weights, its weighted error, and the resulting vote weight alpha.
type Round struct {
	Stump Stump
	Err   float64
	Alpha float64
}

// Run performs numRounds rounds of AdaBoost on (xs, ys), starting from
// uniform weights 1/n. It returns the sequence of rounds in order and the
// weights history: weightsHistory[0] is the initial uniform weights, and
// weightsHistory[t] for t>0 is what round t produced (and what round t+1,
// if any, was fit on). Pure and deterministic: same inputs always produce
// the same rounds, no randomness -- unlike random-forest's resampling,
// every round here is a direct function of the round before it.
func Run(xs []float64, ys []int, numRounds int) (rounds []Round, weightsHistory [][]float64) {
	weights := make([]float64, len(xs))
	for i := range weights {
		weights[i] = 1.0 / float64(len(xs))
	}
	weightsHistory = append(weightsHistory, append([]float64(nil), weights...))

	for t := 0; t < numRounds; t++ {
		st, err := FitWeightedStump(xs, ys, weights)
		alpha := Alpha(err)
		weights = UpdateWeights(ys, predictAll(st, xs), weights, alpha)
		rounds = append(rounds, Round{Stump: st, Err: err, Alpha: alpha})
		weightsHistory = append(weightsHistory, append([]float64(nil), weights...))
	}
	return
}

// Predict returns the ensemble's predicted label at x: the sign of the
// alpha-weighted sum of every round's stump prediction. A sum of exactly
// 0 (never reached on this dataset, but defined for completeness) resolves
// to +1.
func Predict(rounds []Round, x float64) int {
	var sum float64
	for _, r := range rounds {
		sum += r.Alpha * float64(StumpPredict(r.Stump, x))
	}
	if sum < 0 {
		return -1
	}
	return 1
}

// maxRounds is the deepest boosting sequence Render ever needs. Fitting it
// once and slicing prefixes is equivalent to re-running Run at each smaller
// round count, since round t's stump only ever depends on the rounds
// before it.
const maxRounds = 3

// lineColors gives each active round's threshold line (and matching
// readout text) a distinct, consistent color across the whole rounds
// slider, so round 1's line is always the same color whether it's the
// only one showing or one of three.
var lineColors = []string{viz.Good, viz.Warm, viz.Accent}

func render(p map[string]float64) string {
	rounds := int(p["rounds"])

	allRounds, weightsHistory := Run(Xs, Ys, maxRounds)
	active := allRounds[:rounds]
	weights := weightsHistory[rounds]

	// Scale point size against the largest weight seen across ANY round,
	// not just this frame's, so a point's square is directly comparable
	// as the rounds slider moves -- watch x=5,6 balloon after round 1, then
	// shrink back as x=3,4 balloon after round 2.
	globalMax := 0.0
	for _, ws := range weightsHistory {
		for _, w := range ws {
			if w > globalMax {
				globalMax = w
			}
		}
	}

	c := viz.New(680, 460, 0.3, 8.7, 0, 1)
	c.PadT = 130
	c.PadB = 40
	c.Axes()
	for x := 1.0; x <= 8; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	c.Text(10, c.Y(0.72), "+1", 11, viz.Accent, "start")
	c.Text(10, c.Y(0.32), "-1", 11, viz.Warm, "start")

	for i, r := range active {
		c.VLine(r.Stump.Threshold, lineColors[i%len(lineColors)], true)
	}

	nWrong := 0
	for i, x := range Xs {
		y, classColor := 0.3, viz.Warm
		if Ys[i] == 1 {
			y, classColor = 0.7, viz.Accent
		}
		px, py := c.X(x), c.Y(y)

		size := 10.0
		if globalMax > 0 {
			size = 10 + 24*(weights[i]/globalMax)
		}

		if rounds > 0 {
			ring := viz.Good
			if Predict(active, x) != Ys[i] {
				ring = viz.Bad
				nWrong++
			}
			c.Rect(px-size/2, py-size/2, size, size, ring, 1)
			inner := size - 6
			if inner < 2 {
				inner = 2
			}
			c.Rect(px-inner/2, py-inner/2, inner, inner, classColor, 1)
		} else {
			c.Rect(px-size/2, py-size/2, size, size, classColor, 1)
		}
	}

	c.Text(16, 20, fmt.Sprintf("rounds = %d of %d   square size = current point weight", rounds, maxRounds), 14, viz.Ink, "start")
	ty := 44.0
	for i, r := range active {
		c.Text(16, ty, fmt.Sprintf("round %d: split at x=%.1f, weighted error=%.3f, alpha=%.3f",
			i+1, r.Stump.Threshold, r.Err, r.Alpha), 13, lineColors[i%len(lineColors)], "start")
		ty += 20
	}
	if rounds == 0 {
		c.Text(16, ty, "no classifier yet -- every point starts at weight 1/8", 13, viz.Muted, "start")
	} else {
		c.Text(16, ty, fmt.Sprintf("ensemble vote so far: %d/%d correct   (green ring = correct, red ring = wrong)",
			len(Xs)-nWrong, len(Xs)), 13, viz.Ink, "start")
	}

	return c.String()
}
