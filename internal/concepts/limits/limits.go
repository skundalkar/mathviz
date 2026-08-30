// Package limits visualizes the limit of a function as x approaches a
// point a: the number f(x) closes in on, which is a completely separate
// question from what f(a) itself equals -- if anything. The running
// examples both sit at a=1: f(x) = (x^2-1)/(x-1), which is undefined
// exactly at x=1 (0/0) but has a well-defined limit of 2 there, and a step
// function whose limit fails to exist at x=1 because its left- and
// right-hand limits disagree, even though the function itself IS defined
// there. `derivative` and `integral` both quietly lean on this same
// "does f(x) settle down as x approaches something" question -- for a
// secant slope and a running sum, respectively -- without spelling it out;
// this concept spells it out directly for a plain function value.
package limits

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "limits",
		Seq:   80,
		Title: "Limits (what f(x) is closing in on)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Compute (x^2-1)/(x-1) at x=1: plug it in and you get 0/0 -- undefined, a " +
						"calculator error. But try values near x=1 instead: at x=1.1 the formula " +
						"gives (1.21-1)/0.1 = 2.1; at x=1.01 it gives 2.01; at x=0.99 it gives 1.99. " +
						"The formula refuses to answer at exactly x=1, but every value near x=1 is " +
						"crowding tighter and tighter around 2. Is there a number this function " +
						"legitimately 'means' to output at x=1, even though plugging x=1 in " +
						"directly breaks the formula -- and is there a way to tell when a function " +
						"does NOT have an honest answer to sneak up on?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Factor the formula: (x^2-1)/(x-1) = (x-1)(x+1)/(x-1) = x+1, valid for every " +
						"x != 1 (canceling (x-1) is only legal when it isn't zero). So for any x " +
						"close to 1 but not equal to 1, f(x) exactly equals x+1 -- and as x -> 1, " +
						"x+1 -> 2. That number, written lim(x->1) f(x) = 2, is the limit: it exists " +
						"and equals 2, despite f(1) itself being undefined. The graph has a single-" +
						"point hole at (1,2) -- a 'removable' discontinuity -- but is otherwise the " +
						"plain line y=x+1.",
					"Not every function behaves that nicely. Take a step function g: g(x) = -1 " +
						"for x<1, and g(x) = 1 for x>=1. Approach from the left (x=0.9, 0.99, " +
						"0.999, ...) and g(x) is always -1 -- that's the left-hand limit. Approach " +
						"from the right (x=1.1, 1.01, 1.001, ...) and g(x) is always 1 -- the " +
						"right-hand limit. Since the two one-sided limits disagree, the two-sided " +
						"limit lim(x->1) g(x) simply does not exist: no single number captures " +
						"'the value g is sneaking up on,' because it depends entirely on which side " +
						"you approach from. Note that g(1) IS defined here (g(1)=1) -- whether " +
						"f(a) is defined and whether the limit at a exists are two separate " +
						"questions, and this example has one without the other.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"h drags x = 1+h along the curve, right of a=1 for h>0 and left of it for " +
						"h<0. mode switches which function is drawn: mode=0 is (x^2-1)/(x-1), with " +
						"an open ring at (1,2) marking the removable hole -- the point the function " +
						"never actually reaches but every nearby point crowds toward. mode=1 is the " +
						"step function, drawn as two separate horizontal rays with a visible jump " +
						"at x=1. A dashed guide line connects the marked point (1+h, f(1+h)) over to " +
						"the y-axis so you can read off how close it's landing to whichever limit " +
						"is relevant on that side. The readout reports the left-hand limit, the " +
						"right-hand limit, and whether they agree -- in mode=0 they always agree " +
						"(both 2, no matter how far h has to shrink); in mode=1 they never do (-1 " +
						"vs 1), which is exactly why that limit doesn't exist.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Say precisely whether a function has a well-behaved value near a troublesome " +
						"point even when you can't literally evaluate it there (division by zero, " +
						"or the formula is simply undefined at that x) -- and keep that question " +
						"entirely separate from asking whether the function is defined there at " +
						"all. This is the exact tool `derivative` needed to make 'instantaneous " +
						"slope' rigorous (the secant slope as its gap h shrinks to 0) and " +
						"`integral` needed to make 'area as slab count grows' rigorous (a sum as " +
						"the slab width shrinks to 0) -- both are this same 'does f(x) settle down " +
						"as x approaches something' question, just applied to a slope-ratio or a " +
						"running sum instead of a plain function value.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Removable holes: combining two resistors in parallel gives a formula that " +
						"looks like it breaks down as one resistance approaches zero, but the " +
						"circuit's actual behavior in that limit is perfectly well-defined (it just " +
						"shorts to the other resistor's value) -- physics and engineering formulas " +
						"are full of these 'divide by zero on paper, fine in reality' spots. Jump " +
						"discontinuities: a thermostat snaps the heater from off to on the instant " +
						"the temperature crosses a set point, a late fee jumps from $0 to $25 the " +
						"moment a payment is one second late, and a step-function shipping-rate " +
						"table jumps to the next bracket the instant a package crosses a weight " +
						"threshold -- in all three, nudging the input by an infinitesimal amount " +
						"doesn't nudge the output, it teleports it, exactly the failure mode of a " +
						"limit not existing.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the limit is what f(x) is closing in on as x approaches a " +
						"-- a completely separate question from what f(a) itself equals, if " +
						"anything.'",
					"Not like this: assuming a limit exists just because a formula looks like it " +
						"should simplify nicely, without checking both directions -- always check " +
						"that the left-hand and right-hand limits agree before calling a two-sided " +
						"limit's value settled. Also not like this: reading f(a) straight off the " +
						"function and calling that 'the limit' -- the step-function example has " +
						"g(1)=1 defined, yet the limit at 1 doesn't exist at all; the two only have " +
						"to agree when the function is continuous there.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "h", Label: "Distance from a=1 (h, sign = direction)", Min: -1, Max: 1, Step: 0.05, Def: 0.5},
			{Key: "mode", Label: "Function (0 = removable hole, 1 = jump discontinuity)", Min: 0, Max: 1, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 400, -1, 3, -2, 4).String()
}
