package scales

import (
	"fmt"
	"runtime"
)

type Factor string

const (
	NonSublinear Factor = "non-sublinear" // superlinear or linear
	Sublinear    Factor = "sublinear"
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
	if t > 1e-5 {
		return Sublinear, nil
	} else {
		return NonSublinear, nil
	}
}

func totalAllocs() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}

type allocations struct {
	Sizes, Allocs []float64
}

func Allocations(steps int, prep func() (float64, error), perform func() error) (Factor, allocations, error) {
	a := allocations{
		Sizes:  []float64{},
		Allocs: []float64{},
	}
	for i := range steps {
		size, err := prep()
		if err != nil {
			return "", allocations{}, fmt.Errorf("prep(%d): %w", i, err)
		}
		before := totalAllocs()
		if err := perform(); err != nil {
			return "", allocations{}, fmt.Errorf("perform(%d): %w", i, err)
		}
		delta := float64(totalAllocs() - before)
		a.Sizes = append(a.Sizes, size)
		a.Allocs = append(a.Allocs, delta)
		fmt.Printf("Input / Total alloc: %.2f MB => %.2f MB\n", size/1024/1024, delta/1024/1024)
	}
	sl, err := factorize(a.Allocs, a.Sizes)
	if err != nil {
		return "", allocations{}, fmt.Errorf("factorizing: %w", err)
	}
	return sl, a, nil
}
