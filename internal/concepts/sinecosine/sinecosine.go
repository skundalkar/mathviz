// Package sinecosine visualizes sine and cosine as the y- and x-coordinates
// of a point going around the unit circle, and "unrolls" that circular
// motion into the familiar wave shape by plotting the same height against
// how far the point has swept around (the angle), instead of against its
// sideways position.
package sinecosine

import (
	"fmt"
	"math"

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

// CirclePoint returns the (x, y) position of a point that has swept
// through angle theta (radians) counterclockwise from 3-o'clock around a
// circle of the given radius centered at the origin.
func CirclePoint(theta, radius float64) (x, y float64) {
	return radius * math.Cos(theta), radius * math.Sin(theta)
}

// SineWave is the "unrolled" height of CirclePoint at angle theta: the
// y-coordinate, plotted against theta itself instead of against x.
func SineWave(theta, radius float64) float64 {
	return radius * math.Sin(theta)
}

// CosineWave is the "unrolled" sideways position of CirclePoint at angle
// theta: the x-coordinate, plotted against theta itself.
func CosineWave(theta, radius float64) float64 {
	return radius * math.Cos(theta)
}

// Revolutions reports how many full turns around the circle theta
// represents (theta / 2π).
func Revolutions(theta float64) float64 {
	return theta / (2 * math.Pi)
}

// waveMax is the right edge of the wave panel: two full revolutions.
const waveMax = 4 * math.Pi

func render(p map[string]float64) string {
	theta, radius := p["theta"], p["radius"]

	c := viz.New(760, 380, -5, waveMax+0.3, -1.75, 1.75)

	// Zero baseline across the whole canvas, and a divider at x=0 between
	// the circle inset (negative x, otherwise unused space) and the wave
	// panel (x in [0, waveMax]).
	c.Path([][2]float64{{c.XMin, 0}, {c.XMax, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, c.YMin}, {0, c.YMax}}, viz.Muted, 1)

	// --- Circle inset ---------------------------------------------------
	// Convert a desired pixel radius into data-space radii that compensate
	// for the canvas's (very different) x and y scales, so the circle
	// renders as an actual circle instead of a squashed ellipse.
	pxPerUnitX := (c.W - c.PadL - c.PadR) / (c.XMax - c.XMin)
	pxPerUnitY := (c.H - c.PadT - c.PadB) / (c.YMax - c.YMin)
	const insetPixelRadius = 55
	rx := insetPixelRadius / pxPerUnitX * radius
	ry := insetPixelRadius / pxPerUnitY * radius
	const cx, cy = -2.6, 0.0

	circle := make([][2]float64, 0, 101)
	for i := 0; i <= 100; i++ {
		phi := 2 * math.Pi * float64(i) / 100
		circle = append(circle, [2]float64{cx + rx*math.Cos(phi), cy + ry*math.Sin(phi)})
	}
	c.Path(circle, viz.Ink, 1.5)
	c.Path([][2]float64{{cx - rx*1.15, cy}, {cx + rx*1.15, cy}}, viz.Muted, 1)
	c.Path([][2]float64{{cx, cy - ry*1.15}, {cx, cy + ry*1.15}}, viz.Muted, 1)

	dotX, dotY := cx+rx*math.Cos(theta), cy+ry*math.Sin(theta)
	c.Path([][2]float64{{cx, cy}, {dotX, dotY}}, viz.Accent, 2)
	c.Path([][2]float64{{dotX, dotY}, {cx, dotY}}, viz.Warm, 1) // height guide
	px, py := c.X(dotX), c.Y(dotY)
	c.Rect(px-4, py-4, 8, 8, viz.Accent, 1)

	// --- Wave panel -------------------------------------------------------
	for k := 0.0; k <= 4; k++ {
		label := "0"
		switch k {
		case 1:
			label = "π"
		case 2:
			label = "2π"
		case 3:
			label = "3π"
		case 4:
			label = "4π"
		}
		c.Tick(k*math.Pi, label)
	}

	sine := viz.Sample(0, waveMax, 300, func(t float64) float64 { return SineWave(t, radius) })
	cosine := viz.Sample(0, waveMax, 300, func(t float64) float64 { return CosineWave(t, radius) })
	c.Path(sine, viz.Accent, 2)
	c.Path(cosine, viz.Warm, 2)

	sinNow, cosNow := SineWave(theta, radius), CosineWave(theta, radius)
	c.VLine(theta, viz.Muted, true)
	c.Path([][2]float64{{0, sinNow}, {theta, sinNow}}, viz.Accent, 1)
	c.Path([][2]float64{{0, cosNow}, {theta, cosNow}}, viz.Warm, 1)
	sx, sy := c.X(theta), c.Y(sinNow)
	c.Rect(sx-4, sy-4, 8, 8, viz.Accent, 1)
	cxPx, cyPx := c.X(theta), c.Y(cosNow)
	c.Rect(cxPx-4, cyPx-4, 8, 8, viz.Warm, 1)

	c.Text(20, 24, fmt.Sprintf("θ = %.2f rad (%.2f revolutions)    radius = %.2f", theta, Revolutions(theta), radius),
		14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("sin θ = %.3f (blue)    cos θ = %.3f (orange)", sinNow, cosNow), 13, viz.Muted, "start")
	c.Text(20, 62, "left: the circle the angle sweeps around    right: that same height/position, unrolled against θ",
		12, viz.Muted, "start")

	return c.String()
}
