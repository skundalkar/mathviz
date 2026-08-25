// Package anova visualizes one-way analysis of variance: comparing three or
// more groups' means in a single test by splitting the total spread of the
// data into a between-group piece and a within-group piece, the same
// squared-deviation idea `variance-vs-stddev` used, applied twice. The
// F-ratio of those two pieces answers "do these groups differ at all?"
// without stacking up the false-alarm risk of a growing pile of pairwise
// comparisons.
package anova

import (
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

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 4, 0, 100).String()
}
