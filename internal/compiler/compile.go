package compiler

import (
	"fmt"

	"github.com/ufukty/kask/internal/compiler/builder"
	"github.com/ufukty/kask/internal/compiler/builder/directory"
)

func Compile(dst, src string, s *builder.UserSettings) error {
	dir, err := directory.Inspect(src)
	if err != nil {
		return fmt.Errorf("directory.Inspect: %w", err)
	}
	err = builder.Build(dst, dir, s)
	if err != nil {
		return fmt.Errorf("builder.Build: %w", err)
	}
	return nil
}
