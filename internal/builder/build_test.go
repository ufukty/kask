package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/ufukty/kask/internal/builder/markdown"
)

// returns the DFS forest
func forest(n *Node) []*Node {
	f := []*Node{n}
	for _, c := range n.Children {
		f = append(f, forest(c)...)
	}
	return f
}

func fmap[S, T any](ss []S, m func(S) T) []T {
	ts := make([]T, len(ss))
	for i, s := range ss {
		ts[i] = m(s)
	}
	return ts
}

func check(tmp, path string) bool {
	_, err := os.Stat(filepath.Join(tmp, path))
	return err == nil
}

func ancestry(n *Node) []*Node {
	ancestry := []*Node{}
	for n := n; n != nil; n = n.Parent {
		ancestry = append(ancestry, n)
	}
	slices.Reverse(ancestry)
	return ancestry
}

func titlePaths(f []*Node) []string {
	ss := []string{}
	for _, n := range f {
		ss = append(ss, strings.Join(fmap(ancestry(n), func(n *Node) string { return n.Title }), "/"))
	}
	return ss
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
		expected := []string{
			"/.",                                   // "."
			"/career",                              // "./career",
			"/docs",                                // "./Docs",
			"/docs/birdseed.html",                  // "./Docs/ACME Bird Seed",
			"/docs/download.html",                  // "./Docs/Download",
			"/docs/magnet.html",                    // "./Docs/ACME Magnet",
			"/docs/tutorials/getting-started.html", // "./Docs/tutorials/Getting Started"
			"/products",                            // "./products",
		}

		got := fmap(forest(b.root3), func(n *Node) string { return n.Href })
		for _, e := range expected {
			t.Run("sitemap for "+strings.ReplaceAll(e, "/", "\\"), func(t *testing.T) {
				if !slices.Contains(got, e) {
					t.Errorf("not found")
				}
			})
		}
	})

	t.Run("breadcrumbs", func(t *testing.T) {
		expected := []string{
			"Acme",
			"Acme/career",
			"Acme/Docs",
			"Acme/Docs/ACME Bird Seed",
			"Acme/Docs/Download",
			"Acme/Docs/ACME Magnet",
			"Acme/Docs/tutorials",
			"Acme/Docs/tutorials/Getting Started",
			"Acme/products",
		}

		got := titlePaths(forest(b.root3))
		for _, e := range expected {
			t.Run("sitemap for "+strings.ReplaceAll(e, "/", "\\"), func(t *testing.T) {
				if !slices.Contains(got, e) {
					t.Errorf("not found")
				}
			})
		}
	})
}
