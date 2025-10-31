package directory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Meta struct {
	Title            string `yaml:"title"`
	PreserveOrdering bool   `yaml:"preserve-ordering"`
}

func readMeta(path string) (*Meta, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer fh.Close()

	meta := &Meta{}
	err = yaml.NewDecoder(fh).Decode(meta)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return meta, nil
}
