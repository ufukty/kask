package memory

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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
	path = filepath.Clean(path)
	node, err := locate(d, filepath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	dir, ok := node.(*Dir)
	if !ok {
		return nil, fmt.Errorf("destination should be a directory")
	}
	name := filepath.Base(path)
	if name == "" {
		return nil, fmt.Errorf("unexpected empty name")
	}
	if _, ok := (*dir)[name]; ok {
		return nil, fmt.Errorf("exists")
	}
	f := &File{}
	(*dir)[name] = f
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
	path = filepath.Clean(path)
	var cursor *Dir
	if filepath.IsAbs(path) {
		root, err := findRoot(d)
		if err != nil {
			return fmt.Errorf("finding root: %w", err)
		}
		cursor = root
		path = rewriteByTheRoot(path)
	} else {
		cursor = d
	}
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	ss := strings.Split(path, "/")
	for i, s := range ss {
		if s == "" {
			return fmt.Errorf("unexpected empty name")
		}
		node, ok := (*cursor)[s]
		if ok {
			dir, ok := node.(*Dir)
			if !ok {
				return fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
			}
			cursor = dir
		} else {
			child := &Dir{}
			(*cursor)[s] = child
			(*child)["."] = child
			(*child)[".."] = cursor
			cursor = child
		}
	}
	return nil
}

// As in [disk.WriteFS]
func (d *Dir) WriteFile(name string, data []byte) error {
	name = filepath.Clean(name)
	f, err := d.Create(name)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func New() *Dir {
	d := Dir{}
	d["."] = &d
	return &d
}
