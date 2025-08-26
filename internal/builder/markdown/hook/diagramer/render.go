package diagramer

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/ufukty/diagramer/pkg/sequence/parser"
)

func Render(w io.Writer, node *ast.CodeBlock, entering bool) (ast.WalkStatus, bool) {
	parser.Reader(node.Literal)
}
