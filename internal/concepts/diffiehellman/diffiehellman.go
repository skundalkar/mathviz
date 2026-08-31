// Package diffiehellman visualizes Diffie-Hellman key exchange: two parties
// who agree on nothing secret in advance -- only a public prime p and a
// public base g -- each combine their own private number with the public
// values to land on the exact same shared secret, without either private
// number ever crossing the public channel. The running example is the
// classic small-number walkthrough (p=23, g=5), so every value in the
// pipeline can be checked by hand.
package diffiehellman

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// The worked example's public parameters, small enough to check by hand.
// P is the public prime modulus; G is the public base (a primitive root
// mod P, so its powers cycle through every nonzero residue mod P).
const (
	P = 23
	G = 5
)

// ModPow computes base^exp mod m using square-and-multiply, so it stays
// cheap even for exponents far larger than this worked example uses -- the
// same trick `modular-arithmetic` and `rsa-encryption` both rely on.
func ModPow(base, exp, m int) int {
	if m == 1 {
		return 0
	}
	base = base % m
	if base < 0 {
		base += m
	}
	result := 1
	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % m
		}
		exp >>= 1
		base = (base * base) % m
	}
	return result
}

// PublicValue returns the number a party publishes: g^private mod p. It's
// safe to send in the clear -- recovering `private` from it is the
// discrete-log problem BruteForceDiscreteLog demonstrates below.
func PublicValue(private, g, p int) int {
	return ModPow(g, private, p)
}

// SharedSecret returns what one party computes after receiving the other
// party's public value: otherPublic^myPrivate mod p. Both parties end up
// computing g^(a*b) mod p this way, from opposite directions, without
// either ever learning the other's private number.
func SharedSecret(otherPublic, myPrivate, p int) int {
	return ModPow(otherPublic, myPrivate, p)
}

// BruteForceDiscreteLog searches every exponent in [1, p-1] for one that
// reproduces target = g^x mod p, returning the first match and whether one
// was found. This is the only general way to invert PublicValue known --
// there's no shortcut -- so it stands in for "recover the private number
// from the public one," concretely, at a scale small enough to actually
// run: exhaustive search over p-1 candidates is instant here (p=23) and
// would still be instant at p in the low thousands, but at the sizes real
// Diffie-Hellman uses (p hundreds of digits long) the same loop would
// still be running long after the sun burns out.
func BruteForceDiscreteLog(g, target, p int) (x int, ok bool) {
	for x := 1; x < p; x++ {
		if ModPow(g, x, p) == target {
			return x, true
		}
	}
	return 0, false
}

func init() {
	concept.Register(concept.Concept{
		ID:    "diffie-hellman-key-exchange",
		Seq:   83,
		Title: "Diffie-Hellman key exchange (a shared secret over a public channel)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body:    []string{"placeholder"},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "a", Label: "Alice's private number a", Min: 1, Max: 21, Step: 1, Def: 6},
			{Key: "b", Label: "Bob's private number b", Min: 1, Max: 21, Step: 1, Def: 15},
		},
		Render: render,
	})
}

// Layout constants for the two-row exchange diagram, in pixels. Row 1 is
// each party computing and sending their public value; row 2 is each party
// combining the *other* party's public value with their own private one to
// land on the shared secret.
const (
	boxW          = 210.0
	boxH          = 110.0
	leftX, rightX = 20.0, 530.0
	midX          = 275.0
	row1Y, row2Y  = 70.0, 300.0
)

func render(p map[string]float64) string {
	a := int(p["a"])
	b := int(p["b"])

	A := PublicValue(a, G, P)
	B := PublicValue(b, G, P)
	secretAlice := SharedSecret(B, a, P)
	secretBob := SharedSecret(A, b, P)
	match := secretAlice == secretBob

	guessedA, _ := BruteForceDiscreteLog(G, A, P)

	// A canvas whose Path/Axes/Sample side is never used -- every draw call
	// below is Rect/Text in raw pixel space, since the diagram is a
	// labeled pipeline, not a plot (same approach as `rsa-encryption`).
	c := viz.New(760, 480, 0, 1, 0, 1)

	c.Text(20, 24, fmt.Sprintf("Public, known to everyone (including an eavesdropper): p=%d, g=%d", P, G),
		14, viz.Ink, "start")

	// Row 1: each party turns their private number into a public value and
	// sends it across the channel.
	c.Rect(leftX, row1Y, boxW, boxH, viz.Faint, 1)
	c.Text(leftX+boxW/2, row1Y+26, "Alice (private)", 13, viz.Muted, "middle")
	c.Text(leftX+boxW/2, row1Y+52, fmt.Sprintf("a = %d", a), 16, viz.Accent, "middle")
	c.Text(leftX+boxW/2, row1Y+78, fmt.Sprintf("A = g^a mod p = %d", A), 14, viz.Ink, "middle")

	c.Rect(rightX, row1Y, boxW, boxH, viz.Faint, 1)
	c.Text(rightX+boxW/2, row1Y+26, "Bob (private)", 13, viz.Muted, "middle")
	c.Text(rightX+boxW/2, row1Y+52, fmt.Sprintf("b = %d", b), 16, viz.Warm, "middle")
	c.Text(rightX+boxW/2, row1Y+78, fmt.Sprintf("B = g^b mod p = %d", B), 14, viz.Ink, "middle")

	c.Rect(midX, row1Y, boxW, boxH, viz.Bad, 0.08)
	c.Text(midX+boxW/2, row1Y+20, "Public channel", 13, viz.Ink, "middle")
	c.Text(midX+boxW/2, row1Y+44, fmt.Sprintf("A=%d   B=%d", A, B), 15, viz.Ink, "middle")
	c.Text(midX+boxW/2, row1Y+68, "an eavesdropper sees exactly", 12, viz.Muted, "middle")
	c.Text(midX+boxW/2, row1Y+86, "p, g, A, B -- never a or b", 12, viz.Muted, "middle")

	arrow(c, leftX+boxW+6, row1Y+boxH/2-10, midX-6, viz.Accent)
	c.Text((leftX+boxW+midX)/2, row1Y-10, "sends A", 12, viz.Accent, "middle")
	arrow(c, rightX-6, row1Y+boxH/2+10, midX+boxW+6, viz.Warm)
	c.Text((midX+boxW+rightX)/2, row1Y-10, "sends B", 12, viz.Warm, "middle")

	// Row 2: each party combines the *other* party's public value with
	// their own still-private number.
	c.Rect(leftX, row2Y, boxW, boxH, viz.Faint, 1)
	c.Text(leftX+boxW/2, row2Y+26, "Alice computes", 13, viz.Muted, "middle")
	c.Text(leftX+boxW/2, row2Y+52, "B^a mod p", 14, viz.Ink, "middle")
	c.Text(leftX+boxW/2, row2Y+80, fmt.Sprintf("= %d", secretAlice), 18, viz.Accent, "middle")

	c.Rect(rightX, row2Y, boxW, boxH, viz.Faint, 1)
	c.Text(rightX+boxW/2, row2Y+26, "Bob computes", 13, viz.Muted, "middle")
	c.Text(rightX+boxW/2, row2Y+52, "A^b mod p", 14, viz.Ink, "middle")
	c.Text(rightX+boxW/2, row2Y+80, fmt.Sprintf("= %d", secretBob), 18, viz.Warm, "middle")

	secretColor := viz.Good
	verdict := "MATCH -- shared secret"
	if !match {
		secretColor = viz.Bad
		verdict = "mismatch (shouldn't happen)"
	}
	c.Rect(midX, row2Y, boxW, boxH, secretColor, 0.18)
	c.Text(midX+boxW/2, row2Y+26, "shared secret", 13, viz.Muted, "middle")
	c.Text(midX+boxW/2, row2Y+56, fmt.Sprintf("%d", secretAlice), 26, secretColor, "middle")
	c.Text(midX+boxW/2, row2Y+82, verdict, 13, secretColor, "middle")

	arrow(c, leftX+boxW+6, row2Y+boxH/2-10, midX-6, viz.Accent)
	arrow(c, rightX-6, row2Y+boxH/2+10, midX+boxW+6, viz.Warm)

	c.Text(20, 430, fmt.Sprintf(
		"Eve, watching the channel, would have to brute-force search up to p-1=%d exponents to recover a from A -- found a=%d here in that search.",
		P-1, guessedA), 12, viz.Muted, "start")
	c.Text(20, 452, "At real Diffie-Hellman sizes (p hundreds of digits) that same search is the discrete-log problem, believed to take longer than the age of the universe.",
		12, viz.Muted, "start")

	return c.String()
}

// arrow draws a horizontal shaft between x0 and x1 at pixel height y,
// capped with a triangle glyph pointing toward x1 -- built from Rect and
// Text (both raw pixel space) so it works regardless of which side of the
// pipeline diagram it's connecting.
func arrow(c *viz.Canvas, x0, y, x1 float64, color string) {
	const headW = 14.0
	lo, hi := x0, x1
	glyph := "▶"
	if x1 < x0 {
		lo, hi = x1, x0
		glyph = "◀"
	}
	if x1 >= x0 {
		c.Rect(lo, y-1.5, hi-lo-headW, 3, color, 1)
		c.Text(hi-headW/2, y+5, glyph, 14, color, "middle")
	} else {
		c.Rect(lo+headW, y-1.5, hi-lo-headW, 3, color, 1)
		c.Text(lo+headW/2, y+5, glyph, 14, color, "middle")
	}
}
