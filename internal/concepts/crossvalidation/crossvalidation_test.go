package crossvalidation

import (
	"math"
	"testing"
)

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestFullFitSlopeIntercept(t *testing.T) {
	slope := Slope(DataX, DataY)
	intercept := Intercept(DataX, DataY)
	if !approxEqual(slope, 2.2448, 0.001) {
		t.Errorf("full-fit slope = %.4f, want ~2.2448", slope)
	}
	if !approxEqual(intercept, -0.5067, 0.001) {
		t.Errorf("full-fit intercept = %.4f, want ~-0.5067", intercept)
	}
}

func TestKFoldIndicesEvenSplit(t *testing.T) {
	folds := KFoldIndices(10, 5)
	want := [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}}
	if len(folds) != len(want) {
		t.Fatalf("got %d folds, want %d", len(folds), len(want))
	}
	for i, f := range folds {
		if len(f) != len(want[i]) {
			t.Fatalf("fold %d = %v, want %v", i, f, want[i])
		}
		for j := range f {
			if f[j] != want[i][j] {
				t.Errorf("fold %d = %v, want %v", i, f, want[i])
			}
		}
	}
}

func TestKFoldIndicesUnevenSplit(t *testing.T) {
	folds := KFoldIndices(10, 3)
	sizes := []int{}
	total := 0
	for _, f := range folds {
		sizes = append(sizes, len(f))
		total += len(f)
	}
	if total != 10 {
		t.Fatalf("fold sizes %v sum to %d, want 10", sizes, total)
	}
	wantSizes := []int{3, 3, 4}
	for i, s := range sizes {
		if s != wantSizes[i] {
			t.Errorf("fold %d size = %d, want %d", i, s, wantSizes[i])
		}
	}
}

// TestFoldFitOutlierFold checks the fold (k=5, fold index 3) whose held-out
// points are {x=7, x=8} — the fold containing the deliberate outlier at
// x=8. Trained without ever seeing the outlier, the line predicts a clean
// ~16 at x=8 against the actual 24, so this fold's MSE should spike well
// above the other folds'.
func TestFoldFitOutlierFold(t *testing.T) {
	folds := KFoldIndices(10, 5)
	slope, intercept, mse := FoldFit(DataX, DataY, folds, 3)
	if !approxEqual(slope, 2.0, 0.01) {
		t.Errorf("fold 3 slope = %.4f, want ~2.0", slope)
	}
	if !approxEqual(intercept, 0.025, 0.01) {
		t.Errorf("fold 3 intercept = %.4f, want ~0.025", intercept)
	}
	if !approxEqual(mse, 31.82, 0.1) {
		t.Errorf("fold 3 MSE = %.4f, want ~31.82", mse)
	}
}

// TestFoldFitCleanFold checks a fold nowhere near the outlier (fold index
// 0, held-out {x=1, x=2}) has a far smaller MSE than the outlier's fold —
// the whole point being that a single split's verdict depends heavily on
// which fold you happened to land on.
func TestFoldFitCleanFold(t *testing.T) {
	folds := KFoldIndices(10, 5)
	_, _, mse := FoldFit(DataX, DataY, folds, 0)
	if !approxEqual(mse, 0.23, 0.05) {
		t.Errorf("fold 0 MSE = %.4f, want ~0.23", mse)
	}
}

func TestCrossValidateAverage(t *testing.T) {
	foldMSEs, avg := CrossValidate(DataX, DataY, 5)
	if len(foldMSEs) != 5 {
		t.Fatalf("got %d fold MSEs, want 5", len(foldMSEs))
	}
	if !approxEqual(avg, 10.469, 0.01) {
		t.Errorf("CV average MSE = %.4f, want ~10.469", avg)
	}
	// The outlier's fold must dominate the spread: its MSE should be more
	// than an order of magnitude above the smallest fold's.
	min, max := foldMSEs[0], foldMSEs[0]
	for _, m := range foldMSEs {
		if m < min {
			min = m
		}
		if m > max {
			max = m
		}
	}
	if max < 10*min {
		t.Errorf("fold MSEs %v don't show a large spread (min=%.2f max=%.2f)", foldMSEs, min, max)
	}
}

func TestMeanPredictMSEBasics(t *testing.T) {
	if !approxEqual(Mean([]float64{1, 2, 3}), 2, 1e-9) {
		t.Error("Mean([1,2,3]) should be 2")
	}
	if Predict(2, 1, 3) != 7 {
		t.Errorf("Predict(2,1,3) = %v, want 7", Predict(2, 1, 3))
	}
	// A perfect fit has zero MSE.
	xs, ys := []float64{0, 1, 2}, []float64{1, 3, 5}
	if mse := MSE(xs, ys, 2, 1); !approxEqual(mse, 0, 1e-9) {
		t.Errorf("MSE of a perfect fit = %v, want 0", mse)
	}
}

func TestRenderProducesSVG(t *testing.T) {
	svg := render(map[string]float64{"k": 5, "highlightFold": 3})
	if len(svg) == 0 {
		t.Fatal("render returned empty string")
	}
	if svg[:4] != "<svg" {
		t.Errorf("render output doesn't start with <svg: %q...", svg[:20])
	}
}
