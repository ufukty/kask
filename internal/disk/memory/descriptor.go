package memory

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

var ErrNoSpace = fmt.Errorf("no space")

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
