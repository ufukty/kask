package copy

import (
	"fmt"
	"io"
	"io/fs"

	"go.ufukty.com/kask/internal/disk"
)

func File(dstFs disk.WriteFS, dst string, srcFs disk.ReadFS, src string) error {
	srcfile, err := srcFs.Open(src)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer srcfile.Close()
	dstfile, err := dstFs.Create(dst)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer dstfile.Close()
	_, err = io.Copy(dstfile, srcfile)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return nil
}

func Dir(dstFs disk.WriteFS, dst string, srcFs disk.ReadFS, src string) error {
	fi, err := srcFs.Stat(src)
	if err != nil {
		return fmt.Errorf("stat: %w", err)
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %s", src)
	}
	err = fs.WalkDir(srcFs, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("internal: %w", err)
		}
		if d.IsDir() {
			if err := dstFs.MkdirAll(path); err != nil {
				return fmt.Errorf("mkdir: %w", err)
			}
		} else {
			if err := File(dstFs, dst, srcFs, path); err != nil {
				return fmt.Errorf("file: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk: %w", err)
	}
	return nil
}
