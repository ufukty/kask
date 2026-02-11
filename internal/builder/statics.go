package builder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"go.ufukty.com/kask/internal/builder/bundle"
	"go.ufukty.com/kask/internal/builder/copy"
)

func (b *builder) write(dst, content string) error {
	if b.args.Verbose {
		fmt.Println("writing", dst)
	}
	err := os.MkdirAll(filepath.Dir(dst), 0o755)
	if err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(f, strings.NewReader(content))
	if err != nil {
		return fmt.Errorf("copying: %w", err)
	}
	return nil
}

func (b *builder) bundleAndPropagateStylesheets(d *dir2, toPropagate []string) error {
	d.stylesheets = slices.Clone(toPropagate)

	if d.original.Kask != nil && d.original.Kask.Propagate != nil && len(d.original.Kask.Propagate.Css) > 0 {
		css, err := bundle.Files(d.original.Kask.Propagate.Css)
		if err != nil {
			return fmt.Errorf("bundling propagated css file: %w", err)
		}
		dst := "/" + filepath.Join(d.paths.dst, "styles.propagate.css")
		if err := b.write(filepath.Join(b.args.Dst, dst), css); err != nil {
			return fmt.Errorf("writing propagated css file: %w", err)
		}
		d.stylesheets = append(d.stylesheets, dst)
		toPropagate = append(toPropagate, dst)
	}

	if d.original.Kask != nil && len(d.original.Kask.Css) > 0 {
		css, err := bundle.Files(d.original.Kask.Css)
		if err != nil {
			return fmt.Errorf("bundling at-level css file: %w", err)
		}
		dst := "/" + filepath.Join(d.paths.dst, "styles.css")
		if err := b.write(filepath.Join(b.args.Dst, dst), css); err != nil {
			return fmt.Errorf("writing at-level css file: %w", err)
		}
		d.stylesheets = append(d.stylesheets, dst)
	}

	for _, subdir := range d.subdirs {
		if err := b.bundleAndPropagateStylesheets(subdir, slices.Clone(toPropagate)); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.src), err)
		}
	}

	return nil
}

func (b *builder) copyAssetsFolders(d *dir2) error {
	if d.original.Assets {
		err := os.MkdirAll(filepath.Join(b.args.Dst, d.paths.dst), 0o755)
		if err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}

		dst := filepath.Join(b.args.Dst, d.paths.dst, ".assets")
		src := filepath.Join(b.args.Src, d.paths.src, ".assets")
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
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.src), err)
		}
	}

	return nil
}
