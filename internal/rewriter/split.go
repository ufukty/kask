package rewriter

import (
	"cmp"
	"strings"
)

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
	prePath string // domain or slash
	path    string
	assets  string
	tail    string
}

func (rw Rewriter) split(path string) splits {
	var (
		aDomain = cmp.Or(afterPrefix(path, rw.contentDir.Url), afterPrefix(path, "/"))
		bAssets = before(path, ".assets")
		bAnchor = before(path, "#")
		bQuery  = before(path, "?")
	)
	return splits{
		prePath: between(path, 0, aDomain),
		path:    between(path, aDomain, min(bAssets, bAnchor, bQuery)),
		assets:  between(path, bAssets, min(bQuery, bAnchor)),
		tail:    between(path, min(bQuery, bAnchor), len(path)),
	}
}
