package bootstrap

import (
	"math"
	"testing"

	"mathviz/internal/concept"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestMean(t *testing.T) {
	if got := Mean([]float64{12, 15, 15, 22, 41}); !near(got, 21) {
		t.Errorf("Mean = %v, want 21", got)
	}
	if got := Mean(nil); got != 0 {
		t.Errorf("Mean(nil) = %v, want 0", got)
	}
}

func TestBootstrapMeansCount(t *testing.T) {
	means := BootstrapMeans([]float64{1, 2, 3})
	if got, want := len(means), 27; got != want { // 3^3
		t.Errorf("len(BootstrapMeans) = %d, want %d", got, want)
	}
}

func TestBootstrapMeansEndpoints(t *testing.T) {
	// The all-same-value resample (index picked 5 times) always appears, so
	// the min/max of the returned means must equal the min/max of the data.
	data := []float64{12, 15, 15, 22, 41}
	means := BootstrapMeans(data)
	minM, maxM := means[0], means[0]
	for _, m := range means {
		if m < minM {
			minM = m
		}
		if m > maxM {
			maxM = m
		}
	}
	if !near(minM, 12) {
		t.Errorf("min bootstrap mean = %v, want 12", minM)
	}
	if !near(maxM, 41) {
		t.Errorf("max bootstrap mean = %v, want 41", maxM)
	}
}

// TestBootstrapMeansKnownFrequency checks a hand-countable case: for
// data=[1,2,3], a resample of size 3 has mean exactly 2 iff its three draws
// sum to 6. The sequences from {1,2,3}^3 summing to 6 are the 6 permutations
// of (1,2,3) plus (2,2,2) itself: 7 sequences out of 27.
func TestBootstrapMeansKnownFrequency(t *testing.T) {
	means := BootstrapMeans([]float64{1, 2, 3})
	count := 0
	for _, m := range means {
		if near(m, 2.0) {
			count++
		}
	}
	if want := 7; count != want {
		t.Errorf("count of resamples with mean==2 = %d, want %d", count, want)
	}
}

func TestPercentileIntervalKnownValues(t *testing.T) {
	means := BootstrapMeans(Sample(41)) // [12,15,15,22,41]

	lo, hi := PercentileInterval(means, 0.90)
	if !near(round1(lo), 14.4) || !near(round1(hi), 30.0) {
		t.Errorf("90%% interval = [%.2f, %.2f], want [14.4, 30.0]", lo, hi)
	}

	lo80, hi80 := PercentileInterval(means, 0.80)
	if !near(round1(lo80), 15.2) || !near(round1(hi80), 26.8) {
		t.Errorf("80%% interval = [%.2f, %.2f], want [15.2, 26.8]", lo80, hi80)
	}
}

func TestPercentileIntervalWidensWithOutlier(t *testing.T) {
	loSmall, hiSmall := PercentileInterval(BootstrapMeans(Sample(41)), 0.90)
	loBig, hiBig := PercentileInterval(BootstrapMeans(Sample(70)), 0.90)
	widthSmall := hiSmall - loSmall
	widthBig := hiBig - loBig
	if widthBig <= widthSmall {
		t.Errorf("expected a bigger outlier to widen the interval: %v -> %v", widthSmall, widthBig)
	}
}

func TestPercentileIntervalClampsLevel(t *testing.T) {
	means := []float64{1, 2, 3, 4, 5}
	lo, hi := PercentileInterval(means, 1.5) // clamps to 1.0 -> full range
	if !near(lo, 1) || !near(hi, 5) {
		t.Errorf("PercentileInterval with level>1 = [%v, %v], want [1, 5]", lo, hi)
	}
}

func round1(x float64) float64 {
	return math.Round(x*10) / 10
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("bootstrap-resampling")
	if !ok {
		t.Fatal("concept not registered")
	}
	svg := c.Render(c.Defaults())
	if len(svg) < 20 || svg[:4] != "<svg" {
		t.Errorf("Render did not produce an SVG document")
	}
}
