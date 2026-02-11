package kask

import (
	"time"
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

// MarkdownTocNode represents a node in the table of contents tree.
type MarkdownTocNode struct {
	Title    string             // A Markdown page section title
	ID       string             // The identifier can be used for creating anchor links
	Level    int                // Section level, `3` => `###`
	Children []*MarkdownTocNode // Subsections
}

// Markdown contains all the information for a Markdown based page
type Markdown struct {
	Content string           // Rendered HTML for content
	Toc     *MarkdownTocNode // Table-of-Contents root. Start printing ToC from root's children.
}

// TemplateContent provides the dynamic information a template
// may need to render a page.
type TemplateContent struct {
	Stylesheets []string // Needs to be linked in the page.
	Node        *Node    // Represents the page currently rendered.
	Root        *Node    // Represents site root. Used for rendering sitemap.
	Markdown    *Markdown
	Time        time.Time // The time first rendering started
}
