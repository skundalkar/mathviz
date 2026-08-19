// Package pca visualizes Principal Component Analysis: given a cloud of
// correlated points, it finds the one direction the cloud spreads out along
// the most (not necessarily the x-axis or the y-axis — those were just an
// arbitrary choice of how to draw the plot) and reports how much of the
// cloud's total variance that single direction accounts for.
package pca

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "principal-component-analysis",
		Seq:   48,
		Title: "Principal Component Analysis (the direction that explains the most)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`covariance` and `correlation` only ever describe two variables against " +
						"each other, and even then the answer depends on which two axes you " +
						"happened to draw the scatter on. Take a cloud of (hours studied, test " +
						"score) points stretched out diagonally, leaning toward neither axis — the " +
						"x-spread alone (variance of hours) and the y-spread alone (variance of " +
						"score) each miss most of the story, because the cloud isn't actually " +
						"spread out along x or along y, it's spread out along its own diagonal. " +
						"`eigenvectors-eigenvalues` already hinted that a covariance matrix has " +
						"'special directions that only stretch' — but which direction is that, for " +
						"an actual cloud of data, and how much of the cloud's total spread does " +
						"that one direction alone account for?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Same five students as `covariance`: hours studied x = 1,2,3,4,5, test score " +
						"(in tens of points) y = 2,4,5,4,5. Their covariance matrix collects three " +
						"numbers into one 2x2 table — variance of x, variance of y, and their " +
						"covariance filling both off-diagonal slots (the matrix is symmetric, since " +
						"Cov(x,y) = Cov(y,x)): [[Var(x), Cov(x,y)], [Cov(x,y), Var(y)]] = " +
						"[[2, 1.2], [1.2, 1.2]] — Var(x)=2 and Cov(x,y)=1.2 exactly as computed in " +
						"`covariance`'s worked example, plus Var(y)=1.2.",
					"Feed that matrix into the exact same closed-form eigenvalue/eigenvector " +
						"formulas `eigenvectors-eigenvalues` used (λ = mid ± disc, θ = " +
						"atan2(2b, a−d)/2): with a=2, b=1.2, d=1.2, mid=(2+1.2)/2=1.6, " +
						"half=(2−1.2)/2=0.4, disc=√(0.4²+1.2²)=√1.6≈1.265, so λ1≈2.865, λ2≈0.335, " +
						"and θ1 = atan2(2.4, 0.8)/2 ≈ 35.8°.",
					"λ1's direction (35.8° off the x-axis — not 0° or 90°, i.e. not aligned to " +
						"either original axis) is the first principal component: the single " +
						"direction the cloud spreads out along the most. Its eigenvalue, 2.865, is " +
						"exactly the variance of the data once you re-measure it along that " +
						"direction instead of along x or y. Add both eigenvalues: 2.865+0.335=3.2 " +
						"— exactly Var(x)+Var(y)=2+1.2=3.2, confirming the eigen-decomposition just " +
						"redistributes the same total variance across two new, uncorrelated " +
						"directions rather than creating or destroying any of it. 2.865/3.2 ≈ 89.5% " +
						"of the cloud's total spread lives along that one 35.8° direction alone.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The r slider sets how correlated x and y are (same meaning as in " +
						"`covariance` and `correlation`); the stretch slider re-stretches the " +
						"x-axis, growing Var(x) the same way `covariance`'s scale slider did. The " +
						"green arrow is the first principal component (PC1) — the direction of " +
						"maximum spread — drawn one standard deviation long (√λ1); the shorter " +
						"perpendicular arrow is the second principal component (PC2), one √λ2 " +
						"long. Watch PC1 rotate as you change the sliders: increasing stretch alone " +
						"drags PC1 toward the x-axis, because x now carries more of the variance; " +
						"increasing r alone drags it toward the 45° diagonal. The readout reports " +
						"what fraction of the total variance PC1 alone explains — climbing toward " +
						"100% as the cloud gets thinner and more line-like, and dropping toward 50% " +
						"as the cloud rounds out toward a circle (r near 0, stretch near 1).",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Summarize a cloud of correlated points with one number and one direction " +
						"instead of two separate, axis-tied variances: instead of reporting " +
						"'Var(x)=2, Var(y)=1.2' — two numbers that depend on which way you happened " +
						"to draw the axes — report 'this data spreads out about 89.5% along one " +
						"specific direction, and barely at all perpendicular to it,' a fact about " +
						"the cloud itself, not about the axes. That's also literally how PCA " +
						"compresses data with many variables: instead of keeping all of them, keep " +
						"only the first few principal components, whichever few directions capture " +
						"the most variance, and accept a small, quantifiable loss (the leftover " +
						"eigenvalues) in exchange for far fewer numbers per data point.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A face-recognition system built on 'eigenfaces' represents each face image " +
						"(thousands of raw pixel values) as a handful of principal-component " +
						"coordinates instead, because most of the pixel-to-pixel variation across a " +
						"set of faces packs into a small number of directions. A genetics study " +
						"with measurements across thousands of genes uses just the first two or " +
						"three principal components to plot samples on a 2D scatterplot a human can " +
						"actually look at, trusting that those few directions carry most of the " +
						"real signal. Finance runs PCA on many correlated stock returns to find a " +
						"small number of independent 'risk factors' driving the whole market, " +
						"instead of tracking every stock's covariance with every other stock " +
						"individually.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: PC1 is the direction of maximum variance, and its " +
						"eigenvalue tells you exactly how much of the total variance that direction " +
						"accounts for — a number you can quote precisely, like '89.5% of variance " +
						"explained by PC1.'",
					"Not like this: assuming the first principal component always lines up with " +
						"one of your original axes, or that it always means 'the answer' in some " +
						"deeper sense — PC1 is only ever the direction of maximum spread in the " +
						"data you fed it; it has no idea which axis you'd find meaningful. Feed it " +
						"a perfectly round, uncorrelated cloud (r=0, stretch=1) and PC1 and PC2 tie " +
						"at 50% each — the 'first' direction becomes arbitrary, since any direction " +
						"is as good as any other once the cloud has no elongation left to point " +
						"along.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.7},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
			{Key: "stretch", Label: "x-axis stretch (×)", Min: 1, Max: 4, Step: 0.25, Def: 2},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value in
// [0, 1). It is a fixed function of its input, not a random-number generator,
// so GeneratePoints stays pure. Identical to correlation.hash01 and
// covariance.hash01 — each concept package is self-contained (see
// BUILD_CYCLE.md) — reused here to build the same kind of scatter cloud.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// GeneratePoints returns n (x, y) pairs whose population correlation is r,
// both standardized to mean 0, variance 1 before any Scale is applied.
// Identical in spirit to covariance.GeneratePoints.
func GeneratePoints(r float64, n int) (xs, ys []float64) {
	if n < 1 {
		n = 1
	}
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	xs = make([]float64, n)
	ys = make([]float64, n)
	tail := math.Sqrt(1 - r*r)
	for i := 0; i < n; i++ {
		u1 := hash01(2*i + 1)
		u2 := hash01(2*i + 2)
		if u1 < 1e-9 {
			u1 = 1e-9 // keep log() finite
		}
		mag := math.Sqrt(-2 * math.Log(u1))
		z1 := mag * math.Cos(2*math.Pi*u2)
		z2 := mag * math.Sin(2*math.Pi*u2)
		xs[i] = z1
		ys[i] = r*z1 + tail*z2
	}
	return xs, ys
}

// Scale multiplies every x value by factor — relabeling the x-axis's units,
// exactly like covariance.Scale.
func Scale(xs []float64, factor float64) []float64 {
	out := make([]float64, len(xs))
	for i, x := range xs {
		out[i] = x * factor
	}
	return out
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

// Covariance returns the population covariance of xs and ys: the average of
// (x - mean x) * (y - mean y). Returns 0 if the slices are empty or of
// mismatched length. A variable's covariance with itself is its variance.
func Covariance(xs, ys []float64) float64 {
	n := len(xs)
	if n == 0 || len(ys) != n {
		return 0
	}
	mx, my := Mean(xs), Mean(ys)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += (xs[i] - mx) * (ys[i] - my)
	}
	return sum / float64(n)
}

// CovarianceMatrix returns the three distinct entries of the symmetric 2x2
// covariance matrix [[varX, covXY], [covXY, varY]] for (xs, ys).
func CovarianceMatrix(xs, ys []float64) (varX, covXY, varY float64) {
	return Covariance(xs, xs), Covariance(xs, ys), Covariance(ys, ys)
}

// Eigenvalues returns the two eigenvalues of the symmetric 2x2 matrix
// [[a, b], [b, d]], larger first (lambda1 >= lambda2). Identical closed-form
// formula to eigen.Eigenvalues — every real symmetric matrix (a covariance
// matrix always is one) has real eigenvalues, so no complex numbers needed.
func Eigenvalues(a, b, d float64) (lambda1, lambda2 float64) {
	mid := (a + d) / 2
	half := (a - d) / 2
	disc := math.Sqrt(half*half + b*b)
	return mid + disc, mid - disc
}

// PrincipalAngle returns the direction (radians, in [0, π)) of the
// eigenvector belonging to lambda1 — the larger eigenvalue from
// Eigenvalues — i.e. the first principal component. Identical closed-form
// formula to eigen.EigenvectorAngle. The second principal component always
// sits exactly perpendicular to this one, at angle+π/2.
func PrincipalAngle(a, b, d float64) float64 {
	return math.Atan2(2*b, a-d) / 2
}

// ExplainedVarianceRatio returns lambda1/(lambda1+lambda2): the fraction of
// the cloud's total variance the first principal component alone accounts
// for. Returns 0 for a degenerate, zero-variance cloud (both eigenvalues 0)
// rather than dividing by zero.
func ExplainedVarianceRatio(lambda1, lambda2 float64) float64 {
	total := lambda1 + lambda2
	if total <= 0 {
		return 0
	}
	return lambda1 / total
}

func render(params map[string]float64) string {
	r, n, stretch := params["r"], int(params["n"]), params["stretch"]
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	if n < 2 {
		n = 2
	}
	if stretch < 1 {
		stretch = 1
	}

	rawXs, ys := GeneratePoints(r, n)
	xs := Scale(rawXs, stretch)

	varX, covXY, varY := CovarianceMatrix(xs, ys)
	l1, l2 := Eigenvalues(varX, covXY, varY)
	theta1 := PrincipalAngle(varX, covXY, varY)
	theta2 := theta1 + math.Pi/2
	ratio := ExplainedVarianceRatio(l1, l2)

	// y stays standardized (variance 1); x stretches with stretch, so the
	// canvas widens along x the same way covariance's picture does.
	const ylim = 3.6
	xlim := ylim * stretch
	c := viz.New(680, 460, -xlim, xlim, -ylim, ylim)
	c.Path([][2]float64{{-xlim, 0}, {xlim, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, -ylim}, {0, ylim}}, viz.Muted, 1)

	// The scatter cloud.
	for i := range xs {
		px, py := c.X(xs[i]), c.Y(ys[i])
		c.Rect(px-2.5, py-2.5, 5, 5, viz.Accent, 0.5)
	}

	// PC1 (green) drawn one standard deviation long, PC2 (orange) likewise
	// -- both as full lines through the origin so their directions read
	// clearly against the cloud, plus an arrowhead marking each end.
	l1sd, l2sd := math.Sqrt(math.Max(l1, 0)), math.Sqrt(math.Max(l2, 0))
	pc1x, pc1y := l1sd*math.Cos(theta1), l1sd*math.Sin(theta1)
	pc2x, pc2y := l2sd*math.Cos(theta2), l2sd*math.Sin(theta2)
	arrow(c, 0, 0, pc1x, pc1y, viz.Good, 2.5)
	arrow(c, 0, 0, -pc1x, -pc1y, viz.Good, 1.5)
	arrow(c, 0, 0, pc2x, pc2y, viz.Warm, 2.5)
	arrow(c, 0, 0, -pc2x, -pc2y, viz.Warm, 1.5)

	c.Text(16, 24, fmt.Sprintf("covariance matrix: [[%.2f, %.2f], [%.2f, %.2f]]", varX, covXY, covXY, varY),
		14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("PC1 (green): θ=%.1f°, λ1=%.2f", theta1*180/math.Pi, l1),
		14, viz.Good, "start")
	c.Text(16, 64, fmt.Sprintf("PC2 (orange): θ=%.1f°, λ2=%.2f", theta2*180/math.Pi, l2),
		14, viz.Warm, "start")
	c.Text(16, 84, fmt.Sprintf("PC1 explains %.1f%% of total variance", ratio*100),
		14, viz.Ink, "start")

	return c.String()
}

// arrow draws a straight line from (x0,y0) to (x1,y1) in data space, with a
// small V-shaped arrowhead at the end. Identical to eigen.arrow.
func arrow(c *viz.Canvas, x0, y0, x1, y1 float64, color string, width float64) {
	c.Path([][2]float64{{x0, y0}, {x1, y1}}, color, width)

	dx, dy := x1-x0, y1-y0
	length := math.Hypot(dx, dy)
	if length < 1e-9 {
		return
	}
	ux, uy := dx/length, dy/length
	const headLen = 0.2
	const headAngle = 0.5 // radians, ~29 degrees off the shaft on each side

	barb := func(t float64) (float64, float64) {
		cos, sin := math.Cos(t), math.Sin(t)
		bx, by := -ux, -uy
		return bx*cos - by*sin, bx*sin + by*cos
	}
	b1x, b1y := barb(headAngle)
	b2x, b2y := barb(-headAngle)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b1x, y1 + headLen*b1y}}, color, width)
	c.Path([][2]float64{{x1, y1}, {x1 + headLen*b2x, y1 + headLen*b2y}}, color, width)
}
