package builder

import (
	"bytes"
	"fmt"
	"html/template"
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
	base := filename
	if strippedOrdering {
		base = stripOrdering(base)
	}
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	base = url.PathEscape(base)
	return "/" + filepath.Join(dstPathEncoded, base+".html")
}

func targetFromFilename(dst, folderpath, filename string, strippedOrdering bool) string {
	base := filename
	if strippedOrdering {
		base = stripOrdering(base)
	}
	base = strings.TrimSuffix(base, filepath.Ext(filename))
	return filepath.Join(dst, folderpath, base+".html")
}

var regexpMarkdown = regexp.MustCompile(`(?m)^#\s+(.+)$`)

type extractor struct{}

func (e extractor) FromWeb(path string) (string, error) {
	tmpl, err := template.New("").ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if tmpl.Lookup("title") == nil {
		return "", nil
	}
	b := bytes.NewBufferString("")
	if err = tmpl.ExecuteTemplate(b, "title", nil); err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}
	return b.String(), nil
}

func (e extractor) FromMarkdown(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	ms := regexpMarkdown.FindStringSubmatch(string(f))
	if len(ms) < 2 {
		return "", nil
	}
	return ms[1], nil
}

func (e extractor) FromFile(path string) (string, error) {
	switch ext := filepath.Ext(path); ext {
	case ".tmpl":
		p, err := e.FromWeb(path)
		if err != nil {
			return "", fmt.Errorf("markdown: %w", err)
		}
		return p, nil
	case ".md":
		p, err := e.FromMarkdown(path)
		if err != nil {
			return "", fmt.Errorf("web: %w", err)
		}
		return p, nil
	default:
		return "", fmt.Errorf("unknown file extension: %s", ext)
	}
}

var theExtractor = extractor{}

// 1. title from content, if available
// 2. title from file name, if visitable
// 3. title from folder name
func decideOnTitle(src, ext string, strippedOrdering bool) (string, error) {
	title, err := theExtractor.FromFile(src)
	if err != nil {
		return "", fmt.Errorf("reading: %w", err)
	}
	if title != "" {
		return title, nil
	}
	return titleFromFilename(src, ext, strippedOrdering), nil
}
