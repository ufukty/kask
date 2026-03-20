package memory

import (
	"fmt"
	"io"
	"io/fs"
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

func (d *Dir) findDir(ss []string) (*Dir, error) {
	p := d
	for i, s := range ss {
		if s == "" {
			return nil, fmt.Errorf("empty name")
		}
		if s == "." {
			continue
		}
		// if s == ".." {
		// 	if d.Parent != nil {
		// 		d = d.Parent
		// 	}
		// 	continue
		// }
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
	return p, nil
}

// As in [fs.FS]
func (d *Dir) Open(path string) (fs.File, error) {
	if path == "" {
		return nil, fmt.Errorf("file path can't be empty")
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
	inode, ok := (*p)[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	f, ok := inode.(*File)
	if !ok {
		return nil, fs.ErrInvalid
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
			sys:     ss,
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
	if path == "" {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: fmt.Errorf("file path can't be empty")}
	}
	ss := strings.Split(path, "/")
	p, err := d.findDir(ss[:len(ss)-1])
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: err}
	}
	name := ss[len(ss)-1]
	if name == "" {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: fmt.Errorf("unexpected empty name")}
	}
	inode, ok := (*p)[name]
	if !ok {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist}
	}
	f, isFile := inode.(*File)
	if isFile {
		return fileInfo{
			name:    name,
			size:    int64(len(*f)),
			mode:    0o666,
			modTime: time.Now(),
			isDir:   false,
		}, nil
	} else {
		return fileInfo{
			name:    name,
			size:    0,
			mode:    fs.ModeDir | 0o755,
			modTime: time.Now(),
			isDir:   true,
		}, nil
	}
}

// As in [fs.ReadDirFS]
func (d *Dir) ReadDir(path string) ([]fs.DirEntry, error) {
	if path == "" {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: fmt.Errorf("file path can't be empty")}
	}
	p, err := d.findDir(strings.Split(path, "/"))
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: path, Err: err}
	}
	ds := []fs.DirEntry{}
	for name := range *p {
		fi, err := d.Stat(name)
		if err != nil {
			return nil, &fs.PathError{Op: "stat", Path: path, Err: fmt.Errorf("stat: %w", err)}
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
