package hook

import (
	"net/url"
	"path/filepath"
	"strings"
)

func isExternal(url string) bool {
	return false ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")
}

func isRelative(url string) bool {
	return !filepath.IsAbs(url)
}

func has(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}

func unescape(target string) string {
	t2, err := url.PathUnescape(target)
	if err != nil {
		return target
	}
	return t2
}

func isDir(target string) bool {
	return !strings.Contains(filepath.Base(target), ".")
}

// TODO: check path handling with query and browser parameters
func rewrite(target, currentdir string, rewrites map[string]string) string {
	if isExternal(target) {
		return target
	}
	target = unescape(target)
	isDir := isDir(target)
	if isRelative(target) {
		target = filepath.Join(currentdir, target)
	}
	target = filepath.Clean(target)
	if target == "." {
		target = ""
	}
	if isRelative(target) {
		target = "/" + target
	}
	if has(rewrites, target) {
		target = rewrites[target]
	}
	if isDir && !strings.HasSuffix(target, "/") {
		target = target + "/"
	}
	return target
}
