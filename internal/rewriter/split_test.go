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
	tcs := map[struct{ domain, url string }]struct{ domain, path, assets, tail string }{
		{"/", "#title"}:                       {"", "", "", "#title"},
		{"/", ".assets/img.jpg"}:              {"", "", ".assets/img.jpg", ""},
		{"/", ".assets/img.jpg#title"}:        {"", "", ".assets/img.jpg", "#title"},
		{"/", "/a/b/c"}:                       {"/", "a/b/c", "", ""},
		{"/", "/a/b/c#title"}:                 {"/", "a/b/c", "", "#title"},
		{"/", "/a/b/c/.assets/img.jpg#title"}: {"/", "a/b/c/", ".assets/img.jpg", "#title"},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c"}:                       {"https://kask.ufukty.com/", "a/b/c", "", ""},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c#title"}:                 {"https://kask.ufukty.com/", "a/b/c", "", "#title"},
		{"https://kask.ufukty.com/", "https://kask.ufukty.com/a/b/c/.assets/img.jpg#title"}: {"https://kask.ufukty.com/", "a/b/c/", ".assets/img.jpg", "#title"},
	}

	for input, expected := range tcs {
		t.Run(testname(input.domain, input.url), func(t *testing.T) {
			d, p, a, tl := New(paths.Paths{Url: input.domain}).split(input.url)
			if expected.domain != d {
				t.Errorf("assert domain: expected %q, got %q", expected.domain, d)
			}
			if expected.path != p {
				t.Errorf("assert path: expected %q, got %q", expected.path, p)
			}
			if expected.assets != a {
				t.Errorf("assert assets: expected %q, got %q", expected.assets, a)
			}
			if expected.tail != tl {
				t.Errorf("assert query: expected %q, got %q", expected.tail, tl)
			}
		})
	}
}
