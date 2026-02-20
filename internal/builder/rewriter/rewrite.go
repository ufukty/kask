package rewriter

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"go.ufukty.com/kask/internal/builder/paths"
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

func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if isExternal(linked) || strings.HasPrefix(linked, "#") || strings.HasPrefix(linked, "?") {
		return linked, nil
	}
	linked = unescape(linked)
	if dst, ok := rw.isValid(linked, linker.Src); ok {
		return dst, nil
	}
	linked = assureAbsolute(linked, linker.Src)
	linked = filepath.Clean(linked)
	linked, query := splitQuery(linked)
	linked = strings.TrimPrefix(linked, "/")
	linked = strings.TrimSuffix(linked, "/")
	if linked == "" {
		linked = "."
	}
	if !has(rw.links, linked) {
		return "", ErrInvalidTarget
	}
	linked = rw.links[linked]
	return linked + query, nil
}
