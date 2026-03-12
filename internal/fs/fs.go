package fs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Supplementary to [fs.FS]
// Needed for testability
type Writable interface {
	MkdirAll(path string) error
	WriteFile(path string, data []byte) error
}

type Real struct {
	root string
}

var (
	_ Writable     = (*Real)(nil)
	_ fs.FS        = (*Real)(nil)
	_ fs.ReadDirFS = (*Real)(nil)
)

func NewReal(root string) Real {
	return Real{root: root}
}

func (r Real) MkdirAll(path string) error {
	return os.MkdirAll(filepath.Join(r.root, path), 0o755)
}

func (r Real) WriteFile(name string, data []byte) error {
	return os.WriteFile(filepath.Join(r.root, name), data, 0o600)
}

func (r Real) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(r.root, name))
}

func (r Real) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(filepath.Join(r.root, name))
}
