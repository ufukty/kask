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
			got := titleFromFilename(input, ".tmpl")
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
			got := hrefFromFilename(input.dstPathEncoded, input.filename)
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
			got := targetFromFilename(input.dst, input.dstPath, input.filename)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestExtractors(t *testing.T) {

	tcs := map[string]string{ // paths to expected titles
		"career/index.tmpl":                 "Careers",
		"docs/birdseed.md":                  "ACME Bird Seed",
		"docs/download.md":                  "Download",
		"docs/magnet.md":                    "ACME Magnet",
		"docs/README.md":                    "Docs",
		"docs/tutorials/getting-started.md": "Getting Started",
		"index.tmpl":                        "Acme",
		"products/index.tmpl":               "Products",
	}

	for path, expected := range tcs {
		t.Run(strings.ReplaceAll(path, "/", "\\"), func(t *testing.T) {
			c, err := os.ReadFile(filepath.Join("testdata/acme", path))
			if err != nil {
				t.Fatalf("prep, read file: %v", err)
			}
			got := titleFromContent(string(c), filepath.Ext(path))
			if got != expected {
				t.Errorf("assert, expected %q got %q", expected, got)
			}
		})
	}
}
