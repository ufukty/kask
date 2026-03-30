// Package memory contains a file system implementation [Dir] stores data on
// memory, and its utilities that together allow its use with the standard
// library utilities such as [fs.WalkDir] or [template.Template]. It is
// intended to be used in testing only.
//
// [Dir] partially conforms the fstest.TestFS expectations. It supports
// absolute paths and redundant segments as builder utilizes. It doesn't
// support file timestamps and custom permissions.
package memory

import (
	"fmt"
	"io/fs"
	"maps"
	"path/filepath"
	"slices"
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

func readDir(d *Dir, path string, pos, n int) ([]fs.DirEntry, error) {
	node, err := locate(d, path)
	if err != nil {
		return nil, fmt.Errorf("locate: %w", err)
	}
	dir, ok := node.(*Dir)
	if !ok {
		return nil, ErrIsFile
	}
	ds := []fs.DirEntry{}
	for _, name := range slices.Sorted(maps.Keys(dir.entries))[pos:min(pos+n, len(dir.entries))] {
		if name == "." || name == ".." {
			continue
		}
		node := dir.entries[name]
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
