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
		"index.html":           "Index",
		"getting started.html": "Getting Started",
	}
	for input, expected := range tcs {
		t.Run(input, func(t *testing.T) {
			got := pageTitleFromFilename(input)
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
