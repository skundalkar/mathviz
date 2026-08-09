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
		Title: "Bayes' theorem",
		Blurb: "A test for a rare disease is 99% accurate, and you just tested positive. " +
			"What's the chance you actually have it? Most people's gut answer is '99%.' It's " +
			"wrong: with 100 people and only 1 actually sick, the test almost certainly catches " +
			"that 1 (99% sensitivity) but also wrongly flags about 5% of the 99 healthy people " +
			"— roughly 5 false alarms. That's 6 positive results total, only 1 real: about 17%, " +
			"not 99%. The formula gives the same answer directly: (0.99x0.01) / (0.99x0.01 + " +
			"0.05x0.99) = 0.0099/0.0594 = 16.7% — two routes, same number. Below, a population " +
			"of 100 is regrouped by test result: green squares in 'tested positive' are real " +
			"cases, red squares are false alarms — and red usually outnumbers green when the " +
			"condition is rare, even for a pretty accurate test.",
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
