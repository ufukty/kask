package memory

import (
	"io/fs"
	"os"
	"time"

	"go.ufukty.com/kask/internal/disk"
)

type fileInfo struct {
	name    string
	size    int64
	mode    uint32
	modTime time.Time
	isDir   bool
	sys     any
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() fs.FileMode  { return fs.FileMode(fi.mode) }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return fi.isDir }
func (fi fileInfo) Sys() any           { return fi.sys }

type File struct {
	info    fileInfo
	content []byte
}

var _ disk.ReadWriteFile = (*File)(nil)

func newFile(name string) *File {
	return &File{
		content: []byte{},
		info:    fileInfo{name: name},
	}
}

func (f File) Stat() (fs.FileInfo, error) {
	return f.info, nil
}

func (f File) Read([]byte) (int, error)

func (f File) Close() error

func (f File) Write(p []byte) (n int, err error)

type Dir struct {
	files map[string]File
}

var _ disk.ReadWriteFS = (*Dir)(nil)

func New() *Dir {
	return &Dir{
		files: map[string]File{},
	}
}

func has[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]
	return ok
}

func (r Dir) Open(name string) (fs.File, error) {
	if !has(r.files, name) {
		return nil, os.ErrNotExist
	}
	return r.files[name], nil
}

func (r Dir) ReadFile(name string) ([]byte, error)

func (r Dir) ReadDir(name string) ([]os.DirEntry, error)

func (r Dir) Create(name string) (disk.ReadWriteFile, error)

func (r Dir) Stat(name string) (os.FileInfo, error)

func (r Dir) MkdirAll(path string) error

func (r Dir) WriteFile(name string, data []byte) error
