package memory

import (
	"cmp"
	"slices"
	"strings"
)

const (
	red   = "\033[31m"
	bold  = "\033[1m"
	reset = "\033[0m"
)

func highlight(ss []string, i int) string {
	return strings.Join(slices.Concat(
		ss[:max(0, i)], []string{red + bold + cmp.Or(ss[i], "<empty>") + reset}, ss[min(i+1, len(ss)):],
	), "/")
}
