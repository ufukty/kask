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
	_ fs.StatFS     = (*Dir)(nil) // read
	_ fs.ReadFileFS = (*Dir)(nil) // read
	_ fs.ReadDirFS  = (*Dir)(nil) // read
	_ fs.FS         = (*Dir)(nil) // read
	_ disk.ReadFS   = (*Dir)(nil) // read
	_ disk.WriteFS  = (*Dir)(nil) // write
)

func New() *Dir {
	d := Dir{}
	d["."] = &d
	return &d
}

func rewriteByTheRoot(path string) string {
	if path == "/" {
		return "."
	} else {
		return strings.TrimPrefix(path, "/")
	}
}

func findRoot(dir *Dir) (*Dir, error) {
	for range 100 {
		node, ok := (*dir)[".."]
		if !ok {
			return dir, nil
		}
		parent, ok := (node).(*Dir)
		if !ok {
			return nil, fmt.Errorf("unexpected non-dir parent: %T", node)
		}
		dir = parent
	}
	return nil, fmt.Errorf("directory depth limit is exceeded")
}

// Allows [*File] only at leaves.
// Returns either the [*Dir] or [*File] pointed by [string].
func locate(entry *Dir, path string) (any, error) {
	path = filepath.Clean(path)
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	var cursor any
	if filepath.IsAbs(path) {
		root, err := findRoot(entry)
		if err != nil {
			return nil, fmt.Errorf("finding root: %w", err)
		}
		cursor = root
		path = rewriteByTheRoot(path)
	} else {
		cursor = entry
	}
	ss := strings.Split(path, "/")
	for i, s := range ss {
		if s == "" {
			return nil, fmt.Errorf("destination passes through a node with empty name")
		}
		dir, ok := cursor.(*Dir)
		if !ok {
			// destination has passed through a [*File] in previous iteration
			return nil, fmt.Errorf("destination passes through a file: %s", highlight(ss, i-1))
		}
		cursor, ok = (*dir)[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
	}
	return cursor, nil
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
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
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
	f, ok := p.(*File)
	if !ok {
		ss := strings.Split(path, "/")
		return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, len(ss)-1))
	}
	name := filepath.Base(path)
	if name == "" {
		return nil, fmt.Errorf("unexpected empty name")
	}
	fd := &descriptor{
		file: f,
		pos:  0,
		info: fileInfo{
			name:    name,
			size:    int64(len(*f)),
			mode:    0o666,
			modTime: time.Time{},
			isDir:   false,
			sys:     nil,
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

// As in [fs.StatFS]
func (d *Dir) Stat(path string) (fs.FileInfo, error) {
	node, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	file, isFile := node.(*File)
	if isFile {
		return fileInfo{
			name:    filepath.Base(path),
			size:    int64(len(*file)),
			mode:    0o666,
			modTime: time.Now(),
			isDir:   false,
		}, nil
	} else {
		return fileInfo{
			name:    filepath.Base(path),
			size:    0,
			mode:    fs.ModeDir | 0o755,
			modTime: time.Now(),
			isDir:   true,
		}, nil
	}
}

type dirEntry struct {
	name  string
	isDir bool
	typee fs.FileMode
	info  fs.FileInfo
}

// As in [fs.DirEntry]
func (de dirEntry) Name() string               { return de.name }
func (de dirEntry) IsDir() bool                { return de.isDir }
func (de dirEntry) Type() fs.FileMode          { return de.typee }
func (de dirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

var _ fs.DirEntry = (*dirEntry)(nil)

// As in [fs.ReadDirFS]
func (d *Dir) ReadDir(path string) ([]fs.DirEntry, error) {
	node, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	dir, ok := node.(*Dir)
	if !ok {
		return nil, fmt.Errorf("not a directory")
	}
	ds := []fs.DirEntry{}
	for name := range *dir {
		if name == "." || name == ".." {
			continue
		}
		fi, err := dir.Stat(name)
		if err != nil {
			return nil, fmt.Errorf("stat: %w", err)
		}
		di := dirEntry{
			name:  name,
			isDir: fi.IsDir(),
			typee: fi.Mode(),
			info:  fi,
		}
		ds = append(ds, di)
	}
	return ds, nil
}
