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

func getBestLexerFor(source, lang []byte) chroma.Lexer {
	var langstr = string(lang)
	var lexer = lexers.Get(langstr)
	if lexer == nil {
		lexer = lexers.Analyse(string(source))
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	return chroma.Coalesce(lexer)
}

func (r *Renderer) htmlHighlight(w io.Writer, lexer chroma.Lexer, source []byte) error {
	var iterator, err = lexer.Tokenise(nil, string(source))
	if err != nil {
		return fmt.Errorf("calling tokenise: %w", err)
	}
	err = r.formatter.Format(w, r.style, iterator)
	if err != nil {
		return fmt.Errorf("calling formatter: %w", err)
	}
	return nil
}

func (r *Renderer) RenderNodeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		var lexer = getBestLexerFor(code.Literal, code.Info)
		r.htmlHighlight(w, lexer, code.Literal)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
