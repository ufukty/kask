package memory

import (
	"io/fs"
	"os"

	"go.ufukty.com/kask/internal/disk"
)

type File struct {
	content []byte
}

var _ disk.File = (*File)(nil)

func newFile() *File {
	return &File{content: []byte{}}
}

func (fom File) Stat() (fs.FileInfo, error)

func (fom File) Read([]byte) (int, error)

func (fom File) Close() error

func (fom File) Write(p []byte) (n int, err error)

type Dir struct {
	files map[string]File
}

var _ disk.ReadWriteFS = (*Dir)(nil)

func newDir() *Dir {
	return &Dir{}
}

func (r Dir) Open(name string) (fs.File, error)

func (r Dir) ReadFile(name string) ([]byte, error)

func (r Dir) ReadDir(name string) ([]os.DirEntry, error)

func (r Dir) Create(name string) (disk.File, error)

func (r Dir) Stat(name string) (os.FileInfo, error)

func (r Dir) MkdirAll(path string) error

func (r Dir) WriteFile(name string, data []byte) error
