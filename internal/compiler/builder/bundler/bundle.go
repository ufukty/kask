package bundler

import (
	"fmt"
	"os"
	"path/filepath"
)

func BundleCss(dstdir, filebase string, files []string) (string, error) {
	id, err := newid()
	if err != nil {
		return "", fmt.Errorf("newid: %w", err)
	}
	filename := fmt.Sprintf("%s-%s.css", filebase, id)
	dst := filepath.Join(dstdir, filename)
	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("creating destination file to write into: %w", err)
	}
	defer f.Close()
	for _, path := range files {
		appendfile(f, path)
	}
	return filename, nil
}
