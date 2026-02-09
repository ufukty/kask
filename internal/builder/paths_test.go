package builder

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
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
			got := dirUri(i.parent, i.child, true)
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
			got := fileUri(i.parent, i.child, false)
			if got != o {
				t.Errorf("expected %q got %q", o, got)
			}
		})
	}
}

func TestPaths_Subdir(t *testing.T) {
	parent := paths{
		src: "",
		dst: "",
		url: "/",
	}
	got := parent.subdir("c", true)
	expected := paths{
		src: "c",
		dst: "c",
		url: "/c/",
	}
	if got != expected {
		t.Errorf("assert, expected %#v, got %#v", expected, got)
	}
}

func TestPaths_File(t *testing.T) {
	parent := paths{
		src: "/a",
		dst: "/a",
		url: "/a/",
	}
	type tc struct {
		inputBasename string
		inputStripped bool
		expected      paths
	}
	tcs := map[string]tc{
		"index file md with stripped ordering":                 {inputBasename: "README.md", inputStripped: true, expected: paths{src: "/a/README.md", dst: "/a/index.html", url: "/a/"}},
		"index file md":                                        {inputBasename: "README.md", inputStripped: false, expected: paths{src: "/a/README.md", dst: "/a/index.html", url: "/a/"}},
		"index file tmpl with stripped ordering":               {inputBasename: "index.tmpl", inputStripped: true, expected: paths{src: "/a/index.tmpl", dst: "/a/index.html", url: "/a/"}},
		"index file tmpl":                                      {inputBasename: "index.tmpl", inputStripped: false, expected: paths{src: "/a/index.tmpl", dst: "/a/index.html", url: "/a/"}},
		"non-index file with stripped ordering":                {inputBasename: "3.page.tmpl", inputStripped: true, expected: paths{src: "/a/3.page.tmpl", dst: "/a/page.html", url: "/a/page.html"}},
		"non-index file with whitespace and stripped ordering": {inputBasename: "3.pge .tmpl", inputStripped: true, expected: paths{src: "/a/3.pge .tmpl", dst: "/a/pge .html", url: "/a/pge%20.html"}},
		"non-index file with whitespace":                       {inputBasename: "3.pge .tmpl", inputStripped: false, expected: paths{src: "/a/3.pge .tmpl", dst: "/a/3.pge .html", url: "/a/3.pge%20.html"}},
		"non-index file":                                       {inputBasename: "3.page.tmpl", inputStripped: false, expected: paths{src: "/a/3.page.tmpl", dst: "/a/3.page.html", url: "/a/3.page.html"}},
	}
	for _, tn := range slices.Sorted(maps.Keys(tcs)) {
		tc := tcs[tn]
		got, err := parent.file(tc.inputBasename, tc.inputStripped)
		if err != nil {
			t.Run(filepath.Join(tn, "act"), func(t *testing.T) { t.Fatalf("unexpected error: %v", err) })
		}
		if got.src != tc.expected.src {
			t.Run(filepath.Join(tn, "assert src"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.src, got.src) })
		}
		if got.dst != tc.expected.dst {
			t.Run(filepath.Join(tn, "assert dst"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.dst, got.dst) })
		}
		if got.url != tc.expected.url {
			t.Run(filepath.Join(tn, "assert url"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.url, got.url) })
		}
	}
}

func TestPaths_File_preserveEncodedParent(t *testing.T) {
	parent := paths{
		src: "/a ",
		dst: "/a ",
		url: "/a%20/",
	}
	got, err := parent.file("b.tmpl", false)
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	if filepath.Dir(got.url) != "/a%20" {
		t.Errorf("assert, expected the parent path to stay encoded: %q, got. %q", "/a%20", got.url)
	}
}

func TestPaths_Dir_preserveEncodedParent(t *testing.T) {
	parent := paths{
		src: "/a ",
		dst: "/a ",
		url: "/a%20/",
	}
	got := parent.subdir("b", false)
	if filepath.Dir(filepath.Dir(got.url)) != "/a%20" {
		t.Errorf("assert, expected the parent path to stay encoded: %q, got. %q", "/a%20", got.url)
	}
}
