package memory

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
)

var ErrNoSpace = fmt.Errorf("no space")

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return fi.isDir }
func (fi fileInfo) Sys() any           { return fi.sys }

func (fd *descriptor) Stat() (fs.FileInfo, error) { return fd.info, nil }

// [dile.Read] writes the unread portion of file content shorter
// than the len(p). It returns [io.EOF] when there is nothing
// to return. Thus, it may return nil with data less than len(p).
func (fd *descriptor) Read(p []byte) (int, error) {
	rem := len(*fd.file) - fd.pos
	if rem > 0 && len(p) == 0 {
		return 0, ErrNoSpace
	}
	if fd.pos >= len(*fd.file) {
		return 0, io.EOF
	}
	p = (*fd.file)[:len(p)]
	fd.pos += len(p)
	return len(*fd.file), nil
}

func (d *Dir) Open(path string) (fs.File, error) {
	if path == "" {
		return nil, fmt.Errorf("file path can't be empty")
	}
	ss := strings.Split(path, "/")
	p := d
	for i, s := range ss[:len(ss)-1] {
		n, ok := (*p)[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
		d, ok := n.(*Dir)
		if !ok {
			return nil, fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
		}
		p = d
	}
	name := ss[len(ss)-1]
	if name == "" {
		return nil, fmt.Errorf("unexpected empty name")
	}
	if _, ok := (*p)[name]; ok {
		return nil, fmt.Errorf("target already exists: %s", highlight(ss, len(ss)-1))
	}
	f := &File{}
	(*p)[name] = f
	fd := &descriptor{file: f, pos: 0} // FIXME: add [FileInfo]
	return fd, nil
}

func (d *Dir) ReadFile(name string) ([]byte, error) {
	f, err := d.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}
	b := make([]byte, s.Size())
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	return b, nil
}
