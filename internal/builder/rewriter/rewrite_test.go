package rewriter

import (
	"cmp"
	"iter"
	"maps"
	"slices"
	"strings"
	"testing"
)

// /a and /a/b/c is not visitable
var rw = Rewriter{
	links: map[string]string{
		"README.md":         "/",
		"page.md":           "/page.html",
		"a/b/index.tmpl":    "/a/b/",
		"a/b/page.tmpl":     "/a/b/page.html",
		"a/b/c/page.md":     "/a/b/c/page.html",
		"a/b/c/d/README.md": "/a/b/c/d/",

		// visitable directories (src => url)
		".":       "/",
		"a/b":     "/a/b/",
		"a/b/c/d": "/a/b/c/d/",
	},
}

func testname(tn string) string {
	tn = strings.ReplaceAll(tn, "/", "\\")
	tn = strings.ReplaceAll(tn, "%20", " ")
	return tn
}

func sorted[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range slices.Sorted(maps.Keys(m)) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

func TestRewrite_linksToParents(t *testing.T) {
	tcs := map[string]string{
		"../../":          "/",
		"../../page.md":   "/page.html",
		"./../../":        "/",
		"./../../page.md": "/page.html",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "a/b/page.tmpl")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksToSubdirs(t *testing.T) {
	tcs := map[string]string{
		"a/b":               "/a/b/",
		"a/b/c/d":           "/a/b/c/d/",
		"a/b/c/d/README.md": "/a/b/c/d/",
		"a/b/page.tmpl":     "/a/b/page.html",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "page.tmpl")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksWithReduntantSegments(t *testing.T) {
	tcs := map[string]string{
		"a/../a/b":               "/a/b/",
		"a/../a/b/c/d":           "/a/b/c/d/",
		"a/../a/b/c/d/README.md": "/a/b/c/d/",
		"a/../a/b/page.tmpl":     "/a/b/page.html",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "page.tmpl")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_absoluteLinks(t *testing.T) {
	tcs := map[string]string{
		"/":               "/",
		"/README.md":      "/",
		"/page.md":        "/page.html",
		"/a/b":            "/a/b/",
		"/a/b/index.tmpl": "/a/b/",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "page.tmpl")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_internalAnchorLinks(t *testing.T) {
	tcs := map[string]string{
		"#":      "#",
		"#title": "#title",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "page.md")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_externalAnchorLinks(t *testing.T) {
	tcs := map[string]string{
		"/#title":                 "/#title",
		"/README.md#title":        "/#title",
		"a/b/#title":              "/a/b/#title",
		"a/b/c/d/#title":          "/a/b/c/d/#title",
		"a/b/c/d/README.md#title": "/a/b/c/d/#title",
		"a/b/c/page.md#title":     "/a/b/c/page.html#title",
		"a/b/index.tmpl#title":    "/a/b/#title",
		"a/b/page.tmpl#title":     "/a/b/page.html#title",
		"README.md#title":         "/#title",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, "page.md")
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksToUnvisitableDirs(t *testing.T) {
	tcs := []string{
		"..",      // /a/
		"../",     // /a/
		"../b/c",  // /a/b/c/
		"../b/c/", // /a/b/c/
		"./../",   // /a/
		"./c",     // /a/b/c/
		"./c/",    // /a/b/c/
		"c",       // /a/b/c/
		"c/",      // /a/b/c/
	}
	for _, input := range tcs {
		t.Run(testname(input), func(t *testing.T) {
			_, err := rw.Rewrite(input, "a/b/page.tmpl")
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}

func TestRewrite_linksToUnexistingNodes(t *testing.T) {
	tcs := []string{
		"../../..",     // /../
		"../x",         // /x/
		"../../x.tmpl", // /x.html
	}
	for _, input := range tcs {
		t.Run(testname(input), func(t *testing.T) {
			_, err := rw.Rewrite(input, "a/b/page.tmpl")
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}
