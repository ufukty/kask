package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"go.ufukty.com/kask/internal/builder/copy"
	"go.ufukty.com/kask/internal/paths"
)

func (b *builder) copyAssetDir(path paths.Paths) error {
	err := os.MkdirAll(path.Dst, 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	es, err := os.ReadDir(filepath.Join(b.args.Src, path.Src))
	if err != nil {
		return fmt.Errorf("listing: %w", err)
	}
	for _, e := range es {
		if e.IsDir() {
			err := b.copyAssetDir(path.Subdir(e.Name(), false))
			if err != nil {
				return fmt.Errorf("%q: %w", e.Name(), err)
			}
		} else {
			file := path.File(e.Name(), false, paths.UrlModeDefault)
			if b.args.Verbose {
				fmt.Printf("copying %q => %q", file.Src, file.Dst)
			}
			b.rw.Bank(file.Src, file.Url)
			err := copy.File(filepath.Join(b.args.Dst, file.Dst), filepath.Join(b.args.Src, file.Src))
			if err != nil {
				return fmt.Errorf("copying file %q: %w", e.Name(), err)
			}
		}
	}
	return nil
}

func (b *builder) assets(d *dir2) error {
	if d.original.Assets {
		err := b.copyAssetDir(d.paths.Subdir(".assets", false))
		if err != nil {
			return fmt.Errorf("traverse: %w", err)
		}
	}
	for _, c := range d.subdirs {
		err := b.assets(c)
		if err != nil {
			return fmt.Errorf("%q: %w", c.original.Name, err)
		}
	}
	return nil
}
