package rewriter

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"go.ufukty.com/kask/internal/paths"
)

var (
	ErrInvalidTarget = fmt.Errorf("invalid link: internal target doesn't exist.")
	ErrOutsideTarget = fmt.Errorf("outside target: link escapes the content directory.")
)

func has[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func isExternal(url string) bool {
	return false ||
		strings.HasPrefix(url, "http://") ||
		strings.HasPrefix(url, "https://")
}

type Rewriter struct {
	links   map[string]string // src -> url
	targets map[string]any    // urls
}

func New() *Rewriter {
	return &Rewriter{
		links:   map[string]string{},
		targets: map[string]any{},
	}
}

func (rw *Rewriter) Bank(src, url string) {
	rw.links[src] = url
	rw.targets[url] = nil
}

func split(path string) (string, string, string) {
	assets := strings.Index(path, ".assets")
	anchor := strings.Index(path, "#")
	query := strings.Index(path, "?")
	tail := -1
	if anchor != -1 && query != -1 {
		tail = min(anchor, query)
	} else if anchor != -1 || query != -1 {
		tail = max(anchor, query)
	}
	if assets == -1 && tail == -1 {
		return path, "", ""
	} else if assets == -1 {
		return path[:tail], "", path[tail:]
	} else if tail == -1 {
		return path[:assets], path[assets:], ""
	} else {
		return path[:assets], path[assets:tail], path[tail:]
	}
}

// TODO: take as parameter to support domains
var contentDirectory = paths.Paths{
	Src: ".",
	Dst: ".",
	Url: "/",
}

// toRelative rewrites absolute paths as if they're relative to the root
func toRelative(linked, linker, root string) (string, string) {
	if filepath.IsAbs(linked) {
		return strings.TrimPrefix(linked, "/"), root
	}
	return linked, linker
}

func joinSrcPaths(dst, src string) string {
	if dst == "" { // same-page anchor links
		return src
	}
	return filepath.Join(filepath.Dir(src), dst)
}

func (rw Rewriter) rewriteByContentDir(linked string, linker string) (string, bool, error) {
	linked, assets, query := split(linked)
	linked, linker = toRelative(linked, linker, contentDirectory.Src)
	linked = joinSrcPaths(linked, linker)
	if strings.HasPrefix(linked, "..") {
		return "", false, ErrOutsideTarget
	}
	dst, ok := rw.links[linked]
	return dst + assets + query, ok, nil
}

// RFC 3986 Section 5.2
func join3986(dst, src string) (string, error) {
	d, err := url.Parse(dst)
	if err != nil {
		return "", fmt.Errorf("parsing destination url: %w", err)
	}
	s, err := url.Parse(src)
	if err != nil {
		return "", fmt.Errorf("parsing source url: %w", err)
	}
	return s.ResolveReference(d).String(), nil
}

func (rw Rewriter) canonicalizeIfUrl(linked, linker string) (string, bool, error) {
	linked, assets, query := split(linked)
	linked, linker = toRelative(linked, linker, contentDirectory.Url)
	linked, err := join3986(linked, linker)
	if err != nil {
		return "", false, fmt.Errorf("join: %w", err)
	}
	return linked + assets + query, has(rw.targets, linked), nil
}

// Idempotent.
// Input can be absolute/relative local-path/URL of target.
// Return value is absolute and encoded URL.
func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if isExternal(linked) {
		return linked, nil
	}
	dst, ok, err := rw.rewriteByContentDir(linked, linker.Src)
	if err != nil {
		return "", fmt.Errorf("rewriting by content directory: %w", err)
	} else if ok {
		return dst, nil
	}
	dst, ok, err = rw.canonicalizeIfUrl(linked, linker.Url) // for idempotency
	if err != nil {
		return "", fmt.Errorf("canonicalizing: %w", err)
	} else if ok {
		return dst, nil
	}
	return dst, ErrInvalidTarget
}
