// Package crossvalidation visualizes k-fold cross-validation: instead of
// judging a model on one train/test split, rotate which slice of the data is
// held out, score the model on each slice in turn, and average the scores.
// The picture drives home why that average is a far more honest number than
// whatever a single lucky (or unlucky) split happens to hand you.
package crossvalidation

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "cross-validation",
		Seq:   85,
		Title: "Cross-validation (rotating the held-out fold)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You fit a line to 10 data points, hold out 2 of them as a \"test set\" the " +
						"model never saw during training, and check how far off the predictions " +
						"are. The number comes back tiny — the model looks great. Ship it? Here's " +
						"the uncomfortable question: would you have gotten the same verdict if " +
						"you'd held out a *different* 2 points instead? If the answer is no — if " +
						"which 2 points you happened to pick changes the story from \"great model\" " +
						"to \"terrible model\" — then the number you just trusted was measuring " +
						"luck, not the model.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take 10 points with a roughly-linear trend, y ≈ 2x, except one deliberately " +
						"planted outlier: at x=8 the trend predicts ~16, but the actual value is " +
						"24. Instead of picking one train/test split, k-fold cross-validation " +
						"(here k=5) splits all 10 points into 5 folds of 2 points each, in order: " +
						"{x=1,2}, {x=3,4}, {x=5,6}, {x=7,8}, {x=9,10}. Then it rotates: for each " +
						"fold in turn, fit a least-squares line on the other 8 points and measure " +
						"the mean squared error (MSE) on just that fold's 2 held-out points.",
					"• Fold {1,2} held out: trained on the other 8 points, slope≈2.29, " +
						"intercept≈-0.86 → predicts [1.43, 3.73] against actual [2.0, 4.1] → " +
						"MSE=0.23.",
					"• Fold {3,4} held out: MSE=0.19.",
					"• Fold {5,6} held out: MSE=1.31.",
					"• Fold {7,8} held out — this is the outlier's fold: trained on the other 8 " +
						"points (which never include the outlier), slope=2.00, intercept≈0.03 → " +
						"predicts a clean [14.02, 16.02] against actual [14.2, 24.0] → MSE=31.82, " +
						"over 100× fold {1,2}'s error.",
					"• Fold {9,10} held out: this time the outlier IS in the training data, and " +
						"it drags the fitted line off course — slope≈2.67, intercept≈-1.97 → " +
						"predicts [22.04, 24.71] against actual [18.0, 20.1] → MSE=18.80. Neither " +
						"held-out point is unusual; the fold still scores badly because the " +
						"outlier corrupted the line that had to predict them.",
					"Read the five scores together, not one at a time: they range from 0.19 to " +
						"31.82 — a roughly 165× spread — for the exact same model on the exact " +
						"same data, depending only on which 2 points got held out. Averaging all " +
						"five gives (0.23+0.19+1.31+31.82+18.80)/5 = 10.47, the k-fold " +
						"cross-validation estimate: one number that isn't hostage to any single " +
						"split's luck, plus a spread across folds that's informative in its own " +
						"right (it's telling you this model handles most of the data fine but " +
						"struggles badly near x=7–8).",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The \"Number of folds (k)\" slider re-splits the same 10 points into k " +
						"folds; \"Held-out fold\" picks which one is currently drawn as the test " +
						"set (it wraps around if you set it past the current k). Orange squares " +
						"are the held-out fold, black squares are the training points, the blue " +
						"line is fit only on the training points, and the red lines are residuals " +
						"— shown only for the held-out points, since those are the only ones the " +
						"model didn't get to see. The top readout shows the current fold's test " +
						"MSE and the running average across all k folds — the CV estimate. The " +
						"default (k=5, fold 3) lands directly on the outlier fold from the worked " +
						"example above: drag through folds 0-4 and watch the test MSE jump from " +
						"0.19-1.31 on the clean folds to 31.82 the moment the outlier's fold is " +
						"held out. The bottom line freezes the k=5, fold-0-vs-fold-3 comparison " +
						"regardless of the slider position, so the \"lucky vs. unlucky split\" " +
						"contrast from the next section is always visible.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Report \"cross-validated MSE ≈ 10.47\" instead of a single, split-dependent " +
						"number — a score that reflects how the model does averaged across many " +
						"different held-out slices, not whichever slice you happened to pick. You " +
						"also get the per-fold spread for free, which tells you *where* the model " +
						"struggles (here: the region around the outlier), something a single " +
						"train/test split can't reveal even in principle, since it only ever shows " +
						"you one slice.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Any time you see \"cross-validated accuracy\" reported for a machine-learning " +
						"model, this is the literal mechanism, not a figure of speech. The everyday " +
						"version: judging a restaurant off one visit vs. several visits spread " +
						"across different nights — a place that's fantastic on a slow Tuesday but " +
						"chaotic on a packed Saturday is exactly a model that aces one fold and " +
						"bombs another, and only rotating through several visits (folds) tells you " +
						"which is closer to the truth. Same idea grading a student off one exam " +
						"vs. their average across several — one test can land unusually easy or " +
						"unusually hard for that particular student.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: \"cross-validation reports the average test score across k " +
						"rotating splits, so no single split's luck decides the verdict — and the " +
						"spread across folds is itself informative, not noise to average away.\" " +
						"Not like this: trusting a low error from any one split as proof the model " +
						"is good — fold {1,2} above scores MSE=0.23 and fold {7,8} scores " +
						"MSE=31.82 for the identical model on the identical data; a single split " +
						"could have shown you either number and you'd have no way to tell which " +
						"one you got. Also not like this: assuming cross-validation removes the " +
						"need for a separate final test set — CV is for comparing and tuning " +
						"models using the training data's own rotation; once you've picked a model " +
						"this way, you still want one untouched holdout set, never used in any " +
						"fold, for the final honest number.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "k", Label: "Number of folds (k)", Min: 2, Max: 10, Step: 1, Def: 5},
			{Key: "highlightFold", Label: "Held-out fold", Min: 0, Max: 9, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// DataX and DataY are the fixed 10-point worked example every Section walks
// through: a roughly-linear trend (y ≈ 2x) with one deliberate outlier at
// x=8 (y=24 instead of the ~16 the trend predicts) so that which points a
// split happens to hold out visibly changes the verdict.
var (
	DataX = []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	DataY = []float64{2.0, 4.1, 5.9, 8.2, 10.1, 11.8, 14.2, 24.0, 18.0, 20.1}
)

// Mean returns the arithmetic mean of v.
func Mean(v []float64) float64 {
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

// Slope returns the least-squares slope of the line fit to (xs, ys).
func Slope(xs, ys []float64) float64 {
	xBar, yBar := Mean(xs), Mean(ys)
	var sxy, sxx float64
	for i := range xs {
		dx := xs[i] - xBar
		sxy += dx * (ys[i] - yBar)
		sxx += dx * dx
	}
	return sxy / sxx
}

// Intercept returns the least-squares intercept of the line fit to (xs, ys).
func Intercept(xs, ys []float64) float64 {
	return Mean(ys) - Slope(xs, ys)*Mean(xs)
}

// Predict evaluates a line (given its slope and intercept) at x.
func Predict(slope, intercept, x float64) float64 {
	return slope*x + intercept
}

// MSE is the mean (not sum) of the squared residuals of a line (slope,
// intercept) across every (xs[i], ys[i]) pair — the quantity a fold's test
// score is measured in.
func MSE(xs, ys []float64, slope, intercept float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	var sse float64
	for i := range xs {
		r := ys[i] - Predict(slope, intercept, xs[i])
		sse += r * r
	}
	return sse / float64(len(xs))
}

// KFoldIndices partitions the indices [0, n) into k folds as evenly as
// possible: fold f gets the contiguous slice [f*n/k, (f+1)*n/k). With n not
// evenly divisible by k, the last folds absorb the remainder (n=10, k=3
// gives fold sizes 3, 3, 4).
func KFoldIndices(n, k int) [][]int {
	folds := make([][]int, k)
	for f := 0; f < k; f++ {
		lo := f * n / k
		hi := (f + 1) * n / k
		fold := make([]int, 0, hi-lo)
		for i := lo; i < hi; i++ {
			fold = append(fold, i)
		}
		folds[f] = fold
	}
	return folds
}

// FoldFit fits a least-squares line on every (xs[i], ys[i]) pair EXCEPT the
// indices in folds[foldIdx], then scores that line's MSE on exactly the
// held-out fold's points — the core operation of one round of
// cross-validation: train on the rest, test on what you held out.
func FoldFit(xs, ys []float64, folds [][]int, foldIdx int) (slope, intercept, mse float64) {
	test := map[int]bool{}
	for _, i := range folds[foldIdx] {
		test[i] = true
	}
	var trainX, trainY, testX, testY []float64
	for i := range xs {
		if test[i] {
			testX = append(testX, xs[i])
			testY = append(testY, ys[i])
		} else {
			trainX = append(trainX, xs[i])
			trainY = append(trainY, ys[i])
		}
	}
	slope, intercept = Slope(trainX, trainY), Intercept(trainX, trainY)
	mse = MSE(testX, testY, slope, intercept)
	return slope, intercept, mse
}

// CrossValidate runs full k-fold cross-validation: for every fold in turn,
// fit on the rest and score on the held-out fold, then average the k test
// scores into a single CV estimate.
func CrossValidate(xs, ys []float64, k int) (foldMSEs []float64, avg float64) {
	folds := KFoldIndices(len(xs), k)
	foldMSEs = make([]float64, k)
	for f := range folds {
		_, _, mse := FoldFit(xs, ys, folds, f)
		foldMSEs[f] = mse
	}
	avg = Mean(foldMSEs)
	return foldMSEs, avg
}

func render(p map[string]float64) string {
	k := int(p["k"])
	if k < 2 {
		k = 2
	}
	if k > len(DataX) {
		k = len(DataX)
	}
	folds := KFoldIndices(len(DataX), k)
	hf := int(p["highlightFold"]) % k
	if hf < 0 {
		hf += k
	}

	slope, intercept, mse := FoldFit(DataX, DataY, folds, hf)
	_, avg := CrossValidate(DataX, DataY, k)

	// Fixed k=5 reference numbers for the "a single split can mislead"
	// callout at the bottom, independent of whatever k the slider is
	// currently set to — fold 0 (test={x=1,2}) is the lucky split, fold 3
	// (test={x=7,8}) is the one that lands on the deliberate outlier.
	refFolds := KFoldIndices(len(DataX), 5)
	_, _, cleanMSE := FoldFit(DataX, DataY, refFolds, 0)
	_, _, outlierMSE := FoldFit(DataX, DataY, refFolds, 3)

	testSet := map[int]bool{}
	for _, i := range folds[hf] {
		testSet[i] = true
	}

	const xmin, xmax = 0.0, 11.0
	const ymin, ymax = 0.0, 27.0
	c := viz.New(680, 480, xmin, xmax, ymin, ymax)
	c.PadT = 84
	c.Axes()
	for x := 0.0; x <= xmax; x += 2 {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	line := viz.Sample(xmin, xmax, 2, func(x float64) float64 { return Predict(slope, intercept, x) })

	// Residuals for the held-out (test) points only — training points get
	// no residual line since the model was fit to match them closely.
	for i := range DataX {
		if !testSet[i] {
			continue
		}
		predY := Predict(slope, intercept, DataX[i])
		c.Path([][2]float64{{DataX[i], DataY[i]}, {DataX[i], predY}}, viz.Bad, 1.5)
	}

	c.Path(line, viz.Accent, 2.5)

	for i := range DataX {
		px, py := c.X(DataX[i]), c.Y(DataY[i])
		color := viz.Ink
		if testSet[i] {
			color = viz.Warm
		}
		c.Rect(px-4, py-4, 8, 8, color, 0.9)
	}

	c.Text(20, 22, fmt.Sprintf("k=%d folds — held-out fold %d of %d: test MSE = %.2f", k, hf, k-1, mse), 14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("line trained without the held-out points: y = %.2f + %.2f·x", intercept, slope), 13, viz.Muted, "start")
	c.Text(20, 64, fmt.Sprintf("average test MSE across all %d folds (the CV estimate) = %.2f", k, avg), 14, viz.Good, "start")
	c.Text(20, 440, "orange squares = held-out fold    black squares = training points    red = residuals on the held-out fold", 11, viz.Muted, "start")
	c.Text(20, 460, fmt.Sprintf("a single split can mislead: test={x=1,2} -> MSE=%.2f (looks great) vs test={x=7,8} -> MSE=%.2f (looks terrible) -- same model, same data", cleanMSE, outlierMSE), 11, viz.Muted, "start")

	return c.String()
}
