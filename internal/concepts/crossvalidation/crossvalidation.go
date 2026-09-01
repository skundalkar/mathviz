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
				Body:    []string{"placeholder"},
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
