package visitor

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/paths"
)

// unsafe for concurrency
type Visitor struct {
	cf   *codefenceRenderer
	Page paths.Paths
}

func New() *Visitor {
	return &Visitor{cf: newCodefenceRenderer()}
}

func (v *Visitor) Prepare(page paths.Paths) {
	v.Page = page
}

func (v *Visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		if entering {
			return v.cf.RenderNodeHook(w, node, entering)
		}
	}
	return ast.GoToNext, false
}
