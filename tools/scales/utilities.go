package scales

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.ufukty.com/kask/internal/builder"
	"go.ufukty.com/kask/internal/builder/copy"
)

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// it linearly increases the content folder root as a section (1x, 2x, 3x...)
//
//	docs    docs            docs
//	|-a     |-a             |-a
//	|-b  => |-b          => |-b
//	        |-docs-again    |-docs-again
//	          |-a             |-a
//	          |-b             |-b
//	                          |-docs-again
//	                            |-a
//	                            |-b
func doubleUpContentFolder(path string) (float64, error) {
	copy.Dir(filepath.Join(path, "docs-again"), path)
	s, err := dirSize(path)
	if err != nil {
		return -1, fmt.Errorf("dirSize: %w", err)
	}
	return float64(s), nil
}

func toPythonArraySyntax(vs []float64) string {
	ss := []string{}
	for _, v := range vs {
		ss = append(ss, strconv.FormatFloat(v, 'f', 2, 64))
	}
	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}

func mkTempDir() (string, error) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-scales-*")
	if err != nil {
		return "", fmt.Errorf("os.MkdirTemp: %w", err)
	}
	return tmp, nil
}

func main(docssite string) error {
	tmp, err := mkTempDir()
	if err != nil {
		return fmt.Errorf("creating the test directory: %w", err)
	}
	fmt.Println("Using:", tmp)

	i := 0
	f, a, err := Allocations(20, func() (float64, error) {
		err = copy.Dir(filepath.Join(tmp, fmt.Sprintf("docs-%d", i)), docssite)
		if err != nil {
			return -1, fmt.Errorf("copying the docs site: %w", err)
		}
		return 0, nil
	}, func() error {
		args := builder.Args{Src: tmp, Dst: mkTempDir(), Domain: "/"}
		err := builder.Build(args)
		if err != nil {
			return fmt.Errorf("builder.Build: %w", err)
		}
		return nil
	})

	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	} else if f == NonSublinear {
		fmt.Printf("x=%s\n", toPythonArraySyntax(a.Sizes))
		fmt.Printf("y=%s\n", toPythonArraySyntax(a.Allocs))
		t.Fatal("assert memory allocation can't scale superlinear")
	}
}
