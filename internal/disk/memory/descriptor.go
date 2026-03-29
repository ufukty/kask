package memory

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

var ErrNoSpace = fmt.Errorf("no space")

// As in [io.Writer]
// TODO: consider forwarding [fd.pos] as bytes written
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
	n := copy(p, (*fd.file)[fd.pos:])
	fd.pos += n
	fd.info.size += int64(n)
	return n, nil
}
