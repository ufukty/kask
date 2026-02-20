package rewriter

import (
	"cmp"
	"iter"
	"maps"
	"slices"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/builder/paths"
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
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[string]string{
		"../../":               "/",
		"../../page.md":        "/page.html",
		"./../../":             "/",
		"./../../page.md":      "/page.html",
		"./../../page.md#home": "/page.html#home",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "page.tmpl", Dst: "page.html", Url: "/page.html"}
	tcs := map[string]string{
		"a/b":                "/a/b/",
		"a/b/c/d":            "/a/b/c/d/",
		"a/b/c/d/README.md":  "/a/b/c/d/",
		"a/b/page.tmpl":      "/a/b/page.html",
		"a/b/page.tmpl#home": "/a/b/page.html#home",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "page.tmpl", Dst: "page.html", Url: "/page.html"}
	tcs := map[string]string{
		"a/../a/b":                "/a/b/",
		"a/../a/b/c/d":            "/a/b/c/d/",
		"a/../a/b/c/d/README.md":  "/a/b/c/d/",
		"a/../a/b/page.tmpl":      "/a/b/page.html",
		"a/../a/b/page.tmpl#home": "/a/b/page.html#home",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "page.tmpl", Dst: "page.html", Url: "/page.html"}
	tcs := map[string]string{
		"/":                    "/",
		"/README.md":           "/",
		"/page.md":             "/page.html",
		"/a/b":                 "/a/b/",
		"/a/b/index.tmpl":      "/a/b/",
		"/a/b/index.tmpl#home": "/a/b/#home",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "page.tmpl", Dst: "page.html", Url: "/page.html"}
	tcs := map[string]string{
		"#":      "#",
		"#title": "#title",
	}
	for input, expected := range sorted(tcs) {
		t.Run(testname(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
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
			got, err := rw.Rewrite(input, linker)
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
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
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
			_, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}

func TestRewrite_linksToUnexistingNodes(t *testing.T) {
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := []string{
		"../../..",     // /../
		"../x",         // /x/
		"../../x.tmpl", // /x.html
	}
	for _, input := range tcs {
		t.Run(testname(input), func(t *testing.T) {
			_, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}
