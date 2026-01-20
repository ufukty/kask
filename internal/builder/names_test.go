package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripOrdering(t *testing.T) {
	tcs := map[string]string{
		"1.contacts":     "contacts",
		"10.contacts":    "contacts",
		"10. contacts":   "contacts",
		"001.contacts":   "contacts",
		"001 - contacts": "contacts",
		"001 contacts":   "contacts",
		"001.  contacts": "contacts",
		"001.. contacts": "contacts",
	}

	for input, expected := range tcs {
		t.Run(input, func(t *testing.T) {
			got := stripOrdering(input)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestTitleFromFilename(t *testing.T) {
	tcs := map[string]string{
		"index.tmpl":             "Index",
		"3 getting started.tmpl": "Getting Started",
	}
	for input, expected := range tcs {
		t.Run(input, func(t *testing.T) {
			got := titleFromFilename(input, ".tmpl", true)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestHrefFromFilename(t *testing.T) {
	type tc struct {
		testname       string
		dstPathEncoded string
		filename       string
	}
	tcs := map[tc]string{
		{"extension repl", "a/b/c", "d.md"}:   "/a/b/c/d.html",
		{"encoded path", "a/b%20/c", "d.md"}:  "/a/b%20/c/d.html",
		{"filename enc", "a/b/c", "d .md"}:    "/a/b/c/d%20.html",
		{"strip ordering", "a/b/c", "3.d.md"}: "/a/b/c/d.html",
	}
	for input, expected := range tcs {
		t.Run(input.testname, func(t *testing.T) {
			got := hrefFromFilename(input.dstPathEncoded, input.filename, true)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestTargetFromFilename(t *testing.T) {
	type tc struct {
		testname               string
		dst, dstPath, filename string
	}
	tcs := map[tc]string{
		{"repl extension", "/a", "b/c", "d.md"}:   "/a/b/c/d.html",
		{"strip ordering", "/a", "b/c", "3 d.md"}: "/a/b/c/d.html",
	}
	for input, expected := range tcs {
		t.Run(input.testname, func(t *testing.T) {
			got := targetFromFilename(input.dst, input.dstPath, input.filename, true)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestTitleFromContent(t *testing.T) {
	type tc struct {
		path, ext string
	}
	tcs := map[tc]string{ // paths to expected titles
		{"1.career/index.tmpl", ".tmpl"}:                     "Careers at ACME",
		{"2.products/index.tmpl", ".tmpl"}:                   "ACME Products",
		{"3.docs/101 tutorials/1.getting-started.md", ".md"}: "Getting Started",
		{"3.docs/birdseed.md", ".md"}:                        "ACME Bird Seed",
		{"3.docs/download.md", ".md"}:                        "Download",
		{"3.docs/magnet.md", ".md"}:                          "ACME Magnet",
		{"3.docs/README.md", ".md"}:                          "Docs",
		{"index.tmpl", ".tmpl"}:                              "Acme",
	}

	for tc, expected := range tcs {
		t.Run(strings.ReplaceAll(tc.path, "/", "\\"), func(t *testing.T) {
			c, err := os.ReadFile(filepath.Join("testdata/acme", tc.path))
			if err != nil {
				t.Fatalf("prep, read file: %v", err)
			}
			got := titleFromContent(string(c), tc.ext)
			if got != expected {
				t.Errorf("assert, expected %q got %q", expected, got)
			}
		})
	}
}
