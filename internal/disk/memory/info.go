package memory

import (
	"fmt"
	"io/fs"
	"time"
)

type info struct {
	name string
	node any // [*Dir] | [*File]
}

var _ fs.FileInfo = (*info)(nil)

// Unlike other [info] methods, [info.Name] returns the value saved
// at the stat time.
// As in [fs.FileInfo.Name]
func (i info) Name() string {
	return i.name
}

// As in [fs.FileInfo.Size]
func (i info) Size() int64 {
	switch node := i.node.(type) {
	case *Dir:
		return int64(len(node.entries))
	case *File:
		return int64(len(node.data))
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.Mode]
func (i info) Mode() fs.FileMode {
	switch node := i.node.(type) {
	case *Dir:
		return node.mode
	case *File:
		return node.mode
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.ModTime]
func (i info) ModTime() time.Time {
	switch node := i.node.(type) {
	case *Dir:
		return node.modTime
	case *File:
		return node.modTime
	default:
		panic(fmt.Sprintf("unexpected node type %T", i.node))
	}
}

// As in [fs.FileInfo.IsDir]
func (i info) IsDir() bool {
	_, ok := i.node.(*Dir)
	return ok
}

// As in [fs.FileInfo.Sys]
func (i info) Sys() any {
	return nil
}

type entry struct {
	name string
	node any // [*Dir] | [*File]
}

var _ fs.DirEntry = (*entry)(nil)

// As in [fs.DirEntry.Name]
func (e entry) Name() string {
	return e.name
}

// As in [fs.DirEntry.IsDir]
func (e entry) IsDir() bool {
	return mode(e.node).IsDir()
}

// As in [fs.DirEntry.Type]
func (e entry) Type() fs.FileMode {
	return mode(e.node).Type()
}

// As in [fs.DirEntry.Info]
func (e entry) Info() (fs.FileInfo, error) {
	return info{name: e.name, node: e.node}, nil
}
