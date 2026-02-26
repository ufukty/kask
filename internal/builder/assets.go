package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"go.ufukty.com/kask/internal/builder/copy"
)

func (b *builder) copyAssetsFolders(d *dir2) error {
	if d.original.Assets {
		err := os.MkdirAll(filepath.Join(b.args.Dst, d.paths.Dst), 0o755)
		if err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		dst := filepath.Join(b.args.Dst, d.paths.Dst, ".assets")
		src := filepath.Join(b.args.Src, d.paths.Src, ".assets")
		if b.args.Verbose {
			fmt.Println("copying", dst)
		}
		err = copy.Dir(dst, src)
		if err != nil {
			return fmt.Errorf("copy dir: %w", err)
		}
	}
	for _, subdir := range d.subdirs {
		if err := b.copyAssetsFolders(subdir); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.Src), err)
		}
	}
	return nil
}
