// Package simpsonsparadox visualizes Simpson's paradox: a comparison that
// goes one way inside every subgroup can flip the other way once the
// subgroups are combined, purely because of how many cases of each kind
// landed in each group. It uses the real 1986 kidney-stone-treatment dataset
// (Charig et al., BMJ) — Treatment A beats Treatment B on both small stones
// and large stones, yet loses overall, because doctors used A far more often
// on the harder, large-stone cases.
package simpsonsparadox

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "simpsons-paradox",
		Seq:   54,
		Title: "Simpson's paradox",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Two treatments for kidney stones are compared in a 1986 study: Treatment A " +
						"(open surgery) cures 93% of small-stone patients and 73% of large-stone " +
						"patients; Treatment B (a less invasive procedure) cures 87% of small-stone " +
						"patients and 69% of large-stone patients. A wins both comparisons. So which " +
						"treatment should a doctor pick? The obvious instinct is 'A, obviously — it " +
						"won every category.' But when the two hospitals tallied up all 700 patients " +
						"treated, Treatment A's overall cure rate was 78% and Treatment B's was 83% " +
						"— B looked better once you stopped splitting by stone size. How can a " +
						"treatment lose on every single subgroup and still win overall?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"It happens because the two treatments weren't handed identical caseloads. " +
						"Large-stone cases are simply harder to cure than small-stone cases — 73% and " +
						"69% are both worse than 93% and 87% — and doctors reached for open surgery " +
						"(A) disproportionately for those harder cases:",
					"• A treated 350 patients: 263 large-stone (75%) and only 87 small-stone (25%).",
					"• B treated 350 patients: 80 large-stone (23%) and 270 small-stone (77%).",
					"So A's overall number is dragged down because three-quarters of its patients " +
						"were the harder-to-cure kind, while B's is propped up because three-quarters " +
						"of its patients were the easier kind. Weigh it out: A's overall rate is " +
						"0.75×73% + 0.25×93% ≈ 78%. B's overall rate is 0.23×69% + 0.77×87% ≈ 83%. " +
						"Same per-group numbers as above, but now the mix of easy vs. hard cases each " +
						"treatment happened to see is doing most of the work — not which treatment is " +
						"actually better.",
					"That mix is the whole mechanism: 'severity mix' (what share of a treatment's " +
						"patients were the harder kind) determines how much its overall number gets " +
						"pulled toward its harder-case rate versus its easier-case rate. Change only " +
						"the mix, keep every per-group cure rate fixed, and the overall ranking can " +
						"flip without a single patient's outcome changing.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Two fixed bar-pairs show each treatment's real, unchanging cure rate on small " +
						"stones and on large stones — A wins both, always. A third bar-pair shows the " +
						"overall rate, which moves as you drag each treatment's 'share of large-stone " +
						"patients' slider. At the real study's values (A: 75% large-stone, B: 23% " +
						"large-stone) the overall bars cross over and B's overall bar is taller, " +
						"reproducing the paradox. Pull both sliders to the same value and the overall " +
						"bars snap back into agreement with the two subgroup bars — A wins again, " +
						"because now both treatments faced the same mix of easy and hard cases.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Before checking whether groups being compared had a similar mix of easy and " +
						"hard cases, an overall percentage — a treatment's cure rate, a school's pass " +
						"rate, a hiring pipeline's offer rate — cannot safely be trusted. Now you know " +
						"the fix: don't stop at the combined number, and don't stop at the per-group " +
						"numbers alone either — check whether the groups being combined actually had " +
						"a comparable mix of the underlying cases. If they didn't, the combined " +
						"ranking can be telling you about the mix, not about which option is better.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"The most famous case is the 1973 UC Berkeley graduate admissions data: men " +
						"were admitted at a higher overall rate than women, which looked like bias — " +
						"but within almost every individual department, women were admitted at an " +
						"equal or higher rate. Women had disproportionately applied to the most " +
						"competitive departments, which pulled their overall number down. The same " +
						"shape shows up in batting averages (a player can hit better than a teammate " +
						"in both halves of a season and still trail in the full-season average, if " +
						"the number of at-bats differed between halves), in A/B tests split unevenly " +
						"across mobile and desktop traffic, and in hospital survival-rate comparisons " +
						"that don't account for patient severity at admission.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'Treatment A wins both subgroups, but B has the better overall " +
						"number because A took on far more of the hard cases' — naming the case-mix " +
						"difference as the reason the two rankings disagree.",
					"Not like this: 'The overall number is what matters, so B is the better " +
						"treatment' — treating the combined percentage as the final word without " +
						"asking whether it was computed over comparable groups. It's also a mistake " +
						"to swing the other way and say 'always trust the subgroup numbers instead' " +
						"— the right response is to check whether the case mix differed, not to " +
						"blindly prefer one level of aggregation over the other.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "mixA", Label: "Treatment A: share of large-stone patients", Min: 0, Max: 1, Step: 0.05, Def: 0.75},
			{Key: "mixB", Label: "Treatment B: share of large-stone patients", Min: 0, Max: 1, Step: 0.05, Def: 0.25},
		},
		Render: render,
	})
}

// Real per-subgroup cure rates from the 1986 Charig et al. kidney-stone
// study (BMJ). Treatment A (open surgery) beats Treatment B (percutaneous
// nephrolithotomy) on both stone sizes.
const (
	RateASmall = 0.93 // Treatment A, small stones
	RateALarge = 0.73 // Treatment A, large stones
	RateBSmall = 0.87 // Treatment B, small stones
	RateBLarge = 0.69 // Treatment B, large stones
)

// CombinedRate blends a treatment's small-stone and large-stone cure rates
// into a single overall rate, weighted by mixLarge — the fraction of that
// treatment's patients who had a large stone. mixLarge=0 returns rateSmall
// exactly (no large-stone patients); mixLarge=1 returns rateLarge exactly
// (no small-stone patients); values in between interpolate linearly. This is
// the whole mechanism behind Simpson's paradox: two treatments can have
// identical per-group rates yet different overall rates purely because
// mixLarge differs between them.
func CombinedRate(mixLarge, rateSmall, rateLarge float64) float64 {
	if mixLarge < 0 {
		mixLarge = 0
	}
	if mixLarge > 1 {
		mixLarge = 1
	}
	return (1-mixLarge)*rateSmall + mixLarge*rateLarge
}

func render(p map[string]float64) string {
	mixA, mixB := p["mixA"], p["mixB"]

	overallA := CombinedRate(mixA, RateASmall, RateALarge)
	overallB := CombinedRate(mixB, RateBSmall, RateBLarge)

	c := viz.New(680, 430, 0, 1, 0, 1)

	const chartTop, chartBottom = 56.0, 330.0
	const chartHeight = chartBottom - chartTop
	const groupW, barW, barGap = 200.0, 64.0, 16.0
	const startX = 40.0

	bar := func(x, rate float64, color string) {
		h := rate * chartHeight
		y := chartBottom - h
		c.Rect(x, y, barW, h, color, 0.9)
		c.Text(x+barW/2, y-8, fmt.Sprintf("%.0f%%", rate*100), 13, viz.Ink, "middle")
	}

	groups := []struct {
		label        string
		rateA, rateB float64
	}{
		{"Small stones", RateASmall, RateBSmall},
		{"Large stones", RateALarge, RateBLarge},
		{"Overall", overallA, overallB},
	}

	for i, g := range groups {
		gx := startX + float64(i)*groupW
		bar(gx, g.rateA, viz.Accent)
		bar(gx+barW+barGap, g.rateB, viz.Warm)
		c.Text(gx+barW+(barW+barGap)/2, chartBottom+22, g.label, 13, viz.Ink, "middle")
	}

	// Baseline under all bars.
	c.Rect(startX-10, chartBottom, groupW*3-10, 1.5, viz.Ink, 1)

	// Legend.
	c.Rect(470, 8, 14, 14, viz.Accent, 0.9)
	c.Text(490, 19, "Treatment A", 13, viz.Ink, "start")
	c.Rect(470, 28, 14, 14, viz.Warm, 0.9)
	c.Text(490, 39, "Treatment B", 13, viz.Ink, "start")

	c.Text(20, 19, fmt.Sprintf("A's patients: %.0f%% large-stone", mixA*100), 13, viz.Muted, "start")
	c.Text(20, 39, fmt.Sprintf("B's patients: %.0f%% large-stone", mixB*100), 13, viz.Muted, "start")

	note := "No paradox: A wins both subgroups and overall"
	noteColor := viz.Good
	if overallB > overallA {
		note = "Paradox: A wins both subgroups, but B wins overall"
		noteColor = viz.Bad
	}
	c.Text(340, 400, note, 14, noteColor, "middle")

	return c.String()
}
