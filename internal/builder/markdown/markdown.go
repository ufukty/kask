package markdown

import (
	"fmt"
	"io/fs"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/pkg/kask"
)

type Renderer struct {
	src      fs.FS
	renderer *html.Renderer
	visitor  *visitor
}

func New(src fs.FS, domain string) *Renderer {
	v := newVisitor(domain)
	r := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: v.Visit,
	})
	return &Renderer{
		src:      src,
		visitor:  v,
		renderer: r,
	}
}

func (r Renderer) ToHtml(page paths.Paths) (*kask.Markdown, error) {
	c, err := fs.ReadFile(r.src, page.Src)
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
	toc := r.getTableOfContent(n)
	m := &kask.Markdown{
		Content: string(html),
		Toc:     toc,
	}
	return m, nil
}
