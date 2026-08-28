// Package rsa visualizes RSA public-key encryption: a public key anyone can
// use to lock a message, and a private key -- built from the same two
// numbers using the extended Euclidean algorithm -- that only its holder
// can use to unlock it. The running example is the classic small-number
// walkthrough (p=61, q=53, e=17), so every value in the pipeline can be
// checked by hand.
package rsa

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

// The worked example's key pair, small enough to check by hand. p and q are
// the two primes; N is their product (the public modulus); Phi is Euler's
// totient of N; E is the public exponent. D is derived from E and Phi below
// via the extended Euclidean algorithm, the same way real RSA key
// generation computes it.
const (
	p, q = 61, 53
	N    = p * q             // 3233
	Phi  = (p - 1) * (q - 1) // 3120
	E    = 17
)

// D is the private exponent: E's modular inverse mod Phi. Computed once at
// package init so every render call reuses the same value instead of
// re-running the extended Euclidean algorithm each time.
var D = privateExponent()

func privateExponent() int {
	d, ok := ModInverse(E, Phi)
	if !ok {
		panic("rsa: E and Phi are not coprime -- the worked example's constants are wrong")
	}
	return d
}

// ModPow computes base^exp mod m using square-and-multiply, so it stays
// fast (and stays within normal integer range) even for exponents in the
// thousands -- the naive "multiply base by itself exp times" approach would
// overflow long before reaching a real RSA exponent.
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

// ExtendedGCD returns g=gcd(a,b) along with Bézout coefficients x,y such
// that a*x + b*y = g -- the extended form of `euclidean-algorithm`'s plain
// GCD, which only returns g. RSA key generation needs the coefficients, not
// just the GCD itself: x is exactly a's modular inverse mod b when g=1.
func ExtendedGCD(a, b int) (g, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := ExtendedGCD(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

// ModInverse returns d such that e*d ≡ 1 (mod n), and whether one exists
// (it does exactly when gcd(e,n)=1). This is the extended Euclidean
// algorithm's payoff: RSA's private exponent is nothing more than e's
// modular inverse mod φ(n).
func ModInverse(e, n int) (int, bool) {
	g, x, _ := ExtendedGCD(e, n)
	if g != 1 {
		return 0, false
	}
	d := x % n
	if d < 0 {
		d += n
	}
	return d, true
}

// Encrypt applies the public key (n,e) to message m: c = m^e mod n.
func Encrypt(m, e, n int) int { return ModPow(m, e, n) }

// Decrypt applies a decryption exponent d to ciphertext c: m = c^d mod n.
// Passing the true private exponent recovers the original message;
// anything else recovers an unrelated number.
func Decrypt(c, d, n int) int { return ModPow(c, d, n) }

func init() {
	concept.Register(concept.Concept{
		ID:    "rsa-encryption",
		Seq:   76,
		Title: "RSA encryption (public keys built from modular arithmetic)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`modular-arithmetic` ended by pointing at an operation that's easy to do " +
						"forward but hard to undo as the key to keeping a secret online; " +
						"`euclidean-algorithm` ended by pointing at its extended version's " +
						"modular inverses as exactly what that key needs. But there's a more " +
						"basic problem neither concept solved yet: a shared-secret cipher " +
						"requires both sides to already agree on one secret key -- and if the " +
						"only channel two people have is public, with an eavesdropper reading " +
						"every byte, how do they agree on a secret without ever sending it in the " +
						"clear? Is there a way to let anyone lock a message for you using only " +
						"public information, so that only you -- holding one piece of information " +
						"nobody else has -- can unlock it?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Pick two prime numbers, p=61 and q=53, and multiply them: n=p×q=3233. " +
						"Compute φ(n)=(p−1)(q−1)=60×52=3120 -- the count of numbers below n that " +
						"share no factor with n. Pick a public exponent e that shares no factor " +
						"with φ(n); e=17 works (gcd(17,3120)=1, checked the `euclidean-algorithm` " +
						"way). Now find d, the number that undoes e: e×d ≡ 1 (mod φ(n)). That's " +
						"exactly a modular inverse, found by running the extended Euclidean " +
						"algorithm on e and φ(n) and back-substituting through its remainder " +
						"chain -- it turns out d=2753. The public key is the pair (n,e)=(3233,17); " +
						"the private key is (n,d)=(3233,2753). Note what's missing from the public " +
						"key: p, q, and φ(n) themselves. Finding d from e alone means factoring n " +
						"back into p and q first -- trivial for n=3233 by hand, but the same " +
						"multiplication is one-way in practice once p and q are hundreds of " +
						"digits long.",
					"Encrypt a message m=65 with the public key: c = m^e mod n = 65^17 mod 3233 " +
						"= 2790. Decrypt with the private key: m' = c^d mod n = 2790^2753 mod " +
						"3233 = 65 -- back to the original message. This isn't a coincidence: " +
						"e and d were built so that e×d = 1 + k·φ(n) for some whole number k, and " +
						"a theorem of Euler's says m^φ(n) ≡ 1 (mod n) whenever m and n share no " +
						"factor -- so m^(ed) = m¹⁺ᵏᵠ⁽ⁿ⁾ = m·(m^φ(n))^k ≡ m·1^k = m (mod n). Raising " +
						"to the e-th power and then the d-th power always lands back on m, by " +
						"construction.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"m sets the message being sent (any number from 0 to n−1); dOffset nudges " +
						"the decryption exponent away from the true d. The pipeline shows m " +
						"encrypted into c using the public key (n,e), then c decrypted back using " +
						"exponent d+dOffset. At dOffset=0 the recovered value always matches m " +
						"exactly, shown in green. Nudge dOffset to anything else and the recovered " +
						"value jumps to something with no visible relationship to m, shown in red " +
						"-- modular exponentiation doesn't degrade gracefully the way, say, a " +
						"slightly-off measurement would; one wrong exponent is as good as a " +
						"completely wrong key.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Publish a key that lets anyone encrypt a message for you, without you and " +
						"them ever having met to agree on a shared secret -- exactly the gap " +
						"section 1 opened. Whoever encrypted it can't decrypt it back themselves " +
						"(they never had d), which is also how RSA signs things: encrypt a " +
						"message digest with your private key, and anyone can verify it came from " +
						"you by decrypting it with your public key -- the same one-way pipeline, " +
						"run in the other direction.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Every HTTPS connection starts with the server proving its identity using an " +
						"RSA (or similar) public key, often used to establish a shared symmetric " +
						"key for the rest of the session. PGP/GPG encrypted email and SSH host " +
						"keys use the same public/private pair. Software is often 'signed' with a " +
						"publisher's private key so anyone can verify, using the matching public " +
						"key, that the file wasn't tampered with in transit.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'anyone can encrypt with (n,e), but only whoever computed " +
						"d from p and q can decrypt' -- the asymmetry is the entire point, not a " +
						"detail. Not like this: assuming the public key (n,e) alone gives away d " +
						"-- computing d requires φ(n), which requires knowing n's prime factors p " +
						"and q, not just n itself; that gap (easy to multiply p×q, hard to factor " +
						"n back apart at real key sizes) is the one-way street the whole scheme " +
						"stands on. A second slip: expecting a slightly-wrong decryption key to " +
						"produce a slightly-wrong message -- as the picture shows, it produces a " +
						"completely unrelated number instead.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "m", Label: "m (message, 0..3232)", Min: 0, Max: 3232, Step: 1, Def: 65},
			{Key: "dOffset", Label: "dOffset (0 = correct private key)", Min: 0, Max: 5, Step: 1, Def: 0},
		},
		Render: render,
	})
}

// Layout constants for the encrypt/decrypt pipeline diagram, in pixels.
const (
	boxY, boxH   = 170.0, 90.0
	boxW         = 150.0
	mBoxX, cBoxX = 20.0, 285.0
	outBoxX      = 550.0
	outBoxW      = boxW + 20
)

func render(p map[string]float64) string {
	m := int(p["m"])
	dOffset := int(p["dOffset"])
	if m < 0 {
		m = 0
	}
	if m > N-1 {
		m = N - 1
	}

	cipher := Encrypt(m, E, N)
	dUsed := D + dOffset
	recovered := Decrypt(cipher, dUsed, N)
	correct := recovered == m

	// A canvas whose Path/Axes/Sample side is never used -- every draw call
	// below is Rect/Text, both raw-pixel-space regardless of the canvas's
	// data range, since the diagram is a labeled pipeline, not a plot.
	c := viz.New(720, 400, 0, 1, 0, 1)
	midY := boxY + boxH/2

	c.Text(20, 24, fmt.Sprintf("Public key: (n=%d, e=%d)    Private key: (n=%d, d=%d)   [phi(n)=%d]",
		N, E, N, D, Phi), 14, viz.Ink, "start")
	if dOffset != 0 {
		c.Text(20, 44, fmt.Sprintf("Decrypting with d+dOffset = %d+%d = %d instead of the true d", D, dOffset, dUsed),
			13, viz.Warm, "start")
	} else {
		c.Text(20, 44, "Decrypting with the true private exponent d", 13, viz.Muted, "start")
	}

	// Box 1: the plaintext message.
	c.Rect(mBoxX, boxY, boxW, boxH, viz.Faint, 1)
	c.Text(mBoxX+boxW/2, boxY+30, "message m", 13, viz.Muted, "middle")
	c.Text(mBoxX+boxW/2, boxY+58, fmt.Sprintf("%d", m), 20, viz.Ink, "middle")

	// Arrow + label: encrypt.
	arrow(c, mBoxX+boxW+6, midY, cBoxX-6, viz.Accent)
	c.Text((mBoxX+boxW+cBoxX)/2, boxY-14, "encrypt: c = m^e mod n", 12, viz.Accent, "middle")

	// Box 2: the ciphertext.
	c.Rect(cBoxX, boxY, boxW, boxH, viz.Faint, 1)
	c.Text(cBoxX+boxW/2, boxY+30, "cipher c", 13, viz.Muted, "middle")
	c.Text(cBoxX+boxW/2, boxY+58, fmt.Sprintf("%d", cipher), 20, viz.Ink, "middle")

	// Arrow + label: decrypt.
	arrow(c, cBoxX+boxW+6, midY, outBoxX-6, viz.Warm)
	c.Text((cBoxX+boxW+outBoxX)/2, boxY-14, fmt.Sprintf("decrypt: m' = c^%d mod n", dUsed), 12, viz.Warm, "middle")

	// Box 3: the recovered value, green if it matches m, red if not.
	outColor := viz.Good
	verdict := "matches m -- correct key"
	if !correct {
		outColor = viz.Bad
		verdict = "!= m -- wrong key"
	}
	c.Rect(outBoxX, boxY, outBoxW, boxH, outColor, 0.18)
	c.Text(outBoxX+outBoxW/2, boxY+30, "recovered m'", 13, viz.Muted, "middle")
	c.Text(outBoxX+outBoxW/2, boxY+58, fmt.Sprintf("%d", recovered), 20, outColor, "middle")
	c.Text(outBoxX+outBoxW/2, boxY+boxH+20, verdict, 13, outColor, "middle")

	return c.String()
}

// arrow draws a horizontal shaft from x0 to x1 at pixel height y, capped
// with a right-pointing triangle glyph -- built from Rect and Text (both
// raw pixel space) rather than Path, which stays in the canvas's mapped
// data space and would need the shaft's endpoints expressed there instead.
func arrow(c *viz.Canvas, x0, y, x1 float64, color string) {
	const headW = 14.0
	c.Rect(x0, y-1.5, x1-x0-headW, 3, color, 1)
	c.Text(x1-headW/2, y+5, "▶", 14, color, "middle")
}
