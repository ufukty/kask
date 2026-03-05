package builder

import (
	"errors"
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
	rw := rewriter.New(paths.Paths{Src: ".", Dst: ".", Url: "https://kask.ufukty.com/"})
	m := map[string]string{
		// leaves
		"a/page.tmpl":   "https://kask.ufukty.com/a/page.html",
		"a/index.tmpl":  "https://kask.ufukty.com/a/",
		"a/b/README.md": "https://kask.ufukty.com/a/b/",
		"a/b/page.md":   "https://kask.ufukty.com/a/b/page.html",

		// visitable dirs:
		".":   "https://kask.ufukty.com/",
		"a/":  "https://kask.ufukty.com/a/",
		"a/b": "https://kask.ufukty.com/a/b/",

		// assets
		".assets/font.woff2":             "https://kask.ufukty.com/.assets/font.woff2",
		"a/.assets/img.jpg":              "https://kask.ufukty.com/a/.assets/img.jpg",
		"a/.assets/img@2x.jpg":           "https://kask.ufukty.com/a/.assets/img%402x.jpg",
		"a/.assets/img@3x.jpg":           "https://kask.ufukty.com/a/.assets/img%403x.jpg",
		"a/.assets/poster.jpg":           "https://kask.ufukty.com/a/.assets/poster.jpg",
		"a/.assets/video.mp4":            "https://kask.ufukty.com/a/.assets/video.mp4",
		"a/.assets/og.jpg":               "https://kask.ufukty.com/a/.assets/og.jpg",
		"a/.assets/embedded-player.html": "https://kask.ufukty.com/a/.assets/embedded-player.html",
	}
	for s, u := range m {
		rw.Bank(s, u)
	}
	return &builder{rw: rw, incorrectLinks: map[string][]string{}}
}

func TestBuilder_htmlPostProcess(t *testing.T) {
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
		"link href to domain": {
			input:    `<link rel="canonical" href="/" />`,
			expected: `<link rel="canonical" href="https://kask.ufukty.com/" />`,
		},
		"link href to asset with additional attributes": {
			input:    `<link rel="preload" href="/.assets/font.woff2" as="font" />`,
			expected: `<link rel="preload" href="https://kask.ufukty.com/.assets/font.woff2" as="font" />`,
		},
		"meta tag for og:image": {
			input:    `<meta property="og:image" content=".assets/og.jpg" />`,
			expected: `<meta property="og:image" content="https://kask.ufukty.com/a/.assets/og.jpg" />`,
		},
		"meta tag for og:url": {
			input:    `<meta property="og:url" content="" />`,
			expected: `<meta property="og:url" content="https://kask.ufukty.com/a/page.html" />`,
		},
		"meta tag for twitter:image": {
			input:    `<meta name="twitter:image" content=".assets/og.jpg" />`,
			expected: `<meta name="twitter:image" content="https://kask.ufukty.com/a/.assets/og.jpg" />`,
		},
		"meta tag for twitter:url": {
			input:    `<meta name="twitter:url" content="" />`,
			expected: `<meta name="twitter:url" content="https://kask.ufukty.com/a/page.html" />`,
		},
		"iframe src": {
			input:    `<iframe src=".assets/embedded-player.html"></iframe>`,
			expected: `<iframe src="https://kask.ufukty.com/a/.assets/embedded-player.html"></iframe>`,
		},
		"video and source": {
			input:    `<video poster=".assets/poster.jpg"><source src=".assets/video.mp4" type="video/mp4"></video>`,
			expected: `<video poster="https://kask.ufukty.com/a/.assets/poster.jpg"><source src="https://kask.ufukty.com/a/.assets/video.mp4" type="video/mp4"></video>`,
		},
	}
	b := fixture()
	linker := paths.Paths{Src: "a/page.tmpl", Dst: "a/page.html", Url: "/a/page.html"}
	for tn, tc := range tcs {
		t.Run(tn, func(t *testing.T) {
			got, err := b.htmlPostProcess(linker, []byte(tc.input))
			if err != nil {
				if errors.Is(err, ErrIncorrectLinks) {
					t.Errorf("act, unexpected error: %v", err)
				} else {
					t.Fatalf("act, unexpected error: %v", err)
				}
			}
			if tc.expected != string(got) {
				t.Errorf("assert,\nexpected: %s\ngot:      %s", tc.expected, string(got))
			}
		})
	}
}
