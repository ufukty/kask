package markdown

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.ufukty.com/kask/internal/builder/markdown/visitor"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
	"go.ufukty.com/kask/pkg/kask"
)

type Renderer struct {
	src      string
	renderer *html.Renderer
	visitor  *visitor.Visitor
}

func New(src string, rw *rewriter.Rewriter) *Renderer {
	v := visitor.New(rw)
	r := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: v.Visit,
	})
	return &Renderer{
		src:      src,
		visitor:  v,
		renderer: r,
	}
}

func (r Renderer) ToHtml(page paths.Paths) (*kask.Markdown, error) {
	c, err := os.ReadFile(filepath.Join(r.src, page.Src))
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}
	r.visitor.Prepare(page)
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.Attributes |
			parser.AutoHeadingIDs |
			parser.NoEmptyLineBeforeBlock |
			parser.Mmark |
			parser.MathJax,
	)
	n := p.Parse(c).(*ast.Document)
	html := markdown.Render(n, r.renderer)
	err = r.visitor.Error()
	if err != nil {
		return nil, fmt.Errorf("found links to invalid target(s): %s", err)
	}
	toc := r.getTableOfContent(n)
	m := &kask.Markdown{
		Content: string(html),
		Toc:     toc,
	}
	return m, nil
}
