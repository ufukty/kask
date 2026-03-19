package memory

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
)

var ErrNoSpace = fmt.Errorf("no space")

// As in [fs.FileInfo]
func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return fi.size }
func (fi fileInfo) Mode() fs.FileMode  { return fi.mode }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return fi.isDir }
func (fi fileInfo) Sys() any           { return fi.sys }

// As in [fs.StatFS]
func (fd *descriptor) Stat() (fs.FileInfo, error) {
	return fd.info, nil
}

// Read writes the unread portion of file content shorter than the len(p).
// It returns [io.EOF] when there is nothing to return.
// Thus, it may return nil with data less than len(p).
// As in [io.Reader]
func (fd *descriptor) Read(p []byte) (int, error) {
	if fd.file == nil {
		return 0, fmt.Errorf("closed")
	}
	rem := len(*fd.file) - fd.pos
	if rem > 0 && len(p) == 0 {
		return 0, ErrNoSpace
	}
	if fd.pos >= len(*fd.file) {
		return 0, io.EOF
	}
	start, end := fd.pos, min(fd.pos+len(p), len(*fd.file))
	for i := start; i < end; i++ {
		p[i] = (*fd.file)[i]
	}
	fd.pos += len(p)
	return len(*fd.file), nil
}

// As in [fs.FS]
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
	_, ok := (*p)[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	f := &File{}
	(*p)[name] = f
	fd := &descriptor{
		file: f,
		pos:  0,
		info: fileInfo{
			name:    name,
			size:    int64(len(*f)),
			mode:    fs.ModeAppend,
			modTime: time.Time{},
			isDir:   false,
			sys:     ss,
		},
	}
	return fd, nil
}

// As in [fs.ReadFileFS]
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
