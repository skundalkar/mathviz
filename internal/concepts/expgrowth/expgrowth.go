// Package expgrowth visualizes compound (exponential) growth against the
// straight-line growth your intuition tends to reach for instead. The two
// start out looking almost identical — that's exactly the trap. Compounding
// keeps growing off of an ever-larger base, so it eventually pulls away and
// keeps accelerating, while linear growth just keeps adding the same fixed
// amount forever.
package expgrowth

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 380, 0, 1, 0, 1).String()
}
