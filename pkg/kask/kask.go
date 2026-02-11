package kask

import (
	"time"

	"github.com/ufukty/kask/internal/builder/markdown"
)

// represents a sitemap node which can be either of:
//   - non-visitable directories
//   - directories with "index.tmpl" or "README.md" file
//   - pages corresponding to .tmpl or .md files
type Node struct {
	Title    string // markdown h1, meta.yml title or the file name
	Href     string // Visitable when filled
	Parent   *Node
	Children []*Node
}

// template files should access necessary information through
// the fields of this struct
type TemplateContent struct {
	Stylesheets []string
	Node, Root  *Node
	Markdown    *markdown.Page
	Time        time.Time
}
