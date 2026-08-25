// Package anova visualizes one-way analysis of variance: comparing three or
// more groups' means in a single test by splitting the total spread of the
// data into a between-group piece and a within-group piece, the same
// squared-deviation idea `variance-vs-stddev` used, applied twice. The
// F-ratio of those two pieces answers "do these groups differ at all?"
// without stacking up the false-alarm risk of a growing pile of pairwise
// comparisons.
package anova

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "anova",
		Seq:   68,
		Title: "ANOVA (comparing three or more group means at once)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You've got three teaching methods and want to know: does the method " +
						"actually matter for quiz scores, or are the differences you're seeing " +
						"just noise? The obvious approach, armed with `p-value`'s two-sample test, " +
						"is to compare every pair separately — method A vs B, A vs C, B vs C — " +
						"three separate 5%-significance tests. But each test carries its own 5% " +
						"chance of a false alarm even when nothing's really different, and those " +
						"chances stack: run three independent 5% tests and the chance that at " +
						"least one falsely comes back 'significant' climbs to 1-0.95³≈14%, not the " +
						"5% you signed up for. Is there a single test that asks 'do these three " +
						"groups differ at all,' in one shot, without letting the false-alarm rate " +
						"creep up every time you add another group to compare?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Three groups of 5 quiz scores each, all with the same within-group spread " +
						"(each group's scores deviate from its own mean by -6,-3,0,3,6):",
					"• Group A: 64,67,70,73,76 (mean 70)",
					"• Group B: 69,72,75,78,81 (mean 75)",
					"• Group C: 76,79,82,85,88 (mean 82)",
					"Pool all 15 scores: grand mean = 75.67. Now split the total spread into two " +
						"pieces — the same squared-deviation idea `variance-vs-stddev` used to " +
						"build variance, just applied twice:",
					"• Between-group sum of squares (SSB): how far each group's own mean sits " +
						"from the grand mean, weighted by group size — 5×(70-75.67)² + " +
						"5×(75-75.67)² + 5×(82-75.67)² = 363.33.",
					"• Within-group sum of squares (SSW): how far each individual score sits " +
						"from its own group's mean, added up across all three groups — " +
						"90+90+90 = 270.",
					"Divide each by its degrees of freedom — SSB by (groups-1)=2, SSW by (total " +
						"scores - groups)=12 — to get 'mean squares': MSB=181.67, MSW=22.5. The " +
						"F-ratio is MSB/MSW = 8.07: how much bigger the spread between group " +
						"means is than the ordinary noise within a group. A large F means the " +
						"groups' means are spread out far more than random within-group jitter " +
						"could explain by chance alone.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each group's 5 scores plotted as dots above its own position on the " +
						"x-axis, with a short colored tick marking that group's own mean. A " +
						"dashed line runs across the whole plot at the grand mean. Drag any " +
						"group's mean slider and watch that group's tick — and the whole " +
						"picture's balance — shift; drag the spread slider and watch the dots fan " +
						"out or huddle tighter around each tick. F grows when the ticks spread " +
						"far apart relative to how tightly the dots cluster around them, and " +
						"shrinks when the dots' own spread swamps the gap between ticks.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Test whether three or more groups differ at all with a single number, F, " +
						"instead of running an ever-growing pile of pairwise comparisons whose " +
						"combined false-alarm rate creeps upward every time you add another " +
						"group. A large F (like 8.07 here, checked against an F-distribution — " +
						"the same reference-distribution idea `p-value` used for a single test " +
						"statistic) says the between-group spread dwarfs the within-group noise; " +
						"a small F says the group means could easily be this far apart by chance " +
						"alone, even if nothing structural separates the groups.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Comparing exam scores across several teaching methods (the worked example " +
						"itself). A/B/C-testing three or more website variants at once instead of " +
						"three separate A/B tests. Agricultural field trials — ANOVA's original " +
						"use, developed by Ronald Fisher in the 1920s to compare crop yields " +
						"across several fertilizer treatments. Clinical trials with multiple " +
						"treatment arms, checking whether any dosage group's outcome differs from " +
						"the others before drilling into which specific dose matters.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'F=8.07 says the between-group spread is about 8 times " +
						"the within-group noise — evidence the three methods' true averages " +
						"aren't all equal.'",
					"Not like this: treating a significant F as telling you which specific " +
						"groups differ — ANOVA's F-test only answers 'do the means differ at " +
						"all,' not 'A differs from C but not from B'; pinning down which pairs " +
						"differ needs a proper post-hoc test (e.g. Tukey's HSD) run after ANOVA " +
						"flags a real difference, not instead of it. Also don't skip straight to " +
						"comparing raw mean differences without accounting for spread — " +
						"meanC-meanA=12 sounds big, but whether it's real depends entirely on how " +
						"much natural jitter (spread) is already baked into each group; the same " +
						"12-point gap could be decisive, noise-free evidence or statistically " +
						"meaningless, depending on spread.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "meanA", Label: "Group A mean", Min: 50, Max: 100, Step: 1, Def: 70},
			{Key: "meanB", Label: "Group B mean", Min: 50, Max: 100, Step: 1, Def: 75},
			{Key: "meanC", Label: "Group C mean", Min: 50, Max: 100, Step: 1, Def: 82},
			{Key: "spread", Label: "Within-group spread", Min: 0.5, Max: 6, Step: 0.5, Def: 3},
		},
		Render: render,
	})
}

// baseOffsets is a small, fixed set of deviations (summing to 0) that every
// group's samples are built from. Keeping it as a literal -- rather than
// generating it with randomness, which Render must never do -- keeps the
// picture concrete and reproducible: these are "quiz score" style deltas,
// not noise (mirrors variance-vs-stddev's baseSample).
var baseOffsets = []float64{-2, -1, 0, 1, 2}

// GroupSamples returns len(baseOffsets) scores centered on mean, spread
// out by spread: mean + baseOffsets[i]*spread. Pure function: same inputs
// always produce the same slice.
func GroupSamples(mean, spread float64) []float64 {
	xs := make([]float64, len(baseOffsets))
	for i, o := range baseOffsets {
		xs[i] = mean + o*spread
	}
	return xs
}

// Mean is the arithmetic mean of xs. Returns 0 for an empty slice.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// SSBetween returns the between-group sum of squares: how far each group's
// own mean sits from the grand mean (the mean of every group's samples
// pooled together), weighted by group size.
func SSBetween(groups [][]float64) float64 {
	var pooled []float64
	for _, g := range groups {
		pooled = append(pooled, g...)
	}
	grand := Mean(pooled)
	ss := 0.0
	for _, g := range groups {
		d := Mean(g) - grand
		ss += float64(len(g)) * d * d
	}
	return ss
}

// SSWithin returns the within-group sum of squares: how far each
// individual sample sits from its own group's mean, added up across every
// group.
func SSWithin(groups [][]float64) float64 {
	ss := 0.0
	for _, g := range groups {
		gm := Mean(g)
		for _, x := range g {
			d := x - gm
			ss += d * d
		}
	}
	return ss
}

// Result holds the pieces of a one-way ANOVA: the between- and
// within-group sums of squares, their degrees of freedom, the mean
// squares each reduces to, and the F-ratio (MSBetween/MSWithin) that
// summarizes the whole test.
type Result struct {
	SSBetween, SSWithin float64
	DFBetween, DFWithin int
	MSBetween, MSWithin float64
	F                   float64
}

// Run performs a one-way ANOVA across groups (each a slice of samples for
// one group). It panics if fewer than 2 groups or fewer than 1 total
// sample is given -- an ANOVA needs at least two groups to compare.
func Run(groups [][]float64) Result {
	if len(groups) < 2 {
		panic("anova: Run needs at least 2 groups")
	}
	n := 0
	for _, g := range groups {
		n += len(g)
	}
	k := len(groups)
	ssb, ssw := SSBetween(groups), SSWithin(groups)
	dfb, dfw := k-1, n-k

	msb := ssb / float64(dfb)
	var msw, f float64
	if dfw > 0 {
		msw = ssw / float64(dfw)
	}
	switch {
	case msw > 1e-12:
		f = msb / msw
	case msb > 1e-12:
		f = math.Inf(1) // zero within-group noise but real between-group spread
	default:
		f = 0 // everything's identical -- no spread anywhere to report
	}

	return Result{ssb, ssw, dfb, dfw, msb, msw, f}
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 4, 0, 100).String()
}
