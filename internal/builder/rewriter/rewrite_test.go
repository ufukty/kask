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
		".":                 "/",
		"README.md":         "/",
		"page.md":           "/page.html",
		"a/b/index.tmpl":    "/a/b/",
		"a/b/page.tmpl":     "/a/b/page.html",
		"a/b/c/page.md":     "/a/b/c/page.html",
		"a/b/c/d/README.md": "/a/b/c/d/",

		// visitable directories (src => url)
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
			got := rw.Rewrite(input, "a/b/page.tmpl")
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
			got := rw.Rewrite(input, "page.tmpl")
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
			got := rw.Rewrite(input, "page.tmpl")
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_externalInPageLinks(t *testing.T) {
	tcs := map[string]string{
		"/#title":             "/#title",
		"#title":              "/page.html#title",
		"a/b/#title":          "/a/b/#title",
		"a/b/c/d/#title":      "/a/b/c/d/#title",
		"a/b/c/page.md#title": "/a/b/c/page.html#title",
		"a/b/page.tmpl#title": "/a/b/page.html#title",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got := rw.Rewrite(input, "page.md")
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}
