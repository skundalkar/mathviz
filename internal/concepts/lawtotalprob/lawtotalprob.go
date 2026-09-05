// Package lawtotalprob visualizes the law of total probability: splitting a
// hard-to-compute overall probability into a probability-weighted
// combination of easier per-scenario pieces. The running example is a
// factory floor with three machines, each producing a different share of
// total output at a different defect rate -- what fraction of ALL widgets
// coming off the floor are defective?
package lawtotalprob

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "law-of-total-probability",
		Seq:   91,
		Title: "Law of total probability (combining per-scenario pieces)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "shareA", Label: "Machine A's share of output", Min: 0.5, Max: 0.96, Step: 0.01, Def: 0.90},
			{Key: "shareB", Label: "Machine B's share of output", Min: 0.02, Max: 0.40, Step: 0.01, Def: 0.08},
			{Key: "rateB", Label: "Machine B's defect rate", Min: 0.02, Max: 0.50, Step: 0.01, Def: 0.20},
		},
		Render: render,
	})
}

// Shares normalizes the two free machine shares (A and B) into a valid
// three-way partition (a, b, c) that sums to exactly 1, with the remainder
// going to machine C. If shareA+shareB would leave machine C with less
// than 1% of output, both are scaled down proportionally so c stays at
// least 0.01 -- every scenario in the partition keeps some nonzero weight.
func Shares(shareA, shareB float64) (a, b, c float64) {
	if shareA < 0 {
		shareA = 0
	}
	if shareB < 0 {
		shareB = 0
	}
	if shareA+shareB > 0.99 {
		scale := 0.99 / (shareA + shareB)
		shareA *= scale
		shareB *= scale
	}
	return shareA, shareB, 1 - shareA - shareB
}

// TotalProbability implements the law of total probability itself:
// P(event) = Sum_i P(event|scenario_i) * P(scenario_i). Given each
// scenario's share of the population (shares, which should sum to 1) and
// the event's probability within that scenario (rates), it returns the
// event's overall probability across the whole population. Mismatched
// slice lengths or empty input return 0.
func TotalProbability(shares, rates []float64) float64 {
	sum := 0.0
	n := len(shares)
	if len(rates) < n {
		n = len(rates)
	}
	for i := 0; i < n; i++ {
		sum += shares[i] * rates[i]
	}
	return sum
}

// NaiveAverage returns the plain (unweighted) arithmetic mean of the
// per-scenario rates, ignoring how common each scenario actually is. It
// exists to demonstrate the common mistake of treating every scenario as
// equally likely instead of weighting by its real share -- see the
// "What's the common mistake here?" section.
func NaiveAverage(rates []float64) float64 {
	if len(rates) == 0 {
		return 0
	}
	sum := 0.0
	for _, r := range rates {
		sum += r
	}
	return sum / float64(len(rates))
}

// rateA and rateC are held fixed: both are reliable machines with the same
// low defect rate, so the picture isolates the effect of the mix (shareA,
// shareB) and the problem machine's own rate (rateB) without adding a
// fourth and fifth slider.
const rateA, rateC = 0.01, 0.01

func render(params map[string]float64) string {
	shareA, shareB, shareC := Shares(params["shareA"], params["shareB"])
	rateB := params["rateB"]
	if rateB < 0 {
		rateB = 0
	}

	shares := []float64{shareA, shareB, shareC}
	rates := []float64{rateA, rateB, rateC}
	weighted := TotalProbability(shares, rates)
	naive := NaiveAverage(rates)

	c := viz.New(680, 420, 0, 1, 0, 1)

	c.Text(16, 24, "A factory floor: three machines, each a different share of output, each its own defect rate",
		13, viz.Muted, "start")

	// Stacked bar: one horizontal segment per machine, width proportional
	// to its share of total output, filled more deeply red the higher its
	// own defect rate.
	const barX0, barY0, barW, barH = 20.0, 50.0, 640.0, 60.0
	names := []string{"A", "B", "C"}
	maxRate := rateA
	for _, r := range rates {
		if r > maxRate {
			maxRate = r
		}
	}
	x := barX0
	for i, share := range shares {
		w := share * barW
		opacity := 0.15 + 0.85*(rates[i]/maxRate)
		c.Rect(x, barY0, w, barH, viz.Bad, opacity)
		if w > 30 {
			c.Text(x+w/2, barY0+barH/2+5, names[i], 14, viz.Ink, "middle")
		}
		c.Text(x+w/2, barY0+barH+18, fmt.Sprintf("share=%.0f%%", share*100), 12, viz.Muted, "middle")
		c.Text(x+w/2, barY0+barH+34, fmt.Sprintf("rate=%.0f%%", rates[i]*100), 12, viz.Muted, "middle")
		x += w
	}

	c.Text(16, 160, fmt.Sprintf("P(defect) = %.2f×%.2f + %.2f×%.2f + %.2f×%.2f = %.2f%%",
		shareA, rateA, shareB, rateB, shareC, rateC, weighted*100), 14, viz.Ink, "start")

	// A small number line contrasting the correctly weighted total
	// probability against the naive (wrong) unweighted average of the
	// three rates -- the "What's the common mistake here?" section made
	// visible.
	const axisX0, axisX1, axisY = 40.0, 640.0, 260.0
	maxPct := naive * 100 * 1.3
	if weighted*100*1.3 > maxPct {
		maxPct = weighted * 100 * 1.3
	}
	if maxPct < 10 {
		maxPct = 10
	}
	toX := func(pct float64) float64 {
		return axisX0 + (pct/maxPct)*(axisX1-axisX0)
	}
	c.Rect(axisX0, axisY, axisX1-axisX0, 2, viz.Muted, 1)
	for pct := 0.0; pct <= maxPct; pct += maxPct / 5 {
		c.Text(toX(pct), axisY+20, fmt.Sprintf("%.0f%%", pct), 11, viz.Muted, "middle")
	}

	weightedX := toX(weighted * 100)
	c.Rect(weightedX-2, axisY-16, 4, 16, viz.Good, 1)
	c.Text(weightedX, axisY-22, fmt.Sprintf("correct: %.2f%%", weighted*100), 12, viz.Good, "middle")

	naiveX := toX(naive * 100)
	c.Rect(naiveX-2, axisY-16, 4, 16, viz.Warm, 0.6)
	c.Text(naiveX, axisY+50, fmt.Sprintf("naive unweighted avg: %.2f%% (wrong)", naive*100), 12, viz.Warm, "middle")

	return c.String()
}
