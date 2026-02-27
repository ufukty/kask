package builder

import (
	"fmt"
	"slices"
	"strings"

	"go.ufukty.com/kask/internal/builder/narrowing"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
)

// TODO: add support for `<video>` tags
var linkMatchers = []narrowing.Matchers{
	narrowing.MustCompile(`<a[^>]*>[^<]*</a>`, `href="([^"]*)"`),                                    // <a href=
	narrowing.MustCompile(`<img[^>]*/?>`, `src="([^"]*)"`),                                          // <img src=
	narrowing.MustCompile(`<img[^>]*/?>`, `srcset="\s*([^"]*)\s*"`, `([^\s]+)\s+\d+(?:\.\d+)?[wx]`), // <img srcset=
}

func (b *builder) rewriteLinksInRanges(ranges []narrowing.Range, page paths.Paths, bs []byte) ([]byte, error) {
	invTargets := []string{}
	delta := 0
	for _, r := range ranges {
		start, end := r.Start+delta, r.End+delta
		m1 := bs[start:end]
		n, err := b.rw.Rewrite(string(m1), page)
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

func findLinkRanges(bs []byte) []narrowing.Range {
	ranges := []narrowing.Range{}
	for _, lm := range linkMatchers {
		ranges = append(ranges, lm.FindAll(bs)...)
	}
	return ranges
}

func (b *builder) htmlContent(page paths.Paths, bs []byte) ([]byte, error) {
	return b.rewriteLinksInRanges(findLinkRanges(bs), page, bs)
}
