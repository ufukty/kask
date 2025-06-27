package hook

import (
	"cmp"
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

func links(node *ast.Link) {
	dest := string(node.Destination)
	isExternal := cmp.Or(
		strings.HasPrefix(dest, "http://"),
		strings.HasPrefix(dest, "https://"),
		strings.HasPrefix(dest, "/"),
	)
	if isExternal {
		return
	}
	// dest = filepath.Clean(filepath.Join(r.folderpath, dest))
	// fmt.Println(dest)
}
