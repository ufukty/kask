package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Dir struct {
	Name          string
	Assets        string
	Subdirs       []*Dir
	PagesMarkdown []string
	PagesHtml     []string
	Kask          *Kask
}

func (d *Dir) subtree() int {
	c := len(d.PagesHtml) + len(d.PagesMarkdown)
	for _, s := range d.Subdirs {
		c += s.subtree()
	}
	return c
}

func inspect(root, path string) (*Dir, error) {
	d := &Dir{
		Name:    filepath.Base(path),
		Subdirs: []*Dir{},
	}

	entries, err := os.ReadDir(filepath.Join(root, path))
	if err != nil {
		return nil, fmt.Errorf("listing directory entries: %w", err)
	}

	subdirs := []string{}

	for _, entry := range entries {
		var name, isDir = entry.Name(), entry.IsDir()

		switch {
		case !isDir && strings.HasSuffix(name, ".html"):
			d.PagesHtml = append(d.PagesHtml, filepath.Join(path, name))

		case !isDir && strings.HasSuffix(name, ".md"):
			d.PagesMarkdown = append(d.PagesMarkdown, filepath.Join(path, name))

		case isDir && name == ".kask":
			d.Kask, err = inspectKaskFolder(filepath.Join(root, path))
			if err != nil {
				return nil, fmt.Errorf("kask folder: %w", err)
			}

		case isDir && name == ".assets":
			d.Assets = filepath.Join(path, ".assets")

		case isDir:
			subdirs = append(subdirs, filepath.Join(path, name))
		}
	}

	for _, subdir := range subdirs {
		sub, err := inspect(root, subdir)
		if err != nil {
			return nil, fmt.Errorf("inspecting %s: %w", path, err)
		}
		if sub.subtree() > 0 {
			d.Subdirs = append(d.Subdirs, sub)
		}
	}

	return d, nil
}

func Inspect(path string) (*Dir, error) {
	root, err := inspect(path, ".")
	if err != nil {
		return nil, fmt.Errorf("inspect: %w", err)
	}
	return root, nil
}
