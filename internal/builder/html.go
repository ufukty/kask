package builder

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
)

var anchorHref = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>[^<]*</a>`)

func rewriteLinksInHtmlPage(rw *rewriter.Rewriter, page paths.Paths, bs []byte) ([]byte, error) {
	invTargets := []string{}
	delta := 0
	for _, matches := range anchorHref.FindAllSubmatchIndex(bs, -1) {
		if len(matches) < 4 {
			continue
		}
		start, end := matches[2]+delta, matches[3]+delta
		m1 := bs[start:end]
		n, err := rw.Rewrite(string(m1), page)
		if err == rewriter.ErrInvalidTarget {
			invTargets = append(invTargets, fmt.Sprintf("%q", string(m1)))
			continue
		} else if err != nil {
			return nil, fmt.Errorf("rewriting link %q: %w", string(m1), err)
		}
		m2 := []byte(n)
		bs = slices.Replace(bs, start, end, m2...)
		delta += len(m2) - len(m1)
	}
	if len(invTargets) > 0 {
		return nil, fmt.Errorf("found links to invalid target(s): %s", strings.Join(invTargets, ", "))
	}
	return bs, nil
}
