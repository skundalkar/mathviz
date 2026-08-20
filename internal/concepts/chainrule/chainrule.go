// Package chainrule visualizes the chain rule: when a composed function
// f(x) = g(h(x)) is built from a simple inner function h and a simple outer
// function g, its derivative is the product of each part's own derivative,
// f'(x) = g'(h(x)) · h'(x). The running example is a classic related-rates
// problem — a sphere's volume growing as its radius grows at a constant
// rate — where V'(t) = V'(r(t)) · r'(t).
package chainrule

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "chain-rule",
		Seq:   53,
		Title: "Chain rule (multiplying rates through a composition)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"`derivative` showed how to find the slope of a curve like p(t)=t² " +
						"directly. But some quantities aren't handed to you as one formula in t " +
						"— they're built from two simpler, separately-known rates chained " +
						"together. A spherical raindrop's radius grows at a known, constant rate " +
						"as it falls, and geometry separately tells you how a sphere's volume " +
						"responds to a change in radius (dV/dr = 4πr²) — but neither of those " +
						"numbers, on its own, answers 'how fast is the volume growing right now, " +
						"in cm³ per second?' Is there a way to combine two known rates — how fast " +
						"r changes with t, and how fast V changes with r — into the one rate you " +
						"actually want, how fast V changes with t, without first grinding out an " +
						"explicit formula for V(t) and differentiating that from scratch?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Split the problem into two simple, separately-known layers: the inner " +
						"function r(t) = rate·t (radius grows at a constant rate) and the outer " +
						"function V(r) = (4/3)πr³ (a sphere's volume from its radius). Say rate = " +
						"0.5 cm/s, and look at the moment t=4s:",
					"• Inner layer: r(4) = 0.5·4 = 2 cm, and its rate of change is constant: " +
						"dr/dt = 0.5 cm/s at every t.",
					"• Outer layer, evaluated at that radius: dV/dr = 4πr² = 4π(2)² = 16π ≈ " +
						"50.265 — how fast volume responds to radius, right at r=2cm.",
					"• Chain the two rates together by multiplying them: dV/dt = dV/dr · dr/dt " +
						"= 50.265 × 0.5 ≈ 25.133 cm³/s. That's the volume's growth rate in time, " +
						"built entirely from two separately-known, simpler rates — no need to " +
						"expand V(t) = (4/3)π(0.5t)³ into one formula and differentiate that " +
						"directly (though doing so and checking gives the same 25.133 cm³/s).",
					"In general, for f(x) = g(h(x)): f'(x) = g'(h(x)) · h'(x) — take the outer " +
						"function's derivative, evaluate it at the inner function's value (not at " +
						"x directly), and multiply by the inner function's own derivative.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"rate sets how fast the radius grows (dr/dt); t0 picks the moment in time. " +
						"The black curve is the composed function V(t); the orange line is its " +
						"tangent at t0, with slope V'(t0). The readout reports the same three " +
						"numbers as the worked example — r(t0), dV/dr at that radius, dr/dt — and " +
						"their product, which always exactly matches the tangent line's slope. " +
						"Drag rate and watch dr/dt, and therefore the whole tangent slope, scale " +
						"with it; drag t0 and watch dV/dr change (since the sphere is now a " +
						"different size) while dr/dt stays fixed.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Differentiate a function built out of layers by handling each layer on its " +
						"own — however simple or complicated — and multiplying together each " +
						"layer's local derivative, instead of needing one explicit combined " +
						"formula for the whole thing before you can even start.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Related-rates problems are exactly this: a ripple's radius growing on a " +
						"pond and asking how fast its area grows, or a shadow's length changing " +
						"as the sun's angle changes; converting a rate through a chain of units " +
						"(a car's engine RPM to wheel RPM to road speed, each stage its own " +
						"known ratio) chains rates the same way. Most consequentially, training a " +
						"neural network runs this exact rule, layer by layer, backwards through " +
						"the whole network (\"backpropagation\") to find how the final loss " +
						"depends on a weight buried many layers in — each layer contributes one " +
						"more local derivative multiplied into the chain.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'differentiate the outside, leave the inside alone, then " +
						"multiply by the derivative of the inside.' Not like this: differentiating " +
						"only the outer layer and forgetting the inner derivative factor entirely " +
						"— treating V(r(t)) as if it already were a function of t with no inner " +
						"layer to account for is the single most common mistake with composed " +
						"functions, and it silently drops the dr/dt factor (here, 0.5) that the " +
						"whole point of this rule was to include.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "rate", Label: "Radius growth rate (dr/dt, cm/s)", Min: 0.2, Max: 1.5, Step: 0.05, Def: 0.5},
			{Key: "t0", Label: "Moment t0 (s)", Min: 0.2, Max: 6, Step: 0.1, Def: 4},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(560, 400, 0, 6, 0, 30)
	c.Axes()
	return c.String()
}
