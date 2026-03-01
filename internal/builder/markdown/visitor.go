package markdown

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/internal/paths"
)

// unsafe for concurrency
type visitor struct {
	cf     *codefenceRenderer
	Page   paths.Paths
	domain string
}

func newVisitor(domain string) *visitor {
	return &visitor{
		cf:     newCodefenceRenderer(),
		domain: domain,
	}
}

func (v *visitor) Prepare(page paths.Paths) {
	v.Page = page
}

func (v *visitor) isExternal(dst string) bool {
	if strings.HasPrefix(dst, "http://") || strings.HasPrefix(dst, "https://") {
		return !strings.HasPrefix(dst, v.domain)
	}
	return false
}

func (v *visitor) Visit(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		if entering {
			return v.cf.RenderNodeHook(w, node, entering)
		}
	case *ast.Link:
		if entering {
			if v.isExternal(string(node.Destination)) {
				node.AdditionalAttributes = append(node.AdditionalAttributes, `target="_blank"`)
			}
		}
	}
	return ast.GoToNext, false
}
