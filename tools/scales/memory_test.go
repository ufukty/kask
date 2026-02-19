package main

import "testing"

func TestFactorize_sublinear(t *testing.T) {
	ys := []uint64{0, 1, 2, 3, 4, 5}
	xs := []uint64{1, 10, 100, 1000, 10000, 100000}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != Sublinear {
		t.Errorf("assert, expected %s got %s", Sublinear, f)
	}
}

func TestFactorize_linear(t *testing.T) {
	ys := []uint64{1, 2, 3, 4, 5, 6}
	xs := []uint64{1, 2, 3, 4, 5, 6}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != NonSublinear {
		t.Errorf("assert, expected %s got %s", NonSublinear, f)
	}
}

func TestFactorize_superlinear(t *testing.T) {
	ys := []uint64{1, 2, 4, 8, 16, 32}
	xs := []uint64{1, 2, 3, 4, 5, 6}
	f, err := factorize(ys, xs)
	if err != nil {
		t.Error("act, unexpected error: %w", err)
	} else if f != NonSublinear {
		t.Errorf("assert, expected %s got %s", NonSublinear, f)
	}
}
