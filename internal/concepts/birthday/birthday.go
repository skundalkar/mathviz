// Package birthday visualizes the birthday paradox: how the probability
// that at least two people in a room share a birthday climbs to "more
// likely than not" far sooner than most people's gut instinct expects.
package birthday

import (
	"fmt"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "birthday-paradox",
		Seq:   36,
		Title: "Birthday paradox (collisions sooner than you'd guess)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"You're at a party or in a 30-person classroom, and someone raises the " +
						"question: what are the odds two of us share a birthday? Gut instinct: " +
						"with 365 days in a year and only 30 people, that seems like a long shot — " +
						"maybe something like 30 out of 365, roughly 8%. That instinct is comparing " +
						"the number of people to the number of days, as if the only way to get a " +
						"match is for one specific person you already have in mind — say, you — to " +
						"share a birthday with someone else. But nobody asked about you " +
						"specifically; the question is whether ANY two of the 30 people match, and " +
						"with 30 people there are far more than 30 chances for a pair to collide. " +
						"How many pairs are actually in play, and does counting pairs instead of " +
						"people change the answer?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Count the opposite instead: the probability that everyone's birthday is " +
						"different, then subtract that from 1.",
					"• Person 1 can have any birthday — no constraint yet.",
					"• Person 2 has to dodge person 1's day: 364/365 ≈ 99.73% chance of landing " +
						"on a still-free day.",
					"• Person 3 has to dodge both already-taken days: 363/365 ≈ 99.45%.",
					"Multiply those shrinking fractions together as more people join, and the " +
						"product falls faster than intuition expects, because each new person has " +
						"one more taken day to avoid than the last. Doing that multiplication out " +
						"to real group sizes: n=10 people gives an 11.7% chance of a match; n=20 " +
						"gives 41.1%; n=23 gives 50.7% — the first point where a shared birthday " +
						"becomes more likely than not; n=30 gives 70.6%; n=50 gives 97.0%.",
					"The reason it climbs this fast isn't about matching any one specific day — " +
						"it's about how many PAIRS of people exist to potentially match. With n " +
						"people there are n(n-1)/2 distinct pairs, and at n=23 that's " +
						"23×22/2 = 253 separate pairs, each an independent-ish shot at a 1-in-365 " +
						"coincidence. 253 chances at long odds adds up to roughly even overall odds " +
						"— the naive '30 out of 365' guess counted people, not pairs, and pairs is " +
						"what actually drives this number.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The People in the room slider sets n; Possible birthdays sets how many " +
						"equally-likely days there are to collide on (365 by default, but drag it " +
						"down to model a smaller pool — a lottery draw, a hash space, a smaller " +
						"calendar). The blue curve plots P(at least one shared birthday) against n " +
						"for the current day count; the orange square marks exactly where the " +
						"current n sits on that curve, with a dashed line dropping to the axis. The " +
						"faint horizontal line marks the 50% mark, and the green vertical line " +
						"marks the smallest n that crosses it for the current day count — drag " +
						"Possible birthdays down and watch that green line, and the required n, " +
						"slide left.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Recognize when a 'that seems unlikely' reaction is comparing the wrong " +
						"quantities — anytime the real question is 'does ANY pair among many " +
						"options match' rather than 'does one specific thing match a specific " +
						"other thing,' this math applies, and the true probability runs far higher " +
						"than counting individuals suggests. You can also see the effect scale: " +
						"drop Possible birthdays from 365 to 100 (a smaller pool — say, a raffle " +
						"with 100 possible winning numbers instead of 365 calendar days) and the " +
						"50%-crossing point drops from 23 people to just 13, because n(n-1)/2 pairs " +
						"only has to catch up to a smaller pool to become likely.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"The classic version — any room of 23 or more people is more likely than not " +
						"to have a shared birthday — surprises people at parties and in classrooms " +
						"precisely because of the mistake in section one. The same 'ANY pair, not " +
						"one specific pair' logic shows up wherever many independent draws come " +
						"from a shared pool: website session IDs or short URLs, generated at " +
						"random from a large space, only need on the order of √(pool size) IDs " +
						"generated before two are likely to collide — far fewer than the full pool " +
						"size naive intuition might suggest is 'obviously safe.' Cryptographic hash " +
						"functions are designed with exactly this 'birthday bound' in mind: finding " +
						"any two inputs that hash to the same value (not one specific pre-chosen " +
						"pair) takes on the order of 2^(n/2) attempts for an n-bit hash, not 2^n, " +
						"which is why hash outputs need to be roughly twice as many bits as the " +
						"'obvious' security level suggests.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the question is whether ANY two people match, so the " +
						"number of pairs — n(n-1)/2 — is what should scale against the day count, " +
						"not the number of people, n, directly.' Not like this: comparing n to days " +
						"(like '30 people out of 365 days, so about 8%') as if the chances grow " +
						"only linearly with people — that undercounts by roughly a factor of n/2, " +
						"since pairs, not people, are what's really being compared against days. " +
						"Also not like this: confusing 'a shared birthday with YOU specifically' " +
						"(which really does stay a long shot — only about 6% with 22 other people " +
						"in the room) with 'a shared birthday between ANY two people in the room' " +
						"(the 50.7% headline number at n=23) — these are two different questions " +
						"with very different answers, and the whole 'paradox' lives entirely in " +
						"mixing them up.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "People in the room", Min: 1, Max: 100, Step: 1, Def: 23},
			{Key: "days", Label: "Possible birthdays", Min: 2, Max: 365, Step: 1, Def: 365},
		},
		Render: render,
	})
}

// ProbNoCollision returns the probability that n people, each with a
// birthday drawn independently and uniformly at random from `days`
// possible days, all have distinct birthdays. It's the product of n
// shrinking fractions: the 1st person's birthday is free (days/days), the
// 2nd person must dodge 1 already-taken day ((days-1)/days), the 3rd must
// dodge 2 ((days-2)/days), and so on. n<=1 always has zero chance of a
// collision (1.0); n>days is impossible to keep collision-free at all (the
// pigeonhole principle: you can't fit n people each into their own day
// when there are fewer days than people), so it returns 0 rather than
// running the product into negative factors.
func ProbNoCollision(n, days int) float64 {
	if n <= 1 {
		return 1
	}
	if n > days {
		return 0
	}
	p := 1.0
	for i := 0; i < n; i++ {
		p *= float64(days-i) / float64(days)
	}
	return p
}

// ProbCollision returns the probability that at least two of n people
// (birthdays drawn independently and uniformly from `days` possible days)
// share a birthday: the complement of ProbNoCollision.
func ProbCollision(n, days int) float64 {
	return 1 - ProbNoCollision(n, days)
}

// MinPeopleForProbability returns the smallest n for which
// ProbCollision(n, days) >= threshold, searching n = 1, 2, 3, .... Returns
// days+1 (one past the guaranteed-collision point) if no n up to that
// bound reaches the threshold, which only happens for a threshold above 1.
func MinPeopleForProbability(days int, threshold float64) int {
	for n := 1; n <= days+1; n++ {
		if ProbCollision(n, days) >= threshold {
			return n
		}
	}
	return days + 1
}

// maxN is the fixed right edge of the curve's x-axis, matching the "People
// in the room" slider's own Max so the whole slider range is always
// visible on the plot.
const maxN = 100

func render(p map[string]float64) string {
	n := int(p["n"])
	days := int(p["days"])
	if days < 2 {
		days = 2
	}

	c := viz.New(680, 420, 0, maxN, 0, 1)
	c.Axes()
	for x := 0.0; x <= maxN; x += 20 {
		c.Tick(x, fmt.Sprintf("%.0f", x))
	}

	// The collision-probability curve for the current `days`, over the
	// full n range the slider can reach.
	curve := make([][2]float64, maxN+1)
	for i := 0; i <= maxN; i++ {
		curve[i] = [2]float64{float64(i), ProbCollision(i, days)}
	}
	c.Path(curve, viz.Accent, 2.5)

	// The 50%-chance reference line and the smallest n that crosses it,
	// for this many possible days.
	c.Path([][2]float64{{0, 0.5}, {maxN, 0.5}}, viz.Muted, 1)
	half := MinPeopleForProbability(days, 0.5)
	if half <= maxN {
		c.Path([][2]float64{{float64(half), 0}, {float64(half), 1}}, viz.Good, 1.5)
		c.Text(c.X(float64(half))+4, c.Y(0.06), fmt.Sprintf("n=%d crosses 50%%", half), 12, viz.Good, "start")
	}

	// The current (n, P) point, highlighted.
	prob := ProbCollision(n, days)
	px, py := c.X(float64(n)), c.Y(prob)
	c.Rect(px-4, py-4, 8, 8, viz.Warm, 1)
	c.Path([][2]float64{{float64(n), 0}, {float64(n), prob}}, viz.Warm, 1)

	c.Text(16, 24, fmt.Sprintf("%d possible birthdays, %d people in the room", days, n), 14, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("P(at least one shared birthday) = %.1f%%", prob*100), 15, viz.Warm, "start")
	if n > days {
		c.Text(16, 64, "n exceeds the number of possible days -- a shared birthday is guaranteed (pigeonhole principle).",
			13, viz.Muted, "start")
	} else {
		c.Text(16, 64, fmt.Sprintf("It only takes %d people to make a shared birthday more likely than not.", half),
			13, viz.Muted, "start")
	}

	return c.String()
}
