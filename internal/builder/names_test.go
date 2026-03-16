package builder

import (
	"path/filepath"
	"testing"

	"go.ufukty.com/kask/internal/disk"
)

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
		"page.tmpl": "Page title",
		"page.md":   "Page title",
	}
	for path, expected := range tcs {
		t.Run(filepath.Ext(path), func(t *testing.T) {
			got, err := theExtractor.FromFile(disk.NewReal("testdata/extractors"), path)
			if err != nil {
				t.Fatalf("act, unexpected error: %v", err)
			}
			if expected != got {
				t.Fatalf("assert, expected: %q, got: %q", expected, got)
			}
		})
	}
}
