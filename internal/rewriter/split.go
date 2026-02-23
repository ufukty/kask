package rewriter

import "strings"

func before(s string, substr string) int {
	i := strings.Index(s, substr)
	if i == -1 {
		return len(s)
	}
	return i
}

func after(s, substr string) int {
	i := strings.Index(s, substr)
	if i == -1 {
		return 0
	}
	return i + len(substr)
}

func between(s string, before, after int) string {
	if before < after {
		return s[before:after]
	}
	return ""
}

func (rw Rewriter) split(url string) (string, string, string, string) {
	var (
		aDomain = after(url, rw.contentDir.Url)
		bAssets = before(url, ".assets")
		bAnchor = before(url, "#")
		bQuery  = before(url, "?")
	)
	var (
		domain = between(url, 0, aDomain)
		path   = between(url, aDomain, min(bAssets, bAnchor, bQuery))
		assets = between(url, bAssets, min(bQuery, bAnchor))
		tail   = between(url, min(bQuery, bAnchor), len(url))
	)
	return domain, path, assets, tail
}
