package paths

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"strings"
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

func TestPaths_File(t *testing.T) {
	parent := Paths{
		Src: "/a",
		Dst: "/a",
		Url: "/a/",
	}
	type tc struct {
		inputBasename string
		inputStripped bool
		expected      Paths
	}
	tcs := map[string]tc{
		"index file md with stripped ordering":                 {inputBasename: "README.md", inputStripped: true, expected: Paths{Src: "/a/README.md", Dst: "/a/index.html", Url: "/a/"}},
		"index file md":                                        {inputBasename: "README.md", inputStripped: false, expected: Paths{Src: "/a/README.md", Dst: "/a/index.html", Url: "/a/"}},
		"index file tmpl with stripped ordering":               {inputBasename: "index.tmpl", inputStripped: true, expected: Paths{Src: "/a/index.tmpl", Dst: "/a/index.html", Url: "/a/"}},
		"index file tmpl":                                      {inputBasename: "index.tmpl", inputStripped: false, expected: Paths{Src: "/a/index.tmpl", Dst: "/a/index.html", Url: "/a/"}},
		"non-index file with stripped ordering":                {inputBasename: "3.page.tmpl", inputStripped: true, expected: Paths{Src: "/a/3.page.tmpl", Dst: "/a/page.html", Url: "/a/page.html"}},
		"non-index file with whitespace and stripped ordering": {inputBasename: "3.pge .tmpl", inputStripped: true, expected: Paths{Src: "/a/3.pge .tmpl", Dst: "/a/pge .html", Url: "/a/pge%20.html"}},
		"non-index file with whitespace":                       {inputBasename: "3.pge .tmpl", inputStripped: false, expected: Paths{Src: "/a/3.pge .tmpl", Dst: "/a/3.pge .html", Url: "/a/3.pge%20.html"}},
		"non-index file":                                       {inputBasename: "3.page.tmpl", inputStripped: false, expected: Paths{Src: "/a/3.page.tmpl", Dst: "/a/3.page.html", Url: "/a/3.page.html"}},
	}
	for _, tn := range slices.Sorted(maps.Keys(tcs)) {
		tc := tcs[tn]
		got := parent.File(tc.inputBasename, tc.inputStripped)
		if got.Src != tc.expected.Src {
			t.Run(filepath.Join(tn, "assert src"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Src, got.Src) })
		}
		if got.Dst != tc.expected.Dst {
			t.Run(filepath.Join(tn, "assert dst"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Dst, got.Dst) })
		}
		if got.Url != tc.expected.Url {
			t.Run(filepath.Join(tn, "assert url"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Url, got.Url) })
		}
	}
}

func TestPaths_Dir(t *testing.T) {
	parent := Paths{
		Src: "/a",
		Dst: "/a",
		Url: "/a/",
	}
	type tc struct {
		inputBasename string
		inputStripped bool
		expected      Paths
	}
	tcs := map[string]tc{
		"subdir":                                      {inputBasename: "1.b", inputStripped: false, expected: Paths{Src: "/a/1.b", Dst: "/a/1.b", Url: "/a/1.b/"}},
		"subdir with special char":                    {inputBasename: "1.b ", inputStripped: false, expected: Paths{Src: "/a/1.b ", Dst: "/a/1.b ", Url: "/a/1.b%20/"}},
		"subdir with stripped ordering":               {inputBasename: "1.b", inputStripped: true, expected: Paths{Src: "/a/1.b", Dst: "/a/b", Url: "/a/b/"}},
		"subdir with special char and strip ordering": {inputBasename: "1.b ", inputStripped: true, expected: Paths{Src: "/a/1.b ", Dst: "/a/b ", Url: "/a/b%20/"}},
	}
	for _, tn := range slices.Sorted(maps.Keys(tcs)) {
		tc := tcs[tn]
		got := parent.Subdir(tc.inputBasename, tc.inputStripped)
		if got.Src != tc.expected.Src {
			t.Run(filepath.Join(tn, "src"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Src, got.Src) })
		}
		if got.Dst != tc.expected.Dst {
			t.Run(filepath.Join(tn, "dst"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Dst, got.Dst) })
		}
		if got.Url != tc.expected.Url {
			t.Run(filepath.Join(tn, "url"), func(t *testing.T) { t.Errorf("expected %q got %q", tc.expected.Url, got.Url) })
		}
	}
}

func TestPaths_File_preserveEncodedParent(t *testing.T) {
	parent := Paths{
		Src: "/a ",
		Dst: "/a ",
		Url: "/a%20/",
	}
	got := parent.File("b.tmpl", false)
	if filepath.Dir(got.Url) != "/a%20" {
		t.Errorf("assert, expected the parent path to stay encoded: %q, got. %q", "/a%20", got.Url)
	}
}

func TestPaths_Dir_preserveEncodedParent(t *testing.T) {
	parent := Paths{
		Src: "/a ",
		Dst: "/a ",
		Url: "/a%20/",
	}
	got := parent.Subdir("b", false)
	if filepath.Dir(filepath.Dir(got.Url)) != "/a%20" {
		t.Errorf("assert, expected the parent path to stay encoded: %q, got. %q", "/a%20", got.Url)
	}
}
