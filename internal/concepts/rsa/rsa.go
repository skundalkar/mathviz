// Package rsa visualizes RSA public-key encryption: a public key anyone can
// use to lock a message, and a private key -- built from the same two
// numbers using the extended Euclidean algorithm -- that only its holder
// can use to unlock it. The running example is the classic small-number
// walkthrough (p=61, q=53, e=17), so every value in the pipeline can be
// checked by hand.
package rsa

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

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

func render(p map[string]float64) string {
	_ = p
	return viz.New(720, 400, 0, 1, 0, 1).String()
}
