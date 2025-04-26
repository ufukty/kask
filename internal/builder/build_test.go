package builder

import (
	"fmt"
	"os"
	"testing"

	"github.com/ufukty/kask/internal/builder/markdown"
)

func TestBuild(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		t.Fatal(fmt.Errorf("os.MkdirTemp: %w", err))
	}
	fmt.Println("temp folder:", tmp)

	err = Build(Args{
		Domain:  "http://localhost:8080",
		Dev:     false,
		Src:     "testdata/acme",
		Dst:     tmp,
		Verbose: true,
	})
	if err != nil {
		t.Fatal(fmt.Errorf("act, Build: %w", err))
	}
}

func TestBuilder(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		t.Fatal(fmt.Errorf("os.MkdirTemp: %w", err))
	}
	fmt.Println("temp folder:", tmp)

	b := builder{
		args: Args{
			Domain:  "http://localhost:8080",
			Dev:     false,
			Src:     "testdata/acme",
			Dst:     tmp,
			Verbose: true,
		},
		assets:        []string{},
		pagesMarkdown: map[string]*markdown.Page{},
		leaves:        map[pageref]*Node{},
	}

	err = b.Build()

	if err != nil {
		t.Fatal(fmt.Errorf("act, Build: %w", err))
	}
}
