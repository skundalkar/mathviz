package lln

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestOutcomeIsZeroOrOne(t *testing.T) {
	for trial := 1; trial <= 300; trial++ {
		got := Outcome(trial, 0.5)
		if got != 0 && got != 1 {
			t.Fatalf("Outcome(%d, 0.5) = %v, want 0 or 1", trial, got)
		}
	}
}

func TestOutcomeIsDeterministic(t *testing.T) {
	for trial := 1; trial <= 20; trial++ {
		a := Outcome(trial, 0.5)
		b := Outcome(trial, 0.5)
		if a != b {
			t.Errorf("Outcome(%d,0.5) not deterministic: %v then %v", trial, a, b)
		}
	}
}

func TestOutcomeAlwaysHeadsAtP1(t *testing.T) {
	for trial := 1; trial <= 50; trial++ {
		if got := Outcome(trial, 1.0); got != 1 {
			t.Errorf("Outcome(%d, 1.0) = %v, want 1 (certain heads)", trial, got)
		}
	}
}

func TestOutcomeNeverHeadsAtP0(t *testing.T) {
	for trial := 1; trial <= 50; trial++ {
		if got := Outcome(trial, 0.0); got != 0 {
			t.Errorf("Outcome(%d, 0.0) = %v, want 0 (certain tails)", trial, got)
		}
	}
}

func TestRunningAveragesLength(t *testing.T) {
	if got := len(RunningAverages(300, 0.5)); got != 300 {
		t.Errorf("len(RunningAverages(300,0.5)) = %d, want 300", got)
	}
	if got := len(RunningAverages(0, 0.5)); got != 1 {
		t.Errorf("len(RunningAverages(0,0.5)) = %d, want 1 (clamped)", got)
	}
}

func TestRunningAveragesStaysInUnitRange(t *testing.T) {
	for _, avg := range RunningAverages(300, 0.5) {
		if avg < 0 || avg > 1 {
			t.Errorf("running average %v out of [0,1]", avg)
		}
	}
}

func TestRunningAverageMatchesRunningAveragesLastEntry(t *testing.T) {
	for _, n := range []int{1, 5, 50, 300} {
		series := RunningAverages(n, 0.5)
		want := series[len(series)-1]
		if got := RunningAverage(n, 0.5); !almostEqual(got, want, 1e-12) {
			t.Errorf("RunningAverage(%d,0.5) = %v, want %v (RunningAverages' last entry)", n, got, want)
		}
	}
}

func TestRunningAverageKnownEarlyValues(t *testing.T) {
	// Pinned to the deterministic hash01 sequence at p=0.5, so a change to
	// the sequence (or a broken running-average calculation) shows up
	// immediately. See LEARNINGS.md for the same worked numbers.
	cases := map[int]float64{1: 0.0, 10: 0.5, 50: 0.52}
	for n, want := range cases {
		if got := RunningAverage(n, 0.5); !almostEqual(got, want, 1e-9) {
			t.Errorf("RunningAverage(%d,0.5) = %v, want %v", n, got, want)
		}
	}
}

func TestRunningAverageVarianceShrinksWithN(t *testing.T) {
	// The defining behavior of the law of large numbers: the running
	// average bounces around a lot early on and settles down later. Check
	// it on this concrete sequence rather than just asserting it in prose:
	// the spread (variance) of the running average over the first 20
	// flips should be much larger than its spread over the last 20 of 300.
	series := RunningAverages(300, 0.5)
	variance := func(vals []float64) float64 {
		mean := 0.0
		for _, v := range vals {
			mean += v
		}
		mean /= float64(len(vals))
		sq := 0.0
		for _, v := range vals {
			sq += (v - mean) * (v - mean)
		}
		return sq / float64(len(vals))
	}
	early := variance(series[:20])
	late := variance(series[280:])
	if late >= early {
		t.Errorf("variance did not shrink: early(n<=20)=%v, late(n in 281..300)=%v", early, late)
	}
	// Not just smaller -- an order of magnitude smaller, matching the LLN's
	// claim that spread shrinks roughly like 1/n, not just "a bit."
	if late*10 >= early {
		t.Errorf("variance shrank less than 10x: early=%v, late=%v", early, late)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("law-of-large-numbers")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
