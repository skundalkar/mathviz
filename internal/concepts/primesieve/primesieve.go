// Package primesieve visualizes the Sieve of Eratosthenes as a grid of
// integers that gradually get crossed off as multiples of each prime found
// so far, until everything left standing is guaranteed prime — showing why
// the sweep only ever needs to run up to √N.
package primesieve

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "prime-sieve",
		Seq:   27,
		Title: "Prime sieve (the Sieve of Eratosthenes)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You want a list of every prime number up to 100 — for a puzzle, a " +
						"factoring problem, or just curiosity. Gut instinct: check each number one " +
						"at a time, from scratch — to confirm 91 is prime or not, try dividing it " +
						"by 2, 3, 4, 5, ... up until you either find a divisor or give up. Do that " +
						"for every number from 2 to 100 and you technically get the right answer, " +
						"but you've quietly redone the same work over and over — every single " +
						"multiple of 7 gets independently re-tested against 7, one number at a " +
						"time, instead of that one fact ('7 divides these') ever being reused. Is " +
						"there a way to find every prime up to N in one sweep, without re-deriving " +
						"the same divisibility fact again for every single multiple?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Write out every number from 2 to 30 and cross off multiples of each prime " +
						"you find, in order, instead of testing each number against everything " +
						"below it:",
					"• p=2 (the first uncrossed number, so it's prime): cross off every multiple " +
						"of 2 from 4 up — 4, 6, 8, 10, ..., 30. That's 14 numbers eliminated in one " +
						"pass, using a single fact ('divisible by 2') instead of 14 separate tests.",
					"• p=3 (the next number not yet crossed off): cross off multiples of 3 from " +
						"9 up — 9, 15, 21, 27 are newly crossed (6, 12, 18, 24, 30 were already " +
						"gone, crossed off by 2's pass).",
					"• p=5 (next uncrossed number — 4 was already crossed, so it's skipped): " +
						"cross off multiples of 5 from 25 up — only 25 is newly crossed (10, 15, " +
						"20, 30 were already gone).",
					"• p=7: the next multiple of 7 to check is 7×7=49, which is already past 30 " +
						"— so there's nothing left to cross off, and nothing ever will be. Every " +
						"number still uncrossed (2, 3, 5, 7, 11, 13, 17, 19, 23, 29) is prime.",
					"That stopping point isn't a coincidence: any composite number n has a " +
						"smallest prime factor q, and q×q can never exceed n (if both of n's " +
						"factors were bigger than √n, their product would exceed n). So n always " +
						"gets crossed off during q's own pass, at the latest — once the sweep's " +
						"prime p exceeds √N, every number still standing has no factor left that " +
						"could possibly cross it off, and is confirmed prime.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Every integer from 2 to N, in a grid. A gray cell has been crossed off — " +
						"it's composite, and the sweep has already reached its smallest prime " +
						"factor. A white cell is still an open candidate — the sweep hasn't reached " +
						"a factor that would eliminate it yet, even if it secretly is composite. A " +
						"green cell is a confirmed prime: the swept prime has passed √N, so nothing " +
						"left uncrossed can be anything but prime. The highlighted cell is the " +
						"prime currently being swept; drag it up and watch whole diagonal bands of " +
						"multiples turn gray at once, and watch every remaining white cell flip to " +
						"green the moment the sweep clears √N.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Generate every prime up to N in one sweep instead of testing each number " +
						"individually against everything below it, and know exactly when to stop: " +
						"a single number's primality only ever needs checking against divisors up " +
						"to its own square root, not up to the number itself or even half of it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Cryptographic key generation starts by sieving out small-prime candidates " +
						"cheaply before running slower, heavier-duty primality tests on the huge " +
						"numbers actually used for something like RSA. Hash tables are often sized " +
						"to a prime number of buckets specifically because it spreads keys out more " +
						"evenly and cuts down on collisions. And a 'sieve' is a literal kitchen " +
						"tool — sifting flour lets the fine grains fall through while catching the " +
						"lumps — which is exactly this algorithm's shape: let the composites fall " +
						"through, keep what's left.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'to check whether a number is prime, you only need to try " +
						"dividing it by numbers up to its square root' — 97 only needs testing " +
						"against 2, 3, 5, 7 (up to √97 ≈ 9.8), not against 96 different numbers. Not " +
						"like this: assuming you need to check divisibility all the way up to N " +
						"(or N/2) to be sure, or forgetting that 1 is not prime at all (primes are " +
						"defined as greater than 1) while 2 — despite being the only even one — " +
						"absolutely is.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Sieve up to N", Min: 20, Max: 150, Step: 10, Def: 100},
			{Key: "step", Label: "Swept every prime ≤", Min: 1, Max: 150, Step: 1, Def: 3},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(588, 554, 0, 1, 0, 1)
	return c.String()
}
