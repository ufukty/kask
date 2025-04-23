package directory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Meta struct {
	Title       string `yaml:"title"`
	Shortname   string `yaml:"short"`
	Breadcrumbs bool   `yaml:"breadcrumbs"`
	Hidden      bool   `yaml:"hidden"`
}

func readMeta(path string) (*Meta, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reading meta file: %w", err)
	}
	defer fh.Close()

	meta := &Meta{}
	err = yaml.NewDecoder(fh).Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}
	return meta, nil
}
