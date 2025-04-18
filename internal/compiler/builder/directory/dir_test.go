package directory

import (
	"fmt"
	"testing"
)

func spaces(depth int) string {
	s := ""
	for i := 0; i < depth; i++ {
		s += "    "
	}
	return s
}

func nodetree(n *Node, depth int) {
	fmt.Printf("%s%s (PageType:%s) (Visitable:%t) (Self:%p) (Parent:%p) (Subpages:%d) (Subsection:%d) (SrcFilename:%s) (TargetInSitePath:%s)\n", spaces(depth),
		n.Title,
		n.PageType,
		n.Visitable,
		n,
		n.Parent,
		len(n.Subpages),
		len(n.Subsections),
		n.SrcFilename,
		n.TargetInSitePath,
	)
	for _, s := range n.Subpages {
		nodetree(s, depth+1)
	}
	for _, s := range n.Subsections {
		nodetree(s, depth+1)
	}
}

func Test_Inspect(t *testing.T) {
	root, err := Inspect("testdata/acme")
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
	nodetree(root.Node, 0)
	fmt.Println(root)
}
