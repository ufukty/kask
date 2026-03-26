package memory

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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

// As in [fs.DirEntry]
func (de dirEntry) Name() string               { return de.name }
func (de dirEntry) IsDir() bool                { return de.isDir }
func (de dirEntry) Type() fs.FileMode          { return de.typee }
func (de dirEntry) Info() (fs.FileInfo, error) { return de.info, nil }

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
	for i := 0; i < end-start; i++ {
		p[i] = (*fd.file)[i+start]
	}
	fd.pos += len(p)
	return len(*fd.file), nil
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
		// shouldn't execute when it is a [*File],
		// which is only possible at a leaf.
		dir, ok := cursor.(*Dir)
		if !ok {
			return nil, fmt.Errorf("destination passes through a file: %s", highlight(ss, i))
		}
		cursor, ok = (*dir)[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
	}
	return cursor, nil
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
