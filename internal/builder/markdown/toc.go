package markdown

import (
	"bytes"
	"fmt"

	"github.com/gomarkdown/markdown/ast"
	"go.ufukty.com/kask/pkg/kask"
)

func (r Renderer) renderHeadingText(h *ast.Heading) string {
	var buf bytes.Buffer
	for _, child := range h.Children {
		err := r.renderer.RenderNode(&buf, child, true)
		if err != ast.Terminate {
			continue
		}
	}
	return buf.String()
}

func (r Renderer) getTableOfContent(doc *ast.Document) *kask.MarkdownTocNode {
	root := &kask.MarkdownTocNode{Title: "root", Level: 0}
	stack := []*kask.MarkdownTocNode{root}
	headingCount := 0

	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		h, ok := node.(*ast.Heading)
		if ok && entering && !h.IsTitleblock {
			if h.HeadingID == "" {
				h.HeadingID = fmt.Sprintf("toc%d", headingCount)
			}
			headingCount++

			title := r.renderHeadingText(h)
			newNode := &kask.MarkdownTocNode{
				Title: title,
				ID:    h.HeadingID,
				Level: h.Level,
			}

			for len(stack) > 1 && stack[len(stack)-1].Level >= h.Level {
				stack = stack[:len(stack)-1]
			}

			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, newNode)
			stack = append(stack, newNode)
		}
		return ast.GoToNext
	})

	return root
}
