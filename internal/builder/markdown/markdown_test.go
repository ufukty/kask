package markdown

import (
	"fmt"
	"testing"
)

func Test_FromFile(t *testing.T) {
	p, err := ToHtml("testdata", "input.md", map[string]string{})
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	fmt.Println(p.Content)
	fmt.Println(p.Toc)
}
