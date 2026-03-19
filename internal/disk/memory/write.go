package memory

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
)

// As in [io.Writer]
func (fd *descriptor) Write(p []byte) (n int, err error) {
	if fd.file == nil {
		return 0, fmt.Errorf("closed")
	}
	*fd.file = append(*fd.file, p...)
	return len(p), nil
}

// As in [io.Closer]
func (fd *descriptor) Close() error {
	fd.file = nil
	return nil
}

// As in [disk.WriteFS]
func (d *Dir) Create(path string) (io.WriteCloser, error) {
	if path == "" {
		return nil, fmt.Errorf("file name can't be empty")
	}
	ss := strings.Split(path, "/")
	p, err := d.findDir(ss[:len(ss)-1])
	if err != nil {
		return nil, err
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
	fd := &descriptor{
		file: f,
		pos:  0,
		info: fileInfo{
			name:    name,
			size:    int64(len(*f)),
			mode:    fs.ModeAppend,
			modTime: time.Now(),
			isDir:   false,
			sys:     nil,
		},
	}
	return fd, nil
}

// As in [disk.WriteFS]
func (d *Dir) MkdirAll(path string) error {
	if path == "" {
		return fmt.Errorf("file name can't be empty")
	}
	ss := strings.Split(path, "/")
	p := d
	for i, s := range ss {
		if s == "" {
			return fmt.Errorf("unexpected empty name")
		}
		n, ok := (*p)[s]
		if ok {
			d, ok := n.(*Dir)
			if !ok {
				return fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
			}
			p = d
		} else {
			c := &Dir{}
			(*p)[s] = c
			p = c
		}
	}
	return nil
}

// As in [disk.WriteFS]
func (d *Dir) WriteFile(name string, data []byte) error {
	f, err := d.Create(name)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	f.Write(data)
	return nil
}
