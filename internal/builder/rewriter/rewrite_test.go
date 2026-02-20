package rewriter

import (
	"fmt"
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

func TestRewrite_Rewrite(t *testing.T) {
	type tc struct {
		linker paths.Paths
		linked string
	}
	d0 := paths.Paths{Src: "page.tmpl", Dst: "page.html", Url: "/page.html"}
	d2 := paths.Paths{Src: "a/b/page.tmpl", Dst: "a/b/page.html", Url: "/a/b/page.html"}

	tcs := map[tc]string{
		{linker: d0, linked: "/"}:                     "/",
		{linker: d0, linked: "/#title"}:               "/#title",
		{linker: d0, linked: "/a/b"}:                  "/a/b/",
		{linker: d0, linked: "/a/b/index.tmpl"}:       "/a/b/",
		{linker: d0, linked: "/a/b/index.tmpl#title"}: "/a/b/#title",
		{linker: d0, linked: "/page.md"}:              "/page.html",
		{linker: d0, linked: "/README.md"}:            "/",
		{linker: d0, linked: "/README.md#title"}:      "/#title",

		{linker: d0, linked: "#"}:      "/#",
		{linker: d0, linked: "#title"}: "/#title",
		{linker: d2, linked: "#"}:      "/a/b/page.html/#",
		{linker: d2, linked: "#title"}: "/a/b/page.html/#title",

		{linker: d0, linked: "a/../a/b"}:                 "/a/b/",
		{linker: d0, linked: "a/../a/b/c/d"}:             "/a/b/c/d/",
		{linker: d0, linked: "a/../a/b/c/d/README.md"}:   "/a/b/c/d/",
		{linker: d0, linked: "a/../a/b/page.tmpl"}:       "/a/b/page.html",
		{linker: d0, linked: "a/../a/b/page.tmpl#title"}: "/a/b/page.html#title",

		{linker: d0, linked: "a/b"}:                     "/a/b/",
		{linker: d0, linked: "a/b/#title"}:              "/a/b/#title",
		{linker: d0, linked: "a/b/c/d"}:                 "/a/b/c/d/",
		{linker: d0, linked: "a/b/c/d/#title"}:          "/a/b/c/d/#title",
		{linker: d0, linked: "a/b/c/d/README.md"}:       "/a/b/c/d/",
		{linker: d0, linked: "a/b/c/d/README.md#title"}: "/a/b/c/d/#title",
		{linker: d0, linked: "a/b/c/page.md#title"}:     "/a/b/c/page.html#title",
		{linker: d0, linked: "a/b/index.tmpl#title"}:    "/a/b/#title",
		{linker: d0, linked: "a/b/page.tmpl"}:           "/a/b/page.html",
		{linker: d0, linked: "a/b/page.tmpl#title"}:     "/a/b/page.html#title",
		{linker: d0, linked: "README.md#title"}:         "/#title",

		{linker: d2, linked: "../../"}:                "/",
		{linker: d2, linked: "../../page.md"}:         "/page.html",
		{linker: d2, linked: "./../../"}:              "/",
		{linker: d2, linked: "./../../page.md"}:       "/page.html",
		{linker: d2, linked: "./../../page.md#title"}: "/page.html#title",
		{linker: d2, linked: "/"}:                     "/",
	}
	for tc, te := range tcs {
		t.Run(fmt.Sprintf("%s=>%s", testname(tc.linker.Dst), testname(tc.linked)), func(t *testing.T) {
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

func TestRewrite_Rewrite_linksToUnexistingNodes(t *testing.T) {
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
