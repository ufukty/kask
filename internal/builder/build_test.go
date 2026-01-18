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
			"products/1-touch-pro.html",
			"docs",
			"docs/.assets",
			"docs/styles.propagate.css",
			"docs/index.html", // README.md
			"docs/101 tutorials/getting-started.html", // deep levels
			"docs/101 tutorials/install.html",
			"docs/101 tutorials/contribute.html",
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
	// /products/1-touch-pro.html
	// /docs/
	// /docs/birdseed.html
	// /docs/download.html
	// /docs/magnet.html
	//
	// /docs/101%20tutorials/getting-started.html
	// /docs/101%20tutorials/install.html
	// /docs/101%20tutorials/contribute.html
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
	// Acme/ACME Products/1-Touch Pro
	// Acme/Docs
	// Acme/Docs/ACME Bird Seed
	// Acme/Docs/Download
	// Acme/Docs/ACME Magnet
	// Acme/Docs/101 Tutorials
	// Acme/Docs/101 Tutorials/Getting Started
	// Acme/Docs/101 Tutorials/How to install
	// Acme/Docs/101 Tutorials/How to contribute
}

func TestBuilder_propagated(t *testing.T) {
	tcs := []string{"web", "mixed", "markdown"}
	for _, tc := range tcs {
		t.Run(tc, func(t *testing.T) {
			tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
			if err != nil {
				t.Errorf("os.MkdirTemp: %v", err)
			}
			fmt.Println("temp folder:", tmp)

			a := Args{
				Src:     filepath.Join("testdata/propagated", tc),
				Dst:     tmp,
				Dev:     true,
				Verbose: true,
			}
			err = Build(a)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
		})
	}
}
