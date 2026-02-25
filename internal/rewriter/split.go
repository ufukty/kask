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
	base string // domain or slash
	ref  string
	tail string
}

func (rw Rewriter) split(path string) splits {
	var (
		aDomain = cmp.Or(afterPrefix(path, rw.contentDir.Url), afterPrefix(path, "/"))
		bAnchor = before(path, "#")
		bQuery  = before(path, "?")
	)
	return splits{
		base: between(path, 0, aDomain),
		ref:  between(path, aDomain, min(bAnchor, bQuery)),
		tail: between(path, min(bQuery, bAnchor), len(path)),
	}
}
