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

func dfs(n []*Node, f func([]*Node)) {
	f(n)
	for _, c := range n[len(n)-1].Children {
		dfs(append(slices.Clone(n), c), f)
	}
}

func each(ns []*Node, f func(*Node) string) []string {
	ss := make([]string, 0, len(ns))
	for _, n := range ns {
		ss = append(ss, f(n))
	}
	return ss
}

func buildTestSite(path string) (*builder, string) {
	tmp, err := os.MkdirTemp(os.TempDir(), "kask-test-build-*")
	if err != nil {
		panic(fmt.Errorf("buildTestSite: os.MkdirTemp: %w", err))
	}
	b := newBuilder(Args{Src: path, Dst: tmp, Dev: true, Verbose: false})
	err = b.Build()
	if err != nil {
		panic(fmt.Errorf("buildTestSite: b.Build: %w", err))
	}
	return b, tmp
}

func TestBuilder_renderWebPages(t *testing.T) {
	_, tmp := buildTestSite("testdata/web")
	f, err := os.ReadFile(filepath.Join(tmp, "index.html"))
	if err != nil {
		t.Errorf("prep, read file: %v", err)
	}
	if !strings.Contains(string(f), "<h1>Title</h1>") {
		t.Error("expected markdown content not found.")
	}
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

func ExampleBuilder_strippedOrderingHrefs() {
	b, _ := buildTestSite("testdata/stripped-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Href) })
	// Output:
	// /
	// /career.html
	// /docs.html
	// /products.html
	// /about/
	// /contact/
}

func ExampleBuilder_strippedOrderingTitles() {
	b, _ := buildTestSite("testdata/stripped-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Title) })
	// Output:
	// Website Title
	// Career Title
	// Docs Title
	// Products Title
	// About Title
	// Contact Title
}

func ExampleBuilder_strippedOrderingBreadcrumbs() {
	b, _ := buildTestSite("testdata/stripped-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) {
		bs := each(n, func(n *Node) string { return n.Title })
		fmt.Println(strings.Join(bs, " / "))
	})
	// Output:
	// Website Title
	// Website Title / Career Title
	// Website Title / Docs Title
	// Website Title / Products Title
	// Website Title / About Title
	// Website Title / Contact Title
}

func ExampleBuilder_preservedOrderingHrefs() {
	b, _ := buildTestSite("testdata/preserved-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Href) })
	// Output:
	// /
	// /1.career.html
	// /2.docs.html
	// /3.products.html
	// /1.about
	// /2.contact
}

func ExampleBuilder_preservedOrderingTitles() {
	b, _ := buildTestSite("testdata/preserved-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Title) })
	// Output:
	// Website Title
	// Career Title
	// Docs Title
	// Products Title
	// About Title
	// Contact Title
}

func ExampleBuilder_preservedOrderingBreadcrumbs() {
	b, _ := buildTestSite("testdata/preserved-ordering")
	dfs([]*Node{b.root3}, func(n []*Node) {
		bs := each(n, func(n *Node) string { return n.Title })
		fmt.Println(strings.Join(bs, " / "))
	})
	// Output:
	// Website Title
	// Website Title / Career Title
	// Website Title / Docs Title
	// Website Title / Products Title
	// Website Title / About Title
	// Website Title / Contact Title
}

func TestBuilder_cssSplitting(t *testing.T) {
	_, tmp := buildTestSite("testdata/css-splitting")
	fmt.Println(tmp)

	t.Run("linking the scorrect stylesheets", func(t *testing.T) {
		f, err := os.ReadFile(filepath.Join(tmp, "a/index.html"))
		if err != nil {
			t.Errorf("prep, read file: %v", err)
		}
		content := string(f)
		if strings.Contains(content, `"/styles.css"`) {
			t.Error("page of subsection should NOT link the root section's non-propagated styles")
		}
		if !strings.Contains(content, `"/styles.propagate.css"`) {
			t.Error("page of subsection should link the root section's propagated styles")
		}
		if !strings.Contains(content, `"/a/styles.propagate.css"`) {
			t.Error("page should link its section's propagated styles")
		}
		if !strings.Contains(content, `"/a/styles.css"`) {
			t.Error("page should link its section's non-propagated styles")
		}
	})

	tcs := map[string]string{
		"/styles.css":             ":root {}",
		"/styles.propagate.css":   ":root.propagated {}",
		"/a/styles.css":           ".a {}",
		"/a/styles.propagate.css": ".a.propagated {}",
	}

	for ss, expected := range tcs {
		t.Run(fmt.Sprintf("stylesheet contents/%s", strings.ReplaceAll(ss, "/", "\\")), func(t *testing.T) {
			f, err := os.ReadFile(filepath.Join(tmp, ss))
			if err != nil {
				t.Errorf("prep, read stylesheet: %v", err)
			}
			if !strings.Contains(string(f), "") {
				t.Errorf("assert, not found: %q", expected)
			}
		})
	}
}

func ExampleBuilder_titles() {
	b, _ := buildTestSite("testdata/titles")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Title) })
	// Output:
	// .
	// An Anonymous Web Page
	// A beautiful tmpl file with title tag
	// An Anonymous Markdown Page
	// A beatiful markdown file with title
}

func ExampleBuilder_metaTitle() {
	b, _ := buildTestSite("testdata/meta-title")
	dfs([]*Node{b.root3}, func(n []*Node) { fmt.Println(n[len(n)-1].Title) })
	// Output:
	// My Beautiful Site
	// Page
	// My Beautiful Section
	// Page
}

func assertfile(t *testing.T, tmp, path string) {
	t.Run(strings.ReplaceAll(path, "/", "\\"), func(t *testing.T) {
		if !check(tmp, path) {
			t.Errorf("assert, file not found: %s", path)
		}
	})
}

func TestBuilder_assets(t *testing.T) {
	_, tmp := buildTestSite("testdata/assets")
	tcs := []string{".assets/sample.txt", "section/.assets/subsample.txt"}
	for _, tc := range tcs {
		assertfile(t, tmp, tc)
	}
}
