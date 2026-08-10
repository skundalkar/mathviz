// Package expgrowth visualizes compound (exponential) growth against the
// straight-line growth your intuition tends to reach for instead. The two
// start out looking almost identical — that's exactly the trap. Compounding
// keeps growing off of an ever-larger base, so it eventually pulls away and
// keeps accelerating, while linear growth just keeps adding the same fixed
// amount forever.
package expgrowth

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "exponential-growth",
		Seq:   19,
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
	rate := p["rate"]
	if rate < 0.01 {
		rate = 0.01
	}
	periods := p["periods"]
	if periods < 1 {
		periods = 1
	}

	finalExp := Value(periods, rate)
	finalLin := Linear(periods, rate)
	ymax := finalExp
	if finalLin > ymax {
		ymax = finalLin
	}
	ymax *= 1.08

	c := viz.New(680, 400, 0, periods, 0, ymax)
	c.Axes()
	for t := 0.0; t <= periods; t += periods / 5 {
		c.Tick(t, fmt.Sprintf("%.0f", t))
	}

	// The naive straight-line projection: same starting point, same initial
	// slope, but it never accelerates.
	linear := viz.Sample(0, periods, 200, func(t float64) float64 { return Linear(t, rate) })
	c.Path(linear, viz.Muted, 2)

	// The real, compounding curve.
	curve := viz.Sample(0, periods, 200, func(t float64) float64 { return Value(t, rate) })
	c.Path(curve, viz.Accent, 2.5)

	// Mark every doubling: t=0 (value 1), then each successive doubling time,
	// as long as it still falls within the visible window.
	dt := DoublingTime(rate)
	if !math.IsInf(dt, 1) {
		for k := 0; ; k++ {
			t := dt * float64(k)
			if t > periods {
				break
			}
			v := Value(t, rate)
			px, py := c.X(t), c.Y(v)
			c.Rect(px-3, py-3, 6, 6, viz.Warm, 1)
			c.VLine(t, viz.Warm, true)
		}
	}

	c.Text(20, 24, fmt.Sprintf("growth rate = %.0f%% per period    doubling time ≈ %.2f periods (rule of 70 ≈ %.1f)",
		rate, dt, Rule70(rate)), 13, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("after %.0f periods: compounding = %.2fx   straight-line guess = %.2fx",
		periods, finalExp, finalLin), 13, viz.Muted, "start")
	c.Text(20, 62, "solid = compounding value(t)=(1+rate)^t    thin = linear guess from the same starting slope",
		12, viz.Muted, "start")

	return c.String()
}
