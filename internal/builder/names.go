package builder

import (
	"bytes"
	"fmt"
	"html/template"
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

var regexpMarkdown = regexp.MustCompile(`(?m)^#\s+(.+)$`)

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

var titler = cases.Title(language.Und, cases.NoLower)

func pageTitleFromFilename(base string) string {
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return titler.String(base)
}

func pageTitle(src string, p paths) (string, error) {
	title, err := theExtractor.FromFile(filepath.Join(src, p.src))
	if err != nil {
		return "", fmt.Errorf("extracting from file: %w", err)
	}
	if title != "" {
		return title, nil
	}
	return pageTitleFromFilename(filepath.Base(p.dst)), nil
}
