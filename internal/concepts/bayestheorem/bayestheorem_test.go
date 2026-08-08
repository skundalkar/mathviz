package bayestheorem

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestPosteriorPositiveKnownValue(t *testing.T) {
	// prior=1%, sensitivity=99%, specificity=95% is a clean case:
	// P(C|+) = (0.99*0.01) / (0.99*0.01 + 0.05*0.99) = 0.0099/0.0594 = 1/6.
	got := PosteriorPositive(0.01, 0.99, 0.95)
	want := 1.0 / 6.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("PosteriorPositive(0.01, 0.99, 0.95) = %v, want %v", got, want)
	}
}

func TestPosteriorPositiveIncreasesWithPrior(t *testing.T) {
	prev := PosteriorPositive(0.001, 0.9, 0.9)
	for _, prior := range []float64{0.01, 0.1, 0.3, 0.5} {
		got := PosteriorPositive(prior, 0.9, 0.9)
		if got <= prev {
			t.Fatalf("PosteriorPositive should rise with prior: prior=%v got=%v <= prev=%v", prior, got, prev)
		}
		prev = got
	}
}

func TestPosteriorPositiveIncreasesWithSpecificity(t *testing.T) {
	// Fewer false alarms (higher specificity) should make a positive test
	// more trustworthy, holding prior and sensitivity fixed.
	prev := PosteriorPositive(0.02, 0.9, 0.5)
	for _, spec := range []float64{0.7, 0.9, 0.95, 0.99} {
		got := PosteriorPositive(0.02, 0.9, spec)
		if got <= prev {
			t.Fatalf("PosteriorPositive should rise with specificity: spec=%v got=%v <= prev=%v", spec, got, prev)
		}
		prev = got
	}
}

func TestPosteriorPositiveBounds(t *testing.T) {
	if got := PosteriorPositive(0, 0.99, 0.95); got != 0 {
		t.Errorf("PosteriorPositive with prior=0 = %v, want 0", got)
	}
	// A perfect test (no false positives) should make a positive result
	// certain, regardless of how rare the condition is.
	if got := PosteriorPositive(0.001, 0.99, 1.0); math.Abs(got-1) > 1e-9 {
		t.Errorf("PosteriorPositive with specificity=1 = %v, want 1", got)
	}
}

func TestPosteriorNegativeKnownValue(t *testing.T) {
	got := PosteriorNegative(0.01, 0.99, 0.95)
	want := (0.95 * 0.99) / (0.95*0.99 + 0.01*0.01)
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("PosteriorNegative(0.01, 0.99, 0.95) = %v, want %v", got, want)
	}
	if got < 0.999 {
		t.Errorf("PosteriorNegative(0.01, 0.99, 0.95) = %v, want > 0.999 (a negative should be reassuring when the condition is rare)", got)
	}
}

func TestCountsSumToPopulation(t *testing.T) {
	cases := []struct {
		prior, sens, spec float64
		n                 int
	}{
		{0.01, 0.99, 0.95, 100},
		{0.5, 0.8, 0.8, 100},
		{0.2, 0.7, 0.9, 250},
		{0.001, 0.999, 0.999, 100},
	}
	for _, c := range cases {
		tp, fn, fp, tn := Counts(c.prior, c.sens, c.spec, c.n)
		if got := tp + fn + fp + tn; got != c.n {
			t.Errorf("Counts(%v,%v,%v,%d) sums to %d, want %d", c.prior, c.sens, c.spec, c.n, got, c.n)
		}
		if tp < 0 || fn < 0 || fp < 0 || tn < 0 {
			t.Errorf("Counts(%v,%v,%v,%d) produced a negative count: tp=%d fn=%d fp=%d tn=%d",
				c.prior, c.sens, c.spec, c.n, tp, fn, fp, tn)
		}
	}
}

func TestCountsBalancedSplit(t *testing.T) {
	// prior=0.5, perfect test: half the population is tp, half is tn.
	tp, fn, fp, tn := Counts(0.5, 1.0, 1.0, 100)
	if tp != 50 || fn != 0 || fp != 0 || tn != 50 {
		t.Errorf("Counts(0.5, 1, 1, 100) = tp=%d fn=%d fp=%d tn=%d, want 50/0/0/50", tp, fn, fp, tn)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("bayes-theorem")
	if !ok {
		t.Fatal("concept \"bayes-theorem\" not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
