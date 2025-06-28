package hook

import (
	"io"
	"path/filepath"

	"github.com/gomarkdown/markdown/ast"
	"github.com/ufukty/kask/internal/builder/markdown/hook/codefence"
)

type visitor struct {
	cf      *codefence.Renderer
	pageDir string
}

func NewVisitor(page string) *visitor {
	return &visitor{
		cf:      codefence.NewRenderer(),
		pageDir: filepath.Dir(page),
	}
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
