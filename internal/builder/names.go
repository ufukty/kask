package builder

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var orderingStripper = regexp.MustCompile(`^(\d+[\-., ]*)?(.*)$`)

func stripOrdering(s string) string {
	return orderingStripper.FindStringSubmatch(s)[2]
}

var titler = cases.Title(language.Und, cases.NoLower)

func titleFromFilename(base, ext string, strippedOrdering bool) string {
	base = filepath.Base(base)
	if strippedOrdering {
		base = stripOrdering(base)
	}
	return titler.String(strings.TrimSuffix(base, ext))
}

func hrefFromFilename(dstPathEncoded, filename string, strippedOrdering bool) string {
	base := ""
	if strippedOrdering {
		base = stripOrdering(filename)
	}
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	base = url.PathEscape(base)
	return "/" + filepath.Join(dstPathEncoded, base+".html")
}

func targetFromFilename(dst, dstPath, filename string) string {
	base := stripOrdering(filename)
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	return filepath.Join(dst, dstPath, base+".html")
}

var titleExtractors = map[string]*regexp.Regexp{
	".md":   regexp.MustCompile(`(?m)^#\s+(.+)$`),
	".tmpl": regexp.MustCompile(`(?i)<title>(.*?)</title>`),
}

func titleFromContent(content, ext string) string {
	extractor, ok := titleExtractors[ext]
	if !ok {
		return ""
	}
	submatches := extractor.FindStringSubmatch(content)
	if len(submatches) < 2 {
		return ""
	}
	return submatches[1]
}

// 1. title from content, if available
// 2. title from file name, if visitable
// 3. title from folder name
func decideOnTitle(src, ext string, strippedOrdering bool) (string, error) {
	bs, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	title := titleFromContent(string(bs), ext)
	if title != "" {
		return title, nil
	}
	return titleFromFilename(src, ext, strippedOrdering), nil
}
