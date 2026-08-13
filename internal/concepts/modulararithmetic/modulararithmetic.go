// Package modulararithmetic visualizes "clock arithmetic" — numbers that
// wrap back around to 0 once they hit a modulus n, instead of growing
// forever — as chords drawn on a circle of n points: one chord from each
// point k to (k × multiplier) mod n, all redrawn at once as the sliders
// move.
package modulararithmetic

import (
	"fmt"
	"math"
	"strings"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "modular-arithmetic",
		Seq:   28,
		Title: "Modular arithmetic (the clock that makes cryptography work)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"It's 9 o'clock, and someone says 'let's meet in 5 hours.' Add it up the " +
						"obvious way and you get 9+5=14 — but no clock face shows 14, so without " +
						"even thinking about it you say '2 o'clock' instead. That trick (numbers " +
						"wrapping back around once they hit some limit, instead of growing forever) " +
						"feels like a harmless clock quirk. But the same wraparound turns out to be " +
						"the load-bearing idea behind keeping a secret online: ordinary addition " +
						"and multiplication grow without bound and are trivial to undo — given a " +
						"sum and one addend, subtraction instantly hands you the other. Is there a " +
						"kind of arithmetic that closes back on itself the way a clock does, and — " +
						"more importantly — can it be built so that undoing one specific operation " +
						"is hard even though doing it forward is easy?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"'a mod n' means: take a, subtract the largest multiple of n that's still " +
						"≤ a, and whatever's left over is the answer — always landing somewhere in " +
						"[0, n). Check the clock scenario from section 1: 14 mod 12 — the largest " +
						"multiple of 12 not exceeding 14 is 12 itself, 14−12=2, so '2 o'clock' was " +
						"exactly 14 mod 12 all along, not a special clock-only rule. The same idea " +
						"works for multiplication, not just addition — walk the times table for " +
						"n=12, multiplier=5:",
					"• 0 × 5 = 0 → 0 mod 12 = 0",
					"• 1 × 5 = 5 → 5 mod 12 = 5",
					"• 2 × 5 = 10 → 10 mod 12 = 10",
					"• 3 × 5 = 15 → 15 mod 12 = 3 — past 12 for the first time, so it wraps",
					"• 4 × 5 = 20 → 20 mod 12 = 8",
					"• 5 × 5 = 25 → 25 mod 12 = 1",
					"Keep going through k=0..11 and every remainder still lands somewhere in " +
						"0-11 no matter how large k×5 gets — the wraparound never lets the result " +
						"escape the clock face. Draw a straight line from each number k to its " +
						"remainder k×5 mod 12, for every k at once, and that set of lines is " +
						"exactly the picture on this page.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"n points evenly spaced around a circle, labeled 0 through n−1 like a clock " +
						"face that starts counting at 0 instead of 12. For every point k, a line is " +
						"drawn from k to (k × multiplier) mod n — the same computation walked in " +
						"section 2, now run for every point at once instead of just the first six. " +
						"Drag the multiplier and every line redraws simultaneously into a " +
						"completely different pattern — multiplier=1 leaves every point connected " +
						"only to itself (no visible lines at all), while other settings sweep out " +
						"stars, loops, and braided patterns as the wraparound sends different " +
						"points to different neighbors. Drag n and the number of points on the " +
						"circle — and therefore how far each one can be sent before it wraps — " +
						"changes with it.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"The same wraparound scales up to modular exponentiation — repeatedly " +
						"multiplying a number by itself, mod n — and that operation is the actual " +
						"trapdoor cryptography leans on. Worked example: 3^4 mod 7. Instead of " +
						"computing 3^4=81 and then reducing, square as you go and reduce at each " +
						"step: 3^2=9, 9 mod 7=2; then 2^2=4, so 3^4 mod 7 = 4 — a few " +
						"multiplications total, however large the exponent, because each " +
						"intermediate result gets pulled back into [0,7) before it's squared " +
						"again. That's cheap to compute forward for enormous exponents. Going " +
						"backward — given that the answer is 4, and the base was 3 and the modulus " +
						"was 7, recover the exponent — is called the discrete logarithm problem, " +
						"and for a large enough modulus no one has found a way to do it that's " +
						"anywhere near as fast as the forward direction. That gap between 'easy " +
						"forward, hard backward' is exactly what Diffie–Hellman key exchange and " +
						"RSA are built on.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Clocks and calendars (day-of-week arithmetic: today's weekday plus some " +
						"number of days, mod 7); hash tables, which drop a key into bucket (key " +
						"mod table size) so every key has a fast, predictable home; check digits " +
						"on ISBNs and credit card numbers (the Luhn algorithm), which catch typos " +
						"by verifying a mod-10 checksum; music theory, where notes an octave apart " +
						"share the same 'pitch class' because pitches are compared mod 12; and " +
						"cryptography — RSA and Diffie–Hellman key exchange both run on modular " +
						"exponentiation, as in section 4.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'mod n means the remainder after dividing by n, always " +
						"landing in [0, n) — 14 mod 12 is 2, not 14, and not 1.16 (that's 14÷12, a " +
						"completely different operation).' Not like this: treating mod as rounding " +
						"or truncating division (14 mod 12 is not 'about 1'), or assuming negative " +
						"numbers go negative — they still wrap into [0, n): −1 mod 12 is 11, the " +
						"hour right before midnight, not −1.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Clock size n (modulus)", Min: 3, Max: 30, Step: 1, Def: 12},
			{Key: "multiplier", Label: "Multiplier a", Min: 0, Max: 29, Step: 1, Def: 5},
		},
		Render: render,
	})
}

// Mod returns a mod n, always in [0, n) — including for negative a, unlike
// Go's built-in % which can return a negative result (e.g. -1 % 12 == -1 in
// Go, but Mod(-1, 12) == 11, the hour right before midnight). n must be
// positive.
func Mod(a, n int) int {
	r := a % n
	if r < 0 {
		r += n
	}
	return r
}

// AddMod returns (a + b) mod n.
func AddMod(a, b, n int) int {
	return Mod(a+b, n)
}

// MulMod returns (a * b) mod n.
func MulMod(a, b, n int) int {
	return Mod(a*b, n)
}

// PowMod returns (base ^ exp) mod n, for exp >= 0, using repeated squaring
// so each intermediate result is reduced mod n before it's squared again —
// the same trick that makes modular exponentiation cheap even for huge
// exponents (see the concept's "What can you do now" section).
func PowMod(base, exp, n int) int {
	if n == 1 {
		return 0
	}
	result := 1
	base = Mod(base, n)
	for exp > 0 {
		if exp&1 == 1 {
			result = MulMod(result, base, n)
		}
		base = MulMod(base, base, n)
		exp >>= 1
	}
	return result
}

// clockPoint returns the (x, y) position of point k on a unit circle of n
// evenly-spaced points, starting at the top (12-o'clock) and going
// clockwise like a real clock face — the same layout described in the
// "What does the picture show?" section.
func clockPoint(k, n int) (x, y float64) {
	theta := math.Pi/2 - 2*math.Pi*float64(k)/float64(n)
	return math.Cos(theta), math.Sin(theta)
}

func render(p map[string]float64) string {
	n := int(math.Round(p["n"]))
	mult := int(math.Round(p["multiplier"]))

	// 534x604 with PadT/PadB overridden below keeps the plot area an exact
	// 468x468 square (534-48-18 = 604-96-40 = 468) so a unit circle drawn
	// through data-space (x,y) renders as an actual circle, not an ellipse
	// — the same trick used in sinecosine's inset circle.
	c := viz.New(534, 604, -1.3, 1.3, -1.3, 1.3)
	c.PadT = 96 // room for three lines of header text above the circle
	c.PadB = 40 // room for one line of footer text below it

	// The clock face itself.
	circle := make([][2]float64, 0, 121)
	for i := 0; i <= 120; i++ {
		phi := 2 * math.Pi * float64(i) / 120
		circle = append(circle, [2]float64{math.Cos(phi), math.Sin(phi)})
	}
	c.Path(circle, viz.Faint, 1.5)

	// One chord per point k, from k to (k*mult) mod n — skipped when a
	// point maps to itself (so multiplier=1 draws no chords at all), and
	// deduplicated so a mutual pair (k maps to j, and j maps back to k)
	// only draws its shared segment once instead of stacking two
	// identical lines on top of each other.
	drawn := map[[2]int]bool{}
	chords := 0
	for k := 0; k < n; k++ {
		target := MulMod(k, mult, n)
		if target == k {
			continue
		}
		key := [2]int{k, target}
		if key[0] > key[1] {
			key[0], key[1] = key[1], key[0]
		}
		if drawn[key] {
			continue
		}
		drawn[key] = true

		x0, y0 := clockPoint(k, n)
		x1, y1 := clockPoint(target, n)
		c.Path([][2]float64{{x0, y0}, {x1, y1}}, viz.Accent, 1)
		chords++
	}

	// Points and their labels, drawn last so they sit on top of the chords.
	for k := 0; k < n; k++ {
		x, y := clockPoint(k, n)
		px, py := c.X(x), c.Y(y)
		c.Rect(px-3, py-3, 6, 6, viz.Ink, 1)

		lx, ly := clockPoint(k, n)
		lpx, lpy := c.X(lx*1.16), c.Y(ly*1.16)
		c.Text(lpx, lpy+4, fmt.Sprintf("%d", k), 12, viz.Muted, "middle")
	}

	// A few worked mappings, matching the section-2 walkthrough, computed
	// for whatever n and multiplier are currently dialed in.
	shown := n
	if shown > 6 {
		shown = 6
	}
	pairs := make([]string, 0, shown)
	for k := 0; k < shown; k++ {
		pairs = append(pairs, fmt.Sprintf("%d→%d", k, MulMod(k, mult, n)))
	}
	trail := ""
	if n > shown {
		trail = " ..."
	}

	c.Text(20, 26, fmt.Sprintf("Clock size n=%d, multiplier a=%d — chord from k to (k × %d) mod %d", n, mult, mult, n),
		15, viz.Ink, "start")
	c.Text(20, 48, fmt.Sprintf("k → k×%d mod %d:  %s%s", mult, n, strings.Join(pairs, "  "), trail),
		13, viz.Muted, "start")
	c.Text(20, 70, fmt.Sprintf("%d chords drawn out of %d points", chords, n), 13, viz.Muted, "start")

	c.Text(20, c.H-14, "multiplier=1 draws no chords at all — every point maps only to itself",
		12, viz.Muted, "start")

	return c.String()
}
