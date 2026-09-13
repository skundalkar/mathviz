// Package autocorr visualizes autocorrelation: how well a time series
// matches up with a lagged copy of itself, computed the same way
// correlation measures how two variables move together -- except here the
// "second variable" is the same series, shifted in time. Sweeping the lag
// traces a correlogram, and where that correlogram peaks reveals a series's
// hidden period even when the naked eye can't quite pin it down through the
// noise.
package autocorr

import (
	"fmt"
	"math"

	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "autocorrelation",
		Seq:   103,
		Title: "Autocorrelation (correlating a series with its lagged self)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"Placeholder -- filled in once the math and picture exist.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "lag", Label: "Lag (k)", Min: 0, Max: 20, Step: 1, Def: 3},
		},
		Render: render,
	})
}

// SeriesLength is the number of points in the fixed example series every
// Section walks through.
const SeriesLength = 40

// Period is the hidden repeating cycle baked into Series -- roughly 7 steps,
// deliberately not a clean divisor of SeriesLength so the noise doesn't
// happen to line up in a special way.
const Period = 7.0

// NoiseAmplitude bounds the uniform noise added to Series: each point is
// perturbed by an amount in [-NoiseAmplitude, +NoiseAmplitude].
const NoiseAmplitude = 3.0

// MaxLag is the largest lag every Section's slider sweeps through.
const MaxLag = 20

// hash01 turns an integer seed into a deterministic pseudo-random number in
// [0, 1) via a fixed irrational-multiplier trick -- not a real random number
// generator, just a fixed, reproducible sequence that merely looks
// patternless, which is exactly what Render needs (pure: same input, same
// output, no math/rand, no time). Same trick as internal/concepts/correlation.
func hash01(seed int) float64 {
	x := math.Sin(float64(seed)*12.9898) * 43758.5453123
	_, frac := math.Modf(x)
	if frac < 0 {
		frac += 1
	}
	return frac
}

// Series returns the fixed SeriesLength-point example every Section walks
// through: a period-7 sine wave (amplitude 10) plus bounded pseudo-random
// noise. Pure and deterministic -- the same 40 numbers every call.
func Series() []float64 {
	vals := make([]float64, SeriesLength)
	for t := 0; t < SeriesLength; t++ {
		signal := 10 * math.Sin(2*math.Pi*float64(t)/Period)
		noise := NoiseAmplitude * (2*hash01(t) - 1)
		vals[t] = signal + noise
	}
	return vals
}

// PearsonR is the sample Pearson correlation coefficient of xs and ys --
// the same formula internal/concepts/correlation defines, reproduced here so
// this package stays self-contained. Returns 0 if either series has zero
// variance (undefined correlation, rather than dividing by zero).
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

// ACF is the sample autocorrelation of series at the given lag: the Pearson
// correlation between the series and a copy of itself shifted forward by
// lag steps -- pairing series[lag], series[lag+1], ... against series[0],
// series[1], .... ACF(series, 0) is always 1 (a series is perfectly
// correlated with an unshifted copy of itself). A lag of len(series) or
// more, or a negative lag, leaves no pairs to compare and returns 0.
func ACF(series []float64, lag int) float64 {
	n := len(series)
	if lag < 0 || lag >= n {
		return 0
	}
	return PearsonR(series[lag:], series[:n-lag])
}

func render(p map[string]float64) string {
	lag := int(p["lag"] + 0.5)
	if lag < 0 {
		lag = 0
	}
	if lag > MaxLag {
		lag = MaxLag
	}

	series := Series()
	seriesMin, seriesMax := series[0], series[0]
	for _, v := range series {
		if v < seriesMin {
			seriesMin = v
		}
		if v > seriesMax {
			seriesMax = v
		}
	}
	pad := (seriesMax - seriesMin) * 0.15
	seriesMin -= pad
	seriesMax += pad

	// The series overlay lives in the lower band of the canvas, using the
	// canvas's one data-space mapping (x = time step, y = series value);
	// the correlogram above it is drawn in raw pixel space, same technique
	// dynamic-programming-knapsack and multi-armed-bandit use to keep two
	// panels from fighting over one coordinate system.
	c := viz.New(700, 560, 0, float64(SeriesLength-1), seriesMin, seriesMax)
	c.PadT = 360
	c.PadB = 40

	// The correlogram: one bar per lag 0..MaxLag, height/direction given by
	// ACF, drawn in raw pixel space above the series panel.
	const zeroY, scale = 200.0, 80.0
	const barW, spacing = 18.0, 30.0
	leftMargin := (700.0 - float64(MaxLag)*spacing) / 2

	c.Rect(leftMargin-16, zeroY-1, float64(MaxLag)*spacing+32, 2, viz.Muted, 0.9)
	c.Text(leftMargin-26, zeroY+4, "0", 11, viz.Muted, "end")
	c.Text(leftMargin-26, zeroY-scale+4, "+1", 11, viz.Muted, "end")
	c.Text(leftMargin-26, zeroY+scale+4, "-1", 11, viz.Muted, "end")

	acf := ACF(series, lag)
	for l := 0; l <= MaxLag; l++ {
		x := leftMargin + float64(l)*spacing
		lacf := ACF(series, l)
		barY := zeroY - lacf*scale
		top, height := barY, zeroY-barY
		if height < 0 {
			top, height = zeroY, -height
		}
		color, opacity := viz.Accent, 0.6
		if l == lag {
			color, opacity = viz.Warm, 0.95
		}
		c.Rect(x-barW/2, top, barW, height, color, opacity)
		if l%5 == 0 {
			c.Text(x, zeroY+scale+18, fmt.Sprintf("%d", l), 11, viz.Muted, "middle")
		}
	}
	c.Text(leftMargin, zeroY+scale+36, "lag", 11, viz.Muted, "start")

	c.Axes()
	for t := 0; t < SeriesLength; t += 5 {
		c.Tick(float64(t), fmt.Sprintf("t=%d", t))
	}

	// Original series, solid.
	pts := make([][2]float64, SeriesLength)
	for t, v := range series {
		pts[t] = [2]float64{float64(t), v}
	}
	c.Path(pts, viz.Accent, 2)

	// The same series shifted forward by lag steps, in orange: at time t
	// this plots series[t-lag], so it visually lines up what ACF is
	// comparing -- how well the value `lag` steps ago predicts the current
	// one.
	if lag > 0 && lag < SeriesLength {
		shifted := make([][2]float64, SeriesLength-lag)
		for t := lag; t < SeriesLength; t++ {
			shifted[t-lag] = [2]float64{float64(t), series[t-lag]}
		}
		c.Path(shifted, viz.Warm, 2)
	}

	verdict := "weak: barely related to a copy of itself this far off"
	switch {
	case acf > 0.7:
		verdict = "strong positive: a lagged copy tracks the original closely"
	case acf < -0.7:
		verdict = "strong negative: a lagged copy runs opposite the original"
	case acf > 0.3:
		verdict = "moderate positive"
	case acf < -0.3:
		verdict = "moderate negative"
	}
	c.Text(16, 24, fmt.Sprintf("lag = %d    ACF(%d) = %.3f    %s", lag, lag, acf, verdict), 13, viz.Ink, "start")
	c.Text(16, 44, fmt.Sprintf("period baked into this series ≈ %.0f steps -- correlogram peaks near lag %.0f and %.0f",
		Period, Period, 2*Period), 13, viz.Muted, "start")
	c.Text(16, 64, "orange bar = current lag    blue line = original series    orange line = series shifted by that lag",
		12, viz.Muted, "start")

	return c.String()
}
