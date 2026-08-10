// Package correlation shows what the Pearson correlation coefficient r
// actually looks like: a scatter cloud that goes from a tight downhill line
// (r = -1) through a shapeless blob (r = 0) to a tight uphill line (r = +1).
// The picture also carries the standard warning label — a strong r only says
// two variables move together, never that one causes the other.
package correlation

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "correlation",
		Seq:   4,
		Title: "Correlation (r)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Every summer, ice cream sales climb — and every summer, drowning deaths " +
						"climb too. Plot one against the other and you'd find a strong positive " +
						"relationship: when one goes up, so does the other, reliably. It's tempting " +
						"to read that as 'ice cream causes drowning' (or the reverse) — but neither " +
						"is true. Both are being pulled along by a third thing entirely: hot weather " +
						"means more swimming and more ice cream at the same time.",
					"That raises two questions a gut reaction to a scatter of points can't answer " +
						"on its own: how tightly do two variables actually move together — is it a " +
						"strong, reliable pattern or just a handful of coincidental points — and, " +
						"separately, does 'moving together' even mean one causes the other? You need " +
						"a number that answers the first question precisely before you can even " +
						"start on the second.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"That number is r, and here's the whole scale in real numbers. Four students, " +
						"hours studied vs. test score:",
					"• 1,2,3,4 hours vs. 60,70,80,90 — each extra hour is worth exactly +10 " +
						"points, and the points sit dead on a straight uphill line: r = +1, perfect " +
						"positive correlation.",
					"• Hours of sleep lost vs. focus score: 0,1,2,3 vs. 90,80,70,60 — same " +
						"perfectly straight line, only downhill: r = -1, perfect negative correlation.",
					"• Shoe size vs. test score: 8,9,10,11 vs. 72,58,81,65 — no pattern at all; " +
						"run the numbers and r comes out to about 0.03, essentially r = 0.",
					"r isn't measuring whether some deep relationship exists — only how tightly " +
						"the points hug a straight line. The ice-cream/drowning cloud from the " +
						"opening sits up near +1 for exactly the same reason the hours-studied cloud " +
						"does: real, tight, predictable co-movement. It just doesn't say why the " +
						"line is there.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"Slide r and watch the scatter cloud tighten into a line or spread into a " +
						"shapeless blob — the same landmarks as the worked examples above: near " +
						"r = +1 it looks like the hours-studied cloud, near r = -1 like the " +
						"sleep-lost cloud, near r = 0 like the shoe-size blob. The orange line is " +
						"the trend r implies; both axes are standardized, so that line always runs " +
						"through the origin with slope r, which is why it visibly steepens as r " +
						"moves toward ±1. Sample size n controls how closely the plotted cloud " +
						"matches the target: with a small n, the sample r printed above the plot " +
						"visibly drifts from the target r, a reminder that a correlation measured " +
						"from one sample is itself just an estimate of the true relationship, not " +
						"the relationship.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now put one precise number on how tightly two variables move " +
						"together instead of eyeballing a scatterplot — and, just as importantly, " +
						"you know exactly what that number does and doesn't tell you: a strong r " +
						"(like the ice-cream/drowning cloud) is a real, useful signal that something " +
						"is going on, but on its own it never tells you what — that still takes more " +
						"information than the correlation itself can supply.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'Cities with more police have more crime' (both driven by population " +
						"density), a stock-picking strategy that correlates with past returns " +
						"purely by chance, or a feature in a machine-learning model that correlates " +
						"with the label without causing it — swap it out in production and the " +
						"correlation, having no causal footing, quietly disappears.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: 'X correlates with Y' is a purely descriptive claim — they " +
						"tend to move together — and stops exactly there; it says nothing about " +
						"which one, if either, is causing the other.",
					"Not like this: 'it's just correlation' is said so reflexively it's become the " +
						"opposite mistake just as often — dismissing a strong, useful signal as " +
						"meaningless because no mechanism has been proven yet, when a strong " +
						"correlation is usually worth investigating, not shrugging off.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "r", Label: "Target correlation (r)", Min: -1, Max: 1, Step: 0.05, Def: 0.6},
			{Key: "n", Label: "Sample size (n)", Min: 20, Max: 300, Step: 10, Def: 150},
		},
		Render: render,
	})
}

// hash01 turns an integer seed into a deterministic, evenly-scattered value in
// [0, 1). It is a fixed function of its input — not a random-number generator
// — so GeneratePoints stays pure: same (r, n) always produces the same cloud.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// GeneratePoints returns n (x, y) pairs whose population correlation is r.
// Both x and y are standardized (mean 0, variance 1): x is standard-normal
// noise, and y = r*x + sqrt(1-r^2)*noise, so the correlation is exactly r in
// expectation and the theoretical trend line has slope r through the origin.
// The "noise" comes from hash01 fed through a Box-Muller transform, so the
// whole function is pure and deterministic — no math/rand, no time.
func GeneratePoints(r float64, n int) (xs, ys []float64) {
	if n < 1 {
		n = 1
	}
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	xs = make([]float64, n)
	ys = make([]float64, n)
	tail := math.Sqrt(1 - r*r)
	for i := 0; i < n; i++ {
		u1 := hash01(2*i + 1)
		u2 := hash01(2*i + 2)
		if u1 < 1e-9 {
			u1 = 1e-9 // keep log() finite
		}
		mag := math.Sqrt(-2 * math.Log(u1))
		z1 := mag * math.Cos(2*math.Pi*u2)
		z2 := mag * math.Sin(2*math.Pi*u2)
		xs[i] = z1
		ys[i] = r*z1 + tail*z2
	}
	return xs, ys
}

// PearsonR is the sample Pearson correlation coefficient of xs and ys. Pure
// math, unit-tested below. Returns 0 if either series has zero variance
// (undefined correlation, rather than dividing by zero).
func PearsonR(xs, ys []float64) float64 {
	n := len(xs)
	if n == 0 || len(ys) != n {
		return 0
	}
	var sx, sy float64
	for i := 0; i < n; i++ {
		sx += xs[i]
		sy += ys[i]
	}
	mx, my := sx/float64(n), sy/float64(n)

	var cov, vx, vy float64
	for i := 0; i < n; i++ {
		dx, dy := xs[i]-mx, ys[i]-my
		cov += dx * dy
		vx += dx * dx
		vy += dy * dy
	}
	if vx <= 0 || vy <= 0 {
		return 0
	}
	return cov / math.Sqrt(vx*vy)
}

func render(p map[string]float64) string {
	r, n := p["r"], int(p["n"])
	if r > 1 {
		r = 1
	}
	if r < -1 {
		r = -1
	}
	if n < 2 {
		n = 2
	}

	xs, ys := GeneratePoints(r, n)
	sampleR := PearsonR(xs, ys)

	const lim = 3.6
	c := viz.New(680, 400, -lim, lim, -lim, lim)
	c.Axes()
	for x := -3.0; x <= 3.0; x++ {
		c.Tick(x, fmt.Sprintf("%g", x))
	}

	// Reference lines through the origin (both axes are standardized).
	c.Path([][2]float64{{-lim, 0}, {lim, 0}}, viz.Muted, 1)
	c.Path([][2]float64{{0, -lim}, {0, lim}}, viz.Muted, 1)

	// The scatter cloud itself: one small square per point.
	for i := range xs {
		px, py := c.X(xs[i]), c.Y(ys[i])
		c.Rect(px-2.5, py-2.5, 5, 5, viz.Accent, 0.55)
	}

	// The trend line r implies: slope r through the origin, since x and y are
	// both standardized to mean 0, variance 1.
	c.Path([][2]float64{{-lim, -lim * r}, {lim, lim * r}}, viz.Warm, 2.5)

	c.Text(20, 24, fmt.Sprintf("target r = %.2f    sample r ≈ %.2f    n = %d", r, sampleR, n),
		14, viz.Ink, "start")
	c.Text(20, 44, "r = -1 downhill line   r = 0 no linear pattern   r = +1 uphill line",
		12, viz.Muted, "start")
	c.Text(20, 62, "a strong r means the points move together — it never proves one causes the other",
		12, viz.Muted, "start")
	return c.String()
}
