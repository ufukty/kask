package builder

import (
	"net/url"
	"path/filepath"
	"strings"
)

func withStripping(path string, toStrip bool) string {
	if toStrip {
		return stripOrdering(path)
	}
	return path
}

func uri(parent, child string, isDir bool) string {
	parent = strings.TrimSuffix(parent, "/")
	uri := url.PathEscape(child)
	if isDir && !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	if parent != "" {
		uri = parent + "/" + uri
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

func (p paths) sub(basename string, isToStrip bool) paths {
	strpd := withStripping(basename, isToStrip)
	return paths{
		src: filepath.Join(p.src, basename),
		dst: filepath.Join(p.dst, strpd),
		url: uri(p.url, strpd, true),
	}
}
