package markdown

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func ToHtml(src string) (string, *TocNode, error) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock |
			parser.Mmark |
			parser.MathJax,
	)
	c, err := os.ReadFile(src)
	if err != nil {
		return "", nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	n := p.Parse(c).(*ast.Document)

	v := visitor{
		cf: newCodefenceRenderer(),
	}
	r := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: v.Hook,
	})
	h := markdown.Render(n, r)
	toc := getTableOfContent(n, r)
	return string(h), toc, nil
}
