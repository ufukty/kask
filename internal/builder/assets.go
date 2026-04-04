package builder

import (
	"fmt"

	"go.ufukty.com/kask/internal/disk/copy"
	"go.ufukty.com/kask/internal/paths"
)

func (b *builder) copyAssetDir(path paths.Paths) error {
	err := b.args.Dst.MkdirAll(path.Dst, 0755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	es, err := b.args.Src.ReadDir(path.Src)
	if err != nil {
		return fmt.Errorf("listing: %w", err)
	}
	for _, e := range es {
		if e.IsDir() {
			err := b.copyAssetDir(path.AssetDir(e.Name()))
			if err != nil {
				return fmt.Errorf("%q: %w", e.Name(), err)
			}
		} else {
			file := path.AssetFile(e.Name())
			if b.args.Verbose {
				fmt.Printf("copying asset %-30q => %-30q\n", file.Src, file.Dst)
			}
			b.rw.Bank(file.Src, file.Url)
			err := copy.File(b.args.Dst, file.Dst, b.args.Src, file.Src)
			if err != nil {
				return fmt.Errorf("copying file %q: %w", e.Name(), err)
			}
		}
	}
	return nil
}

func (b *builder) assets(d *dir2) error {
	if d.original.Assets {
		err := b.copyAssetDir(d.paths.AssetDir(".assets"))
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
