package builder

import (
	"fmt"
	"path/filepath"
	"slices"

	"go.ufukty.com/kask/internal/builder/bundle"
)

func (b *builder) write(dst, content string) error {
	if b.args.Verbose {
		fmt.Println("writing", dst)
	}
	if err := b.args.Dst.MkdirAll(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	if err := b.args.Dst.WriteFile(dst, []byte(content)); err != nil {
		return fmt.Errorf("writing: %w", err)
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
		p := d.paths.Stylesheet(true)
		b.rw.Bank(p.Src, p.Url)
		if err := b.write(p.Dst, css); err != nil {
			return fmt.Errorf("writing propagated css file: %w", err)
		}
		d.stylesheets = append(d.stylesheets, p.Url)
		toPropagate = append(toPropagate, p.Url)
	}

	if d.original.Kask != nil && len(d.original.Kask.Css) > 0 {
		css, err := bundle.Files(d.original.Kask.Css)
		if err != nil {
			return fmt.Errorf("bundling at-level css file: %w", err)
		}
		p := d.paths.Stylesheet(false)
		b.rw.Bank(p.Src, p.Url)
		if err := b.write(p.Dst, css); err != nil {
			return fmt.Errorf("writing at-level css file: %w", err)
		}
		d.stylesheets = append(d.stylesheets, p.Url)
	}

	for _, subdir := range d.subdirs {
		if err := b.bundleAndPropagateStylesheets(subdir, slices.Clone(toPropagate)); err != nil {
			return fmt.Errorf("%q: %w", filepath.Base(subdir.paths.Src), err)
		}
	}

	return nil
}
