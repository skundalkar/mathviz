package birthday

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestProbNoCollisionSinglePerson(t *testing.T) {
	if got := ProbNoCollision(1, 365); got != 1 {
		t.Errorf("ProbNoCollision(1,365) = %v, want 1", got)
	}
	if got := ProbNoCollision(0, 365); got != 1 {
		t.Errorf("ProbNoCollision(0,365) = %v, want 1", got)
	}
}

func TestProbNoCollisionTwoPeople(t *testing.T) {
	// The 2nd person just has to dodge the 1st's one taken day: 364/365.
	want := 364.0 / 365.0
	if got := ProbNoCollision(2, 365); !almostEqual(got, want, 1e-9) {
		t.Errorf("ProbNoCollision(2,365) = %v, want %v", got, want)
	}
}

func TestProbCollisionKnown23People(t *testing.T) {
	// The textbook result: with 365 days, 23 people crosses 50% -- this is
	// the whole "paradox," so pin the number down precisely.
	got := ProbCollision(23, 365)
	if got < 0.5 || got > 0.51 {
		t.Errorf("ProbCollision(23,365) = %v, want in [0.50, 0.51]", got)
	}
}

func TestProbCollision22PeopleIsBelow50Percent(t *testing.T) {
	got := ProbCollision(22, 365)
	if got >= 0.5 {
		t.Errorf("ProbCollision(22,365) = %v, want < 0.5", got)
	}
}

func TestProbNoCollisionMoreThanDaysIsImpossible(t *testing.T) {
	// Pigeonhole principle: with only 5 possible days, 6 people can't all
	// have distinct birthdays.
	if got := ProbNoCollision(6, 5); got != 0 {
		t.Errorf("ProbNoCollision(6,5) = %v, want 0", got)
	}
	if got := ProbCollision(6, 5); got != 1 {
		t.Errorf("ProbCollision(6,5) = %v, want 1", got)
	}
}

func TestProbNoCollisionExactlyDaysPeople(t *testing.T) {
	// n == days is still possible (everyone takes a distinct day, no
	// spares) -- should be strictly positive, not zero.
	got := ProbNoCollision(5, 5)
	if got <= 0 || got > 1 {
		t.Errorf("ProbNoCollision(5,5) = %v, want in (0, 1]", got)
	}
}

func TestProbCollisionIsMonotonicInN(t *testing.T) {
	// Adding another person can only make a collision more likely, never
	// less -- ProbCollision(n) should never decrease as n grows.
	prev := 0.0
	for n := 1; n <= 60; n++ {
		got := ProbCollision(n, 365)
		if got < prev-1e-12 {
			t.Errorf("ProbCollision(%d,365) = %v, less than ProbCollision(%d,365) = %v", n, got, n-1, prev)
		}
		prev = got
	}
}

func TestProbCollisionStaysInUnitRange(t *testing.T) {
	for n := 0; n <= 400; n++ {
		got := ProbCollision(n, 365)
		if got < 0 || got > 1 {
			t.Errorf("ProbCollision(%d,365) = %v, want in [0,1]", n, got)
		}
	}
}

func TestMinPeopleForProbabilityMatchesTextbook23(t *testing.T) {
	if got := MinPeopleForProbability(365, 0.5); got != 23 {
		t.Errorf("MinPeopleForProbability(365,0.5) = %d, want 23", got)
	}
}

func TestMinPeopleForProbabilitySmallDayCount(t *testing.T) {
	// With only 5 possible days, collisions should cross 50% much sooner
	// than with 365.
	got := MinPeopleForProbability(5, 0.5)
	if got < 2 || got > 5 {
		t.Errorf("MinPeopleForProbability(5,0.5) = %d, want a small n in [2,5]", got)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("birthday-paradox")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
