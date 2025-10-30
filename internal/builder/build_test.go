package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// returns the DFS forest
func forest(n *Node) []*Node {
	f := []*Node{n}
	for _, c := range n.Children {
		f = append(f, forest(c)...)
	}
	return f
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

func TestBuild(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		t.Errorf("os.MkdirTemp: %v", err)
	}
	fmt.Println("temp folder:", tmp)

	b := newBuilder(Args{
		Domain:  "http://localhost:8080",
		Dev:     false,
		Src:     "testdata/acme",
		Dst:     tmp,
		Verbose: true,
	})

	t.Run("build", func(t *testing.T) {
		err = b.Build()
		if err != nil {
			t.Fatal(fmt.Errorf("act, Build: %w", err))
		}
	})

	t.Run("stat files", func(t *testing.T) {
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
}

func ExampleBuild_sitemap() {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		panic(fmt.Errorf("os.MkdirTemp: %w", err))
	}

	b := newBuilder(Args{
		Domain:  "http://localhost:8080",
		Dev:     false,
		Src:     "testdata/acme",
		Dst:     tmp,
		Verbose: false,
	})

	err = b.Build()
	if err != nil {
		panic(fmt.Errorf("act, Build: %w", err))
	}

	for _, node := range forest(b.root3) {
		fmt.Println(node.Href)
	}
	// Output:
	// /
	// /career/
	// /products/
	// /docs/
	// /docs/birdseed.html
	// /docs/download.html
	// /docs/magnet.html
	//
	// /docs/tutorials/getting-started.html
}

func ExampleBuild_breadcrumbs() {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		panic(fmt.Errorf("os.MkdirTemp: %w", err))
	}

	b := newBuilder(Args{
		Domain:  "http://localhost:8080",
		Dev:     false,
		Src:     "testdata/acme",
		Dst:     tmp,
		Verbose: false,
	})

	err = b.Build()
	if err != nil {
		panic(fmt.Errorf("act, Build: %w", err))
	}

	for _, node := range forest(b.root3) {
		path := []string{}
		for _, p := range ancestry(node) {
			path = append(path, p.Title)
		}
		fmt.Println(strings.Join(path, "/"))
	}
	// Output:
	// Acme
	// Acme/Careers at ACME
	// Acme/ACME Products
	// Acme/Docs
	// Acme/Docs/ACME Bird Seed
	// Acme/Docs/Download
	// Acme/Docs/ACME Magnet
	// Acme/Docs/tutorials
	// Acme/Docs/tutorials/Getting Started
}
