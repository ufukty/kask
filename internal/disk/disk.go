package disk

import (
	"io"
	"io/fs"
	_ "os"
)

// [ReadFS] implementations need to guarantee interoperability with standard
// library utilities such as [fs.WalkDir] and [template.template.ParseFS].
// TODO: [fs.GlobFS], [fs.ReadLinkFS], [fs.SubFS]
type ReadFS interface {
	fs.FS
	fs.ReadFileFS
	fs.ReadDirFS
	fs.StatFS
}

// TODO: remove once the [fs.File] added a [Write] and update [WriteFS.Create].
type File interface {
	fs.File
	io.Writer
}

// [WriteFS] implementations need to guarantee interoperability with the builder
// package. Methods need to share the same constraints with the [fs.ValidPath]
// on inputs.
type WriteFS interface {
	Create(name string) (File, error)                           // As in [os.Create].
	MkdirAll(path string, perm fs.FileMode) error               // As in [os.MkdirAll]
	WriteFile(path string, data []byte, perm fs.FileMode) error // As in [os.WriteFile]
}
