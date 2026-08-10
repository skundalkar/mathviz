// Package meanmedianmode visualizes how a skewed distribution pulls the mean,
// median, and mode apart. All three coincide for a symmetric distribution, but
// as skew increases they peel off in a fixed order: mode stays at the peak,
// the median sits at the halfway point, and the mean gets dragged toward the
// long tail. A log-normal curve is used because all three statistics have
// simple closed forms, so the picture is exact rather than sampled.
package meanmedianmode

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "mean-median-mode",
		Seq:   5,
		Title: "Mean, median & mode",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Bill Gates walks into a small bar with 20 regulars, each with a net worth of " +
						"about $80,000. Before he sits down, the total net worth in the room is " +
						"20 × $80,000 = $1,600,000, so the average (mean) is $1,600,000 ÷ 20 = " +
						"$80,000 — a fair summary of 'the typical person here.'",
					"Now Bill Gates sits down. His real net worth is roughly $100 billion — round " +
						"it to exactly $100,000,000,000. The room's total is now $1,600,000 + " +
						"$100,000,000,000 = $100,001,600,000, split across 21 people: " +
						"$100,001,600,000 ÷ 21 ≈ $4.76 billion. That's the new mean — even though " +
						"not one of the 20 regulars got a dollar richer or poorer. 'The average " +
						"person in this bar' now sounds like a multi-billionaire. If the mean can " +
						"say that with a straight face, is it even still answering 'what's typical " +
						"here' — and if not, what would?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"The median doesn't have this problem. Sort all 21 net worths from lowest to " +
						"highest: twenty $80,000 values, then Gates's $100 billion alone at the far " +
						"end. The middle of that line — the 11th of 21 — is still one of the " +
						"ordinary $80,000 regulars, because Gates is just one more name standing at " +
						"the far end, not a number that gets averaged in. Median: still $80,000, " +
						"barely moved. The mode — the single most common value — doesn't move " +
						"either: $80,000 is still shared by 20 of the 21 people.",
					"That gap generalizes into a fixed order. Five salaries: $40k, $40k, $45k, " +
						"$50k, $200k:",
					"• Mode = $40k (the only value that repeats)",
					"• Median = $45k (sort them — 40k, 40k, 45k, 50k, 200k — and take the " +
						"middle, 3rd of 5)",
					"• Mean = (40k+40k+45k+50k+200k) ÷ 5 = 375k ÷ 5 = $75k",
					"$40k < $45k < $75k — mode, then median, then mean, every time a distribution " +
						"skews, because the mean is the only one of the three that does arithmetic " +
						"with the outlier instead of just counting or ranking it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The curve is a log-normal distribution, whose mean, median, and mode all " +
						"have simple closed forms, so the picture is exact rather than sampled. " +
						"Drag the skew (σ) slider and watch mode, median, and mean peel apart from " +
						"a single starting point in that same fixed order — mode < median < mean — " +
						"as the tail stretches out, exactly like the $40k / $45k / $75k salary " +
						"example above; the log-location (μ) slider shifts the whole curve left or " +
						"right without changing that order.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"You can now tell which 'typical' value a number is quoting, and reason about " +
						"whether it's the one you actually want — when someone reports 'the " +
						"average,' you know to ask which statistic they mean and whether a skewed " +
						"tail (like Gates in the bar) is inflating it.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"'Average household income' headlines are almost always higher than what a " +
						"typical household actually earns, because a small number of very high " +
						"earners drag the mean up while the median barely moves — exactly why " +
						"income and home-price statistics are usually reported as medians, and why " +
						"a company's 'average' salary can look great in a press release while most " +
						"employees make noticeably less.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: whenever someone says 'the average,' it's worth asking " +
						"which one they mean — 'average household income' almost always means the " +
						"mean, and it's usually well above what a typical household actually earns.",
					"Not like this: 'half of people earn less than the average' — only guaranteed " +
						"true if 'average' means the median (the middle value by definition); it's " +
						"often false for the mean, which a handful of outliers can drag well above " +
						"where most people actually sit.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "sigma", Label: "Skew (σ)", Min: 0.05, Max: 1.2, Step: 0.05, Def: 0.6},
			{Key: "mu", Label: "Log-location (μ)", Min: -0.5, Max: 0.5, Step: 0.05, Def: 0},
		},
		Render: render,
	})
}

// LogNormalPDF is the probability density of a log-normal distribution with
// log-space parameters mu, sigma. Pure math, unit-tested below.
func LogNormalPDF(x, mu, sigma float64) float64 {
	if x <= 0 || sigma <= 0 {
		return 0
	}
	z := (math.Log(x) - mu) / sigma
	return math.Exp(-0.5*z*z) / (x * sigma * math.Sqrt(2*math.Pi))
}

// Mean is the closed-form mean of a log-normal(mu, sigma) distribution.
func Mean(mu, sigma float64) float64 {
	return math.Exp(mu + sigma*sigma/2)
}

// Median is the closed-form median of a log-normal(mu, sigma) distribution.
func Median(mu, sigma float64) float64 {
	return math.Exp(mu)
}

// Mode is the closed-form mode (location of peak density) of a
// log-normal(mu, sigma) distribution.
func Mode(mu, sigma float64) float64 {
	return math.Exp(mu - sigma*sigma)
}

func render(p map[string]float64) string {
	mu, sigma := p["mu"], p["sigma"]
	if sigma <= 0 {
		sigma = 0.05
	}

	mean, median, mode := Mean(mu, sigma), Median(mu, sigma), Mode(mu, sigma)

	xmax := math.Exp(mu+3.2*sigma) * 1.1
	peak := LogNormalPDF(mode, mu, sigma)
	c := viz.New(680, 340, 0, xmax, 0, peak*1.15)
	c.Axes()

	step := xmax / 8
	for x := step; x <= xmax; x += step {
		c.Tick(x, fmt.Sprintf("%.1f", x))
	}

	curve := viz.Sample(0.001, xmax, 300, func(x float64) float64 {
		return LogNormalPDF(x, mu, sigma)
	})
	c.Path(curve, viz.Accent, 2.5)

	c.VLine(mode, viz.Good, true)
	c.VLine(median, viz.Warm, true)
	c.VLine(mean, viz.Bad, false)

	c.Text(c.X(mode), c.PadT+12, "mode", 12, viz.Good, "middle")
	c.Text(c.X(median), c.PadT+28, "median", 12, viz.Warm, "middle")
	c.Text(c.X(mean), c.PadT+44, "mean", 12, viz.Bad, "middle")

	c.Text(20, 24, fmt.Sprintf("mode = %.2f    median = %.2f    mean = %.2f", mode, median, mean),
		14, viz.Ink, "start")
	c.Text(20, 44, "right-skewed: the long tail drags the mean past the median, past the mode",
		12, viz.Muted, "start")

	return c.String()
}
