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

func (rw Rewriter) isExternal(url string) bool {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return !strings.HasPrefix(url, rw.contentDir.Url)
	}
	return false
}

type Rewriter struct {
	links      map[string]string // src -> url
	targets    map[string]any    // urls
	contentDir paths.Paths
}

func New(contentDir paths.Paths) *Rewriter {
	return &Rewriter{
		links:      map[string]string{},
		targets:    map[string]any{},
		contentDir: contentDir,
	}
}

func (rw *Rewriter) Bank(src, url string) {
	rw.links[src] = url
	rw.targets[url] = nil
}

// toRelative rewrites absolute paths as if they're relative to the root
func toRelative(linked, linker, root string) (string, string) {
	if filepath.IsAbs(linked) {
		return strings.TrimPrefix(linked, "/"), root
	}
	return linked, linker
}

func join(a, b string) string {
	if strings.HasSuffix(a, "/") || strings.HasPrefix(b, "/") {
		return a + b
	}
	return a + "/" + b
}

func joinSrcPaths(dst, src string) string {
	if dst == "" { // same-page anchor links
		return src
	}
	return join(filepath.Dir(src), dst)
}

func (rw Rewriter) rewriteByContentDir(linked string, linker string) (string, bool, error) {
	domain, linked, assets, query := rw.split(linked)
	linked, linker = toRelative(linked, linker, rw.contentDir.Src)
	linked = joinSrcPaths(linked, linker)
	if strings.HasPrefix(linked, "..") {
		return "", false, ErrOutsideTarget
	}
	dst, ok := rw.links[linked]
	return domain + dst + assets + query, ok, nil
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
	domain, linked, assets, query := rw.split(linked)
	linked, linker = toRelative(linked, linker, rw.contentDir.Url)
	linked, err := join3986(linked, linker)
	if err != nil {
		return "", false, fmt.Errorf("join: %w", err)
	}
	return domain + linked + assets + query, has(rw.targets, linked), nil
}

// Idempotent.
// Input can be absolute/relative local-path/URL of target.
// Return value is absolute and encoded URL.
func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if rw.isExternal(linked) {
		return linked, nil
	}
	dst, ok, err := rw.rewriteByContentDir(linked, linker.Src)
	if err != nil {
		return "", fmt.Errorf("rewriting by content directory: %w", err)
	} else if ok {
		return dst, nil
	}
	dst, ok, err = rw.canonicalizeIfUrl(linked, linker.Url) // for idempotency + allow writing url-like links
	if err != nil {
		return "", fmt.Errorf("canonicalizing: %w", err)
	} else if ok {
		return dst, nil
	}
	return dst, ErrInvalidTarget
}
