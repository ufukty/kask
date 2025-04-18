package build

import (
	"flag"
	"fmt"

	"github.com/ufukty/kask/internal/compiler"
	"github.com/ufukty/kask/internal/compiler/builder"
)

type args struct {
	In     string
	Out    string
	Domain string
}

var zero args

func readargs() (*args, error) {
	a := &args{}
	flag.StringVar(&a.In, "in", "", "input directory path")
	flag.StringVar(&a.Out, "out", "", "output directory path")
	flag.StringVar(&a.Domain, "domain", "", "domain that will be used to prefix each link to static assets, pages and css files")

	flag.Parse()

	if *a == zero {
		flag.PrintDefaults()
		return nil, fmt.Errorf("all arguments are set to zero values")
	}
	return a, nil
}

func Run() error {
	a, err := readargs()
	if err != nil {
		return fmt.Errorf("readargs: %w", err)
	}
	err = compiler.Compile(a.Out, a.In, &builder.UserSettings{
		Domain: a.Domain,
	})
	if err != nil {
		return fmt.Errorf("compiler.Compile: %w", err)
	}
	return nil
}
