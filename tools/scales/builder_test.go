package scales

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

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

// it duplicates the content folder root as a section
//
//	docs    docs
//	|-a     |-a
//	|-b  => |-b
//	        |-docs
//	          |-a
//	          |-b
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

// TestBuilder_docsSite performs a series of builds with the docs site which
// its pages are de-duplicated 2*i times each run to check if the memory
// consumption scales sublinear.
func TestBuilder_docsSiteAllocationScaling(t *testing.T) {
	tmp := t.TempDir()
	fmt.Println(tmp)
	err := copy.Dir(tmp, "../../docs")
	if err != nil {
		t.Errorf("prep, copying docs site into the test directory to recursively duplicate its contents: %v", err)
	}
	f, a, err := Allocations(20, func() (float64, error) {
		return doubleUpContentFolder(tmp)
	}, func() error {
		return builder.Build(builder.Args{Src: tmp, Dst: t.TempDir(), Domain: "/"})
	})
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	} else if f == NonSublinear {
		fmt.Printf("x=%s\n", toPythonArraySyntax(a.Sizes))
		fmt.Printf("y=%s\n", toPythonArraySyntax(a.Allocs))
		t.Fatal("assert memory allocation can't scale superlinear")
	}
}
