package memory

import (
	"maps"
	"path/filepath"
	"slices"

	"go.ufukty.com/gommons/pkg/tree"
)

func (d Dir) strings() []string {
	ss := []string{}
	for _, name := range slices.Sorted(maps.Keys(d)) {
		c := d[name]
		if d, ok := c.(*Dir); ok {
			ss = append(ss, tree.List(name, d.strings()))
		} else if _, ok := c.(*File); ok {
			ss = append(ss, name)
		}
	}
	return ss
}

func (d Dir) String() string {
	return tree.List(".", d.strings())
}

// use dot for path
func walkDir(root any, path string, v func(string, any) bool) bool {
	if !v(path, root) {
		return false
	}
	if d, ok := root.(*Dir); ok {
		for name, sub := range *d {
			if name == "." || name == ".." {
				continue
			}
			if !walkDir(sub, filepath.Join(path, name), v) {
				return false
			}
		}
	}
	return true
}

func find(d *Dir) []string {
	ss := []string{}
	walkDir(d, ".", func(s string, a any) bool {
		ss = append(ss, s)
		return true
	})
	return ss
}
