package directory

import (
	"fmt"
	"strings"
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
		"(Tmpl:", strings.Join(d.PagesTmpl, ","), ") ",
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

func ExampleInspect() {
	root, err := Inspect("testdata/acme")
	if err != nil {
		panic(fmt.Errorf("act: %w", err))
	}
	printTree(root, 0)
	// Output:
	// Name:. (Tmpl:index.tmpl) (Markdown:) (Assets:false) (Kask:false)
	//     Name:career (Tmpl:career/index.tmpl) (Markdown:) (Assets:false) (Kask:false)
	//     Name:docs (Tmpl:) (Markdown:docs/birdseed.md,docs/magnet.md) (Assets:false) (Kask:true)
	//     Name:products (Tmpl:products/index.tmpl) (Markdown:) (Assets:false) (Kask:false)

}
