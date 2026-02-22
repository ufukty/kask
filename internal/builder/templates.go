package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go.ufukty.com/kask/internal/paths"
	"go.ufukty.com/kask/pkg/kask"
)

func newTemplateSet() *template.Template {
	tmpl := template.New("page")
	tmpl.Funcs(template.FuncMap{
		"trustedCss":      func(s string) template.CSS { return template.CSS(s) },
		"trustedHtml":     func(s string) template.HTML { return template.HTML(s) },
		"trustedHtmlAttr": func(s string) template.HTMLAttr { return template.HTMLAttr(s) },
		"trustedJs":       func(s string) template.JS { return template.JS(s) },
		"trustedJsStr":    func(s string) template.JSStr { return template.JSStr(s) },
		"trustedSrcSet":   func(s string) template.Srcset { return template.Srcset(s) },
		"trustedUrl":      func(s string) template.URL { return template.URL(s) },
	})
	return tmpl
}

func (b *builder) propagateTemplates(d *dir2, toPropagate *template.Template) error {
	var err error
	if d.original.Kask != nil && d.original.Kask.Propagate != nil && len(d.original.Kask.Propagate.Tmpl) > 0 {
		toPropagate, err = toPropagate.ParseFiles(d.original.Kask.Propagate.Tmpl...)
		if err != nil {
			return fmt.Errorf("parsing to-propagate template files: %w", err)
		}
	}
	atLevel, err := toPropagate.Clone()
	if err != nil {
		return fmt.Errorf("cloning propagated: %w", err)
	}
	if d.original.Kask != nil && len(d.original.Kask.Tmpl) > 0 {
		atLevel, err = atLevel.ParseFiles(d.original.Kask.Tmpl...)
		if err != nil {
			return fmt.Errorf("parsing at-level template files: %w", err)
		}
	}
	d.templates = atLevel
	for _, subdir := range d.subdirs {
		if err := b.propagateTemplates(subdir, toPropagate); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.Src), err)
		}
	}
	return nil
}

func pageTemplateName(path string) string {
	if filepath.Ext(path) == ".md" {
		return "markdown-page"
	}
	return "page"
}

func (b *builder) prepareTemplates(d *dir2, p paths.Paths) (*template.Template, error) {
	t, err := d.templates.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}
	if filepath.Ext(p.Src) == ".tmpl" {
		t, err = t.ParseFiles(filepath.Join(b.args.Src, p.Src))
		if err != nil {
			return nil, fmt.Errorf("parsing itself: %w", err)
		}
	}
	return t, nil
}

func (b *builder) executeTemplates(p paths.Paths, t *template.Template, c *kask.TemplateContent) error {
	if b.args.Verbose {
		fmt.Printf("printing %s\n", p.Dst)
	}
	buf := bytes.NewBuffer([]byte{})
	if _, err := fmt.Fprintln(buf, fileheader); err != nil {
		return fmt.Errorf("writing the autogen notice: %w", err)
	}
	if err := t.ExecuteTemplate(buf, pageTemplateName(p.Src), c); err != nil {
		return fmt.Errorf("executing: %w", err)
	}
	bs := buf.Bytes()
	if filepath.Ext(p.Src) == ".tmpl" {
		var err error
		bs, err = rewriteLinksInHtmlPage(b.rw, p, bs)
		if err != nil {
			return fmt.Errorf("rewriting the links found at the page: %w", err)
		}
	}
	err := os.WriteFile(filepath.Join(b.args.Dst, p.Dst), bs, 0o666)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	return nil
}

func (b *builder) execPage(d *dir2, p paths.Paths) error {
	c := &kask.TemplateContent{
		Stylesheets: d.stylesheets,
		Node:        b.leaves[p.Url],
		Root:        b.root3,
		Markdown:    b.markdown[p.Src], // otherwise `nil`
		Time:        b.start,
	}
	t, err := b.prepareTemplates(d, p)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	err = b.executeTemplates(p, t, c)
	if err != nil {
		return fmt.Errorf("template: %w", err)
	}
	return nil
}

func (b *builder) execDir(d *dir2) error {
	err := os.MkdirAll(filepath.Join(b.args.Dst, d.paths.Dst), 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	for _, page := range d.original.Pages {
		p := d.paths.File(page, d.original.IsToStrip(), b.args.Provider.UrlMode())
		if err := b.execPage(d, p); err != nil {
			return fmt.Errorf("%q: %w", page, err)
		}
	}
	for _, subdir := range d.subdirs {
		if err := b.execDir(subdir); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.Src), err)
		}
	}
	return nil
}
