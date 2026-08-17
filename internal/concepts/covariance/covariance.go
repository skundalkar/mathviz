// Package covariance visualizes the raw building block underneath
// correlation: how two variables move together in whatever units they
// happen to be measured in. Unlike correlation's r, covariance isn't
// dimensionless — its number changes if you simply relabel an axis from
// centimeters to millimeters, even though the actual relationship between
// the two variables hasn't changed at all.
package covariance

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "covariance",
		Seq:   44,
		Title: "Covariance (and why we normalize it)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`correlation` gives you r, a single number from -1 to +1 that says how " +
						"tightly two variables move together, no matter what they're measured in. " +
						"But r isn't computed directly from the raw data — underneath it is a " +
						"simpler-looking quantity, the average of (x - mean x) times " +
						"(y - mean y), that seems like it should already answer 'do these move " +
						"together?' on its own. Take five people's height in centimeters and " +
						"weight in kilograms, compute that number, then remeasure the exact same " +
						"five people's height in millimeters instead. Same people, same weights, " +
						"same real relationship — does that number change, and if so, is it still " +
						"telling you anything about how strongly height and weight move together?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Five students, hours studied (x) vs. test score out of 100 (y): x = 1,2,3,4,5, " +
						"y = 2,4,5,4,5 (scaled to a 0-100 test, read as tens of points). Mean x = 3, " +
						"mean y = 4. For each student, multiply (x - 3) by (y - 4):",
					"• Student 1: (1-3)(2-4) = (-2)(-2) = 4. Student 2: (2-3)(4-4) = (-1)(0) = 0.",
					"• Student 3: (3-3)(5-4) = (0)(1) = 0. Student 4: (4-3)(4-4) = (1)(0) = 0.",
					"• Student 5: (5-3)(5-4) = (2)(1) = 2.",
					"Average those five products: (4+0+0+0+2)/5 = 1.2 — that's the covariance. " +
						"Now relabel x as 'tens of minutes studied' instead of 'hours' (x times 10): " +
						"1,2,3,4,5 becomes 10,20,30,40,50, mean x becomes 30, and every (x-mean) term " +
						"is now 10x bigger, so the covariance becomes 10 x 1.2 = 12 — ten times " +
						"bigger, purely from relabeling a unit, with the actual pattern in the data " +
						"completely unchanged. Divide covariance by both standard deviations " +
						"(std x = 1.414, std y = 1.095 in the original units) and you get " +
						"1.2 / (1.414 x 1.095) = 0.775 — correlation r. Redo that division after " +
						"the x10 relabeling: std x also grows by exactly 10x (to 14.14), so the 10s " +
						"cancel and r comes back out to 0.775, unchanged.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The r slider sets the target correlation of a scatter cloud, same as " +
						"`correlation`'s picture. The scale slider relabels the x-axis units (x1 " +
						"through x10, like switching from centimeters to millimeters) — the cloud's " +
						"shape and its correlation r never change, but the axis numbers stretch out " +
						"and the raw covariance readout grows right along with them, in exact " +
						"proportion to scale. Watch the two numbers at the top: correlation r stays " +
						"pinned to the slider's target the whole time; covariance climbs steadily as " +
						"scale increases, even though nothing about how tightly the points hug the " +
						"trend line has changed.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Recognize why correlation, not covariance, is the number to report or compare " +
						"across different studies or different units: a covariance of 12 between " +
						"height (in mm) and weight (in kg) can't be compared to a covariance of 1.2 " +
						"between height (in cm) and weight (in kg) — same relationship, different " +
						"number, purely from unit choice. Correlation divides that unit-dependence " +
						"back out, which is exactly why `correlation`'s r stays meaningful across " +
						"any two datasets no matter what scale each was measured on. Covariance " +
						"still has its own use once you already know the units matter: its raw " +
						"sign tells you the direction of the relationship (positive vs. negative) " +
						"just as reliably as r's sign does, before you've bothered to normalize it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A finance textbook's covariance matrix between two stocks' returns is in " +
						"units of 'percent squared,' which is why practitioners almost always " +
						"convert it to a correlation before comparing two different stock pairs. A " +
						"scientific paper reporting 'covariance between temperature (Celsius) and " +
						"ice cream sales (dollars)' is reporting a number tied to those exact units " +
						"— switch to Fahrenheit or another currency and the number changes even " +
						"though the underlying relationship is identical. `eigenvectors-eigenvalues` " +
						"applied to a covariance matrix is also the core machinery behind Principal " +
						"Component Analysis, a widely used technique for compressing many correlated " +
						"variables (like pixels in an image or genes in an experiment) down to a few " +
						"informative directions.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: covariance's sign (positive, negative, or near zero) tells " +
						"you the direction of a linear relationship, exactly like correlation's " +
						"sign does — but its magnitude depends on the units of both variables, so " +
						"a covariance of 12 isn't automatically 'more' relationship than a " +
						"covariance of 1.2 unless you already know they're measured the same way.",
					"Not like this: comparing two raw covariance numbers across different datasets " +
						"or units to decide which relationship is 'stronger' — that comparison is " +
						"only valid after normalizing to correlation, same as it would be a mistake " +
						"to compare a `variance-vs-stddev` in squared units directly against a " +
						"standard deviation in linear units without converting one to match the " +
						"other first.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.7},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
			{Key: "scale", Label: "x-axis units (×)", Min: 1, Max: 10, Step: 1, Def: 1},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value in
// [0, 1). It is a fixed function of its input — not a random-number generator
// — so GeneratePoints stays pure: same (r, n) always produces the same cloud.
// Identical to correlation.hash01 (each concept package is self-contained;
// see BUILD_CYCLE.md), reused here to build the same kind of scatter cloud.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// GeneratePoints returns n (x, y) pairs whose population correlation is r.
// Both x and y are standardized (mean 0, variance 1) before any Scale is
// applied — see Scale below for the unit-relabeling step that keeps
// correlation fixed while covariance moves.
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

// Scale multiplies every x value by factor — a stand-in for relabeling the
// x-axis's units (centimeters to millimeters, hours to minutes). It changes
// nothing about the actual relationship between x and y, only the number
// each x is reported as.
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
// mismatched length.
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

// StdDev returns the population standard deviation of xs: sqrt(Covariance(xs, xs)),
// since a variable's covariance with itself is exactly its variance.
func StdDev(xs []float64) float64 {
	return math.Sqrt(Covariance(xs, xs))
}

// Correlation returns the Pearson correlation coefficient of xs and ys —
// covariance normalized by both standard deviations, which is what makes it
// invariant to rescaling either variable (see LEARNINGS.md). Returns 0 if
// either series has zero variance (undefined correlation, rather than
// dividing by zero).
func Correlation(xs, ys []float64) float64 {
	sx, sy := StdDev(xs), StdDev(ys)
	if sx <= 0 || sy <= 0 {
		return 0
	}
	return Covariance(xs, ys) / (sx * sy)
}

func render(params map[string]float64) string {
	r, n, scale := params["r"], int(params["n"]), params["scale"]
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	if n < 2 {
		n = 2
	}
	if scale < 1 {
		scale = 1
	}

	rawXs, ys := GeneratePoints(r, n)
	xs := Scale(rawXs, scale)

	cov := Covariance(xs, ys)
	corr := Correlation(xs, ys)

	// Both axes are standardized before scaling (variance 1), so the y
	// range stays fixed at a comfortable +-3.6 while the x range stretches
	// with scale — the same cloud, relabeled onto wider and wider units.
	const ylim = 3.6
	xlim := ylim * scale
	c := viz.New(680, 400, -xlim, xlim, -ylim, ylim)
	c.Axes()
	xTick := xlim / 3
	for x := -xlim; x <= xlim+1e-9; x += xTick {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// Reference lines through the origin.
	c.Path([][2]float64{{-xlim, 0}, {xlim, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, -ylim}, {0, ylim}}, viz.Muted, 1)

	// The scatter cloud itself: one small square per point.
	for i := range xs {
		px, py := c.X(xs[i]), c.Y(ys[i])
		c.Rect(px-2.5, py-2.5, 5, 5, viz.Accent, 0.55)
	}

	// The trend line r implies: slope r/scale in these axis units, since y
	// is standardized but x has been stretched by scale.
	slope := r / scale
	c.Path([][2]float64{{-xlim, -xlim * slope}, {xlim, xlim * slope}}, viz.Warm, 2.5)

	c.Text(20, 24, fmt.Sprintf("x-axis units: ×%.0f    covariance = %.3f    correlation r = %.3f", scale, cov, corr),
		14, viz.Ink, "start")
	c.Text(20, 44, "covariance grows with the x-axis units — correlation r does not",
		12, viz.Muted, "start")
	c.Text(20, 62, fmt.Sprintf("target r = %.2f    n = %d", r, n),
		12, viz.Muted, "start")
	return c.String()
}
