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

type info struct {
	name string
	node any // [*Dir] | [*File]
}

var _ fs.FileInfo = (*info)(nil)

// Unlike other [info] methods, [info.Name] returns the value saved
// at the stat time.
// As in [fs.FileInfo.Name]
func (i info) Name() string {
	return i.name
}

// As in [fs.FileInfo.Size]
func (i info) Size() int64 {
	switch node := i.node.(type) {
	case *Dir:
		return int64(len(node.entries))
	case *File:
		return int64(len(node.data))
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.Mode]
func (i info) Mode() fs.FileMode {
	switch node := i.node.(type) {
	case *Dir:
		return node.mode
	case *File:
		return node.mode
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.ModTime]
func (i info) ModTime() time.Time {
	switch node := i.node.(type) {
	case *Dir:
		return node.modTime
	case *File:
		return node.modTime
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.IsDir]
func (i info) IsDir() bool {
	_, ok := i.node.(*Dir)
	return ok
}

// As in [fs.FileInfo.Sys]
func (i info) Sys() any {
	return nil
}

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
		return -1, ErrUninitialized
	}
	if d.data == nil {
		return 0, ErrClosed
	}
	f, ok := d.data.(*File)
	if !ok {
		return 0, ErrIsDir
	}
	f.data = append(f.data, p...)
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
	return info{name: d.name, node: d.data}, nil
}

// Read writes the unread portion of file content shorter than the len(p).
// It returns [io.EOF] when there is nothing to return.
// Thus, it may return nil with data less than len(p).
// As in [io.Reader] and [fs.File]
func (d *handle) Read(p []byte) (int, error) {
	if d == nil {
		return -1, ErrUninitialized
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
