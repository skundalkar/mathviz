// Package bayestheorem visualizes why a positive test for a rare condition is
// so often wrong: a population of 100 people is split into those who tested
// positive and those who tested negative, colored by whether the test got it
// right. When the base rate is low, most of the "positive" group turns out to
// be false alarms even for a fairly accurate test — Bayes' theorem is just
// the arithmetic that explains why.
package bayestheorem

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bayes-theorem",
		Seq:   11,
		Title: "Bayes' theorem",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A test for a rare disease is 99% accurate, and you just tested positive. " +
						"What's the chance you actually have it?",
					"Most people's gut answer is '99%' — the test is 99% accurate, after all, what " +
						"else would it mean? That instinct is wrong, often dramatically so, and the " +
						"gap between the gut answer and the real one is exactly what Bayes' theorem " +
						"measures.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Imagine 100 people take the test, and the disease is rare — say only 1 person " +
						"actually has it:",
					"• The test almost certainly catches that 1 true case (99% sensitivity).",
					"• But it also wrongly flags about 5% of the 99 healthy people — roughly 5 " +
						"false alarms.",
					"• Count everyone who tested positive: 1 real case plus roughly 5 false alarms " +
						"— 6 positive results total, only 1 of which is real: about 17%, not 99%.",
					"The formula gives the same answer directly, computed independently of the " +
						"counting story: P(condition | positive) = (sensitivity x prior) / " +
						"(sensitivity x prior + false-positive-rate x (1 - prior)) = (0.99x0.01) / " +
						"(0.99x0.01 + 0.05x0.99) = 0.0099/0.0594 = 16.7% — two different routes, " +
						"same number.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Below, a population of 100 is regrouped by test result instead of true status, " +
						"at the prior/sensitivity/specificity you set. In the 'tested positive' strip, " +
						"green squares are true positives and red squares are false alarms — at the " +
						"default numbers (1% base rate, 99% sensitivity, 95% specificity) the strip " +
						"is mostly red: about 5 false alarms for every true positive, matching the " +
						"1-real-out-of-6 count above. Raise the base rate (prior) slider and green " +
						"takes over the strip; raise specificity and the false alarms vanish, because " +
						"specificity directly controls how many healthy people get swept into a " +
						"positive result in the first place.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now correct the gut-instinct '99% accurate means 99% likely' reasoning " +
						"by actually accounting for how rare the condition is before trusting a " +
						"positive result — updating your belief from a positive test now requires the " +
						"base rate, not just the test's accuracy.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Screening for a rare disease, a rare fraud pattern, or a rare security alert: a " +
						"'99% accurate' flag sounds like a near-certainty, but if the thing being " +
						"flagged is rare, most flags are still noise. It's why doctors order a " +
						"second, more specific test before acting on one positive screen, and why " +
						"'the model flagged it' needs a base rate attached before anyone should trust " +
						"the flag.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'don't ignore the base rate' — before updating your belief " +
						"off one new piece of evidence, remember how rare or common the thing was to " +
						"begin with; a positive test for a rare condition means far less than gut " +
						"instinct suggests.",
					"Not like this: 'the test is 99% accurate, so a positive result means a 99% " +
						"chance I have it' — that skips the base rate entirely, and for a rare " +
						"enough condition, the base rate is often most of what determines the real " +
						"answer.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "prior", Label: "Base rate (prior)", Min: 0.001, Max: 0.5, Step: 0.001, Def: 0.01},
			{Key: "sensitivity", Label: "Sensitivity (catches true cases)", Min: 0.5, Max: 0.999, Step: 0.001, Def: 0.99},
			{Key: "specificity", Label: "Specificity (clears true negatives)", Min: 0.5, Max: 0.999, Step: 0.001, Def: 0.95},
		},
		Render: render,
	})
}

// PosteriorPositive is P(condition | positive test) via Bayes' theorem:
//
//	P(C|+) = P(+|C)*P(C) / [ P(+|C)*P(C) + P(+|~C)*P(~C) ]
//
// where P(+|C) is the sensitivity and P(+|~C) is the false-positive rate
// (1 - specificity).
func PosteriorPositive(prior, sensitivity, specificity float64) float64 {
	fpr := 1 - specificity
	num := sensitivity * prior
	den := num + fpr*(1-prior)
	if den <= 0 {
		return 0
	}
	return num / den
}

// PosteriorNegative is P(no condition | negative test): given a negative
// result, how likely is the person actually healthy?
func PosteriorNegative(prior, sensitivity, specificity float64) float64 {
	fnr := 1 - sensitivity
	num := specificity * (1 - prior)
	den := num + fnr*prior
	if den <= 0 {
		return 0
	}
	return num / den
}

// Counts splits a population of n people into the four Bayes outcomes —
// true positive, false negative, false positive, true negative — given a
// prior prevalence, sensitivity and specificity. Counts are whole people,
// rounded to the nearest integer; any rounding slack is absorbed into tn
// (typically the largest group) so the four counts always sum to n.
func Counts(prior, sensitivity, specificity float64, n int) (tp, fn, fp, tn int) {
	if n < 0 {
		n = 0
	}
	diseased := prior * float64(n)
	healthy := float64(n) - diseased

	tp = roundHalfUp(diseased * sensitivity)
	fn = roundHalfUp(diseased * (1 - sensitivity))
	fp = roundHalfUp(healthy * (1 - specificity))
	tn = n - tp - fn - fp
	if tn < 0 {
		tn = 0
	}
	return
}

func roundHalfUp(x float64) int {
	if x < 0 {
		return 0
	}
	return int(math.Floor(x + 0.5))
}

func render(p map[string]float64) string {
	prior := clampProb(p["prior"])
	sens := clampProb(p["sensitivity"])
	spec := clampProb(p["specificity"])

	const n = 100
	tp, fn, fp, tn := Counts(prior, sens, spec, n)
	postPos := PosteriorPositive(prior, sens, spec)
	postNegHealthy := PosteriorNegative(prior, sens, spec)

	c := viz.New(680, 360, 0, 1, 0, 1)

	posColors := make([]string, 0, tp+fp)
	for i := 0; i < tp; i++ {
		posColors = append(posColors, viz.Good)
	}
	for i := 0; i < fp; i++ {
		posColors = append(posColors, viz.Bad)
	}

	negColors := make([]string, 0, fn+tn)
	for i := 0; i < fn; i++ {
		negColors = append(negColors, viz.Warm)
	}
	for i := 0; i < tn; i++ {
		negColors = append(negColors, viz.Faint)
	}

	c.Text(20, 24, fmt.Sprintf("population of %d   prior = %.1f%%   sensitivity = %.1f%%   specificity = %.1f%%",
		n, prior*100, sens*100, spec*100), 13, viz.Ink, "start")

	y := drawStrip(c, 46, fmt.Sprintf("tested positive (%d)  —  green = true positive, red = false alarm", len(posColors)), posColors)
	y = drawStrip(c, y+28, fmt.Sprintf("tested negative (%d)  —  orange = missed case, gray = true negative", len(negColors)), negColors)

	c.Text(20, y+26, fmt.Sprintf("P(condition | positive test) = %d/%d ≈ %.0f%%", tp, tp+fp, postPos*100), 15, viz.Bad, "start")
	c.Text(20, y+46, fmt.Sprintf("P(healthy | negative test) = %d/%d ≈ %.1f%%", tn, fn+tn, postNegHealthy*100), 13, viz.Muted, "start")

	return c.String()
}

// drawStrip renders a labeled row of small squares, wrapping to a new pixel
// row every maxCols squares, and returns the pixel y just below the strip so
// the caller can stack the next one underneath.
func drawStrip(c *viz.Canvas, top float64, label string, colors []string) (bottom float64) {
	c.Text(20, top, label, 12, viz.Muted, "start")

	const cell, gap, maxCols = 16.0, 2.0, 25
	startY := top + 10
	for i, col := range colors {
		row := i / maxCols
		colIdx := i % maxCols
		x := 20 + float64(colIdx)*(cell+gap)
		y := startY + float64(row)*(cell+gap)
		c.Rect(x, y, cell, cell, col, 0.9)
	}

	rows := (len(colors) + maxCols - 1) / maxCols
	if rows == 0 {
		rows = 1
	}
	return startY + float64(rows)*(cell+gap)
}

func clampProb(x float64) float64 {
	if x < 0.001 {
		return 0.001
	}
	if x > 0.999 {
		return 0.999
	}
	return x
}
