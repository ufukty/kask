package builder

import (
	"net/url"
	"path/filepath"
	"strings"
)

func uri(dst, dir string) string {
	dst = strings.TrimSuffix(dst, "/")
	uri := url.PathEscape(dir)
	if !strings.HasSuffix(uri, "/") {
		uri = uri + "/"
	}
	if dst != "" {
		uri = dst + "/" + uri
	}
	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}
	if uri == "/./" {
		uri = "/"
	}
	return uri
}

type paths struct {
	src string
	dst string // stripped ordering
	url string // path encoded
}

func (p paths) subdir(dir string, isToStrip bool) paths {
	dirStripped := dir
	if isToStrip {
		dirStripped = stripOrdering(dirStripped)
	}
	return paths{
		src: filepath.Join(p.src, dir),
		dst: filepath.Join(p.dst, dirStripped),
		url: uri(p.url, dirStripped),
	}
}
