package builder

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/ufukty/kask/internal/builder/rewriter"
)

var anchorHref = regexp.MustCompile(`<a[^>]*href="([^"]*)"[^>]*>[^<]*</a>`)

func rewriteLinksInHtmlPage(rw *rewriter.Rewriter, page string, bs []byte) ([]byte, error) {
	delta := 0
	for _, matches := range anchorHref.FindAllSubmatchIndex(bs, -1) {
		if len(matches) < 4 {
			continue
		}
		start, end := matches[2]+delta, matches[3]+delta
		m1 := bs[start:end]
		n, err := rw.Rewrite(string(m1), page)
		if err != nil {
			return nil, fmt.Errorf("found the link %q: %w", string(m1), err)
		}
		m2 := []byte(n)
		bs = slices.Replace(bs, start, end, m2...)
		delta += len(m2) - len(m1)
	}
	return bs, nil
}
