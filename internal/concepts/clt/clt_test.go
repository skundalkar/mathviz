package clt

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestExponentialPDFKnownValues(t *testing.T) {
	if got, want := ExponentialPDF(0, 1), 1.0; math.Abs(got-want) > 1e-12 {
		t.Errorf("ExponentialPDF(0,1) = %v, want %v", got, want)
	}
	if got, want := ExponentialPDF(1, 1), math.Exp(-1); math.Abs(got-want) > 1e-12 {
		t.Errorf("ExponentialPDF(1,1) = %v, want %v", got, want)
	}
	if got := ExponentialPDF(-1, 1); got != 0 {
		t.Errorf("ExponentialPDF(-1,1) = %v, want 0", got)
	}
	if got := ExponentialPDF(1, 0); got != 0 {
		t.Errorf("ExponentialPDF(1,0) = %v, want 0 (invalid rate)", got)
	}
}

func TestSampleMeanPDFAtNEqualsOneIsThePopulation(t *testing.T) {
	// A "sample" of size 1 is just the population itself.
	for _, lambda := range []float64{0.5, 1, 2} {
		for _, x := range []float64{0.1, 0.5, 1, 2, 4} {
			got, want := SampleMeanPDF(x, 1, lambda), ExponentialPDF(x, lambda)
			if math.Abs(got-want) > 1e-9 {
				t.Errorf("SampleMeanPDF(%v,1,%v) = %v, want %v", x, lambda, got, want)
			}
		}
	}
}

func TestSampleMeanPDFInvalidInputs(t *testing.T) {
	if got := SampleMeanPDF(1, 0, 1); got != 0 {
		t.Errorf("SampleMeanPDF with n=0 = %v, want 0", got)
	}
	if got := SampleMeanPDF(1, 5, 0); got != 0 {
		t.Errorf("SampleMeanPDF with lambda=0 = %v, want 0", got)
	}
	if got := SampleMeanPDF(-1, 5, 1); got != 0 {
		t.Errorf("SampleMeanPDF with m<0 = %v, want 0", got)
	}
}

func TestSampleMeanPDFIntegratesToOne(t *testing.T) {
	for _, tc := range []struct {
		n      int
		lambda float64
	}{
		{1, 1}, {2, 1}, {5, 0.5}, {15, 2}, {30, 1},
	} {
		lo, hi, steps := 1e-6, 50.0/tc.lambda, 40000
		h := (hi - lo) / float64(steps)
		sum := 0.0
		for i := 0; i <= steps; i++ {
			x := lo + float64(i)*h
			w := 1.0
			if i == 0 || i == steps {
				w = 0.5
			}
			sum += w * SampleMeanPDF(x, tc.n, tc.lambda)
		}
		if area := sum * h; math.Abs(area-1) > 1e-2 {
			t.Errorf("n=%d lambda=%v: integral = %v, want ≈ 1", tc.n, tc.lambda, area)
		}
	}
}

func TestSampleMeanVarianceShrinksWithN(t *testing.T) {
	lambda := 1.5
	v1 := SampleMeanVariance(1, lambda)
	v10 := SampleMeanVariance(10, lambda)
	v30 := SampleMeanVariance(30, lambda)
	if !(v1 > v10 && v10 > v30) {
		t.Fatalf("variance should shrink as n grows: v1=%v v10=%v v30=%v", v1, v10, v30)
	}
	want1 := 1 / (lambda * lambda)
	if math.Abs(v1-want1) > 1e-12 {
		t.Errorf("SampleMeanVariance(1,%v) = %v, want %v", lambda, v1, want1)
	}
	want30 := 1 / (30 * lambda * lambda)
	if math.Abs(v30-want30) > 1e-12 {
		t.Errorf("SampleMeanVariance(30,%v) = %v, want %v", lambda, v30, want30)
	}
}

func TestSampleMeanPDFApproachesNormalAsNGrows(t *testing.T) {
	// The central limit theorem: at large n the sample-mean density near its
	// own mean should be close to the matching normal's peak density.
	lambda := 1.0
	mean := 1 / lambda
	sigma30 := math.Sqrt(SampleMeanVariance(30, lambda))
	normalPeak30 := NormalPDF(mean, mean, sigma30)
	samplePeak30 := SampleMeanPDF(mean, 30, lambda)
	if rel := math.Abs(samplePeak30-normalPeak30) / normalPeak30; rel > 0.05 {
		t.Errorf("n=30 peak %v too far from normal peak %v (rel err %v)", samplePeak30, normalPeak30, rel)
	}

	// n=1 (the raw exponential) should NOT be close to its "matching normal":
	// the whole point is that small-n starts out far from normal.
	sigma1 := math.Sqrt(SampleMeanVariance(1, lambda))
	normalPeak1 := NormalPDF(mean, mean, sigma1)
	samplePeak1 := SampleMeanPDF(mean, 1, lambda)
	if rel := math.Abs(samplePeak1-normalPeak1) / normalPeak1; rel < 0.05 {
		t.Errorf("n=1 peak %v unexpectedly close to normal peak %v", samplePeak1, normalPeak1)
	}
}

func TestNormalPDFKnownValues(t *testing.T) {
	got := NormalPDF(0, 0, 1)
	want := 1 / math.Sqrt(2*math.Pi)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("NormalPDF(0,0,1) = %v, want %v", got, want)
	}
	if got := NormalPDF(0, 0, 0); got != 0 {
		t.Errorf("NormalPDF with sigma=0 = %v, want 0", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("central-limit-theorem")
	if !ok {
		t.Fatal("concept \"central-limit-theorem\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
