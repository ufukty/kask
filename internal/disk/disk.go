package disk

import (
	"io"
	"io/fs"
)

type ReadFS interface {
	fs.ReadFileFS
	fs.ReadDirFS
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
