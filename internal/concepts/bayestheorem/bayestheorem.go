// Package bayestheorem visualizes why a positive test for a rare condition is
// so often wrong: a population of 100 people is split into those who tested
// positive and those who tested negative, colored by whether the test got it
// right. When the base rate is low, most of the "positive" group turns out to
// be false alarms even for a fairly accurate test — Bayes' theorem is just
// the arithmetic that explains why.
package bayestheorem

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bayes-theorem",
		Title: "Bayes' theorem",
		Blurb: "A population of 100 people is split into a 'tested positive' group and a " +
			"'tested negative' group. Within 'tested positive', green squares are people who " +
			"really are sick (true positives) and red squares are people who aren't (false " +
			"alarms). When the condition is rare, red usually outnumbers green even for a " +
			"pretty accurate test — that's Bayes' theorem: a positive result updates your " +
			"belief, but it starts from how rare the condition was in the first place.",
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
	_ = p
	return viz.New(680, 360, 0, 1, 0, 1).String()
}
