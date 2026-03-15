package main

import (
	"fmt"
	"os"

	"go.ufukty.com/kask/internal/cmd"
)

func main() {
	err := cmd.Dispatch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
