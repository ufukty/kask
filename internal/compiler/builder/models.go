package builder

import (
	"time"

	"github.com/ufukty/kask/internal/compiler/builder/directory"
	"github.com/ufukty/kask/internal/compiler/builder/markdown"
)

type Data struct {
	Stylesheets     []string
	Node            *directory.Node
	WebSiteRoot     *directory.Node
	MarkdownContent string
	MarkdownTOC     *markdown.TocNode
	Time            time.Time
	Dir             *directory.Dir
}
