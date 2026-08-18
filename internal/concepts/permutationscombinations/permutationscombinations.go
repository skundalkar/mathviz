// Package permutationscombinations visualizes the difference between
// permutations (ordered arrangements) and combinations (unordered groups):
// the same "choose k of n" question can have two very different answers
// depending on whether swapping two chosen items counts as a different
// outcome. `pascals-triangle` already covers the combination count, C(n,k);
// this concept adds the permutation count, P(n,k), and the exact factor —
// k! — that separates them.
package permutationscombinations

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "permutations-vs-combinations",
		Seq:   46,
		Title: "Permutations vs. combinations (does order matter?)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You're organizing an 8-racer event. Awarding gold, silver, and bronze needs " +
						"an ordered outcome — coming in 2nd is a different result from coming in " +
						"1st, even for the exact same three racers. Picking a 3-person relay team " +
						"from the same 8 racers doesn't care about order at all; the team is the " +
						"team, however you name the members. `pascals-triangle` already gave you " +
						"'n choose k,' C(n,k), for exactly the second question — but if you reach " +
						"for that same formula for the medal question, are you counting the right " +
						"thing? Does caring about order change the count, and if so, by how much?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Count medal orderings directly: 1st place can be any of the 8 racers, 2nd " +
						"place any of the remaining 7, 3rd place any of the remaining 6. " +
						"8 × 7 × 6 = 336 distinct ordered outcomes — that's the permutation count, " +
						"P(8,3). Now count relay teams, which ignore order entirely: those same 336 " +
						"orderings only describe 336 / 6 = 56 distinct teams, because any specific " +
						"3-person team, say {A, B, C}, can be arranged into 3! = 6 different medal " +
						"orders — ABC, ACB, BAC, BCA, CAB, CBA — and all six are the same team, just " +
						"a different assignment of who's on which step of the podium. 56 is exactly " +
						"`pascals-triangle`'s C(8,3).",
					"• P(8,3) = 8 × 7 × 6 = 336 ordered outcomes.",
					"• Each unordered team of 3 corresponds to 3! = 6 of those orderings.",
					"• C(8,3) = P(8,3) / 3! = 336 / 6 = 56 distinct teams.",
					"The general rule: Permutations = Combinations × k!, because k! is exactly the " +
						"number of ways to reorder the k items you've already chosen among " +
						"themselves — a fact that doesn't depend on n at all, only on how many " +
						"items were chosen.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"For the n set by the top slider, one pair of bars per k from 0 to n: blue " +
						"for P(n,k), orange for C(n,k), both height-scaled by square root so the " +
						"much larger permutation count stays on the same chart as the smaller " +
						"combination count. The k slider highlights one pair with a bold outline and " +
						"reads off both exact numbers plus their ratio, k!. Push k up from 0 and " +
						"watch the gap between the blue and orange bars widen — the ratio between " +
						"them keeps climbing (k! grows fast) even as both individual counts rise and " +
						"then fall back down toward n's edges.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Decide, before counting anything, whether a scenario needs P(n,k) or C(n,k) " +
						"by asking one question: does swapping two of the chosen items produce a " +
						"different outcome? Yes means permutations; no means combinations. And once " +
						"you know which one you need, you can get the other for free — they're " +
						"always exactly k! apart, so you never have to rederive either count from " +
						"scratch once you have one of them.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A lottery draw where six numbered balls are picked and their draw order " +
						"doesn't affect whether you won: combinations. A 4-digit PIN code, where " +
						"1234 and 4321 unlock completely different accounts: permutations (and with " +
						"repeated digits allowed, an even bigger count than P(n,k) alone). Assigning " +
						"a soccer team's starting lineup to specific numbered positions is a " +
						"permutation; simply picking which 11 players make the squad, position " +
						"unassigned, is a combination. Counting anagrams of a word — how many " +
						"distinct letter orderings exist — is a permutation count.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: ask 'does swapping two chosen items change the outcome?' — " +
						"yes means permutations, no means combinations — and Permutations = " +
						"Combinations × k!, never just 'a bit more.'",
					"Not like this: assuming the smaller `pascals-triangle` 'n choose k' number " +
						"always applies, even to a scenario where order genuinely matters — that " +
						"undercounts by a factor of k!, exactly the mistake `pascals-triangle`'s own " +
						"'common mistake' section already flagged with a smaller worked example (5 " +
						"people choose 2: 20 ordered pairs vs. 10 unordered pairs). Also not like " +
						"this: assuming the two counts are 'close' for large k — the gap is " +
						"multiplicative (k!), so it grows explosively: at k=5 they differ by a " +
						"factor of 120, not by some small fixed amount.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Items available (n)", Min: 1, Max: 10, Step: 1, Def: 8},
			{Key: "k", Label: "Highlighted choice size (k)", Min: 0, Max: 10, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// Factorial returns n! for n>=0 as an int64. Negative n returns 0 —
// undefined. n stays small in this package (the n Param caps at 10, so
// 10! = 3,628,800), well within int64 range.
func Factorial(n int) int64 {
	if n < 0 {
		return 0
	}
	f := int64(1)
	for i := 2; i <= n; i++ {
		f *= int64(i)
	}
	return f
}

// Permutations returns P(n,k) — the number of ordered arrangements of k
// items chosen from n distinct items: n! / (n-k)!, computed as a direct
// product n × (n-1) × ... × (n-k+1) so it never divides a huge factorial by
// another huge factorial. Returns 0 for an out-of-range k (k<0 or k>n).
func Permutations(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	p := int64(1)
	for i := 0; i < k; i++ {
		p *= int64(n - i)
	}
	return p
}

// Combinations returns C(n,k), "n choose k" — the number of unordered
// groups of k items chosen from n distinct items: Permutations(n,k) / k!.
// Returns 0 for an out-of-range k (k<0 or k>n). Same value pascals-triangle
// builds by addition; this package derives it from Permutations instead,
// to keep the Permutations-to-Combinations relationship (divide out the k!
// orderings) explicit in the code, not just in prose.
func Combinations(n, k int) int64 {
	if n < 0 || k < 0 || k > n {
		return 0
	}
	return Permutations(n, k) / Factorial(k)
}

func render(params map[string]float64) string {
	n := int(params["n"] + 0.5)
	if n < 1 {
		n = 1
	}
	if n > 10 {
		n = 10
	}
	k := int(params["k"] + 0.5)
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}

	// Bar heights are sqrt-scaled: P(n,k) can run into the millions while
	// C(n,k) stays far smaller, and a plain linear axis would flatten every
	// combination bar to an invisible sliver next to its permutation pair.
	maxSqrt := 0.0
	for i := 0; i <= n; i++ {
		if s := math.Sqrt(float64(Permutations(n, i))); s > maxSqrt {
			maxSqrt = s
		}
	}
	if maxSqrt <= 0 {
		maxSqrt = 1
	}
	yMax := maxSqrt * 1.2

	c := viz.New(680, 420, -0.5, float64(n)+0.5, 0, yMax)
	c.Axes()
	for i := 0; i <= n; i++ {
		c.Tick(float64(i), fmt.Sprintf("%d", i))
	}

	for i := 0; i <= n; i++ {
		if i == k {
			// Backdrop band marking the highlighted k, drawn before the
			// bars so it sits behind them.
			bx0, bx1 := c.X(float64(i)-0.46), c.X(float64(i)+0.46)
			c.Rect(bx0, c.Y(yMax), bx1-bx0, c.Y(0)-c.Y(yMax), viz.Faint, 1)
		}

		p, comb := Permutations(n, i), Combinations(n, i)

		// Permutation bar, left half of this k's slot.
		px0, px1 := c.X(float64(i)-0.42), c.X(float64(i)-0.03)
		py0, py1 := c.Y(0), c.Y(math.Sqrt(float64(p)))
		c.Rect(px0, py1, px1-px0, py0-py1, viz.Accent, 0.85)

		// Combination bar, right half of this k's slot.
		cx0, cx1 := c.X(float64(i)+0.03), c.X(float64(i)+0.42)
		cy0, cy1 := c.Y(0), c.Y(math.Sqrt(float64(comb)))
		c.Rect(cx0, cy1, cx1-cx0, cy0-cy1, viz.Warm, 0.85)

		if i == k {
			c.Text((px0+px1)/2, py1-6, fmt.Sprintf("%d", p), 11, viz.Accent, "middle")
			c.Text((cx0+cx1)/2, cy1-6, fmt.Sprintf("%d", comb), 11, viz.Warm, "middle")
		}
	}

	permK, combK := Permutations(n, k), Combinations(n, k)
	c.Text(16, 24, fmt.Sprintf("n = %d    k = %d    (blue = permutations P(n,k), orange = combinations C(n,k))", n, k),
		14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("P(%d,%d) = %d    C(%d,%d) = %d    ratio = P/C = %d = %d!",
		n, k, permK, n, k, combK, Factorial(k), k), 15, viz.Warm, "start")

	return c.String()
}
