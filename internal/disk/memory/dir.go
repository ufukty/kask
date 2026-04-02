package memory

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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

func New() *Dir {
	d := Dir{
		entries: map[string]any{},
		index:   []string{},
		mode:    fs.ModeDir,
		modTime: time.Now(),
	}
	return &d
}

// As in [disk.WriteFS]
func (d *Dir) Create(path string) (disk.File, error) {
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
	if _, ok := dir.entries[name]; ok {
		return nil, fmt.Errorf("exists")
	}
	f := &File{}
	dir.entries[name] = f
	dir.insertIndex(name)
	fi, err := fileInfo(f, name)
	if err != nil {
		return nil, fmt.Errorf("fileInfo: %w", err)
	}
	fd := &handle{
		data: f,
		pos:  0,
		info: fi,
	}
	return fd, nil
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
			child := New()
			cursor.entries[s] = child
			cursor.insertIndex(s)
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
	f, err := d.Create(filepath.Clean(name))
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func fileInfo(node any, base string) (info, error) {
	switch node := node.(type) {
	case *File:
		return info{
			name:    base,
			size:    int64(len(node.data)),
			mode:    node.mode,
			modTime: node.modTime,
			isDir:   false,
		}, nil
	case *Dir:
		return info{
			name:    base,
			size:    0,
			mode:    node.mode,
			modTime: node.modTime,
			isDir:   true,
		}, nil
	default:
		return info{}, fmt.Errorf("unknown type: %T", node)
	}
}

// As in [fs.FS]
func (d *Dir) Open(path string) (fs.File, error) {
	p, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	fi, err := fileInfo(p, filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("fileInfo: %w", err)
	}
	return &handle{data: p, pos: 0, info: fi}, nil
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
	fi, err := fileInfo(node, filepath.Base(path))
	if err != nil {
		return nil, fmt.Errorf("fileInfo: %w", err)
	}
	return fi, nil
}

type entry struct {
	name  string
	isDir bool
	mode  fs.FileMode
	info  info
}

// As in [fs.DirEntry]
func (e entry) Name() string               { return e.name }
func (e entry) IsDir() bool                { return e.isDir }
func (e entry) Type() fs.FileMode          { return e.mode }
func (e entry) Info() (fs.FileInfo, error) { return e.info, nil }

var _ fs.DirEntry = (*entry)(nil)

func isDir(node any) bool {
	_, ok := node.(*Dir)
	return ok
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
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("entries: %w", err)
	}
	return es, nil
}
