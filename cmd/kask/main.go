package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"go.ufukty.com/kask/cmd/kask/commands/build"
	"go.ufukty.com/kask/cmd/kask/commands/version"
)

func Main() error {
	cmdmap := map[string]func() error{
		"build":   build.Run,
		"version": version.Run,
	}

	if len(os.Args) < 2 {
		return fmt.Errorf("usage:\n\tkask [ %s ]",
			strings.Join(slices.Collect(maps.Keys(cmdmap)), " | "),
		)
	}

	command := os.Args[1]

	runner, ok := cmdmap[command]
	if !ok {
		return fmt.Errorf("command not found: %s", command)
	}

	os.Args = os.Args[1:]
	err := runner()
	if err != nil {
		return fmt.Errorf("%s: %w", command, err)
	}

	return nil
}

func main() {
	err := Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
