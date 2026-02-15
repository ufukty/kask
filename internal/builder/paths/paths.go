package paths

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

var orderingStripper = regexp.MustCompile(`^(\d+[\-., ]*)?(.*)$`)

func stripOrdering(s string) string {
	return orderingStripper.FindStringSubmatch(s)[2]
}

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

type Paths struct {
	Src string
	Dst string // stripped ordering
	Url string // path encoded
}

func (p Paths) Subdir(basename string, strip bool) Paths {
	return Paths{
		Src: filepath.Join(p.Src, basename),
		Dst: dirDst(p.Dst, basename, strip),
		Url: dirUri(p.Url, basename, strip),
	}
}

func (p Paths) File(basename string, strip bool) Paths {
	return Paths{
		Src: filepath.Join(p.Src, basename),
		Dst: fileDst(p.Dst, basename, strip),
		Url: fileUri(p.Url, basename, strip),
	}
}
