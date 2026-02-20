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
	cutoff := max(strings.Index(path, "#"), strings.Index(path, "?"))
	if cutoff == -1 {
		cutoff = len(path)
	}
	url, query := path[:cutoff], path[cutoff:]
	if url == "" {
		url = "."
	}
	return path, query
}

func assureAbsolute(dst, src string) string {
	if filepath.IsAbs(dst) {
		return dst
	}
	return filepath.Join(filepath.Dir(src), dst)
}

// returns the absolute URL for the linked resource in content directory
func (rw Rewriter) locateByContentDir(linked string, linker paths.Paths) (string, bool) {
	if !filepath.IsAbs(linked) {
		linked = filepath.Join(linker.Src, linked)
	}
	dst, ok := rw.links[linked]
	return dst, ok
}

// returns the absolute URL for the linked resource in build directory
func (rw Rewriter) locateByBuildDir(linked string, linker paths.Paths) (string, bool) {
	if !filepath.IsAbs(linked) {
		linked = filepath.Join(linker.Dst, linked)
	}
	_, ok := rw.targets[linked]
	if ok {
		return linked, true
	}
	return "", false
}

func (rw Rewriter) locate(linked string, linker paths.Paths) (string, bool) {
	byContentDir, ok := rw.locateByContentDir(linked, linker)
	if ok {
		return byContentDir, true
	}
	byBuildDir, ok := rw.locateByBuildDir(linked, linker)
	if ok {
		return byBuildDir, true
	}
	return "", false
}

// user can link a page by
//   - its path in local content directory
//   - its path in local content directory, with encoding special characters
//   - its path in build directory
//   - its path in build directory, with encoding special characters
//
// either relative to the linker page, or in absolute form. using leading slash for content directory root.
// at either combination, [Rewriter.Rewrite] creates encoded absolute URLs (except the domain).
// The bank may contain a rule for the parent or ancestor directory of linked page instead of the direct file.
func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if isExternal(linked) {
		return linked, nil
	}
	linked = unescape(linked)
	linked, query := splitQuery(linked)
	dst, ok := rw.locate(linked, linker)
	if !ok {
		return "", ErrInvalidTarget
	}
	return dst + query, nil
}
