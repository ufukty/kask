package hook

import (
	"path/filepath"
	"strings"
)

func rewrite(url, pagedir string, rewrites map[string]string) string {
	isExternal := false ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://") ||
		strings.HasPrefix(url, "/")
	if isExternal {
		return url
	}

	url = strings.TrimSuffix(url, "README.md")
	url = strings.TrimSuffix(url, "index.tmpl")

	// TODO: absolute paths => trim the prefix WORKING DIR in the PATH (?)
	url = filepath.Clean(filepath.Join(pagedir, url))

	if rewritten, ok := rewrites[url]; ok {
		url = rewritten
	}

	if url == "." {
		url = "/"
	}

	return url
}
