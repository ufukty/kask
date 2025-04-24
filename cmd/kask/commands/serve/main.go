package serve

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ufukty/kask/internal/builder"
)

// run from project root
const pub = "certs/localhost.crt" // TODO: take as args
const pri = "certs/localhost.key"

type args struct {
	In   string
	Port int
	Dev  bool
}

var zero args

func readargs() (*args, error) {
	a := &args{}
	flag.StringVar(&a.In, "in", "", "input directory path")
	flag.IntVar(&a.Port, "p", 0, "port")
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

	dst, err := os.MkdirTemp(os.TempDir(), "kask-serve-*")
	if err != nil {
		return fmt.Errorf("creating temporary build folder: %w", err)
	}

	addr := fmt.Sprintf("localhost:%d", a.Port)
	err = builder.Build(dst, a.In, builder.Args{
		Domain: addr,
		Dev:    a.Dev,
	})
	if err != nil {
		return fmt.Errorf("building: %w", err)
	}

	err = http.ListenAndServeTLS(addr, pub, pri, http.FileServer(http.Dir(dst)))
	if err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}

	return nil
}
