package version

import (
	"fmt"
	"runtime/debug"
)

func Run() error {
	v, err := DigBuildInfo()
	if err != nil {
		return fmt.Errorf("digging build details: %w", err)
	}
	fmt.Println(v)
	return nil
}

func DigBuildInfo() (string, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", fmt.Errorf("build info is not available")
	}
	return bi.Main.Version, nil
}
