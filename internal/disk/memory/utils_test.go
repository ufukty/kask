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
		if d, ok := c.(Dir); ok {
			ss = append(ss, tree.List(name, d.strings()))
		} else if _, ok := c.(File); ok {
			ss = append(ss, name)
		}
	}
	return ss
}

func (d Dir) String() string {
	return tree.List(".", d.strings())
}

func find(d Dir) []string {
	ss := []string{""}
	for name, s := range d {
		switch s := s.(type) {
		case Dir:
			for _, c := range find(s) {
				ss = append(ss, filepath.Join(name, c))
			}
		case File:
			ss = append(ss, name)
		}
	}
	return ss
}
