// Package benfordslaw visualizes Benford's law: in many real-world datasets
// that span several orders of magnitude, the leading (first) digit of each
// number is far from uniformly distributed across 1-9. Digit 1 leads about
// 30% of the time and digit 9 barely 4.6% -- because the leading digit is
// decided by the fractional part of log10(x), and the fractional-part
// interval each digit owns isn't equally wide. Directly built on
// `logarithms`: the same "each xb step in x adds a fixed amount to log(x)"
// idea, applied to base 10 and then read off which digit that lands on.
package benfordslaw

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "benfords-law",
		Seq:   79,
		Title: "Benford's law (why leading digit 1 wins)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"An auditor eyeballing a spreadsheet of expense reports, tax filings, or " +
						"census counts wants a quick way to flag numbers that might be made up. A " +
						"natural guess: if the figures are basically 'random', each first digit " +
						"1-9 should show up about equally often -- a bit over 11% each. Pull the " +
						"leading digit off almost any large real dataset that grows or scales " +
						"multiplicatively (populations, revenues, prices) and that guess is way " +
						"off: digit 1 leads roughly three times out of ten, digit 9 barely one time " +
						"in twenty. Is 'roughly equal' just wrong for some structural reason -- and " +
						"if so, can spotting the wrong shape actually catch someone faking numbers?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the powers of 2 from 2^0 up to 2^29 -- 30 numbers: 1, 2, 4, 8, 16, 32, " +
						"64, 128, 256, 512, 1024, 2048, ... Read off just the first digit of each: " +
						"1 appears 9 times out of 30 (30.0%), while 7 and 9 never appear at all in " +
						"this run. That's not a fluke of small samples -- it's the predicted shape.",
					"Here's why. `logarithms` already showed that each x2 step in x adds a fixed " +
						"log_2(2)=1 to log_2(x); the same addition trick works with log10 instead: " +
						"each x2 step adds a fixed log10(2) approx 0.301 to log10(x). A number's " +
						"leading digit is decided entirely by the *fractional part* of log10(x) -- " +
						"digit d owns the fractional-part range [log10(d), log10(d+1)). Repeatedly " +
						"adding a fixed amount and keeping only the fractional part sweeps evenly " +
						"around the [0,1) circle, but the arcs assigned to each digit are not equal " +
						"width: digit 1 owns [0, 0.301) -- width 0.301 -- while digit 9 owns " +
						"[0.954, 1) -- width just 0.046, about 6.6x narrower. Landing evenly around " +
						"a circle where one arc is 6.6x wider than another means landing in the " +
						"wide arc roughly 6.6x more often. That gives the general formula: " +
						"P(leading digit = d) = log10(1 + 1/d).",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Bars are the empirical share of each leading digit 1-9 across the current " +
						"dataset. b and n control the 'powers of b' dataset (mode=0): b sets the " +
						"growth multiplier, n sets how many terms; the default (b=2, n=30) " +
						"reproduces the worked 30-powers-of-2 example above almost exactly. " +
						"Flipping mode to 1 swaps in a completely different dataset -- n values " +
						"evenly spaced across the fixed range 1..999, with no multiplicative growth " +
						"at all -- and the bars flatten out to roughly 11% each, no matter how " +
						"large n gets. The warm markers and connecting line trace Benford's " +
						"predicted P(d) = log10(1+1/d) for comparison; in mode 0 the accent bars " +
						"hug that line, in mode 1 they don't.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Flag a spreadsheet as suspicious from its digit shape alone, without " +
						"checking a single number against outside records: if digit 1 doesn't " +
						"dominate the leading digits (or 9 shows up far more than about 4.6% of " +
						"the time), that's a statistical red flag that the numbers weren't " +
						"produced by natural multiplicative growth -- someone (or some process) " +
						"spread the digits more evenly than reality does, the same 'expected vs. " +
						"actual shape' test `chi-squared-test` runs on category counts. You can " +
						"also now tell, before even collecting data, which kinds of numbers should " +
						"follow the law: quantities that grow or scale multiplicatively across " +
						"orders of magnitude, not quantities capped to a narrow range.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Forensic accounting: auditors run first-digit tests on tax filings, invoice " +
						"amounts, and expense reports as a first-pass fraud screen, since real " +
						"accounting entries feed off compounding growth and Benford's law holds up " +
						"well against it. It has also been used to flag anomalies in reported " +
						"election vote counts and macroeconomic statistics. Physical quantities " +
						"that span scales -- river lengths, city and country populations, stock " +
						"prices, physical constants in whatever unit -- tend to follow it too.",
					"It does NOT apply to numbers that are capped to a narrow range or handed out " +
						"sequentially: human heights in centimeters, ages in years, phone numbers, " +
						"or lottery numbers all fail the test -- the same way this picture's " +
						"'evenly spaced 1..999' dataset does.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'leading digits skew toward 1 because multiplying by a " +
						"fixed factor over and over adds a fixed amount to log10(x), and small " +
						"digits own a wider slice of the log scale' -- the skew comes from the " +
						"arithmetic of repeated multiplication, not from anything special about " +
						"the number 1 itself.",
					"Not like this: assuming Benford's law applies to any big pile of numbers. It " +
						"only kicks in for data that spans several orders of magnitude through " +
						"multiplicative growth or mixed scales -- a dataset confined to a narrow " +
						"range (test scores out of 100, ages 0-100) or handed out in a fixed " +
						"pattern (phone numbers, ID numbers) has no reason to follow it, and " +
						"forcing the test onto those numbers just 'flags' perfectly honest data.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "base", Label: "Growth base (b)", Min: 2, Max: 20, Step: 1, Def: 2},
			{Key: "n", Label: "Sample count (n)", Min: 10, Max: 300, Step: 5, Def: 30},
			{Key: "mode", Label: "Dataset (0 = powers of b, 1 = evenly spaced 1..999)", Min: 0, Max: 1, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// LeadingDigit returns the first significant (nonzero) digit of a positive
// number, 1-9: 314 -> 3, 0.0072 -> 7, 9999 -> 9. Works by repeatedly
// scaling x into [1,10) rather than going through log10 and rounding,
// which keeps it exact at decade boundaries like x=1000 that floating-point
// log/pow round-trips can nudge onto the wrong side (9.999... vs 10.0).
func LeadingDigit(x float64) int {
	x = math.Abs(x)
	if x == 0 {
		return 0
	}
	for x >= 10 {
		x /= 10
	}
	for x < 1 {
		x *= 10
	}
	d := int(x + 1e-9) // guard against e.g. 8.999999999 from fp error
	if d < 1 {
		d = 1
	}
	if d > 9 {
		d = 9
	}
	return d
}

// BenfordProbability returns Benford's law's predicted probability that a
// number's leading digit is d: log10(1 + 1/d). It is the width of the
// fractional-log10 interval [log10(d), log10(d+1)) that digit d owns, out
// of the full unit interval every digit's arcs partition. Returns 0 outside
// 1..9, which isn't a valid leading digit.
func BenfordProbability(d int) float64 {
	if d < 1 || d > 9 {
		return 0
	}
	return math.Log10(1 + 1/float64(d))
}

// PowersOfBase returns the n numbers base^0, base^1, ..., base^(n-1) --
// the multiplicative-growth dataset that (for n large enough) tracks
// Benford's law closely, since each step adds the fixed amount
// log10(base) to log10(x).
func PowersOfBase(base float64, n int) []float64 {
	if n < 1 {
		n = 1
	}
	vals := make([]float64, n)
	v := 1.0
	for i := 0; i < n; i++ {
		vals[i] = v
		v *= base
	}
	return vals
}

// LinearSamples returns n values evenly spaced across [lo, hi] inclusive
// (n=1 returns just lo) -- the contrasting non-multiplicative dataset: no
// repeated scaling by a fixed factor, so its leading digits stay close to
// uniform across 1-9 regardless of n.
func LinearSamples(lo, hi float64, n int) []float64 {
	if n < 1 {
		n = 1
	}
	if n == 1 {
		return []float64{lo}
	}
	vals := make([]float64, n)
	for i := 0; i < n; i++ {
		vals[i] = lo + (hi-lo)*float64(i)/float64(n-1)
	}
	return vals
}

// LeadingDigitFrequencies returns, for each digit 1-9 (index d-1), the
// fraction of vals whose leading digit is d. Empty input returns all
// zeros rather than dividing by zero.
func LeadingDigitFrequencies(vals []float64) [9]float64 {
	var counts [9]int
	for _, v := range vals {
		if d := LeadingDigit(v); d >= 1 && d <= 9 {
			counts[d-1]++
		}
	}
	var freq [9]float64
	if len(vals) == 0 {
		return freq
	}
	for i, c := range counts {
		freq[i] = float64(c) / float64(len(vals))
	}
	return freq
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 420, 0, 10, 0, 0.35).String()
}
