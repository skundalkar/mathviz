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
