package version

import (
	"fmt"
	"testing"

	"go.ufukty.com/kask/internal/version"
)

func TestDigBuildInfo(t *testing.T) {
	v, err := version.OfBuild()
	if err != nil {
		t.Errorf("act, unexpected error: %v", err)
	}
	if v == "" {
		t.Errorf("assert, unexpected empty value")
	}
	fmt.Println(v)
}
