package copy

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func File(dst, src string) error {
	srcfile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfile.Close()

	dstfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstfile.Close()

	_, err = io.Copy(dstfile, srcfile)
	if err != nil {
		return err
	}
	s, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, s.Mode())
}

func Dir(dst, src string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	err = os.MkdirAll(dst, s.Mode())
	if err != nil {
		return err
	}
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return File(targetPath, path)
	})

	return err
}
