package memory

import (
	"io/fs"
	"maps"
	"slices"

	"go.ufukty.com/gommons/pkg/tree"
)

func (d Dir) strings() []string {
	ss := []string{}
	for _, name := range slices.Sorted(maps.Keys(d.entries)) {
		c := d.entries[name]
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

func find(d *Dir) ([]string, error) {
	ss := []string{}
	err := fs.WalkDir(d, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		ss = append(ss, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ss, nil
}
