package assert

import (
	"io/fs"
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

func EachNamedResult(t *testing.T, expected map[string]string, got []string) {
	if len(expected) != len(got) {
		t.Errorf("assert lengths: expected %d, got %d", len(expected), len(got))
	}
	for tn, expected := range expected {
		t.Run(tescape(tn), func(t *testing.T) {
			if !slices.Contains(got, expected) {
				t.Errorf("assert, expected item: %s", expected)
			}
		})
	}
}

func ResultInFile(t *testing.T, expected string, fs fs.ReadFileFS, path string) {
	c, err := fs.ReadFile(path)
	if err != nil {
		t.Fatalf("assert prep, reading file: %v", err)
	}
	s := string(c)
	if !strings.Contains(s, expected) {
		t.Errorf("assert, expected item: %s", expected)
	}
	if t.Failed() {
		t.Logf("got:\n\n%s", s)
	}
}

func EachNamedResultInFile(t *testing.T, expected map[string]string, fs fs.ReadFileFS, path string) {
	c, err := fs.ReadFile(path)
	if err != nil {
		t.Fatalf("assert prep, reading file: %v", err)
	}
	s := string(c)
	for tn, expected := range expected {
		t.Run(tescape(tn), func(t *testing.T) {
			if !strings.Contains(s, expected) {
				t.Errorf("assert, expected item: %s", expected)
			}
		})
	}
	if t.Failed() {
		t.Logf("got:\n\n%s", s)
	}
}
