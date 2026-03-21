package disk

import (
	"io"
	"io/fs"
)

// [ReadFS] needs to conform [fs.FS] to work with [templates.Template.ParseFS]
// and other stdlib symbols.
// Which is unfortunate because [fs.File] is read-only.
// [WriteFS] is needed to make units testable.
// Hence, they are split into two interfaces.
type (
	ReadFS interface {
		fs.ReadFileFS
		fs.ReadDirFS
		fs.StatFS
	}
	WriteFS interface {
		Create(name string) (io.WriteCloser, error)
		MkdirAll(path string) error
		WriteFile(path string, data []byte) error
	}
)
