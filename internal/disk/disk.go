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

// Writable [fs.FS] for unit testing.
type WriteFS interface {
	Create(name string) (io.WriteCloser, error)
	MkdirAll(path string) error
	WriteFile(path string, data []byte) error
}
