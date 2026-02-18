package visitor

import (
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/builder/markdown/visitor/codefence"
	"go.ufukty.com/kask/internal/builder/paths"
	"go.ufukty.com/kask/internal/builder/rewriter"
)

// unsafe for concurrency
type Visitor struct {
	rw *rewriter.Rewriter
	cf *codefence.Renderer

	InvTargets []string
	Page       paths.Paths
}

func New(rw *rewriter.Rewriter) *Visitor {
	return &Visitor{
		rw: rw,
		cf: codefence.NewRenderer(),
	}
}

func (v *Visitor) Prepare(page paths.Paths) {
	v.Page = page
	v.InvTargets = []string{}
}

func serializeInvalidTargets(ts []string) string {
	for i := range len(ts) {
		ts[i] = fmt.Sprintf("%q", ts[i])
	}
	return strings.Join(ts, ", ")
}

func (v *Visitor) Error() error {
	if len(v.InvTargets) > 0 {
		return fmt.Errorf("%s", serializeInvalidTargets(v.InvTargets))
	}
	return nil
}

func (v *Visitor) links(node *ast.Link, entering bool) (ast.WalkStatus, bool) {
	if !entering {
		return ast.GoToNext, false
	}
	h2, err := v.rw.Rewrite(string(node.Destination), v.Page.Src)
	if err == rewriter.ErrInvalidTarget {
		v.InvTargets = append(v.InvTargets, string(node.Destination))
		return ast.GoToNext, false
	}
	node.Destination = []byte(h2)
	return ast.GoToNext, false
}

func (v *Visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
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
