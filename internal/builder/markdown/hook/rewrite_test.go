package hook

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

var rewrites = map[string]string{
	"/a.md":                                      "/a.html",
	"/README.md":                                 "/",
	"/subdir/a.md":                               "/subdir/a.html",
	"/subdir/README.md":                          "/subdir/",
	"/subdir/subsubdir/a.md":                     "/subdir/subsubdir/a.html",
	"/subdir/subsubdir/README.md":                "/subdir/subsubdir/",
	"/subdir/subsubdir/a/b.md":                   "/subdir/subsubdir/a/b.html",
	"/subdir/subsubdir/a/README.md":              "/subdir/subsubdir/a/",
	"/subdir/subsubdir/3. sit":                   "/subdir/subsubdir/sit/",
	"/subdir/subsubdir/3. sit/2. consectetur.md": "/subdir/subsubdir/sit/consectetur.html",
	"/subdir/subsubdir/1. lorem":                 "/subdir/subsubdir/lorem/",
	"/subdir/subsubdir/1. lorem/1. ipsum.md":     "/subdir/subsubdir/lorem/ipsum.html",
}

func testname(tn string) string {
	tn = strings.ReplaceAll(tn, "/", "\\")
	tn = strings.ReplaceAll(tn, "%20", " ")
	return tn
}

func TestRewrite_linksToParents(t *testing.T) {
	tcs := map[string]string{
		"..":                "/subdir/",
		"../..":             "/",
		"../../":            "/",
		"../../a.md":        "/a.html",
		"../../README.md":   "/",
		"../":               "/subdir/",
		"../a.md":           "/subdir/a.html",
		"../README.md":      "/subdir/",
		".":                 "/subdir/subsubdir/",
		"./..":              "/subdir/",
		"./../..":           "/",
		"./../../":          "/",
		"./../../a.md":      "/a.html",
		"./../../README.md": "/",
		"./../":             "/subdir/",
		"./../a.md":         "/subdir/a.html",
		"./../README.md":    "/subdir/",
	}

	for _, link := range slices.Sorted(maps.Keys(tcs)) {
		t.Run(testname(link), func(t *testing.T) {
			got := rewrite(link, "/subdir/subsubdir", rewrites)
			expected := tcs[link]
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksToSubdirs(t *testing.T) {
	tcs := map[string]string{
		"./a.md":        "/subdir/subsubdir/a.html",
		"./a":           "/subdir/subsubdir/a/",
		"./a/b.md":      "/subdir/subsubdir/a/b.html",
		"./a/README.md": "/subdir/subsubdir/a/",
		"a.md":          "/subdir/subsubdir/a.html",
		"a":             "/subdir/subsubdir/a/",
		"a/b.md":        "/subdir/subsubdir/a/b.html",
		"a/README.md":   "/subdir/subsubdir/a/",
	}

	for _, link := range slices.Sorted(maps.Keys(tcs)) {
		t.Run(testname(link), func(t *testing.T) {
			got := rewrite(link, "/subdir/subsubdir", rewrites)
			expected := tcs[link]
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksWithReduntantSegments(t *testing.T) {
	tcs := map[string]string{
		"../subsubdir/a.md":          "/subdir/subsubdir/a.html",
		"../subsubdir/a/b.md":        "/subdir/subsubdir/a/b.html",
		"../subsubdir/a/README.md":   "/subdir/subsubdir/a/",
		"./../subsubdir/a.md":        "/subdir/subsubdir/a.html",
		"./../subsubdir/a/b.md":      "/subdir/subsubdir/a/b.html",
		"./../subsubdir/a/README.md": "/subdir/subsubdir/a/",
	}

	for _, link := range slices.Sorted(maps.Keys(tcs)) {
		t.Run(testname(link), func(t *testing.T) {
			got := rewrite(link, "/subdir/subsubdir", rewrites)
			expected := tcs[link]
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestRewrite_linksWithPathsWithStrippedOrdering(t *testing.T) {
	tcs := map[string]string{
		"../subsubdir/1.%20lorem/":                  "/subdir/subsubdir/lorem/",
		"../subsubdir/1.%20lorem/1.%20ipsum.md":     "/subdir/subsubdir/lorem/ipsum.html",
		"../subsubdir/3.%20sit":                     "/subdir/subsubdir/sit/",
		"../subsubdir/3.%20sit/2.%20consectetur.md": "/subdir/subsubdir/sit/consectetur.html",

		"./1.%20lorem/":                  "/subdir/subsubdir/lorem/",
		"./1.%20lorem/1.%20ipsum.md":     "/subdir/subsubdir/lorem/ipsum.html",
		"./3.%20sit":                     "/subdir/subsubdir/sit/",
		"./3.%20sit/2.%20consectetur.md": "/subdir/subsubdir/sit/consectetur.html",

		"1.%20lorem/":                  "/subdir/subsubdir/lorem/",
		"1.%20lorem/1.%20ipsum.md":     "/subdir/subsubdir/lorem/ipsum.html",
		"3.%20sit":                     "/subdir/subsubdir/sit/",
		"3.%20sit/2.%20consectetur.md": "/subdir/subsubdir/sit/consectetur.html",
	}

	for _, link := range slices.Sorted(maps.Keys(tcs)) {
		t.Run(testname(link), func(t *testing.T) {
			got := rewrite(link, "/subdir/subsubdir", rewrites)
			expected := tcs[link]
			if expected != got {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}
