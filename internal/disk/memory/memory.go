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
	dirEntry struct {
		name  string
		isDir bool
		typee fs.FileMode
		info  fs.FileInfo
	}
	descriptor struct {
		file *File
		pos  int
		info fileInfo
	}
	File []byte
	Dir  map[string]any
)

var (
	_ fs.FileInfo = (*fileInfo)(nil)
	_ fs.DirEntry = (*dirEntry)(nil)
)

var (
	// write
	_ io.WriteCloser = (*descriptor)(nil)
	// read
	_ fs.File       = (*descriptor)(nil)
	_ io.ReadCloser = (*descriptor)(nil)
)

var (
	// write
	_ disk.WriteFS = (*Dir)(nil)
	// read
	_ fs.FS         = (*Dir)(nil)
	_ fs.ReadFileFS = (*Dir)(nil)
	_ fs.ReadDirFS  = (*Dir)(nil)
	_ fs.StatFS     = (*Dir)(nil)
	_ disk.ReadFS   = (*Dir)(nil)
)
