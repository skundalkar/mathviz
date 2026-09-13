package bandit

import (
	"math"
	"strings"
	"testing"

	"mathviz/internal/concept"
)

func TestArmOutcomeDeterministic(t *testing.T) {
	for arm := 0; arm < 3; arm++ {
		for occ := 0; occ < 10; occ++ {
			a := ArmOutcome(arm, occ)
			b := ArmOutcome(arm, occ)
			if a != b {
				t.Errorf("ArmOutcome(%d,%d) not deterministic: %v then %v", arm, occ, a, b)
			}
		}
	}
}

// TestArmOutcomeMatchesTrueProbApproximately checks that, over many
// occurrences, the fraction of wins ArmOutcome reports for an arm lands
// close to that arm's TrueProbs entry -- ArmOutcome is a Bernoulli draw at
// that rate, not an arbitrary deterministic pattern.
func TestArmOutcomeMatchesTrueProbApproximately(t *testing.T) {
	const n = 5000
	for arm, want := range TrueProbs {
		wins := 0
		for occ := 0; occ < n; occ++ {
			if ArmOutcome(arm, occ) {
				wins++
			}
		}
		got := float64(wins) / float64(n)
		if math.Abs(got-want) > 0.03 {
			t.Errorf("arm %d: observed win rate %.3f over %d pulls, want close to true prob %.2f", arm, got, n, want)
		}
	}
}

func TestSimulateLength(t *testing.T) {
	if got := len(Simulate(0.2, 60)); got != 60 {
		t.Errorf("len(Simulate(0.2,60)) = %d, want 60", got)
	}
	if got := len(Simulate(0.2, 0)); got != 0 {
		t.Errorf("len(Simulate(0.2,0)) = %d, want 0", got)
	}
	if got := len(Simulate(0.2, -5)); got != 0 {
		t.Errorf("len(Simulate(0.2,-5)) = %d, want 0 (clamped)", got)
	}
}

// TestSimulateForcedFirstPulls checks that the first len(TrueProbs) pulls
// always visit each arm once, in order, and never count as "explored" --
// even at epsilon=1.0, where exploreDecision would otherwise always fire.
func TestSimulateForcedFirstPulls(t *testing.T) {
	steps := Simulate(1.0, len(TrueProbs))
	for i, s := range steps {
		if s.Arm != i {
			t.Errorf("forced pull %d: Arm = %d, want %d", i, s.Arm, i)
		}
		if s.Explored {
			t.Errorf("forced pull %d: Explored = true, want false", i)
		}
	}
}

// TestSimulateCountsSumToSteps checks that every pull is attributed to
// exactly one arm.
func TestSimulateCountsSumToSteps(t *testing.T) {
	for _, eps := range []float64{0, 0.2, 0.5, 1.0} {
		steps := Simulate(eps, 137)
		last := steps[len(steps)-1]
		sum := last.Counts[0] + last.Counts[1] + last.Counts[2]
		if sum != 137 {
			t.Errorf("epsilon=%.2f: counts sum to %d, want 137", eps, sum)
		}
	}
}

// TestSimulateCumRegretNonDecreasing checks the running cost of not always
// pulling the best arm only ever grows -- every pull of a suboptimal arm
// adds a non-negative amount (best-chosenProb >= 0 by definition of best),
// and pulling the best arm itself adds exactly 0.
func TestSimulateCumRegretNonDecreasing(t *testing.T) {
	steps := Simulate(0.3, 200)
	prev := 0.0
	for i, s := range steps {
		if s.CumRegret < prev-1e-12 {
			t.Fatalf("step %d: CumRegret decreased from %v to %v", i, prev, s.CumRegret)
		}
		prev = s.CumRegret
	}
}

// TestSimulateKnownEarlyValues is pinned to the deterministic hash01
// sequence at epsilon=0.2, the concept's default -- the same numbers walked
// through in LEARNINGS.md and the concept's own Sections.
func TestSimulateKnownEarlyValues(t *testing.T) {
	steps := Simulate(0.2, 3)
	want := []struct {
		arm      int
		reward   int
		explored bool
	}{
		{0, 1, false}, // t=1: forced pull of A, wins
		{1, 0, false}, // t=2: forced pull of B, loses
		{2, 1, false}, // t=3: forced pull of C, wins
	}
	for i, w := range want {
		s := steps[i]
		if s.Arm != w.arm || s.Reward != w.reward || s.Explored != w.explored {
			t.Errorf("step %d = {arm:%d reward:%d explored:%v}, want {arm:%d reward:%d explored:%v}",
				i, s.Arm, s.Reward, s.Explored, w.arm, w.reward, w.explored)
		}
	}
	if got := steps[2].CumRegret; math.Abs(got-0.45) > 1e-9 {
		t.Errorf("CumRegret after 3 forced pulls = %v, want 0.45 (0.30 for A + 0 for B + 0.15 for C)", got)
	}
}

// TestSimulateExplorePullAtStepFifteen is pinned to the same sequence:
// pull 15 at epsilon=0.2 is the first exploring pull, and it happens to
// land on B (the arm epsilon-greedy would otherwise have starved after B's
// unlucky first pull -- see TestEpsilonZeroNeverRevisitsUnluckyArm).
func TestSimulateExplorePullAtStepFifteen(t *testing.T) {
	steps := Simulate(0.2, 15)
	s := steps[14]
	if !s.Explored {
		t.Errorf("step 15: Explored = false, want true")
	}
	if s.Arm != 1 {
		t.Errorf("step 15: Arm = %d, want 1 (B)", s.Arm)
	}
}

// TestEpsilonZeroNeverRevisitsUnluckyArm demonstrates the greedy trap: B's
// only forced pull (occurrence 0) happens to lose (ArmOutcome(1,0) is
// false), so its estimate is pinned at 0. With epsilon=0, exploreDecision
// never fires, so argmaxEstimate never picks an arm whose estimate isn't
// currently the highest -- and since every other arm's estimate is >= 0,
// B (estimate exactly 0, tie broken toward lower indices anyway) can never
// win the argmax again. B's count should stay at exactly 1 forever, no
// matter how many pulls follow.
func TestEpsilonZeroNeverRevisitsUnluckyArm(t *testing.T) {
	if ArmOutcome(1, 0) {
		t.Fatal("test assumption violated: B's forced first pull is expected to lose")
	}
	for _, n := range []int{3, 10, 60, 300} {
		steps := Simulate(0, n)
		if got := steps[len(steps)-1].Counts[1]; got != 1 {
			t.Errorf("epsilon=0, n=%d: B pulled %d times, want exactly 1 (never revisited)", n, got)
		}
	}
}

// TestEpsilonPointTwoEventuallyFindsBestArm is pinned to the same
// deterministic sequence over a long horizon: despite B's unlucky start,
// enough exploring pulls at epsilon=0.2 eventually let its true 0.50 rate
// show through the noise, and by pull 1000 it dominates the pull counts.
func TestEpsilonPointTwoEventuallyFindsBestArm(t *testing.T) {
	steps := Simulate(0.2, 1000)
	last := steps[len(steps)-1]
	want := [3]int{91, 758, 151}
	if last.Counts != want {
		t.Errorf("Simulate(0.2,1000) final Counts = %v, want %v", last.Counts, want)
	}
	if last.Counts[1] <= last.Counts[0] || last.Counts[1] <= last.Counts[2] {
		t.Errorf("B (the true best arm) should have the highest pull count by n=1000, got %v", last.Counts)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("multi-armed-bandit")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
