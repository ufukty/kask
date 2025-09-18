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

func titleFromFilename(src, ext string) string {
	return titler.String(strings.TrimSuffix(stripOrdering(filepath.Base(src)), ext))
}

func hrefFromFilename(dstPathEncoded, filename string) string {
	base := stripOrdering(filename)
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	base = url.PathEscape(base)
	return "/" + filepath.Join(dstPathEncoded, base+".html")
}

func targetFromFilename(dst, dstPath, filename string) string {
	base := stripOrdering(filename)
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	return filepath.Join(dst, dstPath, base+".html")
}

var titleSelector = regexp.MustCompile(`^# (.*)$`)

func titleFromMarkdownPage(src string) (string, error) {
	bs, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	ss := titleSelector.FindSubmatch(bs)
	if len(ss) == 0 {
		return "", nil
	}
	return string(ss[1]), nil
}

func titleForMarkdownPage(src string) (string, error) {
	title, err := titleFromMarkdownPage(src)
	if err != nil {
		return "", fmt.Errorf("from page: %w", err)
	}
	if title != "" {
		return title, nil
	}
	return titleFromFilename(src, ".md"), nil
}
