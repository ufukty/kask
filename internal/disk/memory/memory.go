package memory

import (
	"fmt"
	"io"
	"strings"

	"go.ufukty.com/kask/internal/disk"
)

var (
	_ io.WriteCloser = (*File)(nil)
	_ disk.WriteFS   = (*Dir)(nil)
)

type File []byte

func (f *File) Write(p []byte) (n int, err error) {
	if *f == nil {
		return 0, fmt.Errorf("closed")
	}
	*f = append(*f, p...)
	return len(p), nil
}

func (f *File) Close() error {
	*f = nil
	return nil
}

type Dir map[string]any

func (d *Dir) Create(path string) (io.WriteCloser, error) {
	if path == "" {
		return nil, fmt.Errorf("file name can't be empty")
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
	return f, nil
}

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

func (d *Dir) WriteFile(name string, data []byte) error {
	f, err := d.Create(name)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	f.Write(data)
	return nil
}
