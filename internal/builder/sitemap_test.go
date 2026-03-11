package builder

import (
	"os"
	"testing"

	"go.ufukty.com/kask/pkg/kask"
)

func TestWriteSitemap(t *testing.T) {
	var (
		c    = &kask.Node{Href: "/c.html"}
		b    = &kask.Node{Href: "/a/b.html"}
		a    = &kask.Node{Href: "/a", Children: []*kask.Node{b}}
		root = &kask.Node{Href: "/", Children: []*kask.Node{a, c}}
	)
	err := writeSitemap(os.Stdout, root)
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
}
