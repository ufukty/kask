package memory

import (
	"fmt"
	"io/fs"
	"strings"
)

var (
	_ fs.FS         = (*Dir)(nil)
	_ fs.ReadFileFS = (*Dir)(nil)
)

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
	return f, nil
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
