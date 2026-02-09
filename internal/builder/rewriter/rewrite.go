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

func splitQuery(path string) (string, string) {
	query := max(strings.Index(path, "#"), strings.Index(path, "?"))
	if query != -1 {
		return path[:query], path[query:]
	}
	return path, ""
}

func assureAbsolute(cwd, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(cwd, path)
}

func (rw Rewriter) Rewrite(dst, src string) string {
	if isExternal(dst) {
		return dst
	}
	dst = unescape(dst)
	dst = strings.TrimSuffix(dst, "/")
	dst = assureAbsolute(filepath.Dir(src), dst)
	dst = filepath.Clean(dst)
	dst, query := splitQuery(dst)
	if has(rw.links, dst) {
		dst = rw.links[dst]
	}
	return dst + query
}
