package markdown

import (
	"fmt"
	"testing"
)

func Test_FromFile(t *testing.T) {
	h, toc, err := ToHtml("testdata/input.md")
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	fmt.Println(h)
	fmt.Println(toc)
}
