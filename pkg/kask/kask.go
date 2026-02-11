// Package kask provides documentation for user-facing Kask code.
// The information here are supposed to be accessed by Kask users in their templates.
package kask

import (
	"time"
)

// Node represents a sitemap node which can be either of:
//   - Unvisitable folder
//   - Visitable folder (has `index.tmpl` or `README.md` file)
//   - Page of a `.tmpl` or `.md` file
//
// # Printing sitemaps
//
// Use a recursive template that calls itself on each children:
//
//	{{define "sitemap"}}
//	  <li>
//	    <a href="{{.Href}}">{{.Title}}</a>
//		  {{with .Children}}
//		  <ul>
//		    {{range .}}
//		    {{template "sitemap" .}}
//		    {{end}}
//		  </ul>
//		  {{end}}
//	  </li>
//	{{end}}
//
// # Printing breadcrumbs
//
// Use a recursive template but invoke it with the [TemplateContent.Node],
// instead of [TemplateContent.Root] and call itself on the parent until
// parent is `nil`.
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
//
// # Linking stylesheets
//
// TemplateContent contains the links for smallest set of Css bundles the
// current page needs to link. Since Kask bundles stylesheets at each level
// of directory structure separately, the pages need to link stylesheets
// created at parent sections when available.
//
//	{{range .Stylesheets}}
//	<link rel="stylesheet" href="{{.}}">
//	{{end}}
type TemplateContent struct {
	Stylesheets []string // Needs to be linked in the page.
	Node        *Node    // Represents the page currently rendered.
	Root        *Node    // Represents site root. Used for rendering sitemap.
	Markdown    *Markdown
	Time        time.Time // The time first rendering started
}
