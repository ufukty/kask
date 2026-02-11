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

// .  => /
// a  => /a
// /a => /a
func assureLeadingSlash(path string) string {
	if path == "." {
		return "/"
	} else if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

// .  => /
// a  => a/
// a/ => a/
func assureTrailingSlash(path string) string {
	if path == "." {
		return "/"
	} else if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}

func dirDst(parent, child string, strip bool) string {
	child = withStripping(child, strip)
	dst := filepath.Join(parent, child)
	return dst
}

func fileDst(parent, child string, strip bool) string {
	if child == "README.md" || child == "index.tmpl" {
		return filepath.Join(parent, "index.html")
	} else {
		ext := filepath.Ext(child)
		child = strings.TrimSuffix(child, ext) + ".html"
		child = withStripping(child, strip)
		dst := filepath.Join(parent, child)
		return dst
	}
}

func dirUri(parent, child string, strip bool) string {
	child = withStripping(child, strip)
	child = url.PathEscape(child)
	uri := filepath.Join(parent, child)
	uri = assureLeadingSlash(uri)
	uri = assureTrailingSlash(uri)
	return uri
}

func fileUri(parent, child string, strip bool) string {
	if child == "README.md" || child == "index.tmpl" {
		return parent
	} else {
		ext := filepath.Ext(child)
		child = strings.TrimSuffix(child, ext)
		child = withStripping(child, strip)
		child = url.PathEscape(child)
		child = child + ".html"
		uri := filepath.Join(parent, child)
		uri = assureLeadingSlash(uri)
		return uri
	}
}

type paths struct {
	src string
	dst string // stripped ordering
	url string // path encoded
}

func (p paths) subdir(basename string, strip bool) paths {
	return paths{
		src: filepath.Join(p.src, basename),
		dst: dirDst(p.dst, basename, strip),
		url: dirUri(p.url, basename, strip),
	}
}

func (p paths) file(basename string, strip bool) paths {
	return paths{
		src: filepath.Join(p.src, basename),
		dst: fileDst(p.dst, basename, strip),
		url: fileUri(p.url, basename, strip),
	}
}
