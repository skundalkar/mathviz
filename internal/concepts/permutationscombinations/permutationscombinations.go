// Package permutationscombinations visualizes the difference between
// permutations (ordered arrangements) and combinations (unordered groups):
// the same "choose k of n" question can have two very different answers
// depending on whether swapping two chosen items counts as a different
// outcome. `pascals-triangle` already covers the combination count, C(n,k);
// this concept adds the permutation count, P(n,k), and the exact factor —
// k! — that separates them.
package permutationscombinations

import (
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

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, -1, 1, -1, 1).String()
}
