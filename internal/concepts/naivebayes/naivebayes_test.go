package naivebayes

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

const epsilon = 1e-9

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestOddsProbabilityFromOddsRoundTrip(t *testing.T) {
	for _, p := range []float64{0.1, 0.25, 0.5, 0.75, 0.9153} {
		got := ProbabilityFromOdds(Odds(p))
		if !almostEqual(got, p, 1e-9) {
			t.Errorf("ProbabilityFromOdds(Odds(%v)) = %v, want %v", p, got, p)
		}
	}
}

func TestOddsKnownValues(t *testing.T) {
	if got := Odds(0.5); !almostEqual(got, 1, epsilon) {
		t.Errorf("Odds(0.5) = %v, want 1", got)
	}
	if got := Odds(0.9); !almostEqual(got, 9, epsilon) {
		t.Errorf("Odds(0.9) = %v, want 9", got)
	}
}

func TestWordLikelihoodRatioWorkedExample(t *testing.T) {
	// From the section-2 walkthrough: "free" present -> ratio 6; "meeting"
	// absent -> ratio 1.8.
	if got := WordLikelihoodRatio(true, PFreeGivenSpam, PFreeGivenHam); !almostEqual(got, 6, epsilon) {
		t.Errorf("WordLikelihoodRatio(true, free) = %v, want 6", got)
	}
	if got := WordLikelihoodRatio(false, PMeetingGivenSpam, PMeetingGivenHam); !almostEqual(got, 1.8, epsilon) {
		t.Errorf("WordLikelihoodRatio(false, meeting) = %v, want 1.8", got)
	}
}

func TestWordLikelihoodRatioAbsentIsComplement(t *testing.T) {
	// Present + absent ratios for the same word use complementary
	// probabilities, not the same numbers.
	present := WordLikelihoodRatio(true, PFreeGivenSpam, PFreeGivenHam)
	absent := WordLikelihoodRatio(false, PFreeGivenSpam, PFreeGivenHam)
	wantAbsent := (1 - PFreeGivenSpam) / (1 - PFreeGivenHam)
	if !almostEqual(absent, wantAbsent, epsilon) {
		t.Errorf("WordLikelihoodRatio(false, free) = %v, want %v", absent, wantAbsent)
	}
	if almostEqual(present, absent, 1e-6) {
		t.Error("present and absent likelihood ratios for the same word should differ")
	}
}

func TestPosteriorProbabilityWorkedExample(t *testing.T) {
	// From the section-2 walkthrough: 50% prior, "free" present, "meeting"
	// absent -> posterior odds 10.8:1 -> posterior probability ~91.5%.
	ratioFree := WordLikelihoodRatio(true, PFreeGivenSpam, PFreeGivenHam)
	ratioMeeting := WordLikelihoodRatio(false, PMeetingGivenSpam, PMeetingGivenHam)
	got := PosteriorProbability(0.5, ratioFree, ratioMeeting)
	want := 10.8 / 11.8
	if !almostEqual(got, want, 1e-9) {
		t.Errorf("PosteriorProbability(0.5, 6, 1.8) = %v, want %v (~91.5%%)", got, want)
	}
}

func TestPosteriorProbabilityNoEvidenceReturnsPrior(t *testing.T) {
	for _, prior := range []float64{0.1, 0.5, 0.9} {
		if got := PosteriorProbability(prior); !almostEqual(got, prior, 1e-9) {
			t.Errorf("PosteriorProbability(%v) with no ratios = %v, want %v", prior, got, prior)
		}
	}
}

func TestPosteriorProbabilityIsMonotonicInRatio(t *testing.T) {
	prev := PosteriorProbability(0.5, 0.1)
	for _, r := range []float64{0.5, 1, 2, 6, 20} {
		cur := PosteriorProbability(0.5, r)
		if cur <= prev {
			t.Errorf("PosteriorProbability(0.5, %v) = %v, want > previous %v", r, cur, prev)
		}
		prev = cur
	}
}

func TestPosteriorProbabilityRatioOneLeavesPriorUnchanged(t *testing.T) {
	// A likelihood ratio of 1 means the evidence is equally likely under
	// both classes -- it carries no information and shouldn't move the
	// posterior at all.
	if got := PosteriorProbability(0.3, 1, 1, 1); !almostEqual(got, 0.3, 1e-9) {
		t.Errorf("PosteriorProbability(0.3, 1, 1, 1) = %v, want 0.3", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("naive-bayes")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
