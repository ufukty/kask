package main

import (
	"fmt"
	"os"

	"go.ufukty.com/kask/cmd/kask"
)

func main() {
	err := kask.Main()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
