// Package fibonacci visualizes the ratio of consecutive Fibonacci numbers,
// F(n+1)/F(n), oscillating in on the golden ratio φ = (1+√5)/2 as n grows —
// a fixed point that quadratic-formula solves for directly, in one step,
// once the Fibonacci recurrence is rewritten as an equation in the ratio
// itself.
package fibonacci

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// Phi is the golden ratio, (1+√5)/2 — the positive root of x²−x−1=0, the
// fixed point of x ↦ 1+1/x that the Fibonacci ratio converges to.
var Phi = (1 + math.Sqrt(5)) / 2

func init() {
	concept.Register(concept.Concept{
		ID:    "fibonacci-golden-ratio",
		Seq:   62,
		Title: "Fibonacci sequence & the golden ratio (φ)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`quadratic-formula` just gave us a one-step way to solve any equation " +
						"shaped like ax²+bx+c=0. Where would an equation like that show up on its " +
						"own, outside a textbook? Look at the Fibonacci sequence — 1, 1, 2, 3, 5, 8, " +
						"13, 21, 34... each term the sum of the two before it — and its ratio of " +
						"consecutive terms: 1/1=1, 2/1=2, 3/2=1.5, 5/3≈1.667, 8/5=1.6, 13/8=1.625, " +
						"21/13≈1.615, 34/21≈1.619... It's obviously homing in on some number near " +
						"1.618, bouncing above and below it rather than climbing smoothly. What " +
						"number is it converging to, and why that number specifically?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Start from the Fibonacci rule itself, F(n+1) = F(n) + F(n−1), and divide both " +
						"sides by F(n): F(n+1)/F(n) = 1 + F(n−1)/F(n) = 1 + 1/(F(n)/F(n−1)). Calling " +
						"the ratio r(n) = F(n+1)/F(n), this says r(n) = 1 + 1/r(n−1) — each ratio is " +
						"completely determined by the one before it through this one rule. Check it " +
						"directly: r(1)=1, r(2)=1+1/1=2 ✓ (matches 2/1 above), r(3)=1+1/2=1.5 ✓ " +
						"(matches 3/2), r(4)=1+1/1.5≈1.667 ✓.",
					"If the ratios are settling down at all, they must be settling on a value x " +
						"where feeding x into that rule gives back x itself: x = 1 + 1/x. Multiply " +
						"both sides by x: x² = x + 1, i.e. x² − x − 1 = 0. That's exactly the shape " +
						"`quadratic-formula` solves in one step: a=1, b=−1, c=−1, so D = (−1)² − " +
						"4×1×(−1) = 1+4 = 5, and x = (1 ± √5)/2. The positive root, φ = (1+√5)/2 ≈ " +
						"1.6180339887, is the golden ratio — the only positive number the always-" +
						"positive Fibonacci ratios could possibly be converging to; the negative " +
						"root, (1−√5)/2 ≈ −0.618, solves the same equation but isn't where a ratio " +
						"of positive numbers can land.",
					"• n=1: F(1)=1, F(2)=1, ratio=1.000, off from φ by −0.618.",
					"• n=5: F(5)=5, F(6)=8, ratio=1.600, off from φ by −0.018.",
					"• n=10: F(10)=55, F(11)=89, ratio=1.61818, off from φ by +0.00015.",
					"• n=15: F(15)=610, F(16)=987, ratio=1.6180328, off by −0.0000012.",
					"Each step's overshoot/undershoot shrinks by roughly the same factor (≈1/φ²) " +
						"that produced the last one, and flips sign every time — a spiral closing in " +
						"on φ from alternating sides, not a straight climb toward it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Each point is r(n) = F(n+1)/F(n); the dashed line marks φ itself. Early terms " +
						"swing widely above and below the line (r(1)=1, r(2)=2), then the swings " +
						"visibly tighten in — by around n=10 the curve is already indistinguishable " +
						"from flat at this scale. The n slider highlights one term and reads off the " +
						"exact fraction F(n+1)/F(n), its decimal value, and how far it currently sits " +
						"from φ.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Approximate φ to arbitrary precision using nothing but integer addition — no " +
						"square root needed — by just running the Fibonacci recurrence far enough " +
						"and taking a late ratio; useful anywhere division and √ are expensive or " +
						"unavailable but addition is cheap. You can also derive φ's own algebraic " +
						"quirks directly from x²=x+1 instead of memorizing them: dividing by x gives " +
						"x = 1 + 1/x, so 1/φ = φ−1 ≈ 0.618 (φ is the only positive number whose " +
						"reciprocal is exactly itself minus 1); and φ² = φ+1 ≈ 2.618 falls straight " +
						"out of the original equation.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Sunflower seed heads and pinecones pack their spirals in counts that are " +
						"consecutive Fibonacci numbers (commonly 34 and 55, or 55 and 89) because " +
						"packing new growth at a turning angle of 360°/φ² per step — the 'golden " +
						"angle,' ≈137.5° — leaves the least overlap between successive seeds, a " +
						"real, measurable botanical pattern (phyllotaxis). Technical stock-market " +
						"analysis uses 'Fibonacci retracement' levels (23.6%, 38.2%, 61.8%) derived " +
						"from φ and its powers to guess where a price move might pause. Some " +
						"seashells (like the chambered nautilus) grow in a logarithmic spiral whose " +
						"growth factor is close to, though not exactly, φ per quarter-turn.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the ratio of consecutive Fibonacci terms approaches φ as n " +
						"grows, getting closer and closer but never exactly equal to it for any " +
						"finite n — 21/13 is close to φ, not equal to it.",
					"Not like this: treating φ as a universal design law baked into nature and " +
						"historic art — 'the Parthenon's proportions are the golden ratio,' 'the " +
						"perfect human face/body follows φ' — claims that don't hold up under actual " +
						"measurement and mostly reflect after-the-fact rectangle-fitting rather than " +
						"documented intent or a real pattern. The sunflower/pinecone spiral counts " +
						"above are a solid, independently verified example; treat 'I found something " +
						"shaped roughly 1.6-to-1 in this famous object' with real skepticism instead " +
						"of taking it as more evidence for the same claim.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Term index (n)", Min: 1, Max: 20, Step: 1, Def: 8},
		},
		Render: render,
	})
}

// Fibonacci returns F(n) for n>=0: F(0)=0, F(1)=1, F(n)=F(n-1)+F(n-2).
// Iterative, not recursive, so it stays O(n) with no repeated work.
func Fibonacci(n int) int64 {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// Ratio returns F(n+1)/F(n), the ratio of consecutive Fibonacci terms. n
// must be >=1 (F(0)=0 makes the ratio undefined at n=0); returns 0 for n<1.
func Ratio(n int) float64 {
	if n < 1 {
		return 0
	}
	return float64(Fibonacci(n+1)) / float64(Fibonacci(n))
}

// FixedPointStep applies x ↦ 1+1/x, the map the Fibonacci ratio recurrence
// reduces to (r(n) = 1 + 1/r(n-1)). Its only fixed point among positive
// numbers is Phi: FixedPointStep(Phi) == Phi, up to floating-point error.
func FixedPointStep(x float64) float64 {
	if x == 0 {
		return math.Inf(1)
	}
	return 1 + 1/x
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 380, 0, 20, 0, 2).String()
}
