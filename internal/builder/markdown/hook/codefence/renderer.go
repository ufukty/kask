package codefence

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type Renderer struct {
	formatter *html.Formatter
	style     *chroma.Style
}

func NewRenderer() *Renderer {
	return &Renderer{
		formatter: html.New(html.WithClasses(true), html.TabWidth(4)),
		style:     styles.Get("monokailight"),
	}
}

func getBestLexerForLang(source, lang []byte) chroma.Lexer {
	if l := lexers.Get(string(lang)); l != nil {
		return l
	}
	if l := lexers.Analyse(string(source)); l != nil {
		return l
	}
	return lexers.Fallback
}

func (r *Renderer) htmlHighlight(w io.Writer, lexer chroma.Lexer, source []byte) error {
	iterator, err := lexer.Tokenise(nil, string(source))
	if err != nil {
		return fmt.Errorf("calling tokenise: %w", err)
	}
	err = r.formatter.Format(w, r.style, iterator)
	if err != nil {
		return fmt.Errorf("calling formatter: %w", err)
	}
	return nil
}

func (r *Renderer) RenderNodeHook(w io.Writer, node *ast.CodeBlock, entering bool) (ast.WalkStatus, bool) {
	lexer := chroma.Coalesce(getBestLexerForLang(node.Literal, node.Info))
	r.htmlHighlight(w, lexer, node.Literal)
	return ast.GoToNext, true
}
