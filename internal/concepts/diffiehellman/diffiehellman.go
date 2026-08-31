// Package diffiehellman visualizes Diffie-Hellman key exchange: two parties
// who agree on nothing secret in advance -- only a public prime p and a
// public base g -- each combine their own private number with the public
// values to land on the exact same shared secret, without either private
// number ever crossing the public channel. The running example is the
// classic small-number walkthrough (p=23, g=5), so every value in the
// pipeline can be checked by hand.
package diffiehellman

import (
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

func render(p map[string]float64) string {
	c := viz.New(760, 480, 0, 1, 0, 1)
	return c.String()
}
