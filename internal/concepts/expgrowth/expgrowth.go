// Package expgrowth visualizes compound (exponential) growth against the
// straight-line growth your intuition tends to reach for instead. The two
// start out looking almost identical — that's exactly the trap. Compounding
// keeps growing off of an ever-larger base, so it eventually pulls away and
// keeps accelerating, while linear growth just keeps adding the same fixed
// amount forever.
package expgrowth

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "exponential-growth",
		Title: "Exponential growth & doubling time",
		Blurb: "A lily pad patch doubles in size every day, and takes 30 days to cover an " +
			"entire pond. On what day is the pond half covered? Almost everyone's gut answer " +
			"is 'day 15' — the halfway mark. The real answer is day 29: doubling from half-full " +
			"to full only takes one more day, because each day's growth is a multiple of " +
			"whatever's already there, not a fixed amount. For the first 25 days the pond looks " +
			"basically empty (under 10%) — anyone watching would swear nothing much is " +
			"happening — then the last handful of days it explodes, because 'a little more than " +
			"before' becomes a much bigger absolute number once the base is already large. " +
			"That's compounding: value(t) = (1+rate)^t. Doubling time isn't a fixed number of " +
			"periods you add, it's driven by the rate itself: doubling time = ln(2)/ln(1+rate) " +
			"— higher rate, shorter doubling time, and the doublings just keep stacking. Compare " +
			"the curve to the dashed straight line (same starting slope) and watch how close " +
			"they track at first, then how fast the gap opens up.",
		Params: []concept.ParamSpec{
			{Key: "rate", Label: "Growth rate per period (%)", Min: 1, Max: 50, Step: 1, Def: 10},
			{Key: "periods", Label: "Periods shown", Min: 5, Max: 40, Step: 1, Def: 20},
		},
		Render: render,
	})
}

// Value returns the compounded value at period t, starting from 1, growing
// ratePct percent every period: (1 + ratePct/100)^t. Pure math — t need not
// be an integer, so the curve can be sampled anywhere along it.
func Value(t, ratePct float64) float64 {
	return math.Pow(1+ratePct/100, t)
}

// Linear returns the straight-line projection that starts at the same point
// (1) and matches the curve's initial slope: 1 + rate*t. It's the naive
// "extend the current trend" guess intuition tends to reach for, and it
// only agrees with Value near t=0.
func Linear(t, ratePct float64) float64 {
	return 1 + (ratePct/100)*t
}

// DoublingTime returns how many periods it takes the compounding value to
// double: ln(2)/ln(1+rate). Returns +Inf for a non-positive rate, since a
// value that never grows never doubles.
func DoublingTime(ratePct float64) float64 {
	rate := ratePct / 100
	if rate <= 0 {
		return math.Inf(1)
	}
	return math.Log(2) / math.Log(1+rate)
}

// Rule70 is the classic mental-math shortcut for doubling time: 70 divided
// by the growth rate as a whole number (e.g. 7% growth ≈ 70/7 = 10 periods
// to double). It's an approximation of DoublingTime, accurate to within a
// few percent for rates under about 15%, and drifts further off at higher
// rates. Returns +Inf for a non-positive rate.
func Rule70(ratePct float64) float64 {
	if ratePct <= 0 {
		return math.Inf(1)
	}
	return 70 / ratePct
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 380, 0, 1, 0, 1).String()
}
