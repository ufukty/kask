package rewriter

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

var ErrInvalidTarget = fmt.Errorf("invalid link: internal target doesn't exist.")

func isExternal(url string) bool {
	return false ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")
}

func has[C comparable, V any](m map[C]V, k C) bool {
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
	links   map[string]string // src path -> url
	targets map[string]any    // urls
}

func New() *Rewriter {
	return &Rewriter{
		links:   map[string]string{},
		targets: map[string]any{},
	}
}

func (rw *Rewriter) Bank(src, dst string) {
	rw.links[src] = dst
	rw.targets[dst] = nil
}

func splitQuery(path string) (string, string) {
	query := max(strings.Index(path, "#"), strings.Index(path, "?"))
	if query != -1 {
		return path[:query], path[query:]
	}
	return path, ""
}

func assureAbsolute(dst, src string) string {
	if filepath.IsAbs(dst) {
		return dst
	}
	return filepath.Join(filepath.Dir(src), dst)
}

func assureLeadingSlash(path string) string {
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func (rw Rewriter) isValid(dst, src string) (string, bool) {
	src = assureLeadingSlash(src)
	dst = assureAbsolute(dst, src)
	dst = filepath.Clean(dst)
	return dst, has(rw.targets, dst)
}

func (rw Rewriter) Rewrite(dst, src string) (string, error) {
	if isExternal(dst) || strings.HasPrefix(dst, "#") || strings.HasPrefix(dst, "?") {
		return dst, nil
	}
	dst = unescape(dst)
	if dst, ok := rw.isValid(dst, src); ok {
		return dst, nil
	}
	dst = assureAbsolute(dst, src)
	dst = filepath.Clean(dst)
	dst, query := splitQuery(dst)
	dst = strings.TrimPrefix(dst, "/")
	dst = strings.TrimSuffix(dst, "/")
	if dst == "" {
		dst = "."
	}
	if !has(rw.links, dst) {
		return "", ErrInvalidTarget
	}
	dst = rw.links[dst]
	return dst + query, nil
}
