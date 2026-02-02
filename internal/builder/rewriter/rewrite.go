package rewriter

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

func New() *Rewriter {
	return &Rewriter{
		links: map[string]string{},
	}
}

func (rw *Rewriter) Bank(src, dst string) {
	rw.links[src] = dst
}

// TODO: check path handling with query and browser parameters
func (rw Rewriter) Rewrite(dst, src string) string {
	if isExternal(dst) {
		return dst
	}
	dst = unescape(dst)
	if isRelative(dst) {
		dst = filepath.Join(src, dst)
	}
	dst = filepath.Clean(dst)
	if dst == "." {
		dst = ""
	}
	if isRelative(dst) {
		dst = "/" + dst
	}
	if has(rw.links, dst) {
		dst = rw.links[dst]
	}
	if isDir(dst) && !strings.HasSuffix(dst, "/") {
		dst = dst + "/"
	}
	return dst
}
