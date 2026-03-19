package memory

import (
	"io"
	"io/fs"
	"time"

	"go.ufukty.com/kask/internal/disk"
)

type (
	fileInfo struct {
		name    string
		size    int64
		mode    fs.FileMode
		modTime time.Time
		isDir   bool
		sys     any
	}
	descriptor struct {
		file *File
		pos  int
		info fileInfo
	}
	File []byte
	Dir  map[string]any
)

// write
var (
	_ io.WriteCloser = (*descriptor)(nil)
	_ disk.WriteFS   = (*Dir)(nil)
)

// read
var (
	_ fs.FileInfo   = (*fileInfo)(nil)
	_ fs.File       = (*descriptor)(nil)
	_ io.Reader     = (*descriptor)(nil)
	_ io.ReadCloser = (*descriptor)(nil)
	_ fs.FS         = (*Dir)(nil)
	_ fs.ReadFileFS = (*Dir)(nil)
)
