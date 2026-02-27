package builder

import (
	"fmt"
	"testing"

	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
)

func Example_patterns_validateLinkMatchers() {
	example := []byte(`<a href="anchor-target">Lorem ipsum dolor sit amet.</a><img src="img-source" srcset="img-source-set-2x 2x, img-source-set-3x 3x, img-source-set-wide 1000w">`)
	for _, lm := range linkMatchers {
		for _, m := range lm.FindAll(example) {
			fmt.Println(string(example[m.Start:m.End]))
		}
	}
	// Output:
	// anchor-target
	// img-source
	// img-source-set-2x
	// img-source-set-3x
	// img-source-set-wide
}

func fixture() *builder {
	domain := "https://kask.ufukty.com/"
	rw := rewriter.New(paths.Paths{Src: ".", Dst: ".", Url: "/"})
	rw.Bank("a/b/README.md", domain+"a/b/")
	rw.Bank("a/b/page.md", domain+"a/b/page.html")
	rw.Bank("a/index.tmpl", domain+"a/")
	// visitable dirs:
	rw.Bank("a/", domain+"a/")
	rw.Bank("a/b", domain+"a/b/")
	// assets
	rw.Bank("a/.assets/img.jpg", domain+"a/.assets/img.jpg")
	rw.Bank("a/.assets/img@2x.jpg", domain+"a/.assets/img%402x.jpg")
	rw.Bank("a/.assets/img@3x.jpg", domain+"a/.assets/img%403x.jpg")
	return &builder{rw: rw}
}

// TODO: add asset linking. eg. <img>
func TestBuilder_htmlContent(t *testing.T) {
	type tc struct {
		input, expected string
	}
	tcs := map[string]tc{
		"anchor href with redundant traverse": {
			input:    `<a href="../a/b/README.md#Title"></a>`,
			expected: `<a href="https://kask.ufukty.com/a/b/#Title"></a>`,
		},
		"anchor href with anchor target": {
			input:    `<a href="../a/b/page.md#Title"></a>`,
			expected: `<a href="https://kask.ufukty.com/a/b/page.html#Title"></a>`,
		},
		"anchor href to index page and anchor target": {
			input:    `<a href="../a/index.tmpl#Title"></a>`,
			expected: `<a href="https://kask.ufukty.com/a/#Title"></a>`,
		},
		"img src and srcset": {
			input:    `<img src=".assets/img.jpg" srcset=".assets/img@2x.jpg 2x, .assets/img@3x.jpg 3x">`,
			expected: `<img src="https://kask.ufukty.com/a/.assets/img.jpg" srcset="https://kask.ufukty.com/a/.assets/img%402x.jpg 2x, https://kask.ufukty.com/a/.assets/img%403x.jpg 3x">`,
		},
	}
	b := fixture()
	linker := paths.Paths{Src: "a/page.tmpl", Dst: "a/page.html", Url: "/a/page.html"}
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			got, err := b.htmlContent(linker, []byte(tc.input))
			if err != nil {
				t.Fatalf("act, unexpected error: %v", err)
			}
			if tc.expected != string(got) {
				t.Errorf("assert, expected:\n  %s, got:\n  %s", tc.expected, string(got))
			}
		})
	}
}
