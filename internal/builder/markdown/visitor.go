package markdown

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/paths"
)

// unsafe for concurrency
type visitor struct {
	cf   *codefenceRenderer
	Page paths.Paths
}

func newVisitor() *visitor {
	return &visitor{cf: newCodefenceRenderer()}
}

func (v *visitor) Prepare(page paths.Paths) {
	v.Page = page
}

func (v *visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		if entering {
			return v.cf.RenderNodeHook(w, node, entering)
		}
	}
	return ast.GoToNext, false
}
