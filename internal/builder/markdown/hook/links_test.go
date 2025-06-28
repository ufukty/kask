package hook

import (
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown/ast"
)

func TestVisitor_links(t *testing.T) {
	tcs := map[string]struct {
		page, link, expected string
	}{
		"subdir":           {"README.md", "a", "a"},
		"subdir's readme'": {"README.md", "a/README.md", "a"},
		"subdir's page":    {"README.md", "a/b.md", "a/b.html"},

		"parent dir":                   {"a/b/c.md", "..", "a"},
		"parent dir's readme'":         {"a/b/c.md", "../README.md", "a"},
		"parent dir is root":           {"a/b.md", "..", "/"},
		"parent dir is root's readme'": {"a/b.md", "../README.md", "/"},
		"parent dir's page in root":    {"a/b.md", "../d.md", "d.html"},

		"double parent dir":                   {"a/b/c/d.md", "../..", "a"},
		"double parent dir's readme'":         {"a/b/c/d.md", "../../README.md", "a"},
		"double parent dir is root":           {"a/b/c.md", "../..", "/"},
		"double parent dir is root's readme'": {"a/b/c.md", "../../README.md", "/"},
		"double parent dir's page in root":    {"a/b/c.md", "../../d.md", "d.html"},

		"prefixed subdir":           {"README.md", "./a", "a"},
		"prefixed subdir's readme'": {"README.md", "./a/README.md", "a"},
		"prefixed subdir's page":    {"README.md", "./a/b.md", "a/b.html"},

		"prefixed parent dir":                   {"a/b/c.md", "./..", "a"},
		"prefixed parent dir's readme'":         {"a/b/c.md", "./../README.md", "a"},
		"prefixed parent dir is root":           {"a/b.md", "./..", "/"},
		"prefixed parent dir is root's readme'": {"a/b.md", "./../README.md", "/"},
		"prefixed parent dir's page in root":    {"a/b.md", "./../d.md", "d.html"},

		"prefixed double parent dir":                   {"a/b/c/d.md", "./../..", "a"},
		"prefixed double parent dir's readme'":         {"a/b/c/d.md", "./../../README.md", "a"},
		"prefixed double parent dir is root":           {"a/b/c.md", "./../..", "/"},
		"prefixed double parent dir is root's readme'": {"a/b/c.md", "./../../README.md", "/"},
		"prefixed double parent dir's page in root":    {"a/b/c.md", "./../../d.md", "d.html"},
	}

	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			n := &ast.Link{Destination: []byte(tc.link)}
			NewVisitor(filepath.Dir(tc.page)).links(n)
			if string(n.Destination) != tc.expected {
				t.Errorf("expected %q got %q", tc.expected, string(n.Destination))
			}
		})
	}
}
