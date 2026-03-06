package paths

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

type UrlMode int

const (
	UrlModeDefault = UrlMode(iota) // eg. "/", "/dir/", "/page.html"
	UrlModeExtless                 // eg. "/", "/dir/", "/page"
)

var orderingPrefixMatcher = regexp.MustCompile(`^\d+\s*[\-.,]*\s*`)

func stripOrdering(s string) string {
	p := orderingPrefixMatcher.FindString(s)
	if len(p) == len(s) {
		return s
	}
	return strings.TrimPrefix(s, p)
}

func withStripping(path string, toStrip bool) string {
	if toStrip {
		return stripOrdering(path)
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
		child = strings.TrimSuffix(child, ext)
		child = withStripping(child, strip)
		child += ".html"
		dst := filepath.Join(parent, child)
		return dst
	}
}

func dirUri(parent, child string, strip bool) string {
	child = withStripping(child, strip)
	child = url.PathEscape(child)
	uri, _ := url.JoinPath(parent, child)
	uri = assureTrailingSlash(uri)
	return uri
}

func assureExtension(path string, um UrlMode) string {
	if um == UrlModeExtless {
		return path
	}
	return path + ".html"
}

func fileUri(parent, child string, strip bool, um UrlMode) string {
	if child == "README.md" || child == "index.tmpl" {
		return parent
	} else {
		ext := filepath.Ext(child)
		child = strings.TrimSuffix(child, ext)
		child = withStripping(child, strip)
		child = url.PathEscape(child)
		child = assureExtension(child, um)
		uri, _ := url.JoinPath(parent, child)
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

func (p Paths) File(basename string, strip bool, um UrlMode) Paths {
	cleaned := strings.TrimPrefix(basename, ".")
	return Paths{
		Src: filepath.Join(p.Src, basename),
		Dst: fileDst(p.Dst, cleaned, strip),
		Url: fileUri(p.Url, cleaned, strip, um),
	}
}

func (p Paths) AssetDir(basename string) Paths {
	return Paths{
		Src: filepath.Join(p.Src, basename),
		Dst: dirDst(p.Dst, basename, false),
		Url: dirUri(p.Url, basename, false),
	}
}

func (p Paths) AssetFile(basename string) Paths {
	u, _ := url.JoinPath(p.Url, basename)
	ss := strings.Split(u, "/")
	if len(ss) > 0 {
		ss[len(ss)-1] = strings.ReplaceAll(ss[len(ss)-1], "@", "%40")
	}
	u = strings.Join(ss, "/")
	return Paths{
		Src: filepath.Join(p.Src, basename),
		Dst: filepath.Join(p.Dst, basename),
		Url: u,
	}
}

func cssBundleName(propagated bool) string {
	if propagated {
		return "styles.propagate.css"
	} else {
		return "styles.css"
	}
}

func (p Paths) Stylesheet(propagated bool) Paths {
	base := cssBundleName(propagated)
	u, _ := url.JoinPath(p.Url, base)
	return Paths{
		Src: filepath.Join(p.Src, base),
		Dst: filepath.Join(p.Dst, base),
		Url: u,
	}
}
