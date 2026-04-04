package memory

import (
	"fmt"
	"io"
	"io/fs"
	slashpath "path"
	"slices"
	"strings"
	"time"

	"go.ufukty.com/kask/internal/disk"
)

type File struct {
	data    []byte
	mode    fs.FileMode
	modTime time.Time
}

// use [New] to construct
type Dir struct {
	entries map[string]any // [*Dir] | [*File]
	index   []string
	mode    fs.FileMode
	modTime time.Time
}

var (
	_ disk.WriteFS = (*Dir)(nil)
	_ disk.ReadFS  = (*Dir)(nil)
)

func newDir(perm fs.FileMode) *Dir {
	d := Dir{
		entries: map[string]any{},
		index:   []string{},
		mode:    fs.ModeDir | perm,
		modTime: time.Now(),
	}
	return &d
}

func New() *Dir {
	return newDir(0o755)
}

func (d *Dir) create(path string, perm fs.FileMode) (disk.File, error) {
	node, err := locate(d, slashpath.Dir(path))
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	dir, ok := node.(*Dir)
	if !ok {
		return nil, fmt.Errorf("destination should be a directory")
	}
	name := slashpath.Base(path)
	if name == "" {
		return nil, fmt.Errorf("unexpected empty name")
	}
	if _, ok := dir.entries[name]; ok {
		return nil, fmt.Errorf("exists")
	}
	f := &File{mode: perm, modTime: time.Now()}
	dir.entries[name] = f
	dir.insertIndex(name)
	dir.modTime = time.Now()
	fd := &handle{name: name, data: f, pos: 0}
	return fd, nil
}

// As in [disk.WriteFS]
func (d *Dir) Create(path string) (disk.File, error) {
	return d.create(path, 0o666)
}

// As in [disk.WriteFS]
func (d *Dir) MkdirAll(path string, perm fs.FileMode) error {
	if path == "." {
		return nil
	}
	cursor := d
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	ss := strings.Split(path, "/")
	for i, s := range ss {
		if isForbidden(s) {
			return fmt.Errorf("destination passes through an invalid node: %s", highlight(ss, i))
		}
		node, ok := cursor.entries[s]
		if ok {
			dir, ok := node.(*Dir)
			if !ok {
				return fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
			}
			cursor = dir
		} else {
			child := newDir(perm)
			cursor.entries[s] = child
			cursor.insertIndex(s)
			cursor.modTime = time.Now()
			cursor = child
		}
	}
	return nil
}

func (d *Dir) insertIndex(name string) {
	i := 0
	for ; i < len(d.index) && d.index[i] < name; i++ {
		/* i like to move it */
	}
	d.index = slices.Insert(d.index, i, name)
}

// As in [disk.WriteFS]
func (d *Dir) WriteFile(name string, data []byte, perm fs.FileMode) error {
	f, err := d.create(name, perm)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// As in [fs.FS]
func (d *Dir) Open(path string) (fs.File, error) {
	p, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	return &handle{name: slashpath.Base(path), data: p, pos: 0}, nil
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

// As in [fs.StatFS]
func (d *Dir) Stat(path string) (fs.FileInfo, error) {
	node, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	return info{name: slashpath.Base(path), node: node}, nil
}

// As in [fs.ReadDirFS]
func (d *Dir) ReadDir(path string) ([]fs.DirEntry, error) {
	node, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	dir, ok := node.(*Dir)
	if !ok {
		return nil, ErrIsFile
	}
	es, err := entries(dir, 0, -1)
	if err == io.EOF {
		return []fs.DirEntry{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("entries: %w", err)
	}
	return es, nil
}
