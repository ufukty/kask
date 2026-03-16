package directory

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.ufukty.com/kask/internal/disk"
)

type Dir struct {
	Name    string
	Assets  bool
	Subdirs []*Dir
	Pages   []string // .md + .tmpl
	Kask    *Kask
	Meta    *Meta
}

func (d *Dir) subtree() int {
	c := len(d.Pages)
	for _, s := range d.Subdirs {
		c += s.subtree()
	}
	return c
}

func inspect(fs disk.ReadFS, path string) (*Dir, error) {
	d := &Dir{
		Name:    filepath.Base(path),
		Subdirs: []*Dir{},
	}
	entries, err := fs.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("listing directory entries: %w", err)
	}
	subdirs := []string{}
	for _, entry := range entries {
		name, isDir := entry.Name(), entry.IsDir()
		switch {
		case !isDir && (strings.HasSuffix(name, ".tmpl") || strings.HasSuffix(name, ".md")):
			d.Pages = append(d.Pages, name)
		case !isDir && name == ".kask.yml":
			d.Meta, err = readMeta(fs, filepath.Join(path, name))
			if err != nil {
				return nil, fmt.Errorf("reading meta file %q: %w", filepath.Join(path, name), err)
			}
		case isDir && name == ".kask":
			d.Kask, err = inspectKaskFolder(fs, path)
			if err != nil {
				return nil, fmt.Errorf("inspecting kask folder: %w", err)
			}
		case isDir && name == ".assets":
			d.Assets = true
		case isDir:
			subdirs = append(subdirs, name)
		}
	}
	for _, subdir := range subdirs {
		sub, err := inspect(fs, filepath.Join(path, subdir))
		if err != nil {
			return nil, fmt.Errorf("%q: %w", path, err)
		}
		if sub.subtree() > 0 {
			d.Subdirs = append(d.Subdirs, sub)
		}
	}
	return d, nil
}

func Inspect(fs disk.ReadFS) (*Dir, error) {
	return inspect(fs, ".")
}

func (d *Dir) IsToStrip() bool {
	return d == nil || d.Meta == nil || !d.Meta.PreserveOrdering
}
