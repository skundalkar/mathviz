package hmm

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestTransitionMatchesMarkovStepConvention(t *testing.T) {
	trans := Transition(0.6, 0.7)
	want := [2][2]float64{
		{0.6, 0.4},
		{0.3, 0.7},
	}
	for s := 0; s < 2; s++ {
		for d := 0; d < 2; d++ {
			if !almostEqual(trans[s][d], want[s][d], 1e-9) {
				t.Errorf("Transition(0.6,0.7)[%d][%d] = %v, want %v", s, d, trans[s][d], want[s][d])
			}
		}
	}
}

func TestTransitionRowsSumToOne(t *testing.T) {
	for _, a := range []float64{0.05, 0.5, 0.95} {
		for _, b := range []float64{0.05, 0.5, 0.95} {
			trans := Transition(a, b)
			for s := 0; s < 2; s++ {
				sum := trans[s][0] + trans[s][1]
				if !almostEqual(sum, 1, 1e-9) {
					t.Errorf("Transition(%v,%v)[%d] sums to %v, want 1", a, b, s, sum)
				}
			}
		}
	}
}

// TestViterbiKnownValues is pinned to the textbook Walk/Shop/Clean example
// worked step by step in LEARNINGS.md and the concept's own Sections.
func TestViterbiKnownValues(t *testing.T) {
	trans := Transition(0.6, 0.7)
	delta, path, prob := Viterbi(trans, Emission, InitialDist, Sequence)

	wantDelta := [][2]float64{
		{0.2400, 0.0600},
		{0.0432, 0.0384},
		{0.002592, 0.01344},
	}
	if len(delta) != len(wantDelta) {
		t.Fatalf("len(delta) = %d, want %d", len(delta), len(wantDelta))
	}
	for day, want := range wantDelta {
		for s := 0; s < 2; s++ {
			if !almostEqual(delta[day][s], want[s], 1e-6) {
				t.Errorf("delta[%d][%d] = %v, want %v", day, s, delta[day][s], want[s])
			}
		}
	}

	wantPath := []int{Sunny, Rainy, Rainy}
	if len(path) != len(wantPath) {
		t.Fatalf("len(path) = %d, want %d", len(path), len(wantPath))
	}
	for i, w := range wantPath {
		if path[i] != w {
			t.Errorf("path[%d] = %d, want %d", i, path[i], w)
		}
	}

	if !almostEqual(prob, 0.01344, 1e-6) {
		t.Errorf("prob = %v, want 0.01344", prob)
	}
}

func TestViterbiEmptyObservations(t *testing.T) {
	trans := Transition(0.6, 0.7)
	delta, path, prob := Viterbi(trans, Emission, InitialDist, nil)
	if delta != nil || path != nil || prob != 0 {
		t.Errorf("Viterbi(nil obs) = (%v, %v, %v), want (nil, nil, 0)", delta, path, prob)
	}
}

// TestViterbiMatchesBruteForce checks Viterbi's decoded path and probability
// against an exhaustive search over every possible hidden-state sequence,
// across several transition settings and observation sequences -- the
// defining property of the algorithm, independent of any one pinned example.
func TestViterbiMatchesBruteForce(t *testing.T) {
	obsSeqs := [][]int{
		{Walk, Shop, Clean},
		{Clean, Clean, Clean},
		{Walk, Walk, Shop, Clean},
	}
	settings := []struct{ a, b float64 }{
		{0.6, 0.7}, {0.2, 0.9}, {0.5, 0.5}, {0.9, 0.1},
	}

	for _, s := range settings {
		trans := Transition(s.a, s.b)
		for _, obs := range obsSeqs {
			_, gotPath, gotProb := Viterbi(trans, Emission, InitialDist, obs)
			bruteProb, bruteAnyPath := bruteForceBest(trans, Emission, InitialDist, obs)

			if !almostEqual(gotProb, bruteProb, 1e-12) {
				t.Errorf("a=%v b=%v obs=%v: Viterbi prob = %v, brute force best = %v",
					s.a, s.b, obs, gotProb, bruteProb)
			}
			gotProbOfPath := pathProbability(trans, Emission, InitialDist, obs, gotPath)
			if !almostEqual(gotProbOfPath, bruteProb, 1e-12) {
				t.Errorf("a=%v b=%v obs=%v: Viterbi's own path %v scores %v, want brute force best %v (any optimal path %v)",
					s.a, s.b, obs, gotPath, gotProbOfPath, bruteProb, bruteAnyPath)
			}
		}
	}
}

// pathProbability computes P(path, obs) for one specific fully-specified
// hidden-state path -- the same product Viterbi's delta update accumulates,
// evaluated directly instead of via dynamic programming.
func pathProbability(trans [2][2]float64, emit [2][3]float64, initial [2]float64, obs []int, path []int) float64 {
	p := initial[path[0]] * emit[path[0]][obs[0]]
	for t := 1; t < len(obs); t++ {
		p *= trans[path[t-1]][path[t]] * emit[path[t]][obs[t]]
	}
	return p
}

// bruteForceBest enumerates every possible hidden-state sequence of the
// same length as obs and returns the highest P(path, obs) found, plus one
// path that achieves it.
func bruteForceBest(trans [2][2]float64, emit [2][3]float64, initial [2]float64, obs []int) (float64, []int) {
	n := len(obs)
	best := -1.0
	var bestPath []int
	for mask := 0; mask < (1 << n); mask++ {
		path := make([]int, n)
		for i := 0; i < n; i++ {
			path[i] = (mask >> i) & 1
		}
		p := pathProbability(trans, emit, initial, obs, path)
		if p > best {
			best, bestPath = p, path
		}
	}
	return best, bestPath
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("hidden-markov-models")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
