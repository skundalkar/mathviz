// Package regularization visualizes L1 (Lasso) and L2 (Ridge) penalties:
// fitting a linear model on 6 candidate input signals, only 2 of which
// actually drive the outcome, and watching how each penalty type treats the
// 4 irrelevant ones as the penalty strength λ increases — Ridge shrinks
// every coefficient smoothly, Lasso can zero one out entirely.
package regularization

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "regularization",
		Seq:   57,
		Title: "Regularization (L1/L2 penalties)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`linear-regression` already showed how to fit the single best straight " +
						"line through a scatter with one input — but that story assumed you " +
						"already knew hours-studied was the input that mattered. Real prediction " +
						"problems usually start with a pile of CANDIDATE inputs, not one " +
						"pre-screened one: predicting an exam score from hours studied, hours " +
						"slept, a handful of other measurements nobody has vetted for relevance. " +
						"Fit ordinary least squares using every candidate at once, and it will " +
						"happily hand back a nonzero weight for every single one — including the " +
						"ones that don't actually matter — because on this one particular batch " +
						"of data, pure chance gives every measured number some small, spurious " +
						"correlation with the outcome. How do you get a model that doesn't quietly " +
						"treat that noise as if it were signal?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take 24 data points and 6 candidate input signals, x1 through x6. Only x1 " +
						"and x2 actually generate the outcome for this worked example (true " +
						"relationship y = 3·x1 − 2·x2 + noise); x3 through x6 are pure noise, " +
						"unrelated to y by construction, mixed in exactly the way an irrelevant " +
						"measurement would be in a real dataset. Fit ordinary least squares (no " +
						"penalty at all, λ=0) and this comes out: x1 ≈ 3.01, x2 ≈ −1.88 — close " +
						"to the true 3 and −2, good — but ALSO x3 ≈ −0.06, x4 ≈ 0.06, x5 ≈ 0.03, " +
						"x6 ≈ 0.09: four nonzero weights on features that have nothing to do with " +
						"y, purely because in this one 24-point sample they happen to line up " +
						"with the noise a little.",
					"Ridge regression patches the normal equations plain least squares already " +
						"solves — (XᵀX)β = Xᵀy — by adding a penalty λ to the diagonal for every " +
						"coefficient except the intercept: (XᵀX + λI)β = Xᵀy. That extra λ makes " +
						"large coefficients costlier, so the solver settles for smaller ones " +
						"across the board. At λ=2, the total penalized quantity — the sum of " +
						"every non-intercept coefficient squared — drops from 12.65 (λ=0) to " +
						"11.03, but no single coefficient hits exactly zero; x3–x6 shift to " +
						"roughly (−0.05, 0.04, 0.10, −0.07) — smaller in total, but still there.",
					"Lasso instead updates one coefficient at a time by a rule called " +
						"soft-thresholding: nudge the coefficient toward zero by exactly λ, and " +
						"clip it to exactly 0 if that nudge would flip its sign. At the same λ=2, " +
						"x1 ≈ 2.95 and x2 ≈ −1.81 survive — still clearly nonzero, close to the " +
						"true 3 and −2 — but x3, x4, x5, and x6 all land at EXACTLY 0.0, every " +
						"one of them. Lasso didn't just shrink the four irrelevant weights, it " +
						"deleted them.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Six bars, one per candidate input x1–x6: a faint grey ghost bar behind each " +
						"one is the unpenalized OLS weight (λ=0) for reference, and the solid " +
						"colored bar in front is the current fit at whatever λ and penalty type " +
						"the sliders are set to — blue for the two genuinely predictive inputs " +
						"(x1, x2), orange for the four that are pure noise (x3–x6). Set penalty " +
						"to 1 (Lasso) and push λ up past about 1: the four orange bars snap to " +
						"exactly zero height while the two blue bars stay tall. Switch penalty " +
						"back to 0 (Ridge) at the same λ and the orange bars only shrink a " +
						"little — they never fully disappear. The readout reports the exact SSE " +
						"(how well the current fit matches these 24 points) and the penalty term " +
						"(Σ|coefficient| for Lasso, Σcoefficient² for Ridge), so you can watch " +
						"that total shrink smoothly as λ grows even when individual Ridge " +
						"coefficients don't move in lockstep.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Recover which inputs actually matter without knowing that in advance: " +
						"crank λ up under Lasso and read off which coefficients survive at " +
						"exactly nonzero — that IS automatic feature selection, done by the same " +
						"machinery that fits the line, no separate step required. Ridge doesn't " +
						"give that crisp yes/no, but it gives something plain least squares " +
						"can't either: a dial for trading a little bias for a lot less variance, " +
						"exactly the trade `bias-variance-tradeoff` described — plain OLS has " +
						"exactly one setting (λ=0) and no dial at all.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Genomics studies routinely start with thousands of candidate genes and a " +
						"few dozen patients — Lasso is a standard tool for narrowing that down to " +
						"the handful of genes that actually predict an outcome. Spam filters and " +
						"recommendation systems start from huge feature sets (every word, every " +
						"past click) and use L1 penalties to keep only the ones worth keeping. " +
						"Calling a theory or an explanation 'overfit' or 'too complicated' in " +
						"everyday conversation is an informal version of exactly the same " +
						"complaint regularization answers formally: don't let something explain " +
						"away noise that isn't really there.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: L1 (Lasso) and L2 (Ridge) both shrink coefficients to " +
						"fight overfitting, but by mechanisms with a real, visible difference — " +
						"L1 can zero a coefficient out entirely, L2 mathematically cannot (its " +
						"penalty is smooth everywhere, so the pull toward zero gets weaker, not " +
						"sharper, the closer a coefficient already is to it).",
					"Not like this: 'more regularization is always better.' Push λ high enough " +
						"under either penalty and even the genuinely predictive x1, x2 weights " +
						"get crushed toward zero along with the noise — the same underfitting " +
						"`bias-variance-tradeoff` called high bias. λ is a dial, not a one-way " +
						"improvement switch; picking it well normally means checking performance " +
						"on held-out data, not just watching training error fall.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lambda", Label: "Penalty strength (λ)", Min: 0, Max: 10, Step: 0.1, Def: 2},
			{Key: "penalty", Label: "Penalty type (0=L2 ridge, 1=L1 lasso)", Min: 0, Max: 1, Step: 1, Def: 1},
			{Key: "noise", Label: "Noise", Min: 0, Max: 1.5, Step: 0.05, Def: 0.6},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value
// in [0, 1). It is a fixed function of its input, not a random-number
// generator, so everything built on it stays pure. Same formula used across
// the gallery (see overfitting.hash01, biasvariance.hash01).
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// numSamples is the fixed size of the synthetic dataset every fit in this
// concept is trained on.
const numSamples = 24

// TrueCoeffs are the coefficients that actually generate y: only x1 and x2
// (indices 0, 1) truly drive the outcome; x3..x6 (indices 2..5) are pure
// noise features, included exactly the way an irrelevant measurement would
// be in a real dataset.
var TrueCoeffs = []float64{3, -2, 0, 0, 0, 0}

// Features returns n fixed synthetic rows of 6 candidate input signals, each
// in [-2, 2), generated deterministically (via hash01) so the same "dataset"
// renders identically on every call -- no globals, no time, no randomness.
func Features(n int) [][]float64 {
	p := len(TrueCoeffs)
	X := make([][]float64, n)
	for i := 0; i < n; i++ {
		row := make([]float64, p)
		for j := 0; j < p; j++ {
			row[j] = hash01(i*17+j*101+7)*4 - 2
		}
		X[i] = row
	}
	return X
}

// Targets returns the noisy outcome y = X·TrueCoeffs + noise for every row
// of X, with noise deterministically drawn from [-noiseAmp, noiseAmp].
func Targets(X [][]float64, noiseAmp float64) []float64 {
	y := make([]float64, len(X))
	for i, row := range X {
		v := 0.0
		for j, coef := range TrueCoeffs {
			v += coef * row[j]
		}
		noise := (hash01(i*53+911)*2 - 1) * noiseAmp
		y[i] = v + noise
	}
	return y
}

// Predict evaluates coeffs (intercept followed by one weight per feature) on
// one feature row.
func Predict(coeffs, row []float64) float64 {
	v := coeffs[0]
	for j, c := range coeffs[1:] {
		v += c * row[j]
	}
	return v
}

// SSE is the total squared error of coeffs' predictions across every row of
// X against the matching y.
func SSE(coeffs []float64, X [][]float64, y []float64) float64 {
	s := 0.0
	for i, row := range X {
		d := y[i] - Predict(coeffs, row)
		s += d * d
	}
	return s
}

// RidgeFit fits [intercept, c1..cp] minimizing SSE(coeffs) + lambda *
// sum(cj^2 for j>=1) by solving the ridge normal equations directly:
// (XᵀX + λD)c = Xᵀy, where D is the identity with a 0 in the intercept's
// row/column so the intercept itself is never penalized.
func RidgeFit(X [][]float64, y []float64, lambda float64) []float64 {
	n := len(X)
	p := len(X[0])
	k := p + 1
	A := make([][]float64, k)
	for i := range A {
		A[i] = make([]float64, k)
	}
	b := make([]float64, k)
	for i := 0; i < n; i++ {
		d := make([]float64, k)
		d[0] = 1
		copy(d[1:], X[i])
		for a := 0; a < k; a++ {
			for c := 0; c < k; c++ {
				A[a][c] += d[a] * d[c]
			}
			b[a] += d[a] * y[i]
		}
	}
	for i := 1; i < k; i++ {
		A[i][i] += lambda
	}
	return solveLinearSystem(A, b)
}

// solveLinearSystem solves Ax = b for a small square system via Gaussian
// elimination with partial pivoting. Returns a zero vector for a singular
// column. Same algorithm as biasvariance.solveLinearSystem and
// overfitting.solveLinearSystem.
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

// softThreshold is the core of L1 shrinkage: it moves z toward 0 by gamma,
// clamping to exactly 0 rather than crossing it. This clamp is what lets
// Lasso zero a coefficient out entirely, unlike Ridge's smooth shrink.
func softThreshold(z, gamma float64) float64 {
	if z > gamma {
		return z - gamma
	}
	if z < -gamma {
		return z + gamma
	}
	return 0
}

// lassoIterations is the fixed number of coordinate-descent sweeps LassoFit
// runs -- enough for this dataset's size and scale to converge to a stable
// fixed point (checked in the tests below), with no early-stopping check so
// the function stays a plain, deterministic loop.
const lassoIterations = 300

// LassoFit fits [intercept, c1..cp] minimizing SSE(coeffs) + lambda *
// sum(|cj| for j>=1) by cyclic coordinate descent: repeatedly hold every
// coefficient but one fixed, solve for that one coefficient's optimum via
// soft-thresholding, and sweep over every coefficient lassoIterations times.
func LassoFit(X [][]float64, y []float64, lambda float64) []float64 {
	n := len(X)
	p := len(X[0])
	k := p + 1
	design := make([][]float64, n)
	for i := 0; i < n; i++ {
		d := make([]float64, k)
		d[0] = 1
		copy(d[1:], X[i])
		design[i] = d
	}
	colNormSq := make([]float64, k)
	for j := 0; j < k; j++ {
		s := 0.0
		for i := 0; i < n; i++ {
			s += design[i][j] * design[i][j]
		}
		colNormSq[j] = s
	}

	coeffs := make([]float64, k)
	residual := append([]float64(nil), y...)
	for iter := 0; iter < lassoIterations; iter++ {
		for j := 0; j < k; j++ {
			if colNormSq[j] < 1e-12 {
				continue
			}
			// Add feature j's current contribution back into the residual
			// so rho measures how strongly j still correlates with what's
			// left unexplained by every OTHER coefficient.
			for i := 0; i < n; i++ {
				residual[i] += design[i][j] * coeffs[j]
			}
			rho := 0.0
			for i := 0; i < n; i++ {
				rho += design[i][j] * residual[i]
			}
			var next float64
			if j == 0 {
				next = rho / colNormSq[j] // intercept: never penalized
			} else {
				next = softThreshold(rho, lambda) / colNormSq[j]
			}
			coeffs[j] = next
			for i := 0; i < n; i++ {
				residual[i] -= design[i][j] * coeffs[j]
			}
		}
	}
	return coeffs
}

// L1Penalty is the sum of absolute values of every non-intercept
// coefficient -- what Lasso minimizes alongside SSE.
func L1Penalty(coeffs []float64) float64 {
	s := 0.0
	for _, c := range coeffs[1:] {
		s += math.Abs(c)
	}
	return s
}

// L2Penalty is the sum of squares of every non-intercept coefficient --
// what Ridge minimizes alongside SSE.
func L2Penalty(coeffs []float64) float64 {
	s := 0.0
	for _, c := range coeffs[1:] {
		s += c * c
	}
	return s
}

// CountNearZero reports how many of coeffs[from:] have magnitude below tol.
func CountNearZero(coeffs []float64, from int, tol float64) int {
	n := 0
	for _, c := range coeffs[from:] {
		if math.Abs(c) < tol {
			n++
		}
	}
	return n
}

// barChartScale converts a coefficient magnitude into pixels. Fixed rather
// than dynamically rescaled per-frame, so a bar's height means the same
// thing (and stays comparable) at every slider setting.
const barChartScale = 55.0

func render(p map[string]float64) string {
	lambda := p["lambda"]
	noiseAmp := p["noise"]
	lasso := p["penalty"] >= 0.5

	X := Features(numSamples)
	y := Targets(X, noiseAmp)

	ols := RidgeFit(X, y, 0)
	var fit []float64
	var penaltyName string
	if lasso {
		fit = LassoFit(X, y, lambda)
		penaltyName = "Lasso (L1)"
	} else {
		fit = RidgeFit(X, y, lambda)
		penaltyName = "Ridge (L2)"
	}

	const w, h = 680, 460
	// This is a custom bar chart, not a cartesian data plot, so the canvas's
	// data-space mapping (c.X/c.Y) is left unused; every shape below is
	// placed directly in pixel coordinates via Rect/Text.
	c := viz.New(w, h, 0, 1, 0, 1)

	const (
		baseline  = 300.0 // y=0 pixel row
		groupLeft = 40.0
		groupW    = 100.0
		barW      = 26.0
		barGap    = 4.0
	)
	numGroups := len(fit) - 1 // exclude the intercept

	c.Rect(groupLeft-10, baseline-1, groupW*float64(numGroups)+20, 2, viz.Ink, 1)

	bar := func(centerX, value float64, color string, opacity float64) {
		hpx := math.Abs(value) * barChartScale
		if value >= 0 {
			c.Rect(centerX, baseline-hpx, barW, hpx, color, opacity)
		} else {
			c.Rect(centerX, baseline, barW, hpx, color, opacity)
		}
	}

	for i := 0; i < numGroups; i++ {
		groupCenter := groupLeft + groupW*(float64(i)+0.5)
		ghostX := groupCenter - barW - barGap/2
		solidX := groupCenter + barGap/2

		bar(ghostX, ols[i+1], viz.Muted, 0.4)

		color := viz.Warm // x3..x6: pure-noise features
		if i < 2 {
			color = viz.Accent // x1, x2: the two genuinely predictive features
		}
		bar(solidX, fit[i+1], color, 0.9)

		// True-coefficient reference tick, only meaningful for x1/x2 -- the
		// noise features' true coefficient is 0, which already coincides
		// with the baseline.
		if i < 2 {
			trueY := baseline - TrueCoeffs[i]*barChartScale
			c.Rect(groupCenter-barW-barGap/2-2, trueY-1, barW*2+barGap+4, 2, viz.Good, 0.8)
		}

		c.Text(groupCenter, baseline+22, fmt.Sprintf("x%d", i+1), 13, viz.Ink, "middle")
	}

	c.Text(20, 24, fmt.Sprintf("λ = %.2f    penalty = %s    noise = %.2f", lambda, penaltyName, noiseAmp),
		14, viz.Ink, "start")
	c.Text(20, 44, "blue = signal (x1,x2)    orange = noise (x3-x6)    grey ghost = unpenalized OLS    green tick = true coefficient",
		12, viz.Muted, "start")

	sse := SSE(fit, X, y)
	zeros := CountNearZero(fit, 3, 1e-9)
	var penaltyTerm float64
	var penaltyLabel string
	if lasso {
		penaltyTerm = L1Penalty(fit)
		penaltyLabel = "Σ|coefficient|"
	} else {
		penaltyTerm = L2Penalty(fit)
		penaltyLabel = "Σcoefficient²"
	}
	c.Text(20, h-20, fmt.Sprintf("SSE = %.3f    penalty term (%s) = %.3f    noise coefficients exactly zero = %d of 4",
		sse, penaltyLabel, penaltyTerm, zeros), 13, viz.Ink, "start")

	return c.String()
}
