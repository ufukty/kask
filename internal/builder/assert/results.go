package assert

import (
	"slices"
	"strings"
	"testing"
)

func tescape(s string) string {
	return strings.ReplaceAll(s, "/", "\\")
}

func EachResult(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Errorf("assert lengths: expected %d, got %d", len(expected), len(got))
	}
	for _, expected := range expected {
		t.Run(tescape(expected), func(t *testing.T) {
			if !slices.Contains(got, expected) {
				t.Errorf("assert, expected item: %s", expected)
			}
		})
	}
}
