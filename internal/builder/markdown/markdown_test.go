package markdown

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_FromFile(t *testing.T) {
	p, err := ToHtml("testdata", "input.md")
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	fmt.Println(p.Content)
	fmt.Println(p.Toc)
}

func unmarshal(content string) map[string]string {
	matches := regexp.MustCompile(`(?m)^.*<a href="([^"]*)"[^>]*>([^<]*)</a>.*$`).FindAllStringSubmatch(content, -1)
	tcs := map[string]string{} // input => got
	for _, match := range matches {
		tcs[match[2]] = match[1]
	}
	return tcs
}

func TestToHtml_links(t *testing.T) {
	p, err := ToHtml("testdata", "subdir/subsubdir/README.md")
	if err != nil {
		panic(fmt.Errorf("act, ToHtml: %w", err))
	}
	expected := map[string]string{
		// links to parents
		"..":                "/subdir",
		"../..":             "/",
		"../../":            "/",
		"../../a.md":        "/a.html",
		"../../README.md":   "/",
		"../":               "/subdir",
		"../a.md":           "/subdir/a.html",
		"../README.md":      "/subdir",
		".":                 "/subdir/subsubdir",
		"./..":              "/subdir",
		"./../..":           "/",
		"./../../":          "/",
		"./../":             "/subdir",
		"./../a.md":         "/subdir/a.html",
		"./../README.md":    "/subdir",
		"./../../a.md":      "/a.html",
		"./../../README.md": "/",

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
	}
	got := unmarshal(p.Content)
	if len(expected) != len(got) {
		t.Errorf("expected len(expected) = len(got) got %d != %d", len(expected), len(got))
	}
	for link, expected := range expected {
		if _, ok := got[link]; !ok {
			t.Errorf("for %q expected %q got nothing", link, expected)
		} else if expected != got[link] {
			t.Errorf("for %q expected %q got %q", link, expected, got[link])
		}
	}
}
