package rewriter

import (
	"fmt"
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
	if !strings.HasSuffix(contentDir.Url, "/") {
		contentDir.Url += "/"
	}
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

// Idempotent.
// Input can be absolute or relative local-path of the target.
// Return value is absolute and encoded URL.
func (rw Rewriter) Rewrite(linked string, linker paths.Paths) (string, error) {
	if rw.isExternal(linked) {
		return linked, nil
	}

	linkedS := rw.split(linked)

	if has(rw.targets, linkedS.base+linkedS.ref) { // idempotency
		return linked, nil
	}

	var base string
	if linkedS.base == "" && linkedS.ref == "" { // same page
		base = linker.Src
	} else if linkedS.base == "" { // same dir
		base = filepath.Dir(linker.Src)
	} else { // absolute
		base = "."
	}

	resolved := filepath.Join(base, linkedS.ref)

	if strings.HasPrefix(resolved, "..") { // warn escaping
		return "", ErrOutsideTarget
	}

	rewritten, ok := rw.links[resolved]
	if !ok {
		return "", ErrInvalidTarget
	}

	// TODO: validate asset existence
	// TODO: validate anchor target existence (via ToC)
	return rewritten + linkedS.tail, nil
}
