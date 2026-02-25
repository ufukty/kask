package rewriter

import (
	"testing"

	"go.ufukty.com/kask/internal/paths"
)

func TestSplit(t *testing.T) {
	// path can be:
	//   - <assets>
	//   - <page/dir>
	//   - <page/dir> <query>
	//   - <page/dir> <assets>
	//   - <page/dir> <assets> <query>
	//   - <query>
	//   - <domain> <query>
	//   - <domain> <assets>
	//   - <domain> <page/dir>
	//   - <domain> <page/dir> <query>
	//   - <domain> <page/dir> <assets>
	//   - <domain> <page/dir> <assets> <query>
	tcs := map[struct{ domain, url string }]splits{
		{"/", "#title"}:                       {"", "", "#title"},
		{"/", ".assets/img.jpg"}:              {"", ".assets/img.jpg", ""},
		{"/", ".assets/img.jpg#title"}:        {"", ".assets/img.jpg", "#title"},
		{"/", "/a/b/c"}:                       {"/", "a/b/c", ""},
		{"/", "/a/b/c#title"}:                 {"/", "a/b/c", "#title"},
		{"/", "/a/b/c/.assets/img.jpg#title"}: {"/", "a/b/c/.assets/img.jpg", "#title"},
		{"/", "a/b/c"}:                        {"", "a/b/c", ""},
		{"/", "a/b/c#title"}:                  {"", "a/b/c", "#title"},
		{"/", "a/b/c/.assets/img.jpg#title"}:  {"", "a/b/c/.assets/img.jpg", "#title"},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c"}:                       {"https://kask.ufukty.com/", "a/b/c", ""},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c#title"}:                 {"https://kask.ufukty.com/", "a/b/c", "#title"},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c/.assets/img.jpg#title"}: {"https://kask.ufukty.com/", "a/b/c/.assets/img.jpg", "#title"},
		{"https://kask.ufukty.com/", "/a/b/c"}:                                              {"/", "a/b/c", ""},
		{"https://kask.ufukty.com/", "/a/b/c#title"}:                                        {"/", "a/b/c", "#title"},
		{"https://kask.ufukty.com/", "/a/b/c/.assets/img.jpg#title"}:                        {"/", "a/b/c/.assets/img.jpg", "#title"},
		{"https://kask.ufukty.com/", "a/b/c"}:                                               {"", "a/b/c", ""},
		{"https://kask.ufukty.com/", "a/b/c#title"}:                                         {"", "a/b/c", "#title"},
		{"https://kask.ufukty.com/", "a/b/c/.assets/img.jpg#title"}:                         {"", "a/b/c/.assets/img.jpg", "#title"},
	}

	for input, expected := range tcs {
		t.Run(testname(input.domain, input.url), func(t *testing.T) {
			splits := New(paths.Paths{Url: input.domain}).split(input.url)
			if expected.base != splits.base {
				t.Errorf("assert base: expected %q, got %q", expected.base, splits.base)
			}
			if expected.ref != splits.ref {
				t.Errorf("assert ref: expected %q, got %q", expected.ref, splits.ref)
			}
			if expected.tail != splits.tail {
				t.Errorf("assert tail: expected %q, got %q", expected.tail, splits.tail)
			}
		})
	}
}
