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

func formatstr(expected, got string) string {
	if len(expected) > 10 || len(got) > 10 {
		return "assert, values:\nexp: %s\ngot: %s"
	}
	return "assert, values: expected: %s got: %s"
}

func Results(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf(formatstr(expected, got), expected, got)
	}
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
	if t.Failed() {
		t.Logf("got:\n\n%s", strings.Join(got, "\n"))
	}
}

func Order(t *testing.T, expected, got []string) {
	if len(expected) != len(got) {
		t.Errorf("assert lengths: expected %d, got %d", len(expected), len(got))
	}
	for i := range len(expected) {
		if expected[i] != got[i] {
			t.Fatalf("order is broken at item %d:\n\texp: %s\n\tgot: %s", i, strings.Join(expected, ", "), strings.Join(got, ", "))
		}
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
	if t.Failed() {
		t.Logf("got:\n\n%s", strings.Join(got, "\n"))
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

// expected is files to contents
func ResultsInFiles(t *testing.T, fs fs.ReadFileFS, expected map[string]string) {
	for file, content := range expected {
		c, err := fs.ReadFile(file)
		if err != nil {
			t.Fatalf("assert, prep: reading file: %v", err)
		}
		s := string(c)
		if !strings.Contains(s, content) {
			t.Errorf("assert, expected item: %s", content)
		}
		if t.Failed() {
			t.Logf("file contents for %s:\n\n%s", file, s)
		}
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

func NamedResultsInFile(t *testing.T, expected, unexpected map[string]string, fs fs.ReadFileFS, path string) {
	c, err := fs.ReadFile(path)
	if err != nil {
		t.Fatalf("assert prep, reading file: %v", err)
	}
	s := string(c)
	t.Run("positive", func(t *testing.T) {
		for tn, expected := range expected {
			t.Run(tescape(tn), func(t *testing.T) {
				if !strings.Contains(s, expected) {
					t.Errorf("assert, expected item: %s", expected)
				}
			})
		}
	})
	t.Run("negative", func(t *testing.T) {
		for tn, unexpected := range unexpected {
			t.Run(tescape(tn), func(t *testing.T) {
				if strings.Contains(s, unexpected) {
					t.Errorf("assert, unexpected item: %s", unexpected)
				}
			})
		}
	})
	if t.Failed() {
		t.Logf("got:\n\n%s", s)
	}
}
