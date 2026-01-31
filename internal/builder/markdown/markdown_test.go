package markdown

import (
	"fmt"
	"testing"

	"github.com/ufukty/kask/internal/builder/markdown/hook"
)

func Test_FromFile(t *testing.T) {
	p, err := ToHtml("testdata", "input.md", hook.NewRewriter())
	if err != nil {
		t.Fatal(fmt.Errorf("act, ToHtml: %w", err))
	}
	fmt.Println(p.Content)
	fmt.Println(p.Toc)
}
