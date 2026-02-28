package visitor

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/builder/markdown/visitor/codefence"
	"go.ufukty.com/kask/internal/paths"
)

// unsafe for concurrency
type Visitor struct {
	cf   *codefence.Renderer
	Page paths.Paths
}

func New() *Visitor {
	return &Visitor{cf: codefence.NewRenderer()}
}

func (v *Visitor) Prepare(page paths.Paths) {
	v.Page = page
}

// TODO: rewrite links inside the html blocks
func (v *Visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		return v.cf.RenderNodeHook(w, node, entering)
	}
	return ast.GoToNext, false
}
