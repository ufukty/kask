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

func unescape(target string) string {
	t2, err := url.PathUnescape(target)
	if err != nil {
		return target
	}
	return t2
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

func splitQuery(path string) (string, string) {
	i := max(strings.Index(path, "#"), strings.Index(path, "?"))
	if i == -1 {
		i = len(path)
	}
	return path[:i], path[i:]
}

// returns the absolute URL for the linked resource in content directory
func (rw Rewriter) locateByContentDir(linked string, linker paths.Paths) (string, bool) {
	if linked == "" { // same-page anchor links
		linked = linker.Src
	} else if !filepath.IsAbs(linked) {
		linked = filepath.Join(filepath.Dir(linker.Src), linked)
	} else {
		linked = strings.TrimPrefix(linked, "/")
	}
	dst, ok := rw.links[linked]
	if ok {
		return dst, true
	}
	return "", false
}

// returns the absolute URL for the linked resource in sitemap
func (rw Rewriter) locateByUrl(linked string, linker paths.Paths) (string, bool) {
	lr, err := url.Parse(linker.Url)
	if err != nil {
		return "", false
	}
	ld, err := url.Parse(linked)
	if err != nil {
		return "", false
	}
	ld = lr.ResolveReference(ld)
	_, ok := rw.targets[ld.String()]
	if ok {
		return ld.String(), true
	}
	return "", false
}

func (rw Rewriter) locate(linked string, linker paths.Paths) (string, bool) {
	byContentDir, ok := rw.locateByContentDir(linked, linker)
	if ok {
		return byContentDir, true
	}
	byBuildDir, ok := rw.locateByUrl(linked, linker)
	if ok {
		return byBuildDir, true
	}
	return "", false
}

func isEvil(linked string, linker paths.Paths) bool {
	return strings.HasPrefix(filepath.Join(filepath.Dir(linker.Src), linked), "..")
}

// Rewriter returns an absolute and encoded URL for a resource which the user linked it either:
//   - by its path in the content or build directory;
//   - with its absolute path (by the content directory root) or relative (to the linker page);
//   - with/out url encoded path segments.
func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if isExternal(linked) {
		return linked, nil
	}
	if isEvil(linked, linker) {
		return "", ErrInvalidTarget
	}
	linked = unescape(linked)
	linked, query := splitQuery(linked)
	dst, ok := rw.locate(linked, linker)
	if !ok {
		return "", ErrInvalidTarget
	}
	return dst + query, nil
}
