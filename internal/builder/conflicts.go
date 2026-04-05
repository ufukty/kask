package builder

import (
	"fmt"
	"strings"

	"go.ufukty.com/kask/internal/builder/directory"
)

func (b *builder) checkCompetingEntries(dir *directory.Dir) error {
	children := map[string]int{}
	for _, subdir := range dir.Subdirs {
		children[subdir.Name] = 1
	}
	for _, page := range dir.Pages {
		if has(children, page) {
			children[page] = -1
		}
		children[page]++
	}
	duplicates := []string{}
	for child, freq := range children {
		if freq > 1 {
			duplicates = append(duplicates, child)
		}
	}
	if len(duplicates) > 0 {
		return fmt.Errorf("multiple entries sharing the same path for those: %s", strings.Join(duplicates, ", "))
	}
	for _, sub := range dir.Subdirs {
		if err := b.checkCompetingEntries(sub); err != nil {
			return fmt.Errorf("%q: %w", sub.Name, err)
		}
	}
	return nil
}
