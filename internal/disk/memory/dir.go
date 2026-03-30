package memory

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"go.ufukty.com/kask/internal/disk"
)

type File []byte

// use [New] to construct
type Dir map[string]any

var (
	_ disk.WriteFS  = (*Dir)(nil) // write
	_ disk.ReadFS   = (*Dir)(nil) // read
	_ fs.FS         = (*Dir)(nil) // read
	_ fs.ReadDirFS  = (*Dir)(nil) // read
	_ fs.ReadFileFS = (*Dir)(nil) // read
	_ fs.StatFS     = (*Dir)(nil) // read
)

func New() *Dir {
	d := Dir{}
	d["."] = &d
	return &d
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
		data: f,
		pos:  0,
		info: fileInfo(f, name),
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

func fileInfo(node any, base string) info {
	file, isFile := node.(*File)
	if isFile {
		return info{
			name:    base,
			size:    int64(len(*file)),
			mode:    0o666,
			modTime: time.Now(),
			isDir:   false,
		}
	}
	return info{
		name:    base,
		size:    0,
		mode:    fs.ModeDir | 0o755,
		modTime: time.Now(),
		isDir:   true,
	}
}

// As in [fs.FS]
func (d *Dir) Open(path string) (fs.File, error) {
	p, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	fi := fileInfo(p, filepath.Base(path))
	return &descriptor{data: p, pos: 0, info: fi}, nil
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
	return fileInfo(node, filepath.Base(path)), nil
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
	ds := []fs.DirEntry{}
	for name, node := range *dir {
		if name == "." || name == ".." {
			continue
		}
		fi := fileInfo(node, name)
		di := entry{
			name:  name,
			isDir: isDir(node),
			mode:  fi.Mode(),
			info:  fi,
		}
		ds = append(ds, di)
	}
	return ds, nil
}
