package builder

import (
	"fmt"

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

// TODO: add asset linking. eg. <img>
func Example_builder_rewriteLinksInHtmlPage() {
	input := `<a href="../a/b/README.md#Title">link with redundant traverse</a>
<a href="../a/b/page.md#Title">link with redundant traverse</a>
<a href="../a/index.tmpl#Title">link with redundant traverse</a>`

	rw := rewriter.New(paths.Paths{Src: ".", Dst: ".", Url: "/"})
	rw.Bank("a/b/README.md", "/a/b/")
	rw.Bank("a/b/page.md", "/a/b/page.html")
	rw.Bank("a/index.tmpl", "/a/")
	// visitable dirs:
	rw.Bank("a/", "/a/")
	rw.Bank("a/b", "/a/b/")
	b := builder{rw: rw}

	page := paths.Paths{Src: "a/page.tmpl", Dst: "a/page.html", Url: "/a/page.html"}
	new, _ := b.htmlContent(page, []byte(input))
	fmt.Println(string(new))
	// Output:
	// <a href="/a/b/#Title">link with redundant traverse</a>
	// <a href="/a/b/page.html#Title">link with redundant traverse</a>
	// <a href="/a/#Title">link with redundant traverse</a>
}
