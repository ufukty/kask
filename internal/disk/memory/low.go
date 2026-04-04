package memory

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"time"
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
			return nil, fmt.Errorf("destination passes through an missing directory: %s", highlight(ss, i))
		}
	}
	return cursor, nil
}

func size(node any) int64 {
	switch node := node.(type) {
	case *Dir:
		return int64(len(node.entries))
	case *File:
		return int64(len(node.data))
	default:
		panic(fmt.Sprintf("unexpected node type %T", node))
	}
}

func mode(node any) fs.FileMode {
	switch node := node.(type) {
	case *Dir:
		return node.mode
	case *File:
		return node.mode
	default:
		panic(fmt.Sprintf("unexpected node type %T", node))
	}
}

func modTime(node any) time.Time {
	switch node := node.(type) {
	case *Dir:
		return node.modTime
	case *File:
		return node.modTime
	default:
		panic(fmt.Sprintf("unexpected node type %T", node))
	}
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
		ds = append(ds, entry{name: name, node: d.entries[name]})
	}
	return ds, nil
}
