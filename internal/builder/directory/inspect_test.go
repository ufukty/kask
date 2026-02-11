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
		"(Pages:", strings.Join(d.Pages, ","), ") ",
		"(Assets:", d.Assets, ") ",
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
	// Name:. (Pages:index.tmpl) (Assets:false) (Kask:false)
	//     Name:career (Pages:career/index.tmpl) (Assets:false) (Kask:false)
	//     Name:docs (Pages:docs/birdseed.md,docs/magnet.md) (Assets:false) (Kask:true)
	//     Name:products (Pages:products/index.tmpl) (Assets:false) (Kask:false)
}
