package cosinesimilarity

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

func TestDotKnownValues(t *testing.T) {
	if got := Dot(4, 2, 2, 1); !almostEqual(got, 10, epsilon) {
		t.Errorf("Dot(4,2,2,1) = %v, want 10", got)
	}
	if got := Dot(1, 0, 0, 1); !almostEqual(got, 0, epsilon) {
		t.Errorf("Dot(1,0,0,1) = %v, want 0", got)
	}
}

func TestMagnitudeKnownValues(t *testing.T) {
	if got := Magnitude(3, 4); !almostEqual(got, 5, epsilon) {
		t.Errorf("Magnitude(3,4) = %v, want 5", got)
	}
	if got := Magnitude(0, 0); !almostEqual(got, 0, epsilon) {
		t.Errorf("Magnitude(0,0) = %v, want 0", got)
	}
}

func TestCosineSimilaritySameDirectionIsOne(t *testing.T) {
	// v = (2,1) is exactly half of u = (4,2): same direction, different
	// length -- cosine similarity should be exactly 1 regardless.
	if got := CosineSimilarity(4, 2, 2, 1); !almostEqual(got, 1, 1e-9) {
		t.Errorf("CosineSimilarity(4,2,2,1) = %v, want 1", got)
	}
	// Scaling u further shouldn't change it either.
	if got := CosineSimilarity(4, 2, 40, 20); !almostEqual(got, 1, 1e-9) {
		t.Errorf("CosineSimilarity(4,2,40,20) = %v, want 1", got)
	}
}

func TestCosineSimilarityOrthogonalIsZero(t *testing.T) {
	// u=(4,2), w=(2,-4): Dot = 8-8 = 0 -> orthogonal.
	if got := CosineSimilarity(4, 2, 2, -4); !almostEqual(got, 0, 1e-9) {
		t.Errorf("CosineSimilarity(4,2,2,-4) = %v, want 0", got)
	}
}

func TestCosineSimilarityOppositeIsNegativeOne(t *testing.T) {
	if got := CosineSimilarity(4, 2, -4, -2); !almostEqual(got, -1, 1e-9) {
		t.Errorf("CosineSimilarity(4,2,-4,-2) = %v, want -1", got)
	}
}

func TestCosineSimilarityWorkedExample(t *testing.T) {
	// u=(4,2), C=(1,5): Dot=4+10=14, |u|=sqrt(20), |C|=sqrt(26).
	got := CosineSimilarity(4, 2, 1, 5)
	want := 14 / (math.Sqrt(20) * math.Sqrt(26))
	if !almostEqual(got, want, 1e-9) {
		t.Errorf("CosineSimilarity(4,2,1,5) = %v, want %v (~0.614)", got, want)
	}
	if !almostEqual(want, 0.6139406135, 1e-6) {
		t.Errorf("sanity: worked-example cosine ~0.614, got %v", want)
	}
}

func TestCosineSimilarityZeroVectorReturnsZero(t *testing.T) {
	if got := CosineSimilarity(0, 0, 3, 4); got != 0 {
		t.Errorf("CosineSimilarity with zero vector = %v, want 0", got)
	}
}

func TestAngleDegreesKnownValues(t *testing.T) {
	if got := AngleDegrees(1, 0, 0, 1); !almostEqual(got, 90, 1e-6) {
		t.Errorf("AngleDegrees(1,0,0,1) = %v, want 90", got)
	}
	if got := AngleDegrees(4, 2, 2, 1); !almostEqual(got, 0, 1e-4) {
		t.Errorf("AngleDegrees(4,2,2,1) = %v, want ~0 (same direction)", got)
	}
	if got := AngleDegrees(4, 2, -4, -2); !almostEqual(got, 180, 1e-4) {
		t.Errorf("AngleDegrees(4,2,-4,-2) = %v, want ~180 (opposite)", got)
	}
}

func TestEuclideanDistanceContrastsWithCosine(t *testing.T) {
	// u=(4,2), v=(2,1): cosine similarity is 1 (identical direction) but
	// the points themselves are not the same point -- distance is nonzero.
	dist := EuclideanDistance(4, 2, 2, 1)
	want := math.Sqrt(5)
	if !almostEqual(dist, want, 1e-9) {
		t.Errorf("EuclideanDistance(4,2,2,1) = %v, want %v", dist, want)
	}
	if dist <= 0 {
		t.Error("expected nonzero distance even though cosine similarity is 1")
	}
}

func TestRenderProducesSVG(t *testing.T) {
	c, ok := concept.Get("cosine-similarity")
	if !ok {
		t.Fatal("concept not registered")
	}
	out := c.Render(c.Defaults())
	if !strings.HasPrefix(out, "<svg") || !strings.HasSuffix(out, "</svg>") {
		t.Fatal("render did not produce a well-formed svg")
	}
}
