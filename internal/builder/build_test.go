package builder

import (
	"fmt"
	"os"
	"testing"

	"github.com/ufukty/kask/internal/builder/markdown"
)

var defaults = Args{
	Domain: "http://localhost:8080",
	Dev:    false,
}

func TestBuild(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		t.Fatal(fmt.Errorf("os.MkdirTemp: %w", err))
	}
	fmt.Println("temp folder:", tmp)

	err = Build(tmp, "testdata/acme", defaults)
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
		args:          defaults,
		assets:        []string{},
		stylesheets:   map[string]string{},
		pagesMarkdown: map[string]*markdown.Page{},
	}

	err = b.Build(tmp, "testdata/acme")

	if err != nil {
		t.Fatal(fmt.Errorf("act, Build: %w", err))
	}
}
