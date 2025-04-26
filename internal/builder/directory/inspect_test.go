package directory

import (
	"fmt"
	"strings"
	"testing"
)

func spaces(depth int) string {
	s := ""
	for range depth {
		s += "    "
	}
	return s
}

func (d *Dir) String() string {
	return fmt.Sprint(
		"Name:", d.Name, " ",
		"(Html:", strings.Join(d.PagesHtml, ","), ") ",
		"(Markdown:", strings.Join(d.PagesMarkdown, ","), ") ",
		"(Assets:", d.Assets != "", ") ",
		"(Kask:", d.Kask != nil, ")",
	)
}

func printTree(d *Dir, depth int) {
	fmt.Printf("%s%s\n", spaces(depth), d)
	for _, s := range d.Subdirs {
		printTree(s, depth+1)
	}
}

func Test_Inspect(t *testing.T) {
	root, err := Inspect("testdata/acme")
	if err != nil {
		t.Fatal(fmt.Errorf("act: %w", err))
	}
	printTree(root, 0)
}
