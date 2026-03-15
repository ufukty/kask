package cmd

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"go.ufukty.com/kask/internal/cmd/build"
	"go.ufukty.com/kask/internal/cmd/version"
)

func Dispatch() error {
	commands := map[string]func() error{
		"build":   build.Run,
		"version": version.Run,
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("usage:\n\tkask [ %s ]",
			strings.Join(slices.Collect(maps.Keys(commands)), " | "),
		)
	}

	pick := os.Args[1]

	command, ok := commands[pick]
	if !ok {
		return fmt.Errorf("command not found: %s", pick)
	}

	os.Args = os.Args[1:]
	err := command()
	if err != nil {
		return fmt.Errorf("%s: %w", pick, err)
	}

	return nil
}
