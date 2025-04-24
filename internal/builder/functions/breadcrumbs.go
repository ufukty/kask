package functions

import (
	"slices"

	"github.com/ufukty/kask/internal/compiler/builder/directory"
)

func Breadcrumbs(n *directory.Node) []*directory.Node {
	l := []*directory.Node{}
	cursor := n
	for cursor != nil {
		l = append(l, cursor)
		cursor = cursor.Parent
	}
	slices.Reverse(l)
	return l
}
