package visitor

import (
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/builder/markdown/visitor/codefence"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
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

func (v *Visitor) rewrite(dst []byte) []byte {
	dst2, err := v.rw.Rewrite(string(dst), v.Page)
	if err == rewriter.ErrInvalidTarget {
		v.InvTargets = append(v.InvTargets, string(dst))
		return dst
	}
	return []byte(dst2)
}

// TODO: rewrite links inside the html blocks
func (v *Visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		return v.cf.RenderNodeHook(w, node, entering)
	case *ast.Image:
		if entering {
			node.Destination = v.rewrite(node.Destination)
		}
	case *ast.Link:
		if entering {
			node.Destination = v.rewrite(node.Destination)
		}
	}
	return ast.GoToNext, false
}
