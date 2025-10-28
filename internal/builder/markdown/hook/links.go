package hook

import (
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

func (v visitor) links(node *ast.Link) (ast.WalkStatus, bool) {
	dest := string(node.Destination)

	isExternal := false ||
		strings.HasPrefix(dest, "http://") ||
		strings.HasPrefix(dest, "https://") ||
		strings.HasPrefix(dest, "/")
	if isExternal {
		return ast.GoToNext, false
	}

	dest = strings.TrimSuffix(dest, "README.md")
	dest = strings.TrimSuffix(dest, "index.tmpl")

	// TODO: absolute paths => trim the prefix WORKING DIR in the PATH (?)
	dest = filepath.Clean(filepath.Join(v.pagedir, dest))

	if rewritten, ok := v.links_[dest]; ok {
		dest = rewritten
	}

	if dest == "." {
		dest = "/"
	}
	node.Destination = []byte(dest)

	return ast.GoToNext, false
}
