package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/ufukty/kask/internal/builder/markdown/hook"
	"github.com/ufukty/kask/internal/builder/rewriter"
	"github.com/ufukty/kask/pkg/kask"
)

func serializeInvalidTargets(ts []string) string {
	for i := range len(ts) {
		ts[i] = fmt.Sprintf("%q", ts[i])
	}
	return strings.Join(ts, ", ")
}

func ToHtml(root, page string, rw *rewriter.Rewriter) (*kask.Markdown, error) {
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

	v := hook.NewVisitor(page, rw)
	r := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: v.Visit,
	})
	h := markdown.Render(n, r)
	if len(v.InvTargets) > 0 {
		return nil, fmt.Errorf("found links to invalid target(s): %s", serializeInvalidTargets(v.InvTargets))
	}
	toc := getTableOfContent(n, r)
	return &kask.Markdown{string(h), toc}, nil
}
