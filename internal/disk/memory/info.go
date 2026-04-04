package memory

import (
	"fmt"
	"io/fs"
	"time"
)

func size(node any) int64 {
	switch node := node.(type) {
	case *Dir:
		return 0
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

func is[T any](a any) bool {
	_, ok := a.(T)
	return ok
}

type info struct {
	name string
	node any // [*Dir] | [*File]
}

var _ fs.FileInfo = (*info)(nil)

// Unlike other [info] methods, [info.Name] returns the last value saved at the stat time.
func (i info) Name() string       { return i.name }           // As in [fs.FileInfo.Name]
func (i info) Size() int64        { return size(i.node) }     // As in [fs.FileInfo.Size]
func (i info) Mode() fs.FileMode  { return mode(i.node) }     // As in [fs.FileInfo.Mode]
func (i info) ModTime() time.Time { return modTime(i.node) }  // As in [fs.FileInfo.ModTime]
func (i info) IsDir() bool        { return is[*Dir](i.node) } // As in [fs.FileInfo.IsDir]
func (i info) Sys() any           { return nil }              // As in [fs.FileInfo.Sys]

type entry struct {
	name string
	node any // [*Dir] | [*File]
}

var _ fs.DirEntry = (*entry)(nil)

// Unlike other [entry] methods, [entry.Name] returns the last value saved at the stat time.
func (e entry) Name() string               { return e.name }                                // As in [fs.DirEntry.Name]
func (e entry) IsDir() bool                { return mode(e.node).IsDir() }                  // As in [fs.DirEntry.IsDir]
func (e entry) Type() fs.FileMode          { return mode(e.node).Type() }                   // As in [fs.DirEntry.Type]
func (e entry) Info() (fs.FileInfo, error) { return info{name: e.name, node: e.node}, nil } // As in [fs.DirEntry.Info]
