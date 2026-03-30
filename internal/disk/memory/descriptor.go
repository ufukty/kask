package memory

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

var (
	ErrClosed        = fmt.Errorf("closed")
	ErrIsDir         = fmt.Errorf("node is a directory")
	ErrIsFile        = fmt.Errorf("node is a file")
	ErrNoSpace       = fmt.Errorf("no space")
	ErrUninitialized = fmt.Errorf("uninitialized")
)

type info struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
	sys     any
}

// As in [fs.FileInfo]
func (fi info) Name() string       { return fi.name }
func (fi info) Size() int64        { return fi.size }
func (fi info) Mode() fs.FileMode  { return fi.mode }
func (fi info) ModTime() time.Time { return fi.modTime }
func (fi info) IsDir() bool        { return fi.isDir }
func (fi info) Sys() any           { return fi.sys }

var _ fs.FileInfo = (*info)(nil)

// used for both the files and "dir files".
type descriptor struct {
	data any // [*Dir] | [*File]
	pos  int // a byte offset or dir entry.
	info info
}

var (
	_ io.WriteCloser = (*descriptor)(nil)
	_ fs.File        = (*descriptor)(nil)
	_ io.ReadCloser  = (*descriptor)(nil)
	_ fs.ReadDirFile = (*descriptor)(nil)
)

// As in [io.Writer]
// TODO: consider forwarding [fd.pos] as bytes written
func (d *descriptor) Write(p []byte) (n int, err error) {
	if d == nil {
		return -1, ErrUninitialized
	}
	f, ok := d.data.(*File)
	if !ok {
		return 0, ErrIsDir
	}
	if d.data == nil {
		return 0, ErrClosed
	}
	*f = append(*f, p...)
	d.info.size += int64(len(p))
	return len(p), nil
}

// As in [io.Closer]
func (d *descriptor) Close() error {
	if d == nil {
		return nil
	}
	d.data = nil
	return nil
}

// As in [fs.StatFS]
func (d *descriptor) Stat() (fs.FileInfo, error) {
	return d.info, nil
}

// Read writes the unread portion of file content shorter than the len(p).
// It returns [io.EOF] when there is nothing to return.
// Thus, it may return nil with data less than len(p).
// As in [io.Reader] and [fs.File]
func (d *descriptor) Read(p []byte) (int, error) {
	if d == nil {
		return -1, ErrUninitialized
	}
	f, ok := d.data.(*File)
	if !ok {
		return 0, ErrIsDir
	}
	if d.data == nil {
		return 0, ErrClosed
	}
	rem := len(*f) - d.pos
	if rem > 0 && len(p) == 0 {
		return 0, ErrNoSpace
	}
	if d.pos >= len(*f) {
		return 0, io.EOF
	}
	n := copy(p, (*f)[d.pos:])
	d.pos += n
	return n, nil
}

// As in [fs.dir]
func (d *descriptor) ReadDir(n int) ([]fs.DirEntry, error) {
	if d == nil {
		return nil, ErrUninitialized
	}
	di, ok := d.data.(*Dir)
	if !ok {
		return nil, ErrIsFile
	}
	if d.data == nil {
		return nil, ErrClosed
	}
	return di.ReadDir(".")
}
