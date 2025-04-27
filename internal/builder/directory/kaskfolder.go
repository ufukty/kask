package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ready-to-use propagate folder content
type propagate struct {
	Css  []string
	Tmpl []string
	Page string
}

func inspectPropagateDir(dir string) (*propagate, error) {
	entries, err := os.ReadDir(filepath.Join(dir, ".kask/propagate"))
	if err != nil {
		return nil, fmt.Errorf("listing directory: %w", err)
	}

	prop := &propagate{}
	for _, entry := range entries {
		name, isDir := entry.Name(), entry.IsDir()

		switch {
		case !isDir && strings.HasSuffix(name, ".css"):
			prop.Css = append(prop.Css, filepath.Join(dir, ".kask/propagate", name))

		case !isDir && strings.HasSuffix(name, ".tmpl"):
			prop.Tmpl = append(prop.Tmpl, filepath.Join(dir, ".kask/propagate", name))

		}
	}

	return prop, nil
}

// ready-to-use Kask folder content (contains propagate content at the folder)
type Kask struct {
	Propagate *propagate

	Css  []string
	Tmpl []string
	Page string
	Meta *Meta
}

func inspectKaskFolder(dir string) (*Kask, error) {
	entries, err := os.ReadDir(filepath.Join(dir, ".kask"))
	if err != nil {
		return nil, fmt.Errorf("listing directory: %w", err)
	}

	kask := &Kask{}
	for _, entry := range entries {
		name, isDir := entry.Name(), entry.IsDir()

		switch {
		case !isDir && strings.HasSuffix(name, ".css"):
			kask.Css = append(kask.Css, filepath.Join(dir, ".kask", name))

		case !isDir && strings.HasSuffix(name, ".tmpl"):
			kask.Tmpl = append(kask.Tmpl, filepath.Join(dir, ".kask", name))

		case !isDir && name == "meta.yml":
			kask.Meta, err = readMeta(filepath.Join(dir, ".kask/meta.yml"))
			if err != nil {
				return nil, fmt.Errorf("reading meta file %s: %w", filepath.Join(dir, ".kask/meta.yml"), err)
			}

		case isDir && name == "propagate":
			kask.Propagate, err = inspectPropagateDir(dir)
			if err != nil {
				return nil, fmt.Errorf("propagate folder: %w", err)
			}
		}
	}

	return kask, nil
}
