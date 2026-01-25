package builder

import (
	"path/filepath"
	"testing"

	"github.com/ufukty/kask/internal/builder/directory"
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
			got := pageTitleFromFilename(input, ".tmpl", true)
			if got != expected {
				t.Errorf("expected %q got %q", expected, got)
			}
		})
	}
}

func TestHrefFromFilename(t *testing.T) {
	type tc struct {
		testname         string
		dirUrl           string
		filename         string
		preserveOrdering bool
	}
	tcs := map[tc]string{
		{"encoded path with ordering", "a/b%20/c", "d.md", true}:  "/a/b%20/c/d.html",
		{"encoded path", "a/b%20/c", "d.md", false}:               "/a/b%20/c/d.html",
		{"extension repl with ordering", "a/b/c", "d.md", true}:   "/a/b/c/d.html",
		{"extension repl", "a/b/c", "d.md", false}:                "/a/b/c/d.html",
		{"filename enc with ordering", "a/b/c", "d .md", true}:    "/a/b/c/d%20.html",
		{"filename enc", "a/b/c", "d .md", false}:                 "/a/b/c/d%20.html",
		{"strip ordering with ordering", "a/b/c", "3.d.md", true}: "/a/b/c/3.d.html",
		{"strip ordering", "a/b/c", "3.d.md", false}:              "/a/b/c/d.html",
	}
	for input, expected := range tcs {
		t.Run(input.testname, func(t *testing.T) {
			d := &dir2{
				meta:  &directory.Meta{PreserveOrdering: input.preserveOrdering},
				paths: paths{url: input.dirUrl},
			}
			got := pageLinkFromFilename(d, input.filename)
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
			got := pageDestFromFilename(input.dst, input.dstPath, input.filename, true)
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
