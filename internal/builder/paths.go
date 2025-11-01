package builder

import (
	"net/url"
	"path/filepath"
)

type paths struct {
	src string
	dst string // stripped ordering
	url string // path encoded
}

func (p paths) withChild(item string, isToStrip bool) paths {
	dst := item
	if isToStrip {
		dst = stripOrdering(dst)
	}
	return paths{
		src: filepath.Join(p.src, item),
		dst: filepath.Join(p.dst, dst),
		url: filepath.Join(p.url, url.PathEscape(dst)),
	}
}
