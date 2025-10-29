package hook

import (
	"path/filepath"
	"strings"
)

func isExternal(url string) bool {
	return false ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")
}

func isRelative(url string) bool {
	return !strings.HasPrefix(url, "/")
}

func has(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}

func rewrite(url, currentdir string, rewrites map[string]string) string {
	if isExternal(url) {
		return url
	}
	if isRelative(url) {
		url = filepath.Join(currentdir, url)
	}
	url = filepath.Clean(url)
	if url == "." {
		url = ""
	}
	if isRelative(url) {
		url = "/" + url
	}
	if has(rewrites, url) {
		url = rewrites[url]
	}
	return url
}
