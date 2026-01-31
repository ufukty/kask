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

type Rewriter struct {
	links map[string]string // src path -> uri
}

func NewRewriter() *Rewriter {
	return &Rewriter{
		links: map[string]string{},
	}
}

func (rw *Rewriter) Bank(src, dst string) {
	rw.links[src] = dst
}

// TODO: check path handling with query and browser parameters
func (rw Rewriter) rewrite(target, currentdir string) string {
	if isExternal(target) {
		return target
	}
	target = unescape(target)
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
	if has(rw.links, target) {
		target = rw.links[target]
	}
	if isDir(target) && !strings.HasSuffix(target, "/") {
		target = target + "/"
	}
	return target
}
