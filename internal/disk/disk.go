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

type ReadWriteFile interface {
	io.Writer
	fs.File
}

// Writable [fs.FS] for unit testing.
type WriteFS interface {
	Create(name string) (ReadWriteFile, error)
	MkdirAll(path string) error
	WriteFile(path string, data []byte) error
}

type ReadWriteFS interface {
	ReadFS
	WriteFS
}
