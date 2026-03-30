package memory

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

func rewriteByTheRoot(path string) string {
	if path == "/" {
		return "."
	} else {
		return strings.TrimPrefix(path, "/")
	}
}

func findRoot(dir *Dir) (*Dir, error) {
	for range 100 {
		node, ok := dir.entries[".."]
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
		cursor, ok = dir.entries[s]
		if !ok {
			return nil, fmt.Errorf("destination passes through an unexisting directory: %s", highlight(ss, i))
		}
	}
	return cursor, nil
}

func entries(d *Dir, pos, n int) []fs.DirEntry {
	ds := []fs.DirEntry{}
	for _, name := range d.index[pos:min(pos+n, len(d.index))] {
		node := d.entries[name]
		fi := fileInfo(node, name)
		di := entry{
			name:  name,
			isDir: isDir(node),
			mode:  fi.Mode(),
			info:  fi,
		}
		ds = append(ds, di)
	}
	return ds
}
