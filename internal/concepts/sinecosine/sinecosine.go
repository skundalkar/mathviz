// Package sinecosine visualizes sine and cosine as the y- and x-coordinates
// of a point going around the unit circle, and "unrolls" that circular
// motion into the familiar wave shape by plotting the same height against
// how far the point has swept around (the angle), instead of against its
// sideways position.
package sinecosine

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "sine-cosine",
		Seq:   24,
		Title: "Sine & cosine (the unit circle unrolled)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A Ferris wheel car starts at 3-o'clock and turns steadily counterclockwise " +
						"around its axle. After it's swept through some angle θ, how high off the " +
						"axle's own height is the car? Gut instinct: the wheel turns at a constant " +
						"rate, so height should climb at a constant rate too — height ∝ θ. That " +
						"instinct is wrong: near the very start and near the very top, the car's " +
						"height barely changes as θ ticks forward; almost all of its rise happens " +
						"through the middle of the swing. A single angle measurement doesn't come " +
						"with a height attached in any obvious way — you need a rule that converts " +
						"'how far around' into 'how high,' and the rule can't be a straight line.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Put the wheel's axle at the center of a circle of radius r. At angle θ, " +
						"measured counterclockwise from 3-o'clock, the car sits at (r·cos θ, r·sin " +
						"θ) — cos θ is how far right/left of center, sin θ is how far up/down. Track " +
						"just the height (sin θ, with r=1) through the first quarter turn:",
					"• θ=0° (the start): height = sin(0°) = 0 — level with the axle.",
					"• θ=10°: height = sin(10°) ≈ 0.174 — barely risen, even though the wheel has " +
						"already turned 10° of its 90° trip to the top.",
					"• θ=60°: height = sin(60°) ≈ 0.866 — two-thirds of the way through the angle, " +
						"but already 87% of the way to the top.",
					"• θ=90° (top of the quarter turn): height = sin(90°) = 1 — the full radius.",
					"The rise from 0° to 10° was only 0.174, but the rise from 50° to 60° " +
						"(sin(60°)−sin(50°) ≈ 0.866−0.766 = 0.10) is nearly as large over the same " +
						"10° step — height changes fastest through the middle of the swing, exactly " +
						"the non-constant pattern the gut instinct missed. sin θ is defined as " +
						"exactly this y-coordinate at every angle θ; cos θ is the matching " +
						"x-coordinate. Past θ=360° (2π radians) the point just goes around again, so " +
						"both repeat forever — that's why they're called periodic.",
					"'Unrolling' means plotting that same height, not against the car's sideways " +
						"position, but against how far around it has swept so far (θ, left to right) " +
						"— like unrolling a paper tape that was wrapped around the wheel. The result " +
						"is the sine wave.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The small circle on the left is the wheel itself, with a dot at the current " +
						"angle θ and a line from the center out to it. The wide panel on the right " +
						"plots that same dot's height (blue, sin θ) and sideways position (orange, " +
						"cos θ) continuously as θ sweeps from 0 all the way to two full turns (4π). " +
						"Drag θ and the marker slides along both waves in lockstep with the dot on " +
						"the circle; drag the radius slider and both the circle and the wave heights " +
						"grow or shrink together — the wave's amplitude is exactly the circle's " +
						"radius.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Predict the exact height or sideways position of anything moving in a circle " +
						"or swinging back and forth, at any angle or any moment, by reading sin θ or " +
						"cos θ off a formula — without needing to draw the circle and measure it " +
						"each time.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A Ferris wheel car's height, a swinging pendulum's sideways displacement, the " +
						"alternating voltage in household AC electricity, a speaker cone's position " +
						"as it produces sound, and the length of daylight through the year all trace " +
						"out this same wave shape. Calling something a 'sine wave' informally — a " +
						"smoothly oscillating stock price, a heart-rate signal — is a real reference " +
						"to this exact shape, not just a figure of speech.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'height rises fastest near the middle of the swing (θ near " +
						"90°) and barely changes near the top and bottom' — the corrected version of " +
						"the gut instinct from section 1. Not like this: assuming a steadily turning " +
						"wheel produces a steadily rising height (height ∝ θ), or mixing up degrees " +
						"and radians — sin(30°) = 0.5, but sin(30 radians) ≈ −0.988, a completely " +
						"different number, so always check which unit a formula or function expects " +
						"before plugging in.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "theta", Label: "Angle θ", Min: 0, Max: 12.566, Step: 0.05, Def: 1.047, Unit: "rad"},
			{Key: "radius", Label: "Circle radius", Min: 0.4, Max: 1.6, Step: 0.05, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(760, 380, 0, 1, -1, 1).Axes().String()
}
