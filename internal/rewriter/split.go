package rewriter

import "strings"

func before(s string, substr string) int {
	i := strings.Index(s, substr)
	if i == -1 {
		return len(s)
	}
	return i
}

func afterPrefix(s, domain string) int {
	if strings.HasPrefix(s, domain) {
		return len(domain)
	}
	return 0
}

func between(s string, before, after int) string {
	if before < after {
		return s[before:after]
	}
	return ""
}

type splits struct {
	domain string
	path   string
	assets string
	tail   string
}

func (rw Rewriter) split(url string) splits {
	var (
		aDomain = afterPrefix(url, rw.contentDir.Url)
		bAssets = before(url, ".assets")
		bAnchor = before(url, "#")
		bQuery  = before(url, "?")
	)
	return splits{
		domain: between(url, 0, aDomain),
		path:   between(url, aDomain, min(bAssets, bAnchor, bQuery)),
		assets: between(url, bAssets, min(bQuery, bAnchor)),
		tail:   between(url, min(bQuery, bAnchor), len(url)),
	}
}
