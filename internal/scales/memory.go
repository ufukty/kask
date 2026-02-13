package scales

import (
	"fmt"
	"math"
	"runtime"
)

type Factor string

const (
	Superlinear Factor = "superlinear"
	Linear      Factor = "linear"
	Sublinear   Factor = "sublinear"
)

// This just assumes that the average data point elevation over the line
// which connects the first and last point would be positive for values
// increase sublinearly. Expected to fail at small sets.
//
//	^
//	|                           x
//	|                 x
//	|           x
//	|       x
//	|    x
//	|  x
//	| x
//	+----------------------------->
func factorize(ys, xs []float64) (Factor, error) {
	if len(ys) != len(xs) {
		return "", fmt.Errorf("expected same number of x and y values")
	}
	dy, dx := ys[len(ys)-1]-ys[0], xs[len(xs)-1]-xs[0]
	if dx == 0 {
		return "", fmt.Errorf("constant scaling (impossible, check your code)")
	}
	m := dy / dx
	t := 0.0
	for i := 1; i+1 < len(ys); i++ {
		dy, dx := ys[i]-ys[0], xs[i]-xs[0]
		edy := dx * m // expected dy
		t += dy - edy
	}
	if t == 0 {
		return Linear, nil
	} else if t > 0 {
		return Sublinear, nil
	} else {
		return Superlinear, nil
	}
}

func totalAllocs() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}

func Allocations(maxsize int, prep, perform func(size int) error) (Factor, error) {
	if e := math.Log2(float64(maxsize)); (e - float64(int(e))) > 0.0 {
		return "", fmt.Errorf("max size should be a power of 2")
	}
	inputSizes := []float64{}
	allocs := []float64{}
	for i := 1; i <= maxsize; i *= 2 {
		if err := prep(i); err != nil {
			return "", fmt.Errorf("prep(%d): %w", i, err)
		}
		before := totalAllocs()
		if err := perform(i); err != nil {
			return "", fmt.Errorf("prep(%d): %w", i, err)
		}
		delta := totalAllocs() - before
		inputSizes = append(inputSizes, float64(i))
		allocs = append(allocs, float64(delta))
		deltaMb := delta / 1024 / 1024
		fmt.Printf("Input size: [log2(%3dx) = %.2f], Total Alloc: [log2(%3d MB) = %.2f]\n",
			i, math.Log2(float64(i)),
			deltaMb, math.Log2(max(float64(deltaMb), 1e-9)),
		)
	}
	sl, err := factorize(allocs, inputSizes)
	if err != nil {
		return "", fmt.Errorf("factorizing: %w", err)
	}
	return sl, nil
}
