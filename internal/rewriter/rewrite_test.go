package rewriter

import (
	"cmp"
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
	"testing"

	"go.ufukty.com/kask/internal/paths"
)

func tescape(tn string) string {
	tn = strings.ReplaceAll(tn, "/", "\\")
	tn = strings.ReplaceAll(tn, "%20", " ")
	return tn
}

func testname(a, b string) string {
	return fmt.Sprintf("%s=>%s", tescape(a), tescape(b))
}

// /a and /a/b/c is not visitable
func rewriter(domain string) *Rewriter {
	links := map[string]string{
		// pages
		"README.md":          domain + "",
		"page.md":            domain + "page.html",
		"a/b/index.tmpl":     domain + "a/b/",
		"a/b/page.tmpl":      domain + "a/b/page.html",
		"a/b/c /page.md":     domain + "a/b/c%20/page.html",
		"a/b/c /d/README.md": domain + "a/b/c%20/d/",
		// visitable directories (src => url)
		".":        domain + "",
		"a/b":      domain + "a/b/",
		"a/b/c /d": domain + "a/b/c%20/d/",
		// assets
		".assets/x.png":       domain + ".assets/x.png",
		".assets/x 2.png":     domain + ".assets/x 2.png",
		"a/.assets/x.png":     domain + "a/.assets/x.png",
		"a/.assets/x 2.png":   domain + "a/.assets/x 2.png",
		"a/b/.assets/x.png":   domain + "a/b/.assets/x.png",
		"a/b/.assets/x 2.png": domain + "a/b/.assets/x 2.png",
	}
	r := New(paths.Paths{Src: ".", Dst: ".", Url: domain})
	for src, url := range links {
		r.Bank(src, url)
	}
	return r
}

type tc struct {
	linker paths.Paths
	linked string
}

func sorted(m map[tc]string) iter.Seq2[tc, string] {
	m2 := map[string][]tc{}
	for k, o := range m {
		if _, ok := m2[o]; !ok {
			m2[o] = []tc{}
		}
		m2[o] = append(m2[o], k)
	}
	for k := range m2 {
		slices.SortFunc(m2[k], func(a, b tc) int {
			return cmp.Or(
				cmp.Compare(a.linker.Src, b.linker.Src),
				cmp.Compare(a.linked, b.linked),
			)
		})
	}
	return func(yield func(tc, string) bool) {
		for _, k := range slices.Sorted(maps.Keys(m2)) {
			for _, o := range m2[k] {
				if !yield(o, k) {
					return
				}
			}
		}
	}
}

// TODO: add cases where (linker ⋁ linked) (has ⋁ should contain) encoded parts
func TestRewrite_Rewrite_toVisitableDir(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/"}:                       "/",
		{linker: d0, linked: "/a/b"}:                    "/a/b/",
		{linker: d0, linked: "/a/b/index.tmpl"}:         "/a/b/",
		{linker: d0, linked: "/README.md"}:              "/",
		{linker: d0, linked: "a/../a/b"}:                "/a/b/",
		{linker: d0, linked: "a/../a/b/c /d"}:           "/a/b/c%20/d/",
		{linker: d0, linked: "a/../a/b/c /d/README.md"}: "/a/b/c%20/d/",
		{linker: d0, linked: "a/b"}:                     "/a/b/",
		{linker: d0, linked: "a/b/c /d"}:                "/a/b/c%20/d/",
		{linker: d0, linked: "a/b/c /d/README.md"}:      "/a/b/c%20/d/",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPageURLs(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/a/b/"}:          "/a/b/",
		{linker: d0, linked: "/"}:              "/",
		{linker: d0, linked: "a/../a/b/c /d/"}: "/a/b/c%20/d/",
		{linker: d0, linked: "a/b/c /d/"}:      "/a/b/c%20/d/",
		{linker: d2, linked: "../../"}:         "/",
		{linker: d2, linked: "./../../"}:       "/",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPage(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/page.md"}:           "/page.html",
		{linker: d0, linked: "a/../a/b/page.tmpl"}: "/a/b/page.html",
		{linker: d0, linked: "a/b/page.tmpl"}:      "/a/b/page.html",
		{linker: d2, linked: "../../page.md"}:      "/page.html",
		{linker: d2, linked: "./../../page.md"}:    "/page.html",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toVisitableDirAnchor(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/#title"}:                  "/#title",
		{linker: d0, linked: "/a/b/#title"}:              "/a/b/#title",
		{linker: d0, linked: "/a/b/index.tmpl#title"}:    "/a/b/#title",
		{linker: d0, linked: "/a/b#title"}:               "/a/b/#title",
		{linker: d0, linked: "/README.md#title"}:         "/#title",
		{linker: d0, linked: "a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "a/b/c /d/#title"}:          "/a/b/c%20/d/#title",
		{linker: d0, linked: "a/b/c /d/README.md#title"}: "/a/b/c%20/d/#title",
		{linker: d0, linked: "a/b/index.tmpl#title"}:     "/a/b/#title",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPageAnchor(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "#"}:                        "/page.html#",
		{linker: d0, linked: "#title"}:                   "/page.html#title",
		{linker: d0, linked: "a/../a/b/page.tmpl#title"}: "/a/b/page.html#title",
		{linker: d0, linked: "a/b/c /page.md#title"}:     "/a/b/c%20/page.html#title",
		{linker: d0, linked: "a/b/page.tmpl#title"}:      "/a/b/page.html#title",
		{linker: d0, linked: "page.md#title"}:            "/page.html#title",
		{linker: d2, linked: "./../../page.md#title"}:    "/page.html#title",
		{linker: d2, linked: "/page.md#title"}:           "/page.html#title",
		{linker: d2, linked: "#"}:                        "/a/b/page.html#",
		{linker: d2, linked: "#title"}:                   "/a/b/page.html#title",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_linksToUnvisitableDirs(t *testing.T) {
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := []string{
		"..",      // /a/
		"../",     // /a/
		"../b/c",  // /a/b/c%20/
		"../b/c/", // /a/b/c%20/
		"./../",   // /a/
		"./c",     // /a/b/c%20/
		"./c/",    // /a/b/c%20/
		"c",       // /a/b/c%20/
		"c/",      // /a/b/c%20/
	}
	rw := rewriter("/")
	for _, input := range tcs {
		t.Run(tescape(input), func(t *testing.T) {
			_, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}

func TestRewrite_Rewrite_linksToUnexistingNodes(t *testing.T) {
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := []string{
		"../x",                    // /x/
		"../../x.tmpl",            // /x.html
		"../../..",                // path escape
		"../../../..",             // path escape
		"a/b/c/../../../../../..", // path escape
	}
	rw := rewriter("/")
	for _, input := range tcs {
		t.Run(tescape(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success with value: %s", got)
			}
		})
	}
}

// TODO: add cases where (linker ⋁ linked) (has ⋁ should contain) encoded parts
func TestRewrite_Rewrite_toVisitableDirWithDomain(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "https://kask.ufukty.com/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/"}:                       "https://kask.ufukty.com/",
		{linker: d0, linked: "/a/b"}:                    "https://kask.ufukty.com/a/b/",
		{linker: d0, linked: "/a/b/index.tmpl"}:         "https://kask.ufukty.com/a/b/",
		{linker: d0, linked: "/README.md"}:              "https://kask.ufukty.com/",
		{linker: d0, linked: "a/../a/b"}:                "https://kask.ufukty.com/a/b/",
		{linker: d0, linked: "a/../a/b/c /d"}:           "https://kask.ufukty.com/a/b/c%20/d/",
		{linker: d0, linked: "a/../a/b/c /d/README.md"}: "https://kask.ufukty.com/a/b/c%20/d/",
		{linker: d0, linked: "a/b"}:                     "https://kask.ufukty.com/a/b/",
		{linker: d0, linked: "a/b/c /d"}:                "https://kask.ufukty.com/a/b/c%20/d/",
		{linker: d0, linked: "a/b/c /d/README.md"}:      "https://kask.ufukty.com/a/b/c%20/d/",
	}
	rw := rewriter("https://kask.ufukty.com/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPageURLsWithDomain(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "https://kask.ufukty.com/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "https://kask.ufukty.com/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/a/b/"}:          "https://kask.ufukty.com/a/b/",
		{linker: d0, linked: "/"}:              "https://kask.ufukty.com/",
		{linker: d0, linked: "a/../a/b/c /d/"}: "https://kask.ufukty.com/a/b/c%20/d/",
		{linker: d0, linked: "a/b/c /d/"}:      "https://kask.ufukty.com/a/b/c%20/d/",
		{linker: d2, linked: "../../"}:         "https://kask.ufukty.com/",
		{linker: d2, linked: "./../../"}:       "https://kask.ufukty.com/",
	}
	rw := rewriter("https://kask.ufukty.com/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPageWithDomain(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "https://kask.ufukty.com/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "https://kask.ufukty.com/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/page.md"}:           "https://kask.ufukty.com/page.html",
		{linker: d0, linked: "a/../a/b/page.tmpl"}: "https://kask.ufukty.com/a/b/page.html",
		{linker: d0, linked: "a/b/page.tmpl"}:      "https://kask.ufukty.com/a/b/page.html",
		{linker: d2, linked: "../../page.md"}:      "https://kask.ufukty.com/page.html",
		{linker: d2, linked: "./../../page.md"}:    "https://kask.ufukty.com/page.html",
	}
	rw := rewriter("https://kask.ufukty.com/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toVisitableDirAnchorWithDomain(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "https://kask.ufukty.com/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/#title"}:                  "https://kask.ufukty.com/#title",
		{linker: d0, linked: "/a/b/#title"}:              "https://kask.ufukty.com/a/b/#title",
		{linker: d0, linked: "/a/b/index.tmpl#title"}:    "https://kask.ufukty.com/a/b/#title",
		{linker: d0, linked: "/a/b#title"}:               "https://kask.ufukty.com/a/b/#title",
		{linker: d0, linked: "/README.md#title"}:         "https://kask.ufukty.com/#title",
		{linker: d0, linked: "a/b/#title"}:               "https://kask.ufukty.com/a/b/#title",
		{linker: d0, linked: "a/b/c /d/#title"}:          "https://kask.ufukty.com/a/b/c%20/d/#title",
		{linker: d0, linked: "a/b/c /d/README.md#title"}: "https://kask.ufukty.com/a/b/c%20/d/#title",
		{linker: d0, linked: "a/b/index.tmpl#title"}:     "https://kask.ufukty.com/a/b/#title",
	}
	rw := rewriter("https://kask.ufukty.com/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_toPageAnchorWithDomain(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "https://kask.ufukty.com/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "https://kask.ufukty.com/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "#"}:                        "https://kask.ufukty.com/page.html#",
		{linker: d0, linked: "#title"}:                   "https://kask.ufukty.com/page.html#title",
		{linker: d0, linked: "a/../a/b/page.tmpl#title"}: "https://kask.ufukty.com/a/b/page.html#title",
		{linker: d0, linked: "a/b/c /page.md#title"}:     "https://kask.ufukty.com/a/b/c%20/page.html#title",
		{linker: d0, linked: "a/b/page.tmpl#title"}:      "https://kask.ufukty.com/a/b/page.html#title",
		{linker: d0, linked: "page.md#title"}:            "https://kask.ufukty.com/page.html#title",
		{linker: d2, linked: "./../../page.md#title"}:    "https://kask.ufukty.com/page.html#title",
		{linker: d2, linked: "/page.md#title"}:           "https://kask.ufukty.com/page.html#title",
		{linker: d2, linked: "#"}:                        "https://kask.ufukty.com/a/b/page.html#",
		{linker: d2, linked: "#title"}:                   "https://kask.ufukty.com/a/b/page.html#title",
	}
	rw := rewriter("https://kask.ufukty.com/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("expected %q got %q", te, got)
			}
		})
	}
}

func TestRewrite_Rewrite_linksToUnvisitableDirsWithDomain(t *testing.T) {
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "https://kask.ufukty.com/a/b/page.html"}
	tcs := []string{
		"..",      // /a/
		"../",     // /a/
		"../b/c",  // /a/b/c%20/
		"../b/c/", // /a/b/c%20/
		"./../",   // /a/
		"./c",     // /a/b/c%20/
		"./c/",    // /a/b/c%20/
		"c",       // /a/b/c%20/
		"c/",      // /a/b/c%20/
	}
	rw := rewriter("https://kask.ufukty.com/")
	for _, input := range tcs {
		t.Run(tescape(input), func(t *testing.T) {
			_, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success")
			}
		})
	}
}

func TestRewrite_Rewrite_linksToUnexistingNodesWithDomain(t *testing.T) {
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "https://kask.ufukty.com/a/b/page.html"}
	tcs := []string{
		"../x",                    // /x/
		"../../x.tmpl",            // /x.html
		"../../..",                // path escape
		"../../../..",             // path escape
		"a/b/c/../../../../../..", // path escape
	}
	rw := rewriter("https://kask.ufukty.com/")
	for _, input := range tcs {
		t.Run(tescape(input), func(t *testing.T) {
			got, err := rw.Rewrite(input, linker)
			if err == nil {
				t.Errorf("act, unexpected success with value: %s", got)
			}
		})
	}
}

func TestRewrite_Rewrite_idempotency(t *testing.T) {
	d0 := paths.Paths{Src: "page.md", Dst: "page.html", Url: "/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[tc]string{
		{linker: d0, linked: "/"}:                         "/",
		{linker: d0, linked: "/"}:                         "/",
		{linker: d0, linked: "/#title"}:                   "/#title",
		{linker: d0, linked: "/a/b/"}:                     "/a/b/",
		{linker: d0, linked: "/a/b/"}:                     "/a/b/",
		{linker: d0, linked: "/a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "/a/b/"}:                     "/a/b/",
		{linker: d0, linked: "/a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "/a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "/page.html"}:                "/page.html",
		{linker: d0, linked: "/"}:                         "/",
		{linker: d0, linked: "/#title"}:                   "/#title",
		{linker: d0, linked: "/page.html#"}:               "/page.html#",
		{linker: d0, linked: "/page.html#title"}:          "/page.html#title",
		{linker: d0, linked: "/a/b/"}:                     "/a/b/",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/page.html"}:            "/a/b/page.html",
		{linker: d0, linked: "/a/b/page.html#title"}:      "/a/b/page.html#title",
		{linker: d0, linked: "/a/b/"}:                     "/a/b/",
		{linker: d0, linked: "/a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/c%20/d/#title"}:        "/a/b/c%20/d/#title",
		{linker: d0, linked: "/a/b/c%20/d/"}:              "/a/b/c%20/d/",
		{linker: d0, linked: "/a/b/c%20/d/#title"}:        "/a/b/c%20/d/#title",
		{linker: d0, linked: "/a/b/c%20/page.html#title"}: "/a/b/c%20/page.html#title",
		{linker: d0, linked: "/a/b/#title"}:               "/a/b/#title",
		{linker: d0, linked: "/a/b/page.html"}:            "/a/b/page.html",
		{linker: d0, linked: "/a/b/page.html#title"}:      "/a/b/page.html#title",
		{linker: d0, linked: "/page.html#title"}:          "/page.html#title",
		{linker: d2, linked: "/"}:                         "/",
		{linker: d2, linked: "/page.html"}:                "/page.html",
		{linker: d2, linked: "/"}:                         "/",
		{linker: d2, linked: "/page.html"}:                "/page.html",
		{linker: d2, linked: "/page.html#title"}:          "/page.html#title",
		{linker: d2, linked: "/page.html#title"}:          "/page.html#title",
		{linker: d2, linked: "/a/b/page.html#"}:           "/a/b/page.html#",
		{linker: d2, linked: "/a/b/page.html#title"}:      "/a/b/page.html#title",
	}
	rw := rewriter("/")
	for tc, te := range sorted(tcs) {
		t.Run(testname(tc.linker.Src, tc.linked), func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, tc.linker)
			if err != nil {
				t.Errorf("1st act, unexpected error: %v", err)
			} else if te != got {
				t.Errorf("1st assert, expected: %q got: %q", te, got)
			}
		})
	}
}

func TestAssetLink(t *testing.T) {
	type tc struct{ linked, expected string }
	linker := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}
	tcs := map[string]tc{
		"absolute with special character": {linked: "/.assets/x 2.png", expected: "/.assets/x%202.png"},
		"absolute":                        {linked: "/.assets/x.png", expected: "/.assets/x.png"},
		"relative with special character": {linked: ".assets/x 2.png", expected: "/a/b/.assets/x%202.png"},
		"relative":                        {linked: ".assets/x.png", expected: "/a/b/.assets/x.png"},
		"relative with parent dir and special character": {linked: "../.assets/x 2.png", expected: "/a/.assets/x%202.png"},
		"relative with parent dir":                       {linked: "../.assets/x.png", expected: "/a/.assets/x.png"},
	}
	rw := rewriter("/")
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			got, err := rw.Rewrite(tc.linked, linker)
			if err != nil {
				t.Errorf("act, unexpected error: %v", err)
			} else if tc.expected != got {
				t.Errorf("assert, expected %q got %q", tc.expected, got)
			}
		})
	}
}
