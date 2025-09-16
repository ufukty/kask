package builder

import (
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
