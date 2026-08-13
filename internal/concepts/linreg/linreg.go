// Package linreg visualizes simple linear regression: fitting the one
// straight line through a scatter of points that minimizes the total
// squared vertical distance between the line and every point (the
// "least-squares" line), and lets you try to beat it by hand.
package linreg

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "linear-regression",
		Seq:   30,
		Title: "Linear regression (fitting the least-squares line)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You track 5 students' study hours against their quiz scores: (1, 50), (2, " +
						"55), (3, 65), (4, 70), (5, 80). A new student says they plan to study 3.5 " +
						"hours — what score should they expect? Gut instinct: eyeball a line " +
						"through 'the middle of the dots' by hand. That's fine for a rough guess, " +
						"but ask two different people to eyeball it and you'll get two different " +
						"lines — different slopes, different predictions for 3.5 hours, and no way " +
						"to say which eyeballed line is actually 'best.' Is there a precise, " +
						"repeatable way to find the single line that best fits a scatter of points " +
						"— not just 'looks about right,' but provably the best possible fit by some " +
						"clear standard?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Define a residual as actual y minus predicted y — the vertical gap between " +
						"a real point and the line, at that point's own x. 'Best fit' means: choose " +
						"a slope and intercept that make the total squared residual as small as " +
						"possible (squaring so positive and negative gaps don't cancel out, and so " +
						"a couple of huge misses get penalized far more than several tiny ones) — " +
						"this is the 'least-squares' line. Minimizing that total turns out to have " +
						"a direct formula, no trial and error needed: slope = Σ(x−x̄)(y−ȳ) / " +
						"Σ(x−x̄)², intercept = ȳ − slope·x̄, where x̄ and ȳ are the average x and " +
						"average y.",
					"Walk the 5-student example. x̄ = 3, ȳ = 64.",
					"• deviations from the mean: x−x̄ = −2,−1,0,1,2 and y−ȳ = −14,−9,1,6,16",
					"• Σ(x−x̄)(y−ȳ) = (−2×−14)+(−1×−9)+(0×1)+(1×6)+(2×16) = 28+9+0+6+32 = 75",
					"• Σ(x−x̄)² = 4+1+0+1+4 = 10",
					"• slope = 75/10 = 7.5, intercept = 64 − 7.5×3 = 41.5",
					"So the least-squares line is predicted score = 41.5 + 7.5 × hours studied. " +
						"Check it against the real data: at x=1 it predicts 49 (actual 50, residual " +
						"+1); at x=2 it predicts 56.5 (actual 55, residual −1.5); at x=3 it predicts " +
						"64 (actual 65, residual +1); at x=4 it predicts 71.5 (actual 70, residual " +
						"−1.5); at x=5 it predicts 79 (actual 80, residual +1). Those five residuals " +
						"— +1, −1.5, +1, −1.5, +1 — sum to exactly 0. That's not a coincidence: the " +
						"least-squares line is always forced through the point (x̄, ȳ), the data's " +
						"own center of mass, which pins the positive and negative residuals to " +
						"balance out.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Five fixed points — one per student's (hours studied, quiz score) pair from " +
						"the worked example — plotted as a scatter. The green line is the " +
						"least-squares fit computed directly from the formula in section 2; its " +
						"'best possible' total squared error is shown as SSE(best). The blue line " +
						"is your own guess, built from the slope and intercept sliders, with thin " +
						"lines from each point down to your line showing its own residuals — " +
						"squared and summed into SSE(yours). Drag either slider and SSE(yours) " +
						"changes immediately; no matter how you drag, SSE(yours) never drops below " +
						"SSE(best) — the green line truly is the smallest possible total squared " +
						"error a straight line can achieve on this data.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Predict y for a new x you haven't observed, from a formula instead of a " +
						"guess — the opening question's 3.5-hour student is predicted to score " +
						"41.5 + 7.5×3.5 = 67.75. More generally, any two-variable scatter has " +
						"exactly one best-fit line by the least-squares standard, computable " +
						"directly with no trial and error — and this is the simplest possible case " +
						"of a much bigger idea: fitting parameters to minimize a total error is " +
						"exactly what gradient-descent (elsewhere in this gallery) does for far more " +
						"complex models where no direct formula like this one exists.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Economists predict sales from advertising spend, real-estate sites estimate a " +
						"home's price per square foot, and sports analysts predict performance from " +
						"training load — all starting from a 'line of best fit' through past data. " +
						"Calling something a 'linear relationship' or saying two things are " +
						"'trending together' in everyday conversation is often an informal nod to " +
						"exactly this: a fitted line with a clear, consistent slope.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the least-squares line minimizes the total squared " +
						"vertical distance between the line and the data points' — a specific, " +
						"provable standard, not just 'a line that looks close.' Not like this: " +
						"assuming a well-fit line proves causation (fitting a line to ice-cream " +
						"sales against shark attacks doesn't mean one causes the other — see the " +
						"correlation concept elsewhere in this gallery), or trusting a prediction " +
						"far outside the range of the observed x values just as much as one inside " +
						"it — predicting a score for 50 hours studied by extending this line isn't " +
						"warranted, since nothing in the data says the relationship stays straight " +
						"that far out.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "guessSlope", Label: "Your guess — slope", Min: -2, Max: 15, Step: 0.5, Def: 5},
			{Key: "guessIntercept", Label: "Your guess — intercept", Min: 20, Max: 70, Step: 0.5, Def: 50},
		},
		Render: render,
	})
}

// HoursStudied and QuizScore are the fixed 5-student worked example every
// Section walks through: hours studied (x) against quiz score (y).
var (
	HoursStudied = []float64{1, 2, 3, 4, 5}
	QuizScore    = []float64{50, 55, 65, 70, 80}
)

// Mean returns the arithmetic mean of v.
func Mean(v []float64) float64 {
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

// Slope returns the least-squares slope of the line fit to (xs, ys):
// Σ(x−x̄)(y−ȳ) / Σ(x−x̄)². xs and ys must be the same non-empty length.
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

// Intercept returns the least-squares intercept of the line fit to (xs,
// ys): ȳ − slope·x̄ — the value that forces the line through the data's
// own center of mass (x̄, ȳ).
func Intercept(xs, ys []float64) float64 {
	return Mean(ys) - Slope(xs, ys)*Mean(xs)
}

// Predict evaluates a line (given its slope and intercept) at x.
func Predict(slope, intercept, x float64) float64 {
	return slope*x + intercept
}

// SumSquaredError totals the squared residual (actual y minus predicted y)
// of a line (slope, intercept) across every (xs[i], ys[i]) pair — the
// quantity the least-squares line is built to minimize.
func SumSquaredError(xs, ys []float64, slope, intercept float64) float64 {
	var sse float64
	for i := range xs {
		r := ys[i] - Predict(slope, intercept, xs[i])
		sse += r * r
	}
	return sse
}

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 1, 0, 1)
	return c.String()
}
