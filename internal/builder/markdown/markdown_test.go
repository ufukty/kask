package markdown

import (
	"fmt"
	"testing"

	"go.ufukty.com/kask/internal/builder/rewriter"
)

func Test_FromFile(t *testing.T) {
	r := rewriter.New()
	r.Bank(".assets/img.jpg", "/.assets/img.jpg")
	r.Bank("sibling.md", "/sibling.html")
	p, err := ToHtml("testdata", "input.md", r)
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	fmt.Println(p.Content)
	fmt.Println(p.Toc)
}
