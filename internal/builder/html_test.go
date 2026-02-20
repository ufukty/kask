package builder

import (
	"fmt"

	"go.ufukty.com/kask/internal/builder/paths"
	"go.ufukty.com/kask/internal/builder/rewriter"
)

func Example_builder_rewriteLinksInHtmlPage() {
	input := `<a href="../a/b/README.md#Title">link with redundant traverse</a>
<a href="../a/b/page.md#Title">link with redundant traverse</a>
<a href="../a/index.tmpl#Title">link with redundant traverse</a>`
	rw := rewriter.New()
	rw.Bank("a/b/README.md", "/a/b/")
	rw.Bank("a/b/page.md", "/a/b/page.html")
	rw.Bank("a/index.tmpl", "/a/")
	// visitable dirs:
	rw.Bank("a/", "/a/")
	rw.Bank("a/b", "/a/b/")
	page := paths.Paths{Src: "a/page.tmpl", Dst: "a/page.html", Url: "/a/page.html"}
	new, _ := rewriteLinksInHtmlPage(rw, page, []byte(input))
	fmt.Println(string(new))
	// Output:
	// <a href="/a/b/#Title">link with redundant traverse</a>
	// <a href="/a/b/page.html#Title">link with redundant traverse</a>
	// <a href="/a/#Title">link with redundant traverse</a>
}
