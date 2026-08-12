// Package complexnumbers visualizes complex-number multiplication as
// rotation: multiplying a point z = a+bi by a "rotor" w = r(cos θ + i sin
// θ) turns z by θ and scales it by r, all in one coupled step instead of
// nudging its x and y coordinates independently.
package complexnumbers

import (
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "complex-numbers",
		Seq:   26,
		Title: "Complex numbers (rotation in the plane)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"A game character (or a robot arm, or a radar sweep) needs to turn: rotate " +
						"its heading by some angle, then rotate again, and end up facing exactly " +
						"the sum of those angles — with its speed (or arm length, or beam range) " +
						"completely unchanged, no drift. Gut instinct: since rotating 'moves' a " +
						"point, just nudge its x- and y-coordinates by some amount each frame — add " +
						"a bit to x, add a bit to y — and it'll end up facing the new way. That " +
						"instinct is wrong: nudging x and y independently also changes the point's " +
						"distance from the origin, so after a few turns the character's speed " +
						"silently grows or shrinks even though you only meant to change which way " +
						"it's facing. Turning a point cleanly — changing its direction without " +
						"changing its length — needs x and y to move together in one coupled step, " +
						"not independently.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Write a 2D point (a, b) as a single number a + bi, where i is just a label " +
						"for 'the second coordinate' with one extra rule: i×i = −1. Multiplying two " +
						"of these numbers together turns out to be exactly the coupled x/y update " +
						"rotation needs: (a+bi)(c+di) = (ac − bd) + (ad + bc)i. Try it on the " +
						"simplest case, multiplying by i itself (i = 0+1i):",
					"• z = 3+4i, the point (3,4). Its length is |z| = √(3²+4²) = 5.",
					"• z × i = (3+4i)(0+1i) = (3·0 − 4·1) + (3·1 + 4·0)i = −4 + 3i, the point " +
						"(−4, 3).",
					"• Check the length: √((−4)²+3²) = √25 = 5 — unchanged. Check the angle: " +
						"(3,4) sits at arctan(4/3) ≈ 53.13°; (−4,3) sits at ≈143.13° — exactly " +
						"53.13°+90° further round. Multiplying by i is a clean quarter-turn.",
					"That's not a coincidence of i specifically — for any two complex numbers, " +
						"multiplying them multiplies their lengths and adds their angles:",
					"• z = 3+4i (length 5, angle 53.13°) and w = 4+3i (length √(4²+3²) = 5, " +
						"angle arctan(3/4) ≈ 36.87°).",
					"• z × w = (3+4i)(4+3i) = (3·4 − 4·3) + (3·3 + 4·4)i = 0 + 25i, the point " +
						"(0, 25).",
					"• Length: 5 × 5 = 25 ✓. Angle: 53.13° + 36.87° = 90.00° ✓ — and (0,25) " +
						"does sit exactly on the positive imaginary axis, straight up, 90° round " +
						"from the start.",
					"So a rotor w = r·(cos θ + i sin θ) — length r, angle θ — rotates whatever it " +
						"multiplies by exactly θ and scales it by exactly r; when r = 1 the length " +
						"never changes at all, only the direction.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The blue arrow is z, the point being rotated. The orange arrow is the rotor " +
						"w = r·(cos θ + i sin θ) built from the angle and scale sliders — its own " +
						"direction is the angle everything else gets turned by. The green arrow is " +
						"the product z×w: drag the angle slider and it sweeps around exactly like " +
						"z, but always offset from z by the rotor's own angle, tracing the muted " +
						"arc between them — that arc is 'the angle just added.' The dashed circle " +
						"has radius |z|; with scale r=1 the green arrow's tip never leaves that " +
						"circle no matter what angle you dial in — proof that a pure rotation " +
						"changes direction only. Drag the scale slider away from 1 and the tip " +
						"moves off the circle, growing or shrinking with r.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Compose any sequence of 2D rotations (and optional resizes) by just adding " +
						"their angles and multiplying their sizes — no re-deriving sin/cos formulas " +
						"by hand — and read a rotation straight off a single number's angle and " +
						"length instead of tracking two coupled coordinates separately. This is " +
						"exactly how a graphics program, a game engine, or a robot's control code " +
						"turns 'rotate 15° now, then another 25°' into 'rotate 40°, once' with no " +
						"accumulated drift.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"2D game engines and animation software rotate sprites and shapes with " +
						"exactly this multiplication instead of trigonometric matrices written out " +
						"by hand. Electrical engineers describe alternating current as a rotating " +
						"complex number (a 'phasor') because combining two AC signals is just " +
						"complex addition, and shifting one signal's timing is just complex " +
						"multiplication by a rotor. Audio and signal processing (the Fourier " +
						"transform) represent a sound wave as a sum of rotating complex numbers of " +
						"different speeds. And 'imaginary number' as an everyday phrase for " +
						"something inherently unreal is a bit of a misnomer carried over from " +
						"history — here i is doing very literal, very real work: encoding a 90° " +
						"turn.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'multiplying two complex numbers multiplies their lengths " +
						"and adds their angles' — length and angle are the two numbers that matter, " +
						"not the raw (a,b) coordinates. Not like this: assuming i is 'just' an " +
						"abstract placeholder with no numeric meaning (it's a specific 90° " +
						"rotation, and i×i=−1 is exactly what you'd expect from turning 90° twice " +
						"— you end up facing backward, i.e. multiplied by −1); or assuming any " +
						"coordinate nudge that 'looks like turning' preserves length — only " +
						"multiplication by a fixed-length rotor does that, adding arbitrary amounts " +
						"to x and y generally does not.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "zRe", Label: "z — real part", Min: -4, Max: 4, Step: 0.5, Def: 3},
			{Key: "zIm", Label: "z — imaginary part", Min: -4, Max: 4, Step: 0.5, Def: 4},
			{Key: "angle", Label: "Rotor angle θ", Min: 0, Max: 360, Step: 5, Def: 90, Unit: "°"},
			{Key: "scale", Label: "Rotor length r", Min: 0.25, Max: 2, Step: 0.25, Def: 1},
		},
		Render: render,
	})
}

// Multiply returns the product (re,im) of the complex numbers a+bi and
// c+di: (ac−bd) + (ad+bc)i.
func Multiply(a, b, c, d float64) (re, im float64) {
	return a*c - b*d, a*d + b*c
}

// Modulus returns the length |z| of the complex number a+bi — its
// distance from the origin.
func Modulus(a, b float64) float64 {
	return math.Hypot(a, b)
}

// Argument returns the angle in radians, in (−π, π], that a+bi makes with
// the positive real axis.
func Argument(a, b float64) float64 {
	return math.Atan2(b, a)
}

// FromPolar returns the (re,im) components of the complex number with the
// given modulus r and argument theta (radians).
func FromPolar(r, theta float64) (re, im float64) {
	return r * math.Cos(theta), r * math.Sin(theta)
}

// Rotate multiplies a+bi by the rotor r·(cos θ + i sin θ): it turns
// (a,b) by theta radians and scales it by r in one step.
func Rotate(a, b, theta, r float64) (re, im float64) {
	c, d := FromPolar(r, theta)
	return Multiply(a, b, c, d)
}

func render(p map[string]float64) string {
	c := viz.New(534, 520, -13, 13, -13, 13)
	return c.String()
}
