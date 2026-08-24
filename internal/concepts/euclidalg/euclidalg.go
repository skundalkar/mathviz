// Package euclidalg visualizes the Euclidean algorithm: finding the
// greatest common divisor of two numbers by repeatedly taking a remainder
// (`modular-arithmetic`'s wraparound operation) instead of testing every
// possible divisor. The geometric picture — tiling an a×b rectangle with
// squares, recursing into whatever's left over — makes the algorithm
// visible as a spiral that always ends on a perfectly tiling square.
package euclidalg

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "euclidean-algorithm",
		Seq:   65,
		Title: "The Euclidean algorithm (finding GCD by remainders)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`modular-arithmetic` introduced 'a mod n', the remainder left over once a " +
						"wraps around n. Now put that remainder to work: to simplify the fraction " +
						"1071/462 to lowest terms, or to check whether two numbers picked for a " +
						"cryptographic key accidentally share a common factor, you need the " +
						"greatest common divisor (GCD) of two numbers — the largest number that " +
						"divides both evenly. The obvious approach is to list every divisor of each " +
						"number and compare, but that means testing candidates up to min(a,b) one " +
						"by one. For 1071 and 462 that's already tedious by hand; real " +
						"cryptographic numbers run to hundreds of digits, where a divisor-by-divisor " +
						"search would never finish. Is there a way to find the GCD without checking " +
						"every possible divisor?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Euclid's algorithm (Elements, Book VII, c. 300 BC): gcd(a,b) = gcd(b, a mod " +
						"b), repeated until the remainder hits 0 — the last nonzero value is the " +
						"answer. This works because any number that divides both a and b also " +
						"divides a mod b (since a mod b is just a minus some multiple of b), so " +
						"swapping in the remainder never changes what the GCD is — it only ever " +
						"shrinks the numbers you're working with, and a mod b is always strictly " +
						"smaller than b, so the process can't run forever.",
					"Walk gcd(1071, 462):",
					"• 1071 = 2×462 + 147, so gcd(1071,462) = gcd(462,147)",
					"• 462 = 3×147 + 21, so gcd(462,147) = gcd(147,21)",
					"• 147 = 7×21 + 0 — remainder hits 0, so gcd(1071,462) = 21, the last nonzero " +
						"remainder",
					"Three remainder steps found the GCD of two four-digit-ish numbers — nowhere " +
						"near testing all 462 candidate divisors. The algorithm's worst case is " +
						"exactly `fibonacci-golden-ratio`'s sequence: gcd(13,8) — consecutive " +
						"Fibonacci numbers — takes 5 steps despite being much smaller than 1071 and " +
						"462, because each Fibonacci remainder shrinks as slowly as possible (this " +
						"is Lamé's theorem). Even that worst case only needs a number of steps " +
						"proportional to the number of digits in the smaller number, never anything " +
						"close to a full divisor search.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"An a×b rectangle, tiled with squares: lay b×b squares along the long side " +
						"(one color per division step) until less than a full square's-width is " +
						"left over, then recurse into that leftover strip exactly the way the " +
						"remainder step above does — the same 1071=2×462+147 arithmetic, now drawn " +
						"as '2 orange squares of side 462, with a 147-wide strip left over.' The " +
						"tiling alternates horizontal and vertical as it recurses, spiraling inward, " +
						"and always finishes on one last square — side length exactly gcd(a,b) — " +
						"that tiles its remaining strip perfectly, with nothing left over.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Compute the GCD of numbers with hundreds of digits in a handful of steps " +
						"instead of an impossibly long divisor search — the core speed-up every " +
						"fraction-reduction routine and cryptographic library relies on. It's also " +
						"the base case of the extended Euclidean algorithm, which back-substitutes " +
						"through these same remainder steps to find modular inverses — the exact " +
						"operation `modular-arithmetic`'s real-life section pointed to as load-" +
						"bearing for RSA key generation.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every time software reduces a fraction to lowest terms (3/9 → 1/3), it " +
						"divides both numbers by their GCD, found this way. RSA key generation " +
						"repeatedly needs GCDs (checking that a chosen exponent is coprime to " +
						"another number) and uses the extended version to compute modular " +
						"inverses. It's also one of the oldest algorithms still in everyday use — " +
						"documented over 2,300 years ago and essentially unchanged since.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'gcd(1071,462)=21, because 1071 mod 462=147, 462 mod " +
						"147=21, and 147 mod 21=0 — the last nonzero remainder is the answer,' " +
						"naming the actual remainder chain rather than just the final number.",
					"Not like this: reaching for 'list all the divisors of each number and find " +
						"the biggest one in common' as a general method — it's correct, but it's " +
						"the slow approach this algorithm exists to replace, and it becomes " +
						"completely impractical once the numbers get large. A second common slip: " +
						"forgetting the base case gcd(a,0)=a, which is exactly what the last step " +
						"above (147 mod 21=0) relies on to stop.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "a", Min: 2, Max: 1200, Step: 1, Def: 1071},
			{Key: "b", Label: "b", Min: 1, Max: 1200, Step: 1, Def: 462},
		},
		Render: render,
	})
}

// GCD returns the greatest common divisor of a and b via the Euclidean
// algorithm: repeatedly replace (a,b) with (b, a mod b) until b hits 0.
// GCD(a,0)=a for any a≥0, so GCD(0,0)=0.
func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Step is one division the Euclidean algorithm performs while reducing
// gcd(Dividend,Divisor): Dividend = Quotient×Divisor + Remainder.
type Step struct {
	Dividend, Divisor, Quotient, Remainder int
}

// Steps returns every division step the Euclidean algorithm takes computing
// gcd(a,b), in order, stopping once a remainder of 0 is reached. The last
// step's Divisor is the GCD.
func Steps(a, b int) []Step {
	var steps []Step
	for b != 0 {
		q, r := a/b, a%b
		steps = append(steps, Step{a, b, q, r})
		a, b = b, r
	}
	return steps
}

// Square is one square tile placed by the Euclidean algorithm's geometric
// interpretation: tiling an a×b rectangle with squares the size of the
// smaller side leaves a strictly smaller rectangle (the remainder) to
// recurse into, alternating which side is tiled each time. X,Y is the
// square's top-left corner in the same units as a,b; Step is which
// division step (matching Steps above, 0-based) placed it.
type Square struct {
	X, Y, Side float64
	Step       int
}

// TileSquares returns, in placement order, every square the rectangle-
// tiling interpretation of the Euclidean algorithm places while reducing
// gcd(a,b) -- including the final gcd×gcd square. The squares exactly tile
// the original a×b rectangle with no gaps or overlaps: summing Side*Side
// over the result always equals a*b.
func TileSquares(a, b int) []Square {
	if a <= 0 || b <= 0 {
		return nil
	}
	var squares []Square
	x, y := 0.0, 0.0
	w, h := float64(a), float64(b)
	step := 0
	for w > 1e-9 && h > 1e-9 {
		if w >= h {
			side, q := h, int(w/h)
			for k := 0; k < q; k++ {
				squares = append(squares, Square{x + float64(k)*side, y, side, step})
			}
			w -= float64(q) * side
			x += float64(q) * side
		} else {
			side, q := w, int(h/w)
			for k := 0; k < q; k++ {
				squares = append(squares, Square{x, y + float64(k)*side, side, step})
			}
			h -= float64(q) * side
			y += float64(q) * side
		}
		step++
	}
	return squares
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 320, -1, 1, -1, 1).String()
}
