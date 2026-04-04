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
	ErrUninitialized = fmt.Errorf("uninitialized")
)

// used for both the files and "dir files".
type handle struct {
	name string
	data any // [*Dir] | [*File]
	pos  int // a byte offset or dir entry.
}

var (
	_ fs.File        = (*handle)(nil)
	_ fs.ReadDirFile = (*handle)(nil)
	_ io.Closer      = (*handle)(nil)
	_ io.Reader      = (*handle)(nil)
	_ io.Writer      = (*handle)(nil)
)

// As in [io.Writer]
// TODO: consider forwarding [fd.pos] as bytes written
func (d *handle) Write(p []byte) (n int, err error) {
	if d == nil {
		return 0, ErrUninitialized
	}
	if d.data == nil {
		return 0, ErrClosed
	}
	f, ok := d.data.(*File)
	if !ok {
		return 0, ErrIsDir
	}
	f.data = append(f.data, p...)
	f.modTime = time.Now()
	d.pos += len(p)
	return len(p), nil
}

// As in [io.Closer]
func (d *handle) Close() error {
	if d == nil {
		return nil
	}
	d.data = nil
	return nil
}

// As in [fs.StatFS]
func (d *handle) Stat() (fs.FileInfo, error) {
	if d == nil {
		return nil, ErrUninitialized
	}
	if d.data == nil {
		return nil, ErrClosed
	}
	return info{name: d.name, node: d.data}, nil
}

// Read writes the unread portion of file content shorter than the len(p).
// It returns [io.EOF] when there is nothing to return.
// Thus, it may return nil with data less than len(p).
// As in [io.Reader] and [fs.File]
func (d *handle) Read(p []byte) (int, error) {
	if d == nil {
		return 0, ErrUninitialized
	}
	if d.data == nil {
		return 0, ErrClosed
	}
	f, ok := d.data.(*File)
	if !ok {
		return 0, ErrIsDir
	}
	if d.pos >= len(f.data) {
		return 0, io.EOF
	}
	n := copy(p, f.data[d.pos:])
	d.pos += n
	return n, nil
}

// As in [fs.ReadDirFile]
func (d *handle) ReadDir(n int) ([]fs.DirEntry, error) {
	if d == nil {
		return nil, ErrUninitialized
	}
	if d.data == nil {
		return nil, ErrClosed
	}
	di, ok := d.data.(*Dir)
	if !ok {
		return nil, ErrIsFile
	}
	es, err := entries(di, d.pos, n)
	if err == io.EOF {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("entries: %w", err)
	}
	d.pos += len(es)
	return es, nil
}
