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
				Body: []string{
					"bayes-theorem's denominator quietly did something worth pulling out on its " +
						"own: to get the overall probability of testing positive, it added up " +
						"P(positive|condition)×P(condition) and P(positive|no condition)×P(no " +
						"condition) — two scenarios, combined. But real situations are rarely split " +
						"into just two scenarios, and the combining step is easy to get wrong once " +
						"there are several: a factory floor runs three different machines, each " +
						"producing a different share of total output, and each with its own defect " +
						"rate. The gut-instinct shortcut is to just average the three defect rates " +
						"together — but that treats a machine that makes 2% of your output the same " +
						"as one that makes 90% of it. If one small, badly-tuned machine has a " +
						"terrible defect rate, does that make your *overall* defect rate terrible " +
						"too, or does high-volume production from the good machines dilute it away? " +
						"How do you correctly combine several scenario-specific probabilities, each " +
						"weighted by how common that scenario actually is, into one honest overall " +
						"number?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Three machines share a factory floor: Machine A makes 90% of all output with " +
						"a 1% defect rate. Machine B is a smaller, less reliable machine making 8% " +
						"of output with a 20% defect rate. Machine C makes the remaining 2% with a " +
						"1% defect rate, same as A. What fraction of ALL widgets coming off the " +
						"floor are defective?",
					"• Machine A's contribution: 90% of output × 1% defect rate = 0.90×0.01 = 0.90%.",
					"• Machine B's contribution: 8% of output × 20% defect rate = 0.08×0.20 = 1.60%.",
					"• Machine C's contribution: 2% of output × 1% defect rate = 0.02×0.01 = 0.02%.",
					"• Add the three contributions: 0.90% + 1.60% + 0.02% = 2.52% — the overall " +
						"defect rate across the whole floor.",
					"That's the law of total probability: P(event) = Σ P(event|scenario_i) × " +
						"P(scenario_i) — exactly the same Σxi·pi shape expected-value used, except " +
						"here the 'values' being weighted and averaged are themselves conditional " +
						"probabilities (each machine's own defect rate) rather than dollar payoffs, " +
						"and the 'probabilities' are how much of the whole population falls into " +
						"each scenario (each machine's share of output). bayes-theorem's denominator " +
						"was this exact formula with only two scenarios (condition present / " +
						"condition absent) instead of three.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"One long bar splits into three segments, left to right, one per machine — " +
						"segment width is that machine's share of total output, and its shade of red " +
						"deepens with its own defect rate. shareA and shareB are sliders (machine " +
						"C's share fills in whatever's left over); rateB lets you dial the problem " +
						"machine's defect rate up or down. The arithmetic underneath spells out each " +
						"machine's contribution and their sum. Below that, a number line marks two " +
						"points: the correctly weighted overall rate (green) and the naive " +
						"unweighted average of the three rates (orange, labeled wrong) — watch how " +
						"far apart they sit when shareB is small (machine B's terrible rate barely " +
						"moves the true overall number) versus how close they get if you drag shareB " +
						"up toward equal footing with A.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Decompose any tangled aggregate rate or probability — an overall defect rate, " +
						"an overall disease prevalence, an overall click-through rate — into simpler, " +
						"per-scenario pieces you can actually reason about, then recombine them " +
						"correctly, weighted by how common each scenario really is. You can also now " +
						"see exactly why bayes-theorem's denominator worked the way it did: it wasn't " +
						"an arbitrary formula, it was this same weighted-combination pattern applied " +
						"to two scenarios instead of three (or more).",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Overall disease prevalence across age groups or regions with different local " +
						"rates, overall click-through rate across traffic sources that convert " +
						"differently, overall loan default rate across risk tiers, overall exam " +
						"average across sections taught by different instructors, and insurance loss " +
						"rates across customer risk tiers. Anywhere a single 'overall' number is " +
						"actually made up of several subgroups that don't behave the same way and " +
						"don't make up equal shares of the total.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the overall probability has to be weighted by how often each " +
						"scenario actually occurs, not treated as a plain average across scenarios — " +
						"a scenario that's rare should barely move the overall number, no matter how " +
						"extreme its own rate is.",
					"Not like this: averaging the three defect rates unweighted, (1%+20%+1%)/3 = " +
						"7.33%, and reporting that as 'the' overall defect rate. The correct, " +
						"weighted answer is 2.52% — nearly 3x lower — because Machine B's alarming " +
						"20% rate applies to only 8% of output; treating all three machines as " +
						"equally important wildly overstates how bad the real, whole-floor number " +
						"actually is.",
				},
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
