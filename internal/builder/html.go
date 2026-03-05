package builder

import (
	"fmt"
	"maps"
	"slices"
	"strconv"

	"go.ufukty.com/gommons/pkg/tree"
	"go.ufukty.com/kask/internal/builder/narrowing"
	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/internal/rewriter"
)

// TODO: add support for `<video>` tags
var linkMatchers = []narrowing.Matchers{
	narrowing.MustCompile(`<a[^>]*>[^<]*</a>`, `href="([^"]*)"`),                                    // a[href]
	narrowing.MustCompile(`<img[^>]*/?>`, `src="([^"]*)"`),                                          // img[src]
	narrowing.MustCompile(`<img[^>]*/?>`, `srcset="\s*([^"]*)\s*"`, `([^\s]+)\s+\d+(?:\.\d+)?[wx]`), // img[srcset]
	narrowing.MustCompile(`<link[^>]*/?>`, `href="([^"]*)"`),                                        // link[href]
	narrowing.MustCompile(`<meta[^>]* property="og:image"[^>]*/?>`, `content="([^"]*)"`),            // meta[property="og:image"]
	narrowing.MustCompile(`<meta[^>]* property="og:url"[^>]*/?>`, `content="([^"]*)"`),              // meta[property="og:url"]
	narrowing.MustCompile(`<meta[^>]* name="twitter:image"[^>]*/?>`, `content="([^"]*)"`),           // meta[name="twitter:image"]
	narrowing.MustCompile(`<meta[^>]* name="twitter:url"[^>]*/?>`, `content="([^"]*)"`),             // meta[name="twitter:url"]
}

var ErrIncorrectLinks = fmt.Errorf("found incorrect links")

func (b *builder) rewriteLinksInRanges(ranges []narrowing.Range, page paths.Paths, bs []byte) ([]byte, error) {
	slices.SortFunc(ranges, narrowing.Compare)
	invTargets := []string{}
	delta := 0
	for _, r := range ranges {
		start, end := r.Start+delta, r.End+delta
		m1 := bs[start:end]
		n, err := b.rw.Rewrite(string(m1), page)
		if err == rewriter.ErrInvalidTarget {
			invTargets = append(invTargets, strconv.Quote(string(m1)))
			continue
		} else if err != nil {
			return nil, fmt.Errorf("rewriting link %q: %w", string(m1), err)
		}
		m2 := []byte(n)
		bs = slices.Replace(bs, start, end, m2...)
		delta += len(m2) - len(m1)
	}
	if len(invTargets) > 0 {
		b.incorrectLinks[page.Src] = invTargets
		return nil, ErrIncorrectLinks
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

func (b *builder) htmlPostProcess(page paths.Paths, bs []byte) ([]byte, error) {
	bs, err := b.rewriteLinksInRanges(findLinkRanges(bs), page, bs)
	if err != nil {
		return nil, fmt.Errorf("rewriting links: %w", err)
	}
	return bs, nil
}

func (b *builder) reportIncorrectLinks() {
	for _, page := range slices.Sorted(maps.Keys(b.incorrectLinks)) {
		fmt.Println(tree.List(page, b.incorrectLinks[page]))
	}
}
