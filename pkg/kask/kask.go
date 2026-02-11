package kask

import (
	"time"

	"github.com/ufukty/kask/internal/builder/markdown"
)

// Node represents a sitemap node which can be either of:
//   - Unvisitable folder
//   - Visitable folder (has `index.tmpl` or `README.md` file)
//   - Page of a `.tmpl` or `.md` file
type Node struct {
	Title    string // Sourced either from the file, meta file or file/folder name
	Href     string // Visitable when filled
	Parent   *Node
	Children []*Node
}

// TemplateContent provides the dynamic information a template
// may need to render a page.
type TemplateContent struct {
	Stylesheets []string       // Needs to be linked in the page.
	Node        *Node          // Represents the page currently rendered.
	Root        *Node          // Represents site root. Used for rendering sitemap.
	Markdown    *markdown.Page // Rendered HTML for Markdown page and the Table-of-Contents
	Time        time.Time      // The time first rendering started
}
