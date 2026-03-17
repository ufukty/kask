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

func (f File) Write(p []byte) (n int, err error) {
	f = append(f, p...)
	return len(p), nil
}

func (f File) Close() error {
	return nil
}

type Dir map[string]any

func (r Dir) lastParent(ss []string) (Dir, error) {
	p := r
	for i, s := range ss[:max(0, len(ss)-1)] {
		n, ok := p[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
		d, ok := n.(Dir)
		if !ok {
			return nil, fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
		}
		p = d
	}
	return p, nil
}

func (r Dir) Create(path string) (io.WriteCloser, error) {
	if path == "" {
		return nil, fmt.Errorf("file name can't be empty")
	}
	ss := strings.Split(path, "/")
	p, err := r.lastParent(strings.Split(path, "/"))
	if err != nil {
		return nil, err
	}
	name := ss[len(ss)-1]
	if _, ok := p[name]; ok {
		return nil, fmt.Errorf("target already exists: %s", highlight(ss, len(ss)-1))
	}
	f := File{}
	p[name] = f
	return f, nil
}

func (r Dir) MkdirAll(path string) error {
	return nil
}

func (r Dir) WriteFile(name string, data []byte) error {
	f, err := r.Create(name)
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	f.Write(data)
	return nil
}
