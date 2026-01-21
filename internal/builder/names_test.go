package builder

import (
	"path/filepath"
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

func Test_extractor_FromFile(t *testing.T) {
	tcs := map[string]string{
		"testdata/extractors/page.tmpl": "Page title",
		"testdata/extractors/page.md":   "Page title",
	}
	for path, expected := range tcs {
		t.Run(filepath.Ext(path), func(t *testing.T) {
			got, err := theExtractor.FromFile(path)
			if err != nil {
				t.Fatalf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Fatalf("assert, expected: %q, got: %q", expected, got)
			}
		})
	}
}
