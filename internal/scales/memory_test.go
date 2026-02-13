package scales

import "testing"

func TestFactorize_sublinear(t *testing.T) {
	ys := []float64{0, 1, 2, 3, 4, 5}
	xs := []float64{1, 10, 100, 1000, 10000, 100000}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != Sublinear {
		t.Errorf("assert, expected %s got %s", Sublinear, f)
	}
}

func TestFactorize_linear(t *testing.T) {
	ys := []float64{1, 2, 3, 4, 5, 6}
	xs := []float64{1, 2, 3, 4, 5, 6}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != Linear {
		t.Errorf("assert, expected %s got %s", Linear, f)
	}
}

func TestFactorize_superlinear(t *testing.T) {
	ys := []float64{1, 2, 4, 8, 16, 32}
	xs := []float64{1, 2, 3, 4, 5, 6}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != Superlinear {
		t.Errorf("assert, expected %s got %s", Superlinear, f)
	}
}
