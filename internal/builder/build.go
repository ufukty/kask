package builder

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/ufukty/kask/cmd/kask/commands/version"
	"github.com/ufukty/kask/internal/builder/bundle"
	"github.com/ufukty/kask/internal/builder/directory"
	"github.com/ufukty/kask/internal/builder/markdown"
)

type Args struct {
	Domain string
	Dev    bool // suffixes css bundles with unique ids to bypass browser caching
}

type builder struct {
	args Args

	assets        []string                  // src paths
	stylesheets   map[string]string         // src path -> content
	pagesMarkdown map[string]*markdown.Page // src path -> content
}

func has[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func (b *builder) checkCompetingEntries(dir *directory.Dir) error {
	children := map[string]int{}
	for _, subdir := range dir.Subdirs {
		children[subdir.Name] = 1
	}
	for _, page := range dir.Pages.Html {
		if has(children, page) {
			children[page] = -1
		}
		children[page]++
	}
	for _, page := range dir.Pages.Markdown {
		if has(children, page) {
			children[page] = -1
		}
		children[page]++
	}
	duplicates := []string{}
	for child, freq := range children {
		if freq > 1 {
			duplicates = append(duplicates, child)
		}
	}
	if len(duplicates) > 0 {
		return fmt.Errorf("multiple entries sharing the same path for those: %s", strings.Join(duplicates, ", "))
	}
	for _, sub := range dir.Subdirs {
		if err := b.checkCompetingEntries(sub); err != nil {
			return fmt.Errorf("%q: %w", sub.Name, err)
		}
	}
	return nil
}

// used in assigning destination addresses, bundling css, and propagating tmpl files
type dir2 struct {
	Kask *directory.Kask

	SrcName, SrcPath, SrcAssets string
	DstName, DstPath, DstAssets string // path encoded

	Subdirs []*dir2

	PagesMarkdown []string // src paths
	PagesHtml     []string // src paths
	Stylesheets   []string // dst paths

	Tmpl *template.Template
}

func (b *builder) toDir2(d *directory.Dir, srcparent, dstparent string) *dir2 {
	srcparent = filepath.Join(srcparent, d.Name)
	dstparent = filepath.Join(dstparent, url.PathEscape(d.Name)) // escaped
	d2 := &dir2{
		Kask: d.Kask,

		Subdirs: []*dir2{},

		SrcName:   d.Name,
		SrcPath:   srcparent,
		SrcAssets: d.Assets,

		DstName:   url.PathEscape(d.Name),
		DstPath:   dstparent,
		DstAssets: "",

		PagesMarkdown: d.Pages.Markdown,
		PagesHtml:     d.Pages.Html,
		Stylesheets:   []string{},

		Tmpl: nil,
	}
	for _, subdir := range d.Subdirs {
		d2.Subdirs = append(d2.Subdirs, b.toDir2(subdir, srcparent, dstparent))
	}
	return d2
}

func (b *builder) bundleStylesheets(d *dir2) error {
	if d.Kask != nil && d.Kask.Propagate != nil && len(d.Kask.Propagate.Css) > 0 {
		css, err := bundle.Files(d.Kask.Propagate.Css)
		if err != nil {
			return fmt.Errorf("bundling .kask/propagate/*.css: %w", err)
		}
		dst := filepath.Join(d.SrcPath, "styles.propagate.css")
		b.stylesheets[dst] = css
		d.Stylesheets = append(d.Stylesheets, dst)
	}

	if d.Kask != nil && len(d.Kask.Css) > 0 {
		css, err := bundle.Files(d.Kask.Css)
		if err != nil {
			return fmt.Errorf("bundling .kask/*.css: %w", err)
		}
		dst := filepath.Join(d.SrcPath, "styles.css")
		b.stylesheets[dst] = css
		d.Stylesheets = append(d.Stylesheets, dst)
	}

	for _, subdir := range d.Subdirs {
		if err := b.bundleStylesheets(subdir); err != nil {
			return fmt.Errorf("%s: %w", subdir.SrcName, err)
		}
	}

	return nil
}

func (b *builder) propagateTemplates(d *dir2, toPropagate *template.Template) error {
	var err error

	if d.Kask != nil && d.Kask.Propagate != nil && len(d.Kask.Propagate.Tmpl) > 0 {
		toPropagate, err = d.Tmpl.ParseFiles(d.Kask.Propagate.Tmpl...)
		if err != nil {
			return fmt.Errorf("parsing to-propagate template files: %w", err)
		}
	}

	atLevel, err := toPropagate.Clone()
	if err != nil {
		return fmt.Errorf("cloning propagated: %w", err)
	}

	if d.Kask != nil && len(d.Kask.Tmpl) > 0 {
		atLevel, err = atLevel.ParseFiles(d.Kask.Tmpl...)
		if err != nil {
			return fmt.Errorf("parsing at-level template files: %w", err)
		}
	}

	d.Tmpl = atLevel

	for _, subdir := range d.Subdirs {
		if err := b.propagateTemplates(subdir, toPropagate); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", subdir.SrcName, err)
		}
	}

	return nil
}

func (b *builder) renderMarkdown(d *dir2) error {
	for _, md := range d.PagesMarkdown {
		page, err := markdown.ToHtml(md)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", md, err)
		}
		b.pagesMarkdown[md] = page
	}

	for _, subdir := range d.Subdirs {
		if err := b.renderMarkdown(subdir); err != nil {
			return fmt.Errorf("%s: %w", subdir.SrcName, err)
		}
	}

	return nil
}

type Node struct {
	Children   []*Node
	Name, Path string
	Visitable  bool
}

func isVisitable(d *dir2) bool {
	return slices.ContainsFunc(d.PagesHtml, func(path string) bool {
		return filepath.Base(path) == "index.html"
	}) || slices.ContainsFunc(d.PagesMarkdown, func(path string) bool {
		return filepath.Base(path) == "README.md"
	})
}

var fileheader = fmt.Sprintf("<!-- Do not edit. Auto-generated by Kask %s -->", version.Version)

func (b *builder) absolutify(prefix string, files []string) []string {
	files2 := []string{}
	for _, file := range files {
		files2 = append(files2, filepath.Join(prefix, file))
	}
	return files2
}

// template files should access necessary information through
// the fields of this struct
type TemplateContent struct {
	Stylesheets     []string
	Node            *rendered.Node
	WebSiteRoot     *rendered.Node
	MarkdownContent string
	MarkdownTOC     *markdown.TocNode
	Time            time.Time
	Dir             *directory.Dir
}

func (b *builder) execTemplates(dstroot string, uri string, root, d *rendered.Node, inherited *context, s *Args) error {
	err := os.MkdirAll(filepath.Join(dstroot, d.InSitePath), 0755)
	if err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	tmplPropagate, err := inherited.Template.Clone()
	if err != nil {
		return fmt.Errorf("inherited.Template.Clone: %w", err)
	}
	propagate := &context{
		Template:    tmplPropagate,
		Stylesheets: inherited.Stylesheets,
	}
	if d.Kask != nil && d.Kask.Propagate != nil {
		if len(d.Kask.Propagate.Css) > 0 {
			propcss, err := bundler.BundleCss(filepath.Join(dstroot, uri), "propagate", absolutify(filepath.Join(d.SiteRoot, d.InSitePath), d.Kask.Propagate.Css))
			if err != nil {
				return fmt.Errorf("bundler.BundleCss(.kask/propagate): %w", err)
			}
			url, err := url.JoinPath(s.Domain, d.InSitePath, propcss)
			if err != nil {
				return fmt.Errorf("url.JoinPath(.kask/propagate): %w", err)
			}
			propagate.Stylesheets = append(inherited.Stylesheets, url)
		}
		if len(d.Kask.Propagate.Tmpl) > 0 {
			_, err = propagate.Template.ParseFiles(absolutify(filepath.Join(d.SiteRoot, d.InSitePath), d.Kask.Propagate.Tmpl)...)
			if err != nil {
				return fmt.Errorf("propagate.Template.ParseFiles: %w", err)
			}
		}
	}

	tmplLevel, err := tmplPropagate.Clone()
	if err != nil {
		return fmt.Errorf("tmplPropagate.Clone: %w", err)
	}
	level := &context{
		Template:    tmplLevel,
		Stylesheets: propagate.Stylesheets,
	}
	if d.Kask != nil {
		if len(d.Kask.Css) > 0 {
			levelcss, err := bundler.BundleCss(filepath.Join(dstroot, uri), "styles", absolutify(filepath.Join(d.SiteRoot, d.InSitePath), d.Kask.Css))
			if err != nil {
				return fmt.Errorf("bundler.BundleCss(.kask): %w", err)
			}
			url, err := url.JoinPath(s.Domain, d.InSitePath, levelcss)
			if err != nil {
				return fmt.Errorf("url.JoinPath(.kask): %w", err)
			}
			level.Stylesheets = append(level.Stylesheets, url)
		}
		if len(d.Kask.Tmpl) > 0 {
			_, err = level.Template.ParseFiles(absolutify(filepath.Join(d.SiteRoot, d.InSitePath), d.Kask.Tmpl)...)
			if err != nil {
				return fmt.Errorf("level.Template.ParseFiles: %w", err)
			}
		}
	}

	dt := &TemplateContent{
		Stylesheets:     level.Stylesheets,
		WebSiteRoot:     root,
		MarkdownContent: "",
		MarkdownTOC:     nil,
		Time:            time.Now(),
		Dir:             d,
	}

	queue := []*directory.Node{}
	if d.Node.Visitable {
		queue = append(queue, d.Node)
	}
	queue = append(queue, d.Node.Subpages...)

	for _, page := range queue {
		dt.Node = page
		tmplPage, err := tmplLevel.Clone()
		if err != nil {
			return fmt.Errorf("tmpLevel.Clone: %w", err)
		}
		buf := bytes.NewBuffer([]byte{})
		fmt.Fprintln(buf, fileheader)
		if page.PageType == directory.Html {
			tmplpath := filepath.Join(d.SiteRoot, d.InSitePath, page.SrcFilename)
			_, err = tmplPage.ParseFiles(tmplpath)
			if err != nil {
				return fmt.Errorf("tmplPage.ParseFiles: %w", err)
			}
			err = tmplPage.ExecuteTemplate(buf, "page", dt)
			if err != nil {
				return fmt.Errorf("tmplPage.ExecuteTemplate at %s: %w", tmplpath, err)
			}
		} else if page.PageType == directory.Md {
			_, err = tmplPage.ParseFiles(filepath.Join(d.SiteRoot, d.PageHtml)) // d.PageHtml contains InSitePath
			if err != nil {
				return fmt.Errorf("tmplPage.ParseFiles: %w", err)
			}
			html, toc, err := markdown.ToHtml(filepath.Join(d.SiteRoot, d.InSitePath, page.SrcFilename))
			if err != nil {
				return fmt.Errorf("markdown.ToHtml: %w", err)
			}
			dt.MarkdownContent = html
			dt.MarkdownTOC = toc
			if len(toc.Children) > 0 { // multiple h1's
				page.Title = toc.Children[0].Title
			}
			err = tmplPage.ExecuteTemplate(buf, "page", dt)
			if err != nil {
				return fmt.Errorf("tmplLevel.ExecuteTemplate at %s: %w", filepath.Join(d.InSitePath, page.SrcFilename), err)
			}
		}
		dst, err := os.Create(filepath.Join(dstroot, page.TargetInSitePath))
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		defer dst.Close()
		_, err = io.Copy(dst, buf)
		if err != nil {
			return fmt.Errorf("io.Copy: %w", err)
		}
		dst.Close() // duplicate. because defer won't work until function returns
	}

	if d.Checks.HasAssetDir {
		err := copy.Dir(filepath.Join(dstroot, d.InSitePath, ".assets"), filepath.Join(d.SiteRoot, d.InSitePath, ".assets"))
		if err != nil {
			return fmt.Errorf("copyDir: %w", err)
		}
	}

	for _, sub := range d.Children {
		uri, err := url.JoinPath(uri, url.PathEscape(sub.Dirname))
		if err != nil {
			return fmt.Errorf("url.JoinPath: %w", err)
		}
		err = build(dstroot, uri, root, sub, propagate, s)
		if err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	return nil
}

// [builder.Build] run separate DFS processes because there are many steps
// involving previous's complete results like templates can access to the
// sitemap which contains headers extracted from markdown files
func (b *builder) Build(dst, src string) error {
	root, err := directory.Inspect(src)
	if err != nil {
		return fmt.Errorf("inspecting source directory: %w", err)
	}

	if err := b.checkCompetingEntries(root); err != nil {
		return fmt.Errorf("checking competing files and folders: %w", err)
	}

	root2 := b.toDir2(root, "", "")

	if err := b.bundleStylesheets(root2); err != nil {
		return fmt.Errorf("bundling stylesheets: %w", err)
	}

	if err := b.propagateTemplates(root2, template.New("root")); err != nil {
		return fmt.Errorf("bundling stylesheets: %w", err)
	}

	if err := b.renderMarkdown(root2); err != nil {
		return fmt.Errorf("rendering markdown pages: %w", err)
	}

	err = builder.Build(dst, dir2, args)
	if err != nil {
		return fmt.Errorf("builder.Build: %w", err)
	}

	t := template.New("root")
	t.Funcs(template.FuncMap{
		"breadcrumbs": functions.Breadcrumbs,
		"dict":        functions.Dict,
	})

	err := b.execTemplates(dst, "/", root, root, i, s)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}
	return nil
}

func Build(dst, src string, args Args) error {
	b := &builder{
		stylesheets: map[string]string{},
		args:        args,
	}
	return b.Build(dst, src)
}
