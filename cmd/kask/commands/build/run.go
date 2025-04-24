package build

import (
	"flag"
	"fmt"

	"github.com/ufukty/kask/internal/builder"
)

type args struct {
	In     string
	Out    string
	Domain string
	Dev    bool
}

var zero args

func readargs() (*args, error) {
	a := &args{}
	flag.StringVar(&a.In, "in", "", "input directory path")
	flag.StringVar(&a.Out, "out", "", "output directory path")
	flag.StringVar(&a.Domain, "domain", "", "domain that will be used to prefix each link to static assets, pages and css files")
	flag.BoolVar(&a.Dev, "dev", false, "adds unique suffixes to the bundled CSS to prevent browsers reusing cached stylesheets")

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
		return fmt.Errorf("reading args: %w", err)
	}
	err = builder.Build(a.Out, a.In, builder.Args{
		Domain: a.Domain,
		Dev:    a.Dev,
	})
	if err != nil {
		return fmt.Errorf("building: %w", err)
	}
	return nil
}
