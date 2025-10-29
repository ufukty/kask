package hook

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestRewrite(t *testing.T) {
	r := map[string]string{
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

	tcs := map[string]string{
		// links to parents
		"..":                "/subdir",
		"../..":             "/",
		"../../":            "/",
		"../../a.md":        "/a.html",
		"./../../a.md":      "/",
		"../../README.md":   "/subdir",
		"./../../README.md": "/subdir/a.html",
		"../":               "/subdir",
		"../a.md":           "/subdir/subsubdir",
		"../README.md":      "/subdir",
		".":                 "/",
		"./..":              "/",
		"./../..":           "/subdir",
		"./../../":          "/subdir/a.html",
		"./../":             "/subdir",
		"./../a.md":         "/a.html",
		"./../README.md":    "/",

		// links to subdirs
		"a":             "/subdir/subsubdir/a",
		"./a":           "/subdir/subsubdir/a",
		"./a.md":        "/subdir/subsubdir/a.html",
		"./a/b.md":      "/subdir/subsubdir/a/b.html",
		"./a/README.md": "/subdir/subsubdir/a",
		"a.md":          "/subdir/subsubdir/a.html",
		"a/b.md":        "/subdir/subsubdir/a/b.html",
		"a/README.md":   "/subdir/subsubdir/a",

		// links with redundancy
		"./../subsubdir/a.md":        "/subdir/subsubdir/a.html",
		"./../subsubdir/a/b.md":      "/subdir/subsubdir/a/b.html",
		"./../subsubdir/a/README.md": "/subdir/subsubdir/a",
		"../subsubdir/a.md":          "/subdir/subsubdir/a.html",
		"../subsubdir/a/b.md":        "/subdir/subsubdir/a/b.html",
		"../subsubdir/a/README.md":   "/subdir/subsubdir/a",

		// links to paths with stripped ordering
		"1.%20lorem/":                  "/subdir/subsubdir/lorem/",
		"1.%20lorem/1.%20ipsum.md":     "/subdir/subsubdir/lorem/ipsum.html",
		"3.%20sit":                     "/subdir/subsubdir/sit/",
		"3.%20sit/2.%20consectetur.md": "/subdir/subsubdir/sit/consectetur.html",
	}

	for _, link := range slices.Sorted(maps.Keys(tcs)) {
		expected := tcs[link]
		testname := strings.ReplaceAll(link, "/", "\\")
		testname = strings.ReplaceAll(testname, "%20", " ")
		t.Run(testname, func(t *testing.T) {
			if got := rewrite(link, "/", r); expected != got {
				t.Errorf("for %q expected %q got %q", link, expected, got)
			}
		})
	}
}
