// Package gradientboosting visualizes gradient boosting: instead of
// averaging many independently-trained trees the way `random-forest` does,
// fit one small stump at a time, each aimed specifically at the residual
// error the trees fit so far still leave behind. The running example is
// five points, y = [2,3,3,8,9] at x = [1,2,3,4,5], where a first stump
// captures the big level jump between x=3 and x=4 and a second stump mops
// up most of what's left.
package gradientboosting

import (
	"fmt"
	"sort"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "gradient-boosting",
		Seq:   78,
		Title: "Gradient boosting (sequentially fitting the residual)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`random-forest` averaged many independently-trained trees to smooth out " +
						"noise -- but every one of those trees trained on its own resample, blind " +
						"to what the other trees in the forest got right or wrong. And a single " +
						"stump (one split, two possible outputs) is a blunt tool by itself: fit one " +
						"stump to five points -- y=[2,3,3,8,9] at x=[1,2,3,4,5] -- and its best " +
						"split can only ever produce two flat levels, no matter how many more " +
						"shapes the data actually has. What if, instead of training every tree from " +
						"scratch on the same target, each new tree looked specifically at what the " +
						"trees before it got wrong, and aimed itself at fixing exactly that?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Start with the simplest possible prediction: F0(x) = mean(y) = " +
						"(2+3+3+8+9)/5 = 5 for every x, ignoring x entirely. Its errors " +
						"(residuals), y minus F0, are [-3,-2,-2,3,4] -- exactly the signal a first " +
						"stump should target, not y itself. Fit a stump to predict *that residual* " +
						"from x: scanning every split, threshold=3.5 wins (mean squared error " +
						"1.17, versus 21+ for the other candidates), giving left mean -7/3≈-2.33 " +
						"(x≤3.5) and right mean 3.5 (x>3.5).",
					"Add that stump's output to F0: F1(x) = 5 + (-2.33) = 2.67 for x∈{1,2,3}, and " +
						"5 + 3.5 = 8.5 for x∈{4,5}. Training error (sum of squared residuals) drops " +
						"from 42 (F0 alone) to 1.17 -- the big level jump at x=3.5 is now captured. " +
						"What's left, y minus F1, is [-0.67, 0.33, 0.33, -0.5, 0.5] -- much smaller, " +
						"but not zero: x=1's -0.67 stands out from x=2 and x=3's +0.33. A second " +
						"stump fit to *this* residual splits at threshold=1.5 (left mean -2/3 for " +
						"x=1 alone, right mean 1/6 for x∈{2,3,4,5}). Adding it in gives F2(1)=2.67-" +
						"0.67=2.0 exactly matching y=2, and total squared error keeps falling, from " +
						"1.17 to 0.61. Each stump is cheap and individually weak, but every one is " +
						"aimed at whatever mistake is still on the table.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"numStages adds one more sequential stump (0 = just the flat baseline F0); " +
						"learningRate shrinks how much of each stump's correction actually gets " +
						"added in. The black squares are the five true (x,y) points; the accent " +
						"staircase is the cumulative prediction after the chosen number of stages. " +
						"At numStages=0 it's a single flat line at 5; at numStages=1 it jumps to " +
						"two levels at x=3.5; at numStages=2 it further refines around x=1. The " +
						"readout reports the training sum-of-squared-errors at the current stage, " +
						"falling from 42 towards 0 as more stumps are added.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Approximate a target that no single stump could fit well, by chaining many " +
						"weak learners that each specialize in whatever error is still left over -- " +
						"and control the trade-off explicitly with the learning rate: lr=1.0 fits " +
						"the training residual as fast as possible in the fewest stages, while a " +
						"smaller lr (like 0.3) takes smaller, more cautious steps that need more " +
						"stages to reach the same fit but tend to generalize better on new data, " +
						"the same bias/shrinkage trade-off `regularization` and " +
						"`bias-variance-tradeoff` cover for other model families.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"XGBoost and LightGBM -- both gradient-boosted-tree libraries -- have won a " +
						"large share of Kaggle tabular-data competitions and power production " +
						"systems for click-through-rate prediction, credit risk scoring, and search " +
						"ranking. They run the exact same fit-the-residual loop shown here, just " +
						"with deeper trees than a single stump, many more stages, and losses other " +
						"than squared error.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'stage two's stump is fit to predict stage one's leftover " +
						"error, not to predict y itself' -- every stump after the first is trained " +
						"on a completely different target than the one before it.",
					"Not like this: assuming a higher learning rate is simply 'better' because it " +
						"drives the training error down faster. lr=1.0 reaches the lowest training " +
						"error in the fewest stages here, but on noisier real data that same " +
						"aggressiveness fits the training set's noise, not just its signal -- " +
						"exactly the overfitting risk `regularization` warns about, just reached by " +
						"adding too many full-strength stages instead of too little penalty.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "numStages", Label: "Number of boosting stages", Min: 0, Max: 4, Step: 1, Def: 2},
			{Key: "learningRate", Label: "Learning rate (shrinkage)", Min: 0.1, Max: 1.0, Step: 0.1, Def: 1.0},
		},
		Render: render,
	})
}

// Xs and Ys are the five fixed observations the worked example walks
// through: a big level jump between x=3 and x=4, hand-picked so a single
// stump captures most of the signal and a second stump mops up most of
// what's left.
var (
	Xs = []float64{1, 2, 3, 4, 5}
	Ys = []float64{2, 3, 3, 8, 9}
)

// Stump is one regression split: predict LeftVal for x<=Threshold,
// RightVal otherwise.
type Stump struct {
	Threshold         float64
	LeftVal, RightVal float64
}

// Mean returns the arithmetic mean of vs, or 0 for an empty slice.
func Mean(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}
	var sum float64
	for _, v := range vs {
		sum += v
	}
	return sum / float64(len(vs))
}

// SSE returns the sum of squared errors between actual and predicted,
// index-aligned.
func SSE(actual, predicted []float64) float64 {
	var sum float64
	for i := range actual {
		d := actual[i] - predicted[i]
		sum += d * d
	}
	return sum
}

// sseAroundMean is the sum of squared deviations of vs from its own mean --
// the error a single constant prediction (the mean) leaves behind.
func sseAroundMean(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}
	m := Mean(vs)
	var sum float64
	for _, v := range vs {
		d := v - m
		sum += d * d
	}
	return sum
}

// FitStump fits a regression stump to predict target from xs: it scans
// every midpoint between consecutive distinct sorted x values and keeps the
// threshold whose two-sided split minimizes total sum-of-squared-error,
// predicting each side's mean of target. xs and target must be the same
// length and index-aligned.
func FitStump(xs, target []float64) Stump {
	uniq := append([]float64(nil), xs...)
	sort.Float64s(uniq)
	dedup := uniq[:0]
	for i, x := range uniq {
		if i == 0 || x != dedup[len(dedup)-1] {
			dedup = append(dedup, x)
		}
	}

	bestSSE := sseAroundMean(target) // fallback: one flat level, no split
	bestThreshold := dedup[len(dedup)-1]
	for i := 0; i < len(dedup)-1; i++ {
		th := (dedup[i] + dedup[i+1]) / 2
		var left, right []float64
		for j, x := range xs {
			if x <= th {
				left = append(left, target[j])
			} else {
				right = append(right, target[j])
			}
		}
		sse := sseAroundMean(left) + sseAroundMean(right)
		if sse < bestSSE {
			bestSSE = sse
			bestThreshold = th
		}
	}

	var left, right []float64
	for j, x := range xs {
		if x <= bestThreshold {
			left = append(left, target[j])
		} else {
			right = append(right, target[j])
		}
	}
	return Stump{Threshold: bestThreshold, LeftVal: Mean(left), RightVal: Mean(right)}
}

// StumpPredict returns a stump's output for one x.
func StumpPredict(s Stump, x float64) float64 {
	if x <= s.Threshold {
		return s.LeftVal
	}
	return s.RightVal
}

// Boost runs gradient boosting for squared-error loss: start from the flat
// baseline F0(x) = mean(ys), then repeat numStages times: fit a stump to
// the current residual (ys minus the running prediction) and add
// lr*stump(x) into the running prediction. It returns the baseline and the
// sequence of fitted stumps, in order; Predict sums them back up for a
// given x. Pure and deterministic: same inputs always produce the same
// stumps, no randomness -- unlike random-forest's resampling, every stage
// here is a direct function of the stage before it.
func Boost(xs, ys []float64, numStages int, lr float64) (base float64, stumps []Stump) {
	base = Mean(ys)
	preds := make([]float64, len(xs))
	for i := range preds {
		preds[i] = base
	}
	for s := 0; s < numStages; s++ {
		resid := make([]float64, len(xs))
		for i := range xs {
			resid[i] = ys[i] - preds[i]
		}
		st := FitStump(xs, resid)
		stumps = append(stumps, st)
		for i, x := range xs {
			preds[i] += lr * StumpPredict(st, x)
		}
	}
	return base, stumps
}

// Predict returns the boosted ensemble's prediction at x: the baseline plus
// lr times every stump's output, summed in order.
func Predict(base float64, stumps []Stump, lr float64, x float64) float64 {
	pred := base
	for _, st := range stumps {
		pred += lr * StumpPredict(st, x)
	}
	return pred
}

// maxStages is the deepest boosting sequence Render ever needs. Fitting it
// once and slicing prefixes is equivalent to re-running Boost at each
// smaller stage count, since stage k's stump only ever depends on the
// stages before it.
const maxStages = 4

// stairPoints traces the piecewise-constant prediction base+lr*stumps(x)
// across [xlo, xhi] as an explicit staircase: a horizontal segment per
// region between consecutive stump thresholds, joined by vertical risers,
// so the picture shows the same flat-then-jump shape the stumps actually
// produce instead of a smoothed-over interpolation between them.
func stairPoints(base float64, stumps []Stump, lr, xlo, xhi float64) [][2]float64 {
	cuts := make([]float64, 0, len(stumps))
	for _, s := range stumps {
		if s.Threshold > xlo && s.Threshold < xhi {
			cuts = append(cuts, s.Threshold)
		}
	}
	sort.Float64s(cuts)
	dedup := cuts[:0]
	for i, c := range cuts {
		if i == 0 || c != dedup[len(dedup)-1] {
			dedup = append(dedup, c)
		}
	}
	bounds := append([]float64{xlo}, dedup...)
	bounds = append(bounds, xhi)

	pts := make([][2]float64, 0, 2*len(bounds))
	for i := 0; i < len(bounds)-1; i++ {
		lo, hi := bounds[i], bounds[i+1]
		y := Predict(base, stumps, lr, (lo+hi)/2)
		if i > 0 {
			prevY := Predict(base, stumps, lr, (bounds[i-1]+lo)/2)
			pts = append(pts, [2]float64{lo, prevY})
		}
		pts = append(pts, [2]float64{lo, y})
		pts = append(pts, [2]float64{hi, y})
	}
	return pts
}

func render(p map[string]float64) string {
	numStages := int(p["numStages"])
	lr := p["learningRate"]

	base, allStumps := Boost(Xs, Ys, maxStages, lr)
	active := allStumps[:numStages]

	preds := make([]float64, len(Xs))
	for i, x := range Xs {
		preds[i] = Predict(base, active, lr, x)
	}
	sse := SSE(Ys, preds)
	baselineSSE := SSE(Ys, []float64{base, base, base, base, base})

	const xmin, xmax, ymin, ymax = 0.5, 5.5, -1.0, 11.0
	c := viz.New(600, 420, xmin, xmax, ymin, ymax)
	c.PadT = 70
	c.Axes()
	for x := 1.0; x <= 5.0; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	c.Path([][2]float64{{xmin, base}, {xmax, base}}, viz.Muted, 1)
	c.Path(stairPoints(base, active, lr, xmin, xmax), viz.Accent, 2.5)

	for i, x := range Xs {
		px, py := c.X(x), c.Y(Ys[i])
		c.Rect(px-5, py-5, 10, 10, viz.Ink, 1)
	}

	c.Text(16, 22, fmt.Sprintf("stage %d of %d   learning rate = %.1f", numStages, maxStages, lr), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("training SSE = %.3f   (stage 0 baseline alone: %.0f)", sse, baselineSSE), 14, viz.Accent, "start")
	c.Text(16, 64, "black squares = true (x,y)    gray line = flat baseline F0 = mean(y)    accent staircase = cumulative prediction",
		12, viz.Muted, "start")

	return c.String()
}
