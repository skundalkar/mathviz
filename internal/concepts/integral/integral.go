// Package integral visualizes the definite integral as the limit of a
// Riemann sum: chop an interval into n slabs, treat a changing rate as
// roughly constant across each tiny slab, and add up the resulting little
// rectangles. It's the direct reverse of the derivative concept — where
// that one turned a position into a speed (a slope), this one turns a
// speed back into a distance (an area).
package integral

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "integral",
		Seq:   23,
		Title: "Integral (area under a curve)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"The derivative concept's car had position p(t) = t² meters, and its speed at " +
						"any instant turned out to be v(t) = 2t meters/second. Flip the question " +
						"around: suppose all you're given is the speed reading v(t) = 2t at every " +
						"instant from t=0 to t=2 seconds — how far did the car actually travel? If " +
						"speed were constant this would be trivial (distance = speed × time), but " +
						"here the speed itself keeps changing every instant. Which speed do you even " +
						"multiply by?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Chop the 2-second trip into slabs, pretend speed is constant across each tiny " +
						"slab (using its value at the slab's left edge), multiply, and add up:",
					"• n=1 slab (the whole trip at once): sample speed at t=0, v(0)=0 m/s. Distance " +
						"guess = 0 × 2 = 0 meters. The true distance is 4 meters — using the very " +
						"first, slowest reading for the entire trip is a total miss.",
					"• n=2 slabs, each 1 second wide: sample at t=0 and t=1, v(0)=0, v(1)=2. Guess = " +
						"0×1 + 2×1 = 2 meters. Closer, still half the true 4.",
					"• n=4 slabs, each 0.5s wide: sample at t=0, 0.5, 1, 1.5 → v = 0, 1, 2, 3. Guess " +
						"= 0.5×(0+1+2+3) = 3 meters.",
					"• As n keeps growing the guess keeps climbing toward 4 — in fact for this " +
						"straight-line speed the left-edge guess works out exactly to 4 − 4/n, so " +
						"n=100 gives 3.96, n=10,000 gives 3.9996, and the limit as n→∞ is exactly " +
						"4 meters — matching p(2) − p(0) = 4 − 0 from the derivative concept's " +
						"position function. That's the fundamental link: this 'area under the speed " +
						"curve' (the integral) undoes the derivative and hands back the distance.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"n translucent rectangles sit under the speed line v(t)=2t; drag n up and they " +
						"narrow, hugging the true triangular area tighter as the gap between the " +
						"rectangle tops and the line shrinks. A second slider moves where inside each " +
						"slab the sample is taken, from the left edge (0) to the right edge (1) — drag " +
						"it left and every rectangle sits under the rising line (an underestimate); " +
						"drag it right and every rectangle pokes above it (an overestimate). At 0.5 " +
						"(the midpoint of each slab) the two errors happen to cancel exactly for this " +
						"particular straight-line speed — the sum reads exactly 4 even with n=1.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Total up a continuously changing quantity into one accumulated number — total " +
						"distance from a speed that keeps changing, total cost from a rate that keeps " +
						"changing — instead of being stuck needing a constant rate before you can " +
						"multiply anything.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"An odometer accumulating distance from a speedometer reading that never sits " +
						"still, a bank balance accumulating interest that compounds continuously, a " +
						"rain gauge accumulating total rainfall from a rate that rises and falls with " +
						"the storm, and a company's total revenue accumulating from a sales rate that " +
						"varies hour by hour.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'the integral is the area under the rate curve, and shrinking " +
						"the slabs makes the rectangle approximation converge to the exact value.' " +
						"Not like this: assuming one rough estimate — eyeballing an average rate and " +
						"multiplying by the total time — is just as good as actually summing the " +
						"slabs; that's exactly what a single n=1 rectangle does, and the n=1 case " +
						"above (0 meters instead of 4) shows how badly it can miss the moment the " +
						"rate isn't flat. And the midpoint trick that gets n=1 exactly right here is " +
						"a coincidence of using a straight-line speed for the example, not a general " +
						"rule — for a curving rate the midpoint sample is usually much closer than an " +
						"endpoint, but only shrinking the slabs (large n) guarantees convergence to " +
						"the exact area no matter how the rate curves.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "n", Label: "Number of slabs", Min: 1, Max: 40, Step: 1, Def: 4},
			{Key: "samplePos", Label: "Sample point in slab (0=left, 1=right)", Min: 0, Max: 1, Step: 0.05, Def: 0},
		},
		Render: render,
	})
}

// V is the speed function this concept integrates, v(t) = 2t — the exact
// derivative of the derivative concept's position function p(t) = t².
func V(t float64) float64 {
	return 2 * t
}

// RiemannSum approximates ∫[a,b] V(t) dt by chopping [a,b] into n
// equal-width slabs and, in each one, sampling V at a fixed fraction
// samplePos of the way across the slab (0 = left edge, 1 = right edge),
// treating that sample as the slab's constant height. n below 1 is
// clamped to 1.
func RiemannSum(a, b float64, n int, samplePos float64) float64 {
	if n < 1 {
		n = 1
	}
	w := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i < n; i++ {
		left := a + float64(i)*w
		sample := left + samplePos*w
		sum += V(sample) * w
	}
	return sum
}

// TrueIntegral is the exact, analytic value of ∫[a,b] V(t) dt. Since V(t) =
// 2t has antiderivative t², this is simply b² − a².
func TrueIntegral(a, b float64) float64 {
	return b*b - a*a
}

const (
	windowA, windowB = 0.0, 2.0 // fixed [a,b] this concept sweeps
)

func render(p map[string]float64) string {
	n := int(math.Round(p["n"]))
	if n < 1 {
		n = 1
	}
	samplePos := p["samplePos"]

	c := viz.New(680, 380, windowA, windowB, 0, V(windowB)*1.15)
	c.Axes()
	for x := windowA; x <= windowB+1e-9; x += 0.5 {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	// Draw each slab's rectangle: width = one slab, height = V sampled at
	// samplePos through the slab. Drawn before the curve so the curve's
	// line sits visibly on top of (or above) the rectangle tops.
	w := (windowB - windowA) / float64(n)
	for i := 0; i < n; i++ {
		left := windowA + float64(i)*w
		right := left + w
		height := V(left + samplePos*w)
		leftPx, rightPx := c.X(left), c.X(right)
		topPx, basePx := c.Y(height), c.Y(0)
		c.Rect(leftPx, topPx, rightPx-leftPx, basePx-topPx, viz.Accent, 0.35)
	}

	curve := viz.Sample(windowA, windowB, 100, V)
	c.Path(curve, viz.Ink, 2.5)

	sum := RiemannSum(windowA, windowB, n, samplePos)
	trueVal := TrueIntegral(windowA, windowB)

	c.Text(20, 24, fmt.Sprintf("n = %d slabs    sample = %.2f (0 = left edge, 1 = right edge)", n, samplePos),
		14, viz.Ink, "start")
	c.Text(20, 44, fmt.Sprintf("Riemann sum ≈ %.3f    true area = %.3f    error = %.3f", sum, trueVal, sum-trueVal),
		13, viz.Muted, "start")
	c.Text(20, 62, "x = time t (s), y = speed v(t) = 2t (m/s) — shaded rectangles approximate the distance traveled",
		12, viz.Muted, "start")

	return c.String()
}
