package logscale

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestLogBaseInverseOfPow(t *testing.T) {
	for _, b := range []float64{2, 3, 5, 10} {
		for k := 0; k <= 5; k++ {
			x := math.Pow(b, float64(k))
			if got := LogBase(x, b); math.Abs(got-float64(k)) > 1e-9 {
				t.Errorf("LogBase(%v, %v) = %v, want %d", x, b, got, k)
			}
		}
	}
}

func TestMultiplyingInputAddsOne(t *testing.T) {
	// The core lesson, asserted numerically: log_b(x*b) == log_b(x) + 1.
	b, x := 2.0, 7.0
	got := LogBase(x*b, b) - LogBase(x, b)
	if math.Abs(got-1) > 1e-12 {
		t.Fatalf("×base changed log by %v, want 1", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, _ := concept.Get("logarithms")
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
