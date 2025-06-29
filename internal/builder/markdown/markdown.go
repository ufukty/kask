package markdown

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/ufukty/kask/internal/builder/markdown/hook"
)

type Page struct {
	Content string
	Toc     *TocNode
}

func ToHtml(root, page string) (*Page, error) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.Attributes |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock |
			parser.Mmark |
			parser.MathJax,
	)
	c, err := os.ReadFile(filepath.Join(root, page))
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	n := p.Parse(c).(*ast.Document)

	r := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: hook.NewVisitor(page).Visit,
	})
	h := markdown.Render(n, r)
	toc := getTableOfContent(n, r)
	return &Page{string(h), toc}, nil
}
