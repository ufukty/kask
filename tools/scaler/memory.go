package main

import (
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"go.ufukty.com/kask/internal/builder"
	"go.ufukty.com/kask/internal/builder/copy"
)

func mkTempDir() (string, error) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-scales-*")
	if err != nil {
		return "", fmt.Errorf("os.MkdirTemp: %w", err)
	}
	return tmp, nil
}

func rmTempDir(tmp string) {
	fmt.Printf("deleting: %s\n", tmp)
	err := os.RemoveAll(tmp)
	if err != nil {
		panic(err)
	}
}

type allocations struct {
	Sizes, TotalAllocs, Sys []uint64
}

func prepare(tmp, docssite string, step int) error {
	if step == 0 {
		err := copy.Dir(tmp, docssite)
		if err != nil {
			return fmt.Errorf("initial copying of docs contents: %w", err)
		}
		return nil
	} else {
		for j := range int(math.Pow(1.6, float64(step))) {
			err := copy.Dir(filepath.Join(tmp, fmt.Sprintf("new-section-%d-%d", step, j)), docssite)
			if err != nil {
				return fmt.Errorf("copying the docs site: %w", err)
			}
		}
		return nil
	}
}

func dirSize(path string) (uint64, error) {
	var size uint64
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("filepath.WalkDir: %w", err)
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return fmt.Errorf("d.Info: %w", err)
			}
			size += uint64(info.Size())
		}
		return nil
	})
	return size, err
}

func measure() (uint64, uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc, m.Sys
}

func invoke(tmp string) error {
	dst, err := mkTempDir()
	if err != nil {
		return fmt.Errorf("creating a directory in temp for output: %w", err)
	}
	args := builder.Args{Src: tmp, Dst: dst, Domain: "/"}
	err = builder.Build(args)
	if err != nil {
		return fmt.Errorf("builder.Build: %w", err)
	}
	return nil
}

type factor string

const (
	NonSublinear factor = "non-sublinear" // superlinear or linear
	Sublinear    factor = "sublinear"
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
func factorize(ys, xs []uint64) (factor, error) {
	if len(ys) != len(xs) {
		return "", fmt.Errorf("expected same number of x and y values")
	}
	dy, dx := ys[len(ys)-1]-ys[0], xs[len(xs)-1]-xs[0]
	if dx == 0 {
		return "", fmt.Errorf("constant scaling (impossible, check your code)")
	}
	m := dy / dx
	t := 0
	for i := 1; i+1 < len(ys); i++ {
		dy, dx := ys[i]-ys[0], xs[i]-xs[0]
		edy := dx * m // expected dy
		t += int(dy) - int(edy)
	}
	if t > 0 {
		return Sublinear, nil
	} else {
		return NonSublinear, nil
	}
}

type args struct {
	docspath string
}

func Main() error {
	args := args{}
	flag.StringVar(&args.docspath, "path", "", "path to the docs site")
	flag.Parse()
	if args.docspath == "" {
		return fmt.Errorf("checking -path: missing arg")
	}

	tmp, err := mkTempDir()
	if err != nil {
		return fmt.Errorf("creating the test directory: %w", err)
	}
	defer rmTempDir(tmp)
	fmt.Println("Using:", tmp)
	fmt.Printf("%13s %13s %13s\n", "Input size", "TotalAlloc", "Sys")

	a := allocations{
		Sizes:       []uint64{},
		TotalAllocs: []uint64{},
		Sys:         []uint64{},
	}
	for i := range 10 {
		err := prepare(tmp, args.docspath, i)
		if err != nil {
			return fmt.Errorf("preparing content directory: %w", err)
		}
		size, err := dirSize(tmp)
		if err != nil {
			return fmt.Errorf("sizing content directory: %w", err)
		}
		taBefore, _ := measure()
		if err := invoke(tmp); err != nil {
			return fmt.Errorf("invoking at step %d: %w", i, err)
		}
		taAfter, sysAfter := measure()
		taDelta := taAfter - taBefore
		a.Sizes = append(a.Sizes, size)
		a.TotalAllocs = append(a.TotalAllocs, taDelta)
		a.Sys = append(a.Sys, sysAfter)
		fmt.Printf("%10.2f MB %10.2f MB %10.2f MB\n",
			float64(size)/1024/1024, float64(taDelta)/1024/1024, float64(sysAfter)/1024/1024)
	}

	sl, err := factorize(a.Sys, a.Sizes)
	if err != nil {
		return fmt.Errorf("factorizing: %w", err)
	}

	if sl != Sublinear {
		return fmt.Errorf("not sublinear")
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
