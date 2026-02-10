package hook

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/ufukty/kask/internal/builder/markdown/hook/codefence"
	"github.com/ufukty/kask/internal/builder/rewriter"
)

type visitor struct {
	cf      *codefence.Renderer
	pagedir string
	rw      *rewriter.Rewriter
}

func NewVisitor(page string, rw *rewriter.Rewriter) *visitor {
	return &visitor{
		cf:      codefence.NewRenderer(),
		pagedir: page,
		rw:      rw,
	}
}

func (v visitor) links(node *ast.Link) (ast.WalkStatus, bool) {
	node.Destination = []byte(v.rw.Rewrite(string(node.Destination), v.pagedir))
	return ast.GoToNext, false
}

func (v visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		return v.cf.RenderNodeHook(w, node, entering)
	case *ast.Image:
		// TODO: change destination

	case *ast.Link:
		return v.links(node)
	}
	return ast.GoToNext, false
}
