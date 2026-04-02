package disk

import (
	"io/fs"
	"os"
	"path/filepath"
)

type Real struct {
	root string
}

var (
	_ ReadFS  = (*Real)(nil)
	_ WriteFS = (*Real)(nil)
)

func NewReal(root string) Real {
	return Real{root: root}
}

func (r Real) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(r.root, name))
}

func (r Real) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(r.root, name))
}

func (r Real) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.Join(r.root, name))
}

func (r Real) Create(name string) (File, error) {
	return os.Create(filepath.Join(r.root, name))
}

func (r Real) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(r.root, name))
}

func (r Real) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(r.root, path), perm)
}

func (r Real) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(r.root, name), data, perm)
}
