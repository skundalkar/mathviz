// Package biasvariance visualizes the bias-variance decomposition:
// retraining the same-degree polynomial on many independently resampled
// noisy datasets and looking at how its predictions at one fixed point
// behave, not just once (as `overfitting` does) but across many retrainings
// — separating "consistently wrong" (bias) from "unpredictable" (variance).
package biasvariance

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bias-variance-tradeoff",
		Seq:   50,
		Title: "Bias-variance tradeoff",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`overfitting` already showed that a high-degree curve chases noise instead " +
						"of the true pattern — but that was one look at one noisy dataset. " +
						"Retrain that exact same-degree model on a DIFFERENT random sample of the " +
						"same underlying pattern, and does it land in the same wrong place every " +
						"time, or does it land somewhere different each time? Those are two very " +
						"different kinds of 'wrong,' and a single-dataset picture can't tell them " +
						"apart — is there a way to separate 'the model is consistently biased " +
						"toward the wrong answer' from 'the model is unpredictable, landing " +
						"somewhere different every time it's retrained'?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Same setup as `overfitting`: 12 noisy points sampled from the true pattern " +
						"sin(x), a degree-d polynomial fit by least squares. But now retrain that " +
						"same-degree model on 30 independently resampled noisy datasets (same true " +
						"pattern, fresh noise each time), and look at all 30 fitted curves' " +
						"predictions at one fixed point, x0=1.5 (true value sin(1.5)≈0.997):",
					"• degree=1 (straight line): the 30 predictions cluster tightly together " +
						"(variance≈0.0076 — barely any spread run to run) but their average, " +
						"0.321, sits far from the true 0.997 — bias²≈0.458. A straight line " +
						"structurally can't bend to follow sin's curve near its peak, so it's " +
						"consistently, predictably wrong, no matter which noisy sample it's " +
						"retrained on.",
					"• degree=9 (very flexible): the 30 predictions scatter widely " +
						"(variance≈0.033) but their average, 0.977, lands close to the true 0.997 " +
						"— bias²≈0.0004. Flexible enough to track the true curve on average, but " +
						"exactly which quirks of noise it latches onto varies wildly from one " +
						"training sample to the next.",
					"• degree=3 (moderate): bias²≈0.0026 and variance≈0.0188 — both small, and " +
						"their sum, 0.0213, beats both extremes (0.465 for degree 1, 0.034 for " +
						"degree 9).",
					"Total expected error decomposes exactly into these two pieces: " +
						"E[(prediction - true)²] = bias² + variance. Degree 1 fails almost " +
						"entirely from bias; degree 9 fails almost entirely from variance; " +
						"degree 3 balances both to land at the lowest total.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"degree sets model complexity, the same slider as `overfitting`; noise sets " +
						"how much each training sample's y-values scatter from the true pattern; " +
						"x0 moves the vertical evaluation line where bias and variance are " +
						"measured. Every faint curve is one polynomial fit to one independently " +
						"resampled noisy dataset — 30 of them drawn together as a 'spaghetti " +
						"plot.' At low degree, all 30 strands bunch tightly on top of each other " +
						"but visibly miss the true dashed curve wherever it bends the most (high " +
						"bias, low variance). At high degree, the strands fan out wildly, each " +
						"chasing its own dataset's particular noise (low bias, high variance). The " +
						"dots mark where each strand crosses the x0 line — their spread is the " +
						"variance number in the readout, and how far their average sits from the " +
						"true curve is the bias.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Diagnose why a model is underperforming, not just that it is: a model with " +
						"high bias and low variance is consistently missing the pattern, and needs " +
						"more flexibility (more features, higher degree, a bigger model) to fix; a " +
						"model with low bias and high variance is capable of finding the pattern " +
						"but is unstable, and needs regularization, more training data, or a " +
						"simpler model to fix — opposite prescriptions for what looks, from " +
						"training accuracy alone, like the same kind of failure. Sweep the degree " +
						"slider from 1 to 11 and watch total error trace a U-shape, exactly the " +
						"sweet-spot degree `overfitting` gestured at without naming what was " +
						"actually trading off underneath it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A linear regression model on a genuinely curved relationship (like " +
						"`overfitting`'s straight-line problem) is high-bias: retrain it on new " +
						"data all year and it keeps making the same kind of systematically wrong " +
						"prediction. A deep decision tree with no depth limit is high-variance: " +
						"retrain it on a slightly different sample of the same customers and it " +
						"can produce a noticeably different tree with different predictions, even " +
						"though the underlying pattern didn't change. 'Ensemble' methods exploit " +
						"this asymmetry directly: bagging (as in random forests) averages many " +
						"high-variance, low-bias trees together, and averaging cancels out " +
						"variance — the same reason `law-of-large-numbers` shrinks toward the " +
						"true mean — while leaving each tree's low bias intact.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: bias and variance are two separate failure modes, and the " +
						"fix for one often makes the other worse — a simpler model reduces " +
						"variance but tends to raise bias; a more flexible model reduces bias but " +
						"tends to raise variance, exactly what the degree slider demonstrates " +
						"directly.",
					"Not like this: assuming 'the model has high error' by itself tells you which " +
						"one is the problem, or that adding more training data always helps — more " +
						"data mainly shrinks variance (each individual fit has less noise to " +
						"overfit to), so it does little for a model whose real problem is high " +
						"bias. Degree 1 in this picture stays badly biased no matter how many " +
						"resampled datasets you average over, because a straight line simply " +
						"can't bend.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "degree", Label: "Model complexity (degree)", Min: 1, Max: 11, Step: 1, Def: 3},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1, Step: 0.05, Def: 0.4},
			{Key: "x0", Label: "Evaluation point (x0)", Min: -3.2, Max: 3.2, Step: 0.2, Def: 1.5},
		},
		Render: render,
	})
}

// TrueFunction is the smooth underlying pattern the noisy data is secretly
// following. Identical to overfitting.TrueFunction — each concept package
// is self-contained (see BUILD_CYCLE.md).
func TrueFunction(x float64) float64 {
	return math.Sin(x)
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value
// in [0, 1). It is a fixed function of its input, not a random-number
// generator, so everything built on it stays pure. Identical to
// overfitting.hash01.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// numPoints is the fixed size of each resampled "practice problem" dataset,
// same as overfitting.numPoints.
const numPoints = 12

// DataPoints returns n fixed, evenly-spaced x's over [-3.2, 3.2] and their
// noisy y's for one independently resampled training run: TrueFunction(x)
// plus deterministic pseudo-random noise in [-noiseAmp, noiseAmp]. Unlike
// overfitting.DataPoints, run selects which of many independent noisy
// resamplings to generate — each run uses a disjoint slice of the hash01
// seed space, so different runs are different fixed "datasets," not the
// same one.
func DataPoints(n int, noiseAmp float64, run int) (xs, ys []float64) {
	if n < 2 {
		n = 2
	}
	const lo, hi = -3.2, 3.2
	xs = make([]float64, n)
	ys = make([]float64, n)
	for i := 0; i < n; i++ {
		x := lo + (hi-lo)*float64(i)/float64(n-1)
		noise := (hash01(run*97+i+1)*2 - 1) * noiseAmp
		xs[i] = x
		ys[i] = TrueFunction(x) + noise
	}
	return xs, ys
}

// PolyFit fits a degree-d polynomial y ≈ c0 + c1*x + ... + cd*x^d to (xs, ys)
// by ordinary least squares, returning the coefficients c0..cd. Identical
// algorithm to overfitting.PolyFit (normal equations via Gaussian
// elimination with partial pivoting).
func PolyFit(xs, ys []float64, degree int) []float64 {
	n := len(xs)
	if degree < 0 {
		degree = 0
	}
	if degree > n-1 {
		degree = n - 1 // can't fit more free parameters than data points
	}
	k := degree + 1

	A := make([][]float64, k)
	b := make([]float64, k)
	for i := 0; i < k; i++ {
		A[i] = make([]float64, k)
	}
	for _, x := range xs {
		pow := make([]float64, 2*k-1)
		pow[0] = 1
		for p := 1; p < len(pow); p++ {
			pow[p] = pow[p-1] * x
		}
		for i := 0; i < k; i++ {
			for j := 0; j < k; j++ {
				A[i][j] += pow[i+j]
			}
		}
	}
	for idx, x := range xs {
		xp := 1.0
		for i := 0; i < k; i++ {
			b[i] += ys[idx] * xp
			xp *= x
		}
	}

	return solveLinearSystem(A, b)
}

// solveLinearSystem solves Ax = b for a small square system via Gaussian
// elimination with partial pivoting. Returns a zero vector if A is
// singular. Identical to overfitting.solveLinearSystem.
func solveLinearSystem(A [][]float64, b []float64) []float64 {
	k := len(b)
	m := make([][]float64, k)
	for i := range m {
		m[i] = append([]float64(nil), A[i]...)
	}
	v := append([]float64(nil), b...)

	for col := 0; col < k; col++ {
		pivot := col
		for r := col + 1; r < k; r++ {
			if math.Abs(m[r][col]) > math.Abs(m[pivot][col]) {
				pivot = r
			}
		}
		m[col], m[pivot] = m[pivot], m[col]
		v[col], v[pivot] = v[pivot], v[col]

		if math.Abs(m[col][col]) < 1e-12 {
			continue // singular in this column; leave the coefficient at 0
		}
		for r := col + 1; r < k; r++ {
			f := m[r][col] / m[col][col]
			for c := col; c < k; c++ {
				m[r][c] -= f * m[col][c]
			}
			v[r] -= f * v[col]
		}
	}

	x := make([]float64, k)
	for i := k - 1; i >= 0; i-- {
		sum := v[i]
		for j := i + 1; j < k; j++ {
			sum -= m[i][j] * x[j]
		}
		if math.Abs(m[i][i]) < 1e-12 {
			x[i] = 0
			continue
		}
		x[i] = sum / m[i][i]
	}
	return x
}

// PolyEval evaluates a polynomial with the given coefficients (c0..cd, low
// degree first) at x. Identical to overfitting.PolyEval.
func PolyEval(coeffs []float64, x float64) float64 {
	result := 0.0
	xp := 1.0
	for _, c := range coeffs {
		result += c * xp
		xp *= x
	}
	return result
}

// Mean returns the arithmetic mean of xs, or 0 for an empty slice.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// numRuns is how many independently resampled noisy datasets get retrained
// on, per degree, to estimate bias and variance -- enough runs for the
// spread's shape to read clearly without the SVG growing unreasonably large.
const numRuns = 30

// Predictions retrains a degree-d polynomial on numRuns independently
// resampled noisy datasets (same true pattern, fresh noise each run) and
// returns each retraining's prediction at the fixed evaluation point x0.
func Predictions(x0 float64, degree int, noiseAmp float64) []float64 {
	preds := make([]float64, numRuns)
	for r := 0; r < numRuns; r++ {
		xs, ys := DataPoints(numPoints, noiseAmp, r)
		coeffs := PolyFit(xs, ys, degree)
		preds[r] = PolyEval(coeffs, x0)
	}
	return preds
}

// BiasSquared returns (mean(preds) - trueVal)^2: how far the AVERAGE
// prediction across many retrainings sits from the true value, squared.
// High bias means the model is consistently wrong in the same direction,
// regardless of which particular noisy dataset it was trained on.
func BiasSquared(preds []float64, trueVal float64) float64 {
	d := Mean(preds) - trueVal
	return d * d
}

// Variance returns the average squared deviation of preds from their own
// mean: how much predictions swing from one retraining to the next, with no
// reference to the true value at all. High variance means the model is
// unstable -- it lands somewhere different depending on which particular
// noisy dataset it happened to see.
func Variance(preds []float64) float64 {
	if len(preds) == 0 {
		return 0
	}
	m := Mean(preds)
	sum := 0.0
	for _, p := range preds {
		d := p - m
		sum += d * d
	}
	return sum / float64(len(preds))
}

func render(params map[string]float64) string {
	degree, noise, x0 := int(params["degree"]), params["noise"], params["x0"]
	if x0 < -3.2 {
		x0 = -3.2
	}
	if x0 > 3.2 {
		x0 = 3.2
	}

	const xmin, xmax = -3.5, 3.5
	const ymin, ymax = -2.5, 2.5
	c := viz.New(680, 460, xmin, xmax, ymin, ymax)
	c.PadT = 100
	c.Axes()
	for x := -3.0; x <= 3.0; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	trueCurve := viz.Sample(xmin, xmax, 240, TrueFunction)
	c.Path(trueCurve, viz.Muted, 1.5)

	// One faint curve per independently resampled training run -- the
	// "spaghetti plot." Low degrees bunch tightly (low variance); high
	// degrees fan out wildly (high variance).
	preds := make([]float64, numRuns)
	for r := 0; r < numRuns; r++ {
		xs, ys := DataPoints(numPoints, noise, r)
		coeffs := PolyFit(xs, ys, degree)
		fitCurve := viz.Sample(xmin, xmax, 120, func(x float64) float64 {
			return PolyEval(coeffs, x)
		})
		c.Path(fitCurve, viz.Accent, 1)
		preds[r] = PolyEval(coeffs, x0)
	}

	c.VLine(x0, viz.Warm, true)
	for _, pred := range preds {
		if pred < ymin || pred > ymax {
			continue
		}
		px, py := c.X(x0), c.Y(pred)
		c.Rect(px-2, py-2, 4, 4, viz.Warm, 0.7)
	}

	trueVal := TrueFunction(x0)
	bias2 := BiasSquared(preds, trueVal)
	vari := Variance(preds)

	c.Text(20, 24, fmt.Sprintf("degree = %d    noise = %.2f    x0 = %.1f", degree, noise, x0),
		14, viz.Ink, "start")
	c.Text(20, 46, fmt.Sprintf("true value at x0 = %.3f    mean prediction = %.3f", trueVal, Mean(preds)),
		13, viz.Muted, "start")
	c.Text(20, 68, fmt.Sprintf("bias² ≈ %.4f    variance ≈ %.4f    total ≈ %.4f", bias2, vari, bias2+vari),
		14, viz.Ink, "start")
	c.Text(20, 88, "grey = true pattern    faint blue = 30 resampled fits    orange dots = their predictions at x0",
		12, viz.Muted, "start")

	return c.String()
}
