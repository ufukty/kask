package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ufukty/kask/internal/builder/markdown"
)

func linearize(n *Node) string {
	s := n.Title
	for _, c := range n.Children {
		s += "\n| " + strings.ReplaceAll(linearize(c), "\n", "\n| ")
	}
	return s
}

func check(tmp, path string) bool {
	_, err := os.Stat(filepath.Join(tmp, path))
	return err == nil
}

func TestBuild(t *testing.T) {
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

	t.Run("building", func(t *testing.T) {
		err = b.Build()
		if err != nil {
			t.Fatal(fmt.Errorf("act, Build: %w", err))
		}
	})

	t.Run("stat", func(t *testing.T) {
		expected := []string{
			"index.html",
			"products",
			"docs",
			"docs/.assets",
			"docs/styles.propagate.css",
			"docs/index.html",                     // README.md
			"docs/tutorials/getting-started.html", // deep levels
		}

		for _, f := range expected {
			if !check(tmp, f) {
				t.Errorf("not found: %s", f)
			}
		}
	})

	t.Run("sitemap", func(t *testing.T) {
		const expected = `.
| career
| Docs
| | ACME Bird Seed
| | Download
| | ACME Magnet
| | tutorials
| | | Getting Started
| products`
		if got := linearize(b.root3); got != expected {
			t.Fatalf("expected:\n\n%s\n\ngot:\n%s", expected, got)
		}
	})
}
