package disk

import (
	"io"
	"io/fs"
)

// [ReadFS] implementations need to guarantee interoperability with standard
// library utilities such as [fs.WalkDir] and [template.template.ParseFS].
type ReadFS interface {
	fs.FS
	fs.ReadFileFS
	fs.ReadDirFS
	fs.StatFS
}

// [WriteFS] implementations need to guarantee interoperability with the builder
// package. Methods need to share the same constraints with the [fs.ValidPath]
// on inputs.
type WriteFS interface {
	// As in [os.Create]
	Create(name string) (io.WriteCloser, error)
	// As in [os.MkdirAll]
	MkdirAll(path string, perm fs.FileMode) error
	// As in [os.WriteFile]
	WriteFile(path string, data []byte, perm fs.FileMode) error
}
