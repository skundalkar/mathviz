// Package birthday visualizes the birthday paradox: how the probability
// that at least two people in a room share a birthday climbs to "more
// likely than not" far sooner than most people's gut instinct expects.
package birthday

import (
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
				Body:    []string{"placeholder"},
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

func render(p map[string]float64) string {
	c := viz.New(680, 420, 0, 100, 0, 1)
	return c.String()
}
