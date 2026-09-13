// Package autocorr visualizes autocorrelation: how well a time series
// matches up with a lagged copy of itself, computed the same way
// correlation measures how two variables move together -- except here the
// "second variable" is the same series, shifted in time. Sweeping the lag
// traces a correlogram, and where that correlogram peaks reveals a series's
// hidden period even when the naked eye can't quite pin it down through the
// noise.
package autocorr

import (
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
	return viz.New(700, 460, 0, 1, 0, 1).String()
}
