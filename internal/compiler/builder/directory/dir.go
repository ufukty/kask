package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type propfolder struct {
	Css  []string // *.css
	Tmpl []string // *.tmpl
	Page string   // page.html
}

type kaskfolder struct {
	Propagate *propfolder // propagate/
	Css       []string    // *.css
	Tmpl      []string    // *.tmpl
	Page      string      // page.html
	Meta      string      // meta.yml
}

type dirchecks struct {
	HasIndexFile     bool // either html or md
	HasIndexHtmlFile bool // index.html
	HasIndexMdFile   bool // index.md
	HasMetaFile      bool // .kask/meta.yml
	HasAssetDir      bool // assets
	HasPageHtml      bool // .kask/page.html, .kask/propagate/page.html
}

type PageType string

const (
	Html = PageType("html")
	Md   = PageType("md")
)

// information for templates
type Node struct {
	Parent      *Node
	Subpages    []*Node
	Subsections []*Node

	Title            string
	PageType         PageType
	Visitable        bool
	TargetInSitePath string
	SrcFilename      string
}

type pagesources[T any] struct {
	Html []T
	Md   []T
}

type Dir struct {
	SiteRoot   string
	InSitePath string
	Dirname    string
	PageHtml   string // inherited or declared in .kask or .kask/propagate

	Node     *Node
	Checks   *dirchecks
	Kask     *kaskfolder
	Children []*Dir

	meta *meta
}

func inspectPropagateDir(d *Dir) error {
	entries, err := os.ReadDir(filepath.Join(d.Path(), ".kask/propagate"))
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}
	d.Kask.Propagate = &propfolder{}
	for _, entry := range entries {
		var name, isDir = entry.Name(), entry.IsDir()

		switch {
		case !isDir && name == "page.html":
			d.Checks.HasPageHtml = true
			d.Kask.Propagate.Page = ".kask/propagate/page.html"

		case !isDir && strings.HasSuffix(name, ".css"):
			d.Kask.Propagate.Css = append(d.Kask.Propagate.Css, filepath.Join(".kask/propagate", name))

		case !isDir && strings.HasSuffix(name, ".tmpl"):
			d.Kask.Propagate.Tmpl = append(d.Kask.Propagate.Tmpl, filepath.Join(".kask/propagate", name))

		}
	}
	return nil
}

func inspectKaskFolder(d *Dir) error {
	entries, err := os.ReadDir(filepath.Join(d.Path(), ".kask"))
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}
	d.Kask = &kaskfolder{}
	for _, entry := range entries {
		name, isDir := entry.Name(), entry.IsDir()

		switch {
		case !isDir && name == "page.html":
			d.Checks.HasPageHtml = true
			d.Kask.Page = ".kask/page.html"

		case !isDir && name == "meta.yml":
			d.Checks.HasMetaFile = true
			d.Kask.Meta = ".kask/meta.yml"
			d.meta, err = readMeta(filepath.Join(d.SiteRoot, d.InSitePath, d.Kask.Meta))
			if err != nil {
				return fmt.Errorf("readMeta: %w", err)
			}

		case !isDir && strings.HasSuffix(name, ".css"):
			d.Kask.Css = append(d.Kask.Css, filepath.Join(".kask", name))

		case !isDir && strings.HasSuffix(name, ".tmpl"):
			d.Kask.Tmpl = append(d.Kask.Tmpl, filepath.Join(".kask", name))

		case isDir && name == "propagate":
			err := inspectPropagateDir(d)
			if err != nil {
				return fmt.Errorf("inspectPropagateDir: %w", err)
			}

		}
	}
	return nil
}

func inspect(siteRoot, inSitePath string, inheritedpagehtml string) (*Dir, error) {
	d := &Dir{
		SiteRoot:   siteRoot,
		InSitePath: inSitePath,
		Dirname:    filepath.Base(inSitePath),
		Children:   []*Dir{},
		Checks:     &dirchecks{},
		PageHtml:   inheritedpagehtml,
	}

	entries, err := os.ReadDir(d.Path())
	if err != nil {
		return nil, fmt.Errorf("listing directory entries: %w", err)
	}

	pages := pagesources[string]{}
	subdirs := []string{}

	for _, entry := range entries {
		var name, isDir = entry.Name(), entry.IsDir()

		switch {
		case !isDir && name == "index.html":
			d.Checks.HasIndexHtmlFile = true
			fallthrough
		case !isDir && strings.HasSuffix(name, ".html"):
			pages.Html = append(pages.Html, name)

		case !isDir && name == "index.md":
			d.Checks.HasIndexMdFile = true
			fallthrough
		case !isDir && strings.HasSuffix(name, ".md"):
			pages.Md = append(pages.Md, name)

		case !isDir:
			fmt.Println("ignored file:", filepath.Join(d.InSitePath, name))

		case isDir && name == ".kask":
			err := inspectKaskFolder(d)
			if err != nil {
				return nil, fmt.Errorf("inspectKaskFolder: %w", err)
			}

		case isDir && name == ".assets":
			d.Checks.HasAssetDir = true

		case isDir:
			subdirs = append(subdirs, name)
		}
	}

	d.Checks.HasIndexFile = d.Checks.HasIndexHtmlFile || d.Checks.HasIndexMdFile

	if d.Checks.HasIndexHtmlFile && d.Checks.HasIndexMdFile {
		return nil, fmt.Errorf("directory contains both index.html file and markdown files: %s", d.InSitePath)
	}
	if !d.Checks.HasIndexFile && !d.Checks.HasMetaFile {
		return nil, fmt.Errorf("no index page or meta file found: %s", d.InSitePath)
	}

	if d.Kask != nil && d.Kask.Page != "" {
		d.PageHtml = filepath.Join(d.InSitePath, d.Kask.Page)
	} else if d.Kask != nil && d.Kask.Propagate != nil && d.Kask.Propagate.Page != "" {
		d.PageHtml = filepath.Join(d.InSitePath, d.Kask.Propagate.Page)
	}

	mdPageNeeded := len(pages.Md) > 0
	mdPageAvailable := d.Checks.HasPageHtml || d.PageHtml != ""
	if mdPageNeeded && !mdPageAvailable {
		return nil, fmt.Errorf("no page template found or propagated for markdown files: %s", d.InSitePath)
	}

	title := d.Dirname
	if d.meta != nil && d.meta.Title != "" {
		title = d.meta.Title
	}

	if d.Checks.HasIndexMdFile {
		d.Node = &Node{
			Parent:      nil,
			Subpages:    []*Node{},
			Subsections: []*Node{},

			Title:            title,
			PageType:         Md,
			Visitable:        true,
			TargetInSitePath: filepath.Join(d.InSitePath, "index.html"),
			SrcFilename:      "index.md",
		}
	} else if d.Checks.HasIndexHtmlFile {
		d.Node = &Node{
			Parent:      nil,
			Subpages:    []*Node{},
			Subsections: []*Node{},

			Title:            title,
			PageType:         Html,
			Visitable:        true,
			TargetInSitePath: filepath.Join(d.InSitePath, "index.html"),
			SrcFilename:      "index.html",
		}
	} else {
		d.Node = &Node{
			Parent:      nil,
			Subpages:    []*Node{},
			Subsections: []*Node{},

			Title:            title,
			Visitable:        false,
			TargetInSitePath: d.InSitePath,
		}
	}

	for _, md := range pages.Md {
		if md == "index.md" {
			continue
		}
		page := &Node{
			Parent:      d.Node,
			Subpages:    []*Node{},
			Subsections: []*Node{},

			Title:            strings.TrimSuffix(md, ".md"),
			PageType:         Md,
			Visitable:        true,
			TargetInSitePath: filepath.Join(d.InSitePath, strings.TrimSuffix(md, ".md")+".html"),
			SrcFilename:      md,
		}
		// d.Pagesources.Md = append(d.Pagesources.Md, page)
		d.Node.Subpages = append(d.Node.Subpages, page)
	}
	for _, html := range pages.Html {
		if html == "index.html" {
			continue
		}
		page := &Node{
			Parent:      d.Node,
			Subpages:    []*Node{},
			Subsections: []*Node{},

			Title:            strings.TrimSuffix(html, ".html"),
			PageType:         Html,
			Visitable:        true,
			TargetInSitePath: filepath.Join(d.InSitePath, html),
			SrcFilename:      html,
		}
		// d.Pagesources.Html = append(d.Pagesources.Html, page)
		d.Node.Subpages = append(d.Node.Subpages, page)
	}

	if len(subdirs) > 0 {
		for _, name := range subdirs {
			subdir, err := inspect(d.SiteRoot, filepath.Join(d.InSitePath, name), d.PageHtml)
			if err != nil {
				return nil, fmt.Errorf("inspect: %w", err)
			}
			d.Children = append(d.Children, subdir)
			d.Node.Subsections = append(d.Node.Subsections, subdir.Node)
			subdir.Node.Parent = d.Node
		}
	}
	return d, nil
}

func Inspect(path string) (*Dir, error) {
	root, err := inspect(path, ".", "")
	if err != nil {
		return nil, fmt.Errorf("newDir: %w", err)
	}
	return root, nil
}
