package hook

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/builder/markdown/hook/codefence"
	"go.ufukty.com/kask/internal/builder/paths"
	"go.ufukty.com/kask/internal/builder/rewriter"
)

type visitor struct {
	page       paths.Paths
	rw         *rewriter.Rewriter
	cf         *codefence.Renderer
	InvTargets []string
}

func NewVisitor(page paths.Paths, rw *rewriter.Rewriter) *visitor {
	return &visitor{
		page:       page,
		rw:         rw,
		cf:         codefence.NewRenderer(), // TODO: reuse
		InvTargets: []string{},
	}
}

func (v *visitor) links(node *ast.Link, entering bool) (ast.WalkStatus, bool) {
	if !entering {
		return ast.GoToNext, false
	}
	h2, err := v.rw.Rewrite(string(node.Destination), v.page.Src)
	if err == rewriter.ErrInvalidTarget {
		v.InvTargets = append(v.InvTargets, string(node.Destination))
		return ast.GoToNext, false
	}
	node.Destination = []byte(h2)
	return ast.GoToNext, false
}

func (v *visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		return v.cf.RenderNodeHook(w, node, entering)
	case *ast.Image:
		// TODO: change destination

	case *ast.Link:
		return v.links(node, entering)
	}
	return ast.GoToNext, false
}
