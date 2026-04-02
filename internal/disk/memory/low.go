package memory

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
)

// based on the [fs.ValidPath].
func isForbidden(s string) bool {
	return s == "" || s == "." || s == ".."
}

// Allows [*File] only at leaves.
// Returns either the [*Dir] or [*File] pointed by [string].
func locate(entry *Dir, path string) (any, error) {
	if path == "." {
		return entry, nil
	}
	if path == "" {
		return nil, fmt.Errorf("path is empty")
	}
	var cursor any = entry
	ss := strings.Split(path, "/")
	for i, s := range ss {
		if isForbidden(s) {
			return nil, fmt.Errorf("destination passes through an invalid node: %s", highlight(ss, i))
		}
		dir, ok := cursor.(*Dir)
		if !ok {
			// destination has passed through a [*File] in previous iteration
			return nil, fmt.Errorf("destination passes through a file: %s", highlight(ss, i-1))
		}
		cursor, ok = dir.entries[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
	}
	return cursor, nil
}

func entries(d *Dir, pos, n int) ([]fs.DirEntry, error) {
	ds := []fs.DirEntry{}
	from, to := pos, pos+n
	if n < 0 || len(d.index) < to {
		to = len(d.index)
	}
	if n > 0 && from == to {
		return nil, io.EOF
	}
	for _, name := range d.index[from:to] {
		node := d.entries[name]
		fi, err := fileInfo(node, name)
		if err != nil {
			return nil, fmt.Errorf("fileInfo %q: %w", name, err)
		}
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
