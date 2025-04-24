package markdown

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

type visitor struct {
	cf *codefenceRenderer
}

func (v visitor) hookLink(node *ast.Link) {
	var dest = string(node.Destination)
	if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") || strings.HasPrefix(dest, "/") {
		return
	}
	// dest = filepath.Clean(filepath.Join(r.folderpath, dest))
	// fmt.Println(dest)
}

func (v visitor) Hook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		return v.cf.RenderNodeHook(w, node, entering)
	case *ast.Image:
		// TODO: change destination

	case *ast.Link:
		v.hookLink(node)
		// default:
		// fmt.Println(reflect.TypeOf(node))
	}
	return ast.GoToNext, false
}
