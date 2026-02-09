package builder

import (
	"fmt"
	"strings"
	"testing"
)

func TestUri_dir(t *testing.T) {
	type input struct{ parent, child string }
	type output = string
	tcs := map[input]output{
		{"", "a"}:      "/a/",
		{"/a", "b"}:    "/a/b/",
		{"/a/", "b"}:   "/a/b/",
		{"/a/b/", "c"}: "/a/b/c/",
	}

	for i, o := range tcs {
		tn := fmt.Sprintf("parent=%q dir=%q",
			strings.ReplaceAll(i.parent, "/", "\\"),
			strings.ReplaceAll(i.child, "/", "\\"),
		)
		t.Run(tn, func(t *testing.T) {
			got := uri(i.parent, i.child, true)
			if got != o {
				t.Errorf("expected %q got %q", o, got)
			}
		})
	}
}

func TestUri_file(t *testing.T) {
	type input struct{ parent, child string }
	type output = string
	tcs := map[input]output{
		{"", "a.md"}:      "/a.html",
		{"/a", "b.md"}:    "/a/b.html",
		{"/a/", "b.md"}:   "/a/b.html",
		{"/a/b/", "c.md"}: "/a/b/c.html",
	}

	for i, o := range tcs {
		tn := fmt.Sprintf("parent=%q dir=%q",
			strings.ReplaceAll(i.parent, "/", "\\"),
			strings.ReplaceAll(i.child, "/", "\\"),
		)
		t.Run(tn, func(t *testing.T) {
			got := uri(i.parent, i.child, false)
			if got != o {
				t.Errorf("expected %q got %q", o, got)
			}
		})
	}
}

func TestPaths_Subdir(t *testing.T) {
	parent := paths{
		src: "src",
		dst: "dst",
		url: "/",
	}
	got := parent.sub("c", dir, true)
	expected := paths{
		src: "src/c",
		dst: "dst/c",
		url: "/c/",
	}
	if got != expected {
		t.Errorf("assert, expected %#v, got %#v", expected, got)
	}
}
