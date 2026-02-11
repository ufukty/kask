package builder

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"go.ufukty.com/kask/pkg/kask"
)

func pageTemplateName(path string) string {
	if filepath.Ext(path) == ".md" {
		return "markdown-page"
	}
	return "page"
}

func (b *builder) prepareTemplates(d *dir2, p paths) (*template.Template, error) {
	t, err := d.templates.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}
	if filepath.Ext(p.src) == ".tmpl" {
		t, err = t.ParseFiles(filepath.Join(b.args.Src, p.src))
		if err != nil {
			return nil, fmt.Errorf("parsing itself: %w", err)
		}
	}
	return t, nil
}

func (b *builder) executeTemplates(p paths, t *template.Template, c *kask.TemplateContent) error {
	if b.args.Verbose {
		fmt.Printf("printing %s\n", p.dst)
	}
	buf := bytes.NewBuffer([]byte{})
	if _, err := fmt.Fprintln(buf, fileheader); err != nil {
		return fmt.Errorf("writing the autogen notice: %w", err)
	}
	if err := t.ExecuteTemplate(buf, pageTemplateName(p.src), c); err != nil {
		return fmt.Errorf("executing: %w", err)
	}
	bs := buf.Bytes()
	if filepath.Ext(p.src) == ".tmpl" {
		var err error
		bs, err = rewriteLinksInHtmlPage(b.rw, p.dst, bs)
		if err != nil {
			return fmt.Errorf("rewriting the links found at the page: %w", err)
		}
	}
	err := os.WriteFile(filepath.Join(b.args.Dst, p.dst), bs, 0o666)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	return nil
}

// TODO: also use inside [builder.toNode]
func leafpath(src string) string {
	if base := filepath.Base(src); base == "index.tmpl" || base == "README.md" {
		return ""
	}
	return src
}

func (b *builder) execPage(d *dir2, p paths) error {
	c := &kask.TemplateContent{
		Stylesheets: d.stylesheets,
		Node:        b.leaves[pageref{d, leafpath(p.src)}],
		Root:        b.root3,
		Markdown:    b.markdown[p.src], // otherwise `nil`
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
	err := os.MkdirAll(filepath.Join(b.args.Dst, d.paths.dst), 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	for _, page := range d.original.Pages {
		p := d.paths.file(page, d.original.IsToStrip())
		if err := b.execPage(d, p); err != nil {
			return fmt.Errorf("%q: %w", page, err)
		}
	}
	for _, subdir := range d.subdirs {
		if err := b.execDir(subdir); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.src), err)
		}
	}
	return nil
}
