// Package bonferroni visualizes the Bonferroni correction: when you run many
// hypothesis tests at once, even if every single null hypothesis is
// perfectly true, testing each one at the usual significance level (like
// α=0.05) lets false positives pile up across the whole batch. Dividing α by
// the number of tests, m, before applying it to any one test keeps the
// probability of even one false positive across the whole family bounded by
// the original α.
package bonferroni

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "bonferroni-correction",
		Seq:   104,
		Title: "Bonferroni correction (adjusting for multiple tests)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"p-value established that testing a single hypothesis at significance level " +
						"\u03b1=0.05 means a 5% chance of a false alarm -- flagging something as real when " +
						"nothing actually is -- if the null hypothesis happens to be true. That sounds " +
						"like an acceptably small risk taken on its own.",
					"But real research rarely stops at one test. A genetics study checks thousands " +
						"of genes for a disease link; a company's dashboard checks dozens of metrics " +
						"after a product change; a medical panel checks many possible conditions at " +
						"once. Gut instinct says: 'each individual test is only 5% risky, so I'm still " +
						"safe.' But run enough separate 5%-risk tests, even when NOTHING real is going " +
						"on in any of them, and by sheer chance some of them will cross the 5% line " +
						"anyway -- purely from noise, not from any genuine effect.",
					"How much does that risk actually pile up as you run more tests, and is there a " +
						"way to set each test's own threshold so the risk of even ONE false alarm across " +
						"the WHOLE batch stays controlled, instead of quietly growing with every test " +
						"you add?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Run m independent tests, each at significance level \u03b1, and suppose every " +
						"single null hypothesis is actually true (nothing real is going on anywhere). " +
						"Each test individually has probability (1-\u03b1) of correctly NOT crossing the " +
						"threshold. Independence lets those probabilities multiply across all m tests, " +
						"so the chance that EVERY test stays quiet is (1-\u03b1)^m -- and the chance that " +
						"at least one test falsely fires is 1 minus that: the family-wise error rate, " +
						"FWER = 1-(1-\u03b1)^m.",
					"At \u03b1=0.05, plug in a few values of m:",
					"|m (tests)|FWER = 1-(1-\u03b1)^m|",
					"|1|5.0%|",
					"|5|22.6%|",
					"|10|40.1%|",
					"|20|64.2%|",
					"|30|78.5%|",
					"|50|92.3%|",
					"By 20 tests, the chance of at least one false alarm has already crossed a coin " +
						"flip -- worse than 50/50 -- purely from running more tests at the 'safe-looking' " +
						"5% level, with nothing real happening anywhere. That's the gap Section 1 named, " +
						"now with an exact number attached.",
					"Bonferroni's fix: instead of testing every one of the m tests at \u03b1, test each " +
						"one at the stricter \u03b1/m. Take a concrete batch: simulate 20 p-values, every " +
						"one of them from a genuinely true null (nothing real in any of the 20). At " +
						"\u03b1=0.05, one of those 20 p-values (0.0167) happens to fall below the raw " +
						"threshold -- a false alarm, exactly the kind of noise the 64.2% FWER above says " +
						"to expect fairly often at m=20. At the corrected threshold \u03b1/m = 0.05/20 = " +
						"0.0025, that same p-value (0.0167) is nowhere close to crossing it -- 0 false " +
						"alarms in this batch, even though all 20 tests came from a situation where " +
						"nothing at all was actually going on.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The top curve traces FWER = 1-(1-\u03b1)^m across m=1 to 50 for the current " +
						"\u03b1, with a marker at the m the slider is set to -- drag m up and watch the " +
						"marker climb the curve, the same numbers tabled above. Below it, the current " +
						"batch of m simulated p-values (all drawn from a true null, so any one flagged " +
						"is by definition a false alarm) sits as dots along a 0-to-1 line, next to two " +
						"threshold markers: the raw \u03b1 (orange, fixed) and the corrected \u03b1/m " +
						"(green, sliding left as m grows). Watch a dot that used to sit between the two " +
						"lines -- caught by the raw threshold but not the corrected one -- as exactly the " +
						"false alarm the correction was built to catch.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Run many hypothesis tests within one study or one decision and still control the " +
						"overall risk of ANY false conclusion, instead of unconsciously accepting a risk " +
						"that grows with every extra test -- checking 5, 20, or 50 outcomes at once " +
						"without fooling yourself into 'finding' an effect that's really just noise. " +
						"That comes with a real tradeoff, worth being upfront about: making every " +
						"individual test's bar much stricter also makes it harder to catch a real effect " +
						"when one genuinely exists (Bonferroni is known as a conservative correction) -- " +
						"but for a chosen, fixed overall risk level, it delivers exactly the guarantee it " +
						"promises.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Genetics studies that scan thousands of genes at once for a disease " +
						"association (a genome-wide association study), where the raw per-gene " +
						"threshold would produce a near-certain flood of false hits without correction. " +
						"A/B testing platforms that track dozens of metrics on one experiment and need " +
						"to avoid declaring victory on whichever metric happened to wiggle by chance. A " +
						"doctor ordering a broad panel of tests for a patient's symptoms, where checking " +
						"many possible conditions at once raises the same risk of a false positive " +
						"result on at least one of them. This is also the rigorous fix for exactly the " +
						"'p-hacking' danger p-value warned about: run enough different tests against the " +
						"same data and something will eventually look significant by chance alone -- " +
						"Bonferroni (and its relatives) is how careful analysis accounts for that instead " +
						"of pretending it isn't happening.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: Bonferroni divides your significance threshold by the NUMBER " +
						"OF TESTS you're running (\u03b1/m), controlling the chance of even ONE false " +
						"positive across the whole family of tests together.",
					"Not like this: applying the original threshold \u03b1 to each test separately " +
						"just because 'each individual test is only 5% risky' -- that exact reasoning " +
						"is what lets the family-wise risk balloon to 64.2% at just 20 tests and 92.3% " +
						"at 50, as this concept's own curve shows.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "m", Label: "Number of tests (m)", Min: 1, Max: 50, Step: 1, Def: 20},
			{Key: "alpha", Label: "Significance level (α)", Min: 0.01, Max: 0.20, Step: 0.01, Def: 0.05},
		},
		Render: render,
	})
}

// MaxM is the largest number of tests every Section's slider sweeps through.
const MaxM = 50

// hash01 turns an integer seed into a deterministic pseudo-random number in
// [0, 1) via a fixed irrational-multiplier trick -- not a real random number
// generator, just a fixed, reproducible sequence that merely looks
// patternless, which is exactly what Render needs (pure: same input, same
// output, no math/rand, no time). Same trick as internal/concepts/pvalue's
// sibling concepts (e.g. internal/concepts/correlation).
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// nullPValueOffset shifts every seed hash01 sees away from i=0 (where
// hash01 degenerates to exactly 0, since sin(0)=0) and lands on a batch of
// 20 simulated p-values that makes a genuinely instructive example: one
// clearly beats the raw threshold 0.05 (a false positive, since every one
// of these is drawn from a true null) while none beat the stricter
// Bonferroni threshold -- the exact contrast this concept exists to show.
const nullPValueOffset = 50

// NullPValue returns the (i+1)-th simulated p-value (0-indexed i) from a
// batch of independent tests where every null hypothesis is actually true.
// Under a true null, a p-value is uniformly distributed on [0,1) by
// definition -- that's the whole reason a raw significance level like
// α=0.05 means "a false positive 5% of the time" -- so hash01 is reused
// directly as one. Deterministic and pure: the i-th test always gets the
// same simulated p-value.
func NullPValue(i int) float64 {
	return hash01(i + nullPValueOffset)
}

// NullPValues returns the first m simulated null p-values, NullPValue(0)
// through NullPValue(m-1).
func NullPValues(m int) []float64 {
	if m < 0 {
		m = 0
	}
	vals := make([]float64, m)
	for i := range vals {
		vals[i] = NullPValue(i)
	}
	return vals
}

// CountBelow counts how many p-values fall strictly below threshold -- how
// many of a batch of tests would be called "significant" (a false positive,
// since every one of these is drawn from a true null) at that threshold.
func CountBelow(pvalues []float64, threshold float64) int {
	n := 0
	for _, pv := range pvalues {
		if pv < threshold {
			n++
		}
	}
	return n
}

// FamilyWiseErrorRate is the probability that at least one of m independent
// tests, each run at significance level alpha, comes back a false positive,
// given every null hypothesis is actually true. Each individual test has
// probability (1-alpha) of correctly NOT rejecting; independence lets those
// multiply across all m tests, so 1 minus that product is the chance at
// least one test slips through: FWER = 1 - (1-alpha)^m.
func FamilyWiseErrorRate(alpha float64, m int) float64 {
	return 1 - math.Pow(1-alpha, float64(m))
}

// BonferroniAlpha is the corrected per-test significance level: the target
// family-wise level alpha, split evenly across m tests. Testing each of the
// m tests at this stricter threshold instead of the raw alpha keeps the
// overall family-wise error rate at or below alpha (rather than growing
// with m the way FamilyWiseErrorRate(alpha, m) does at the raw threshold).
func BonferroniAlpha(alpha float64, m int) float64 {
	if m < 1 {
		m = 1
	}
	return alpha / float64(m)
}

func render(p map[string]float64) string {
	m := int(p["m"] + 0.5)
	if m < 1 {
		m = 1
	}
	if m > MaxM {
		m = MaxM
	}
	alpha := p["alpha"]
	if alpha < 0.01 {
		alpha = 0.01
	}
	if alpha > 0.20 {
		alpha = 0.20
	}
	corrected := BonferroniAlpha(alpha, m)
	fwer := FamilyWiseErrorRate(alpha, m)
	pvalues := NullPValues(m)
	rawHits := CountBelow(pvalues, alpha)
	correctedHits := CountBelow(pvalues, corrected)

	// Two panels share one Canvas by mutating its data-space fields
	// in between -- draw the FWER curve first (x = number of tests,
	// y = FWER), then repoint XMin/XMax/YMin/YMax and PadT/PadB at the
	// p-value strip below it. Canvas's X()/Y() read those fields live, so
	// each panel gets its own correct mapping without needing a second
	// Canvas or any nested SVG.
	c := viz.New(700, 560, 1, float64(MaxM), 0, 1.05)
	c.PadT = 110
	c.PadB = 300
	c.Axes()
	for x := 1; x <= MaxM; x += 7 {
		c.Tick(float64(x), fmt.Sprintf("%d", x))
	}
	c.Text((c.W+c.PadL)/2-60, c.H-c.PadB+34, "number of tests (m)", 12, viz.Muted, "start")

	curve := viz.Sample(1, float64(MaxM), MaxM-1, func(mm float64) float64 {
		return FamilyWiseErrorRate(alpha, int(mm+0.5))
	})
	c.Path(curve, viz.Accent, 2.5)
	// Reference: the raw alpha itself, i.e. what the family-wise error rate
	// would be if it didn't grow with m at all.
	c.Path([][2]float64{{1, alpha}, {float64(MaxM), alpha}}, viz.Muted, 1)
	c.VLine(float64(m), viz.Warm, true)
	mx, my := c.X(float64(m)), c.Y(fwer)
	c.Rect(mx-4, my-4, 8, 8, viz.Warm, 1)

	// Second panel: the m simulated p-values (every one drawn from a TRUE
	// null) plotted along [0,1], against both thresholds.
	c.XMin, c.XMax = 0, 1
	c.YMin, c.YMax = 0, 1
	c.PadT = 340
	c.PadB = 40
	c.Axes()
	for x := 0.0; x <= 1.0; x += 0.2 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}
	c.Text((c.W+c.PadL)/2-40, c.H-c.PadB+34, "simulated p-value", 12, viz.Muted, "start")

	const rows = 8
	for i, pv := range pvalues {
		row := float64(i%rows) / float64(rows-1)
		color := viz.Accent
		switch {
		case pv < corrected:
			color = viz.Bad
		case pv < alpha:
			color = viz.Warm
		}
		px, py := c.X(pv), c.Y(row)
		c.Rect(px-3, py-3, 6, 6, color, 0.85)
	}
	c.VLine(alpha, viz.Warm, true)
	c.VLine(corrected, viz.Good, true)

	c.Text(16, 24, fmt.Sprintf("m = %d tests    raw α = %.2f    corrected α/m = %.4f", m, alpha, corrected),
		13, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("family-wise error rate at raw α ≈ %.1f%% (chance of ≥ 1 false positive)",
		fwer*100), 13, viz.Warm, "start")
	c.Text(16, 64, fmt.Sprintf("this batch (all %d from a TRUE null): %d below raw α    %d below corrected α/m",
		m, rawHits, correctedHits), 13, viz.Ink, "start")
	c.Text(16, 84, "orange dashed = raw α    green dashed = corrected α/m    red dot = still false after correction",
		12, viz.Muted, "start")

	return c.String()
}
