package disk

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type ReadFS interface {
	fs.ReadFileFS
	fs.ReadDirFS
	fs.GlobFS
	fs.StatFS
}

type File interface {
	io.Writer
	fs.File
}

// Writable [fs.FS] for unit testing.
type WriteFS interface {
	Create(name string) (File, error)
	MkdirAll(path string) error
	WriteFile(path string, data []byte) error
}

type ReadWriteFS interface {
	ReadFS
	WriteFS
}

type Real struct {
	root string
}

var _ ReadWriteFS = (*Real)(nil)

func NewReal(root string) Real {
	return Real{root: root}
}

func (r Real) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(r.root, name))
}

func (r Real) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(r.root, name))
}

func (r Real) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(filepath.Join(r.root, name))
}

func (r Real) Glob(pattern string) ([]string, error) {
	return fs.Glob(r, pattern)
}

func (r Real) Create(name string) (File, error) {
	return os.Create(name)
}

func (r Real) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (r Real) MkdirAll(path string) error {
	return os.MkdirAll(filepath.Join(r.root, path), 0o755)
}

func (r Real) WriteFile(name string, data []byte) error {
	return os.WriteFile(filepath.Join(r.root, name), data, 0o666)
}
